package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"proj/models"
)

type MainController struct {
	beego.Controller
}


type Respond struct {
	State bool `json:"state"`
	Msg string `json:"msg"`
	Id int64 `json:"id,omitempty"`
	Count int64 `json:"count"`
	CategoryRows []models.Category `json:"category_list,omitempty"`
	QuestionRows []models.Question `json:"question_list,omitempty"`
}

func (c *MainController) Get() {
	o := orm.NewOrm()
	o.Using("default")

	var q models.Question
	q.Category_id = 5
	q.Title = c.GetString("title")
	q.Content = c.GetString("content")
	id, err := o.Insert(&q)

	if err != nil {
		beego.Error(err)
	} else {
		c.Data["Website"] = id
	}
	
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

type ListController struct {
	beego.Controller
}

func (c *ListController) Get() {
	var res Respond
	o := orm.NewOrm()
	o.Using("default")

	op := c.GetString("op")

	switch op {
	case "list":
		list_type := c.GetString("type")
		if (list_type == "category") {
			var items []models.Category
			num, err := o.Raw("SELECT * FROM category").QueryRows(&items)
			if err == nil {
				res.State = true
				res.Count = num
				res.CategoryRows = items
			} else {
				res.State = false
				res.Msg = err.Error()
			}
		} else if (list_type == "question") {
			var items []models.Question
			category_id, err := c.GetInt("category_id")
			if (err == nil) {
				num, err := o.Raw("SELECT id,category_id,title FROM question where category_id = ?", category_id).QueryRows(&items)
				if err == nil {
					res.State = true
					res.Count = num
					res.QuestionRows = items

				} else {
					res.State = false
					res.Msg = err.Error()
				}
			} else {

			}


		}
	case "add_category":
		title := c.GetString("title")
		if (title != "") {
			item := models.Category{Title:title}
			id, err := o.Insert(&item)
			if (err == nil) {
				res.State = true
				res.Id = id
			} else {
				res.State = false
				res.Msg = err.Error()
			}
		} else {
			res.State = false
			res.Msg = "invalid argument"
		}
	case "add":

		category_id, err := c.GetInt("category_id");
		if err == nil {
			var q models.Question
			q.Category_id = category_id
			q.Title = c.GetString("title")
			q.Content = c.GetString("content")
			if (q.Title == "" || q.Content == "") {
				res.State = false
				res.Msg = "invalid argument"
			} else {
				id, err := o.Insert(&q)
				if err == nil {
					res.State = true
					res.Id = id
				} else {
					res.State = false
					res.Msg = err.Error()
				}
			}
		} else {
			res.State = false
			res.Msg = err.Error()
		}

	case "del":

	case "get":
		id, err := c.GetInt("id")
		if (err == nil) {
			var q models.Question
			q.Id = id;
			err = o.Read(&q)
			if (err == nil) {
				c.Data["json"] = &q
				c.ServeJSON()
				return
			} else {
				res.State = false
				res.Msg = err.Error()
			}
		} else {
			res.State = false
			res.Msg = err.Error()
		}

	case "edit":

	}

	c.Data["json"] = &res
	c.ServeJSON()
}
