package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	Q_cfg "Gmx5xx-tty-controller/configs"
	Q_muti "Gmx5xx-tty-controller/muticast_sender"
	Q_met "Gmx5xx-tty-controller/my_serial"

	"go.bug.st/serial"
)

const serverPort = 3000

var Debug = false

func float32_to_float64(value *float32) *float64 {
	if value != nil {
		converted := float64(*value)
		return &converted
	} else {
		return nil
	}
}

func int_to_int(value *int) *int {
	if value != nil {
		converted := int(*value)
		return &converted
	} else {
		return nil
	}
}

func pub_met_data(data Q_cfg.Data, cfg *Q_cfg.Qlog_Cfg) {
	bodyData, err := json.Marshal(data)
	if Debug {
		log.Printf("Data: %s\n", bodyData)
	}
	bodyReader := bytes.NewReader(bodyData)

	requestURL := fmt.Sprintf("http://%s:%d/api/v1/gmxdata", cfg.MetSvCfg.Hostname, serverPort)

	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)

	defer req.Body.Close()

	if err != nil {
		log.Printf("client: could not create request: %s\n", err)
		return
		// os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("client: error making http request: %s\n", err)
		return
		// os.Exit(1)
	}

	if res.StatusCode != 200 {
		log.Printf("client: status code: %d\n", res.StatusCode)
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("client: could not read response body: %s\n", err)
			res.Body.Close()
			return
			// os.Exit(1)
		}
		log.Printf("client: response body: %s\n", resBody)
	}
	res.Body.Close()
}

func restfull_go(met_data *Q_met.GMX_data_struct, cfg *Q_cfg.Qlog_Cfg) {

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		stopChan := make(chan bool)
		for {
			select {
			case <-ticker.C:
				var record_time = time.Now().UTC()
				// log.Println(record_time)
				Upload_data := &Q_cfg.Data{
					Temperature:            float32_to_float64(met_data.Temperature),
					RelativeHumitiy:        float32_to_float64(met_data.Relative_Humidity),
					Dewpoint:               float32_to_float64(met_data.Dewpoint),
					Pressure:               float32_to_float64(met_data.Pressure),
					WindDirection:          float32_to_float64(met_data.Direction),
					WindSpeed:              float32_to_float64(met_data.Speed),
					WindCorrectedDirection: float32_to_float64(met_data.Corrected_Direction),
					GmxSupplyVoltage:       float32_to_float64(met_data.Supply_Voltage),
					GmxStatus:              int_to_int(met_data.Status),
					Timestamp:              record_time,
				}
				go pub_met_data(*Upload_data, cfg)
			case stop := <-stopChan:
				if stop {
					fmt.Println("Ticker Stop! Channel must be closed")
					return
				}
			}
		}
	}()
}

func main() {
	log.Println("GMX5xx tty controller service start")
	cfg, err := Q_cfg.Read_Qlog_Cfg("./configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	log.Printf("Hostname: %s\n", hostname)

	// log.Printf("%#v\n", cfg)

	argsWithProg := os.Args

	if len(argsWithProg) == 2 {
		if argsWithProg[1] == "backup" {
			days, err := strconv.Atoi(argsWithProg[2])
			if err != nil {
				log.Println("Error during conversion")
				return
			}
			log.Printf("met-logger retry upload start for %d days", days)
		} else if argsWithProg[1] == "debug" {
			log.Println("debug mode open")
			Debug = true
		}
	}
	Met_data_st := new(Q_met.GMX_data_struct)

	Q_met.InitMetData(Met_data_st)

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Println("No Serial ports found!")
	}
	for _, port := range ports {
		log.Printf("Found port: %v\n", port)
	}
	go Q_met.Met_go(Met_data_st, &cfg.Meteo_Cfg)

	time.Sleep(1000 * time.Millisecond)

	go Q_muti.Sender(Met_data_st, &cfg.Muticast_iface)

	restfull_go(Met_data_st, cfg)

	select {}
}
