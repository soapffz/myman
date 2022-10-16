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

func UpdateBugCrowdData(db *gorm.DB, serverJkey string) {
	// 更新Hackerone Data

	// 读取本地拉取的arkadiyt/bounty-targets-data文件
	BcjsonFile, _ := os.Open("bounty-targets-data/data/bugcrowd_data.json")

	// 解析数据，直接获取每一个字段值赋值给结构体，避免结构体嵌套
	BcByte, _ := ioutil.ReadAll(BcjsonFile)
	BcAllStr, err := jsonvalue.Unmarshal(BcByte) //反序列化，数据类型必须是[]byte类型
	if err != nil {
		log.Println(err)
	}
	BcAllStr.RangeArray(func(i int, v *jsonvalue.V) bool { // 最外层是数组所以遍历下数组
		var bwdata db_model.ArkadiytBugcrowd

		bwdata.Name, _ = v.GetString("name")
		bwdata.Url, _ = v.GetString("url")
		bwdata.Allowsdisclosure, _ = v.GetBool("allows_disclosure")
		bwdata.Managedbybugcrowd, _ = v.GetBool("managed_by_bugcrowd")
		bwdata.Safeharbor, _ = v.GetString("safe_harbor")
		bwdata.Maxpayout, _ = v.GetInt64("max_payout")

		// 因为每个json中嵌套了in_scope和out_scope类型，所以还需要进入一下每个子json
		// 进入in_scope获取数据后写入
		inscopeObject, _ := v.Get("targets", "in_scope")
		inscopeObject.RangeArray(func(i int, v *jsonvalue.V) bool {
			// 此循环中结构体insope类型均为in_scope
			bwdata.Inscope = true
			bwdata.Type, _ = v.GetString("type")
			bwdata.Target, _ = v.GetString("target")
			WriteBugCrowdData(db, bwdata)
			return true
		})
		// 进入out_of_scope获取数据后写入
		outscopeObject, _ := v.Get("targets", "out_of_scope")
		outscopeObject.RangeArray(func(i int, v *jsonvalue.V) bool {
			// 此循环中结构体insope类型均为out_of_scope
			bwdata.Inscope = false
			bwdata.Type, _ = v.GetString("type")
			bwdata.Target, _ = v.GetString("target")
			WriteBugCrowdData(db, bwdata)
			return true
		})
		return true // true表示继续遍历，false表示停止遍历
	})
}

func WriteBugCrowdData(db *gorm.DB, bwdatastruct db_model.ArkadiytBugcrowd) {
	res := db.Clauses(clause.OnConflict{
		// 发生冲突时,根据name, type, target确认唯一联合索引（数据库中也已建立唯一联合索引），更新updatetime
		Columns:   []clause.Column{{Name: "name"}, {Name: "type"}, {Name: "target"}},
		DoUpdates: clause.AssignmentColumns([]string{"updatetime"}),
	}).Create(&bwdatastruct)
	// 排除ID冲突错误错误后，将其他错误（字段冲突）打印出来
	if res.Error != nil && !strings.Contains(res.Error.Error(), "unique constraint") {
		log.Println("插入出错" + res.Error.Error())
		return
	}
}
