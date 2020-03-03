package controllers

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/context"
    "fmt"
)

type LoginController struct {
    beego.Controller
}

func (this *LoginController) Get() {
    fmt.Println("[********]"+"Login_Get：")
    isExist := this.Input().Get("exit") == "true"
    if isExist {
        fmt.Println("[********]"+"Cookie isExist")
        this.Ctx.SetCookie("uname", "", -1, "/")
        this.Ctx.SetCookie("pwd", "", -1, "/")
        this.Redirect("/", 301)
        return 
    }
    this.TplName = "login.html"
}

func (this *LoginController) Post() {
    uname := this.Input().Get("uname")
    pwd := this.Input().Get("pwd")
    fmt.Println("[********]"+"输入登录用户名："+uname + "输入登录密码：" + pwd)
    autoLogin := this.Input().Get("autoLogin") == "on"
    if beego.AppConfig.String("uname") == uname &&
    beego.AppConfig.String("pwd") == pwd {
        maxAge := 0
        if autoLogin {
            //1左移动31位然后减去1，是一个比较大的数
           maxAge = 1<<31 - 1
        }
        this.Ctx.SetCookie("uname", uname, maxAge, "/")
        this.Ctx.SetCookie("pwd", pwd, maxAge, "/")

    }
    this.Redirect("/", 301)
    return 
}

//显示：undefined：beego.Context 新版本beego的api 导致beego.Context找不到。需要
//import "github.com/astaxie/beego/context"
//原：beego.Context 改为： *context.Context
func checkAccount(ctx *context.Context) bool {
    ck,err := ctx.Request.Cookie("uname")
    if err != nil {
    return false
    }
    uname := ck.Value
    ck,err = ctx.Request.Cookie("pwd")
    pwd := ck.Value
    if beego.AppConfig.String("uname") != uname ||
    beego.AppConfig.String("pwd") != pwd {
        fmt.Println("[********checkAccount]"+"配置文件用户名："+beego.AppConfig.String("uname") + "配置文件登录密码：" + beego.AppConfig.String("pwd"))
        fmt.Println("[********checkAccount]"+"登录用户名："+uname + "登录密码：" + pwd)
        fmt.Println("[********checkAccount]"+"登录失败！")
        return false
    }
    return true
}