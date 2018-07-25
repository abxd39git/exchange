package model

import (
	. "digicon/proto/common"
	proto "digicon/proto/rpc"
	. "digicon/token_service/dao"
	"errors"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

type UserToken struct {
	Id      int64
	Uid     uint64 `xorm:"unique(currency_uid) INT(11)"`
	TokenId int    `xorm:"comment('币种') unique(currency_uid) INT(11)"`
	Balance int64  `xorm:"comment('余额') BIGINT(20)"`
	Frozen  int64  `xorm:"comment('冻结余额') BIGINT(20)"`
	Version int    `xorm:"version"`
}

type UserTokenWithName struct {
	UserToken `xorm:"extends"`
	TokenName string
}

// 用户币币余额列表
func (s *UserToken) GetUserTokenList(filter map[string]interface{}) ([]UserTokenWithName, error) {
	engine := DB.GetMysqlConn()
	query := engine.Where("1=1")

	// 筛选
	if v, ok := filter["uid"]; ok {
		query.And("ut.uid=?", v)
	}
	if _, ok := filter["no_zero"]; ok {
		query.And("ut.balance!=0 OR ut.frozen!=0")
	}
	if v, ok := filter["token_id"]; ok {
		query.And("ut.token_id=?", v)
	}

	var list []UserTokenWithName
	err := query.
		Table(s).
		Alias("ut").
		Select("ut.*, t.name token_name").
		Join("LEFT", []string{new(Tokens).TableName(), "t"}, "t.id=ut.token_id").
		Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

//获取实体
func (s *UserToken) GetUserToken(uid uint64, token_id int) (err error) {
	var ok bool
	ok, err = DB.GetMysqlConn().Where("uid=? and token_id=?", uid, token_id).Get(s)
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	if !ok {
		s.Uid = uid
		s.TokenId = int(token_id)

		_, err = DB.GetMysqlConn().InsertOne(s)
		if err != nil {
			log.Errorln(err.Error())
			return
		}

		ok, err = DB.GetMysqlConn().Where("uid=? and token_id=?", uid, token_id).Get(s)
		if err != nil {
			log.Errorln(err.Error())
			return
		}

		if !ok {
			errors.New("insert user token err")
			return
		}
	}

	return
}

//加代币数量
func (s *UserToken) AddMoney(session *xorm.Session, num int64) (err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":      num,
				"uid":      s.Uid,
				"token_id": s.TokenId,
				"balance":  s.Balance,
			}).Errorf("add  money  error %s", err.Error())
		}
	}()

	s.Balance += num
	_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("balance").Update(s)
	if err != nil {
		return
	}
	return
}

//加冻结代币数量，如：注册赠送代币默认放到冻结代币里
func (s *UserToken) AddFrozen(session *xorm.Session, num int64) (err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":      num,
				"uid":      s.Uid,
				"token_id": s.TokenId,
				"balance":  s.Balance,
			}).Errorf("add frozen money  error %s", err.Error())
		}
	}()

	s.Frozen += num
	_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("frozen").Update(s)
	if err != nil {
		return
	}
	return
}

//减少代币数量
func (s *UserToken) SubMoney(session *xorm.Session, num int64) (ret int32, err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":      num,
				"uid":      s.Uid,
				"token_id": s.TokenId,
				"balance":  s.Balance,
			}).Errorf("sub money info data error %s", err.Error())
		}
	}()
	if s.Balance < num {
		ret = ERR_TOKEN_LESS
	}

	s.Balance -= num
	_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("balance").Update(s)
	if err != nil {
		return
	}
	return
}

