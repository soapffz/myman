package main

// 通过quake获取指定APP每日更新资产，存入数据库，shell执行脚本传入插入语句，不传入的话打印账号信息

import (
	"flag"
	"fmt"
	"log"
	"monsgdata-and-scanwithnuclei/util_scans"
	"os"
	"strconv"
	"time"

	"monsgdata-and-scanwithnuclei/db_model"

	"github.com/soapffz/common-go-functions/pkg"
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

func GetDataFromAsset(xxljob_crontab_second int, args_relatedapp_type string, db *gorm.DB) []db_model.SeacrhEngineAsset {
	// 从资产数据中获取[从指定时间]开始的[指定标签]的数据
	m, _ := time.ParseDuration("-1s")
	begin_time := time.Now().Add(time.Duration(xxljob_crontab_second) * m).Format("2006-01-02 15:04:05")
	now_time := time.Now().Format("2006-01-02 15:04:05")
	// var bountyasset db_model.SeacrhEngineAsset
	var seacrhengineasset []db_model.SeacrhEngineAsset
	// begin_time = "2022-10-04 08:00:00" // 测试时取消这句注释，即可指定较早时间的数据
	db.Where("createtime BETWEEN ? AND ? AND relatedapp = ?", begin_time, now_time, args_relatedapp_type).Find(&seacrhengineasset)
	return seacrhengineasset
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
	db := pkg.GetMysqlConnByGorm(username, password, host, port, dbname)

	//把模型与数据库中的表对应起来
	if e := db.AutoMigrate(&db_model.SeacrhEngineAsset{}); e != nil {
		log.Fatalln(e.Error())
	}

	// server酱推送的key
	serverJkey := viper.GetString("serverJ.serverJkey")

	// 初始化查询时间，默认与xxl-job定时时间保持一致，查询此时间范围段内插入的数据（更新的数据不好统计因为会持续更新）
	xxljob_crontab_second := 600

	// 根据是否要扫描所有app的flag进行判断
	if scan_alltype_flag {
		// 如果要扫描所有app，则从配置文件中读取所有app的关键词及对应nuclei脚本扫描路径
		relatedappandpocpath := viper.GetStringMap("relatedappandpocpath")
		for k, v := range relatedappandpocpath {
			relatedapp := k
			relatedapp_nuclei_poc_path := viper.GetString("nuclei.nuclei_poc_dir_path") + v.(string)
			asset_l := GetDataFromAsset(xxljob_crontab_second, relatedapp, db)
			InitScanWithDataFromDb(db, asset_l, relatedapp, genrepoer_flag, serverJkey, relatedapp_nuclei_poc_path)
		}
	} else {
		// 如果不扫描所有app，则根据输入的关联app从配置文件中读取nuclei扫描脚本文件
		nuclei_poc_path := viper.GetString("relatedappandpocpath." + args_relatedapp_type)
		if nuclei_poc_path == "" {
			fmt.Println("获取对应app的nuclei扫描脚本失败，请检查配置文件")
			os.Exit(0)
		}
		relatedapp_nuclei_poc_path := viper.GetString("nuclei.nuclei_poc_dir_path") + nuclei_poc_path
		asset_l := GetDataFromAsset(xxljob_crontab_second, args_relatedapp_type, db)
		InitScanWithDataFromDb(db, asset_l, args_relatedapp_type, genrepoer_flag, serverJkey, relatedapp_nuclei_poc_path)
	}

	// 资产列表直接为结构体的数组，在漏洞扫描函数中得到漏洞url及对应解析资产后直接写入数据库
}

func InitScanWithDataFromDb(db *gorm.DB, asset_l []db_model.SeacrhEngineAsset, relatedapp string, genrepoer_flag bool, serverJkey string, poc_path string) {
	if len(asset_l) > 0 {
		fmt.Println(relatedapp + " 数据查询成功，共有" + strconv.Itoa(len(asset_l)) + "条数据")

		// 根据查询的app关键词不同，进入不同的扫描函数
		for _, data := range asset_l {
			relatedapp_type := data.Relatedapp
			util_scans.ScanByNuclei(db, relatedapp_type, data, poc_path, genrepoer_flag, serverJkey)
		}
	} else {
		fmt.Println(relatedapp + " 查询无更新数据，请稍后再来～")
	}
}
