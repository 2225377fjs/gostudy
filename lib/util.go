package lib

import (
	"time"
	"encoding/json"
	"strings"
	"net/http"
	"io/ioutil"
	"strconv"
)

import ("crypto/tls"
)

var client = &http.Client{Transport: &http.Transport{
	MaxIdleConnsPerHost: 100,
	IdleConnTimeout: 0,
	TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
},}

type CanBidItem struct {
	ListingId int
	Title string
	CreditCode string
	Amount float32
	Rate float32
	Months int
	PayWay int
	Remainfunding float32
}


type CanBidResponse struct {
	Result int
	ResultMessage string
	ResultCode int
	LoanInfos []*CanBidItem
}

type BidRequest struct {
	ListingId int
	Amount int
}

type BidRequestWithHongBao struct {
	ListingId int
	Amount int
	UseCoupon string
}

type LoanDetail struct {
	RemainFunding float32
	CreditCode string
	ListingId int
	Amount float32
	Months int
	CurrentRate float32
	EducationDegree string
	SuccessCount int
	WasteCount int
	CancelCount int
	FailedCount int
	NormalCount int
	OverdueLessCount int
	OverdueMoreCount int
	OwingPrincipal float32
	OwingAmount float32
	AmountToReceive float32
	FirstSuccessBorrowTime string
	CertificateValidate int
	LastSuccessBorrowTime string
	HighestPrincipal float32
	HighestDebt float32
	TotalPrincipal float32
}

type LoanDetailList struct {
	LoanInfos []LoanDetail
}



func (s *CanBidItem) doBid() {
	if s.CreditCode == "AA"{
		if s.Rate >= 12 {
			BidMoney(s.ListingId, 200, AccessToken, "fjs", false)
		}
	} else {
		DoCreaditBid(s.ListingId, int(s.Amount), s.Remainfunding)
	}

}

func (s *CanBidItem) DoBid() {
	if !proceeIDS.Add(s.ListingId) {
		return
	}
	go s.doBid()
}

//func GetIdsThroughWeb() []int {
//	var out []int
//	request, _ := http.NewRequest("GET", "https://invest.ppdai.com/loan/listnew?LoanCategoryId=4&CreditCodes=4%2C5%2C&ListTypes=&Rates=&Months=&AuthInfo=1%2C&BorrowCount=&didibid=&SortType=0&MinAmount=0&MaxAmount=0", nil)
//	request.Header.Set("Connection", "keep-alive")
//	rep, _ := client.Do(request)
//	defer rep.Body.Close()
//	doc, _ := goquery.NewDocumentFromReader(rep.Body)
//	doc.Find(".title").Each(func(i int,  s *goquery.Selection){
//		url, exist := s.Attr("href")
//		if exist {
//			strs := strings.Split(url, "=")
//			if len(strs) == 2 {
//				listid, err := strconv.Atoi(strs[1])
//				if err == nil {
//					out = append(out, listid)
//				}
//			}
//		}
//	})
//	return out
//}


/**
* 通过api获取可投的标列表
 */
