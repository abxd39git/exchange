package models

import (
	"digicon/wallet_service/utils"
	"time"
)

type TokenChainInout struct {
	Id          int       `xorm:"INT(11)"`
	Txhash      string    `xorm:"comment('交易hash') VARCHAR(255)"`
	From        string    `xorm:"comment('打款地址') VARCHAR(42)"`
	To          string    `xorm:"comment('付款地址') VARCHAR(42)"`
	Value       string    `xorm:"comment('金额') VARCHAR(30)"`
	Contract    string    `xorm:"comment('合约地址') VARCHAR(42)"`
	Chainid     int       `xorm:"comment('链id') INT(11)"`
	Type        int       `xorm:"not null comment('平台转出:0,充值:1') INT(11)"`
	Signtx      string    `xorm:"comment('平台转出记录交易签名') VARCHAR(1024)"`
	Tokenid     int       `xorm:"not null comment('币种id') INT(11)"`
	TokenName   string    `xorm:"not null comment('币名称') VARCHAR(10)"`
	Uid         int       `xorm:"not null comment('用户id') INT(11)"`
	CreatedTime time.Time `xorm:"default 'CURRENT_TIMESTAMP' TIMESTAMP"`
}

func (this *TokenChainInout) Insert(txhash, from, to, value, contract string, chainid int, uid int, tokenid int, tokenname string) (int, error) {
	this.Id = 0
	this.Txhash = txhash
	this.From = from
	this.To = to
	this.Value = value
	this.Contract = contract
	this.Chainid = chainid
	this.Type = 1

	this.Tokenid = tokenid
	this.TokenName = tokenname
	this.Uid = uid
	//utils.Engine_wallet.ShowSQL(true)
	affected, err := utils.Engine_wallet.InsertOne(this)
	//fmt.Println("aaaa",uid,err)
	return int(affected), err
}
func (this *TokenChainInout) TxhashExist(hash string, chainid int) (bool, error) {
	return utils.Engine_wallet.Where("txhash=? and chainid=?", hash, chainid).Get(this)

}
