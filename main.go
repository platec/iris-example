package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"time"
)

func main() {
	app := iris.New()

	app.RegisterView(iris.HTML("./templates", ".html"))
	app.StaticWeb("/static", "./static")

	// 主页
	app.Get("/", func(ctx iris.Context) {
		sid := ctx.GetCookie("sid")
		if sid != "" {
			ctx.Gzip(true)
			ctx.View("index.html")
		} else {
			ctx.Redirect("/login")
		}
	})

	// 登录页
	app.Get("/login", func(ctx iris.Context) {
		ctx.Gzip(true)
		ctx.View("login.html")
	})

	app.Post("/logout", func(ctx iris.Context) {
		ctx.RemoveCookie("sid")
		ctx.Redirect("/login")
	})

	app.Post("/login", func(ctx iris.Context) {
		username := ctx.FormValue("username")
		password := ctx.FormValue("password")
		if username == "test" && password == "12345678" {
			ctx.SetCookieKV("sid", username, context.CookieExpires(time.Duration(2) * time.Minute))
			ctx.Redirect("/")
		} else {
			ctx.Redirect("/login")
		}
	})

	app.Run(iris.Addr(":8080"))
}