package form

import (
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
	"immense-lowlands-91960/model"
	"net/http"
)

//==============================================================文章创建表单

//PlanCreateForm 态度创建表单
type PlanCreateForm struct {
	AuthorID    string       `json:"author_id"` //user object id
	Name        string       `json:"name"`
	First       model.Part       `json:"first"`
	Second      model.Part       `json:"second"`
	F2          model.Part       `json:"f2"`
	F3          model.Part       `json:"f3"`
	Level       string       `json:"level"`
	Feel        string       `json:"feel"`
	Weeks       []model.Week `json:"weeks"`
	IsRecommend bool         `json:"is_recommend"`
	Desc        string       `json:"desc"`
	Img        string       `json:"img"`
}

// FieldMap 数据绑定
func (o *PlanCreateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.AuthorID: binding.Field{
			Form:         "author_id",
			Required:     true,
			ErrorMessage: "请填写作者id",
		},
		&o.Name: binding.Field{
			Form:         "name",
			Required:     true,
			ErrorMessage: "请填写计划名称",
		},
		&o.First: binding.Field{
			Form:         "first",
			ErrorMessage: "请填写计划一级部位",
		},
		&o.Weeks: binding.Field{
			Form:         "weeks",
			ErrorMessage: "请填写计划周",
		},
		&o.Img: binding.Field{
			Form:         "img",
			ErrorMessage: "请填写计划周",
		},
		&o.Desc: binding.Field{
			Form:         "desc",
			ErrorMessage: "请填写计划周",
		},

	}
}

//Validate 数据格式验证
func (o PlanCreateForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.AuthorID) {
		return binding.Errors{
			binding.NewError([]string{"author_id"}, "format error", "id 格式不正确."),
		}
	}

	return nil
}


//==============================================================获取文章列表表单

//PlanListForm
type PlanListForm struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Name     string `json:"name"`
	First     model.Part `json:"first"`
	IsRecommend int `json:"is_recommend"`
}

// FieldMap 数据绑定
func (o *PlanListForm) FieldMap(req *http.Request) binding.FieldMap {
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
		&o.IsRecommend: binding.Field{
			Form: "is_recommend",
		},
		&o.First: binding.Field{
			Form: "first",
		},
	}
}

//Validate 数据格式验证
func (o PlanListForm) Validate(req *http.Request) error {
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
type PlanIdForm struct {
	ID string `json:"id"`
}

// FieldMap 数据绑定
func (o *PlanIdForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID: binding.Field{
			Form:         "id",
			Required:     true,
			ErrorMessage: "请提交文章id",
		},
	}
}

//Validate 数据格式验证
func (o PlanIdForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.ID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
		}
	}
	return nil
}


