/*
	版权所有，侵权必究
	署名-非商业性使用-禁止演绎 4.0 国际
	警告： 以下的代码版权归属hunterhug，请不要传播或修改代码
	你可以在教育用途下使用该代码，但是禁止公司或个人用于商业用途(在未授权情况下不得用于盈利)
	商业授权请联系邮箱：gdccmcm14@live.com QQ:459527502

	All right reserved
	Attribution-NonCommercial-NoDerivatives 4.0 International
	Notice: The following code's copyright by hunterhug, Please do not spread and modify.
	You can use it for education only but can't make profits for any companies and individuals!
	For more information on commercial licensing please contact hunterhug.
	Ask for commercial licensing please contact Mail:gdccmcm14@live.com Or QQ:459527502

	2017.7 by hunterhug
*/
package smartjp

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/hunterhug/parrot/util"
	"fmt"
)

type UrlController struct {
	baseController
}

func (this *UrlController) Monitor() {
	if this.IsAjax() {
		DB := orm.NewOrm()
		err := DB.Using("jpbasicdb")
		if err != nil {
			beego.Error("jpbasicdb err:" + err.Error())
			this.Rsp(false, err.Error())
			return
		}

		var datas []orm.Params
		DB.Raw("SELECT * FROM smart_monitor").Values(&datas)
		this.Data["json"] = &map[string]interface{}{"total": len(datas), "rows": &datas}
		this.ServeJSON()
		this.StopRun()
	}
	this.Layout = this.GetTemplate() + "/base/layout.html"
	this.TplName = this.GetTemplate() + "/url/jmonitor.html"
}

func (this *UrlController) Index() {
	DB := orm.NewOrm()
	err := DB.Using("jpbasicdb")
	if err != nil {
		beego.Error("jpbasicdb err:" + err.Error())
		this.Rsp(false, err.Error())
	}
	var categorys []orm.Params
	DB.Raw("SELECT bigpname as Bigpname,id FROM smart_category where pid=0 group by bigpname,id").Values(&categorys)
	this.Data["category"] = &categorys
	this.Layout = this.GetTemplate() + "/base/layout.html"
	this.TplName = this.GetTemplate() + "/url/jplist.html"

}

func (this *UrlController) Query() {
	DB := orm.NewOrm()
	err := DB.Using("jpbasicdb")
	if err != nil {
		beego.Error("jpbasicdb err:" + err.Error())
		this.Rsp(false, err.Error())
	}
	num := 0
	var maps []orm.Params
	page, _ := this.GetInt("page", 1)
	rows, _ := this.GetInt("rows", 30)
	start := (page - 1) * rows

	name := this.GetString("name")
	next := this.GetString("next")

	if name == "" && next == "" {

		isvalid, _ := this.GetInt("isvalid", 2)
		bigname := this.GetString("bigname")
		small := this.GetString("small")
		level, _ := this.GetInt("level", 1)
		if level > 6 || level < 0 {
			level = 1
		}
		where := []string{}
		wheresql := ""
		if bigname == "" {
		} else {
			where = append(where, `bigpid="`+bigname+`"`)
		}
		if small == "0" || small == "1" {
			where = append(where, "ismall="+small)
		}
		if isvalid == 1 || isvalid == 0 {
			where = append(where, `isvalid=`+util.IS(isvalid))
		}

		if level >= 1 {
			where = append(where, `level=`+util.IS(level))
		}
		if len(where) == 0 {

		} else {
			wheresql = strings.Join(where, " and ")
			wheresql = "where " + wheresql
		}
		dudu := "SELECT * FROM smart_category " + wheresql + " order by createtime limit " + strconv.Itoa(start) + "," + strconv.Itoa(rows) + ";"
		//fmt.Println(dudu)
		DB.Raw(dudu).Values(&maps)

		dudu1 := "SELECT count(*) as num FROM smart_category " + wheresql + ";"

		DB.Raw(dudu1).QueryRow(&num)
	} else if name != "" {
		dudu := "SELECT * FROM smart_category where name=? limit " + strconv.Itoa(start) + "," + strconv.Itoa(rows) + ";"
		DB.Raw(dudu, name).Values(&maps)
		dudu1 := "SELECT count(*) as num FROM smart_category where name=?;"
		DB.Raw(dudu1, name).QueryRow(&num)
	} else if next != "" {
		dudu := "SELECT * FROM smart_category where id regexp \"^%s-[0-9]*$\" limit " + strconv.Itoa(start) + "," + strconv.Itoa(rows) + ";"
		DB.Raw(fmt.Sprintf(dudu, next)).Values(&maps)
		dudu1 := "SELECT count(*) as num FROM smart_category where id regexp \"^%s-[0-9]*$\""
		DB.Raw(fmt.Sprintf(dudu1, next)).QueryRow(&num)
	}

	if len(maps) == 0 {
		this.Data["json"] = &map[string]interface{}{"total": num, "rows": []interface{}{}}
	} else {
		this.Data["json"] = &map[string]interface{}{"total": num, "rows": &maps}
	}
	this.ServeJSON()
}

func (this *UrlController) Update() {
	DB := orm.NewOrm()
	err := DB.Using("jpbasicdb")
	if err != nil {
		beego.Error("jpbasicdb err:" + err.Error())
		this.Rsp(false, err.Error())
	}
	var maps []orm.Params
	isvalid := this.GetString("isvalid")
	id := this.GetString("id")
	page := this.GetString("page")
	if page == "" {
		dudu := "update smart_category set isvalid=? where id=?;"
		_, err := DB.Raw(dudu, isvalid, id).Values(&maps)
		if err == nil {
			this.Rsp(true, "good job")
		} else {
			this.Rsp(false, err.Error())
		}
	} else {
		dudu := "update smart_category set page=? where id=?;"
		_, err := DB.Raw(dudu, page, id).Values(&maps)
		if err == nil {
			this.Rsp(true, "good job")
		} else {
			this.Rsp(false, err.Error())
		}
	}
}
