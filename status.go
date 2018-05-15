package main


import (
	"fmt"
	"encoding/json"
	"time"
	"os/exec"
	"log"
	"bytes"
	"flag"
	"math/rand"
	MQTT "github.com/eclipse/paho.mqtt.golang"
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
}

type DataDevice struct {
	CurrentValue	int	`json:"currentValue"`
	TotalValue	int	`json:"totalValue"`
	Type		string	`json:"type"`
	UnitInformation string	`json:"unitInformation"`
}

type Status struct {
	Gstatus		map[string]interface{}  `json:"gstatus"`
	SnDev		string  `json:"sn-dev"`
	SnModem		string	`json:"sn-modem"`
	SnDisplay	string	`json:"sn-display"`
	TimeStamp	float64	`json:"timeStamp"`
	IpMaskMap	map[string][]string	`json:"ipMaskMap"`
	Hostname	string	`json:"hostname"`
	AppVers		map[string]string	`json:"AppVers"`
	SimStatus	string	`json:"simStatus"`
	SimImei		string	`json:"simImei"`
	UsosTranspCount	int	`json:"usosTranspCount"`
	ErrsTranspCount	int	`json:"errsTranspCount"`
	Volt		map[string]float64	`json:"volt"`
	Mac		string	`json:"mac"`
	AppTablesVers	map[string]string	`json:"AppTablesVers"`
	CpuStatus	[]float64	`json:"cpuStatus"`
	Dns		[]string	`json:"dns"`
	UpTime		string	`json:"upTime"`
	DeviceDataList	[]*DataDevice	`json:"deviceDataList"`
	Gateway		string	`json:"gateway"`
	TypeDev		string `json:"type-dev"` 
}


var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
  fmt.Printf("TOPIC: %s\n", msg.Topic())
  fmt.Printf("MSG: %s\n", msg.Payload())
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


	b2, err := json.Marshal(m)
	//fmt.Printf("salida: %s\n\n\n", b2)

	uptime := getUptime()
	usos, errores := usosTransp()

	m.DeviceDataList[0].CurrentValue = 350000
	m.DeviceDataList[1].CurrentValue = 800
	m.UpTime = uptime
	m.UsosTranspCount = usos
	m.ErrsTranspCount = errores
	m.TimeStamp = float64(time.Now().UnixNano())/1000000000
	m.Volt["currentValue"] = 12
	m.Volt["highestValue"] = 12.1
	m.Volt["lowestValue"] = 11

	m.SnDev = snDev
	m.TypeDev = typeDev
	m.IpMaskMap["eth0"][0] = ipEth0
	m.IpMaskMap["eth1"][0] = ipEth1
	m.Hostname = hostname
	m.Mac = mac
	m.SimImei =  simImei
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



	timeoutSend := time.Tick( time.Duration(timeout) * time.Second)


	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID("go-simple")
  	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
  	c := MQTT.NewClient(opts)
  	if token := c.Connect(); token.Wait() && token.Error() != nil {
   		panic(token.Error())
  	}

	defer c.Disconnect(250)



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
			go publish(&m, c)
		}
	}
}

func publish(m *Status, c MQTT.Client) {

	m.TimeStamp = float64(time.Now().UnixNano())/1000000000
	
	uptime := getUptime()
	usos, errores := usosTransp()

	m.UpTime = uptime
	m.UsosTranspCount = usos
	m.ErrsTranspCount = errores
	
	b3, err := json.Marshal(m)

	if err != nil {
		log.Println(err)
		return
	}

	token := c.Publish("STATUS/state", 0, false, b3)
	token.Wait()
	log.Printf("message: %s\n",b3)
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
			value = float64(rand.Intn(40))
			value5 = value5 + value
			acc5 = acc5 + 1
			value5 = value5 + value
			acc15 = acc15 + 1
			ch <- typeCpu{name: "one", value: value}
		case <-t2:
			if acc5 != 0 {
				ch <- typeCpu{name: "five", value: value5/float64(acc5)}
				value5 = 0
				acc5 = 0
			}
		case <-t3:
			if acc15 != 0 {
				ch <- typeCpu{name: "fifteen", value: value15/float64(acc15)}
				value15 = 0
				acc15 = 0
			}
		case <-talarm:
			value = float64(rand.Intn(120)) + 50
			value5 = value5 + value
			acc5 = acc5 + 1
			value5 = value5 + value
			acc15 = acc15 + 1
			ch <- typeCpu{name: "one", value: value}
		}
	}
}

/**/

type typeCpu struct {
	name	string
	value	float64
}
func randVolt(ch chan map[string]float64) {
	t1 := time.Tick(60 * time.Second)
	volt := make(map[string]float64)
	for {
		select {
		case <-t1:
			volt["current"] = float64(rand.Intn(200))/100 + 11
			volt["high"] = float64(rand.Intn(400))/100 + 12
			volt["low"] = -float64(rand.Intn(500))/100 + 12
			ch <- volt
		}
	}
}

func usosTransp() (usos int, errores int) {
	usos = rand.Intn(300) + 50
	errores = rand.Intn(30)
	return
}

func randMemory(ch chan int) {
	t1 := time.Tick(300 * time.Second)
	t2 := time.Tick(902 * time.Second)
	for {
		select {
		case <- t1:
			ch <- rand.Intn(2000) + 300000
		case <- t2:
			ch <- rand.Intn(650000) + 300000
		}
	}
}

func randTemp(ch chan int) {
	t1 := time.Tick(60 * time.Second)
	t2 := time.Tick(1800 * time.Second)
	for {
		select {
		case <- t1:
			ch <- rand.Intn(5) + 32
		case <- t2:
			ch <- rand.Intn(12) + 35
		}
	}
}

func randSd(ch chan int) {
	t1 := time.Tick(600 * time.Second)
	t2 := time.Tick(1800 * time.Second)
	for {
		select {
		case <- t1:
			ch <- rand.Intn(500) + 500
		case <- t2:
			ch <- rand.Intn(6000) + 1000
		}
	}
}
/**/