func GetCanBidNow(appInfo *PpAppInfo) *CanBidResponse{
	ss, _ := time.ParseDuration("-6s")
	now := time.Now()
	now = now.Add(ss)
	message := BidRequestData{PageIndex: 1, StartDateTime: now.Format("2006-01-02 15:04:05")}
	requestData := GetBidListRequestData(message)
	signedData := appInfo.Signer.SignData(requestData)

	bodyData, _ := json.Marshal(message)
	bodyStr := strings.ToLower(string(bodyData))


	timeStr := time.Now().UTC().Format("2006-01-02 15:04:05")
	reqest, _ := http.NewRequest("POST", "https://openapi.ppdai.com/invest/LLoanInfoService/LoanList", strings.NewReader(bodyStr))
	reqest.Header.Set("Content-Type", "application/json;charset=utf-8")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("X-PPD-APPID", appInfo.Appid)
	reqest.Header.Set("X-PPD-SIGN", signedData)
	reqest.Header.Set("X-PPD-TIMESTAMP", timeStr)
	timeSign := appInfo.Signer.SignData(appInfo.Appid + timeStr)
	reqest.Header.Set("X-PPD-TIMESTAMP-SIGN", timeSign)

	response, err := client.Do(reqest)
	if err != nil{
		Log("http error when get loanlist   ", err)
		return nil
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	obj := &CanBidResponse{}
	json.Unmarshal(body, obj)
	//if obj.Result != 1 {
	//	if obj.ResultCode != 503 {
	//		Log(string(body))
	//	}
	//}

	return obj
}

/**
* 通过api获取标的投资状态
 */
func GetBidStatus(appInfo *PpAppInfo, ids []int) *BidStatusList{
	message := make(map[string][]int)
	message["ListingIds"] = ids

	bodyData, _ := json.Marshal(message)
	bodyStr := strings.ToLower(string(bodyData))

	timeStr := time.Now().UTC().Format("2006-01-02 15:04:05")
	reqest, _ := http.NewRequest("POST", "https://openapi.ppdai.com/listing/openapiNoAuth/batchListingStatusInfo", strings.NewReader(bodyStr))
	reqest.Header.Set("Content-Type", "application/json;charset=utf-8")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("X-PPD-APPID", appInfo.Appid)
	reqest.Header.Set("X-PPD-SIGN", appInfo.EmptySigData)
	reqest.Header.Set("X-PPD-TIMESTAMP", timeStr)
	timeSign := appInfo.Signer.SignData(appInfo.Appid + timeStr)
	reqest.Header.Set("X-PPD-TIMESTAMP-SIGN", timeSign)

	response, _ := client.Do(reqest)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	statusList := &BidStatusList{}
	json.Unmarshal(body, statusList)
	return statusList
}

/**
* 根据api获取标的详情
 */
func GetListDetail(appInfo *PpAppInfo, ids []int) *LoanDetailList{
	message := make(map[string][]int)
	message["ListingIds"] = ids

	bodyData, _ := json.Marshal(message)
	bodyStr := strings.ToLower(string(bodyData))

	timeStr := time.Now().UTC().Format("2006-01-02 15:04:05")
	reqest, _ := http.NewRequest("POST", "https://openapi.ppdai.com/listing/openapiNoAuth/batchListingInfo", strings.NewReader(bodyStr))
	reqest.Header.Set("Content-Type", "application/json;charset=utf-8")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("X-PPD-APPID", appInfo.Appid)
	reqest.Header.Set("X-PPD-SIGN", appInfo.EmptySigData)
	reqest.Header.Set("X-PPD-TIMESTAMP", timeStr)
	timeSign := appInfo.Signer.SignData(appInfo.Appid + timeStr)
	reqest.Header.Set("X-PPD-TIMESTAMP-SIGN", timeSign)

	response, _ := client.Do(reqest)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	detailList := &LoanDetailList{}
	json.Unmarshal(body, detailList)
	return detailList
}

/**
投标接口
@param listid    标的编号
@param mondy     投标金额
@param accessToken   用户访问token
 */
func BidMoney(listid, money int, accessToken , name string, useHongbao bool) {
	Log("bid  ", listid, " ", money, " ", accessToken)

	var signedData string
	if !useHongbao {
		signedData = SingData1("")
	} else {
		signedData = SingData1("UseCoupontrue")
	}

	var bodyData []byte

	if useHongbao {
		message := BidRequestWithHongBao{ListingId:listid, Amount:money, UseCoupon:"true"}
		bodyData, _ = json.Marshal(message)
	} else {
		message := BidRequest{ListingId:listid, Amount:money}
		bodyData, _ = json.Marshal(message)
	}

	bodyStr := strings.ToLower(string(bodyData))


	timeStr := time.Now().UTC().Format("2006-01-02 15:04:05")
	reqest, _ := http.NewRequest("POST", "https://openapi.ppdai.com/invest/BidService/Bidding", strings.NewReader(bodyStr))
	reqest.Header.Set("Content-Type", "application/json;charset=utf-8")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("X-PPD-APPID", AppId1)
	reqest.Header.Set("X-PPD-SIGN", signedData)
	reqest.Header.Set("X-PPD-TIMESTAMP", timeStr)
	timeSign := SingData1("278ae090e15146d0932c45238b3d941a" + timeStr)
	reqest.Header.Set("X-PPD-TIMESTAMP-SIGN", timeSign)
	reqest.Header.Set("X-PPD-ACCESSTOKEN", accessToken)

	response, _ := client.Do(reqest)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Log("bid http error  ", err)
		return
	}

	type bidResponse struct {
		InvestId int
		ListingId int
		ParticipationAmount int
		Result int
	}
	bidRes := &bidResponse{}
	bidRes.Result = -1
	json.Unmarshal(body, bidRes)
	if bidRes.Result == 0 {
		Log("bid ok ", bidRes.ParticipationAmount,  "  ", listid, "   user: ", name, "  ", string(body))
	} else {
		Log("bid fail, ", string(body))
	}

}


