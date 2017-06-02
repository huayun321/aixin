package handler

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"immense-lowlands-91960/form"
	"immense-lowlands-91960/model"
	"immense-lowlands-91960/util"
	"net/http"
	"time"
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

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
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
	err := nms.DB.C("article").Find(q).Sort("-create_time").Skip((page - 1) * pageSize).Limit(pageSize).All(&l)
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

	fmt.Println("articles before find fans", l)

	for i, v := range l {
		fmt.Println("v.fans", v.Fans)
		//v.Fans = []model.User{}
		var u = model.User{}
		err = nms.DB.C("article").Find(bson.M{"_id": v.AuthorId}).One(&u)
		if err != nil && err != mgo.ErrNotFound {
			fmt.Println("=======获取文章列表 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13204, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		u.Password = ""
		l[i].Author = u

		fc, err := nms.DB.C("fan").Find(bson.M{"article_id": v.ID}).Count()
		if err != nil {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13205, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		l[i].FansCount = fc

		cc, err := nms.DB.C("comment").Find(bson.M{"article_id": v.ID}).Count()
		if err != nil {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13206, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		l[i].CommentsCount = cc

		bc, err := nms.DB.C("bookmark").Find(bson.M{"article_id": v.ID}).Count()
		if err != nil {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13207, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		l[i].BookmarkCount = bc

		vc, err := nms.DB.C("view").Find(bson.M{"article_id": v.ID}).Count()
		if err != nil {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13208, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		l[i].ViewCount = vc

		fl := []model.Fan{}
		err = nms.DB.C("fan").Find(bson.M{"article_id": v.ID}).All(&fl)
		if err != nil && err != mgo.ErrNotFound {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13209, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}

		for _, fu := range fl {
			fuu := model.User{}
			err = nms.DB.C("user").FindId(fu.UserID).One(&fuu)
			if err != nil && err != mgo.ErrNotFound {
				fmt.Println("=======获取文章列表数 err: ", err)
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13210, "message": "查询数据库时遇到内部错误", "err": err})
				return
			}
			fmt.Println("v.fans fuu:", fuu)
			l[i].Fans = append(l[i].Fans, fuu)
		}

		fmt.Println("v.fans after find fans", v.Fans)

	}
	fmt.Println("articles after find fans", l)

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

	upsertdata := bson.M{"$set": bson.M{"is_selected": true, "select_time": time.Now().Unix()}}
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

	upsertdata := bson.M{"$set": bson.M{"is_selected": false, "un_select_time": time.Now().Unix()}}
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

	upsertdata := bson.M{"$set": bson.M{"is_deleted": false, "delete_time": time.Now().Unix()}}
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

//LikeArticle
func LikeArticle(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13601, "message": "数据格式错误", "err": errs})
		return
	}

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	a := model.Fan{}
	a.ID = bson.NewObjectId()
	a.ArticleID = bson.ObjectIdHex(f.ID)
	a.UserID = bson.ObjectIdHex(uid)
	a.CreateTime = time.Now().Unix()

	//store to db
	upsertdata := bson.M{"$set": a}
	_, err := nms.DB.C("fan").Upsert(bson.M{"article_id": a.ArticleID, "user_id": a.UserID}, upsertdata)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13602, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//UnLikeArticle
func UnLikeArticle(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13701, "message": "数据格式错误", "err": errs})
		return
	}

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	aid := bson.ObjectIdHex(f.ID)
	ouid := bson.ObjectIdHex(uid)

	//store to db
	err := nms.DB.C("fan").Remove(bson.M{"article_id": aid, "user_id": ouid})
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13702, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//AddView
func AddView(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("CreateArticle: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13801, "message": "用户数据格式错误", "err": errs})
		return
	}

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	a := model.View{}
	a.ID = bson.NewObjectId()
	a.ArticleID = bson.ObjectIdHex(f.ID)
	a.UserID = bson.ObjectIdHex(uid)
	a.CreateTime = time.Now().Unix()

	//store to db
	err := nms.DB.C("view").Insert(a)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13802, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}

//AddBookmark
func AddBookmark(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 13901, "message": "数据格式错误", "err": errs})
		return
	}

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	a := model.Bookmark{}
	a.ID = bson.NewObjectId()
	a.ArticleID = bson.ObjectIdHex(f.ID)
	a.UserID = bson.ObjectIdHex(uid)
	a.CreateTime = time.Now().Unix()

	//store to db
	upsertdata := bson.M{"$set": a}
	_, err := nms.DB.C("bookmark").Upsert(bson.M{"article_id": a.ArticleID, "user_id": a.UserID}, upsertdata)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 13902, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//UnBookmark
