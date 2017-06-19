package handler

import (
	"bytes"
	"fmt"
	"immense-lowlands-91960/model"
	"immense-lowlands-91960/util"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"encoding/json"

	"crypto/md5"
	"github.com/dgrijalva/jwt-go"
	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"immense-lowlands-91960/form"
	"io"
	"strconv"
)

const (
	JWTEXP = 60 * 60 * 24 * 30
	SMSURL = "https://limitless-spire-42314.herokuapp.com/sms"
)

func getRandomString(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func jwtSign(id, nickname, role string, exp int64) (string, error) {
	// Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"nickname": nickname,
		"role":     role,
		"exp":      exp,
	})

	fmt.Println("jwtsin token: ", token)

	// Headers
	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"

	//sign
	tokenString, err := token.SignedString([]byte("My Secret"))
	if err != nil {
		fmt.Printf("token err: %v \n", err)
		return "", err
	}

	return tokenString, nil
}

func sendSMS(code, phone string) error {
	p, err :=strconv.Atoi(phone)
	if err != nil {
		fmt.Println("sendSMS str to int err:", err)
		return err
	}
	sq := model.SMSQuery{Phone: p, Code: code}
	jsq, err := json.Marshal(sq)
	if err != nil {
		fmt.Println("sendSMS Marshal err:", err)
		return err
	}

	req, err := http.NewRequest("POST", SMSURL, bytes.NewBuffer(jsq))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("sendSMS NewRequest err:", err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("post api err :", err)
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}

//GetVerifyCode 获得手机验证码
func GetVerifyCode(w http.ResponseWriter, r *http.Request) {
	// check phone num
	vcf := new(model.VerifyCode)
	if errs := binding.Bind(r, vcf); errs != nil {
		fmt.Println("GetVerifyCode: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10101, "message": "手机号格式错误",
			"err": errs})
		return
	}
	//generate code
	code := getRandomString(6)

	//fmt.Fprintf(w, "From:    %s\n", vcf.Phone)
	// check if now - verify timestamp less than 60second
	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	vc := model.VerifyCode{}
	err := nms.DB.C("verifycode").Find(bson.M{"phone": vcf.Phone}).One(&vc)
	fmt.Println("GetVerifyCode vc finded : ", vc)
	// got err
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("GetVerifyCode err:", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10102, "message": "GetVerifyCode 查询数据库时遇到内部错误"})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("GetVerifyCode err:", err)
		now := time.Now().Unix()
		vc.Phone = vcf.Phone
		vc.VerifyTimestamp = now
		vc.LastVerifyDay = now
		vc.TimesRemainDay = 4
		vc.VerifyCode = code
		//store to db
		err := nms.DB.C("verifycode").Insert(&vc)
		if err != nil {
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10103,
				"message": "GetVerifyCode 插入数据库时遇到内部错误", "err": err})
			return
		}
		//todo send sms
		go sendSMS(code, vcf.Phone)

		util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "verify_code": code})
		return
	}

	// check if less than 60 seconds
	now := time.Now().Unix()
	tsr := now - vc.VerifyTimestamp
	fmt.Printf("GetVerifyCode now:%d-- vc.VerifyTimestamp: %d = %d \n", now, vc.VerifyTimestamp, tsr)
	if tsr < 60 {
		fmt.Println("一分钟后才能发送")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10104, "message": "GetVerifyCode 一分钟后才能发送！"})
		return
	}

	// check if daily
	dsr := now - vc.LastVerifyDay
	fmt.Printf("GetVerifyCode now:%d-- vc.LastVerifyDay: %d = %d \n", now, vc.LastVerifyDay, dsr)
	if dsr < 60*60*24 {
		if vc.TimesRemainDay < 1 {
			fmt.Println("今天已经发送了5条，不能发送了")
			util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10105, "message": "GetVerifyCode 今天的验证码使用次数已经用完，明天再来吧。"})
			return
		}
		//时间是同一天，次数有剩余
		change := mgo.Change{
			Update:    bson.M{"$inc": bson.M{"times_remain_day": -1}, "$set": bson.M{"verify_timestamp": now, "verify_code": code}},
			ReturnNew: false,
		}
		vcr := model.VerifyCode{}
		_, err := nms.DB.C("verifycode").Find(bson.M{"phone": vcf.Phone}).Apply(change, &vcr)
		if err != nil {
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10106,
				"message": "插入数据库时遇到内部错误!"})
			return
		}
	} else {
		//时间不是同一天，剩余次数和日期重制
		vc.TimesRemainDay = 4
		change := mgo.Change{
			Update:    bson.M{"$inc": bson.M{"times_remain_day": 4}, "$set": bson.M{"last_verify_day": now, "verify_timestamp": now, "verify_code": code}},
			ReturnNew: false,
		}
		vcr := model.VerifyCode{}
		_, err := nms.DB.C("verifycode").Find(bson.M{"phone": vcf.Phone}).Apply(change, &vcr)
		if err != nil {
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10107,
				"message": "插入数据库时遇到内部错误!"})
			return
		}
	}
	//todo send sms
	go sendSMS(code, vcf.Phone)
	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "verify_code": code})
	return
}

