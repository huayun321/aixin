package form

import (
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"regexp"
)

//todo password md5

//==============================================================微信注册表单

//SignInWxForm 微信注册登陆表单
type SignWxForm struct {
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	Code         string `json:"verify_code" bson:"verify_code"`
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
			Form: "phone",
		},
		&o.Password: binding.Field{
			Form: "password",
		},
		&o.Code: binding.Field{
			Form: "verify_code",
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
	if len(o.Nickname) < 1 || len(o.Nickname) > 30 {
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

//==============================================================用户手机注册表单

//SignInPhoneForm 用户手机登陆表单
type SignInPhoneForm struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// FieldMap 数据绑定
func (o *SignInPhoneForm) FieldMap(req *http.Request) binding.FieldMap {
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
	}
}

//Validate 数据格式验证
func (o SignInPhoneForm) Validate(req *http.Request) error {
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

	return nil
}

//==============================================================获取用户列表表单

//UserListForm 用户列表表单
type UserListForm struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Phone    string `json:"phone"`
	Sex      uint8  `json:"sex"`
	Nickname string `json:"nickname"`
	IsFrozen bool   `json:"is_frozen"`
}

// FieldMap 数据绑定
func (o *UserListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.Page: binding.Field{
			Form: "page",
		},
		&o.PageSize: binding.Field{
			Form: "page_size",
		},
		&o.Phone: binding.Field{
			Form: "phone",
		},
		&o.Sex: binding.Field{
			Form: "sex",
		},
		&o.Nickname: binding.Field{
			Form: "nickname",
		},
		&o.IsFrozen: binding.Field{
			Form: "is_frozen",
		},
	}
}

//Validate 数据格式验证
func (o UserListForm) Validate(req *http.Request) error {
	//页码
	if o.Page < 0 {
		return binding.Errors{
			binding.NewError([]string{"page"}, "size error", "页数不能是负数."),
		}
	}
	//每页数据
	if o.PageSize < 0 {
		return binding.Errors{
			binding.NewError([]string{"page_size"}, "size error", "每页数据数不能是负数."),
		}
	}

	if o.Phone != "" {
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
	}

	if o.Sex != 0 && o.Sex != 1 && o.Sex != 2 {
		return binding.Errors{
			binding.NewError([]string{"sex"}, "FormatError", "性别格式不对，必须是0，1，2中一个"),
		}
	}

	//检查昵称长度
	if o.Nickname != "" {
		if len(o.Nickname) < 1 || len(o.Nickname) > 30 {
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
	}

	return nil
}

//==============================================================冻结用户表单

//FrozeForm 冻结用户表单
type FrozeForm struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

// FieldMap 数据绑定
func (o *FrozeForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID: binding.Field{
			Form:         "id",
			Required:     true,
			ErrorMessage: "请提交用户id",
		},
		&o.Reason: binding.Field{
			Form:         "reason",
			Required:     true,
			ErrorMessage: "请提交冻结原因",
		},
	}
}

//Validate 数据格式验证
func (o FrozeForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.ID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
		}
	}

	if len(o.Reason) > 500 {
		return binding.Errors{
			binding.NewError([]string{"reason"}, "length error", "冻结原因过长."),
		}
	}

	return nil
}

//==============================================================解冻用户表单

//UnFrozeForm 解冻用户表单
type UnFrozeForm struct {
	ID string `json:"id"`
}

// FieldMap 数据绑定
func (o *UnFrozeForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID: binding.Field{
			Form:         "id",
			Required:     true,
			ErrorMessage: "请提交用户id",
		},
	}
}

//Validate 数据格式验证
func (o UnFrozeForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.ID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
		}
	}
	return nil
}

//SetAdminForm 设置管理员表单
type SetAdminForm struct {
	ID string `json:"id"`
}

// FieldMap 数据绑定
func (o *SetAdminForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID: binding.Field{
			Form:         "id",
			Required:     true,
			ErrorMessage: "请提交用户id",
		},
	}
}

//Validate 数据格式验证
func (o SetAdminForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.ID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
		}
	}
	return nil
}


//GetUserByIDForm 用户详细信息表单
type GetUserByIDForm struct {
	ID string `json:"id"`
}

// FieldMap 数据绑定
func (o *GetUserByIDForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID: binding.Field{
			Form:         "id",
			Required:     true,
			ErrorMessage: "请提交用户id",
		},
	}
}

//Validate 数据格式验证
func (o GetUserByIDForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.ID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
		}
	}
	return nil
}
