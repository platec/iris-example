package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"time"
	"crypto/sha256"
	"encoding/hex"
)

var (
	sess = sessions.New(sessions.Config{Cookie: "sessionid", AllowReclaim: true, Expires: time.Duration(time.Minute * 2)})
	usernameStored string
	passwordStored string
)

func SHA256Str(src string) string {
	h := sha256.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}


func initApp() *iris.Application {
	app := iris.New()

	app.RegisterView(iris.HTML("./public", ".html"))
	app.StaticWeb("/static", "./public/static")

	// 主页
	app.Get("/", func(ctx iris.Context) {
		if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
			//ctx.StatusCode(iris.StatusForbidden)
			ctx.Redirect("/login")
			return
		}
		ctx.View("index.html")
	})

	// 登录页
	app.Get("/login", func(ctx iris.Context) {
		ctx.Gzip(true)
		ctx.View("login.html")
	})

	app.Post("/logout", func(ctx iris.Context) {
		session := sess.Start(ctx)
		session.Set("authenticated", false)
		ctx.Redirect("/login")
	})

	app.Post("/login", func(ctx iris.Context) {
		username := ctx.FormValue("username")
		password := ctx.FormValue("password")
		if username == usernameStored && SHA256Str(password) == passwordStored {
			session := sess.Start(ctx)
			session.Set("authenticated", true)
			ctx.JSON(iris.Map{"message": "", "status": "ok"})
		} else {
			ctx.JSON(iris.Map{"message": "用户名或密码错误", "status": "error"})
		}
	})
	return app
}

func main() {
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		fmt.Println("读取配置文件失败")
		return
	}
	port, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "port")
	usernameStored, _ = cfg.GetValue(goconfig.DEFAULT_SECTION, "username")
	passwordStored, _ = cfg.GetValue(goconfig.DEFAULT_SECTION, "password")

	app := initApp()
	app.Run(iris.Addr(":" + port))
}