//SignUpWithPhone 通过手机号注册
func SignUpWithPhone(w http.ResponseWriter, r *http.Request) {
	// check params
	uf := new(form.SignUpPhoneForm)

	if errs := binding.Bind(r, uf); errs != nil {
		fmt.Println("SignInWithPhone: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10201, "message": "用户数据格式错误", "err": errs})
		return
	}

	// check if verify code match
	vc := model.VerifyCode{}
	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	err := nms.DB.C("verifycode").Find(bson.M{"phone": uf.Phone}).One(&vc)
	// got err
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("SignInWithPhone err:", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10202, "message": "SignInWithPhone 查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("SignInWithPhone err:", err)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10203, "message": "SignInWithPhone 数据库中并没有此电话的验证码", "err": err})
		return
	}

	if uf.Code != vc.VerifyCode {
		fmt.Println("SignInWithPhone 验证码与存储的验证码不匹配:")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10204, "message": "SignInWithPhone 验证码与存储的验证码不匹配"})
		return
	}

	// check timestamp so late
	// 验证码超过1小时为过期
	now := time.Now().Unix()
	st := now - vc.VerifyTimestamp
	if st > 60*60 {
		fmt.Println("SignInWithPhone 验证码已超过一小时:", err)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10205, "message": "SignInWithPhone 验证码已超过一小时"})
		return
	}

	udb := model.User{}
	err = nms.DB.C("user").Find(bson.M{"phone": uf.Phone}).One(&udb)
	// got err
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("SignInWithPhone err:", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10206, "message": "SignInWithPhone 查询数据库时遇到内部错误", "err": err})
		return
	}

	if err == nil {
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10208, "message": "SignInWithPhone 手机号已被注册"})
		return
	}

	u := model.User{}
	u.Phone = uf.Phone
	//u.Password = uf.Password
	u.Avatar = uf.Avatar
	u.Nickname = uf.Nickname
	u.IsFrozen = false
	u.CreateTime = now

	h := md5.New()
	io.WriteString(h, uf.Password)
	u.Password = fmt.Sprintf("%x", h.Sum(nil))

	//store to db
	err = nms.DB.C("user").Insert(u)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10206, "message": "SignInWithPhone 插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "SignInWithPhone 注册成功"})
	return
}

//SignInWithPhone 通过手机号登陆
func SignInWithPhone(w http.ResponseWriter, r *http.Request) {
	// check params
	uf := new(form.SignInPhoneForm)

	if errs := binding.Bind(r, uf); errs != nil {
		fmt.Println("SignInWithPhone: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10401, "message": "用户数据格式错误", "err": errs})
		return
	}

	h := md5.New()
	io.WriteString(h, uf.Password)
	uf.Password = fmt.Sprintf("%x", h.Sum(nil))

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	//验证用户名密码
	udb := model.User{}
	err := nms.DB.C("user").Find(bson.M{"phone": uf.Phone, "password": uf.Password}).One(&udb)
	// got err
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("SignInWithPhone err:", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10402, "message": "SignInWithPhone 查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("SignInWithPhone 用户不存在 err:", err)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10403, "message": "用户名密码错误"})
		return
	}

	//check frozen
	if udb.IsFrozen {
		fmt.Println("SignInWithPhone phone user is frozen")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10404, "message": "该用户已被冻结，请联系管理人员"})
		return
	}

	udb.LastLoginTime = time.Now().Unix()
	//更新最后登陆时间
	upsertdata := bson.M{"$set": bson.M{"last_login_time": udb.LastLoginTime}}
	err = nms.DB.C("user").UpdateId(udb.ID, upsertdata)
	if err != nil {
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10405, "message": "插入数据库时遇到内部错误!", "err": err})
		return
	}

	//token
	tk, err := jwtSign(udb.ID.Hex(), udb.Nickname, udb.Role, time.Now().Unix()+JWTEXP)
	if err != nil {
		fmt.Println("=======SignWithWx 生成token 遇到错误 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10406, "message": "生成token遇到错误!", "err": err})
		return
	}

	udb.Password = ""

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "登陆成功", "user": udb, "token": tk})
	return
}

