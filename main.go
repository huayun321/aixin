package main

import (
	"fmt"
	"immense-lowlands-91960/handler"
	"net/http"
	"os"

	jwtmiddleware "github.com/huayun321/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	"github.com/unrolled/render"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"immense-lowlands-91960/middleware"
)

const VERSION_ONE_PREFIX  = "/v1"

var (
	port   = os.Getenv("PORT")
	dbURL  = os.Getenv("MONGODB_URI")
	dbName = "heroku_90v42m0v"
	dbColl = "user"
	ren    = render.New(render.Options{IndentJSON: true, StreamingJSON: true, IsDevelopment: true})
)


//===================== jwt start
func jwtOnError(w http.ResponseWriter, r *http.Request, err string) {
	ren.JSON(w, http.StatusUnauthorized, map[string]interface{}{"code": 10001, "msg": "未通过认证的请求", "err":
	err})
}

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

		Extractor: jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader,
				     jwtmiddleware.FromParameter("token"),
                                     jwtmiddleware.FromJSON("token")),

		Debug:true,
	})

	//nms middleware
	dbAccessor, err := nigronimgosession.NewDatabaseAccessor(dbURL, dbName, dbColl)
	if err != nil {
		panic(err)
	}

	//cors
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials:true,
		AllowedHeaders:[]string{"*"},
	})

	router := mux.NewRouter()
	router.HandleFunc("/hello", helloHandler).Methods("GET")
	//admin
	subRouter := mux.NewRouter().PathPrefix(VERSION_ONE_PREFIX + "/admin").Subrouter().StrictSlash(true)
	subRouter.HandleFunc("/user/unfroze", handler.UnFrozeUser).Methods("POST")
	subRouter.HandleFunc("/user/froze", handler.FrozeUser).Methods("POST")
	subRouter.HandleFunc("/user/list", handler.GetUsers).Methods("POST")
	subRouter.HandleFunc("/user/index", handler.EnsureIndex).Methods("GET")
	subRouter.HandleFunc("/user/drop", handler.DropUser).Methods("POST")
	subRouter.HandleFunc("/user/code/drop", handler.DropCode).Methods("POST")
	subRouter.HandleFunc("/user/get-by-id", handler.GetUserByID).Methods("POST")
	router.PathPrefix(VERSION_ONE_PREFIX + "/admin").Handler(negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.HandlerFunc(middleware.IsAdminM),
		negroni.Wrap(subRouter),
	))
	//client
	clientRouter := mux.NewRouter().PathPrefix(VERSION_ONE_PREFIX + "/client").Subrouter().StrictSlash(true)
	clientRouter.HandleFunc("/article/create", handler.CreateArticle).Methods("POST")
	router.PathPrefix(VERSION_ONE_PREFIX + "/client").Handler(negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(clientRouter),
	))

	router.HandleFunc(VERSION_ONE_PREFIX + "/user/signin-phone", handler.SignInWithPhone).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX + "/user/signup-phone", handler.SignUpWithPhone).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX + "/user/sign-wx", handler.SignWithWx).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX + "/user/verify", handler.GetVerifyCode).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX + "/user/set-admin", handler.SetAdmin).Methods("POST")

	n.Use(nigronimgosession.NewDatabase(*dbAccessor).Middleware())
	n.Use(c)
	n.UseHandler(router)
	n.Run(":" + port)
}
