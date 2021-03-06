package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"os"
	"time"
)

var (
	usernameStored string
	passwordStored string
	duration int64
	sess *sessions.Sessions
)

func SHA256Str(src string) string {
	h := sha256.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}

func control(ctx iris.Context) {
	// TODO
	ctx.Application().Logger().Infof("控制操作来自 %s %s",  ctx.RemoteAddr(), ctx.GetHeader("User-Agent"))
}

func loginRequired(ctx iris.Context) {
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		ctx.JSON(iris.Map{"message": "认证过期，请重新登录", "status": "error"})
	} else {
		ctx.Next()
	}
}

func initApp() *iris.Application {
	app := iris.New()

	app.RegisterView(iris.HTML("./public", ".html"))
	app.StaticWeb("/static", "./public/static")

	// 主页
	app.Get("/", func(ctx iris.Context) {
		if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
			ctx.Redirect("/login")
		} else {
			ctx.View("index.html")
		}
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

	r := app.Party("/api", loginRequired)

	r.Post("/control", control)

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

func newLogFile() *os.File {
	f, err := os.OpenFile("history.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		fmt.Println("读取配置文件失败[config.ini]")
		return
	}
	port, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "port")
	ip, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "ip")
	usernameStored, _ = cfg.GetValue(goconfig.DEFAULT_SECTION, "username")
	passwordStored, _ = cfg.GetValue(goconfig.DEFAULT_SECTION, "password")
	duration, _ = cfg.Int64(goconfig.DEFAULT_SECTION,"duration")

	sess = sessions.New(sessions.Config{Cookie: "sessionid", AllowReclaim: true, Expires: time.Duration(time.Minute * time.Duration(duration))})

	f := newLogFile()
	defer f.Close()
	app := initApp()
	app.Logger().SetOutput(newLogFile())
	app.Run(iris.Addr(ip + ":" + port))
}
