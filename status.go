package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/dumacp/gpsnmea"
	"github.com/dumacp/pubsub"
)

var timeout int
var mac string
var ipEth0 string
var ipEth1 string
var snDev string
var typeDev string
var hostname string
var simImei string
var simStatus string

//var instance string

func init() {
	flag.IntVar(&timeout, "timeout", 30, "timeout to capture send status")
	flag.StringVar(&mac, "mac", "01:02:03:04:05:06", "device's MAC address (default: 30)")
	flag.StringVar(&ipEth0, "ipEth0", "192.168.188.23/24", "device's ip address")
	flag.StringVar(&ipEth1, "ipEth1", "172.23.99.1/29", "device's ip address")
	flag.StringVar(&snDev, "snDev", "sn0001-TEST-0001", "device's serial ID")
	flag.StringVar(&typeDev, "typeDev", "OMVZ7", "device type")
	flag.StringVar(&hostname, "hostname", "OMVZ7", "device's hostname")
	flag.StringVar(&simImei, "simImei", "123456789ABCD", "modem's SIM IMEI")
	flag.StringVar(&simStatus, "simStatus", "OK", "SIM Status")
	//flag.StringVar(&instance, "instance", "", "instance / multi devices")
}

type DataDevice struct {
	CurrentValue    int    `json:"currentValue"`
	TotalValue      int    `json:"totalValue"`
	Type            string `json:"type"`
	UnitInformation string `json:"unitInformation"`
}

type Status struct {
	Gstatus                     map[string]interface{} `json:"gstatus"`
	SnDev                       string                 `json:"sn-dev"`
	SnModem                     string                 `json:"sn-modem"`
	SnDisplay                   string                 `json:"sn-display"`
	TimeStamp                   float64                `json:"timeStamp"`
	IpMaskMap                   map[string][]string    `json:"ipMaskMap"`
	Hostname                    string                 `json:"hostname"`
	AppVers                     map[string]string      `json:"AppVers"`
	SimStatus                   string                 `json:"simStatus"`
	SimImei                     string                 `json:"simImei"`
	UsosTranspCount             int                    `json:"usosTranspCount"`
	ErrsTranspCount             int                    `json:"errsTranspCount"`
	Volt                        map[string]float64     `json:"volt"`
	Mac                         string                 `json:"mac"`
	AppTablesVers               map[string]string      `json:"AppTablesVers"`
	CpuStatus                   []float64              `json:"cpuStatus"`
	Dns                         []string               `json:"dns"`
	UpTime                      string                 `json:"upTime"`
	DeviceDataList              []*DataDevice          `json:"deviceDataList"`
	Gateway                     string                 `json:"gateway"`
	TypeDev                     string                 `json:"type-dev"`
	TurnstileUpAccum            int                    `json:"turnstileUpAccum"`
	TurnstileDownAccum          int                    `json:"turnstileDownAccum"`
	FrontDoorPassengerUpAccum   int                    `json:"frontDoorPassengerUpAccum"`
	FrontDoorPassengerDownAccum int                    `json:"frontDoorPassengerDownAccum"`
	BackDoorPassengerUpAccum    int                    `json:"backDoorPassengerUpAccum"`
	BackDoorPassengerDownAccum  int                    `json:"backDoorPassengerDownAccum"`
}

