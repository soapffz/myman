package db_model

import "time"

type BountyAsset struct {
	ID         int64     `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null;autoIncrementIncrement:1"` // ID主键，自增步长为1
	Ip         string    `gorm:"index:index_asset;ccolumn:ip;type:varchar(255)"`                                   //  资产ip
	Port       string    `gorm:"index:index_asset;column:port;type:varchar(255)"`                                  //  资产端口
	Rootdomain string    `gorm:"column:rootdomain;type:varchar(255)"`                                              //  根域名
	Domain     string    `gorm:"column:domain;type:varchar(255)"`                                                  //  域名
	AssetUrl   string    `gorm:"column:asset_url;type:varchar(255)"`                                               //  资产url，例如111.111.111.111:8443
	FullUrl    string    `gorm:"column:full_url;type:varchar(255)"`                                                //  url全量链接
	VulnUrl    string    `gorm:"column:vuln_url;type:varchar(255)"`                                                //  其他链接，如漏洞链接
	Relatedapp string    `gorm:"column:relatedapp;type:varchar(255)"`                                              //  关联的app名称
	Activemark int64     `gorm:"column:activemark;default:1"`                                                      //  是否存活，默认1为有效，失效则置为0
	Remark     string    `gorm:"column:remark;type:varchar(255)"`                                                  //  一些自定义备注
	Addsource  string    `gorm:"column:addsource;type:varchar(255);default:quake"`                                 //  添加来源,默认为quake
	CreatedAt  time.Time `gorm:"column:createtime"`                                                                //  创建时间
	UpdatedAt  time.Time `gorm:"column:updatetime"`                                                                //  更新时间
}

func (BountyAsset) TabName() string {
	return "bounty_asset"
}
