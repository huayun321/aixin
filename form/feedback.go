package form

import (
	"github.com/mholt/binding"
	"net/http"
	"regexp"
	"gopkg.in/mgo.v2/bson"
)

//ArticleCreateForm 文章创建表单
type FeedbackCreateForm struct {
	Content string `json:"content"` //content string minLength 10 maxLength 10000
	Phone   string `json:"phone"`
}

// FieldMap 数据绑定
func (o *FeedbackCreateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.Content: binding.Field{
			Form:         "content",
			Required:     true,
			ErrorMessage: "请填写内容",
		},
		&o.Phone: binding.Field{
			Form:         "phone",
			Required:     true,
			ErrorMessage: "请填写电话",
		},
	}
}

//Validate 数据格式验证
func (o FeedbackCreateForm) Validate(req *http.Request) error {

	if len(o.Content) < 10 {
		return binding.Errors{
			binding.NewError([]string{"content"}, "length error", "文章内容过短."),
		}
	}

	if len(o.Content) > 10000 {
		return binding.Errors{
			binding.NewError([]string{"content"}, "length error", "文章内容过长."),
		}
	}

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

	return nil
}

//ReplyCreateForm 文章创建表单
type ReplyCreateForm struct {
	Content string `json:"content"` //content string minLength 10 maxLength 10000
	FeedbackID string `json:"feedback_id"`
}

// FieldMap 数据绑定
func (o *ReplyCreateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.Content: binding.Field{
			Form:         "content",
			Required:     true,
			ErrorMessage: "请填写内容",
		},
		&o.FeedbackID: binding.Field{
			Form:         "feedback_id",
			Required:     true,
			ErrorMessage: "请填写内容",
		},
	}
}

//Validate 数据格式验证
func (o ReplyCreateForm) Validate(req *http.Request) error {

	if len(o.Content) < 1 {
		return binding.Errors{
			binding.NewError([]string{"content"}, "length error", "文章内容过短."),
		}
	}

	if len(o.Content) > 10000 {
		return binding.Errors{
			binding.NewError([]string{"content"}, "length error", "文章内容过长."),
		}
	}


	return nil
}


//FeedbackListForm
type FeedbackListForm struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	IsTracked int    `json:"is_tracked"`
	Phone     string `json:"phone"`
	TimeStart int    `json:"time_start"`
	TimeEnd   int    `json:"time_end"`
	NickName  string `json:"nickname"`
	UserID    string `json:"user_id"`
}

// FieldMap 数据绑定
func (o *FeedbackListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.Page: binding.Field{
			Form: "page",
		},
		&o.PageSize: binding.Field{
			Form: "page_size",
		},
		&o.IsTracked: binding.Field{
			Form: "is_tracked",
		},
		&o.Phone: binding.Field{
			Form: "phone",
		},
		&o.NickName: binding.Field{
			Form: "nickname",
		},
		&o.TimeStart: binding.Field{
			Form: "time_start",
		},
		&o.TimeEnd: binding.Field{
			Form: "time_end",
		},
	}
}

//Validate 数据格式验证
func (o FeedbackListForm) Validate(req *http.Request) error {
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

	//检查手机号长度
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

	if o.UserID != "" {

		if !bson.IsObjectIdHex(o.UserID) {
			return binding.Errors{
				binding.NewError([]string{"user_id"}, "format error", "id 格式不正确."),
			}
		}
	}

	return nil
}
