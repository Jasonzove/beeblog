package controllers

import (
	"github.com/astaxie/beego"
    "beeblog/models"
    "fmt"
    "strconv"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
    fmt.Println("[*****mHomeController]:Get()")
    this.Data["IsHome"] = true
	this.TplName = "home.html"
    this.Data["IsLogin"] = checkAccount(this.Ctx)
    topics,err:= models.GetAllTopics(this.Input().Get("cate"),this.Input().Get("label"),true)
    if err != nil {
        beego.Error(err) 
    } 
    
    this.Data["Topics"] = topics
    categories, err := models.GetAllCategories()
    fmt.Println("[*****categories length]:" + strconv.Itoa(len(categories)))

    for _,data := range categories{
        fmt.Println("[*****models.GetAllCategories]:")
        fmt.Println(*data)
    }
    if err != nil {
        fmt.Println("[*****models.GetAllTopics]:err")
    beego.Error(err)
    }

    this.Data["Categories"] = categories

}
