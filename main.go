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

	lib.SetFastUser(aa)



	var tList []AppInfo
	tData, _ := ioutil.ReadFile("info.json")
	json.Unmarshal(tData, &tList)





	//var fastUseApp *lib.App = nil

	for _, item := range tList {
		nowPpInfo := lib.NewPpAppInfo(item.AppId, item.Key)
		app := lib.NewApp(aa)
		app.SetUseInfo(nowPpInfo)
		app.SetInterval(item.Interval * int64(lib.ServerConfig.ServerNum))
		lib.UseAppInfos = append(lib.UseAppInfos, app.UseAppInfo)

		//before := time.Now().UnixNano()
		//lib.Log(lib.GetListDetail(app.UseAppInfo, []int{124709096, 124709195}))
		//lib.Log((time.Now().UnixNano() - before) / 1000000)
		//
		//time.Sleep(1000 * time.Hour)
		//
		//begin := 124250589
		//for index := 0; index <= 50; index++ {
		//	temp := []int{}
		//	for i := 0; i < 10; i++  {
		//		begin += 1
		//		temp = append(temp, begin)
		//	}
		//	nowInfo := lib.GetListDetail(app.UseAppInfo, temp)
		//	for _, listInfo := range nowInfo.LoanInfos {
		//		lib.Log(listInfo.ListingId, listInfo.Amount, listInfo.CurrentRate)
		//	}
		//	time.Sleep(1 * time.Second)
		//}
		//lib.Log("overvoer")
		//
		//time.Sleep(10000 * time.Hour)


		//if index == 0 {
		//	fastUseApp = app
		//}
		go app.Do()
	}


	//fastBid := lib.NewFastBid(lib.Cookie, aa, fastUseApp)
	//go fastBid.Do()

	go lib.DoFastApi()

	for {
		time.Sleep(10000 * time.Second)
	}


}