//SignInWithWx 通过微信登陆或者注册
func SignWithWx(w http.ResponseWriter, r *http.Request) {
	// check params
	swf := new(form.SignWxForm)

	if errs := binding.Bind(r, swf); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10301, "message": "用户数据格式错误", "err": errs})
		return
	}

	fmt.Println("=======SignWithWx 处理开始")
	fmt.Println("=======SignWithWx 参数:", swf)

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	fmt.Println("=======SignWithWx 获得nms")

	//define user
	u := model.User{}
	u.OpenID = swf.WxOpenID
	u.WxUserInfo.OpenID = swf.WxOpenID
	u.WxUserInfo.Nickname = swf.WxNickname
	u.WxUserInfo.City = swf.WxCity
	u.WxUserInfo.Province = swf.WxProvince
	u.WxUserInfo.Country = swf.WxCountry
	u.WxUserInfo.Headimgurl = swf.WxHeadimgurl
	u.WxUserInfo.Sex = swf.WxSex
	u.WxUserInfo.Unionid = swf.WxUnionid
	u.LastLoginTime = time.Now().Unix()
	u.IsFrozen = false

	// check if verify code match
	if swf.Phone != "" && swf.Password != "" && swf.Code != "" {
		fmt.Println("=======SignWithWx 带有电话密码验证码字段")
		u.Phone = swf.Phone
		h := md5.New()
		io.WriteString(h, swf.Password)
		u.Password = fmt.Sprintf("%x", h.Sum(nil))

		fmt.Println("=======SignWithWx 查询验证码")
		vc := model.VerifyCode{}
		err := nms.DB.C("verifycode").Find(bson.M{"phone": swf.Phone}).One(&vc)
		// got err
		if err != nil && err != mgo.ErrNotFound {
			fmt.Println("=======SignWithWx 查询验证码 失败err：", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10302, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}

		if err != nil && err == mgo.ErrNotFound {
			fmt.Println("=======SignWithWx 查询验证码 验证码不存在：", err)
			util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10303, "message": "数据库中并没有此电话的验证码", "err": err})
			return
		}

		if swf.Code != vc.VerifyCode {
			fmt.Println("=======SignWithWx 查询验证码 验证码与存储的验证码不匹配")
			util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10304, "message": "验证码与存储的验证码不匹配"})
			return
		}

		// check timestamp so late
		// 验证码超过1小时为过期
		now := time.Now().Unix()
		st := now - vc.VerifyTimestamp
		if st > 60*60 {
			fmt.Println("=======SignWithWx 验证码已超过一小时：")
			util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10305, "message": "验证码已超过一小时"})
			return
		}

		//check if phone has signed
		udbp := model.User{}
		fmt.Println("=======SignWithWx 检查电话是否已被注册：")
		fmt.Println("=======SignWithWx 如果已注册则upsert by phone：")

		err = nms.DB.C("user").Find(bson.M{"phone": swf.Phone}).One(&udbp)
		if err != nil && err != mgo.ErrNotFound {
			fmt.Println("=======SignWithWx 查询电话遇到错误 err：", err)
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10306, "message": "查询数据库时遇到内部错误", "err": err})
			return
		}

		if err != nil && err == mgo.ErrNotFound {
			fmt.Println("=======SignWithWx 未查到电话 ")
		}

		fmt.Println("=======SignWithWx 查到电话 phone：", udbp)

		if udbp.IsFrozen {
			fmt.Println("SignWithWx phone user is frozen")
			util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10308, "message": "该用户已被冻结，请联系管理人员"})
			return
		}

		////允许微信登陆更改密码
		//u.Phone = swf.Phone
		//u.Password = swf.Password

		if err == nil {
			//check if openid exist
			//更换微信号
			//添加微信号
			if udbp.OpenID != swf.WxOpenID {
				fmt.Printf("=======SignWithWx 微信号不相同，更改微信号：old %v new %v \n", udbp.OpenID, swf.WxOpenID)
				fmt.Println("=======SignWithWx 微信号不相同，删除所提交微信号的 已注册数据：")
				err = nms.DB.C("user").Remove(bson.M{"openid": swf.WxOpenID})
				if err != nil && err != mgo.ErrNotFound {
					fmt.Println("=======SignWithWx 微信号不相同，删除旧微信号数据 遇到错误err：", err)
					util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10307, "message": "查询数据库时遇到内部错误", "err": err})
					return
				}
			}
			fmt.Println("=======SignWithWx upsert by phone old: ", udbp)
			fmt.Println("=======SignWithWx upsert by phone new: ", u)

			//update
			upsertdata := bson.M{"$set": u}
			_, err := nms.DB.C("user").Upsert(bson.M{"phone": u.Phone}, upsertdata)
			if err != nil {
				fmt.Println("=======SignWithWx upsert by phone 遇到错误 err: ", err)
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10309, "message": "插入数据库时遇到内部错误!"})
				return
			}

			tk, err := jwtSign(udbp.ID.Hex(), udbp.Nickname, udbp.Role, time.Now().Unix()+JWTEXP)
			if err != nil {
				fmt.Println("=======SignWithWx 生成token 遇到错误 err: ", err)
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10310, "message": "生成token遇到错误!", "err": err})
				return
			}

			udbp.OpenID = u.OpenID
			udbp.WxUserInfo = u.WxUserInfo
			result := map[string]interface{}{"code": 0, "message": "操作成功", "token": tk, "user": udbp}
			fmt.Println("=======SignWithWx 成功返回: result", result)
			util.Ren.JSON(w, http.StatusOK, result)
			return
		}
	}

	//todo 最后应该把各流程画个流程图
	fmt.Println("=======SignWithWx 通过openid检测 是否已经注册: ")

	//check if is sign
	udb := model.User{}

	err := nms.DB.C("user").Find(bson.M{"openid": swf.WxOpenID}).One(&udb)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======SignWithWx 通过openid检测 是否已经注册 遇到错误 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10311, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err == nil {
		fmt.Println("=======SignWithWx 通过openid检测 是否已经注册 已注册 user: ", udb)
		if udb.IsFrozen {
			fmt.Println("=======SignWithWx 通过openid检测 是否已经注册 用户已冻结")
			util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10312, "message": "该用户已被冻结，请联系管理人员"})
			return
		}
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======SignWithWx 通过openid检测 是否已经注册 未注册")
		now := time.Now().Unix()
		u.CreateTime = now
		u.ID = bson.NewObjectId()
		fmt.Println("=======SignWithWx  新用户id: ", u.ID)
		udb = u
	}

	//upsert to db
	fmt.Println("=======SignWithWx upsert by openid")
	upsertdata := bson.M{"$set": u}
	_, err = nms.DB.C("user").Upsert(bson.M{"openid": u.OpenID}, upsertdata)
	if err != nil {
		fmt.Println("=======SignWithWx upsert by openid err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10313, "message": "插入数据库时遇到内部错误!"})
		return
	}

	tk, err := jwtSign(udb.ID.Hex(), udb.Nickname, udb.Role, time.Now().Unix()+JWTEXP)
	if err != nil {
		fmt.Println("=======SignWithWx 生成token 遇到错误 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10314, "message": "生成token遇到错误!", "err": err})
		return
	}

	//udb.LastLoginTime = udb.LastLoginTime * 1000
	udb.WxUserInfo = u.WxUserInfo
	result := map[string]interface{}{"code": 0, "message": "操作成功", "token": tk, "user": udb}
	fmt.Println("=======SignWithWx 成功返回: result", result)
	util.Ren.JSON(w, http.StatusOK, result)
	return
}

//GetUsers 获取所有用户
func GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUsers start")
	// check params
	f := new(form.UserListForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10501, "message": "数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("=======SignWithWx 获得nms")

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

	if f.Sex != 0 {
		q["sex"] = f.Sex
	}

	if f.Nickname != "" {
		q["nickname"] = f.Nickname
	}

	if f.IsFrozen {
		q["is_frozen"] = f.IsFrozen
	}

	l := []model.User{}
	err := nms.DB.C("user").Find(q).Sort("create_time").Skip((page - 1) * pageSize).Limit(pageSize).All(&l)
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======GetUsers 获取用户列表 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10502, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======GetUsers 获取用户列表 not found user: ")
	}

	c, err := nms.DB.C("user").Find(q).Count()
	if err != nil {
		fmt.Println("=======GetUsers 获取用户列表数 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10503, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	for i := range l {
		l[i].Password = ""
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "result": l, "total": c})
	return
}

