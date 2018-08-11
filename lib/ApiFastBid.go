package lib

import (
	"sync"
	"time"
	"strings"
	//"sync/atomic"
	"sync/atomic"
	"container/list"
)

var UseAppInfos []*PpAppInfo
var appUseIndex = 0

var bigIdLock = sync.Mutex{}
var BigExistId = 0
var OverIdMap = NewSet()
var users []User


var NearIdsChan *list.List
var nearLock = sync.Mutex{}

var inFastNum int32 = 0

func init() {
	NearIdsChan = list.New()
}


func addWaitNearChan() (chan int, *list.Element){
	nearLock.Lock()
	defer nearLock.Unlock()
	useChan := make(chan int, 100)
	e := NearIdsChan.PushBack(useChan)
	return useChan, e
}

func AddNearId(listId int) {
	nearLock.Lock()
	defer nearLock.Unlock()
	for e := NearIdsChan.Front(); e != nil; e = e.Next() {
		useChan := e.Value.(chan int)
		useChan <- listId
	}
}

func removeUseChan(e *list.Element) {
	nearLock.Lock()
	defer nearLock.Unlock()
	useChan := e.Value.(chan int)
	close(useChan)
	NearIdsChan.Remove(e)
}



func SetFastUser(fastUsers []User) {
	users = fastUsers
}

func getUseAppInfo() *PpAppInfo {
	appUseIndex += 1
	if appUseIndex >= len(UseAppInfos) {
		appUseIndex = 0
	}
	return UseAppInfos[appUseIndex]
}

func SetBigExistId(id int) {
	bigIdLock.Lock()
	defer bigIdLock.Unlock()
	if id > BigExistId {
		BigExistId = id
	}
}

func GetBigExistId() int {
	bigIdLock.Lock()
	defer bigIdLock.Unlock()
	return BigExistId
}

func doFastCheck(ids []int) {
	appInfo := getUseAppInfo()
	infos := GetListDetail(appInfo, ids)
	if infos == nil {
		return
	}
	for _, info := range infos.LoanInfos {

		deadLine := info.DeadLineTimeOrRemindTimeStr
		// 表示已经结束了，可能是审核不通过之类的
		if strings.Contains(deadLine, "/") {
			OverIdMap.Add(info.ListingId, "timeout")
			continue
		}

		if info.CurrentRate < 18 {
			OverIdMap.Add(info.ListingId, "rate < 18")
			continue
		}


		// 判断是否已经处理过了，如果没有，加入id，防止以后再处理
		if !OverIdMap.Add(info.ListingId, "process" + "->   " + time.Now().Format("2006-01-02 15:04:05")) {
			continue
		}
		go func(useInfo LoanDetail){
			if !OverIdMap.AddBid(useInfo.ListingId) {
				return
			}
			Log("ooooooooooooooooooo wawa bid through fast api bid  ", useInfo.ListingId)
			hasEducation, bidMoney := GetCanBidMoneyThroughApiDetail(&useInfo)
			if bidMoney > 0{
				Log("================ can bid ", useInfo.ListingId, "  ", bidMoney)
				if !hasEducation {
					Log("do not have education by fast api  ",useInfo.ListingId)
				} else {

					useChan, element := addWaitNearChan()
					stopChan := make(chan int)
					before := time.Now().UnixNano()
					var bidFlag int32 = 0

					atomic.AddInt32(&inFastNum, 1)    // 等待计数加1

					// 根据near的检测进行快速投标的逻辑
					go func() {
						stop := false
						for {
							if stop {
								break
							}
							select {
							case <- stopChan:
								stop = true
								break
							case nowId, ok := <- useChan:
								if !ok {
									stop = true
									break
								}
								cha := useInfo.ListingId - nowId
								if cha <= 30 && cha >= -30 {
									nearUse := (time.Now().UnixNano() - before) / 1000000
									Log("near use -> ", nearUse)
									Log("@@@@@@@@@@@@@@@@ look bid through near id  ", useInfo.ListingId, nowId)
									Log("@@@@@@@@@@@@@@@@ look bid through near id  ", useInfo.ListingId, nowId)
									if nearUse < 1950 {
										time.Sleep(time.Duration(1950 - nearUse) * time.Millisecond)
									}
									go BidMoney(useInfo.ListingId, bidMoney, users[0].AccessToken, users[0].Name, users[0].UseHongbao)
									go BidMoney(useInfo.ListingId, bidMoney, users[3].AccessToken, users[3].Name, users[3].UseHongbao)
									stop = true
									break
								}
							}
						}
						removeUseChan(element)

					}()

					// 根据api查询status的快速投标逻辑
					for {
						if atomic.LoadInt32(&bidFlag) != 0 {
							break
						}
						go func() {
							status := GetBidStatus(getUseAppInfo(), []int{useInfo.ListingId})
							statusUse := (time.Now().UnixNano() - before) / 1000000
							if status != nil && len(status.Infos) > 0 {
								if atomic.CompareAndSwapInt32(&bidFlag, 0, 1) {
									Log("status uss --> ", statusUse)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									go BidMoney(useInfo.ListingId, bidMoney, users[0].AccessToken, users[0].Name, users[0].UseHongbao)
									go BidMoney(useInfo.ListingId, bidMoney, users[3].AccessToken, users[3].Name, users[3].UseHongbao)
								}
							} else if statusUse > 1950 {
								if atomic.CompareAndSwapInt32(&bidFlag, 0, 1) {
									Log("status uss too much, bid direct , now use--> ", statusUse)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									go BidMoney(useInfo.ListingId, bidMoney, users[0].AccessToken, users[0].Name, users[0].UseHongbao)
									go BidMoney(useInfo.ListingId, bidMoney, users[3].AccessToken, users[3].Name, users[3].UseHongbao)
								}
							}
						}()
						time.Sleep(30 * time.Millisecond)
					}
					stopChan <- 1
					time.Sleep(300 * time.Millisecond)
					atomic.AddInt32(&inFastNum, -1)    // 计数器减1
				}
			}
		}(info)

	}
}

func RefastCheck(listid int) {
	for {
		if atomic.LoadInt32(&inFastNum) == 0 {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	var doLook []int
	for nowID := listid + 1; nowID <= (listid + 100); nowID++ {
		if OverIdMap.Exist(nowID) {
			continue
		}
		doLook = append(doLook, nowID)
		if len(doLook) >= 10 {
			go doFastCheck(doLook)
			doLook = []int{}
		}
	}
	if len(doLook) > 0 {
		go doFastCheck(doLook)
	}
	time.Sleep(1 * time.Second)
}

func DoFastApi() {
	for {
		if atomic.LoadInt32(&inFastNum) != 0 {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		nowBig := GetBigExistId()
		if nowBig == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		var doLook []int
		for i := nowBig + 1; i <= (nowBig + 100); i++ {
			if OverIdMap.Exist(i) {
				continue
			}
			doLook = append(doLook, i)
			if len(doLook) >= 10 {
				go doFastCheck(doLook)
				doLook = []int{}
			}
		}
		if len(doLook) > 0 {
			go doFastCheck(doLook)
		}
		time.Sleep(880 * time.Millisecond)
	}
}