func main() {

	flag.Parse()

	b1 := []byte(`{"gstatus": {"rsrq": null, "temperature": 37, "sinr": null, "band": "WCDMA", "mode": "WCDMA", "tac": null, "cellid": 19692771}, "sn-dev": "sn0001-0001-TEST", "sn-modem": "sn0001-0001-TEST-MODEM", "sn-display": "sn0001-0001-DISP", "timeStamp": 1526041174.061357, "ipMaskMap": {"lo": ["127.0.0.1/8", "::1/128"], "eth1": ["172.28.99.47/27", "fe80::d82a:bdff:fe31:e408/64"], "eth0": ["192.168.188.197/24", "fe80::201:2ff:fe03:405/64"]}, "hostname": "OMVZ7", "AppVers": {"libparamoperacioncliente": "17.10.27.1", "libcontrolregistros": "17.10.31.1", "embedded.libgestionhardware": "18.02.13.3", "libcontrolconsecutivos": "13.07.19.1", "AppUsosTrasnporte": "18.02.07.1", "libcommonentities": "17.12.21.1", "libcontrolmensajeria": "17.10.31.1", "libgestionhardware": "15.05.22.01"}, "simStatus": "OK", "simImei": "359072061300642", "usosTranspCount": 0, "errsTranspCount": 0, "volt": {"currentValue": 12.096, "highestValue": 12.256, "lowestValue": 0.44800000000000001}, "mac": "00:01:02:03:04:05", "AppTablesVers": {"TablaTrayectos": "139", "ListaNegra": "5378.07"}, "cpuStatus": [7.0, 7.0, 4.0], "dns": ["192.168.188.188"], "upTime": "07:19:32 up 32 days, 15 min,  load average: 0.07, 0.07, 0.04", "deviceDataList": [{"currentValue": 319240, "totalValue": 1026448, "type": "MEM", "unitInformation": "KiB"}, {"currentValue": 815, "totalValue": 7446, "type": "SD", "unitInformation": "MiB"}], "gateway": "192.168.188.188", "type-dev": "OMVZ7"}`)
	var m Status
	err := json.Unmarshal(b1, &m)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("%v\n", m)
	//fmt.Printf("cpuStatus: %v\n", m.CpuStatus)

	//b2, err := json.Marshal(m)
	//fmt.Printf("salida: %s\n\n\n", b2)

	uptime := getUptime()
	usos, errores := usosTransp()

	m.DeviceDataList[0].CurrentValue = 350000
	m.DeviceDataList[1].CurrentValue = 800
	m.UpTime = uptime
	m.UsosTranspCount = usos
	m.ErrsTranspCount = errores
	m.TimeStamp = float64(time.Now().UnixNano()) / 1000000000
	m.Volt["currentValue"] = 12
	m.Volt["highestValue"] = 12.1
	m.Volt["lowestValue"] = 11

	m.SnDev = snDev
	m.TypeDev = typeDev
	m.IpMaskMap["eth0"][0] = ipEth0
	m.IpMaskMap["eth1"][0] = ipEth1
	m.Hostname = hostname
	m.Mac = mac
	m.SimImei = simImei
	m.SimStatus = simStatus
	m.Gstatus["temperature"] = int(34)

	chVolt := make(chan map[string]float64)
	go randVolt(chVolt)
	chSd := make(chan int)
	go randSd(chSd)
	chMem := make(chan int)
	go randMemory(chMem)
	chTemp := make(chan int)
	go randTemp(chTemp)
	chCpu := make(chan typeCpu)
	go randCpu(chCpu)

	timeoutSend := time.Tick(time.Duration(timeout) * time.Second)

	//create and start a client using the above ClientOptions
	pub, err := pubsub.NewConnection(fmt.Sprintf("go-status-%v", time.Now().UnixNano()))
	if err != nil {
		log.Fatal(err)
	}
	defer pub.Disconnect()
	msgChan := make(chan string)
	go pub.Publish("STATUS/state", msgChan)
	go func() {
		for v := range pub.Err {
			log.Println(v)
		}
	}()

	evtChan := make(chan string)
	go pub.Publish("EVENTS/event", evtChan)
	go func() {
		for v := range pub.Err {
			log.Println(v)
		}
	}()

	go events(evtChan)

	for {
		rand.Seed(time.Now().UnixNano())
		select {
		case volt := <-chVolt:
			m.Volt["currentValue"] = volt["current"]
			m.Volt["highestValue"] = volt["high"]
			m.Volt["lowestValue"] = volt["low"]
		case sd := <-chSd:
			m.DeviceDataList[1].CurrentValue = sd
		case mem := <-chMem:
			m.DeviceDataList[0].CurrentValue = mem
		case temp := <-chTemp:
			m.Gstatus["temperature"] = temp
		case cpu := <-chCpu:
			switch name := cpu.name; name {
			case "one":
				m.CpuStatus[0] = cpu.value
			case "five":
				m.CpuStatus[1] = cpu.value
			case "fifteen":
				m.CpuStatus[2] = cpu.value
			}
		case <-timeoutSend:
			time.Sleep(time.Second * 1)
			if msg, err := prepare(&m); err != nil {
				log.Println(err)
			} else {
				msgChan <- string(msg)
			}
		}
	}
}

