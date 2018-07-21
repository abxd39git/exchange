package model

import (
	. "digicon/price_service/dao"
	. "digicon/price_service/log"
	proto "digicon/proto/rpc"
)

type ConfigQuenes struct {
	Id           int64  `xorm:"pk autoincr BIGINT(20)"`
	TokenId      int    `xorm:"comment('交易币') unique(union_quene_id) INT(11)"`
	TokenTradeId int    `xorm:"comment('实际交易币') unique(union_quene_id) INT(11)"`
	Switch       int    `xorm:"comment('开关0关1开') TINYINT(4)"`
	Name         string `xorm:"comment('USDT/BTC') VARCHAR(32)"`
}

type ConfigTokenCny struct {
	TokenId int   `xorm:"not null pk comment(' 币类型') INT(10)"`
	Price   int64 `xorm:"comment('人民币价格') BIGINT(20)"`
}

/*
type ConfigQuenes struct {
	Id           int64  `xorm:"pk autoincr BIGINT(20)"`
	TokenId      int    `xorm:"comment('交易币') unique(union_quene_id) INT(11)"`
	TokenTradeId int    `xorm:"comment('实际交易币') unique(union_quene_id) INT(11)"`
	Switch       int    `xorm:"comment('开关0关1开') TINYINT(4)"`
	Price        int64  `xorm:"comment('初始价格') BIGINT(20)"`
	Name         string `xorm:"comment('USDT/BTC') VARCHAR(32)"`
	Scope        string `xorm:"comment('振幅') DECIMAL(6,2)"`
	Low          int64  `xorm:"comment('最低价') BIGINT(20)"`
	High         int64  `xorm:"comment('最高价') BIGINT(20)"`
	Amount       int64  `xorm:"comment('成交量') BIGINT(20)"`
}
*/
func GetAllQuenes() []ConfigQuenes {
	t := make([]ConfigQuenes, 0)
	err := DB.GetMysqlConn2().Where("switch=1").Find(&t)
	if err != nil {
		Log.Errorln(err.Error())
		return nil
	}
	return t
}

func GetCnyQuenes() []ConfigTokenCny {
	t := make([]ConfigTokenCny, 0)
	err := DB.GetMysqlConn2().Find(&t)
	if err != nil {
		Log.Errorln(err.Error())
		return nil
	}
	return t
}

func InitConfig() {
	g := GetAllQuenes()
	for _, v := range g {
		ConfigQueneData[v.Name] = &proto.ConfigQueneBaseData{
			TokenId:      int32(v.TokenId),
			TokenTradeId: int32(v.TokenTradeId),
			Name:         v.Name,
		}
	}

	h := GetCnyQuenes()
	for _, v := range h {
		ConfigCny[int32(v.TokenId)] = &proto.CnyPriceBaseData{
			TokenId:  int32(v.TokenId),
			CnyPrice: v.Price,
		}
	}
}

var ConfigQueneData map[string]*proto.ConfigQueneBaseData

var ConfigCny map[int32]*proto.CnyPriceBaseData
var IsFinish bool

func GetTokenCnyPrice(token_id int32) int64 {
	g, ok := ConfigCny[token_id]
	if ok {
		return g.CnyPrice
	}
	return 0
}

func init() {
	ConfigQueneData = make(map[string]*proto.ConfigQueneBaseData, 0)
	ConfigCny = make(map[int32]*proto.CnyPriceBaseData, 0)
	IsFinish = false
}
