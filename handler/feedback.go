package handler

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"net/http"
	"github.com/joeljames/nigroni-mgo-session"
	"immense-lowlands-91960/model"
	"immense-lowlands-91960/form"
	"github.com/mholt/binding"
	"immense-lowlands-91960/util"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
)

//CreateFeedback 添加反馈
func CreateFeedback(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.FeedbackCreateForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("CreateFeedback: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17001, "message": "用户数据格式错误",
			"err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)
	fmt.Println(uid)

	u := model.User{}
	err := nms.DB.C("user").FindId(bson.ObjectIdHex(uid)).One(&u)
	// got err
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("ResetPassword err:", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 17002, "message":
		"查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("用户不存在 err:", err)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17003, "message": "数据库中并没有此用户," +
			"用户不存在", "err": err})
		return
	}

	if u.IsFrozen {
		fmt.Println("phone user is frozen")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17004, "message":
		"该用户已被冻结，请联系管理人员"})
		return
	}

	a := model.Feedback{}
	a.ID = bson.NewObjectId()
	a.Content = f.Content
	a.Phone = f.Phone
	a.IsTracked = false
	a.CreateTime = time.Now().Unix()
	a.AuthorId = bson.ObjectIdHex(uid)
	a.Author = u

	//store to db
	err = nms.DB.C("feedback").Insert(a)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 17005, "message":
		"插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}

//GetFeedbacks
func GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.FeedbackListForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17101, "message": "数据格式错误",
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

	if f.NickName != "" {
		q["author.nickname"] = f.NickName
	}

	if f.TimeStart > 0 && f.TimeEnd > 0 {
		q["create_time"] = bson.M{"$gte": f.TimeStart, "$lte": f.TimeEnd}
	}

	if f.IsTracked == 1 {
		q["is_tracked"] = true
	}

	if f.IsTracked == 2 {
		q["is_tracked"] = false
	}

	if f.Phone != "" {
		q["phone"] = f.Phone
	}

	if f.UserID != "" {
		q["author_id"] = bson.ObjectIdHex(f.UserID)
	}

	l := []model.Feedback{}
	err := nms.DB.C("feedback").Find(q).Sort("-create_time").Skip((page - 1) * pageSize).Limit(pageSize).All(&l)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 17102, "message":
		"查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 not found: ")
	}

	c, err := nms.DB.C("feedback").Find(q).Count()
	if err != nil {
		fmt.Println("=======获取列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 17103, "message":
		"查询数据库时遇到内部错误", "err": err})
		return
	}


	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": l, "total": c})
	return
}

//ReplyFeedback
func ReplyFeedback(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ReplyCreateForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17201, "message": "数据格式错误",
			"err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	u := model.User{}
	err := nms.DB.C("user").FindId(bson.ObjectIdHex(uid)).One(&u)
	// got err
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("ResetPassword err:", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 17202, "message":
		"查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("用户不存在 err:", err)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17203, "message": "数据库中并没有此用户," +
			"用户不存在", "err": err})
		return
	}

	if u.IsFrozen {
		fmt.Println("phone user is frozen")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17204, "message":
		"该用户已被冻结，请联系管理人员"})
		return
	}

	rp := model.Reply{}
	rp.Content = f.Content
	rp.CreateTime = time.Now().Unix()
	u.Password = ""
	rp.Author = u
	upsertdata := bson.M{"$set": bson.M{"is_replied": true, "reply": rp}}
	err = nms.DB.C("feedback").UpdateId(bson.ObjectIdHex(f.FeedbackID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("======= update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 17205, "message":
		"插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("======= not found : ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17206, "message": "不存在此条数据",
			"err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//TrackFeedback
func TrackFeedback(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.FeedbackIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17301, "message": "数据格式错误",
			"err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	u := model.User{}
	err := nms.DB.C("user").FindId( bson.ObjectIdHex(uid)).One(&u)
	// got err
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("ResetPassword err:", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 17302, "message":
		"查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("用户不存在 err:", err)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17303, "message": "数据库中并没有此用户," +
			"用户不存在", "err": err})
		return
	}

	if u.IsFrozen {
		fmt.Println("phone user is frozen")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17304, "message":
		"该用户已被冻结，请联系管理人员"})
		return
	}

	upsertdata := bson.M{"$set": bson.M{"is_tracked": true}}
	err = nms.DB.C("feedback").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("======= update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 17305, "message":
		"插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("======= not found : ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 17306, "message": "不存在此条数据",
			"err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

