package handler

import (
	"github.com/disintegration/imaging"
	"github.com/satori/go.uuid"
	"image"
	"immense-lowlands-91960/util"
	"net/http"
)

const url = "http://immense-lowlands-91960.herokuapp.com/upload/images/"
const dir = "./public/upload/images/"
const suffix = ".jpg"

//test
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseMultipartForm(100000000)
		if err != nil {
			util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 20001, "message": "图片过大", "err": err.Error()})
			return
		}

		//get a ref to the parsed multipart form
		m := r.MultipartForm

		//image names
		inames := []string{}
		//get the *fileheaders
		files := m.File["myfiles"]
		for i, _ := range files {

			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 20002, "message": "无法读取temp 文件", "err": err.Error()})
				return
			}
			img, _, err := image.Decode(file)
			if err != nil {
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 20003, "message": "无法解码图片文件", "err": err.Error()})
				return
			}

			u1 := uuid.NewV4()

			n := dir + u1.String() + suffix
			iname := url + u1.String() + suffix

			err = imaging.Save(img, n)
			if err != nil {
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 20004, "message": "无法保存图片文件：", "err": err.Error()})
				return
			}
			inames = append(inames, iname)
		}
		//display success message.
		util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "上传成功", "urls": inames})
		return
	} else {
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 20000, "message": "只支持post"})
		return
	}
}


func Upload2(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseMultipartForm(100000000)
		if err != nil {
			util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 20001, "message": "图片过大", "err": err.Error()})
			return
		}

		//get the *fileheaders
		file, _, err := r.FormFile("uploadfile")

			defer file.Close()
			if err != nil {
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 20002, "message": "无法读取temp 文件", "err": err.Error()})
				return
			}
			img, _, err := image.Decode(file)
			if err != nil {
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 20003, "message": "无法解码图片文件", "err": err.Error()})
				return
			}

			u1 := uuid.NewV4()

			n := dir + u1.String() + suffix
			iname := url + u1.String() + suffix

			err = imaging.Save(img, n)
			if err != nil {
				util.Ren.JSON(w, http.StatusInternalServerError, map[string]interface{}{"code": 20004, "message": "无法保存图片文件：", "err": err.Error()})
				return
			}


		//display success message.
		util.Ren.JSON(w, http.StatusOK, map[string]interface{}{"code": 0, "message": "上传成功", "url": iname})
		return
	} else {
		util.Ren.JSON(w, http.StatusBadRequest, map[string]interface{}{"code": 20000, "message": "只支持post"})
		return
	}
}
