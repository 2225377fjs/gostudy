package lib

import (
	"sync"
	"time"
	"strings"
	"sync/atomic"
	"container/list"
	"sort"
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
var fastWaitLock = sync.Mutex{}


func init() {
	NearIdsChan = list.New()
}

func setInFastWait() {
	fastWaitLock.Lock()
	defer fastWaitLock.Unlock()
	inFastNum += 1
}

func releaseFastWait() {
	fastWaitLock.Lock()
	defer fastWaitLock.Unlock()
	inFastNum -= 1
}

func checkInFastWait() bool {
	fastWaitLock.Lock()
	defer fastWaitLock.Unlock()
	if inFastNum > 0 {
		return true
	} else {
		return false
	}
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

/**
判断债务图，最后有几个连续下降的
 */
func CheckHistory(listid int) (int, int) {
	checkNum := 0

	allMonths := 0

	info := GetWebListStatic(listid)
	if info == nil {
		return checkNum, allMonths
	}

	for _, item := range info.ResultContent.LoanerStatistics.PreviousListings {
		allMonths += item.Months
	}
	temp := info.ResultContent.LoanerStatistics.DebtAmountMap

	var itemList HistoryDebtItemSlice
	for key, value := range temp {
		t, _ := time.Parse("2006-01-02", key)
		nowItem := HistoryDebtItem{Time:t, Debt:value}
		itemList = append(itemList, nowItem)
	}
	if len(itemList) < 5 {
		return checkNum, allMonths
	}
	sort.Sort(itemList)
	for index := len(itemList) - 1; index >= 0; index-- {
		if index > 0 {
			if itemList[index].Debt < itemList[index - 1].Debt {
				checkNum += 1
			}  else {
				break
			}
		} else {
			break
		}
	}
	return checkNum, allMonths
}

func doBid(listId int, bidMoney int) {
	go BidMoney(listId, bidMoney, users[0].AccessToken, users[0].Name, users[0].UseHongbao)
	if bidMoney > 70 {
		go BidMoney(listId, bidMoney, users[3].AccessToken, users[3].Name, users[3].UseHongbao)
	}

	go BidMoneyNew(listId, bidMoney, users[0].AccessToken, users[0].Name, users[0].UseHongbao)
	if bidMoney > 70 {
		go BidMoneyNew(listId, bidMoney, users[3].AccessToken, users[3].Name, users[3].UseHongbao)
	}
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
				if !hasEducation && bidMoney < 100{
					Log("do not have education by fast api  ",useInfo.ListingId)
				} else {

					var bidFlag int32 = 0

					realBidMondy :=  int32(bidMoney)
					if !hasEducation {
						bidMoney = 50
						realBidMondy = 50
						Log("$$$$$$$$$$$$  no education   , but can bid some   ", useInfo.ListingId)
						Log("$$$$$$$$$$$$  no education   , but can bid some   ", useInfo.ListingId)
						go func() {
							// 检测最近的债务曲线，看能保持多少个点的下降
							debtDecreaseContinue, allMonths := CheckHistory(useInfo.ListingId)
							Log("$$$$$$$ no education continue debt decrease is ->  ", debtDecreaseContinue, useInfo.ListingId)
							if debtDecreaseContinue >=8 {
								Log("$$$$$$ wa very very nice will add money  -> ", debtDecreaseContinue, useInfo.ListingId)
								Log("$$$$$$ wa very very nice will add money  -> ", debtDecreaseContinue, useInfo.ListingId)
								atomic.AddInt32(&realBidMondy, 40)
							} else if debtDecreaseContinue >= 4 {
								Log("$$$$$$ wa  very nice will add money  -> ", debtDecreaseContinue, useInfo.ListingId)
								Log("$$$$$$ wa  very nice will add money  -> ", debtDecreaseContinue, useInfo.ListingId)
								atomic.AddInt32(&realBidMondy, 20)
							}

							// 还款数目和借款的项目月份相等，这个其实是一个非常好的事情
							if (useInfo.NormalCount + useInfo.OverdueLessCount) == allMonths {
								Log("@@@@@@@@@@@@@@@@ wawa moth equal  ", allMonths, useInfo.ListingId)
								Log("@@@@@@@@@@@@@@@@ wawa moth equal  ", allMonths, useInfo.ListingId)
								atomic.AddInt32(&realBidMondy, 33)
							}
						}()
					}

					useChan, element := addWaitNearChan()
					stopChan := make(chan int, 100)
					before := time.Now().UnixNano()

					setInFastWait()
					defer releaseFastWait()

					// 根据near的检测进行快速投标的逻辑
					go func() {
						stop := false
						defer removeUseChan(element)
						return
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
								if cha <= 15 && cha >= -15 {
									stop = true
									nearUse := (time.Now().UnixNano() - before) / 1000000
									Log("near use -> ", nearUse)
									time.Sleep(100 * time.Millisecond)
									if !atomic.CompareAndSwapInt32(&bidFlag, 0, 1) {
										return
									}
									Log("@@@@@@@@@@@@@@@@ look bid through near id  ", useInfo.ListingId, nowId)
									Log("@@@@@@@@@@@@@@@@ look bid through near id  ", useInfo.ListingId, nowId)
									bidMoney = int(atomic.LoadInt32(&realBidMondy))
									doBid(useInfo.ListingId, bidMoney)
									break
								}
							}
						}
					}()

					// 如果有学历，可以开启直接的sleep的等待快速投标
					if hasEducation {
						go func() {
							time.Sleep(1970 * time.Millisecond)
							bidMoney = int(atomic.LoadInt32(&realBidMondy))
							if !atomic.CompareAndSwapInt32(&bidFlag, 0, 1) {
								return
							}
							Log("!!!!!!!!!!!!!!bid through wait -->  ", useInfo.ListingId)
							Log("!!!!!!!!!!!!!!bid through wait -->  ", useInfo.ListingId)
							Log("!!!!!!!!!!!!!!bid through wait -->  ", useInfo.ListingId)
							doBid(useInfo.ListingId, bidMoney)
						}()
					}

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
									bidMoney = int(atomic.LoadInt32(&realBidMondy))
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									doBid(useInfo.ListingId, bidMoney)
								}
							} else if hasEducation && statusUse > 1965 {
								if atomic.CompareAndSwapInt32(&bidFlag, 0, 1) {
									bidMoney = int(atomic.LoadInt32(&realBidMondy))
									Log("status uss too much, bid direct , now use--> ", statusUse)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									Log("+++++ look fast api bid   ->  ", bidMoney, useInfo.ListingId)
									doBid(useInfo.ListingId, bidMoney)
								}
							}
						}()
						time.Sleep(15 * time.Millisecond)
					}
					stopChan <- 1
					time.Sleep(300 * time.Millisecond)
				}
			}
		}(info)

	}
}

func RefastCheck(listid int) {
	for {
		if checkInFastWait() {
			time.Sleep(200 * time.Millisecond)
		} else {
			break
		}
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
}

func DoFastApi() {
	for {
		if checkInFastWait() {
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
		time.Sleep(500 * time.Millisecond)
	}
}
