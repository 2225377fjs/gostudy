package lib

import "time"

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
	Months   int                        // 月份
}

type HistoryDebtItem struct {
	Time time.Time
	Debt float32
}

type HistoryDebtItemSlice []HistoryDebtItem

func (c HistoryDebtItemSlice) Len() int {
	return len(c)
}
func (c HistoryDebtItemSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c HistoryDebtItemSlice) Less(i, j int) bool {
	return c[i].Time.Sub(c[j].Time).Hours() < 0
}

type lonstatic struct {
	ListingStatics liststatic
	DebtAmountMap map[string]float32    // 那张负债图形
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

type EducationInfo struct {
	EducationDegree string           // 学历，本科，专科
	Graduate        string           // 学校
	StudyStyle      string           // 学习形式，普通，成人

}

type PersonInfo struct {
	BalAmount        float32              // 网贷余额
	UserAuthsList    []AuthsInfo
	EducationInfo    EducationInfo     // 学习形式
	Gender           string            // 性别
}


type FastPersonInfo struct {
	Result int      // 0错误，1成功     -1 异常
	ResultContent PersonInfo
}


// ---------------------------------------------------------

type LoanDetail struct {
	RemainFunding float32                                  // 剩余可投金额
	CreditCode string
	ListingId int
	Amount float32
	Months int
	CurrentRate float32                                    // 利率
	EducationDegree string
	StudyStyle string
	SuccessCount int
	WasteCount int
	CancelCount int
	FailedCount int
	NormalCount int
	OverdueLessCount int                                     // 1-15天内的逾期
	OverdueMoreCount int                                     // 15天以上的逾期
	OwingPrincipal float32                                   // 剩余待还本金
	OwingAmount float32                                      // 待还金额
	AmountToReceive float32                                  // 代收金额
	FirstSuccessBorrowTime string                            // 第一次成功借款时间
	CertificateValidate int                                  // 是否有学历认证
	LastSuccessBorrowTime string                             // 最后一次成功借款时间
	HighestPrincipal float32                                 // 最高单笔借款金额
	HighestDebt float32                                      // 最高负债
	TotalPrincipal float32                                   // 累计借款金额

	DeadLineTimeOrRemindTimeStr string                        // 截止时间   如果正在投标，那么是xx天xx时xx分
}

type LoanDetailList struct {
	Result int
	ResultMessage string
	LoanInfos []LoanDetail
}

// ----------------------------------------------------------------------------------------
type QueryContent struct {
	BidId int
	ListingId int
	ParticipationAmount float32
}

type QueryResponse struct {
	Result int
	ResultMessage string
	ResultContent QueryContent
}
