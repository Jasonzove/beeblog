package models

import (
    "os"
    "path"//取出某一个路径的中的目录路径
    "time"
    "strings"
    "github.com/Unknwon/com"
    "github.com/astaxie/beego/orm"
    //一定要注意前面的下滑杠
    _"github.com/mattn/go-sqlite3"
    "fmt"
    "strconv"//转换包
)

const (
    _DB_NAME        = "data/beeblog.db"
    _SQLITE3_DRIVER = "sqlite3"
)

type Category struct {
 Id              int64
 Title           string 
 Created         time.Time `orm:"index"` 
 Views           int64     `orm:"index"`
 TopicTime       time.Time `orm:"index"`
 TopicCount      int64
 TopicLastUserId int64

}

type Topic struct {
    Id              int64
    Uid             int64
    Title           string
    Category        string
    Labels          string
    Content         string      `orm:"size(5000)"`
    Attachment      string 
    Created         time.Time   `orm:"index"` 
    Updated         time.Time   `orm:"index"` 
    Views           int64       `orm:"index"`
    Author          string
    ReplyTime       time.Time   `orm:"index"` 
    ReplyCount      int64
    ReplyLastUserId int64

}

//评论
type Comment struct {
    Id int64
    Tid int64
    Name string
    //beego 的默认设置是255
    Content string `orm:"size(1000)"`
    Created time.Time `orm:"index"`
}


func RegisterDB(){
    if !com.IsExist(_DB_NAME) {
        os.MkdirAll(path.Dir(_DB_NAME),os.ModePerm)
        os.Create(_DB_NAME)
    }

    
    //注册模型
    orm.RegisterModel(new(Category),new(Topic),new(Comment))
    //注册驱动（“sqlite3” 属于默认注册，此处代码可以省略）
    orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)

    //创建一个默认的数据库
    //orm.RegisterDataBase(aliasName, driverName, dataSource, params)
    orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)


}

func AddTopic(title, category,label, content,attchment string) error{
    //处理标签
    /*
    下面语句的意思是 比如传的字符是"bddgo orm" -->分隔后[beego orm] --> 组合之后$beego#$orm#*/
    label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"
    //空格作为过个标签的分隔符
    //beego存到数据库中是$beego# 如果还有一个标签orm 存到数据库变为$beego#$orm# 这样作为标签是唯一的


    
    o := orm.NewOrm()
    fmt.Println("[*****AddTopic]title:"+title + "category:"+ category +"content:" + content )
    topic := &Topic{
        Title: title,
        Category: category,
        Labels: label,
        Content: content,
        Attachment: attchment,
        Created: time.Now(),
        Updated: time.Now(),
        ReplyTime: time.Now(),
    }
    _, err:=o.Insert(topic)
    err = AddCategory(category)
    if err != nil {
    return err
    }

    //更新分类统计
    //获得一个分类的对象
    cate := new(Category)
    //取得分类相关的querytable
    qs := o.QueryTable("category")
    //找到想要操作的分类，根据title找
    err = qs.Filter("title",category).One(cate)
    if err == nil {
        //r如果不存在，简单的忽略更新操作
        cate.TopicCount++
        _,err=o.Update(cate)
    }
    return err
}

func ModifyTopic(tid, title,category,label,content,attchment string) error{
    fmt.Println("[*****ModifyTopic]:tid" + tid)
    tidNum, err := strconv.ParseInt(tid, 10, 64)
    if err != nil {
        fmt.Println("[*****ModifyTopic]strconv.ParseInt failed" )
    return err
    }

    //处理标签
    /*下面语句的意思是 比如传的字符是"bddgo orm" -->分隔后[beego orm] --> 组合之后$beego#$orm#*/
    label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

    var oldCate, oldAattch string
    o:= orm.NewOrm()
    topic := &Topic{Id: tidNum,Created: time.Now(),Updated: time.Now(),ReplyTime: time.Now()}
    if o.Read(topic) == nil {
        fmt.Println("[*****ModifyTopic]o.Read(topic) == nil" )
        //给old赋值 保存起来
        oldCate = topic.Category
        oldAattch = topic.Attachment
        topic.Title = title
        topic.Category = category
        topic.Labels = label
        topic.Content = content
        topic.Attachment = attchment
        topic.Updated = time.Now()
        o.Update(topic)
    }

    //更新旧的分类统计
    if len(oldCate) > 0{
        cate:= new(Category)
        qs := o.QueryTable("category")
        err := qs.Filter("title",oldCate).One(cate)
        //如果找到了
        if err == nil {
            cate.TopicCount--
            _,err = o.Update(cate)
        }
    }

    //删除旧的附件
    if len(oldAattch)>0{
        os.Remove(path.Join("attachment",oldAattch))
    }

    //更新新的分类统计  
    cate := new(Category)
    qs := o.QueryTable("category")
    err = qs.Filter("title",category).One(cate)
    if err == nil {
        cate.TopicCount++
        _,err = o.Update(cate)
    }
    return nil
}

