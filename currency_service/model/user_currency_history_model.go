package model

import (
	"digicon/common/errors"
	"digicon/currency_service/dao"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type UserCurrencyHistory struct {
	Id           int    `json:"id"                  xorm:"not null pk autoincr comment('ID') INT(10)"`
	Uid          int32  `json:"uid"                 xorm:"not null default 0 INT(10)"`
	TradeUid     int32  `json:"trade_uid"           xorm:"not null default 0 INT(10)"`
	OrderId      string `json:"order_id"            xorm:"not null default '' comment('订单ID') VARCHAR(64)"`
	TokenId      int    `json:"token_id"            xorm:"not null default 0 comment('货币类型') INT(10)"`
	Num          int64  `json:"num"                 xorm:"not null default 0 comment('数量') BIGINT(64)"`
	Fee          int64  `json:"fee"                 xorm:"not null default 0 comment('手续费用') BIGINT(64)"`
	Surplus      int64  `json:"surplus"             xorm:"comment('账户余额') BIGINT(64)"`
	Operator     int    `json:"operator"            xorm:"not null default 0 comment('操作类型 1订单转入 2订单转出 3充币 4提币 5冻结') TINYINT(2)"`
	Address      string `json:"address"             xorm:"not null default '' comment('提币地址') VARCHAR(255)"`
	States       int    `json:"states"              xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(2)"`
	CreatedTime  string `json:"created_time"        xorm:"not null comment('创建时间') DATETIME"`
	UpdatedTime  string `json:"updated_time"        xorm:"comment('修改时间') DATETIME"`
	TransferTime int64  `json:"transfer_time		xorm:"transfer_time"`
}

func (UserCurrencyHistory) TableName() string {
	return "user_currency_history"
}

func (this *UserCurrencyHistory) GetHistory(startTime, endTime string, limit int32) (uhistory []UserCurrencyHistory, err error) {
	now := time.Now()
	if startTime == "" {
		startTime = now.Format("2006-01-02")
	}
	if endTime == "" {
		endTime = now.Format("2006-01-02")
	}
	engine := dao.DB.GetMysqlConn()
	if limit != 0 {
		err = engine.Where("created_time >= ? && created_time <= ?", startTime, endTime).Limit(int(limit)).Find(&uhistory)
	} else {
		err = engine.Where("created_time >= ? && created_time <= ?", startTime, endTime).Find(&uhistory)
	}
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

func (this *UserCurrencyHistory) GetAssetDetail(uid int32, Page uint32, PageNum uint32) (
	uAssetDetails []UserCurrencyHistory, total int64, rPage uint32, rPageNum uint32, err error) {
	if uid <= 0 {
		return
	}
	engine := dao.DB.GetMysqlConn()
	if Page <= 1 {
		Page = 1
	}
	if PageNum <= 0 {
		PageNum = 10
	}
	cAssetDetail := new(UserCurrencyHistory)
	query := engine.Desc("created_time")
	query = query.Where("uid=? and (operator= 1 or operator = 2 )", uid)
	tmpQuery := *query
	countQuery := &tmpQuery
	err = query.Limit(int(PageNum), (int(Page)-1)*int(PageNum)).Find(&uAssetDetails)
	total, _ = countQuery.Count(cAssetDetail)
	if err != nil {
		log.Errorln(err.Error())
		total = 0
		rPage = 0
		rPageNum = 0
	} else {
		total = total
		rPage = Page
		rPageNum = PageNum
	}
	return
}

/*

 */
func (this *UserCurrencyHistory) GetLastPrice(tokenId uint32) (err error, price int64) {
	type NewPrice struct {
		Price int64 `json:"price"`
	}
	//sql := "SELECT  id, price, num FROM g_currency.`order`  WHERE  order_id = ( SELECT  order_id FROM g_currency.`user_currency_history`  where token_id=`?` ORDER BY id DESC LIMIT 1)"
	sql := "SELECT id, price ,num FROM  `order` WHERE token_id=? AND states = 3 ORDER BY confirm_time  DESC  LIMIT 1"
	nprice := NewPrice{}
	engine := dao.DB.GetMysqlConn()
	ok, err := engine.SQL(sql, tokenId).Get(&nprice)
	if err != nil {
		return
	}
	if ok {
		price = nprice.Price
		return
	}
	return
}

// 检查币币划转到法币消息是否已处理
func (this *UserCurrencyHistory) IsTransferFromTokenHandled(transferId int64) (bool, *UserCurrencyHistory, error) {
	has, err := dao.DB.GetMysqlConn().Where(fmt.Sprintf("order_id='%d'", transferId)).And("operator=3").Get(this)
	if err != nil {
		return false, nil, errors.NewSys(err)
	}

	return has, this, nil
}
