package handler

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
	"net/http"
	"github.com/mholt/binding"
	"github.com/joeljames/nigroni-mgo-session"
	"immense-lowlands-91960/model"
	"immense-lowlands-91960/form"
	"immense-lowlands-91960/util"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
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
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16002, "message":
		"插入数据库时遇到内部错误", "err": err})
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
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16102, "message":
		"查询数据库时遇到内部错误", "err": err})
		return
	}
	u.Password = ""

	a := model.Comment{}
	a.ID = bson.NewObjectId()
	a.Content = f.Content
	a.AuthorID = bson.ObjectIdHex(uid)
	a.ArticleID = bson.ObjectIdHex(f.ArticleID)
	a.Author = u
	a.CreateTime = time.Now().Unix()

	//store to db
	err = nms.DB.C("ncomment").Insert(a)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 16103, "message":
		"插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}
