package main

import (
	"runtime"
	"study/lib"
)
import (
	"os"
	"io/ioutil"
	"encoding/json"
	"time"
)

type AppInfo struct {
	AppId string
	Key string
	Interval int64
}



func main() {
	lib.Log("cpu: ", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	lib.Log("use cookie:  ", lib.Cookie)

	var aa []lib.User
	file, _ := os.Open("user.json")
	data, _ := ioutil.ReadAll(file)
	json.Unmarshal(data, &aa)
	lib.Log("use users   ", aa)




	var tList []AppInfo
	tData, _ := ioutil.ReadFile("info.json")
	json.Unmarshal(tData, &tList)



	var fastUseApp *lib.App = nil

	for index, item := range tList {
		nowPpInfo := lib.NewPpAppInfo(item.AppId, item.Key)
		app := lib.NewApp(aa)
		app.SetUseInfo(nowPpInfo)
		app.SetInterval(item.Interval * int64(lib.ServerConfig.ServerNum))

		if index == 0 {
			fastUseApp = app
		}
		go app.Do()
	}

	fastBid := lib.NewFastBid(lib.Cookie, aa, fastUseApp)
	go fastBid.Do()

	for {
		time.Sleep(10000 * time.Second)
	}


}