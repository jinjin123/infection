package request

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var dlSizes = [...]int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
var ulSizes = [...]int{100, 300, 500, 800, 1000, 1500, 2500, 3000, 3500, 4000} //kB
var client = http.Client{}

func DownloadTest(sURL string, latency time.Duration) float64 {
	dlURL := strings.Split(sURL, "/upload")[0]
	log.Printf("Download Test: ")
	wg := new(sync.WaitGroup)

	// Warming up
	sTime := time.Now()
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go dlWarmUp(wg, dlURL)
	}
	wg.Wait()
	fTime := time.Now()
	// 1.125MB for each request (750 * 750 * 2)
	wuSpeed := 1.125 * 8 * 2 / fTime.Sub(sTime.Add(latency)).Seconds()

	// Decide workload by warm up speed
	workload := 0
	weight := 0
	skip := false
	if 10.0 < wuSpeed {
		workload = 16
		weight = 4
	} else if 4.0 < wuSpeed {
		workload = 8
		weight = 4
	} else if 2.5 < wuSpeed {
		workload = 4
		weight = 4
	} else {
		skip = true
	}

	// Main speedtest
	dlSpeed := wuSpeed
	if skip == false {
		sTime = time.Now()
		for i := 0; i < workload; i++ {
			wg.Add(1)
			go downloadRequest(wg, dlURL, weight)
		}
		wg.Wait()
		fTime = time.Now()
		log.Printf("\n")

		reqMB := dlSizes[weight] * dlSizes[weight] * 2 / 1000 / 1000
		dlSpeed = float64(reqMB) * 8 * float64(workload) / fTime.Sub(sTime).Seconds()
	}

	return dlSpeed
}

func UploadTest(sURL string, latency time.Duration) float64 {
	log.Printf("Upload Test: ")
	wg := new(sync.WaitGroup)

	// Warm up
	sTime := time.Now()
	wg = new(sync.WaitGroup)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go ulWarmUp(wg, sURL)
	}
	wg.Wait()
	fTime := time.Now()
	// 1.0 MB for each request
	wuSpeed := 1.0 * 8 * 2 / fTime.Sub(sTime.Add(latency)).Seconds()

	// Decide workload by warm up speed
	workload := 0
	weight := 0
	skip := false
	if 10.0 < wuSpeed {
		workload = 16
		weight = 9
	} else if 4.0 < wuSpeed {
		workload = 8
		weight = 9
	} else if 2.5 < wuSpeed {
		workload = 4
		weight = 5
	} else {
		skip = true
	}

	// Main speedtest
	ulSpeed := wuSpeed
	if skip == false {
		sTime = time.Now()
		for i := 0; i < workload; i++ {
			wg.Add(1)
			go uploadRequest(wg, sURL, weight)
		}
		wg.Wait()
		fTime = time.Now()
		log.Printf("\n")

		reqMB := float64(ulSizes[weight]) / 1000
		ulSpeed = reqMB * 8 * float64(workload) / fTime.Sub(sTime).Seconds()
	}

	return ulSpeed
}

func dlWarmUp(wg *sync.WaitGroup, dlURL string) {
	size := dlSizes[2]
	url := dlURL + "/random" + strconv.Itoa(size) + "x" + strconv.Itoa(size) + ".jpg"

	resp, _ := client.Get(url)
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)

	wg.Done()
}

func ulWarmUp(wg *sync.WaitGroup, ulURL string) {
	size := ulSizes[4]
	v := url.Values{}
	v.Add("content", strings.Repeat("0123456789", size*100-51))

	resp, _ := client.PostForm(ulURL, v)
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)

	wg.Done()
}

func downloadRequest(wg *sync.WaitGroup, dlURL string, w int) {
	size := dlSizes[w]
	url := dlURL + "/random" + strconv.Itoa(size) + "x" + strconv.Itoa(size) + ".jpg"

	resp, _ := client.Get(url)
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)

	log.Printf(".")
	wg.Done()
}

func uploadRequest(wg *sync.WaitGroup, ulURL string, w int) {
	size := ulSizes[9]
	v := url.Values{}
	v.Add("content", strings.Repeat("0123456789", size*100-51))

	resp, _ := client.PostForm(ulURL, v)
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)

	log.Printf(".")
	wg.Done()
}

func PingTest(sURL string) time.Duration {
	pingURL := strings.Split(sURL, "/upload")[0] + "/latency.txt"

	l := time.Duration(100000000000) // 10sec
	for i := 0; i < 3; i++ {
		sTime := time.Now()
		resp, _ := http.Get(pingURL)
		fTime := time.Now()
		defer resp.Body.Close()
		if fTime.Sub(sTime) < l {
			l = fTime.Sub(sTime)
		}
	}

	log.Println("Latency:", (l / 2.0))
	return l / 2.0
}