func prepare(m *Status) ([]byte, error) {

	m.TimeStamp = float64(time.Now().UnixNano()) / 1000000000
	uptime := getUptime()
	//usos, errores := usosTransp()
	counters := contadores()
	m.TurnstileUpAccum += counters[0]
	m.TurnstileDownAccum += counters[1]
	m.FrontDoorPassengerUpAccum += counters[2]
	m.FrontDoorPassengerDownAccum += counters[3]
	m.BackDoorPassengerUpAccum += counters[4]
	m.BackDoorPassengerDownAccum += counters[5]

	m.UpTime = uptime
	m.UsosTranspCount = 0
	m.ErrsTranspCount = 0

	b3, err := json.Marshal(m)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("message: %s\n", b3)
	return b3, nil
}

func getUptime() (uptime string) {

	cmd := exec.Command("uptime")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	uptime = out.String()
	return
}

func randCpu(ch chan typeCpu) {
	t1 := time.Tick(60 * time.Second)
	t2 := time.Tick(300 * time.Second)
	t3 := time.Tick(900 * time.Second)
	talarm := time.Tick(600 * time.Second)

	value := 0.0
	value5 := 0.0
	acc5 := 0
	value15 := 0.0
	acc15 := 0

	for {
		select {
		case <-t1:
			value = float64(rand.Intn(100))
			value5 = value5 + value
			acc5 = acc5 + 1
			value15 = value15 + value
			acc15 = acc15 + 1
			ch <- typeCpu{name: "one", value: value}
		case <-t2:
			if acc5 != 0 {
				ch <- typeCpu{name: "five", value: value5 / float64(acc5)}
				value5 = 0
				acc5 = 0
			}
		case <-t3:
			if acc15 != 0 {
				ch <- typeCpu{name: "fifteen", value: value15 / float64(acc15)}
				value15 = 0
				acc15 = 0
			}
		case <-talarm:
			value = float64(rand.Intn(20)) + 80
			value5 = value5 + value
			acc5 = acc5 + 1
			value15 = value15 + value
			acc15 = acc15 + 1
			ch <- typeCpu{name: "one", value: value}
		}
	}
}

/**/

type typeCpu struct {
	name  string
	value float64
}

func randVolt(ch chan map[string]float64) {
	t1 := time.Tick(60 * time.Second)
	t2 := time.Tick(600 * time.Second)
	volt := make(map[string]float64)
	for {
		select {
		case <-t1:
			volt["current"] = float64(rand.Intn(200))/100 + 11
			volt["high"] = float64(rand.Intn(400))/100 + 12
			volt["low"] = -float64(rand.Intn(500))/100 + 12
			ch <- volt
		case <-t2:
			volt["current"] = float64(rand.Intn(200))/100 + 11
			volt["high"] = float64(rand.Intn(400))/100 + 12
			volt["low"] = -float64(rand.Intn(800))/100 + 3
			ch <- volt
		}
	}
}

func usosTransp() (usos int, errores int) {
	usos = rand.Intn(10)
	errores = 0
	if usos > 5 {
		errores = rand.Intn(3)
	}
	return
}

func contadores() []int {
	turnstileUpCount := rand.Intn(10)
	puertaDelanteraIngresos := turnstileUpCount
	turnstileDownCount := 0
	puertaTraseraSalidas := 0
	puertaDelanteraSalidas := 0
	puertaTraseraIngresos := 0
	if puertaDelanteraIngresos > 3 {
		puertaTraseraSalidas = rand.Intn(5)
	}
	if puertaTraseraSalidas > 3 {
		puertaDelanteraSalidas = rand.Intn(3)
	}
	return []int{turnstileUpCount, turnstileDownCount, puertaDelanteraIngresos, puertaDelanteraSalidas, puertaTraseraIngresos, puertaTraseraSalidas}
}

func randMemory(ch chan int) {
	t1 := time.Tick(60 * time.Second)
	t2 := time.Tick(300 * time.Second)
	t3 := time.Tick(902 * time.Second)
	temp := 350000
	for {
		select {
		case <-t1:
			ch <- rand.Intn(10000) + temp
		case <-t2:
			temp = rand.Intn(80000) + 300000
		case <-t3:
			temp = rand.Intn(650000) + 300000
		}
	}
}

