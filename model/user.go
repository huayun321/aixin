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
	Code          string     `json:"verify_code" bson:"verify_code"`
	OpenID        string     `json:"openid,omitempty" bson:"openid,omitempty"`
	WxUserInfo    wxUserInfo `json:"wx_user_info,omitempty" bson:"wx_user_info,omitempty"`
	IsFrozen      bool       `json:"is_frozen" bson:"is_frozen"`
	CreateTime    int64      `json:"create_time,omitempty" bson:"create_time"`
	LastLoginTime int64      `json:"last_login_time,omitempty" bson:"last_login_time,omitempty"`
}

// FieldMap 数据绑定
func (o *User) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.Phone: binding.Field{
			Form:         "phone",
			Required:     true,
			ErrorMessage: "请提交手机号",
		},
		&o.Password: binding.Field{
			Form:         "password",
			Required:     true,
			ErrorMessage: "请提交用户密码",
		},
		&o.Avatar: binding.Field{
			Form:         "avatar",
			Required:     true,
			ErrorMessage: "请提交用户头像地址",
		},
		&o.Nickname: binding.Field{
			Form:         "nickname",
			Required:     true,
			ErrorMessage: "请提交用户昵称",
		},
		&o.Code: binding.Field{
			Form:         "verify_code",
			Required:     true,
			ErrorMessage: "请提交验证码",
		},
		&o.OpenID: binding.Field{
			Form: "openid",
		},
		&o.WxUserInfo: binding.Field{
			Form: "wx_user_info",
		},
	}
}

//Validate 数据格式验证
func (o User) Validate(req *http.Request) error {
	//检查手机号长度
	if len(o.Phone) < 11 || len(o.Phone) > 11 {
		return binding.Errors{
			binding.NewError([]string{"phone"}, "LengthError", "手机号必须是11位."),
		}
	}
	//检查手机号格式
	var validPhone = regexp.MustCompile(`^1{1}[\d]{10}$`)
	iv := validPhone.MatchString(o.Phone)
	if !iv {
		return binding.Errors{
			binding.NewError([]string{"phone"}, "FormatError", "手机号格式不正确,必须是11位1开头数字。"),
		}
	}

	//检查密码长度
	if len(o.Password) < 6 || len(o.Password) > 30 {
		return binding.Errors{
			binding.NewError([]string{"password"}, "LengthError", "用户密码长度必须大于等于6位，小于等于30位."),
		}
	}

	//检查密码格式
	var validPassword = regexp.MustCompile(`^[[:graph:]]{6,30}$`)
	ivp := validPassword.MatchString(o.Password)
	if !ivp {
		return binding.Errors{
			binding.NewError([]string{"password"}, "FormatError", "密码格式不正确，必须是6至30位alphabetic字母数字或者特殊字符。"),
		}
	}

	//检查头像地址长度
	if len(o.Avatar) < 3 || len(o.Avatar) > 300 {
		return binding.Errors{
			binding.NewError([]string{"avatar"}, "LengthError", "用户头像地址长度，必须大于等于3位，小于等于300位."),
		}
	}

	//检查头像地址格式
	validAvatar := regexp.MustCompile(`^http://|https://[0-9A-Za-z\-.]{1,100}\.([[:alpha:]]{2,10})(/[[:graph:]]*)*$`)
	iva := validAvatar.MatchString(o.Avatar)
	if !iva {
		return binding.Errors{
			binding.NewError([]string{"avatar"}, "FormatError", "头像地址不正确，正确地址例子：http://a.com/a.jpg or https://www.a.com/a.jpg."),
		}
	}

	//检查昵称长度
	if len(o.Avatar) < 1 || len(o.Avatar) > 30 {
		return binding.Errors{
			binding.NewError([]string{"nickname"}, "LengthError", "用户昵称长度，必须大于等于1位，小于等于30位."),
		}
	}

	//检查昵称格式
	validNickname := regexp.MustCompile(`^[a-z0-9A-Z\p{Han}]+(_[a-z0-9A-Z\p{Han}]+)*$`)
	ivn := validNickname.MatchString(o.Nickname)
	if !ivn {
		return binding.Errors{
			binding.NewError([]string{"nickname"}, "FormatError", "昵称格式不正确，正确地址例子：a_bc汉子_汉789字"),
		}
	}

	//检查验证码长度
	if len(o.Code) != 6 {
		return binding.Errors{
			binding.NewError([]string{"password"}, "LengthError", "用户验证码必须等于6位."),
		}
	}

	//检查验证码格式
	var validCode = regexp.MustCompile(`^[0-9a-z]{6}$`)
	ivc := validCode.MatchString(o.Code)
	if !ivc {
		return binding.Errors{
			binding.NewError([]string{"code"}, "FormatError", "验证码格式不正确，必须是6位[0-9a-z]"),
		}
	}
	return nil
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
