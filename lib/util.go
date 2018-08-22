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
	"sync"
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




func (s *CanBidItem) doBid() {
	if s.CreditCode == "AA"{
		if s.Rate >= 12 {
			BidMoney(s.ListingId, 200, AccessToken, "fjs", false)
		}
	}
}

func (s *CanBidItem) DoBid() {
	if !proceeIDS.Add(s.ListingId, "") {
		return
	}
	go s.doBid()
}


var early = -3
var earlyLock sync.Mutex

func getEarly() int {
	earlyLock.Lock()
	defer earlyLock.Unlock()
	if early < -20 {
		early = -3
	} else {
		early -= 1
	}
	return early
}


/**
* 通过api获取可投的标列表
 */
func GetCanBidNow(appInfo *PpAppInfo) *CanBidResponse{
	useCha := getEarly()
	haha := strconv.Itoa(useCha) + "s"
	ss, _ := time.ParseDuration(haha)
	now := time.Now()
	now = now.Add(ss)
	message := BidRequestData{PageIndex: 1, StartDateTime: now.Format("2006-01-02 15:04:05")}
	requestData := GetBidListRequestData(message)
	signedData := appInfo.Signer.SignData(requestData)

	bodyData, _ := json.Marshal(message)
	bodyStr := strings.ToLower(string(bodyData))


	timeStr := time.Now().UTC().Format("2006-01-02 15:04:05")
	reqest, _ := http.NewRequest("POST", "https://openapi.ppdai.com/listing/openapiNoAuth/loanList", strings.NewReader(bodyStr))
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

	response, err := client.Do(reqest)
	if err != nil {
		return nil
	}
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

	response, err := client.Do(reqest)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	body, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		return nil
	}

	detailList := &LoanDetailList{}
	json.Unmarshal(body, detailList)
	if detailList.Result != 1 {
		Log("get detail list error ->  ", string(body[:]))
	}
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

	response, httpErr := client.Do(reqest)
	if httpErr != nil {
		Log("bid http error  ", httpErr)
		return
	}
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


func QueryOrder(orderId, listId, accessToken string) *QueryResponse{
	signedData := SingData1("listingId" + listId + "orderId" + orderId)


	message := make(map[string]string)
	message["orderId"] = orderId
	message["listingId"] = listId
	bodyData, _ := json.Marshal(message)
	bodyStr := string(bodyData)


	timeStr := time.Now().UTC().Format("2006-01-02 15:04:05")
	reqest, _ := http.NewRequest("POST", "https://openapi.ppdai.com/listingbid/openapi/queryBid", strings.NewReader(bodyStr))
	reqest.Header.Set("Content-Type", "application/json;charset=utf-8")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("X-PPD-APPID", AppId1)
	reqest.Header.Set("X-PPD-SIGN", signedData)
	reqest.Header.Set("X-PPD-TIMESTAMP", timeStr)
	timeSign := SingData1("278ae090e15146d0932c45238b3d941a" + timeStr)
	reqest.Header.Set("X-PPD-TIMESTAMP-SIGN", timeSign)
	reqest.Header.Set("X-PPD-ACCESSTOKEN", accessToken)

	response, httpErr := client.Do(reqest)
	if httpErr != nil {
		Log("query http error  ", httpErr)
		return nil
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Log("query http error  ", err)
		return nil
	}


	queryRes := &QueryResponse{}
	queryRes.Result = -1
	json.Unmarshal(body, queryRes)
	return queryRes

}

/**
投标接口
@param listid    标的编号
@param mondy     投标金额
@param accessToken   用户访问token
 */
