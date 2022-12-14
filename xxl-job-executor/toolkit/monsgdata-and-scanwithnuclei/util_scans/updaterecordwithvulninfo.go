package util_scans

import (
	"monsgdata-and-scanwithnuclei/db_model"

	"gorm.io/gorm"
)

func UpdateRecordWithVulnInfo(db *gorm.DB, args_relatedapp_type string, ip string, port string, vul_l string, domain string, root_domain string, root_domain_web_weight int, serverJkey string) {
	// 传入数据库链接及结构体数组，更新数据库
	var searchengineasset db_model.SeacrhEngineAsset

	// 如果成功解析域名则将域名等信息也写入数据库
	if domain != "" {
		searchengineasset.Domain = domain
		searchengineasset.Rootdomain = root_domain
		searchengineasset.WebWeight = int64(root_domain_web_weight)
		// 更新数据库
		db.Model(&searchengineasset).Where("ip = ? AND port = ? AND relatedapp = ?", ip, port, args_relatedapp_type).Updates(map[string]interface{}{"vuln_url": vul_l, "domain": domain, "rootdomain": root_domain, "web_weight": root_domain_web_weight})
	} else {
		// 如果没有解析出域名信息则只写入漏洞url
		db.Model(&searchengineasset).Where("ip = ? AND port = ? AND relatedapp = ?", ip, port, args_relatedapp_type).Update("vuln_url", vul_l)
	}
}
