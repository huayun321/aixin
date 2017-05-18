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

	"github.com/dgrijalva/jwt-go"
	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"immense-lowlands-91960/form"
)

const (
	JWTEXP = 60 * 60 * 24 * 30
	SMSURL = "https://limitless-spire-42314.herokuapp.com/sms"
)

func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func jwtSign(phone, openid, role string, exp int64) (string, error) {
	// Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone":  phone,
		"openid": openid,
		"role":   role,
		"exp":    exp,
	})

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
	sq := model.SMSQuery{Phone: phone, Code: code}
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
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 101, "message": "手机号格式错误"})
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
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 110, "message": "GetVerifyCode 查询数据库时遇到内部错误"})
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
			util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 110, "message": "GetVerifyCode 插入数据库时遇到内部错误", "vc": vc})
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
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 102, "message": "GetVerifyCode 一分钟后才能发送！"})
		return
	}

	// check if daily
	dsr := now - vc.LastVerifyDay
	fmt.Printf("GetVerifyCode now:%d-- vc.LastVerifyDay: %d = %d \n", now, vc.LastVerifyDay, dsr)
	if dsr < 60*60*24 {
		if vc.TimesRemainDay < 1 {
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 103, "message": "GetVerifyCode 今天的验证码使用次数已经用完，明天再来吧。"})
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
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 110, "message": "插入数据库时遇到内部错误!"})
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
			util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 110, "message": "插入数据库时遇到内部错误!"})
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
	u.Password = uf.Password
	u.Avatar = uf.Avatar
	u.Nickname = uf.Nickname
	u.IsFrozen = false
	u.CreateTime = now

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

//todo add user ad

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

	// check if verify code match
	if swf.Phone != "" && swf.Password != "" && swf.Code != "" {
		fmt.Println("=======SignWithWx 带有电话密码验证码字段")
		u.Phone = swf.Phone
		u.Password = swf.Password
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

			tk, err := jwtSign(swf.Phone, "", "user", time.Now().Unix()+JWTEXP)
			if err != nil {
				fmt.Println("=======SignWithWx 生成token 遇到错误 err: ", err)
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10310, "message": "生成token遇到错误!", "err": err})
				return
			}

			udbp.OpenID = u.OpenID
			udbp.WxUserInfo = u.WxUserInfo
			result := map[string]interface{}{"code": 0, "message": "SignInWithPhone 注册成功", "token": tk, "user": udbp}
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
		u.CreateTime = time.Now().Unix()
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

	tk, err := jwtSign("", swf.WxOpenID, "user", time.Now().Unix()+JWTEXP)
	if err != nil {
		fmt.Println("=======SignWithWx 生成token 遇到错误 err: ", err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10314, "message": "生成token遇到错误!", "err": err})
		return
	}

	//udb.LastLoginTime = udb.LastLoginTime * 1000
	udb.WxUserInfo = u.WxUserInfo
	result := map[string]interface{}{"code": 0, "message": "SignInWithPhone 注册成功", "token": tk, "user": udb}
	fmt.Println("=======SignWithWx 成功返回: result", result)
	util.Ren.JSON(w, http.StatusOK, result)
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
		Key: []string{"phone"},
		Unique: true,
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
		Key: []string{"phone"},
		Unique: true,
		DropDups: true,
	}
	err = nms.DB.C("user").EnsureIndex(iup)
	if err != nil {
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "声明索引时遇到错误"})
		return
	}
	iuo := mgo.Index{
		Key: []string{"openid"},
		Unique: true,
		DropDups: true,
	}
	err = nms.DB.C("user").EnsureIndex(iuo)
	if err != nil {
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "声明索引时遇到错误"})
		return
	}
	util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 0, "message": "声明索引成功"})
}
