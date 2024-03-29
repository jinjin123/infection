package transfer

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, Data[0])
}

func PacHandle(PacPort string) {
	switch runtime.GOOS {
	case "darwin":
	case "windows":
		registerPAC()
	case "linux":
	}
	http.HandleFunc("/pac", handler)
	if err := http.ListenAndServe(PacPort, nil); err != nil {
		log.Println("pac Faild", err)
	}
}

func ResgisterAutoRun() {
	path, _ := GetTargetPath()
	log.Println("可执行文件:", path)
	//register table
	key, exists, err := registry.CreateKey(registry.LOCAL_MACHINE, "Software\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	if err != nil {
		log.Println("set auto register table Faild:", err)
		log.Println("请以管理员权限运行一次（仅需运行一次）以开机自启，已执行请忽略。")
		return
	}
	defer key.Close()

	if exists {
		//log.Println("key exits")
	} else {
		//log.Println("table not exits，created",)
	}
	// 写入字符串
	err_write := key.SetStringValue("proxygo", path)
	if err_write != nil {
		log.Println("set auto register table Faild", err_write)
	} else {
		log.Println("set auto register table Success path", path)
	}
}

func registerPAC() {
	//设置注册表
	key, exists, err := registry.CreateKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", registry.ALL_ACCESS)
	if err != nil {
		log.Println(err)
	}
	defer key.Close()

	if exists {
		//log.Println("键已存在")
	} else {
		log.Println("此注册表项不存在，已自动新建")
	}
	// 写入字符串
	key.SetStringValue("AutoConfigURL", "http://127.0.0.1:9999/pac")
	key.SetStringValue("ProxyServer", "127.0.0.1:8888")
	key.SetStringValue("ProxyOverride", "localhost;127.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*;172.32.*;192.168.*;127.0.0.1;<local>")
	// 写入32位整形值
	key.SetDWordValue("MigrateProxy", 0x00000001)
	key.SetDWordValue("ProxyEnable", 0x00000000)
	log.Println("自动设置windows代理相关注册表成功")
}

func GetTargetPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path), nil
}
