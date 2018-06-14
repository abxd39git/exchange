package tools

import (
	"digicon/common/sms"
	. "digicon/proto/common"
	cf "digicon/user_service/conf"
	"digicon/user_service/model"
	"encoding/json"
	"fmt"
	"strconv"
)

func Send253YunSms(phone, code string) (rcode int32, msg string) {
	content := fmt.Sprintf("【253云通讯】您好，您的验证码是%s", code)
	ret, err := sms.Send253Sms(phone, cf.SmsAccount, cf.SmsPwd, content, cf.SmsWebUrl)
	if err != nil {
		rcode = ERRCODE_UNKNOWN
		msg = err.Error()
		return
	}
	p := &model.SmsRet{}
	err = json.Unmarshal([]byte(ret), p)
	if err != nil {
		rcode = ERRCODE_UNKNOWN
		msg = err.Error()
		return
	}

	code_, _ := strconv.Atoi(p.Code)
	msg, ok := CheckErrorMessage(int32(code_))
	if ok {
		rcode = ERRCODE_SUCCESS
		msg = msg
	} else {
		rcode = ERRCODE_UNKNOWN
		msg = p.ErrorMsg
	}
	return
}