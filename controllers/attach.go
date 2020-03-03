package controllers

import (
    "github.com/astaxie/beego"
    "io"
    "net/url"
    "os"
)

type AttachController struct {
    beego.Controller
}

func (this *AttachController) Get() {
filePath,err:=url.QueryUnescape(this.Ctx.Request.RequestURI[1:])
    if err != nil {
        this.Ctx.WriteString(err.Error())
         return 
    }

     f,err := os.Open(filePath)
    if err != nil {
        this.Ctx.WriteString(err.Error())
         return 
    }
    defer f.Close()

    /*为什么要进行[1:]，因为如果得到的url比如是/topic/view/24，不去掉前面的"/"，则会变成一个绝对路径
    直接获取的url如果是中文，就会被编码，获取的编码就会不相符，则会找不到文件。所以要进行一个反编码
    进行一个反编码，将传入进去的字符串，反编码成为一个正常的字符串*/
    _,err = io.Copy(this.Ctx.ResponseWriter, f)
    if err != nil {
        this.Ctx.WriteString(err.Error())
         return 
    }
}
    
   
    



