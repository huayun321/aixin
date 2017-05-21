package middleware

import (
	"net/http"
	"immense-lowlands-91960/util"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

//IsAdminM 验证是否为admin
func IsAdminM(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("IsAdminM")
	user := r.Context().Value("user")
	if user.(*jwt.Token).Claims.(jwt.MapClaims)["role"] != "admin" {
		util.Ren.JSON(w, http.StatusUnauthorized, map[string]interface{}{
			"code":10001,
			"message":"该用户没有此权限",
		})
		return
	}

	next(w, r)
}
