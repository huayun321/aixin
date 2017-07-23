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
	First       int       `json:"first"`
	Second      int       `json:"second"`
	F2          int       `json:"f2"`
	F3          int       `json:"f3"`
	Level       string       `json:"level"`
	Feel        string       `json:"feel"`
	Weeks       []model.Week `json:"weeks"`
	IsRecommend bool         `json:"is_recommend"`
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
		&o.Second: binding.Field{
			Form:         "second",
			Required:     true,
			ErrorMessage: "请填写计划二级部位",
		},
		&o.F2: binding.Field{
			Form:         "f2",
			Required:     true,
			ErrorMessage: "请填写计划二级部位f2",
		},
		&o.F3: binding.Field{
			Form:         "f3",
			Required:     true,
			ErrorMessage: "请填写计划二级部位f3",
		},
		&o.Level: binding.Field{
			Form:         "level",
			Required:     true,
			ErrorMessage: "请填写计划锻炼强度",
		},
		&o.Feel: binding.Field{
			Form:         "feel",
			Required:     true,
			ErrorMessage: "请填写计划病情判断",
		},
		&o.Weeks: binding.Field{
			Form:         "weeks",
			ErrorMessage: "请填写计划周",
		},
		&o.IsRecommend: binding.Field{
			Form:         "is_recommend",
			Required:     true,
			ErrorMessage: "请填写计划类型",
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
	First     int `json:"first"`
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
