package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/basicauth"
	"net/http"
	"time"
)

func main() {
	app := iris.New()

	authConfig := basicauth.Config{
		Users:   map[string]string{"test": "12345678"},
		Realm:   "Authorization Required", // defaults to "Authorization Required"
		Expires: time.Duration(2) * time.Minute,
	}

	authentication := basicauth.New(authConfig)

	needAuth := app.Party("/", authentication)
	{
		needAuth.Get("/", func(ctx iris.Context) {
			ctx.Gzip(true)
			ctx.View("index.html")
		})
	}

	app.RegisterView(iris.HTML("./templates", ".html"))
	app.StaticWeb("/static", "./static")

	// 主页
	//app.Get("/", func(ctx iris.Context) {
	//	ctx.Gzip(true)
	//	ctx.View("index.html")
	//})

	// 登录页
	app.Get("/login", func(ctx iris.Context) {
		ctx.Gzip(true)
		ctx.View("login.html")
	})

	app.Post("/login", func(ctx iris.Context) {
		mobile := ctx.FormValue("mobile")
		password := ctx.FormValue("password")
		fmt.Println(mobile, password)

		cookie := http.Cookie{"sid", "f"}
		ctx.SetCookie()

		ctx.JSON(iris.Map{
			"status": "ok",
			"message": "",
		})
	})

	app.Get("/ping", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "pong",
		})
	})
	app.Run(iris.Addr(":8080"))
}