func BidMoneyNew(listid, money int, accessToken , name string, useHongbao bool) {
	Log("++++++++++ new  bid  ", listid, " ", money, " ", accessToken)

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
	reqest, _ := http.NewRequest("POST", "https://openapi.ppdai.com/listing/openapi/bid", strings.NewReader(bodyStr))
	reqest.Header.Set("Content-Type", "application/json;charset=utf-8")
	reqest.Header.Set("Connection", "keep-alive")
	reqest.Header.Set("X-PPD-APPID", AppId1)
	reqest.Header.Set("X-PPD-SIGN", signedData)
	reqest.Header.Set("X-PPD-TIMESTAMP", timeStr)
	timeSign := SingData1("278ae090e15146d0932c45238b3d941a" + timeStr)
	reqest.Header.Set("X-PPD-TIMESTAMP-SIGN", timeSign)
	reqest.Header.Set("X-PPD-ACCESSTOKEN", accessToken)

	response, httpErr := client.Do(reqest)
	if httpErr != nil {
		Log("bid http error  ", httpErr)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Log("bid http error  ", err)
		return
	}

	type bidResponse struct {
		ListingId int
		Amount int
		Result int
		ResultMessage string
		OrderId string
	}
	bidRes := &bidResponse{}
	bidRes.Result = -1
	json.Unmarshal(body, bidRes)
	if len(bidRes.OrderId) == 0 {
		Log("new bid fail ->  ", string(body))
		return
	}

	time.Sleep(3 * time.Second)
	queryRes := QueryOrder(bidRes.OrderId, strconv.Itoa(listid), accessToken)
	if queryRes == nil {
		Log("++++++++++++ query order error  ", listid, "  user: ", name)
		return
	}
	if queryRes.ResultContent.ParticipationAmount > 0 {
		Log("+++++++++++++++ new bid ok ", queryRes.ResultContent.ParticipationAmount,  "  ", listid, "   user: ", name, "  ", string(body))
	} else {
		Log("++++++++++++++++ new bid fail ", queryRes,  "  ", listid, "   user: ", name)
	}

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
从web上获取标的一些信息
 */
func GetWebListStatic(listid int) *StaticResponse{
	reqeustMessage := &InfoRequest{ListingId:strconv.Itoa(listid), Source:1}
	bodyBytes, _ := json.Marshal(reqeustMessage)
	bodyStr := string(bodyBytes)
	bodyStr = strings.ToLower(bodyStr)
	reqest, _ := http.NewRequest("POST", "https://invest.ppdai.com/api/invapi/LoanDetailPcService/showBorrowerStatistics", strings.NewReader(bodyStr))

	reqest.Header.Set("Cookie", Cookie)
	reqest.Header.Set("Connection", "keep-alive")
	response, err := client.Do(reqest)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)


	aa := &StaticResponse{}

	json.Unmarshal(body, aa)
	return aa
}


/**
从web上拿数据bid，并就算能投多少钱
 */
