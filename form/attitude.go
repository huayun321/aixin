package form

import (
	"github.com/mholt/binding"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

//==============================================================文章创建表单

//AttitudeCreateForm 态度创建表单
type AttitudeCreateForm struct {
	AuthorID string   `json:"author_id"` //user object id
	Name       string        `json:"name"`
	Desc       string        `json:"desc"`
	MainImg    string        `json:"main_img"`
	IconImg    string        `json:"icon_img"`
}

// FieldMap 数据绑定
func (o *AttitudeCreateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.AuthorID: binding.Field{
			Form:         "author_id",
			Required:     true,
			ErrorMessage: "请填写作者id",
		},
		&o.Name: binding.Field{
			Form:         "name",
			Required:     true,
			ErrorMessage: "请填写态度名称",
		},
		&o.Desc: binding.Field{
			Form:         "desc",
			ErrorMessage: "请填写态度描述",
		},
		&o.MainImg: binding.Field{
			Form:         "main_img",
			Required:     true,
			ErrorMessage: "请填写态度主图",
		},
		&o.IconImg: binding.Field{
			Form:         "icon_img",
			Required:     true,
			ErrorMessage: "请填写态度小图",
		},
	}
}

//Validate 数据格式验证
func (o AttitudeCreateForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.AuthorID) {
		return binding.Errors{
			binding.NewError([]string{"author_id"}, "format error", "id 格式不正确."),
		}
	}

	return nil
}



//==============================================================获取文章列表表单

//AttitudeListForm
type AttitudeListForm struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Name     string `json:"name"`
}

// FieldMap 数据绑定
func (o *AttitudeListForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.Page: binding.Field{
			Form: "page",
		},
		&o.PageSize: binding.Field{
			Form: "page_size",
		},
		&o.Name: binding.Field{
			Form: "name",
		},
	}
}

//Validate 数据格式验证
func (o AttitudeListForm) Validate(req *http.Request) error {
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

//AttitudeIdForm
type AttitudeIdForm struct {
	ID string `json:"id"`
}

// FieldMap 数据绑定
func (o *AttitudeIdForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID: binding.Field{
			Form:         "id",
			Required:     true,
			ErrorMessage: "请提交文章id",
		},
	}
}

//Validate 数据格式验证
func (o AttitudeIdForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.ID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
		}
	}
	return nil
}



//AttitudeUpdateForm 动作表单
type AttitudeUpdateForm struct {
	ID       string   `json:"id"`
	AuthorID string   `json:"author_id"` //user object id
	Name       string        `json:"name"`
	Desc       string        `json:"desc"`
	MainImg    string        `json:"main_img"`
	IconImg    string        `json:"icon_img"`
}

// FieldMap 数据绑定
func (o *AttitudeUpdateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID: binding.Field{
			Form:         "id",
			Required:     true,
			ErrorMessage: "请填写id",
		},
		&o.AuthorID: binding.Field{
			Form:         "author_id",
			Required:     true,
			ErrorMessage: "请填写作者id",
		},
	}
}

//Validate 数据格式验证
func (o AttitudeUpdateForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.AuthorID) {
		return binding.Errors{
			binding.NewError([]string{"author_id"}, "format error", "id 格式不正确."),
		}
	}
	fmt.Println("action form id: ", o.ID)

	if !bson.IsObjectIdHex(o.ID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
		}
	}

	return nil
}
