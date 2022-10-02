package db_model

import "time"

type BountyAssetWildcarddomain struct {
	ID         uint8     `gorm:"primary_key;column:id;type:bigint(20) unsigned;not null;autoIncrementIncrement:1"` // ID主键，自增步长为1
	Rootdomain string    `gorm:"column:rootdomain;type:varchar(255);not null"`                                     //  根域名
	CreatedAt  time.Time `gorm:"column:createtime"`                                                                //创建时间，必须名称为CreatedAt且不设置默认值时才能自动生成
	UpdatedAt  time.Time `gorm:"column:updatetime"`                                                                //更新时间，必须名称为UpdatedAt且不设置默认值时才能自动生成
	Remark     string    `gorm:"column:remark;type:varchar(255)"`                                                  //  一些自定义备注
	Activemark uint8     `gorm:"column:activemark;type:tinyint(1);default:1"`                                      //  是否有效，默认1为有效，失效则置为0
	Addsource  string    `gorm:"column:addsource;type:varchar(255);default:'bounty-targets-data'"`                 //  添加来源,默认为bounty-targets-data
}

func (BountyAssetWildcarddomain) TabName() string {
	return "bounty_asset_wildcarddomain"
}
