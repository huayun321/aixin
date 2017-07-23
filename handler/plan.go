package handler

import (
	"net/http"
	"github.com/joeljames/nigroni-mgo-session"
	"immense-lowlands-91960/model"
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
	"immense-lowlands-91960/form"
	"github.com/mholt/binding"
	"immense-lowlands-91960/util"
	"gopkg.in/mgo.v2"
)

//CreatePlan 添加动作
func CreatePlan(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.PlanCreateForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("CreateArticle: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13101, "message": "用户数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	a := model.Plan{}
	a.ID = bson.NewObjectId()
	a.Name = f.Name
	a.First = f.First
	a.Second = f.Second
	a.F2 = f.F2
	a.F3 = f.F3
	a.Level = f.Level
	a.Feel = f.Feel
	a.Weeks = f.Weeks
	a.IsRecommend = f.IsRecommend
	a.CreateTime = time.Now().Unix()
	a.AuthorId = bson.ObjectIdHex(f.AuthorID)

	//store to db
	err := nms.DB.C("plan").Insert(a)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13102, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}

//GetPlans
func GetPlans(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.PlanListForm)

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

	if f.First != 0 {
		q["first"] = f.First
	}

	if f.IsRecommend == 1 {
		q["is_recommend"] = true
	}

	if f.IsRecommend == 2 {
		q["is_recommend"] = false
	}

	l := []model.Plan{}
	err := nms.DB.C("plan").Find(q).Sort("-create_time").Skip((page - 1) * pageSize).Limit(pageSize).All(&l)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13202, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 not found: ")
	}

	c, err := nms.DB.C("plan").Find(q).Count()
	if err != nil {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13203, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": l, "total": c})
	return
}