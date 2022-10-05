package util_scans

import (
	"scan-when-asset-add/db_model"

	"gorm.io/gorm"
)

func UpdateRecordWithVulnInfo(db *gorm.DB, args_relatedapp_type string, ip string, port string, vul_l string, domain string, root_domain string, root_domain_web_weight int, serverJkey string) {
	// 传入数据库链接及结构体数组，更新数据库
	var bountyasset db_model.BountyAsset

	// 如果成功解析域名则将域名等信息也写入数据库
	if domain != "" {
		bountyasset.Domain = domain
		bountyasset.Rootdomain = root_domain
		bountyasset.WebWeight = int64(root_domain_web_weight)
		// 更新数据库
		db.Model(&bountyasset).Where("ip = ? AND port = ? AND relatedapp = ?", ip, port, args_relatedapp_type).Updates(map[string]interface{}{"vuln_url": vul_l, "domain": domain, "rootdomain": root_domain, "web_weight": root_domain_web_weight})
	} else {
		// 如果没有解析出域名信息则只写入漏洞url
		db.Model(&bountyasset).Where("ip = ? AND port = ? AND relatedapp = ?", ip, port, args_relatedapp_type).Update("vuln_url", vul_l)
	}
}
