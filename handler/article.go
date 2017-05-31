package handler

import (
	"net/http"
	"immense-lowlands-91960/model"
	"immense-lowlands-91960/util"
	"fmt"
	"immense-lowlands-91960/form"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
	"time"
	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	"gopkg.in/mgo.v2"
)

//CreateArticle 添加文章
func CreateArticle(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleCreateForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("CreateArticle: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13101, "message": "用户数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	a := model.Article{}
	a.ID = bson.NewObjectId()
	a.Content = f.Content
	a.Images = f.Images
	a.IsSelected = false
	a.IsDeleted = false
	a.CreateTime = time.Now().Unix()
	a.AuthorId = bson.ObjectIdHex(f.AuthorID)

	//store to db
	err := nms.DB.C("article").Insert(a)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13102, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result":a})
	return
}

//GetArticles
func GetArticles(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleListForm)

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

	if f.Phone != "" {
		q["phone"] = f.Phone
	}

	if f.Nickname != "" {
		q["nickname"] = f.Nickname
	}

	if f.IsSelected {
		q["is_selected"] = f.IsSelected
	}

	if f.IsDeleted {
		q["is_deleted"] = f.IsDeleted
	}

	l := []model.Article{}
	err := nms.DB.C("article").Find(q).Sort("create_time").Skip((page - 1) * pageSize).Limit(pageSize).All(&l)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13202, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 not found: ")
	}

	c, err := nms.DB.C("article").Find(q).Count()
	if err != nil {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13203, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	for _, v := range l {
		var u = model.User{}
		err = nms.DB.C("article").Find(bson.M{"_id": v.AuthorId}).One(&u)
		if err != nil && err != mgo.ErrNotFound {
			fmt.Println("=======获取文章列表 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13204, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		v.Author = u
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": l, "total": c})
	return
}


//SelectArticle
func SelectArticle(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13301, "message": "数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	upsertdata := bson.M{"$set": bson.M{"is_selected": true, "select_time":time.Now().Unix()}}
	err := nms.DB.C("article").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("======= update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13302, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("======= not found : ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13303, "message": "不存在此条数据", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}


//UnSelectArticle
func UnSelectArticle(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13401, "message": "数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	upsertdata := bson.M{"$set": bson.M{"is_selected": false, "un_select_time":time.Now().Unix()}}
	err := nms.DB.C("article").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("======= update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13402, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("======= not found : ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13403, "message": "不存在此条数据", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}


//DeleteArticle
func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13501, "message": "数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	upsertdata := bson.M{"$set": bson.M{"is_deleted": false, "delete_time":time.Now().Unix()}}
	err := nms.DB.C("article").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("======= update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13502, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("======= not found : ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13503, "message": "不存在此条数据", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}
