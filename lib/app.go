package lib

import (
	"time"
	"strings"
	"sync/atomic"
)

var proceeIDS = NewSet()

type PpAppInfo struct {
	Appid string
	PrivateKey string
	EmptySigData string
	Signer Signer
}

type User struct {
	Name string
	AccessToken string
	UseHongbao bool
}

func NewPpAppInfo(appid string, key string) *PpAppInfo {
	info := PpAppInfo{Appid:appid, PrivateKey:key}
	info.Signer, _ = GetSigner(key)
	info.EmptySigData = info.Signer.SignData("")
	return &info
}

type App struct {
	UseAppInfo *PpAppInfo
	interval int64
	users []*User
	useIndex int
}

func NewApp(users []User) *App {
	app := App{users: []*User{}}
	for _, user := range users {
		app.AddUser(user)
	}
	app.interval = 130
	return &app
}

func (s *App) SetUseInfo(info *PpAppInfo) {
	s.UseAppInfo = info
}

func (s *App) SetInterval(interval int64) {
	s.interval = interval
}


// 添加一个投标的用户
func (s *App) AddUser(user User) {
	s.users = append(s.users, &user)
}


func (s *App) doApiDetailBid(listIds []int) {
	loanDetailList := GetListDetail(s.UseAppInfo, listIds)
	if loanDetailList == nil {
		Log("&&&&&&&&&&&& get list detail error  ", listIds)
		Log("&&&&&&&&&&&& get list detail error  ", listIds)
		return
	}
	for _, item := range loanDetailList.LoanInfos {
		hasEducation, bidMoney := GetCanBidMoneyThroughApiDetail(&item)
		if bidMoney > 0 {
			if !hasEducation {
				Log("no education  -- normal api ", item.ListingId)
			} else {
				Log("-- - normal api-----------look look bid  ", bidMoney, item.ListingId)
				Log("-- - normal api-----------look look bid  ", bidMoney, item.ListingId)
				Log("-- - normal api-----------look look bid  ", bidMoney, item.ListingId)
				go BidMoney(item.ListingId, bidMoney, s.users[0].AccessToken, s.users[0].Name, s.users[0].UseHongbao)
				go BidMoney(item.ListingId, bidMoney, s.users[3].AccessToken, s.users[3].Name, s.users[3].UseHongbao)
			}
		}

	}
}

func (s *App) doBid(appInfo *PpAppInfo) {
	canBidResponse := GetCanBidNow(appInfo)
	if canBidResponse == nil {
		return
	}
	if canBidResponse.Result != 1 {
		Log(canBidResponse)
		return
	}
	var canUseListID []int
	nowBig := 0
	for _, item := range canBidResponse.LoanInfos {
		if proceeIDS.Add(item.ListingId, "") {
			AddNearId(item.ListingId)
			if strings.Compare(item.CreditCode, "AA") != 0 {
				if item.ListingId > nowBig {
					nowBig = item.ListingId
				}

				if item.Rate < 19 {
					Log("no rate  ", item.Rate, item.ListingId)
					continue
				}
				//go s.doWeb(item.ListingId, item.Amount, item.Remainfunding)

				if OverIdMap.BidExist(item.ListingId) {
					Log("----already bid through fast api  ", item.ListingId)
				} else {
					if OverIdMap.Exist(item.ListingId) {
						Log("+++++++++Fast not process because -> ", OverIdMap.GetReason(item.ListingId), "  ", item.ListingId)
					}
					Log("++++++++++++++ will bid through normal  ", item.ListingId)
					canUseListID = append(canUseListID, item.ListingId)
					if len(canUseListID) >= 10 {
						go s.doApiDetailBid(canUseListID)
						canUseListID = []int{}
					}
				}

			}
		}
	}
	if nowBig > 0 {
		SetBigExistId(nowBig)
		RefastCheck(nowBig)
	}
	if len(canUseListID) > 0 {
		go s.doApiDetailBid(canUseListID)
	}

}


/**
* 先看看有没有之前fast测试的缓存数据，如果有的话就不用重新获取了，直接bid即可
* 没有的话再去获取可以bid的金额
 */
func (s *App) doWeb(listid int, amount, remain float32 ) {
	if money, exist := GetFastInfo(listid); exist {
		if money <= 0 {
			Log("-------  not bid beacuse fast bid info  ", listid, money)
			return
		}
		Log("########## bid through fast info   ", listid)



		if money > 150 {
			go BidMoney(listid, money, s.users[0].AccessToken, s.users[0].Name, s.users[0].UseHongbao)
			go BidMoney(listid, money, s.users[3].AccessToken, s.users[3].Name, s.users[3].UseHongbao)
		}

		return
	}

	begin := time.Now()

	haseEduction, beginAmount := GetCanBidMoney(listid, amount, remain)
	if beginAmount <= 0 {
		return
	}

	afterInfo := time.Now()


	if proceeIDS.AddBid(listid) {
		Log("---- do web bid process info #####  ,", listid, begin, afterInfo)
		if haseEduction {
			go BidMoney(listid, beginAmount, s.users[0].AccessToken, s.users[0].Name, s.users[0].UseHongbao)
			go BidMoney(listid, beginAmount, s.users[3].AccessToken, s.users[3].Name, s.users[3].UseHongbao)
		} else {
			Log("web no education -- ", listid)
		}
	}


}

func (s *App) DoTest(listid, amout int) {
	for _, user := range s.users {
		go BidMoney(listid, amout, user.AccessToken, user.Name, false)
	}
}


func (s *App) Do(){
	for {
		// 测试代码
		//s.doWeb(121711257, 100, 300)
		//time.Sleep(10000 * time.Second)

		// 如果在fast的等待途中，那么降低扫标频率
		if atomic.LoadInt32(&inFastNum) != 0 {
			time.Sleep(200 * time.Millisecond)
		}
		before := time.Now().UnixNano()
		s.doBid(s.UseAppInfo)
		use := (time.Now().UnixNano() - before) / 1000000
		if use < s.interval {
			time.Sleep(time.Duration(s.interval - use) * time.Millisecond)
		}

		//go s.doBid(s.UseAppInfo)
		//time.Sleep(time.Duration(s.interval + int64(rand.Intn(10))) * time.Millisecond)

		//go func(appInfo *PpAppInfo) {
		//	s.doBid(appInfo)
		//}(s.appInfos[s.useIndex])
		//
		//s.useIndex += 1
		//if s.useIndex >= len(s.appInfos) {
		//	s.useIndex = 0
		//}




		//time.Sleep(65 * time.Millisecond)
	}
}

