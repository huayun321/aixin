package form

import (
	"github.com/mholt/binding"
	"net/http"
	"regexp"
)

//NewsCreateForm news创建表单
type NewsCreateForm struct {
	Title    string `json:"title"`
	Content  string `json:"content"`   //content string minLength 10 maxLength 10000
	Image    string `json:"image"`
	Position int    `json:"position"`
}

// FieldMap 数据绑定
func (o *NewsCreateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.Title: binding.Field{
			Form:         "title",
			Required:     true,
			ErrorMessage: "请填写内容",
		},
		&o.Content: binding.Field{
			Form:         "content",
			Required:     true,
			ErrorMessage: "请填写内容",
		},
		&o.Image: binding.Field{
			Form:         "image",
			Required:     true,
			ErrorMessage: "请填写图片地址",
		},
		&o.Position: binding.Field{
			Form:         "position",
			Required:     true,
			ErrorMessage: "请填写图片位置",
		},
	}
}

//Validate 数据格式验证
func (o NewsCreateForm) Validate(req *http.Request) error {
	if len(o.Title) < 1 {
		return binding.Errors{
			binding.NewError([]string{"title"}, "length error", "文章标题内容过短."),
		}
	}

	if len(o.Title) > 100 {
		return binding.Errors{
			binding.NewError([]string{"title"}, "length error", "文章标题内容过长."),
		}
	}

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

	//检查图片地址格式
	validImg := regexp.MustCompile(`^http://|https://[0-9A-Za-z\-.]{1,100}\.([[:alpha:]]{2,10})(/[[:graph:]]*)*$`)

	if len(o.Image) < 10 || len(o.Image) > 300 {
		return binding.Errors{
			binding.NewError([]string{"iamge"}, "FormatError",
				"图片地址长度，必须大于等于10位，小于等于300位."),
		}
	}

	iva := validImg.MatchString(o.Image)
	if !iva {
		return binding.Errors{
			binding.NewError([]string{"imgage"}, "FormatError", "图片地址不正确，正确地址例子：http://a.com/a.jpg or https://www.a.com/a.jpg."),
		}
	}

	return nil
}
