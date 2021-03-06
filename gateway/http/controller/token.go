package controller

import (
	"digicon/common/convert"
	"digicon/gateway/rpc"
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"digicon/common/check"
	cf "digicon/gateway/conf"
)

type TokenGroup struct{}

func (s *TokenGroup) Router(r *gin.Engine) {
	action := r.Group("/token", TokenVerify)
	{
		action.POST("/entrust_order", s.EntrustOrder)
		//action.GET("/market/history/kline", s.HistoryKline)

		action.GET("/self_symbols", s.SelfSymbols)

		action.GET("/entrust_list", s.EntrustList)

		action.GET("/entrust_history", s.EntrustHistory)

		action.GET("/balance", s.TokenBalance)

		action.GET("/balance_list", s.TokenBalanceList)

		action.GET("/trade_list", s.TokenTradeList)

		action.POST("/del_entrust", s.DelEntrust)

		action.POST("/transfer_to_currency", s.TransferToCurrency)

		action.GET("/bibi_history", s.BibiHistory)

		action.GET("/transfer_list", s.TransferList)

		action.GET("/refund_list", s.RefundList)
	}
}

func (s *TokenGroup) EntrustOrder(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type EntrustOrderParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		Symbol  string `form:"symbol" binding:"required"`
		Opt     int32  `form:"opt"  binding:"required"`
		OnPrice string `form:"on_price"  `
		Type    int32  `form:"type" binding:"required"`
		Num     string `form:"num" binding:"required"`
		Sign    string  `form:"sign" binding:"required"`
		NonceStr     string `form:"nonce_str" binding:"required"`
	}
	var param EntrustOrderParam

	if err := c.ShouldBind(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	keys:=make(map[string]interface{},0)
	keys["uid"]=param.Uid
	keys["token"]=param.Token
	keys["symbol"]=param.Symbol
	keys["opt"]=param.Opt
	keys["on_price"]=param.OnPrice
	keys["type"]=param.Type
	keys["num"]=param.Num
	keys["nonce_str"]=param.NonceStr

	str:=check.MakeSignature(keys,cf.AppKey)

	if str!=param.Sign {
		ret.SetErrCode(ERRCODE_SIGN)
		return
	}
	var o int64
	var err error
	if param.Type == int32(proto.ENTRUST_TYPE_LIMIT_PRICE) {
		o, err = convert.StringToInt64By8Bit(param.OnPrice)
		if err != nil {
			ret.SetErrCode(ERRCODE_PARAM, err.Error())
			return
		}
		if o == 0 {
			ret.SetErrCode(ERRCODE_PARAM)
			return
		}
	} else {
		if param.OnPrice == "" {
			o = 0
		}
	}

	n, err := convert.StringToInt64By8Bit(param.Num)
	if err != nil {
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallEntrustOrder(&proto.EntrustOrderRequest{
		Symbol:  param.Symbol,
		Opt:     proto.ENTRUST_OPT(param.Opt),
		OnPrice: o,
		Num:     n,
		Uid:     param.Uid,
		Type:    proto.ENTRUST_TYPE(param.Type),
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}

func (s *TokenGroup) SelfSymbols(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type SelfSymbolsParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		TokenId int32  `form:"token_id" binding:"required"`
	}

	var param SelfSymbolsParam

	if err := c.ShouldBind(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallSelfSymbols(&proto.SelfSymbolsRequest{
		Uid: param.Uid,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection(RET_DATA, rsp.Data)
}

func (s *TokenGroup) EntrustList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type EntrustListParam struct {
		Uid   uint64 `form:"uid" binding:"required"`
		Token string `form:"token" binding:"required"`
		Limit int32  `form:"limit" `
		Page  int32  `form:"page" `
	}

	var param EntrustListParam

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Limit == 0 {
		param.Limit = 5
	}
	if param.Page == 0 {
		param.Page = 1
	}
	rsp, err := rpc.InnerService.TokenService.CallEntrustList(&proto.EntrustHistoryRequest{
		Uid:   param.Uid,
		Limit: param.Limit,
		Page:  param.Page,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list", rsp.Data)
}

func (s *TokenGroup) BibiHistory(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type EntrustListParam struct {
		Uid       uint64 `form:"uid" binding:"required"`
		Token     string `form:"token" binding:"required"`
		Limit     int32  `form:"limit" `
		Page      int32  `form:"page" `
		Symbol    string `form:"symbol"`
		Opt       int32  `form:"opt"`
		States    int32  `form:"states"`
		StartTime int32  `form:"startTime"`
		EndTime   int32  `form:"endTime"`
	}

	var param EntrustListParam

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Limit == 0 {
		param.Limit = 5
	}
	if param.Page == 0 {
		param.Page = 1
	}

	startTime := param.StartTime
	//if startTime == 0 {
	//	startTime = int32(time.Now().Unix()) - 86400
	//}
	endTime := param.EndTime
	//if endTime == 0 {
	//	endTime = int32(time.Now().Unix())
	//}

	rsp, err := rpc.InnerService.TokenService.CallBibiHistory(&proto.BibiHistoryRequest{
		Uid:       param.Uid,
		Limit:     param.Limit,
		Page:      param.Page,
		Symbol:    param.Symbol,
		Opt:       param.Opt,
		States:    param.States,
		StartTime: startTime,
		EndTime:   endTime,
	})

	fmt.Println("打印：", err, rsp)

	type list struct {
		PageIndex int32                                  `json:"page_index"`
		PageSize  int32                                  `json:"page_size"`
		TotalPage int32                                  `json:"total_page"`
		Total     int32                                  `json:"total"`
		Items     []*proto.BibiHistoryResponse_Data_Item `json:"items"`
	}

	newList := &list{
		PageIndex: rsp.Data.PageIndex,
		PageSize:  rsp.Data.PageSize,
		TotalPage: rsp.Data.TotalPage,
		Total:     rsp.Data.Total,
		Items:     rsp.Data.Items,
	}

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Code, rsp.Msg)
	ret.SetDataSection("list", newList)
}

func (s *TokenGroup) EntrustHistory(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type EntrustListParam struct {
		Uid   uint64 `form:"uid" binding:"required"`
		Token string `form:"token" binding:"required"`
		Limit int32  `form:"limit" `
		Page  int32  `form:"page" `
	}

	var param EntrustListParam

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	if param.Limit == 0 {
		param.Limit = 5
	}
	if param.Page == 0 {
		param.Page = 1
	}

	rsp, err := rpc.InnerService.TokenService.CallEntrustHistory(&proto.EntrustHistoryRequest{
		Uid:   param.Uid,
		Limit: param.Limit,
		Page:  param.Page,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list", rsp.Data)
}

func (s *TokenGroup) TokenBalance(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type TokenBalanceParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		TokenId int32  `form:"token_id" binding:"required"`
	}

	var param TokenBalanceParam

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallTokenBalance(&proto.TokenBalanceRequest{
		Uid:     param.Uid,
		TokenId: param.TokenId,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("balance", rsp.Balance)
}

// 代币余额列表
func (s *TokenGroup) TokenBalanceList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type TokenBalanceListParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		TokenId int32  `form:"token_id"`
		NoZero  bool   `form:"no_zero"`
	}

	var param TokenBalanceListParam

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallTokenBalanceList(&proto.TokenBalanceListRequest{
		Uid:     param.Uid,
		NoZero:  param.NoZero,
		TokenId: param.TokenId,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	// 重组list
	type NewList struct {
		TokenId   int32  `json:"token_id"`
		TokenName string `json:"token_name"`
		Balance   string `json:"balance"`
		Frozen    string `json:"frozen"`
		WorthCny  string `json:"worth_cny"`
	}
	newList := make([]*NewList, len(rsp.Data.List))
	for k, v := range rsp.Data.List {
		newList[k] = &NewList{
			TokenId:   v.TokenId,
			TokenName: v.TokenName,
			Balance:   convert.Int64ToStringBy8Bit(v.Balance),
			Frozen:    convert.Int64ToStringBy8Bit(v.Frozen),
			WorthCny:  fmt.Sprintf("%.2f", convert.Int64ToFloat64By8Bit(v.WorthCny)),
		}
	}

	ret.SetErrCode(rsp.Err, rsp.Message)
	ret.SetDataSection("list", newList)
	ret.SetDataSection("total_worth_cny", fmt.Sprintf("%.2f", convert.Int64ToFloat64By8Bit(rsp.Data.TotalWorthCny)))
	ret.SetDataSection("total_worth_btc", convert.Int64ToStringBy8Bit(rsp.Data.TotalWorthBtc))
}

// 代币订单明细
func (s *TokenGroup) TokenTradeList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type TokenTradeListParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		Page    int32  `form:"page" binding:"required"`
		PageNum int32  `form:"page_num"`
	}

	var param TokenTradeListParam

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallTokenTradeList(&proto.TokenTradeListRequest{
		Uid:     param.Uid,
		Page:    param.Page,
		PageNum: param.PageNum,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetErrCode(rsp.Err, rsp.Message)

	// 重组data
	type item struct {
		TradeId   int32  `json:"trade_id"`
		TokenName string `json:"token_name"`
		Opt       int32  `json:"opt"`
		Num       string `json:"num"`
		Fee       string `json:"fee"`
		DealTime  int64  `json:"deal_time"`
	}
	type list struct {
		PageIndex int32   `json:"page_index"`
		PageSize  int32   `json:"page_size"`
		TotalPage int32   `json:"total_page"`
		Total     int32   `json:"total"`
		Items     []*item `json:"items"`
	}

	newItems := make([]*item, len(rsp.Data.Items))
	for k, v := range rsp.Data.Items {
		newItems[k] = &item{
			TradeId:   v.TradeId,
			TokenName: v.TokenName,
			Opt:       v.Opt,
			Num:       convert.Int64ToStringBy8Bit(v.Num + v.Fee),
			Fee:       convert.Int64ToStringBy8Bit(v.Fee),
			DealTime:  v.DealTime,
		}
	}

	newList := &list{
		PageIndex: rsp.Data.PageIndex,
		PageSize:  rsp.Data.PageSize,
		TotalPage: rsp.Data.TotalPage,
		Total:     rsp.Data.Total,
		Items:     newItems,
	}

	ret.SetDataSection("list", newList)
}

func (s *TokenGroup) DelEntrust(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	param := &struct {
		Uid       uint64 `form:"uid" binding:"required"`
		Token     string `form:"token" binding:"required"`
		EntrustId string `form:"entrust_id" binding:"required"`
		Sign    string  `form:"sign" binding:"required"`
		NonceStr     string `form:"nonce_str" binding:"required"`
	}{}

	if err := c.ShouldBind(param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}
	keys:=make(map[string]interface{},0)
	keys["uid"]=param.Uid
	keys["token"]=param.Token
	keys["entrust_id"]=param.EntrustId
	keys["nonce_str"]=param.NonceStr

	str:=check.MakeSignature(keys,cf.AppKey)

	if str!=param.Sign {
		ret.SetErrCode(ERRCODE_SIGN)
		return
	}
	rsp, err := rpc.InnerService.TokenService.CallDelEntrust(&proto.DelEntrustRequest{
		Uid:       param.Uid,
		EntrustId: param.EntrustId,
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}

func (s *TokenGroup) TransferToCurrency(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	param := &struct {
		Uid     uint64  `form:"uid" binding:"required"`
		Token   string  `form:"token" binding:"required"`
		TokenId int32   `form:"token_id" binding:"required"`
		Num     float64 `form:"num" binding:"required"`
	}{}

	if err := c.ShouldBind(param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallTransferToCurrency(&proto.TransferToCurrencyRequest{
		Uid:     param.Uid,
		TokenId: param.TokenId,
		Num:     convert.Float64ToInt64By8Bit(param.Num),
	})

	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}
	ret.SetErrCode(rsp.Err, rsp.Message)
}

// 划转明细
func (s *TokenGroup) TransferList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type TransferParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		Page    int32  `form:"page" binding:"required"`
		PageNum int32  `form:"page_num"`
	}

	var param TransferParam

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallTransferList(&proto.TransferListRequest{
		Uid:     param.Uid,
		Page:    param.Page,
		PageNum: param.PageNum,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetErrCode(rsp.Err, rsp.Message)

	// 重组data
	type item struct {
		Id          int64  `json:"id"`
		TokenId     int32  `json:"token_id"`
		TokenName   string `json:"token_name"`
		Type        int32  `json:"type"`
		Num         string `json:"num"`
		CreatedTime int64  `json:"created_time"`
	}
	type list struct {
		PageIndex int32   `json:"page_index"`
		PageSize  int32   `json:"page_size"`
		TotalPage int32   `json:"total_page"`
		Total     int32   `json:"total"`
		Items     []*item `json:"items"`
	}

	newItems := make([]*item, len(rsp.Data.Items))
	for k, v := range rsp.Data.Items {
		newItems[k] = &item{
			Id:          v.Id,
			TokenId:     v.TokenId,
			TokenName:   v.TokenName,
			Type:        v.Type,
			Num:         convert.Int64ToStringBy8Bit(v.Num),
			CreatedTime: v.TransferTime,
		}
	}

	newList := &list{
		PageIndex: rsp.Data.PageIndex,
		PageSize:  rsp.Data.PageSize,
		TotalPage: rsp.Data.TotalPage,
		Total:     rsp.Data.Total,
		Items:     newItems,
	}

	ret.SetDataSection("list", newList)
}

// 退款明细
func (s *TokenGroup) RefundList(c *gin.Context) {
	ret := NewPublciError()
	defer func() {
		c.JSON(http.StatusOK, ret.GetResult())
	}()

	type RefundListParam struct {
		Uid     uint64 `form:"uid" binding:"required"`
		Token   string `form:"token" binding:"required"`
		Page    int32  `form:"page" binding:"required"`
		PageNum int32  `form:"page_num"`
	}

	var param RefundListParam

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf(err.Error())
		ret.SetErrCode(ERRCODE_PARAM, err.Error())
		return
	}

	rsp, err := rpc.InnerService.TokenService.CallRefundList(&proto.RefundListRequest{
		Uid:     param.Uid,
		Page:    param.Page,
		PageNum: param.PageNum,
	})
	if err != nil {
		ret.SetErrCode(ERRCODE_UNKNOWN, err.Error())
		return
	}

	ret.SetErrCode(rsp.Err, rsp.Message)

	// 重组data
	type item struct {
		Id          int64  `json:"id"`
		TokenId     int32  `json:"token_id"`
		TokenName   string `json:"token_name"`
		Type        int32  `json:"type"`
		Num         string `json:"num"`
		CreatedTime int64  `json:"created_time"`
	}
	type list struct {
		PageIndex int32   `json:"page_index"`
		PageSize  int32   `json:"page_size"`
		TotalPage int32   `json:"total_page"`
		Total     int32   `json:"total"`
		Items     []*item `json:"items"`
	}

	newItems := make([]*item, len(rsp.Data.Items))
	for k, v := range rsp.Data.Items {
		newItems[k] = &item{
			Id:          v.Id,
			TokenId:     v.TokenId,
			TokenName:   v.TokenName,
			Type:        v.Type,
			Num:         convert.Int64ToStringBy8Bit(v.Num),
			CreatedTime: v.CreatedTime,
		}
	}

	newList := &list{
		PageIndex: rsp.Data.PageIndex,
		PageSize:  rsp.Data.PageSize,
		TotalPage: rsp.Data.TotalPage,
		Total:     rsp.Data.Total,
		Items:     newItems,
	}

	ret.SetDataSection("list", newList)
}
