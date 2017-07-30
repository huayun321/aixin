package form

import (
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"fmt"
	"immense-lowlands-91960/model"
)

//==============================================================文章创建表单

//ActionCreateForm 动作创建表单
type ActionCreateForm struct {
	AuthorID string   `json:"author_id"` //user object id
	Name     string   `json:"name"`
	Level    string   `json:"level"`
	Symptom  string   `json:"symptom,omitempty"` //适应症状
	People   string   `json:"people,omitempty"`
	Notice   string   `json:"notice,omitempty"`
	MainImg  string   `json:"main_img,omitempty"`
	Subs  []model.Sub `json:"subs,omitempty"`
	Images   []string `json:"images"`
}

// FieldMap 数据绑定
func (o *ActionCreateForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.AuthorID: binding.Field{
			Form:         "author_id",
			Required:     true,
			ErrorMessage: "请填写作者id",
		},
		&o.Name: binding.Field{
			Form:         "name",
			Required:     true,
			ErrorMessage: "请填写动作名称",
		},
		&o.Level: binding.Field{
			Form:         "level",
			Required:     true,
			ErrorMessage: "请填写动作难度",
		},
		&o.Symptom: binding.Field{
			Form:         "symptom",
			Required:     true,
			ErrorMessage: "请填写适应症状",
		},
		&o.People: binding.Field{
			Form:         "people",
			Required:     true,
			ErrorMessage: "请填写适应人群",
		},
		&o.Notice: binding.Field{
			Form:         "notice",
			Required:     true,
			ErrorMessage: "请填写注意事项",
		},
		&o.MainImg: binding.Field{
			Form:         "main_img",
			Required:     true,
			ErrorMessage: "请填写动作主图",
		},
		&o.Subs: binding.Field{
			Form:         "subs",
			Required:     true,
			ErrorMessage: "请填写动作分解图",
		},
	}
}

//Validate 数据格式验证
func (o ActionCreateForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.AuthorID) {
		return binding.Errors{
			binding.NewError([]string{"author_id"}, "format error", "id 格式不正确."),
		}
	}

	return nil
}

//==============================================================获取文章列表表单

//ActionListForm
type ActionListForm struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Name     string `json:"name"`
	Level    string `json:"level"`
}

// FieldMap 数据绑定
func (o *ActionListForm) FieldMap(req *http.Request) binding.FieldMap {
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
		&o.Level: binding.Field{
			Form: "level",
		},
	}
}

//Validate 数据格式验证
func (o ActionListForm) Validate(req *http.Request) error {
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

//ActionIdForm
type ActionIdForm struct {
	ID string `json:"id"`
}

// FieldMap 数据绑定
func (o *ActionIdForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&o.ID: binding.Field{
			Form:         "id",
			Required:     true,
			ErrorMessage: "请提交文章id",
		},
	}
}

//Validate 数据格式验证
func (o ActionIdForm) Validate(req *http.Request) error {
	if !bson.IsObjectIdHex(o.ID) {
		return binding.Errors{
			binding.NewError([]string{"id"}, "format error", "id 格式不正确."),
		}
	}
	return nil
}

//ActionUpdateForm 动作表单
type ActionUpdateForm struct {
	ID       string   `json:"id"`
	AuthorID string   `json:"author_id"` //user object id
	Name     string   `json:"name"`
	Level    string   `json:"level"`
	Symptom  string   `json:"symptom,omitempty"` //适应症状
	People   string   `json:"people,omitempty"`
	Notice   string   `json:"notice,omitempty"`
	MainImg  string   `json:"main_img,omitempty"`
	StepImg  []string `json:"step_img,omitempty"`
	Key      string   `json:"key,omitempty"` //关键点
	Images   []string `json:"images"`
}

// FieldMap 数据绑定
func (o *ActionUpdateForm) FieldMap(req *http.Request) binding.FieldMap {
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
func (o ActionUpdateForm) Validate(req *http.Request) error {
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