func UnBookmark(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 14001, "message": "数据格式错误", "err": errs})
		return
	}

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("======= 获得nms")

	aid := bson.ObjectIdHex(f.ID)
	ouid := bson.ObjectIdHex(uid)

	//store to db
	err := nms.DB.C("bookmark").Remove(bson.M{"article_id": aid, "user_id": ouid})
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14002, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//GetBookmarks
func GetBookmarks(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.BookmarkListForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 14101, "message": "数据格式错误", "err": errs})
		return
	}

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	q := bson.M{}
	q["user_id"] = bson.ObjectIdHex(uid)
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

	l := []model.Bookmark{}
	err := nms.DB.C("bookmark").Find(q).Sort("-create_time").Skip((page - 1) * pageSize).Limit(pageSize).All(&l)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14102, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 not found: ")
	}

	c, err := nms.DB.C("bookmark").Find(q).Count()
	if err != nil {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14103, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	u := model.User{}
	err = nms.DB.C("user").Find(bson.M{"_id": bson.ObjectIdHex(uid)}).One(&u)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14104, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}
	u.Password = ""

	al := []model.Article{}

	for _, v := range l {
		var la = model.Article{}

		la.Author = u

		fc, err := nms.DB.C("fan").Find(bson.M{"article_id": v.ArticleID}).Count()
		if err != nil {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14105, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		la.FansCount = fc

		cc, err := nms.DB.C("comment").Find(bson.M{"article_id": v.ArticleID}).Count()
		if err != nil {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14106, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		la.CommentsCount = cc

		bc, err := nms.DB.C("bookmark").Find(bson.M{"article_id": v.ArticleID}).Count()
		if err != nil {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14107, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		la.BookmarkCount = bc

		vc, err := nms.DB.C("view").Find(bson.M{"article_id": v.ArticleID}).Count()
		if err != nil {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14108, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		la.ViewCount = vc

		fl := []model.Fan{}
		err = nms.DB.C("fan").Find(bson.M{"article_id": v.ArticleID}).All(&fl)
		if err != nil && err != mgo.ErrNotFound {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14109, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}

		for _, fu := range fl {
			fuu := model.User{}
			err = nms.DB.C("user").FindId(fu.UserID).One(&fuu)
			if err != nil && err != mgo.ErrNotFound {
				fmt.Println("=======获取文章列表数 err: ", err)
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14110, "message": "查询数据库时遇到内部错误", "err": err})
				return
			}
			fmt.Println("v.fans fuu:", fuu)
			la.Fans = append(la.Fans, fuu)
		}

		al = append(al, la)
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": al, "total": c})
	return
}

//GetArticles
func GetArticleByID(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ArticleIdForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 14201, "message": "数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	a := model.Article{}
	err := nms.DB.C("article").Find(bson.M{"_id": bson.ObjectIdHex(f.ID)}).One(&a)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14202, "message": "查询数据库时遇到内部错误", "err": err})
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
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14204, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}
	u.Password = ""
	a.Author = u

	fc, err := nms.DB.C("fan").Find(bson.M{"article_id": a.ID}).Count()
	if err != nil {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14205, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}
	a.FansCount = fc

	cc, err := nms.DB.C("comment").Find(bson.M{"article_id": a.ID}).Count()
	if err != nil {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14206, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}
	a.CommentsCount = cc

	bc, err := nms.DB.C("bookmark").Find(bson.M{"article_id": a.ID}).Count()
	if err != nil {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14207, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}
	a.BookmarkCount = bc

	vc, err := nms.DB.C("view").Find(bson.M{"article_id": a.ID}).Count()
	if err != nil {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14208, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}
	a.ViewCount = vc

	fl := []model.Fan{}
	err = nms.DB.C("fan").Find(bson.M{"article_id": a.ID}).All(&fl)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======获取文章列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14209, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	for _, fu := range fl {
		fuu := model.User{}
		err = nms.DB.C("user").FindId(fu.UserID).One(&fuu)
		if err != nil && err != mgo.ErrNotFound {
			fmt.Println("=======获取文章列表数 err: ", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14210, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}
		fmt.Println("v.fans fuu:", fuu)
		a.Fans = append(a.Fans, fuu)
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}

//CreateComment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.CommentCreateForm)
	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("CreateArticle: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 14301, "message": "用户数据格式错误", "err": errs})
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
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14302, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}
	u.Password = ""

	a := model.Comment{}
	a.ID = bson.NewObjectId()
	a.Content = f.Content
	a.AuthorID = bson.ObjectIdHex(uid)
	a.ArticleID = bson.ObjectIdHex(f.ArticleID)
	a.Author = u
	if f.ReferenceID != "" {
		a.ReferenceID = bson.ObjectIdHex(f.ReferenceID)
	}
	a.CreateTime = time.Now().Unix()

	//store to db
	err = nms.DB.C("comment").Insert(a)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 14303, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": a})
	return
}