func DoCreaditBid(listid int, amount int, remain float32) {
	type infoRequest struct {
		ListingId string
		Source int
	}
	reqeustMessage := &infoRequest{ListingId:strconv.Itoa(listid), Source:1}
	bodyBytes, _ := json.Marshal(reqeustMessage)
	bodyStr := string(bodyBytes)
	bodyStr = strings.ToLower(bodyStr)
	reqest, _ := http.NewRequest("POST", "https://invest.ppdai.com/api/invapi/LoanDetailPcService/showBorrowerStatistics", strings.NewReader(bodyStr))

	reqest.Header.Set("Cookie", Cookie)
	reqest.Header.Set("Connection", "keep-alive")
	response, _ := client.Do(reqest)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	type liststatic struct {
		FirstSuccessDate string
		SuccessNum int
	}

	type previousListItem struct {
		Amount float32
		CreationDate string
	}

	type lonstatic struct {
		ListingStatics liststatic
		Normalnum int                       // 正常还款次数
		OverdueLessNum int                  // 15天内逾期
		OverdueMoreNum int                  // 15天以上逾期
		TotalPrincipal float32              // 总共借款
		OwingAmount float32                 // 剩余欠款
		LoanAmountMax float32               // 最高单次借款
		DebtAmountMax float32               // 最高历史负载
		OverdueDayMap map[string]float32    // 最近几次还款的逾期天数
		PreviousListings []previousListItem  // 最近几次借款的信息
	}

	type staticContent struct {
		LoanerStatistics lonstatic
	}

	type staticResponse struct {
		Result int
		ResultContent staticContent
	}

	aa := &staticResponse{}

	json.Unmarshal(body, aa)
	staticInfo := &aa.ResultContent.LoanerStatistics


	//if amount > 15000 {
	//	Log("no amount , ", amount, listid)
	//	return
	//}

	if staticInfo.ListingStatics.SuccessNum < 2 {
		Log("no success num , ", staticInfo.ListingStatics.SuccessNum, listid)
		return
	}

	if staticInfo.OverdueMoreNum > 0 {
		Log("no overduce more num , ", staticInfo.OverdueMoreNum, listid)
		return
	}

	if staticInfo.OverdueLessNum > (staticInfo.Normalnum * 2 / 10) {
		Log("no overduce less num, ", staticInfo.OverdueLessNum, staticInfo.Normalnum, listid)
		return
	}

	if (float32(amount) + staticInfo.OwingAmount) > (staticInfo.DebtAmountMax * 0.4) {
		Log("no amout + owing ", amount, staticInfo.OwingAmount, staticInfo.DebtAmountMax, listid)
		return
	}

	//if staticInfo.OwingAmount > (staticInfo.LoanAmountMax * 0.3) {
	//	Log("no owing debetamount , ", staticInfo.OwingAmount, staticInfo.LoanAmountMax, listid)
	//	return
	//}

	//if staticInfo.Normalnum < (staticInfo.ListingStatics.SuccessNum * 3) {
	//	Log("no normal num , succs num ", staticInfo.Normalnum, staticInfo.ListingStatics.SuccessNum, listid)
	//	return
	//}

	//if len(staticInfo.PreviousListings) == 0 {
	//	Log("nono previous list ,,, ", listid)
	//	return
	//}

	//nearest := staticInfo.PreviousListings[0].CreationDate
	//t, _ := time.Parse("2006-01-02 15:04:05", nearest)
	//now := time.Now()
	//sub := now.Sub(t)
	//cha := sub.Hours() / 24
	//if cha <= 20 || cha >= 300 {
	//	Log("no111111111    ", cha, listid)
	//	return
	//}
	//
	//if len(staticInfo.OverdueDayMap) == 0 {
	//	return
	//}
	//no3 := false
	//for _, day := range staticInfo.OverdueDayMap {
	//	if day > 0 {
	//		Log("no222222        ", listid, day)
	//		return
	//	} else if day < -5 {
	//		Log("no3333333   ", listid, day)
	//		no3 = true
	//		return
	//	}
	//}q

	//Log(listid, staticInfo.ListingStatics.SuccessNum, staticInfo.ListingStatics.FirstSuccessDate, staticInfo.Normalnum, staticInfo.OverdueLessNum, staticInfo.OverdueMoreNum)
	//Log(staticInfo.TotalPrincipal, staticInfo.OwingAmount, staticInfo.LoanAmountMax, staticInfo.DebtAmountMax)
	//Log(staticInfo.OverdueDayMap)
	//Log(staticInfo.PreviousListings)

	beginAmount := 166

	//if staticInfo.Normalnum > 10 {
	//	beginAmount += 50
	//}
	//
	//if staticInfo.ListingStatics.SuccessNum > 3 {
	//	beginAmount += 60
	//} else if staticInfo.ListingStatics.SuccessNum > 2 {
	//	beginAmount += 30
	//}
	//
	//if staticInfo.Normalnum >= staticInfo.ListingStatics.SuccessNum * 6 {
	//	beginAmount += 100
	//}
	//
	//if staticInfo.OwingAmount > 0 {
	//	beginAmount /= 2
	//	if beginAmount > 100 {
	//		beginAmount = 100
	//	}
	//} else if staticInfo.OwingAmount == 0 {
	//	beginAmount += 50
	//}
	//
	//if beginAmount > 500 {
	//	beginAmount = 500
	//} else if beginAmount < 50 {
	//	beginAmount = 50
	//}
	//
	//if beginAmount > amount * 3 / 10 {
	//	beginAmount = amount * 3 / 10
	//}

	if beginAmount > int(remain) {
		beginAmount = int(remain)
	}

	BidMoney(listid, beginAmount, AccessToken, "fjs", false)
}


