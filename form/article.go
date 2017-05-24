package form

import (
	"github.com/mholt/binding"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"regexp"
)

//==============================================================文章创建表单

//ArticleCreateForm 文章创建表单
type ArticleCreateForm struct {
	AuthorID string   `json:"author_id"`  //user object id
	Content  string   `json:"content"`    //content string minLength 10 maxLength 10000
	Images   []string `json:"images"`
}

// FieldMap 数据绑定
func (o *ArticleCreateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.AuthorID: binding.Field{
			Form: "author_id",
			Required:true,
			ErrorMessage:"请填写作者id",
		},
		&o.Content: binding.Field{
			Form: "content",
			Required:true,
			ErrorMessage:"请填写内容",
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
