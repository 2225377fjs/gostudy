package lib

type InfoRequest struct {
	ListingId string
	Source int
}

//-------------------------------------------------------

type FastListBaseInfoContentListInfo struct {
	Amount float32                     // 借款金额
	CreditCode string                  // a, b, c, d
	StatusId  int                      // 0 审核，1进行中，2over
	ShowRate float32
}

type FastListBaseInfoContent struct {
	Listing FastListBaseInfoContentListInfo
}

type FastListBaseInfo struct {
	Result int
	ResultContent FastListBaseInfoContent

}


// ------------------------------------------------------------------


type liststatic struct {
	FirstSuccessDate string             // 首次借款成功日期
	SuccessNum int                      // 成功借款次数
	WasteNum int                        // 流标次数
}

type previousListItem struct {
	Amount float32                      // 金额
	CreationDate string                 // 日期
	ListType int                        // 201随借随还
	Title string                        // 标题，可以判断是否是随借随还
	StatusId int                        // 12 表示已经结清了
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

type StaticResponse struct {
	Result int
	ResultContent staticContent
}


// ------------------------------------------------------------------

type BidStatusInfo struct {
	ListingId int
	Status int     // 0 流标   1 满标    2 投标中   3 finish    4 审核失败   5  撤标
}

type BidStatusList struct {
	Infos []BidStatusInfo
	Result int      // 0错误，1成功     -1 异常
}


// ------------------------------------------------------------------


type AuthsInfo struct {
	Code int
	Name string
}

type PersonInfo struct {
	BalAmount float32              // 网贷余额
	UserAuthsList []AuthsInfo
}


type FastPersonInfo struct {
	Result int      // 0错误，1成功     -1 异常
	ResultContent PersonInfo
}