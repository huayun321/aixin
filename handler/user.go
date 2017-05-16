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

	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func sendSMS(code, phone string) error {
	url := "https://limitless-spire-42314.herokuapp.com/sms"
	sq := model.SMSQuery{Phone: phone, Code: code}
	jsq, err := json.Marshal(sq)
	if err != nil {
		fmt.Println("sendSMS Marshal err:", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsq))
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

//SignInWithPhone 通过手机号注册
func SignInWithPhone(w http.ResponseWriter, r *http.Request) {
	// check params
	uf := new(model.User)

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
	if st > 60 * 60 {
		fmt.Println("SignInWithPhone 验证码已超过一小时:", err)
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 10205, "message": "SignInWithPhone 验证码已超过一小时"})
		return
	}

	uf.CreateTime = now

	//store to db
	fmt.Println("pre insert uf: ", uf)
	err = nms.DB.C("user").Insert(&uf)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 10206, "message": "SignInWithPhone 插入数据库时遇到内部错误", "err": err})
		return
	}

	util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "SignInWithPhone 注册成功"})
	return
}

//EnsureIndex 声明索引
func EnsureIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	err := nms.DB.C("verifycode").EnsureIndexKey("phone")
	if err != nil {
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "声明索引时遇到错误"})
		return
	}
	err = nms.DB.C("user").EnsureIndexKey("phone", "openid")
	if err != nil {
		util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 120, "message": "声明索引时遇到错误"})
		return
	}
	util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 0, "message": "声明索引成功"})
}
