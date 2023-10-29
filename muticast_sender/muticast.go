package muticast_sender

import (
	Q_cfg "Gmx5xx-tty-controller/configs"
	Q_met "Gmx5xx-tty-controller/my_serial"
	"log"
	"math"

	"fmt"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

const color_red = 0x01
const color_green = 0x02
const color_yellow = 0x03

func get_interface_ip(iface *net.Interface) (net.IP, error) {
	addrs, err := iface.Addrs()
	// handle err
	if err != nil {
		log.Printf("can't find specified interface Addrs %v\n", err)
		return nil, err
	}
	var ip net.IP
	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		// process IP address
		ip = ip.To4()
		if ip == nil {
			continue // not an ipv4 address
		}
		break
	}
	return ip, nil
}

func Sender(met_data *Q_met.GMX_data_struct, ifaces *Q_cfg.Muti_Cfg) {
	var conn_new = make([]*net.UDPConn, 10)
	var pc_new = make([]*ipv4.PacketConn, 10)
	log.Printf("setup interfaces: %v\n", ifaces.Ifaces)
	for i, iface_name := range ifaces.Ifaces {
		iface, err := net.InterfaceByName(iface_name)
		if err != nil {
			log.Printf("can't find specified interface %v\n", err)
			return
		}
		ip, err := get_interface_ip(iface)
		ipv4Addr := &net.UDPAddr{IP: ip, Port: 0}
		raddr := net.UDPAddr{IP: net.ParseIP("232.140.115.65"), Port: 50013}
		connecter, err := net.DialUDP("udp", ipv4Addr, &raddr)
		if err != nil {
			log.Printf("ListenUDP ipv4Addr error %v\n", err)
			return
		}

		pc := ipv4.NewPacketConn(connecter)
		if loop, err := pc.MulticastLoopback(); err == nil {
			log.Printf("%s MulticastLoopback status:%v\n", iface_name, loop)
			if !loop {
				if err := pc.SetMulticastLoopback(true); err != nil {
					log.Printf("SetMulticastLoopback error:%v\n", err)
				}
			}
		}
		conn_new[i] = connecter
		pc_new[i] = pc
	}

	te_gw_id := []byte("te-se-000002\x00       ")
	any_dispaly_header := []byte("\xdf\x00$001,")
	any_dispaly_end := []byte("#")

	WindU := make([]float64, 60)
	WindV := make([]float64, 60)
	for i := 0; i < 60; i++ {
		WindU[i] = 100
		WindV[i] = 100
	}
	WindPointer := 0
	WindSpeed := 0.0
	WindDirection := 0.0
	var windU, windV float64
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		stopChan := make(chan bool)
		for {
			select {
			case <-ticker.C:

				if met_data.Direction != nil && met_data.Speed != nil {
					Degree := (-float64(*met_data.Direction) + 270) * (math.Pi / 180)
					WindU[WindPointer] = math.Cos(Degree) * float64(*met_data.Speed)
					WindV[WindPointer] = math.Sin(Degree) * float64(*met_data.Speed)
				} else {
					WindU[WindPointer] = 100
					WindV[WindPointer] = 100
				}

				if WindPointer <= 58 {
					WindPointer = WindPointer + 1
				} else {
					WindPointer = 0
				}
				windU = 0
				windV = 0
				counter := 0
				for i := 0; i < 60; i++ {
					if WindU[i] < 100 {
						windU = windU + WindU[i]
						windV = windV + WindV[i]
						counter = counter + 1
					}
				}
				windU = windU / float64(counter)
				windV = windV / float64(counter)
				WindSpeed = math.Sqrt(math.Pow(windU, 2) + math.Pow(windV, 2))
				WindDirection = math.Atan2(windU, windV)*(180/math.Pi) + 180 // radians * (180/pi)
				// log.Printf("WS: %f WD: %f", WindSpeed, WindDirection)

				any_display_messsage := append(te_gw_id, any_dispaly_header...)
				any_display_messsage = append(any_display_messsage, color_green)
				if met_data.Temperature != nil {
					any_display_messsage = append(any_display_messsage, []byte(fmt.Sprintf("%3.0f", math.Round(float64(*met_data.Temperature*10))))...)
				} else {
					any_display_messsage = append(any_display_messsage, []byte("---")...)
				}
				any_display_messsage = append(any_display_messsage, []byte("  ")...)
				any_display_messsage = append(any_display_messsage, color_green)
				if met_data.Speed != nil {
					// any_display_messsage = append(any_display_messsage, []byte(fmt.Sprintf("%3.0f", math.Round(float64(*met_data.Speed*10))))...)
					any_display_messsage = append(any_display_messsage, []byte(fmt.Sprintf("%3.0f", math.Round(float64(WindSpeed*10))))...)
				} else {
					any_display_messsage = append(any_display_messsage, []byte("---")...)
				}
				any_display_messsage = append(any_display_messsage, color_green)
				if met_data.Relative_Humidity != nil {
					any_display_messsage = append(any_display_messsage, []byte(fmt.Sprintf("%3.0f", math.Round(float64(*met_data.Relative_Humidity))))...)
				} else {
					any_display_messsage = append(any_display_messsage, []byte("---")...)
				}
				any_display_messsage = append(any_display_messsage, []byte("  ")...)
				any_display_messsage = append(any_display_messsage, color_green)
				if met_data.Direction != nil {
					// any_display_messsage = append(any_display_messsage, []byte(fmt.Sprintf("%3.0f", *met_data.Direction))...)
					any_display_messsage = append(any_display_messsage, []byte(fmt.Sprintf("%3.0f", WindDirection))...)
				} else {
					any_display_messsage = append(any_display_messsage, []byte("---")...)
				}
				any_display_messsage = append(any_display_messsage, any_dispaly_end...)
				any_display_messsage[21] = byte(len(any_display_messsage) - 22)
				for _, iconn := range conn_new {
					if iconn != nil {
						if _, err := iconn.Write(any_display_messsage); err != nil {
							log.Printf("Write failed, %v\n", err)
						}
					}

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
