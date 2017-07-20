package handler

import (
	"fmt"
	"github.com/joeljames/nigroni-mgo-session"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"immense-lowlands-91960/form"
	"immense-lowlands-91960/model"
	"immense-lowlands-91960/util"
	"net/http"
	"time"
)

//CreateAction 添加动作
func CreateAction(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ActionCreateForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("CreateArticle: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13101, "message": "用户数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	a := model.Action{}
	a.ID = bson.NewObjectId()
	a.Name = f.Name
	a.Level = f.Level
	a.Symptom = f.Symptom
	a.People = f.People
	a.Notice = f.Notice
	a.MainImg = f.MainImg
	a.StepImg = f.StepImg
	a.Key = f.Key
	a.CreateTime = time.Now().Unix()
	a.AuthorId = bson.ObjectIdHex(f.AuthorID)

	//store to db
	err := nms.DB.C("action").Insert(a)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13102, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}

//GetActions
func GetActions(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ActionListForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13201, "message": "数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	q := bson.M{}
	var page int
	var pageSize int
	page = 1
	pageSize = 20

	if f.Page != 0 {
		page = f.Page
	}

	if f.PageSize != 0 {
		pageSize = f.PageSize
	}

	if f.Name != "" {
		q["name"] = f.Name
	}

	if f.Level != "" {
		q["level"] = f.Level
	}

	l := []model.Action{}
	err := nms.DB.C("action").Find(q).Sort("-create_time").Skip((page - 1) * pageSize).Limit(pageSize).All(&l)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13202, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 not found: ")
	}

	c, err := nms.DB.C("action").Find(q).Count()
	if err != nil {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13203, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": l, "total": c})
	return
}

//DeleteAction
func DeleteAction(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ActionIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16701, "message": "数据格式错误",
			"err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	err := nms.DB.C("action").RemoveId(bson.ObjectIdHex(f.ID))

	if err != nil {
		fmt.Println("======= update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16702, "message": "删除数据时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//UpdateAction
func UpdateAction(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ActionUpdateForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16601, "message": "数据格式错误",
			"err": errs})
		return
	}

	q := bson.M{}

	if f.Name != "" {
		q["name"] = f.Name
	}
	if f.Level != "" {
		q["level"] = f.Level
	}
	if f.Symptom != "" {
		q["symptom"] = f.Symptom
	}
	if f.People != "" {
		q["people"] = f.People
	}
	if f.Notice != "" {
		q["notice"] = f.Notice
	}
	if f.MainImg != "" {
		q["main_img"] = f.MainImg
	}
	if len(f.StepImg) != 0 {
		q["step_img"] = f.StepImg
	}
	if f.Key != "" {
		q["key"] = f.Key
	}
	if f.AuthorID != "" {
		q["author_id"] = bson.ObjectIdHex(f.AuthorID)
	}

	fmt.Println("update form:", q)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	upsertdata := bson.M{"$set": q}
	err := nms.DB.C("action").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("======= update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16602, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("======= not found : ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16603, "message": "不存在此条数据",
			"err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}


//GetActionByID
func GetActionByID(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ActionIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16201, "message": "数据格式错误",
			"err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	a := model.Action{}
	err := nms.DB.C("action").Find(bson.M{"_id": bson.ObjectIdHex(f.ID)}).One(&a)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16202, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 not found: ")
	}



	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}
