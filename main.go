package main

import (
	"fmt"
	"github.com/urfave/negroni"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"github.com/mholt/binding"
)

var port = os.Getenv("PORT")

//===================== binding
// First define a type to hold the data
// (If the data comes from JSON, see: http://mholt.github.io/json-to-go)
type LoginForm struct {
	Username   string
	Password   string
}

// Then provide a field mapping (pointer receiver is vital)
func (lf *LoginForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&lf.Username: binding.Field{
			Form:     "username",
			Required: true,
			ErrorMessage: "数据格式错误，请提交用户名",
		},
		&lf.Password: binding.Field{
			Form:     "password",
			Required: true,
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


//test
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
	fmt.Fprintln(w, "hello from heroku")
}

func main() {
	n := negroni.Classic()

	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler).Methods("GET")
	router.HandleFunc("/binding", bindingHandler).Methods("POST")

	n.UseHandler(router)
	n.Run(":" + port)
}
