package model

import (
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

/*
const (
	FROZEN_LOGIC_TYPE_ENTRUST = 1
	FROZEN_LOGIC_TYPE_ADMIN   = 2
	FROZEN_LOGIC_TYPE_DEAL    = 3
)
*/
type Frozen struct {
	Uid        uint64 `xorm:"comment('用户ID') BIGINT(20)"`
	Ukey       string `xorm:"comment('流水ID') unique(uni) VARCHAR(128)"`
	Num        int64  `xorm:"comment('数量') BIGINT(20)"`
	TokenId    int    `xorm:"comment('币类型') unique(uni)  INT(11)"`
	Type       int    `xorm:"comment('业务使用类型') INT(11)"`
	CreateTime int64  `xorm:"comment('创建时间')  created BIGINT(20)"`
	Opt        int    `xorm:"comment('操作类型') unique(uni) INT(11)"`
}

//新增一条流水
func (s *Frozen) InsertRecord(session *xorm.Session, p *Frozen) (err error) {
	_, err = session.InsertOne(p)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}
