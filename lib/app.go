package lib

import (
	"time"
	"strings"
	"math/rand"
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


func (s *App) doBid(appInfo *PpAppInfo) {
	canBidResponse := GetCanBidNow(appInfo)
	if canBidResponse == nil {
		return
	}
	canUseListID := []int{}
	nowBig := 0
	for _, item := range canBidResponse.LoanInfos {
		if proceeIDS.Add(item.ListingId) {

			if strings.Compare(item.CreditCode, "AA") != 0 {
				if item.ListingId > nowBig {
					nowBig = item.ListingId
				}

				if item.Rate < 19 {
					Log("no rate  ", item.Rate)
					continue
				}
				go s.doWeb(item.ListingId, item.Amount, item.Remainfunding)


				//canUseListID = append(canUseListID, item.ListingId)
				//if len(canUseListID) >= 10 {
				//	break
				//}
			}
		}
	}
	SetGetBig(nowBig)
	return
	// 暂时不走 api获取详情
	if len(canUseListID) > 0 {
		loanDetailList := GetListDetail(appInfo, canUseListID)
		for _, item := range loanDetailList.LoanInfos {
			if !proceeIDS.AddBid(item.ListingId) {
				continue
			}

			Log("wee use api bid  !!!!!!!!!!!!!  ", item.ListingId)
			lastSuccess, _ := time.Parse("2006-01-02T15:04:05", item.LastSuccessBorrowTime)
			lastCha := time.Now().Sub(lastSuccess).Hours() / 24


			// 上次借款30天以内的，不投
			if lastCha < 30 {
				Log("no last cha,  ", lastCha, item.ListingId)
				continue
			}
			// 成功借款次数小于2次的，不投
			if item.SuccessCount < 2 {
				Log("no success num  ", item.SuccessCount, item.ListingId)
				continue
			}

			// 有15天以上逾期的不投
			if item.OverdueMoreCount > 0 {
				Log("no overduce more num,  ", item.OverdueMoreCount, item.ListingId)
				continue
			}

			// 还款次数小于正常借款次数的3倍不投
			if (item.NormalCount + item.OverdueLessCount) < item.SuccessCount * 3 {
				Log("no normal + overduce ," , item.NormalCount, item.OverdueLessCount, item.SuccessCount)
				continue
			}

			// 每10次正常还款才能有2次15天内逾期
			if item.OverdueLessCount > (item.NormalCount * 2 / 10) {
				Log("no overduce less num,  ", item.OverdueLessCount, item.NormalCount, item.ListingId)
				continue
			}

			// 本次金额加上待还不能超过最高负负债的0.7倍
			if (item.Amount + item.OwingAmount) > (item.HighestDebt * 0.7) {
				Log("no amount + owing  ", item.Amount, item.OwingAmount, item.HighestDebt, item.ListingId)
				continue
			}

			beginAmount := 50        // 初始金额50

			// 如果0逾期，加20
			if item.OverdueLessCount == 0 {
				beginAmount += 20
			}

			// 正常还款次数大于成功借款次数5倍，加20，大于4倍加10，或者正常还款次数大于15，加10
			if item.NormalCount >= item.SuccessCount * 5 {
				beginAmount += 20
			} else if item.NormalCount >= item.SuccessCount * 4 {
				beginAmount += 10
			} else if item.NormalCount > 15 {
				beginAmount += 10
			}

			// 上次借款在90天之前加40，在60天之前加20
			if lastCha > 90 {
				beginAmount += 40
			} else if lastCha > 60 {
				beginAmount += 20
			}

			// 如果本次借款加上待还小于最大单次借款本金，加60
			if (item.Amount + item.OwingAmount) < item.HighestPrincipal {
				beginAmount += 40
			}

			// 如果没有欠款了，加10
			if item.OwingAmount == 0 {
				beginAmount += 10
			}

			// 有学历认证，加20
			if item.CertificateValidate == 1 {
				beginAmount += 20
			}

			// 凡是还有债务的借款，最多投60
			if item.OwingAmount > 0 {
				if beginAmount > 60 {
					beginAmount = 60
				}
			}

			if beginAmount > int(item.RemainFunding) {
				beginAmount = int(item.RemainFunding)
			}

			if beginAmount <= 0 {
				Log("000000000000  ", item.ListingId)
				continue
			}

			go BidMoney(item.ListingId, beginAmount, s.users[0].AccessToken, s.users[0].Name, s.users[0].UseHongbao)
			go BidMoney(item.ListingId, beginAmount, s.users[1].AccessToken, s.users[1].Name, s.users[1].UseHongbao)
			go BidMoney(item.ListingId, beginAmount, s.users[2].AccessToken, s.users[2].Name, s.users[2].UseHongbao)

			if beginAmount > 160 {
				go BidMoney(item.ListingId, beginAmount, s.users[3].AccessToken, s.users[3].Name, s.users[3].UseHongbao)
			}

		}
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


		go BidMoney(listid, money, s.users[0].AccessToken, s.users[0].Name, s.users[0].UseHongbao)
		go BidMoney(listid, money, s.users[1].AccessToken, s.users[1].Name, s.users[1].UseHongbao)
		go BidMoney(listid, money, s.users[2].AccessToken, s.users[2].Name, s.users[2].UseHongbao)

		if money > 160 {
			go BidMoney(listid, money, s.users[3].AccessToken, s.users[3].Name, s.users[3].UseHongbao)
		}

		return
	} else {
		Log("********* no fast info ", listid)
	}

	begin := time.Now()

	beginAmount := GetCanBidMoney(listid, amount, remain)
	if beginAmount <= 0 {
		return
	}

	afterInfo := time.Now()

	Log("bid process info #####  ,", listid, begin, afterInfo)
	go BidMoney(listid, beginAmount, s.users[0].AccessToken, s.users[0].Name, s.users[0].UseHongbao)
	go BidMoney(listid, beginAmount, s.users[1].AccessToken, s.users[1].Name, s.users[1].UseHongbao)
	go BidMoney(listid, beginAmount, s.users[2].AccessToken, s.users[2].Name, s.users[2].UseHongbao)

	if beginAmount > 160 {
		go BidMoney(listid, beginAmount, s.users[3].AccessToken, s.users[3].Name, s.users[3].UseHongbao)
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

		go s.doBid(s.UseAppInfo)
		time.Sleep(time.Duration(s.interval + int64(rand.Intn(10))) * time.Millisecond)

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

