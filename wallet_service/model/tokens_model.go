package models

import(
	. "digicon/wallet_service/utils"
	"errors"
)

type Tokens struct {
	Id        int    `xorm:"not null pk autoincr INT(11)"`
	Name      string `xorm:"default '' comment('货币名称') VARCHAR(64)"`
	Detail    string `xorm:"default '' comment('详情地址') VARCHAR(255)"`
	Signature string `xorm:"default '' comment('签名方式') VARCHAR(255)"`
	Chainid   int    `xorm:"default 0 comment('链id') INT(11)"`
	Github    string `xorm:"default '' comment('项目地址') VARCHAR(255)"`
	Web       string `xorm:"default '' VARCHAR(255)"`
	Mark      string `xorm:"comment('英文标识') CHAR(10)"`
	Logo      string `xorm:"default '' comment('货币logo') VARCHAR(255)"`
	Contract  string `xorm:"default '' comment('合约地址') VARCHAR(255)"`
	Node      string `xorm:"default '' comment('节点地址') VARCHAR(100)"`
	Decimal   int    `xorm:"not null default 1 comment('精度 1个eos最小精度的10的18次方') INT(11)"`
}


func (this *Tokens) GetByid()(bool, error){
	return Engine_common.Id(this.Id).Get(this)
}

func (this *Tokens)GetidByContract(contract string,chainid int) (int,error){
	this.Id = 0
	ok,err := Engine_common.Where("contract=? and chainid=?",contract,chainid).Get(this)
	//fmt.Println(contract,chainid)
	if err != nil{
		return 0, err
	}
	if !ok {
		return 0,errors.New("token not exist")
	}
	return this.Id,nil
}

func (this *Tokens)GetDecimal(id int)(int,error){
	this.Id = id
	ok,err:=this.GetByid()
	if err != nil{
		return 0,err
	}

	if !ok {
		return 0,errors.New("connot find this token")
	}
	return this.Decimal,nil
}