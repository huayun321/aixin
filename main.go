// @APIVersion 1.0.0
// @APITitle My Cool API
// @APIDescription My API usually works as expected.
// @Contact huayun321@gmail.com
// @TermsOfServiceUrl http://pandariel.com
// @BasePath http://localhost:3000

package main

import (
	"fmt"
	"immense-lowlands-91960/handler"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	"github.com/unrolled/render"

	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

var (
	port   = os.Getenv("PORT")
	dbURL  = os.Getenv("MONGODB_URI")
	dbName = "heroku_90v42m0v"
	dbColl = "user"
	ren    = render.New(render.Options{IndentJSON: true, StreamingJSON: true, IsDevelopment: true})
)

type ErrorResp struct {
	Code int
	Msg  string
}

//===================== binding
// First define a type to hold the data
// (If the data comes from JSON, see: http://mholt.github.io/json-to-go)
type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Then provide a field mapping (pointer receiver is vital)
func (lf *LoginForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&lf.Username: binding.Field{
			Form:         "username",
			Required:     true,
			ErrorMessage: "数据格式错误，请提交用户名",
		},
		&lf.Password: binding.Field{
			Form:         "password",
			Required:     true,
			ErrorMessage: "数据格式错误，请提交密码",
		},
	}
}

func (lf LoginForm) Validate(req *http.Request) error {
	if len(lf.Username) < 6 {
		return binding.Errors{
			binding.NewError([]string{"message"}, "LengthError", "用户名不能少于6个字符."),
		}
	}
	return nil
}

// Now your handlers can stay clean and simple
func bindingHandler(resp http.ResponseWriter, req *http.Request) {
	lf := new(LoginForm)
	if errs := binding.Bind(req, lf); errs != nil {
		http.Error(resp, errs.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(resp, "From:    %d\n", lf.Username)
	fmt.Fprintf(resp, "Message: %s\n", lf.Password)
}

//===================== binding end

//===================== nms start

// @Title signIn
// @Description login and get token
// @Accept  json
// @Param   username        query   string  true        "username should not less than 6 characters"
// @Param   password        query   string  true        "password should not less than 6 characters"
// @Router /sign-in [post]
func signInHandler(resp http.ResponseWriter, req *http.Request) {
	lf := new(LoginForm)
	if errs := binding.Bind(req, lf); errs != nil {
		ren.JSON(resp, http.StatusBadRequest, ErrorResp{0001, "bad req"})
		return
	}
	fmt.Fprintf(resp, "From:    %d\n", lf.Username)
	fmt.Fprintf(resp, "Message: %s\n", lf.Password)

	ctx := req.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	nms.DB.C(dbColl).Insert(&lf)
}

// @Title get users
// @Description get all users
// @Router /users [get]
func usersHandler(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	list := []LoginForm{}
	nms.DB.C(dbColl).Find(nil).All(&list)
	ren.JSON(resp, http.StatusOK, struct {
		Users []LoginForm `json:"users"`
	}{list})
}

//===================== nms end

//===================== jwt start
func jwtOnError(w http.ResponseWriter, r *http.Request, err string) {
	ren.JSON(w, http.StatusUnauthorized, map[string]interface{}{"code": 001, "msg": "未通过认证的请求", "err": "unauthorized"})
}

func jwtLoginHandler(resp http.ResponseWriter, req *http.Request) {
	lf := new(LoginForm)
	if errs := binding.Bind(req, lf); errs != nil {
		http.Error(resp, errs.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(lf)

	//check username password
	if lf.Username != "admin123" || lf.Password != "admin123" {
		ren.JSON(resp, http.StatusBadRequest, map[string]interface{}{"code": 001, "msg": "用户名或密码错误", "err": "username password not match"})
		return
	}

	// Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "asdf",
		"admin":   "true",
	})

	// Headers
	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"

	//sign
	tokenString, err := token.SignedString([]byte("My Secret"))
	if err != nil {
		fmt.Fprintf(resp, "token err: %v", err)
		return
	}

	err = ren.JSON(resp, http.StatusOK, map[string]interface{}{"code": 0, "msg": "ok", "token": tokenString})
	if err != nil {
		fmt.Fprint(resp, "something wrong")
	}
}

func jwtSecuredHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Println(user)
	ren.JSON(w, http.StatusOK, "All good. You only get this message if you're authenticated")
}

//===================== jwt end
//test
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
	fmt.Fprintln(w, "hello from heroku")
}

func main() {
	n := negroni.Classic()
	//jwt middleware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,

		ErrorHandler: jwtOnError,
	})

	//nms middleware
	dbAccessor, err := nigronimgosession.NewDatabaseAccessor(dbURL, dbName, dbColl)
	if err != nil {
		panic(err)
	}

	//cors
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler).Methods("GET")
	router.HandleFunc("/binding", bindingHandler).Methods("POST")
	router.HandleFunc("/sign-in", signInHandler).Methods("POST")
	router.HandleFunc("/users", usersHandler).Methods("GET")
	router.HandleFunc("/login", jwtLoginHandler).Methods("POST")
	router.Handle("/secured", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(jwtSecuredHandler)),
	))
	router.HandleFunc("/user/sign-in", handler.SignInWithPhone).Methods("POST")
	router.HandleFunc("/user/signup-phone", handler.SignUpWithPhone).Methods("POST")
	router.HandleFunc("/user/sign-wx", handler.SignWithWx).Methods("POST")
	// router.HandleFunc("/user/signin-wx", handler.SignInWithPhone).Methods("POST")
	router.HandleFunc("/user/verify", handler.GetVerifyCode).Methods("POST")
	router.HandleFunc("/user/index", handler.EnsureIndex).Methods("GET")

	n.Use(nigronimgosession.NewDatabase(*dbAccessor).Middleware())
	n.Use(c)
	n.UseHandler(router)
	n.Run(":" + port)
}