func GetCanBidMoney(listid int, amount, remain float32) (bool, int) {


	personInfoChan := make(chan *FastPersonInfo)
	staticReponseChan := make(chan *StaticResponse)


	// 获取用户信息
	go func() {
		personInfo := GetFastPersonInfo(listid)
		personInfoChan <- personInfo
	}()


	go func() {
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
			staticReponseChan <- nil
			return
		}
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)


		aa := &StaticResponse{}

		json.Unmarshal(body, aa)
		staticReponseChan <- aa

	}()

	aa := <- staticReponseChan
	if aa == nil {
		return false, 0
	}
	staticInfo := &aa.ResultContent.LoanerStatistics

	personInfo := <- personInfoChan
	if personInfo == nil {
		Log("get person info error ", listid)
		return false, 0
	}



	// 检测是否有学历
	hasEducation := false
	for _, item := range personInfo.ResultContent.UserAuthsList {
		if item.Name == "学历认证" {
			hasEducation = true
		}
	}

	isNormalEducation := false
	if hasEducation && personInfo.ResultContent.EducationInfo.StudyStyle == "普通" {
		isNormalEducation = true
	}

	isNormalBachelor := false
	if isNormalEducation && personInfo.ResultContent.EducationInfo.EducationDegree == "本科" {
		isNormalBachelor = true
	}



	//if !proceeIDS.AddBid(listid) {
	//	return
	//}


	if len(staticInfo.PreviousListings) == 0 {
		Log("nono previous list ,,, ", listid)
		return false, 0
	}
	nearest := staticInfo.PreviousListings[0].CreationDate
	t, _ := time.Parse("2006-01-02 15:04:05", nearest)
	now := time.Now()
	sub := now.Sub(t)
	lastCha := sub.Hours() / 24

	// 上次借款45天以内的，不投，学历可以放宽到30天
	if lastCha < 45 {
		if hasEducation && lastCha > 30 {
			Log(" education last cha  okok ", lastCha, listid)
		} else {
			Log("no last cha,  ", lastCha, listid)
			return false, 0
		}
	}
	// 成功借款次数小于2次的，不投，学历放宽到1次
	if staticInfo.ListingStatics.SuccessNum < 2 {
		if hasEducation && staticInfo.ListingStatics.SuccessNum > 1 {
			Log(" education success num okok , ", staticInfo.ListingStatics.SuccessNum, listid)
		} else {
			Log("no success num  ", staticInfo.ListingStatics.SuccessNum, listid)
			return false, 0
		}
	}

	// 有15天以上逾期的不投
	if staticInfo.OverdueMoreNum > 0 {
		Log("no overduce more num,  ", staticInfo.OverdueMoreNum, listid)
		return false, 0
	}

	// 还款次数小于正常借款次数的3倍不投
	if (staticInfo.Normalnum + staticInfo.OverdueLessNum) < staticInfo.ListingStatics.SuccessNum * 3 {
		Log("no normal + overduce ," , staticInfo.Normalnum, staticInfo.OverdueLessNum, staticInfo.ListingStatics.SuccessNum, listid)
		return false, 0
	}

	// 每15次正常还款才能有2次15天内逾期
	if staticInfo.OverdueLessNum > (staticInfo.Normalnum * 2 / 15) {
		Log("no overduce less num,  ", staticInfo.OverdueLessNum, staticInfo.Normalnum, listid)
		return false, 0
	}

	// 本次金额加上待还不能超过最高负负债的0.7倍,普通本科可以放宽到0.9
	if staticInfo.OwingAmount > 0 {
		if (amount + staticInfo.OwingAmount) > (staticInfo.DebtAmountMax * 0.7) {
			if isNormalBachelor && (amount+staticInfo.OwingAmount) < (staticInfo.DebtAmountMax*0.9) {
				Log("普通本科  okok  ", listid)
			} else {
				Log("no amount + owing  ", amount, staticInfo.OwingAmount, staticInfo.DebtAmountMax, listid)
				return false, 0
			}
		}
	} else if amount > staticInfo.LoanAmountMax {
		Log("no amout, lonanmax  ", amount, staticInfo.LoanAmountMax, listid)
		return false, 0
	}



	// 本次借款加上待还不能超过最大单单次借款金额的1.2倍，该条件可以在标很少的时候适当放宽
	if (amount + staticInfo.OwingAmount) > (staticInfo.LoanAmountMax * 1.2) {
		Log("no amout + owing, amount max  ", amount, staticInfo.OwingAmount, staticInfo.LoanAmountMax, listid)
		return false, 0
	}

	// 待还不能超过最大借款本金的0.5倍，该条件可以在标少的时候移除
	if staticInfo.OwingAmount > (staticInfo.LoanAmountMax * 0.5) {
		Log("no owing, loanmax,  ", staticInfo.OwingAmount, staticInfo.LoanAmountMax, listid)
		return false, 0
	}

	// 未结清借款超过1个的不投    这个可以调整
	//unfinish := 0
	for _, item := range staticInfo.PreviousListings {
		if item.StatusId != 12 {
			//unfinish += 1
			// 如果有未还完的随借随还，不投
			if strings.Compare(item.Title, "随借随还") == 0 {
				Log("随借随还，nonon  ", listid)
				return false, 0
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

	// 如果本次借款加上待还小于最大单次借款本金，加10
	if (amount + staticInfo.OwingAmount) < staticInfo.LoanAmountMax {
		beginAmount += 10
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
			return false, 0
		} else if day < -5 {
			Log("no3333333   ", listid, day)
			no3 = true
		}
	}



	// 凡是还有债务的借款，最多投80
	if staticInfo.OwingAmount > 0 {
		beginAmount = 50
	}

	if !no3 {
		beginAmount += 10
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
		return false, 0
	}




	Log("check person info  ", listid)

	// 有学历，普通本科加40，其他加20
	if hasEducation {
		if isNormalBachelor {
			Log("普通本科 add 40  ", listid)
			beginAmount += 40
		} else {
			beginAmount += 20
		}
	}

	if personInfo.ResultContent.BalAmount > 4000 {
		Log("###########have other debet will going to 0   ", listid)
		beginAmount = 0
	} else if personInfo.ResultContent.BalAmount > 1000 {
		Log("###########have other debet   ", listid)
		beginAmount = 52
	}
	if beginAmount > (int(amount) * 3 / 10) {
		beginAmount = int(amount) * 3 / 10
	}

	return hasEducation, beginAmount
}


func GetCanBidMoneyThroughApiDetail(detail *LoanDetail) (bool, int) {

	// 检测是否有学历
	hasEducation := false
	if detail.CertificateValidate == 1 {
		hasEducation = true
	}

	isNormalEducation := false
	if hasEducation && detail.StudyStyle == "普通" {
		isNormalEducation = true
	}

	isNormalBachelor := false
	if isNormalEducation && detail.EducationDegree == "本科" {
		isNormalBachelor = true
	}

	lastSuccess, _ := time.Parse("2006-01-02T15:04:05", detail.LastSuccessBorrowTime)
	lastCha := time.Now().Sub(lastSuccess).Hours() / 24


	if detail.Amount > 15000 {
		Log("amount too much  ", detail.Amount, detail.ListingId)
		return false, 0
	}

	// 上次借款45天以内的，不投，学历可以放宽到30天
	if lastCha < 45 {
		if hasEducation && lastCha > 30 {
			Log(" education last cha  okok ", lastCha, detail.ListingId)
		} else {
			Log("no last cha,  ", lastCha, detail.ListingId)
			return false, 0
		}
	}
	// 成功借款次数小于2次的，不投，普通学历放宽到1次
	if detail.SuccessCount < 2 {
		if isNormalEducation && detail.SuccessCount > 1 {
			Log(" education success num okok , ", detail.SuccessCount, detail.ListingId)
		} else {
			Log("no success num  ", detail.SuccessCount, detail.ListingId)
			return false, 0
		}
	}

	// 有15天以上逾期的不投
	if detail.OverdueMoreCount > 0 {
		Log("no overduce more num,  ", detail.OverdueMoreCount, detail.ListingId)
		return false, 0
	}

	// 还款次数小于正常借款次数的3倍不投
	if (detail.NormalCount + detail.OverdueLessCount) < detail.SuccessCount * 3 {
		Log("no normal + overduce ," , detail.NormalCount, detail.OverdueLessCount, detail.SuccessCount, detail.ListingId)
		return false, 0
	}

	// 每15次正常还款才能有1次15天内逾期
	if detail.OverdueLessCount > (detail.NormalCount * 1 / 15) {
		Log("no overduce less num,  ", detail.OverdueLessCount, detail.NormalCount, detail.ListingId)
		return false, 0
	}

	// 本次金额加上待还不能超过最高负负债的0.7倍,普通本科可以放宽到0.9
	if detail.OwingAmount > 0 {
		if (detail.Amount + detail.OwingAmount) > (detail.HighestDebt * 0.7) {
			if isNormalBachelor && (detail.Amount + detail.OwingAmount) < (detail.HighestDebt * 0.9) {
				Log("普通本科  okok  ", detail.ListingId)
			} else {
				Log("no amount + owing  ", detail.Amount, detail.OwingAmount, detail.HighestDebt, detail.ListingId)
				return false, 0
			}
		}
	} else if detail.Amount > detail.HighestPrincipal {
		Log("no amout, lonanmax  ", detail.Amount, detail.HighestPrincipal, detail.ListingId)
		return false, 0
	}



	// 本次借款加上待还不能超过最大单单次借款金额的1.2倍，该条件可以在标很少的时候适当放宽
	if (detail.Amount + detail.OwingAmount) > (detail.HighestPrincipal * 1.2) {
		Log("no amout + owing, amount max  ", detail.Amount, detail.OwingAmount, detail.HighestPrincipal, detail.ListingId)
		return false, 0
	}

	// 待还不能超过最大借款本金的0.5倍，该条件可以在标少的时候移除
	if detail.OwingAmount > (detail.HighestPrincipal * 0.5) {
		Log("no owing, loanmax,  ", detail.OwingAmount, detail.HighestPrincipal, detail.ListingId)
		return false, 0
	}



	beginAmount := 50        // 初始金额50

	// 如果0逾期，加10
	if detail.OverdueLessCount == 0 {
		beginAmount += 10
	}

	// 正常还款次数大于成功借款次数5倍，+10
	if detail.NormalCount >= detail.SuccessCount * 5 {
		beginAmount += 10
	}

	// 上次借款在90天之前加40，在60天之前加20
	if lastCha > 90 {
		beginAmount += 20
	} else if lastCha > 60 {
		beginAmount += 10
	}

	// 如果本次借款加上待还小于最大单次借款本金，加10
	if (detail.Amount + detail.OwingAmount) < detail.HighestPrincipal {
		beginAmount += 10
	}

	// 最多一次流标，加10
	if detail.WasteCount <= 1 {
		beginAmount += 10
	}

	// 凡是还有债务的借款，最多投80
	if detail.OwingAmount > 0 {
		beginAmount = 50
	}


	// 如果没有欠款了，加20，这个加项放在最后，不会被前面冲突掉
	if detail.OwingAmount == 0 {
		beginAmount += 20
	}


	if beginAmount > int(detail.RemainFunding) {
		beginAmount = int(detail.RemainFunding)
	}
	beginAmount += 2

	if beginAmount <= 0 {
		Log("000000000000  ", detail.ListingId)
		return false, 0
	}




	Log("check person info  ", detail.ListingId)

	// 有学历，普通本科加40，其他加20
	if hasEducation {
		if isNormalBachelor {
			Log("普通本科 add 40  ", detail.ListingId)
			beginAmount += 40
		} else {
			beginAmount += 20
		}
	}

	if beginAmount > (int(detail.Amount) * 3 / 10) {
		beginAmount = int(detail.Amount) * 3 / 10
	}

	return hasEducation, beginAmount


}