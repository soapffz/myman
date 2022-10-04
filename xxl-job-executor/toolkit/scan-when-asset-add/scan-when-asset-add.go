package main

// 通过quake获取指定APP每日更新资产，存入数据库，shell执行脚本传入插入语句，不传入的话打印账号信息

import (
	"flag"
	"fmt"
	"os"
	"scan-when-asset-add/tools"
	"scan-when-asset-add/util_scans"
	"strconv"
	"time"

	"scan-when-asset-add/db_model"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// 配置文件结构体
type Config struct {
	MySQL MySQLConfig
}

type MySQLConfig struct {
	Port     string
	Host     string
	Username string
	Password string
	Dbname   string
}

func init() {
	// 初始化，viper读取配置文件
	var config Config
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("configs")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	viper.Unmarshal(&config) //将配置文件绑定到config上
}

func GetDataFromAsset(begin_time string, args_relatedapp_type string, scan_alltype_flag bool, db *gorm.DB) []db_model.BountyAsset {
	// 从资产数据中获取[从指定时间]开始的[指定标签]的数据
	now_time := time.Now().Format("2006-01-02 15:04:05")
	// var bountyasset db_model.BountyAsset
	var bountyassets []db_model.BountyAsset
	begin_time = "2022-10-04 08:00:00" // 测试时取消这句注释，即可指定较早时间的数据
	if scan_alltype_flag {
		// 若传入的“扫描所有类型数据”为真，则获取指定时间内所有标签数据并返回
		db.Where("createtime BETWEEN ? AND ?", begin_time, now_time).Find(&bountyassets)
	} else {
		db.Where("createtime BETWEEN ? AND ? AND relatedapp = ?", begin_time, now_time, args_relatedapp_type).Find(&bountyassets)
	}
	return bountyassets
}

var args_relatedapp_type = flag.String("vp", "", "对应的app关键词，需与数据库中相同")
var scan_alltype_flag = flag.Bool("sall", false, "扫描数据库中所有已知标签数据，默认关闭")
var genrepoer_flag = flag.Bool("gen", false, "是否生成漏洞提交报告模版，默认关闭")

func main() {
	// 解析命令行传入的查询语句
	flag.Parse()
	args_relatedapp_type := *args_relatedapp_type
	scan_alltype_flag := *scan_alltype_flag
	genrepoer_flag := *genrepoer_flag

	// 初始化后返回数据库链接
	username := viper.GetString("mysql.username") //账号
	password := viper.GetString("mysql.password") //密码
	host := viper.GetString("mysql.host")         //数据库地址，可以是ip或者域名
	port := viper.GetString("mysql.port")         //数据库端口
	dbname := viper.GetString("mysql.dbname")     //数据库名
	db := tools.ConnectMysqlDb(username, password, host, port, dbname)

	// 读取推送消息的serverJ的key
	serverJkey := viper.GetString("serverJ.serverJkey")

	// 查询数据库指定时间内是否有新数据插入，时间周期与xxl-job定时时间保持一致
	xxljob_crontab_second := 600
	m, _ := time.ParseDuration("-1s")
	m1 := time.Now().Add(time.Duration(xxljob_crontab_second) * m)
	begin_time := m1.Format("2006-01-02 15:04:05")
	asset_l := GetDataFromAsset(begin_time, args_relatedapp_type, scan_alltype_flag, db)
	// 资产列表直接为结构体的数组，在漏洞扫描函数中得到漏洞url及对应解析资产后直接写入数据库

	if len(asset_l) > 0 {
		fmt.Println("数据查询成功，共有" + strconv.Itoa(len(asset_l)) + "条数据")
		nuclei_poc_dir_path := viper.GetString("nuclei.nuclei_poc_dir_path")

		// 根据查询的app关键词不同，进入不同的扫描函数
		for _, data := range asset_l {
			relatedapp_type := data.Relatedapp
			switch relatedapp_type {
			case "yongyou_nc":
				{
					yongyou_nc_poc_path := nuclei_poc_dir_path + "/yongyou/yongyou-nc-beanshell-rce.yaml"
					util_scans.ScanByNuclei(db, relatedapp_type, data, yongyou_nc_poc_path, genrepoer_flag, serverJkey)
				}
			default:
				{
					continue
				}
			}

		}
	} else {
		fmt.Println("查询无更新数据，请稍后再来～")
		os.Exit(0)
	}
}
