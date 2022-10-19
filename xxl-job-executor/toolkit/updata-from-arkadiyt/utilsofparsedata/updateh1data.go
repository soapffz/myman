package utilsofparsedata

import (
	"io"
	"log"
	"os"
	"strings"
	"updata-from-arkadiyt/db_model"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateH1Data(db *gorm.DB, serverJkey string) {
	// 更新Hackerone Data

	// 读取本地拉取的arkadiyt/bounty-targets-data文件
	h1jsonFile, _ := os.Open("bounty-targets-data/data/hackerone_data.json")

	// 解析数据，直接获取每一个字段值赋值给结构体，避免结构体嵌套
	h1Byte, _ := io.ReadAll(h1jsonFile)
	h1AllStr, err := jsonvalue.Unmarshal(h1Byte) //反序列化，数据类型必须是[]byte类型
	if err != nil {
		log.Println(err)
	}
	h1AllStr.RangeArray(func(i int, v *jsonvalue.V) bool { // 最外层是数组所以遍历下数组
		var h1data db_model.ArkadiytHackerone
		h1data.Allowsbountysplitting, _ = v.GetBool("allows_bounty_splitting") // 每个Object对象
		h1data.Averagetimetobountyawarded, _ = v.GetFloat64("average_time_to_bounty_awarded")
		h1data.Averagetimetofirstprogramresponse, _ = v.GetFloat64("average_time_to_first_program_response")
		h1data.Averagetimetoreportresolved, _ = v.GetFloat64("average_time_to_report_resolved")
		h1data.Handle, _ = v.GetString("handle")
		h1data.ID, _ = v.GetInt64("id")
		h1data.Managedprogram, _ = v.GetBool("managed_program")
		h1data.Name, _ = v.GetString("name")
		h1data.Offersbounties, _ = v.GetBool("offers_bounties")
		h1data.Offersswag, _ = v.GetBool("offers_swag")
		h1data.Responseefficiencypercentage, _ = v.GetInt64("response_efficiency_percentage")
		h1data.Submissionstate, _ = v.GetString("submission_state")
		h1data.Url, _ = v.GetString("url")
		h1data.Website, _ = v.GetString("website")

		// 因为每个json中嵌套了in_scope和out_scope类型，所以还需要进入一下每个子json
		// 进入in_scope获取数据后写入
		inscopeObject, _ := v.Get("targets", "in_scope")
		inscopeObject.RangeArray(func(i int, v *jsonvalue.V) bool {
			// 此循环中结构体insope类型均为in_scope
			h1data.Inscope = true
			h1data = H1GetInAndOutScopeData(h1data, v)
			WriteHackeroneData(db, h1data)
			return true
		})
		// 进入out_of_scope获取数据后写入
		outscopeObject, _ := v.Get("targets", "out_of_scope")
		outscopeObject.RangeArray(func(i int, v *jsonvalue.V) bool {
			// 此循环中结构体insope类型均为out_of_scope
			h1data.Inscope = false
			h1data = H1GetInAndOutScopeData(h1data, v)
			WriteHackeroneData(db, h1data)
			return true
		})
		return true // true表示继续遍历，false表示停止遍历
	})
}

func H1GetInAndOutScopeData(h1data db_model.ArkadiytHackerone, v *jsonvalue.V) db_model.ArkadiytHackerone {
	// Hackerone数据解析进入in_scope和out_of_scope后取值
	// 因为类型都一样，特意写这个函数减少代码冗余

	// 资产标识符超过100字节就不存储了
	asset_identifier, _ := v.GetString("asset_identifier")
	if len(asset_identifier) > 100 {
		h1data.Assetidentifier = "too long"
	} else {
		h1data.Assetidentifier = asset_identifier
	}
	h1data.Assettype, _ = v.GetString("asset_type")
	h1data.AvailabilityRequirement, _ = v.GetString("availability_requirement")
	h1data.ConfidentialityRequirement, _ = v.GetString("confidentiality_requirement")
	h1data.Eligibleforbounty, _ = v.GetBool("eligible_for_bounty")
	h1data.Eligibleforsubmission, _ = v.GetBool("eligible_for_submission")
	// 介绍长度超过100字节就不存储了
	instruction, _ := v.GetString("instruction")
	if len(instruction) > 100 {
		h1data.Instruction = "too long"
	} else {
		h1data.Instruction = instruction
	}
	h1data.IntegrityRequirement, _ = v.GetString("integrity_requirement")
	h1data.Maxseverity, _ = v.GetString("max_severity")

	return h1data
}

func WriteHackeroneData(db *gorm.DB, h1datastruct db_model.ArkadiytHackerone) {
	res := db.Clauses(clause.OnConflict{
		// 发生冲突时,根据handle,name,assetidentifier确认唯一联合索引（数据库中也已建立唯一联合索引），更新updatetime
		Columns:   []clause.Column{{Name: "handle"}, {Name: "name"}, {Name: "assetidentifier"}},
		DoUpdates: clause.AssignmentColumns([]string{"updatetime"}),
	}).Create(&h1datastruct)
	// 排除ID冲突错误错误后，将其他错误（字段冲突）打印出来
	if res.Error != nil && !strings.Contains(res.Error.Error(), "unique constraint") {
		log.Println("插入出错" + res.Error.Error())
		return
	}
}