/**
根据listid获取基本信息，包括利率，金额等
 */
func GetFastListBaseInfo(listid int) *FastListBaseInfo {
	reqeustMessage := &InfoRequest{ListingId:strconv.Itoa(listid), Source:1}
	bodyBytes, _ := json.Marshal(reqeustMessage)
	bodyStr := string(bodyBytes)
	bodyStr = strings.ToLower(bodyStr)
	reqest, _ := http.NewRequest("POST", "https://invest.ppdai.com/api/invapi/LoanDetailPcService/showListingBaseInfo", strings.NewReader(bodyStr))

	reqest.Header.Set("Cookie", Cookie)
	reqest.Header.Set("Connection", "keep-alive")
	response, err := client.Do(reqest)
	if err != nil {
		Log("http errror when  GetFastListBaseInfo through web     ", err)
		return nil
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	fastBaseInfo := FastListBaseInfo{}
	json.Unmarshal(body, &fastBaseInfo)
	return &fastBaseInfo
}



/**
根据listid获取基本信息，包括利率，金额等
 */
func GetFastPersonInfo(listid int) *FastPersonInfo {
	reqeustMessage := &InfoRequest{ListingId:strconv.Itoa(listid), Source:1}
	bodyBytes, _ := json.Marshal(reqeustMessage)
	bodyStr := string(bodyBytes)
	bodyStr = strings.ToLower(bodyStr)
	reqest, _ := http.NewRequest("POST", "https://invest.ppdai.com/api/invapi/LoanDetailPcService/showBorrowerInfo", strings.NewReader(bodyStr))

	reqest.Header.Set("Cookie", Cookie)
	reqest.Header.Set("Connection", "keep-alive")
	response, err := client.Do(reqest)
	if err != nil {
		Log("http errror when  GetFastListBaseInfo through web     ", err)
		return nil
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	fastInfo := FastPersonInfo{}
	json.Unmarshal(body, &fastInfo)
	return &fastInfo
}


/**
从web上拿数据bid，并就算能投多少钱
 */
func GetCanBidMoney(listid int, amount, remain float32) int {
	reqeustMessage := &InfoRequest{ListingId:strconv.Itoa(listid), Source:1}
	bodyBytes, _ := json.Marshal(reqeustMessage)
	bodyStr := string(bodyBytes)
	bodyStr = strings.ToLower(bodyStr)
	reqest, _ := http.NewRequest("POST", "https://invest.ppdai.com/api/invapi/LoanDetailPcService/showBorrowerStatistics", strings.NewReader(bodyStr))

	reqest.Header.Set("Cookie", Cookie)
	reqest.Header.Set("Connection", "keep-alive")
	response, err := client.Do(reqest)
	if err != nil {
		Log("http errror when get detail through web     ", err)
		return 0
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)


	aa := &StaticResponse{}

	json.Unmarshal(body, aa)
	staticInfo := &aa.ResultContent.LoanerStatistics


	//if !proceeIDS.AddBid(listid) {
	//	return
	//}


	if len(staticInfo.PreviousListings) == 0 {
		Log("nono previous list ,,, ", listid)
		return 0
	}
	nearest := staticInfo.PreviousListings[0].CreationDate
	t, _ := time.Parse("2006-01-02 15:04:05", nearest)
	now := time.Now()
	sub := now.Sub(t)
	lastCha := sub.Hours() / 24

	// 上次借款45天以内的，不投
	if lastCha < 45 {
		Log("no last cha,  ", lastCha, listid)
		return 0
	}
	// 成功借款次数小于2次的，不投
	if staticInfo.ListingStatics.SuccessNum < 2 {
		Log("no success num  ", staticInfo.ListingStatics.SuccessNum, listid)
		return 0
	}

	// 有15天以上逾期的不投
	if staticInfo.OverdueMoreNum > 0 {
		Log("no overduce more num,  ", staticInfo.OverdueMoreNum, listid)
		return 0
	}

	// 还款次数小于正常借款次数的3倍不投
	if (staticInfo.Normalnum + staticInfo.OverdueLessNum) < staticInfo.ListingStatics.SuccessNum * 3 {
		Log("no normal + overduce ," , staticInfo.Normalnum, staticInfo.OverdueLessNum, staticInfo.ListingStatics.SuccessNum, listid)
		return 0
	}

	// 每15次正常还款才能有2次15天内逾期
	if staticInfo.OverdueLessNum > (staticInfo.Normalnum * 2 / 15) {
		Log("no overduce less num,  ", staticInfo.OverdueLessNum, staticInfo.Normalnum, listid)
		return 0
	}

	// 本次金额加上待还不能超过最高负负债的0.7倍
	if (amount + staticInfo.OwingAmount) > (staticInfo.DebtAmountMax * 0.7) {
		Log("no amount + owing  ", amount, staticInfo.OwingAmount, staticInfo.DebtAmountMax, listid)
		return 0
	}

	// 本次借款加上待还不能超过最大单单次借款金额的1.2倍，该条件可以在标很少的时候适当放宽
	if (amount + staticInfo.OwingAmount) > (staticInfo.LoanAmountMax * 1.2) {
		Log("no amout + owing, amount max  ", amount, staticInfo.OwingAmount, staticInfo.LoanAmountMax, listid)
		return 0
	}

	// 待还不能超过最大借款本金的0.5倍，该条件可以在标少的时候移除
	if staticInfo.OwingAmount > (staticInfo.LoanAmountMax * 0.5) {
		Log("no owing, loanmax,  ", staticInfo.OwingAmount, staticInfo.LoanAmountMax, listid)
		return 0
	}

	// 未结清借款超过1个的不投    这个可以调整
	//unfinish := 0
	for _, item := range staticInfo.PreviousListings {
		if item.StatusId != 12 {
			//unfinish += 1
			// 如果有未还完的随借随还，不投
			if strings.Compare(item.Title, "随借随还") == 0 {
				Log("随借随还，nonon  ", listid)
				return 0
			}
		}
	}
	//if unfinish > 1 {
	//	Log("no unfinsh   ", unfinish, listid)
	//	return 0
	//}

	beginAmount := 50        // 初始金额50

	// 如果0逾期，加10
	if staticInfo.OverdueLessNum == 0 {
		beginAmount += 10
	}

	// 正常还款次数大于成功借款次数5倍，+10
	if staticInfo.Normalnum >= staticInfo.ListingStatics.SuccessNum * 5 {
		beginAmount += 10
	}

	// 上次借款在90天之前加40，在60天之前加20
	if lastCha > 90 {
		beginAmount += 20
	} else if lastCha > 60 {
		beginAmount += 10
	}

	// 如果本次借款加上待还小于最大单次借款本金，加20
	if (amount + staticInfo.OwingAmount) < staticInfo.LoanAmountMax {
		beginAmount += 20
	}

	// 最多一次流标，加10
	if staticInfo.ListingStatics.WasteNum <= 1 {
		beginAmount += 10
	}


	// 最近6个月还款不能有逾期，有提前5天以上的最多投80
	no3 := false
	for _, day := range staticInfo.OverdueDayMap {
		if day > 0 {
			Log("no222222        ", listid, day)
			return 0
		} else if day < -5 {
			Log("no3333333   ", listid, day)
			no3 = true
		}
	}



	// 凡是还有债务的借款，或者有提前5天以上的，最多投80
	if staticInfo.OwingAmount > 0 || no3{
		if beginAmount > 80 {
			beginAmount = 80
		}
	}

	// 如果没有欠款了，加20，这个加项放在最后，不会被前面冲突掉
	if staticInfo.OwingAmount == 0 {
		beginAmount += 20
	}


	if beginAmount > int(remain) {
		beginAmount = int(remain)
	}
	beginAmount += 2

	if beginAmount <= 0 {
		Log("000000000000  ", listid)
		return 0
	}
	return beginAmount
}