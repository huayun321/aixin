package form

import (
	"github.com/mholt/binding"
	"net/http"
	"regexp"
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
