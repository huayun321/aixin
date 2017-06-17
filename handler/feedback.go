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
)

//CreateFeedback 添加文章
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

	a := model.Feedback{}
	a.ID = bson.NewObjectId()
	a.Content = f.Content
	a.Phone = f.Phone
	a.IsTracked = false
	a.CreateTime = time.Now().Unix()
	a.AuthorId = bson.ObjectIdHex(uid)

	//store to db
	err := nms.DB.C("feedback").Insert(a)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 17002, "message":
		"插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}

