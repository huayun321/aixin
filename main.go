package main

import (
	"fmt"
	"net/http"
	"os"

	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	"github.com/unrolled/render"

	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/urfave/negroni"
)

var (
	port   = os.Getenv("PORT")
	dbURL  = os.Getenv("MONGODB_URI")
	dbName = "heroku_90v42m0v"
	dbColl = "user"
	ren    = render.New(render.Options{IndentJSON: true, StreamingJSON: true, IsDevelopment: true})
)

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
func signInHandler(resp http.ResponseWriter, req *http.Request) {
	lf := new(LoginForm)
	if errs := binding.Bind(req, lf); errs != nil {
		http.Error(resp, errs.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(resp, "From:    %d\n", lf.Username)
	fmt.Fprintf(resp, "Message: %s\n", lf.Password)

	ctx := req.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)
	nms.DB.C(dbColl).Insert(&lf)
}

func usersHandler(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	nms := ctx.Value(nigronimgosession.KEY).(*nigronimgosession.NMS)

	list := []LoginForm{}
	nms.DB.C(dbColl).Find(nil).All(&list)
	ren.JSON(resp, http.StatusOK, struct {
		Users []LoginForm `json:"users"`
	}{list})
}

//test
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
	fmt.Fprintln(w, "hello from heroku")
}

func main() {
	n := negroni.Classic()

	dbAccessor, err := nigronimgosession.NewDatabaseAccessor(dbURL, dbName, dbColl)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler).Methods("GET")
	router.HandleFunc("/binding", bindingHandler).Methods("POST")
	router.HandleFunc("/sign-in", signInHandler).Methods("POST")
	router.HandleFunc("/users", usersHandler).Methods("GET")

	n.Use(nigronimgosession.NewDatabase(*dbAccessor).Middleware())
	n.UseHandler(router)
	n.Run(":" + port)
}
