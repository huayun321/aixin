package model

import (
	"net/http"
	"regexp"

	"github.com/mholt/binding"
)

//===========================================用户信息部分
//todo trim user field

//User 定义用户信息
type User struct {
	Phone         string     `json:"phone"`
	Password      string     `json:"password"`
	Avatar        string     `json:"avatar"`
	Nickname      string     `json:"nickname"`
	OpenID        string     `json:"openid,omitempty" bson:"openid,omitempty"`
	WxUserInfo    wxUserInfo `json:"wx_user_info,omitempty" bson:"wx_user_info,omitempty"`
	IsFrozen      bool       `json:"is_frozen" bson:"is_frozen"`
	CreateTime    int64      `json:"create_time,omitempty" bson:"create_time"`
	LastLoginTime int64      `json:"last_login_time,omitempty" bson:"last_login_time,omitempty"`
}

//定义微信信息
type wxUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        uint8  `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Headimgurl string `json:"headimgurl"`
	Unionid    string `json:"unionid"`
}

//===========================================验证码部分

//VerifyCode 用于存储用户手机验证信息
type VerifyCode struct {
	Phone           string `json:"phone"`
	VerifyCode      string `json:"verify_code" bson:"verify_code"`           //验证码
	VerifyTimestamp int64  `json:"verify_timestamp,omitempty" bson:"verify_timestamp"` //验证码时间戳
	TimesRemainDay  int    `bson:"times_remain_day,omitempty"`                         //每天限制发五条
	LastVerifyDay   int64  `bson:"last_verify_day,omitempty"`                          //如果是新的一天则重制每天剩余条数
}

//FieldMap 数据绑定验证
func (vc *VerifyCode) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&vc.Phone: binding.Field{
			Form:         "phone",
			Required:     true,
			ErrorMessage: "数据格式错误，请提交手机号",
		},
	}
}

//Validate 数据格式验证
func (vc VerifyCode) Validate(req *http.Request) error {
	//检查手机号长度
	if len(vc.Phone) < 11 || len(vc.Phone) > 11 {
		return binding.Errors{
			binding.NewError([]string{"message"}, "LengthError", "手机号必须是11位."),
		}
	}
	//检查手机号格式
	var validPhone = regexp.MustCompile(`^1{1}[\d]{10}$`)
	iv := validPhone.MatchString(vc.Phone)
	if !iv {
		return binding.Errors{
			binding.NewError([]string{"message"}, "FormatError", "手机号格式不正确."),
		}
	}
	return nil
}

//SMSQuery 向短息服务发送到信息格式
type SMSQuery struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