//FrozeUser 冻结用户
func FrozeUser(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.FrozeForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10601, "message": "数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("=======SignWithWx 获得nms")

	upsertdata := bson.M{"$set": bson.M{"is_frozen": true, "froze_time": time.Now().Unix(), "froze_reason": f.Reason}}
	err := nms.DB.C("user").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======FrozeUser update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10602, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======FrozeUser not found user: ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10603, "message": "不存在此条数据", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//SetAdmin 设置为管理员
func SetAdmin(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.SetAdminForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SetAdmin: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 19001, "message": "数据格式错误",
			"err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("=======SetAdmin 获得nms")

	upsertdata := bson.M{"$set": bson.M{"role": "admin"}}
	err := nms.DB.C("user").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======SetAdmin update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 19002, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======SetAdmin not found user: ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 19003, "message": "不存在此条数据",
			"err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//UnFrozeUser 解冻用户
func UnFrozeUser(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.UnFrozeForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10701, "message": "数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("=======SignWithWx 获得nms")

	upsertdata := bson.M{"$set": bson.M{"is_frozen": false, "un_froze_time": time.Now().Unix()}}
	err := nms.DB.C("user").UpdateId(bson.ObjectIdHex(f.ID), upsertdata)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======FrozeUser update err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10702, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======FrozeUser not found user: ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10703, "message": "不存在此条数据", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//EnsureIndex 声明索引
func EnsureIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	//err := nms.DB.C("verifycode").DropIndex("phone")
	//if err != nil {
	//	fmt.Println("drop verifycode index phone err:", err)
	//	util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "删除索引遇到错误"})
	//	return
	//}

	ivc := mgo.Index{
		Key:      []string{"phone"},
		Unique:   true,
		DropDups: true,
	}
	err := nms.DB.C("verifycode").EnsureIndex(ivc)
	if err != nil {
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "声明索引时遇到错误"})
		return
	}

	//user index
	//err = nms.DB.C("user").DropIndex("phone")
	//if err != nil {
	//	fmt.Println("drop user index phone err:", err)
	//
	//	util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "删除索引遇到错误"})
	//	return
	//}
	//
	//err = nms.DB.C("user").DropIndex("openid")
	//if err != nil {
	//	fmt.Println("drop user index openid err:", err)
	//	util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "删除索引遇到错误"})
	//	return
	//}

	iup := mgo.Index{
		Key:      []string{"phone", "openid"},
		Unique:   true,
		DropDups: true,
	}
	err = nms.DB.C("user").EnsureIndex(iup)
	if err != nil {
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "声明索引时遇到错误"})
		return
	}
	iuo := mgo.Index{
		Key:      []string{"openid"},
		Unique:   true,
		DropDups: true,
	}
	err = nms.DB.C("user").EnsureIndex(iuo)
	if err != nil {
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "声明索引时遇到错误"})
		return
	}

	iupw := mgo.Index{
		Key: []string{"phone", "password"},
	}
	err = nms.DB.C("user").EnsureIndex(iupw)
	if err != nil {
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "声明索引时遇到错误"})
		return
	}

	util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 0, "message": "声明索引成功"})
}

