package main

import (
	//"beeblog/routers"
    "beeblog/controllers"
    "beeblog/models"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "os"

)
func init() {
    models.RegisterDB()
}

func main() {

    //开发着模式，方便调试
    orm.Debug = true
    //自动建表 必须有一个表的名字为default  第二个参数强制删除重建(true）第三参数是否打印相关信息
    orm.RunSyncdb("default", false, true)

    //首页的路由
    beego.Router("/", &controllers.HomeController{})
    
    //注册分类的路由
    beego.Router("/category", &controllers.CategoryController{})
    beego.Router("/topic", &controllers.TopicController{})
    //beego的自动路由
    beego.AutoRouter(&controllers.TopicController{})
    beego.Router("/reply", &controllers.ReplyController{})
    beego.Router("/reply/add", &controllers.ReplyController{},"post:Add")
    beego.Router("/reply/delete", &controllers.ReplyController{},"get:Delete")
    //登录的路由
    beego.Router("/login", &controllers.LoginController{})
    //创建附件目录
    os.Mkdir("attachment", os.ModePerm)
    //有两种方式一种是作为静态文件来处理 一种是作为控制器来出来
    //第一种方法：作为静态文件处理附件
    //beego.SetStaticPath("/attachment", "attachment")
    //第二种方法 作为控制器来处理
    beego.Router("/attachment/:all", &controllers.AttachController{})

    //启动beego
	beego.Run()

}

