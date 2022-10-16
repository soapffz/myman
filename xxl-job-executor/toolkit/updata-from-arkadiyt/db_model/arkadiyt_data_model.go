package db_model

import "time"

// Hackerone
type ArkadiytHackerone struct {
	Uid                               uint8     `gorm:"primary_key;column:uid;type:bigint(20);not null;autoIncrementIncrement:1"` //  自增id，因为原数据中已经有id字段了
	Allowsbountysplitting             bool      `gorm:"column:allowsbountysplitting"`                                             //  是否允许分钱
	Averagetimetobountyawarded        float64   `gorm:"column:averagetimetobountyawarded"`                                        //  平均给钱时间
	Averagetimetofirstprogramresponse float64   `gorm:"column:averagetimetofirstprogramresponse"`                                 //  平均第一反应时间
	Averagetimetoreportresolved       float64   `gorm:"column:averagetimetoreportresolved"`                                       //  平均报告解决时间
	Handle                            string    `gorm:"column:handle;type:varchar(255)"`                                          //  项目简称
	ID                                int64     `gorm:"column:id"`                                                                //  项目id
	Managedprogram                    bool      `gorm:"column:managedprogram"`                                                    //  是否由hackerone托管
	Name                              string    `gorm:"column:name;type:varchar(255)"`                                            //  项目名称
	Offersbounties                    bool      `gorm:"column:offersbounties"`                                                    //  是否提供赏金
	Offersswag                        bool      `gorm:"column:offersswag"`                                                        //  是否提供礼品
	Responseefficiencypercentage      int64     `gorm:"column:responseefficiencypercentage"`                                      //  报告解决率
	Submissionstate                   string    `gorm:"column:submissionstate;type:varchar(255);type:varchar(255)"`               //  可提交状态
	Url                               string    `gorm:"column:url;type:varchar(255)"`                                             //  项目链接
	Website                           string    `gorm:"column:website;type:varchar(255)"`                                         //  项目主要站点
	Inscope                           bool      `gorm:"column:inscope"`                                                           //  是否在范围内
	Assetidentifier                   string    `gorm:"column:assetidentifier"`                                                   //  子资产标识符
	Assettype                         string    `gorm:"column:assettype"`                                                         //  子资产类型
	AvailabilityRequirement           string    `gorm:"column:availabilityrequirement"`
	ConfidentialityRequirement        string    `gorm:"column:confidentialityrequirement"`
	Eligibleforbounty                 bool      `gorm:"column:eligibleforbounty"`     //  是否可以获得赏金
	Eligibleforsubmission             bool      `gorm:"column:eligibleforsubmission"` //  是否可以提交
	Instruction                       string    `gorm:"column:instruction"`           //  子资产介绍
	IntegrityRequirement              string    `gorm:"column:integrityrequirement"`
	Maxseverity                       string    `gorm:"column:maxseverity"`                                               //  子资产最大漏洞等级
	CreatedAt                         time.Time `gorm:"column:createtime"`                                                //  创建时间
	UpdatedAt                         time.Time `gorm:"column:updatetime"`                                                //  更新时间
	Activemark                        bool      `gorm:"column:activemark;type:tinyint(1);default:1"`                      //  资产是否有效，默认1为有效，失效则置为0
	Addsource                         string    `gorm:"column:addsource;type:varchar(255);default:'bounty-targets-data'"` //  添加来源,默认为bounty-targets-data
}

func (ArkadiytHackerone) TableName() string {
	return "arkadiyt_hackerone"
}

// bugcrowd
type ArkadiytBugcrowd struct {
	ID                int64     `gorm:"column:id"`                //  自增id
	Name              string    `gorm:"column:name"`              //  项目名称
	Url               string    `gorm:"column:url"`               //  项目链接
	Allowsdisclosure  bool      `gorm:"column:allowsdisclosure"`  //  是否允许披露漏洞
	Managedbybugcrowd bool      `gorm:"column:managedbybugcrowd"` //  是否由bugcrowd托管
	Safeharbor        string    `gorm:"column:safeharbor"`        //  是否为安全港
	Maxpayout         int64     `gorm:"column:maxpayout"`         //  最高漏洞付款金额
	Inscope           bool      `gorm:"column:inscope"`           //  是否在范围内
	Type              string    `gorm:"column:type"`              //  子项目类型
	Target            string    `gorm:"column:target"`            //  子项目目标
	CreatedAt         time.Time `gorm:"column:createtime"`        //  创建时间
	UpdatedAt         time.Time `gorm:"column:updatetime"`        //  更新时间
	Activemark        int64     `gorm:"column:activemark"`        //  资产是否有效，默认1为有效，失效则置为0
	Addsource         string    `gorm:"column:addsource"`         //  添加来源,默认为bounty-targets-data
}

func (ArkadiytBugcrowd) TableName() string {
	return "arkadiyt_bugcrowd"
}

// Intigriti
type ArkadiytIntigriti struct {
	Uid                  int64     `gorm:"column:uid"`                  //  自增id，因为原数据中已经有id字段了
	ID                   string    `gorm:"column:id"`                   //  项目唯一标识符号
	Name                 string    `gorm:"column:name"`                 //  项目简称
	Companyhandle        string    `gorm:"column:companyhandle"`        //  公司简称
	Handle               string    `gorm:"column:handle"`               //  标识
	Url                  string    `gorm:"column:url"`                  //  项目链接
	Status               string    `gorm:"column:status"`               //  项目状态
	Confidentialitylevel string    `gorm:"column:confidentialitylevel"` //  项目保密等级
	Minbounty            float64   `gorm:"column:minbounty"`            //  最小赏金
	Minbountycurrency    string    `gorm:"column:minbountycurrency"`    //  最小赏金币种
	Maxbounty            float64   `gorm:"column:maxbounty"`            //  最大赏金
	Maxbountycurrency    string    `gorm:"column:maxbountycurrency"`    //  最大赏金币种
	Inscope              bool      `gorm:"column:inscope"`              //  是否在范围内
	Type                 string    `gorm:"column:type"`                 //  子项目类型
	Endpoint             string    `gorm:"column:endpoint"`             //  子项目目标
	Description          string    `gorm:"column:description"`          //  子项目描述
	CreatedAt            time.Time `gorm:"column:createtime"`           //  创建时间
	UpdatedAt            time.Time `gorm:"column:updatetime"`           //  更新时间
	Activemark           int64     `gorm:"column:activemark"`           //  资产是否有效，默认1为有效，失效则置为0
	Addsource            string    `gorm:"column:addsource"`            //  添加来源,默认为bounty-targets-data
}

func (ArkadiytIntigriti) TableName() string {
	return "arkadiyt_intigriti"
}
