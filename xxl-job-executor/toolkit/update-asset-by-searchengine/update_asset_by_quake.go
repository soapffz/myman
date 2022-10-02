package main

// 通过quake获取指定APP每日更新资产，存入数据库，shell执行脚本传入插入语句，不传入的话打印账号信息

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"update-asset-by-searchengine/db_model"
	"update-asset-by-searchengine/tools"

	"github.com/360quake/quake_go/utils"
	"github.com/hpifu/go-kit/hflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

var wg sync.WaitGroup

func main() {
	// 使用配置初始化数据库链接
	username := viper.GetString("mysql.username") //账号
	password := viper.GetString("mysql.password") //密码
	host := viper.GetString("mysql.host")         //数据库地址，可以是ip或者域名
	port := viper.GetString("mysql.port")         //数据库端口
	dbname := viper.GetString("mysql.dbname")     //数据库名
	db := tools.ConnectMysqlDb(username, password, host, port, dbname)

	// 初始化quakego命令行参数解析
	tools.QuakeGoHflagInit()
	quake_token := viper.GetString("quake.quakekey")
	if quake_token == "" {
		log.Fatal("quake key 获取失败，程序退出")
		os.Exit(-1)
	}

	var reqjson utils.Reqjson
	reqjson.Query = hflag.GetString("args")

	// 判断传入参数是否为空，为空则打印账号信息
	if reqjson.Query == "" {
		// 传入为空，则打印当前账号信息
		info := utils.InfoGet(quake_token)
		data_result, user_result := utils.InfoLoadJson(info)
		fmt.Println("#用户名:", user_result["username"])
		fmt.Println("#邮  箱:", user_result["email"])
		fmt.Println("#手机:", data_result["mobile_phone"])
		fmt.Println("#月度积分:", data_result["month_remaining_credit"])
		fmt.Println("#长效积分:", data_result["constant_credit"])
		fmt.Println("#Token:", data_result["token"])
	} else {
		// 根据传入参数内容进行查询
		reqjson.EndTime = hflag.GetTime("end_time")
		reqjson.IgnoreCache = hflag.GetBool("ignore_cache")
		reqjson.Field = hflag.GetString("field")
		reqjson.QueryTxt = hflag.GetString("file_txt")

		reqjson.StartTime = hflag.GetTime("start_time")
		if reqjson.StartTime.Format("2006-01-02") == strconv.Itoa(time.Now().Year())+"-01-01" {
			// 如果时间默认没有更改，则设置为从今天的0点开始
			today := time.Now().Format("2006-01-02")
			parsed, _ := time.Parse("2006-01-02", today)
			reqjson.StartTime = parsed
		}

		relatedapp := hflag.GetString("relatedapp")
		// if relatedapp == "" {
		// 	log.Fatal("请注意没有传入相关app")
		// } else {
		// 	log.Println("传入的关联app参数为" + relatedapp)
		// }

		reqjson.Start = hflag.GetString("start")
		reqjson.Size = hflag.GetString("size")
		downall_flag := hflag.GetBool("downall_flag")
		if !downall_flag {
			body := utils.SearchServicePost(reqjson, quake_token)
			parseQuekeGoDataAndWriteDb(db, reqjson, body, relatedapp)
		} else {
			reqjson.Size = "100"
			reqjson.Start = "0"

			// 如果要下载全部，则没批次设置下载100条数据，则最多循环10次，当单次下载的数据小于100时结束
			for i := 0; i <= 9; i++ {
				// 每次循环重置新的reqjson.Start
				reqjson.Start = strconv.Itoa(i * 100)
				body := utils.SearchServicePost(reqjson, quake_token)
				parseQuekeGoDataAndWriteDb(db, reqjson, body, relatedapp)
			}
		}

	}
}

func parseQuekeGoDataAndWriteDb(db *gorm.DB, reqjson utils.Reqjson, body string, relatedapp string) {
	// 通过每次命令查询后解析quake_go返回的数据，解析后放入数据库
	// 此处如果有报错参考原项目的
	dataResult := utils.RespLoadJson[utils.SearchJson](body).Data
	if reqjson.Field != "" && reqjson.Field != "ip,port" {
		for _, value := range dataResult {
			if value.Service.HTTP[reqjson.Field] == nil {
				fmt.Println(value.IP + ":" + "  " + strconv.Itoa(value.Port))
			} else {
				fmt.Println(value.IP + ":" + strconv.Itoa(value.Port) + "  " + value.Service.HTTP[reqjson.Field].(string))
			}
			writeRecord(db, value.IP, strconv.Itoa(value.Port), relatedapp)
		}
	} else {
		for _, value := range dataResult {
			fmt.Println(value.IP + ":" + strconv.Itoa(value.Port))
			writeRecord(db, value.IP, strconv.Itoa(value.Port), relatedapp)
		}
	}
}

func writeRecord(db *gorm.DB, ip string, port string, relatedapp string) {
	// 解析每条数据到实例上
	var bountyasset db_model.BountyAsset
	bountyasset.Ip = ip
	bountyasset.Port = port
	bountyasset.Relatedapp = relatedapp
	// 插入数据库
	res := db.Clauses(clause.OnConflict{
		// ID发生冲突时候，根据ip和port索引唯一，更新updatetime
		Columns:   []clause.Column{{Name: "ip"}, {Name: "port"}},
		DoUpdates: clause.AssignmentColumns([]string{"updatetime"}),
	}).Create(&bountyasset)
	// 排除ID冲突错误错误后，将其他错误（字段冲突）打印出来
	if res.Error != nil && !strings.Contains(res.Error.Error(), "unique constraint") {
		log.Println("插入出错" + res.Error.Error())
		return
	}
}
