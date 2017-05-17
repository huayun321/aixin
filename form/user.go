package form

import (
	"github.com/mholt/binding"
	"net/http"
	"regexp"
)

//==============================================================微信注册表单

//SignWxForm 微信注册登陆表单
type SignWxForm struct {
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	Code         string `json:"verify_code"`
	WxOpenID     string `json:"openid"`
	WxNickname   string `json:"nickname"`
	WxSex        uint8  `json:"sex"`
	WxProvince   string `json:"province"`
	WxCity       string `json:"city"`
	WxCountry    string `json:"country"`
	WxHeadimgurl string `json:"headimgurl"`
	WxUnionid    string `json:"unionid"`
}

// FieldMap 数据绑定
func (o *SignWxForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.Phone: binding.Field{
			Form:         "phone",
		},
		&o.Password: binding.Field{
			Form:         "password",
		},
		&o.Code: binding.Field{
			Form:         "verify_code",
		},
		&o.WxOpenID: binding.Field{
			Form:         "openid",
			Required:     true,
			ErrorMessage: "请提交 openid",
		},
	}
}

//Validate 数据格式验证
func (o SignWxForm) Validate(req *http.Request) error {
	//检查手机号长度
	if o.Phone != "" && o.Password != "" && o.Code != "" {
		if len(o.Phone) < 11 || len(o.Phone) > 11 {
			return binding.Errors{
				binding.NewError([]string{"phone"}, "LengthError", "手机号必须是11位."),
			}
		}
		//检查手机号格式
		var validPhone = regexp.MustCompile(`^1[\d]{10}$`)
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
	}

	return nil
}






//==============================================================用户手机注册表单

//SignUpPhoneForm 用户手机注册表单
type SignUpPhoneForm struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Code     string `json:"verify_code" bson:"verify_code"`
}

// FieldMap 数据绑定
func (o *SignUpPhoneForm) FieldMap(req *http.Request) binding.FieldMap {
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
	}
}

//Validate 数据格式验证
func (o SignUpPhoneForm) Validate(req *http.Request) error {
	//检查手机号长度
	if len(o.Phone) < 11 || len(o.Phone) > 11 {
		return binding.Errors{
			binding.NewError([]string{"phone"}, "LengthError", "手机号必须是11位."),
		}
	}
	//检查手机号格式
	var validPhone = regexp.MustCompile(`^1[\d]{10}$`)
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