func DeleteTopic(tid string) error{
    tidNum, err := strconv.ParseInt(tid, 10, 64)
    if err != nil {
    return err
    }
    var oldCate string
    o := orm.NewOrm()
    topic := &Topic{Id: tidNum,Created: time.Now(),Updated: time.Now(),ReplyTime:time.Now()}
    if o.Read(topic) == nil{
        oldCate = topic.Category
        _,err = o.Delete(topic)
        if err != nil {
            return err
        }
    }

    if len(oldCate) >0{
        cate := new(Category)
        qs := o.QueryTable("category")
        err = qs.Filter("title",oldCate).One(cate)
        if err == nil {
            cate.TopicCount--
            _,err = o.Update(cate)
        }
    }
    //_,err = o.Delete(topic)
    return err
}

func AddReply(tid,nickname,content string) error {
    tidNum, err := strconv.ParseInt(tid, 10, 64)
    if err != nil {
        return err
    }
    reply := &Comment{
        Tid: tidNum,
        Name: nickname,
        Content: content,
        Created: time.Now(),
    }
    o:= orm.NewOrm()
    _,err = o.Insert(reply)
    if err != nil {
    return err
    }

    topic := &Topic{Id: tidNum, Created: time.Now(),Updated: time.Now(),ReplyTime:time.Now()}
    if o.Read(topic) == nil {
        topic.ReplyTime =time.Now()
        topic.ReplyCount++
        _,err = o.Update(topic)
    }

    return err
    
}

func DeleteReply(rid string) error {
    ridNum, err := strconv.ParseInt(rid, 10, 64)
    if err != nil {
        return err
    }
    o := orm.NewOrm()

    var tidNum int64
    reply := &Comment{
        Id: ridNum,
        Created: time.Now(),
    }
    if o.Read(reply) == nil {
        tidNum = reply.Tid
        _,err = o.Delete(reply)
        if err != nil {
        return err
        }
    }

    replies := make([]*Comment, 0)
    //精确获取最后回复时间，获取所有的回复，然后通过降序排列，取第一个就是最后一个的 回复时间
    qs:= o.QueryTable("comment")
    _,err = qs.Filter("tid",tidNum).OrderBy("-created").All(&replies)
    if err != nil {
    return err
    }
    topic := &Topic{Id: tidNum, Created: time.Now(),Updated: time.Now(),ReplyTime:time.Now()}
     fmt.Println("[*****DeleteReply]topic")
    if o.Read(topic) == nil{
        topic.ReplyTime =replies[0].Created//最后一个精确的创建时间
        topic.ReplyCount = int64(len(replies))
        fmt.Println("[*****DeleteReply]topic.ReplyCount" + strconv.FormatInt(topic.ReplyCount,10) )
        _,err = o.Update(topic)
    }
    return err

}

func GetAllReplies(tid string) (replies []*Comment,err error) {
    tidNum, err := strconv.ParseInt(tid, 10, 64)
    if err != nil {
        return nil, err
    }

    replies = make([]*Comment, 0)
    o := orm.NewOrm()
    qs:= o.QueryTable("comment")
    _, err = qs.Filter("tid",tidNum).All(&replies)
    return replies, err
}

func AddCategory(name string) error {
    o := orm.NewOrm()
    cate := &Category{Title: name,Created: time.Now(),TopicTime: time.Now()}
    qs := o.QueryTable("category")
    err := qs.Filter("title",name).One(cate)
    if err == nil {
        fmt.Println("qs.Filter 失败 err*********")
    return err
    }

    _, err = o.Insert(cate)
    if err != nil {
        fmt.Println("o.Insert 失败 err***********")
    return err
    }
    return nil
}

func DelCategory(id string) error {
    cid, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
    return err
    }
    o := orm.NewOrm()
    cate := &Category{Id: cid,Created: time.Now(),TopicTime: time.Now()}

    _, err = o.Delete(cate)
    return err
}



//參數isDesc是否进行倒序排序。
func GetAllTopics(category,label string,isDesc bool) ([]*Topic,error){
    o := orm.NewOrm()
    topics := make([]*Topic, 0)
    qs := o.QueryTable("topic")

    var err error
    if isDesc {
        if len(category )>0 {
            qs = qs.Filter("category",category)
        }
        if len(label) >0 {
            qs = qs.Filter("labels__contains","$"+label+"#")
        }
        _, err = qs.OrderBy("-created").All(&topics)
    } else {
        _, err = qs.All(&topics)
    }
    
    return topics, err
}

func GetAllCategories()([]*Category,error){
    o := orm.NewOrm()
    cates := make([]*Category, 0)
    qs := o.QueryTable("category")
    _, err := qs.All(&cates)
    return cates, err
}

func GetTopic(tid string) (*Topic, error){
    tidNum, err := strconv.ParseInt(tid, 10, 64)
    if err != nil {
    return nil,err
    }
    o := orm.NewOrm()
    topic := new(Topic)
    qs := o.QueryTable("topic")
    err = qs.Filter("id",tidNum).One(topic)
    if err != nil {
    return nil, err
    }
    topic.Views++
    _, err = o.Update(topic)

    topic.Labels =strings.Replace(strings.Replace(topic.Labels, "#", " ", -1), "$", " ", -1)
    return topic,err
}


