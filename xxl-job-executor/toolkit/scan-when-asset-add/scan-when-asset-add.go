package main

// 通过quake获取指定APP每日更新资产，存入数据库，shell执行脚本传入插入语句，不传入的话打印账号信息

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"scan-when-asset-add/tools"
	"scan-when-asset-add/util_scans"
	"strconv"
	"strings"
	"time"

	"scan-when-asset-add/db_model"

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

func GetDataFromAsset(begin_time string, tags string, db *gorm.DB) []db_model.BountyAsset {
	// 从资产数据中获取[从指定时间]开始的[指定标签]的数据
	now_time := time.Now().Format("2006-01-02 15:04:05")
	// var bountyasset db_model.BountyAsset
	var bountyassets []db_model.BountyAsset
	begin_time = "2022-09-13 01:15:12" // 测试时为查询数据指定较早时间
	db.Where("createtime BETWEEN ? AND ? AND relatedapp = ?", begin_time, now_time, tags).Find(&bountyassets)
	return bountyassets
}

var args_relatedapp_type = flag.String("vp", "", "对应的app关键词，需与数据库中相同")
var genrepoer_flag = flag.Bool("all", false, "是否生成漏洞提交报告模版，默认关闭")

func main() {
	// 解析命令行传入的查询语句
	flag.Parse()
	args_relatedapp_type := *args_relatedapp_type
	if args_relatedapp_type == "" {
		fmt.Println("未传入要扫描的app关键词，请指定！")
		os.Exit(0)
	}
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
	asset_l := GetDataFromAsset(begin_time, args_relatedapp_type, db)

	if len(asset_l) > 0 {
		fmt.Println("数据查询成功，共有" + strconv.Itoa(len(asset_l)) + "条数据")
		nuclei_poc_dir_path := viper.GetString("nuclei.nuclei_poc_dir_path")
		// 根据查询的app关键词不同，进入不同的扫描函数
		switch args_relatedapp_type {
		case "yongyou_nc":
			{
				yongyou_nc_poc_path := nuclei_poc_dir_path + "/yongyou/yongyou-nc-beanshell-rce.yaml"
				vul_result_l := util_scans.Yongyou_nc(asset_l, yongyou_nc_poc_path)
				if vul_result_l != nil {
					// 通知有新的漏洞链接
					pkg.PushMsgByServerJ(serverJkey, "quake监测漏洞通知", "有新的"+args_relatedapp_type+"漏洞，数量个数为"+strconv.Itoa(len(vul_result_l)))
					UpdateRecordWithVulnInfo(db, args_relatedapp_type, vul_result_l, genrepoer_flag)
				}

			}
		default:
			{
				fmt.Println("未匹配到对应app，请检查")
			}
		}

	} else {
		fmt.Println("查询无更新数据，请稍后再来～")
		os.Exit(0)
	}
}

func UpdateRecordWithVulnInfo(db *gorm.DB, args_relatedapp_type string, vul_result_l []string, genrepoer_flag bool) {
	// 根据漏洞类型解析数据更新数据库，并根据是否输出报告模版参数进行操作
	var bountyasset db_model.BountyAsset
	for _, line_vul_l := range vul_result_l {
		vul_l := strings.Replace(line_vul_l, "\n", "", -1)
		// 解析数据
		u, err := url.Parse(vul_l)
		if err != nil {
			log.Fatal(err)
		}
		host := u.Host
		ip_port_l := strings.Split(host, ":")
		ip := ip_port_l[0]
		port := ip_port_l[1]

		// 更新数据库
		db.Model(&bountyasset).Where("ip = ? AND port = ? AND relatedapp = ?", ip, port, args_relatedapp_type).Update("vuln_url", vul_l)

		// 按是否需要生成报告选项进行判断
		if genrepoer_flag == true {
			fmt.Println("报告模版生成在思考中")
		}
	}
}
