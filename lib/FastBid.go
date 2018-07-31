package lib

import (
	"time"
	"sync"
	"math/rand"
	"sync/atomic"
)

var lockBig = sync.Mutex{}
var getBig int = 0

var alreadyFastBid = NewSet()


var lock = sync.Mutex{}
var FastInfo = map[int]int{}


func SetFastInfo(listid, amount int) {
	lock.Lock()
	defer lock.Unlock()
	FastInfo[listid] = amount
}

func GetFastInfo(listid int) (int, bool) {
	lock.Lock()
	defer lock.Unlock()
	amount, exist := FastInfo[listid]
	return amount, exist
}


func SetGetBig(listId int) {
	lockBig.Lock()
	defer lockBig.Unlock()
	if listId > getBig {getBig = listId}
}

func GetGetBig() int {
	lockBig.Lock()
	defer lockBig.Unlock()
	return getBig
}


type FastBid struct {
	sync.Mutex
	bigFast int                   // fastcheck到的最大的listid编号
	useApp *App
	LastListId int
	Cookie string
	users []User
	intFastBidWait bool
	inFastBidNum int
}

func (s *FastBid) setInBidWait() {
	s.Lock()
	defer s.Unlock()
	s.inFastBidNum += 1
}

func (s *FastBid) releaseInBidWait() {
	s.Lock()
	defer s.Unlock()
	s.inFastBidNum -= 1
}

func (s *FastBid) checkInBidWait() bool {
	s.Lock()
	defer s.Unlock()
	if s.inFastBidNum > 0 {
		return true
	} else {
		return false
	}
}



func (s *FastBid) getBigFast() int {
	s.Lock()
	defer s.Unlock()
	return s.bigFast
}

func (s *FastBid) setBigFast(now int)  {
	s.Lock()
	defer s.Unlock()
	if now > s.bigFast {
		s.bigFast = now
	}
}


func (s *FastBid) Do() {
	for {
		nowGetBig := GetGetBig()
		if s.LastListId >= (nowGetBig + 300) || nowGetBig == 0{
			time.Sleep(200 * time.Millisecond)
			continue
		}
		if s.LastListId <= nowGetBig {
			Log("biglistid refresh  -> ", nowGetBig + 1)
			s.LastListId = nowGetBig + 1
		} else {
			s.LastListId += 1
		}

		//if s.LastListId % ServerConfig.ServerNum != ServerConfig.NodeIndex {
		//	continue
		//}


		go s.doCheck(s.LastListId)


	}
}

/**
* 检查一个标是否是信用标，以及其状态情况
 */
func (s *FastBid) doCheck(listid int) {
	firstSleep := rand.Intn(1000)
	time.Sleep(time.Duration(firstSleep) * time.Millisecond)

	for {

		if s.checkInBidWait() {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		getBig := GetGetBig()

		// api扫出来的标号都比这个大了，直接放弃
		if getBig >= listid {
			Log("already late abandon -> ", listid)
			return
		}
		fastBaseInfo := GetFastListBaseInfo(listid)
		if fastBaseInfo == nil || fastBaseInfo.Result == 404 {
			// 如果是404，可能这个标号确实无效，确实无效可以通过更大的标号都已经存在了来确认，否则等待
			bigFast := s.getBigFast()
			if  bigFast > listid {
				return
			}
			nextSleep := 1000 + rand.Intn(1000)
			time.Sleep(time.Duration(nextSleep) * time.Millisecond)
			continue
		} else if fastBaseInfo.ResultContent.Listing.StatusId == 2  {
			// 已经结束的，例如审核不通过之类
			Log("already over -> ", listid)
			return
		} else if fastBaseInfo.Result == 1 {
			// 存在标，先检查标的利率是否合格，如果合格的话接下去检测是否可以买
			if fastBaseInfo.ResultContent.Listing.ShowRate < 19 ||  fastBaseInfo.ResultContent.Listing.ShowRate > 25{
				return
			}
			s.setBigFast(listid)
			s.fastBid(listid,  fastBaseInfo.ResultContent.Listing.Amount)
			return
		}
	}
}

func (s *FastBid) fastBid(listid int, amout float32) {
	canBidMoney := int32(GetCanBidMoney(listid, amout, amout))
	SetFastInfo(listid, int(canBidMoney))
	if canBidMoney <= 0 {
		return
	}

	go func(money int32) {
		personInfo := GetFastPersonInfo(listid)
		for _, item := range personInfo.ResultContent.UserAuthsList {
			if item.Name == "学历认证" {
				Log("############# have degree --> add 40  ", listid)
				money += 40
			}
		}
		if personInfo.ResultContent.BalAmount > 1000 {
			Log("###########have other debet   ", listid)
			money = 50
		}
		if float32(money) > (amout * 3 / 10) {
			money = int32(amout * 3 / 10)
		}
		atomic.StoreInt32(&canBidMoney, money)
		SetFastInfo(listid, int(money))
	}(canBidMoney)

	Log("!!!!!!!!!!!!!!   wait fast bid -->  ", listid)
	Log("!!!!!!!!!!!!!!   wait fast bid -->  ", listid)
	Log("!!!!!!!!!!!!!!   wait fast bid -->  ", listid)
	Log("!!!!!!!!!!!!!!   wait fast bid -->  ", listid)
	for {

		if alreadyFastBid.Exist(listid) {
			return
		}
		go func() {
			fastBaseInfo := GetFastListBaseInfo(listid)
			if fastBaseInfo == nil {
				return
			}
			if fastBaseInfo.ResultContent.Listing.StatusId == 0{
				return
			}
			if fastBaseInfo.ResultContent.Listing.StatusId == 1 {
				if !alreadyFastBid.AddBid(listid) {
					return
				}

				user := s.users[0]

				// 设置标志位，如果持有了快速投标标志位，那么待会需要释放
				s.setInBidWait()

				Log("+++++++++++look  look fast web status is ok now , wait openapi status  ", listid)
				Log("+++++++++++look  look fast web status is ok now , wait openapi status  ", listid)
				Log("+++++++++++look  look fast web status is ok now , wait openapi status  ", listid)
				time.Sleep(1980 * time.Millisecond)
				Log("!!!!!!!!!!!!!!!!!#### look fast bid now   ", listid)
				Log("!!!!!!!!!!!!!!!!!#### look fast bid now   ", listid)
				Log("!!!!!!!!!!!!!!!!!#### look fast bid now   ", listid)

				lastMoney := int(atomic.LoadInt32(&canBidMoney))
				go BidMoney(listid, lastMoney, user.AccessToken, user.Name, user.UseHongbao)

				time.Sleep(300 * time.Millisecond)
				go BidMoney(listid, lastMoney, s.users[1].AccessToken, s.users[1].Name, s.users[1].UseHongbao)

				time.Sleep(40 * time.Millisecond)
				go BidMoney(listid, lastMoney, s.users[2].AccessToken, s.users[2].Name, s.users[2].UseHongbao)

				time.Sleep(350 * time.Millisecond)
				s.releaseInBidWait()
			} else {
				alreadyFastBid.Add(listid)
			}
		}()
		time.Sleep(130 * time.Millisecond)

	}
}



func NewFastBid(cookie string, users []User, useApp *App) *FastBid{
	return &FastBid{LastListId:0, Cookie:cookie, users:users, useApp:useApp, bigFast:0}
}
