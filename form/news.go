package form

import (
	"github.com/mholt/binding"
	"net/http"
	"regexp"
	"gopkg.in/mgo.v2/bson"
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

//NCommentCreateForm 文章创建表单
type NCommentCreateForm struct {
	Content     string `json:"content"` //content string minLength 10 maxLength 1000
	NewsID   string `json:"news_id"`
}

// FieldMap 数据绑定
func (o *NCommentCreateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.NewsID: binding.Field{
			Form:         "article_id",
			Required:     true,
			ErrorMessage: "请填写文章id",
		},
		&o.Content: binding.Field{
			Form:         "content",
			Required:     true,
			ErrorMessage: "请填写内容",
		},
	}
}

//Validate 数据格式验证
func (o NCommentCreateForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.NewsID) {
		return binding.Errors{
			binding.NewError([]string{"article_id"}, "format error", "article_id 格式不正确."),
		}
	}
	if len(o.Content) < 10 {
		return binding.Errors{
			binding.NewError([]string{"content"}, "length error", "内容过短."),
		}
	}

	if len(o.Content) > 1000 {
		return binding.Errors{
			binding.NewError([]string{"content"}, "length error", "内容过长."),
		}
	}

	return nil
}

//NewsIdForm
type NewsIdForm struct {
	ID string `json:"id"`
}

// FieldMap 数据绑定
func (o *NewsIdForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID: binding.Field{
			Form:         "id",
			Required:     true,
			ErrorMessage: "请提交文章id",
		},
	}
}

//Validate 数据格式验证
func (o NewsIdForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.ID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
		}
	}
	return nil
}

//NewsListForm
type NewsListForm struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	IsPublished bool   `json:"is_published"`
	Title     string `json:"title"`
	TimeStart int   `json:"time_start"`
	TimeEnd int `json:"time_end"`
}

// FieldMap 数据绑定
func (o *NewsListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.Page: binding.Field{
			Form: "page",
		},
		&o.PageSize: binding.Field{
			Form: "page_size",
		},
		&o.IsPublished: binding.Field{
			Form: "is_published",
		},
		&o.Title: binding.Field{
			Form: "title",
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
func (o NewsListForm) Validate(req *http.Request) error {
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

	return nil
}

//NewsUpdateForm
type NewsUpdateForm struct {
	ID string `json:"id"`
	Title     string `json:"title"`
	Content string   `json:"content"`
	Image string `json:"image"`
	Position int `json:"position"`
}

// FieldMap 数据绑定
func (o *NewsUpdateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.Title: binding.Field{
			Form: "title",
		},
		&o.Content: binding.Field{
			Form: "content",
		},
		&o.Image: binding.Field{
			Form: "image",
		},
		&o.Position: binding.Field{
			Form: "position",
		},
		&o.ID: binding.Field{
			Form:         "id",
			Required:     true,
			ErrorMessage: "请提交文章id",
		},
	}
}

//Validate 数据格式验证
func (o NewsUpdateForm) Validate(req *http.Request) error {
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

	if !bson.IsObjectIdHex(o.ID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
		}
	}

	return nil
}
