package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

//CreateNews 添加资讯
func CreateNews(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.NewsCreateForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("CreateArticle: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16001, "message": "用户数据格式错误",
			"err": errs})
		return
	}

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	a := model.News{}
	a.ID = bson.NewObjectId()
	a.Title = f.Title
	a.Content = f.Content
	a.Image = f.Image
	a.Position = f.Position
	a.IsPublished = false
	a.CreateTime = time.Now().Unix()
	a.AuthorId = bson.ObjectIdHex(uid)

	//store to db
	err := nms.DB.C("news").Insert(a)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16002, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}

//CreateNComment
func CreateNComment(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.NCommentCreateForm)
	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("CreateArticle: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16101, "message": "用户数据格式错误",
			"err": errs})
		return
	}

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	var u = model.User{}
	err := nms.DB.C("user").Find(bson.M{"_id": bson.ObjectIdHex(uid)}).One(&u)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16102, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}
	u.Password = ""

	a := model.NComment{}
	a.ID = bson.NewObjectId()
	a.Content = f.Content
	a.AuthorID = bson.ObjectIdHex(uid)
	a.NewsID = bson.ObjectIdHex(f.NewsID)
	a.Author = u
	a.CreateTime = time.Now().Unix()

	//store to db
	err = nms.DB.C("ncomment").Insert(a)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16103, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}

//GetNewsByID
func GetNewsByID(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.NewsIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16201, "message": "数据格式错误",
			"err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	a := model.News{}
	err := nms.DB.C("news").Find(bson.M{"_id": bson.ObjectIdHex(f.ID)}).One(&a)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16202, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 not found: ")
	}

	//v.Fans = []model.User{}
	var u = model.User{}
	err = nms.DB.C("user").Find(bson.M{"_id": a.AuthorId}).One(&u)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16203, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}
	u.Password = ""
	a.Author = u

	ncs := []model.NComment{}
	err = nms.DB.C("ncomment").Find(bson.M{"news_id": a.ID}).Sort("-create_time").All(&ncs)
	if err != nil {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16204, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}
	a.NComments = ncs

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}

//GetNews
func GetNews(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.NewsListForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16301, "message": "数据格式错误",
			"err": errs})
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

	if f.Title != "" {
		q["title"] = f.Title
	}

	if f.TimeStart > 0 && f.TimeEnd > 0 {
		q["create_time"] = bson.M{"$gte": f.TimeStart, "$lte": f.TimeEnd}
	}

	if f.IsPublished == 1 {
		q["is_published"] = true
	}

	if f.IsPublished == 2 {
		q["is_published"] = false
	}

	l := []model.News{}
	err := nms.DB.C("news").Find(q).Sort("-create_time").Skip((page - 1) * pageSize).Limit(pageSize).All(&l)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16302, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 not found: ")
	}

	c, err := nms.DB.C("news").Find(q).Count()
	if err != nil {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16303, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": l, "total": c})
	return
}

//PublishNews
func PublishNews(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.NewsIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16401, "message": "数据格式错误",
			"err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	upsertdata := bson.M{"$set": bson.M{"is_published": true, "publish_time": time.Now().Unix()}}
	err := nms.DB.C("news").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("======= update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16402, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("======= not found : ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16403, "message": "不存在此条数据",
			"err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//UnPublishNews
func UnPublishNews(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.NewsIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16501, "message": "数据格式错误",
			"err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	upsertdata := bson.M{"$set": bson.M{"is_published": false, "un_publish_time": time.Now().Unix()}}
	err := nms.DB.C("news").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("======= update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16502, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("======= not found : ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16503, "message": "不存在此条数据",
			"err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//UpdateNews
func UpdateNews(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.NewsUpdateForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16601, "message": "数据格式错误",
			"err": errs})
		return
	}

	q := bson.M{}

	if f.Title != "" {
		q["title"] = f.Title
	}
	if f.Content != "" {
		q["content"] = f.Content
	}
	if f.Position != 0 {
		q["position"] = f.Position
	}
	if f.Image != "" {
		q["image"] = f.Image
	}

	fmt.Println("update form:", q)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	upsertdata := bson.M{"$set": q}
	err := nms.DB.C("news").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

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

//DeleteNews
func DeleteNews(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.NewsIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 16701, "message": "数据格式错误",
			"err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	err := nms.DB.C("news").RemoveId(bson.ObjectIdHex(f.ID))

	if err != nil {
		fmt.Println("======= update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16702, "message": "删除数据时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}