//DropUser 删除user collection
func DropUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	err := nms.DB.C("user").DropCollection()
	if err != nil {
		fmt.Println("DropUser: err :", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 19101, "message": "删除user库时遇到错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
}

func DropCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	err := nms.DB.C("verifycode").DropCollection()
	if err != nil {
		fmt.Println("DropCode: err :", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 19201, "message": "删除库时遇到错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
}

//getUserByID 通过id获取用户信息 admin
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.GetUserByIDForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("SignWithWx: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10801, "message": "数据格式错误", "err": errs})
		return
	}

	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	fmt.Println("=======SignWithWx 获得nms")

	udb := model.User{}
	err := nms.DB.C("user").Find(bson.M{"_id": bson.ObjectIdHex(f.ID)}).One(&udb)

	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("=======GetUserByID  err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10802, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("=======GetUserByID not found user: ")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10803, "message": "不存在此条数据", "err": err})
		return
	}

	udb.Password = ""

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功", "user": udb})
	return
}

//ForgotPassword
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// check params
	f := new(form.ForgotPasswordForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("ResetPassword: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10901, "message": "用户数据格式错误", "err": errs})
		return
	}

	// check if verify code match
	vc := model.VerifyCode{}
	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	err := nms.DB.C("verifycode").Find(bson.M{"phone": f.Phone}).One(&vc)
	// got err
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("ResetPassword err:", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10902, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("ResetPassword err:", err)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10903, "message": "数据库中并没有此电话的验证码", "err": err})
		return
	}

	if f.Code != vc.VerifyCode {
		fmt.Println("ResetPassword 验证码与存储的验证码不匹配:")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10904, "message": "验证码与存储的验证码不匹配"})
		return
	}

	// check timestamp so late
	// 验证码超过1小时为过期
	now := time.Now().Unix()
	st := now - vc.VerifyTimestamp
	if st > 60*60 {
		fmt.Println("ResetPassword 验证码已超过一小时:", err)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10905, "message": "验证码已超过一小时"})
		return
	}

	udb := model.User{}
	err = nms.DB.C("user").Find(bson.M{"phone": f.Phone}).One(&udb)
	// got err
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("ResetPassword err:", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10906, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10907, "message": "手机号不存在"})
		return
	}

	//check frozen
	if udb.IsFrozen {
		fmt.Println("phone user is frozen")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10908, "message": "该用户已被冻结，请联系管理人员"})
		return
	}

	//reset
	h := md5.New()
	io.WriteString(h, f.Password)
	np := fmt.Sprintf("%x", h.Sum(nil))

	//store to db
	upsertdata := bson.M{"$set": bson.M{"password": np}}
	err = nms.DB.C("user").UpdateId(udb.ID, upsertdata)

	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10909, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}

