package util_scans

import (
	"fmt"
	"scan-when-asset-add/db_model"

	"gorm.io/gorm"
)

func UpdateRecordWithVulnInfo(db *gorm.DB, args_relatedapp_type string, ip string, port string, vul_l string, domain string, root_domain string, root_domain_web_weight int, genrepoer_flag bool, serverJkey string) {
	// 传入数据库链接及结构体数组，更新数据库
	var bountyasset db_model.BountyAsset

	// 判断是否解析成功
	if domain != "" {
		bountyasset.Domain = domain
		bountyasset.Rootdomain = root_domain
		bountyasset.WebWeight = int64(root_domain_web_weight)
		// 更新数据库
		db.Model(&bountyasset).Where("ip = ? AND port = ? AND relatedapp = ?", ip, port, args_relatedapp_type).Updates(map[string]interface{}{"vuln_url": vul_l, "domain": domain, "rootdomain": root_domain, "web_weight": root_domain_web_weight})
	} else {
		db.Model(&bountyasset).Where("ip = ? AND port = ? AND relatedapp = ?", ip, port, args_relatedapp_type).Update("vuln_url", vul_l)
	}

	// 按是否需要生成报告选项进行判断
	if genrepoer_flag == true {
		fmt.Println("报告模版生成在思考中")
	}
}
