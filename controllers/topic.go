package controllers

import (
    "github.com/astaxie/beego"
    "beeblog/models"
    "fmt"
    "strings"
    "path"
)

type TopicController struct {
    beego.Controller
}

func (this *TopicController) Get() {
    this.Data["IsLogin"] = checkAccount(this.Ctx)
    this.Data["IsTopic"] = true
    this.TplName = "topic.html"
    topics,err:= models.GetAllTopics("","",false)
    if err != nil {
        beego.Error(err) 
    } else {
        this.Data["Topics"] = topics
    }
}

//使用自动路由访问这个方法
//1,首先需要把需要路由的控制器注册到自动路由中：
//beego.AutoRouter(&controllers.TopicController{})
//2,那么 beego 就会通过反射获取该结构体中所有的实现方法，你就可以通过如下的方式访问到对应的方法中：
//(url:)topic/add   调用 TopicController 中的 add 方法
func (this *TopicController) Add() {
    this.TplName = "topic_add.html"
    // this.Ctx.WriteString("add")
}

func (this *TopicController) Post() {
    if !checkAccount(this.Ctx) {
        //验证不通过就进行重定向，即返回到login页面，继续进行登录验证
        this.Redirect("/login", 302)
        return
    }
    //解析表单
    //对用户提交的表单进行一个解析
    title := this.Input().Get("title")
    content := this.Input().Get("content")
    category := this.Input().Get("category")
    label := this.Input().Get("label")
    //如果有tid代表是修改文章，如果没有tid代表是新增文章
    tid := this.Input().Get("tid")


    //获取附件
    _,fh,err := this.GetFile("attachment")
    if err != nil {
        beego.Error(err)
    }
    var attchment string
    if fh != nil {
        //保存附件
        attchment = fh.Filename
        //把附件信息打印一下
        beego.Info(attchment)
        //假设attachment是tmp.go 则path.Join之后就是attachment/tmp.go
        //这段语句执行之后就会吧tmp.go保存到当前beeblog。exe所在目录的attachment中
        err = this.SaveToFile("attachment", path.Join("attachment",attchment))
        if err != nil {
            beego.Error(err)
        }
    }

    if len(tid) == 0{
        err = models.AddTopic(title,category,label, content, attchment)
    }else{
        err = models.ModifyTopic(tid,title,category,label, content,attchment)
    }
    
    if err != nil {
    beego.Error(err)
    }
    this.Redirect("/topic", 302)

}

func (this *TopicController) View() {
    this.TplName = "topic_view.html"
    topic,err := models.GetTopic(this.Ctx.Input.Params()["0"])
    if err != nil {
        fmt.Println("[********View]"+"models.GetTopic failed!")
        fmt.Println("[********View]"+"this.Ctx.Input.Param:"+ this.Ctx.Input.Params()["0"])
        beego.Error(err)
        this.Redirect("/", 302)
        return 
    }
    this.Data["Topic"] = topic
    this.Data["Labels"] = strings.Split(topic.Labels, " ")

    this.Data["Tid"] = this.Ctx.Input.Params()["0"] 

    tid := this.Ctx.Input.Params()["0"]
    // 获取tda或者用如下方法：
    // 获取的当前页面的url
    // reqUrl := this.Ctx.Request.RequestURI
    // 将url从后往前找最后一个"/"以及之后的字符串
    // i := strings.LastIndex(reqUrl, "/")
     // 取"/"之后的字符串，
    // tid := reqUrl[i+1:]
    replies,err := models.GetAllReplies(tid)
    if err != nil {
        beego.Error(err)
        return 
    }
    this.Data["Replies"] = replies 
    this.Data["IsLogin"] = checkAccount(this.Ctx)


}

func (this *TopicController) Modify() {
    this.TplName = "topic_modify.html"
    tid := this.Input().Get("tid")
    topic, err := models.GetTopic(tid)
    if err != nil {
        beego.Error(err)
        this.Redirect("/", 302)
        return 
    }
    this.Data["Topic"] = topic
    this.Data["Tid"] = tid
}

func (this *TopicController) Delete() {
     if !checkAccount(this.Ctx) {
        //验证不通过就进行重定向，即返回到login页面，继续进行登录验证
        this.Redirect("/login", 302)
        return
    }
    err := models.DeleteTopic(this.Input().Get("tid"))
    if err != nil {
        beego.Error(err)
    }
    this.Redirect("/", 302)

}