//ResetPassword
func ResetPassword(w http.ResponseWriter, r *http.Request) {

	// check params
	f := new(form.ResetPasswordForm)

	if errs := binding.Bind(r, f); errs != nil {
		fmt.Println("ResetPassword: bind err: ", errs)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 11001, "message": "用户数据格式错误", "err": errs})
		return
	}

	user := r.Context().Value("user")
	uid := user.(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	u := model.User{}
	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	hp := md5.New()
	io.WriteString(hp, f.PasswordOld)
	op := fmt.Sprintf("%x", hp.Sum(nil))

	err := nms.DB.C("user").Find(bson.M{"_id": bson.ObjectIdHex(uid), "password": op}).One(&u)
	// got err
	if err != nil && err != mgo.ErrNotFound {
		fmt.Println("ResetPassword err:", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 11002, "message": "查询数据库时遇到内部错误", "err": err})
		return
	}

	if err != nil && err == mgo.ErrNotFound {
		fmt.Println("ResetPassword err:", err)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 11003, "message": "数据库中并没有此用户,或旧密码不匹配", "err": err})
		return
	}

	if u.IsFrozen {
		fmt.Println("phone user is frozen")
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 11004, "message": "该用户已被冻结，请联系管理人员"})
		return
	}

	//reset
	h := md5.New()
	io.WriteString(h, f.PasswordNew)
	np := fmt.Sprintf("%x", h.Sum(nil))

	//store to db
	upsertdata := bson.M{"$set": bson.M{"password": np}}
	err = nms.DB.C("user").UpdateId(bson.ObjectIdHex(uid), upsertdata)

	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 11005, "message": "插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "操作成功"})
	return
}
