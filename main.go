package main

import (
	"fmt"
	"immense-lowlands-91960/handler"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/huayun321/go-jwt-middleware"
	nigronimgosession "github.com/joeljames/nigroni-mgo-session"
	"github.com/unrolled/render"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"immense-lowlands-91960/middleware"
)

const VERSION_ONE_PREFIX = "/v1"

var (
	port   = os.Getenv("PORT")
	dbURL  = os.Getenv("MONGODB_URI")
	dbName = "heroku_90v42m0v"
	dbColl = "user"
	ren    = render.New(render.Options{IndentJSON: true, StreamingJSON: true, IsDevelopment: true})
)

//===================== jwt start
func jwtOnError(w http.ResponseWriter, r *http.Request, err string) {
	ren.JSON(w, http.StatusUnauthorized, map[string]interface{}{"code": 10001, "msg": "未通过认证的请求", "err": err})
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
			jwtmiddleware.FromJSON("token"),
		),

		Debug: true,
	})

	//nms middleware
	dbAccessor, err := nigronimgosession.NewDatabaseAccessor(dbURL, dbName, dbColl)
	if err != nil {
		panic(err)
	}

	//cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
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
	subRouter.HandleFunc("/user/reset-password", handler.ResetPassword).Methods("POST")
	subRouter.HandleFunc("/user/get-by-id", handler.GetUserByID).Methods("POST")
	subRouter.HandleFunc("/article/list", handler.GetArticles).Methods("POST")
	subRouter.HandleFunc("/article/select", handler.SelectArticle).Methods("POST")
	subRouter.HandleFunc("/article/un-select", handler.UnSelectArticle).Methods("POST")
	subRouter.HandleFunc("/article/delete", handler.DeleteArticle).Methods("POST")
	subRouter.HandleFunc("/news/create", handler.CreateNews).Methods("POST")
	subRouter.HandleFunc("/news/list", handler.GetNews).Methods("POST")
	subRouter.HandleFunc("/news/publish", handler.PublishNews).Methods("POST")
	subRouter.HandleFunc("/news/un-publish", handler.UnPublishNews).Methods("POST")
	subRouter.HandleFunc("/news/update", handler.UpdateNews).Methods("POST")
	subRouter.HandleFunc("/news/delete", handler.DeleteNews).Methods("POST")
	subRouter.HandleFunc("/feedback/list", handler.GetFeedbacks).Methods("POST")
	subRouter.HandleFunc("/feedback/reply", handler.ReplyFeedback).Methods("POST")
	subRouter.HandleFunc("/feedback/track", handler.TrackFeedback).Methods("POST")
	subRouter.HandleFunc("/upload", handler.Upload).Methods("POST")
	subRouter.HandleFunc("/action/create", handler.CreateAction).Methods("POST")
	subRouter.HandleFunc("/action/list", handler.GetActions).Methods("POST")
	subRouter.HandleFunc("/action/delete", handler.DeleteAction).Methods("POST")
	subRouter.HandleFunc("/action/update", handler.UpdateAction).Methods("POST")
	subRouter.HandleFunc("/action/get-by-id", handler.GetActionByID).Methods("POST")
	subRouter.HandleFunc("/attitude/create", handler.CreateAttitude).Methods("POST")
	subRouter.HandleFunc("/attitude/list", handler.GetAttitudes).Methods("POST")
	subRouter.HandleFunc("/attitude/delete", handler.DeleteAttitude).Methods("POST")
	subRouter.HandleFunc("/attitude/update", handler.UpdateAttitude).Methods("POST")
	subRouter.HandleFunc("/attitude/get-by-id", handler.GetAttitudeByID).Methods("POST")
	subRouter.HandleFunc("/plan/create", handler.CreatePlan).Methods("POST")
	subRouter.HandleFunc("/plan/list", handler.GetPlans).Methods("POST")


	router.PathPrefix(VERSION_ONE_PREFIX + "/admin").Handler(negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.HandlerFunc(middleware.IsAdminM),
		negroni.Wrap(subRouter),
	))
	//client
	clientRouter := mux.NewRouter().PathPrefix(VERSION_ONE_PREFIX + "/client").Subrouter().StrictSlash(true)
	clientRouter.HandleFunc("/article/create", handler.CreateArticle).Methods("POST")
	clientRouter.HandleFunc("/article/like", handler.LikeArticle).Methods("POST")
	clientRouter.HandleFunc("/article/unlike", handler.UnLikeArticle).Methods("POST")
	clientRouter.HandleFunc("/article/add-bookmark", handler.AddBookmark).Methods("POST")
	clientRouter.HandleFunc("/article/un-bookmark", handler.UnBookmark).Methods("POST")
	clientRouter.HandleFunc("/article/add-comment", handler.CreateComment).Methods("POST")
	clientRouter.HandleFunc("/article/list", handler.GetArticles).Methods("POST")
	clientRouter.HandleFunc("/article/add-view", handler.AddView).Methods("POST")
	clientRouter.HandleFunc("/article/get-bookmarks", handler.GetBookmarks).Methods("POST")
	clientRouter.HandleFunc("/article/get-article-by-id", handler.GetArticleByID).Methods("POST")
	clientRouter.HandleFunc("/news/get-news-by-id", handler.GetNewsByID).Methods("POST")
	clientRouter.HandleFunc("/news/list", handler.GetNews).Methods("POST")
	clientRouter.HandleFunc("/news/add-comment", handler.CreateNComment).Methods("POST")
	clientRouter.HandleFunc("/feedback/create", handler.CreateFeedback).Methods("POST")
	clientRouter.HandleFunc("/feedback/list-by-user-id", handler.GetFeedbacks).Methods("POST")
	clientRouter.HandleFunc("/upload", handler.Upload).Methods("POST")
	clientRouter.HandleFunc("/user/follow", handler.Follow).Methods("POST")
	clientRouter.HandleFunc("/user/unfollow", handler.UnFollow).Methods("POST")

	router.PathPrefix(VERSION_ONE_PREFIX + "/client").Handler(negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.HandlerFunc(middleware.IsUserM),
		negroni.Wrap(clientRouter),
	))

	//guest
	guestRouter := mux.NewRouter().PathPrefix(VERSION_ONE_PREFIX + "/guest").Subrouter().StrictSlash(true)
	guestRouter.HandleFunc("/upload", handler.Upload).Methods("POST")
	guestRouter.HandleFunc("/article/list", handler.GetArticles).Methods("POST")
	guestRouter.HandleFunc("/article/add-view", handler.AddView).Methods("POST")
	guestRouter.HandleFunc("/article/get-bookmarks", handler.GetBookmarks).Methods("POST")
	guestRouter.HandleFunc("/article/get-article-by-id", handler.GetArticleByID).Methods("POST")
	guestRouter.HandleFunc("/news/get-news-by-id", handler.GetNewsByID).Methods("POST")
	guestRouter.HandleFunc("/news/list", handler.GetNews).Methods("POST")

	router.PathPrefix(VERSION_ONE_PREFIX + "/guest").Handler(negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(guestRouter),
	))

	router.HandleFunc(VERSION_ONE_PREFIX+"/user/signin-phone", handler.SignInWithPhone).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX+"/user/signup-phone", handler.SignUpWithPhone).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX+"/user/sign-wx", handler.SignWithWx).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX+"/user/verify", handler.GetVerifyCode).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX+"/user/set-admin", handler.SetAdmin).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX+"/user/forgot", handler.ForgotPassword).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX+"/news/list", handler.GetNews).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX+"/article/list", handler.GetArticles).Methods("POST")
	router.HandleFunc(VERSION_ONE_PREFIX+"/user/signin-guest", handler.SignInGuest).Methods("POST")
	n.Use(nigronimgosession.NewDatabase(*dbAccessor).Middleware())
	n.Use(c)
	n.UseHandler(router)
	n.Run(":" + port)
}