//减代币数量
/*
func (s *UserToken) SubMoney(session *xorm.Session, num int64, ukey string, ty int) (ret int32, err error) {
	m := &MoneyRecord{}
	ok, err := m.CheckExist(ukey, ty)
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	if ok {
		ret = ERR_TOKEN_REPEAT
		return
	}

	if s.Balance < num {
		ret = ERR_TOKEN_LESS
	}
	if session == nil {
		//开始事务入账处理
		session = DB.GetMysqlConn().NewSession()
		defer session.Close()
		err = session.Begin()

		_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Decr("balance", num).Update(&UserToken{})

		if err != nil {
			log.Errorln(err.Error())
			session.Rollback()
			return
		}

		_, err = session.InsertOne(&MoneyRecord{
			Uid:     s.Uid,
			TokenId: s.TokenId,
			Ukey:    ukey,
			Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
			Type:    ty,
		})

		if err != nil {
			log.Errorln(err.Error())
			session.Rollback()
			return
		}

		err = session.Commit()
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	} else {

		//开始事务入账处理
		_, err = session.Where("uid=? and token_id=?", s.Uid, s.TokenId).Decr("balance", num).Update(&UserToken{})

		if err != nil {
			log.Errorln(err.Error())
			return
		}

		_, err = session.InsertOne(&MoneyRecord{
			Uid:     s.Uid,
			TokenId: s.TokenId,
			Ukey:    ukey,
			Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
			Type:    ty,
		})

		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}

	return
}
*/

//冻结资金
func (s *UserToken) SubMoneyWithFronzen(sess *xorm.Session, num int64, entrust_id string, ty int) (ret int32, err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":         num,
				"uid":         s.Uid,
				"token_id":    s.TokenId,
				"balance":     s.Balance,
				"entrusdt_id": entrust_id,
				"ty":          ty,
			}).Errorf("sub  money with fronzen error %s", err.Error())
		}
	}()
	var aff int64
	if s.Balance >= num {
		s.Balance -= num
		s.Frozen += num
		aff, err = sess.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("balance", "frozen").Update(s)
		if err != nil {
			log.Errorln(err.Error())
			ret = ERRCODE_UNKNOWN
			return
		}

		if aff == 0 {
			err = errors.New("update balance err version is wrong")
			ret = ERRCODE_UNKNOWN
			return
		}

		s.Version += 1

		f := Frozen{
			Uid:     s.Uid,
			Ukey:    entrust_id,
			Num:     num,
			TokenId: s.TokenId,
			Type:    ty,
			Opt:     int(proto.TOKEN_OPT_TYPE_ADD),
		}

		_, err = sess.Insert(f)
		if err != nil {
			log.Errorln(err.Error())
			ret = ERRCODE_UNKNOWN
			return
		}
		return
	}

	ret = ERR_TOKEN_LESS
	return
}

//消耗冻结资金
func (s *UserToken) NotifyDelFronzen(sess *xorm.Session, num int64, entrust_id string, ty int) (ret int32, err error) {
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"num":         num,
				"uid":         s.Uid,
				"token_id":    s.TokenId,
				"balance":     s.Balance,
				"entrusdt_id": entrust_id,
				"ty":          ty,
			}).Errorf("notify  money del fronzen error %s", err.Error())
		}
	}()
	if s.Frozen < num {
		ret = ERR_TOKEN_LESS
		return
	}

	var aff int64
	s.Frozen -= num

	aff, err = sess.Where("uid=? and token_id=?", s.Uid, s.TokenId).Cols("frozen").Update(s)
	if err != nil {
		ret = ERRCODE_UNKNOWN
		return
	}
	if aff == 0 {
		err = errors.New("update balance err version is wrong")
		ret = ERRCODE_UNKNOWN
		return
	}
	//s.Version += 1

	f := Frozen{
		Uid:     s.Uid,
		Ukey:    entrust_id,
		Num:     num,
		TokenId: s.TokenId,
		Type:    ty,
		Opt:     int(proto.TOKEN_OPT_TYPE_DEL),
	}

	_, err = sess.Insert(f)
	if err != nil {
		return
	}

	return
}

//返还冻结资金
func (s *UserToken) ReturnFronzen(sess *xorm.Session, num int64, entrust_id string) (err error) {

	return
}

//获取个人资金明细
func (s *UserToken) GetAllToken(uid uint64) []*UserToken {
	r := make([]*UserToken, 0)
	err := DB.GetMysqlConn().Where("uid=?", uid).Find(&r)
	if err != nil {
		log.Errorln(err.Error())
		return nil
	}
	return r
}
