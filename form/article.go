package form

import (
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"regexp"
)

//==============================================================文章创建表单

//ArticleCreateForm 文章创建表单
type ArticleCreateForm struct {
	AuthorID string   `json:"author_id"` //user object id
	Content  string   `json:"content"`   //content string minLength 10 maxLength 10000
	Images   []string `json:"images"`
}

// FieldMap 数据绑定
func (o *ArticleCreateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.AuthorID: binding.Field{
			Form:         "author_id",
			Required:     true,
			ErrorMessage: "请填写作者id",
		},
		&o.Content: binding.Field{
			Form:         "content",
			Required:     true,
			ErrorMessage: "请填写内容",
		},
		&o.Images: binding.Field{
			Form: "images",
		},
	}
}

//Validate 数据格式验证
func (o ArticleCreateForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.AuthorID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
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

	if len(o.Images) > 9 {
		return binding.Errors{
			binding.NewError([]string{"images"}, "size error", "只能上传9张图片."),
		}
	}

	if len(o.Images) > 0 {
		//检查图片地址格式
		validImg := regexp.MustCompile(`^http://|https://[0-9A-Za-z\-.]{1,100}\.([[:alpha:]]{2,10})(/[[:graph:]]*)*$`)

		for i, v := range o.Images {
			if len(v) < 10 || len(v) > 300 {
				return binding.Errors{
					binding.NewError([]string{"imgages[" + string(i) + "]"}, "FormatError", "图片地址长度，必须大于等于10位，小于等于300位."),
				}
			}

			iva := validImg.MatchString(v)
			if !iva {
				return binding.Errors{
					binding.NewError([]string{"imgages[" + string(i) + "]"}, "FormatError", "图片地址不正确，正确地址例子：http://a.com/a.jpg or https://www.a.com/a.jpg."),
				}
			}
		}
	}

	return nil
}

//==============================================================获取文章列表表单

//UserListForm
type ArticleListForm struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Phone      string `json:"phone"`
	Nickname   string `json:"nickname"`
	IsSelected bool   `json:"is_selected"`
	IsDeleted  bool   `json:"is_deleted"`
}

// FieldMap 数据绑定
func (o *ArticleListForm) FieldMap(req *http.Request) binding.FieldMap {
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
		&o.Nickname: binding.Field{
			Form: "nickname",
		},
		&o.IsSelected: binding.Field{
			Form: "is_selected",
		},
		&o.IsDeleted: binding.Field{
			Form: "is_deleted",
		},
	}
}

//Validate 数据格式验证
func (o ArticleListForm) Validate(req *http.Request) error {
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
