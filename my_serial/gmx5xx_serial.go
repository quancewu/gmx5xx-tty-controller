package my_serial

import (
	Q_cfg "Gmx5xx-tty-controller/configs"
	crc "Gmx5xx-tty-controller/samples/crc"
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"sync"
	"time"

	"go.bug.st/serial"
)

var Met_lock sync.RWMutex

var Met_data = make([]float32, 10)

var GMX_data_st = make(chan GMX_data_struct)

type GMX_data_struct struct {
	Direction               *float32
	Speed                   *float32
	Corrected_Direction     *float32
	Pressure                *float32
	Relative_Humidity       *float32
	Temperature             *float32
	Dewpoint                *float32
	Total_Precipitation     *float32
	Precipitation_Intensity *float32
	Supply_Voltage          *float32
	Status                  *int
}

func Met_go(met_data *GMX_data_struct, meteo_cfg *Q_cfg.Meteo_Cfg) {

	mode := &serial.Mode{
		BaudRate: meteo_cfg.MetSerialCfg.BaudRate,
		Parity:   serial.Parity(meteo_cfg.MetSerialCfg.Parity),
		DataBits: meteo_cfg.MetSerialCfg.DataBits,
		StopBits: serial.StopBits(meteo_cfg.MetSerialCfg.StopBits - 1),
	}

	model_para_len := 0

	port2, err := serial.Open(meteo_cfg.MetSerialCfg.Port, mode)

	if err != nil {
		log.Fatal(err)
	}
	port2.SetReadTimeout(time.Duration(meteo_cfg.MetSerialCfg.Timeout) * time.Millisecond)
	modbus_poll := make([]byte, 6)
	if meteo_cfg.Model == "GMX500" {
		modbus_poll = []byte{1, 3, 0, 0, 0, 0x24}
		model_para_len = 6
	} else if meteo_cfg.Model == "GMX550" {
		modbus_poll = []byte{1, 3, 0, 0, 0, 0x28}
		model_para_len = 8
	} else {
		modbus_poll = []byte{1, 3, 0, 0, 0, 0x28}
		model_para_len = 8
	}
	checksum := crc.CheckSum(modbus_poll)

	int16buf := new(bytes.Buffer)

	binary.Write(int16buf, binary.LittleEndian, checksum)

	modbus_poll = append(modbus_poll, int16buf.Bytes()...)

	buff := make([]byte, modbus_poll[5]*2+5)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		stopChan := make(chan bool)
		for {
			select {
			case <-ticker.C:
				port2.ResetInputBuffer()
				_, err := port2.Write(modbus_poll[:])
				if err != nil {
					log.Fatal(err)
				}
				time.Sleep(150 * time.Millisecond)
				n1, err := port2.Read(buff)
				if err != nil {
					log.Fatal(err)
					break
				}
				if n1 == 0 {
					// fmt.Println("\nEOF")
					Met_lock.Lock()
					OffSetMetData(met_data)
					Met_lock.Unlock()
					log.Printf("%s Modbus Timeout", meteo_cfg.Model)
					continue
				}
				dataCRC := crc.CheckSum(buff[:n1-2])
				recive16buf := new(bytes.Buffer)
				crc16buf := new(bytes.Buffer)

				binary.Write(crc16buf, binary.LittleEndian, dataCRC)
				binary.Write(recive16buf, binary.BigEndian, buff[n1-2:n1])
				if bytes.Compare(crc16buf.Bytes(), recive16buf.Bytes()) == 0 {
					// fmt.Printf("checked]\n")
					Met_lock.Lock()
					for i := 0; i <= model_para_len; i = i + 1 {
						bit := binary.BigEndian.Uint32(buff[7+i*4 : 11+i*4])
						Met_data[i] = math.Float32frombits(bit)
					}
					bit := binary.BigEndian.Uint32(buff[7+(model_para_len+9)*4 : 11+(model_para_len+9)*4])
					Met_data[9] = math.Float32frombits(bit)
					bitstatus := binary.BigEndian.Uint32(buff[7+(model_para_len+10)*4 : 11+(model_para_len+10)*4])
					Met_status := int(bitstatus)
					met_data.Direction = &Met_data[0]
					met_data.Speed = &Met_data[1]
					met_data.Corrected_Direction = &Met_data[2]
					met_data.Pressure = &Met_data[3]
					met_data.Relative_Humidity = &Met_data[4]
					met_data.Temperature = &Met_data[5]
					met_data.Dewpoint = &Met_data[6]
					met_data.Total_Precipitation = nil
					met_data.Precipitation_Intensity = nil
					met_data.Supply_Voltage = &Met_data[9]
					met_data.Status = &Met_status
					Met_lock.Unlock()
					// log.Printf("%+v\n", met_data)
					// log.Printf("Dewpoint %+v\n", *met_data.Dewpoint)
				} else {
					Met_lock.Lock()
					OffSetMetData(met_data)
					Met_lock.Unlock()
					log.Printf("%s Modbus CRC error", meteo_cfg.Model)
				}
			case stop := <-stopChan:
				if stop {
					log.Println("Ticker Stop! Channel closed")
					return
				}
			}
		}

	}()
}

func OffSetMetData(met_data *GMX_data_struct) {
	met_data.Direction = nil
	met_data.Speed = nil
	met_data.Corrected_Direction = nil
	met_data.Pressure = nil
	met_data.Relative_Humidity = nil
	met_data.Temperature = nil
	met_data.Dewpoint = nil
	met_data.Total_Precipitation = nil
	met_data.Precipitation_Intensity = nil
	met_data.Supply_Voltage = nil
	met_data.Status = nil
}

func InitMetData(met_data *GMX_data_struct) {
	met_data.Direction = nil
	met_data.Speed = nil
	met_data.Corrected_Direction = nil
	met_data.Pressure = nil
	met_data.Relative_Humidity = nil
	met_data.Temperature = nil
	met_data.Dewpoint = nil
	met_data.Total_Precipitation = nil
	met_data.Precipitation_Intensity = nil
	met_data.Supply_Voltage = nil
	met_data.Status = nil
}