func randTemp(ch chan int) {
	t1 := time.Tick(60 * time.Second)
	t2 := time.Tick(1800 * time.Second)
	for {
		select {
		case <-t1:
			ch <- rand.Intn(5) + 32
		case <-t2:
			ch <- rand.Intn(12) + 38
		}
	}
}

func randSd(ch chan int) {
	t1 := time.Tick(600 * time.Second)
	t2 := time.Tick(1800 * time.Second)
	for {
		select {
		case <-t1:
			ch <- rand.Intn(500) + 500
		case <-t2:
			ch <- rand.Intn(6000) + 1000
		}
	}
}

func events(ch chan string) {
	if len(ruta) <= 0 {
		return
	}
	var1 := make(map[string][][][]float64)
	if err := json.Unmarshal(ruta, &var1); err != nil {
		return
	}

	values, ok := var1["ruta"]
	if !ok {
		return
	}

	itirenario := values[0]
	itirenario = append(itirenario, values[1]...)

	chPoints := make(chan []float64, 0)

	go func() {
		for {
			for _, v := range itirenario {
				chPoints <- v
			}
		}
	}()

	t1 := time.Tick(30 * time.Second)
	t2 := time.Tick(60 * time.Second)
	for {
		select {
		case <-t1:
			tn := time.Now()
			msg := &pubsub.Message{
				Timestamp: float64(tn.UnixNano()) / 1000000000,
				Type:      "GPRMC",
			}
			var lat float64
			var lon float64
			daten := tn.Format("020106")
			timen := tn.Format("150405")
			select {
			case v := <-chPoints:
				lon = v[0]
				lat = v[1]
			}

			frame := fmt.Sprintf("$GPRMC,%v.0,A,%v,%v,%3.1f,0.0,%v,4.7,W,A*",
				timen, gpsnmea.DecimalDegreeToLat(lat), gpsnmea.DecimalDegreeToLon(lon), rand.Float64()*20, daten)

			// frame := "$GPRMC,164016.0,A,0615.179728,N,07535.343742,W,0.0,0.0,100518,4.7,W,A*"

			checksum := byte(0)
			frameB := []byte(frame)
			for _, v := range frameB[1 : len(frameB)-1] {
				checksum = checksum ^ v
			}
			frame = fmt.Sprintf("%v%02X", frame, checksum)

			msg.Value = frame
			v, err := json.Marshal(msg)
			if err != nil {
				break
			}
			fmt.Printf("%s\n", v)
			ch <- string(v)

		case <-t2:
			tn := time.Now()
			msg := &pubsub.Message{
				Timestamp: float64(tn.UnixNano()) / 1000000000,
				Type:      "TURNSTILE",
			}
			var lat float64
			var lon float64
			daten := tn.Format("020106")
			timen := tn.Format("150405")
			select {
			case v := <-chPoints:
				lon = v[0]
				lat = v[1]
			}
			frame := fmt.Sprintf("$GPRMC,%v.0,A,%v,%v,%3.1f,0.0,%v,4.7,W,A*",
				timen, gpsnmea.DecimalDegreeToLat(lat), gpsnmea.DecimalDegreeToLon(lon), rand.Float64()*20, daten)
			// frame := "$GPRMC,164016.0,A,0615.179728,N,07535.343742,W,0.0,0.0,100518,4.7,W,A*"

			checksum := byte(0)
			frameB := []byte(frame)
			for _, v := range frameB[1 : len(frameB)-1] {
				checksum = checksum ^ v
			}
			frame = fmt.Sprintf("%v%02X", frame, checksum)

			val := struct {
				Coord              string `json:"coord"`
				TurnstileUpCount   int    `json:"turnstileUpCount"`
				TurnstileDownCount int    `json:"turnstileDownCount"`
			}{
				frame,
				rand.Intn(15),
				rand.Intn(10),
			}

			msg.Value = val

			v, err := json.Marshal(msg)
			if err != nil {
				break
			}
			fmt.Printf("%s\n", v)
			ch <- string(v)
		}
	}
}

/**/
