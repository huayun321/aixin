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