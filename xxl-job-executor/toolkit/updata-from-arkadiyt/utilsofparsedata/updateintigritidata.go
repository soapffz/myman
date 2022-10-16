package utilsofparsedata

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"updata-from-arkadiyt/db_model"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateIntigritiData(db *gorm.DB, serverJkey string) {
	// 更新Intigriti Data

	// 读取本地拉取的arkadiyt/bounty-targets-data文件
	igJsonFile, _ := os.Open("bounty-targets-data/data/intigriti_data.json")

	// 解析数据，直接获取每一个字段值赋值给结构体，避免结构体嵌套
	igByte, _ := ioutil.ReadAll(igJsonFile)
	jgAllStr, err := jsonvalue.Unmarshal(igByte) //反序列化，数据类型必须是[]byte类型
	if err != nil {
		log.Println(err)
	}
	jgAllStr.RangeArray(func(i int, v *jsonvalue.V) bool { // 最外层是数组所以遍历下数组
		var igdata db_model.ArkadiytIntigriti

		igdata.ID, _ = v.GetString("id")
		igdata.Name, _ = v.GetString("name")
		igdata.Companyhandle, _ = v.GetString("company_handle")
		igdata.Handle, _ = v.GetString("handle")
		igdata.Url, _ = v.GetString("url")
		igdata.Status, _ = v.GetString("status")
		igdata.Confidentialitylevel, _ = v.GetString("confidentiality_level")
		igdata.Minbounty, _ = v.GetFloat64("min_bounty", "value")
		igdata.Minbountycurrency, _ = v.GetString("min_bounty", "currency")
		igdata.Maxbounty, _ = v.GetFloat64("max_bounty", "value")
		igdata.Maxbountycurrency, _ = v.GetString("max_bounty", "currency")

		// 因为每个json中嵌套了in_scope和out_scope类型，所以还需要进入一下每个子json
		// 进入in_scope获取数据后写入
		inScopeObject, _ := v.Get("targets", "in_scope")
		inScopeObject.RangeArray(func(i int, v *jsonvalue.V) bool {
			// 此循环中结构体insope类型均为in_scope
			igdata.Inscope = true
			igdata = IGGetInAndOutScopeData(igdata, v)
			WriteIntigritiData(db, igdata)
			return true
		})
		// 进入out_of_scope获取数据后写入
		outscopeObject, _ := v.Get("targets", "out_of_scope")
		outscopeObject.RangeArray(func(i int, v *jsonvalue.V) bool {
			// 此循环中结构体insope类型均为out_of_scope
			igdata.Inscope = false
			igdata = IGGetInAndOutScopeData(igdata, v)
			WriteIntigritiData(db, igdata)
			return true
		})
		return true // true表示继续遍历，false表示停止遍历
	})
}

func IGGetInAndOutScopeData(igdata db_model.ArkadiytIntigriti, v *jsonvalue.V) db_model.ArkadiytIntigriti {
	// Hackerone数据解析进入in_scope和out_of_scope后取值
	// 因为类型都一样，特意写这个函数减少代码冗余

	igdata.Type, _ = v.GetString("type")
	igdata.Endpoint, _ = v.GetString("endpoint")
	description, _ := v.GetString("description")
	if len(description) > 100 {
		igdata.Description = "too long"
	} else {
		igdata.Description = description
	}
	return igdata
}

func WriteIntigritiData(db *gorm.DB, igdatastruct db_model.ArkadiytIntigriti) {
	res := db.Clauses(clause.OnConflict{
		// 发生冲突时,根据id,type,endpoint确认唯一联合索引（数据库中也已建立唯一联合索引），更新updatetime
		Columns:   []clause.Column{{Name: "id"}, {Name: "type"}, {Name: "endpoint"}},
		DoUpdates: clause.AssignmentColumns([]string{"updatetime"}),
	}).Create(&igdatastruct)
	// 排除ID冲突错误错误后，将其他错误（字段冲突）打印出来
	if res.Error != nil && !strings.Contains(res.Error.Error(), "unique constraint") {
		log.Println("插入出错" + res.Error.Error())
		return
	}
}
