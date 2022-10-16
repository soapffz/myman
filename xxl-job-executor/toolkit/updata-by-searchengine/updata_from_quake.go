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
	"updata-by-searchengine/db_model"
	"updata-by-searchengine/tools"

	"github.com/360quake/quake_go/utils"
	"github.com/hpifu/go-kit/hflag"
	"github.com/soapffz/common-go-functions/pkg"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 配置文件结构体
type Config struct {
	// 这里不一定需要所有配置都写全才能用viper读取
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
	db := pkg.GetMysqlConnByGorm(username, password, host, port, dbname)

	//把模型与数据库中的表对应起来
	if e := db.AutoMigrate(&db_model.SeacrhEngineAsset{}); e != nil {
		log.Fatalln(e.Error())
	}

	// 初始化quakego命令行参数解析
	tools.QuakeGoHflagInit()
	quake_token := viper.GetString("quake.quakekey")
	if quake_token == "" {
		log.Fatal("quake key 获取失败，程序退出")
		os.Exit(-1)
	}

	// 读取参数传入值，初始化
	var reqjson utils.Reqjson
	reqjson.Query = hflag.GetString("args")
	reqjson.EndTime = hflag.GetTime("end_time")
	reqjson.IgnoreCache = hflag.GetBool("ignore_cache")
	reqjson.Field = hflag.GetString("field")
	reqjson.QueryTxt = hflag.GetString("file_txt")
	reqjson.StartTime = hflag.GetTime("start_time")
	reqjson.Start = hflag.GetString("start")
	reqjson.Size = hflag.GetString("size")

	downall_flag := hflag.GetBool("downall_flag")
	updateallapp_flag := hflag.GetBool("updateallapp_flag")
	relatedapp := hflag.GetString("relatedapp")

	// 判断是否需要更新所有数据，若是，则忽略本次传入的其他所有语句
	if updateallapp_flag {
		// 如果指定了要更新所有app，则忽略其他查询语句，从读取配置文件中的kv查询语句
		//k是关联app关键词，v是网络空间搜索引擎查询语句,根据v使用搜索引擎查询后，根据k在数据库中写入对应app的数据
		quake_query_statement := viper.GetStringMap("quake_query_statement")
		for k, v := range quake_query_statement {
			relatedapp := k
			reqjson.Query = v.(string)
			InitQuakeReqjsonAndDownloadData(db, reqjson, downall_flag, quake_token, relatedapp)
		}
	} else {
		// 如果没有指定要更新所有app，则根据传入语句进行判断
		if reqjson.Query == "" && relatedapp == "" {
			// 如果没有传入查询语句，则打印当前账号信息，并直接退出
			fmt.Println("没有传入搜索语句也没传入关联app关键词，打印账号信息并退出")
			info := utils.InfoGet(quake_token)
			data_result, user_result := utils.InfoLoadJson(info)
			fmt.Println("#用户名:", user_result["username"])
			fmt.Println("#邮  箱:", user_result["email"])
			fmt.Println("#手机:", data_result["mobile_phone"])
			fmt.Println("#月度积分:", data_result["month_remaining_credit"])
			fmt.Println("#长效积分:", data_result["constant_credit"])
			fmt.Println("#Token:", data_result["token"])
			os.Exit(0)
		}
		if relatedapp == "" {
			// 没有传入关联的app则提示没有传入app关联标签
			fmt.Println("未传入关联app关键词，本次数据更新将不带标签！")
		}
		if reqjson.Query == "" {
			// 没有传入直接搜索语句则提示没有传入搜索语句
			fmt.Println("没有传入搜索语句，将尝试使用标签进行查询")
			reqjson.Query = viper.GetString("quake_query_statement." + relatedapp)
		}

		// 经过逻辑判断，此处reqjson.Query 传入直接输入的查询语句或经关联app搜索到的查询语句
		// 关联app不一定传入，未传入关联app时会提示本次数据更新不带标签
		InitQuakeReqjsonAndDownloadData(db, reqjson, downall_flag, quake_token, relatedapp)
	}
}

func InitQuakeReqjsonAndDownloadData(db *gorm.DB, reqjson utils.Reqjson, downall_flag bool, quake_token string, relatedapp string) {
	// 传入请求结构体及下载数据的flag等，下载数据
	if reqjson.StartTime.Format("2006-01-02") == strconv.Itoa(time.Now().Year())+"-01-01" {
		// 如果时间默认没有更改，则设置为从2个小时前到现在
		// h, _ := time.ParseDuration("-1h")
		today := time.Now().Format("2006-01-02")
		// yesday := time.Now().AddDate(0, 0, -1).Format("2006-01-02") // 测试时注释掉此语句则查询从昨天到现在的数据
		// lastyear_today := time.Now().AddDate(-1, 0, 0).Format("2006-01-02") // 测试时注释掉此语句则查询去年今天到现在的数据
		parsed, _ := time.Parse("2006-01-02", today)
		reqjson.StartTime = parsed
	}

	if !downall_flag {
		body := utils.SearchServicePost(reqjson, quake_token)
		ParseQuekeGoDataAndWriteDb(db, reqjson, body, relatedapp)
	} else {
		reqjson.Size = "100"
		reqjson.Start = "0"

		// 如果要下载全部，则没批次设置下载100条数据，则最多循环100次，当单次下载的数据返回为空时跳出循环
		for i := 0; i <= 100; i++ {
			// 每次循环重置新的reqjson.Start
			reqjson.Start = strconv.Itoa(i * 100)
			body := utils.SearchServicePost(reqjson, quake_token)
			if len(body) <= 130 {
				// 此处还需要优化，暂时测试没有返回值的quake查询结果长度为121或122，取130作为测试值
				break
			}
			// fmt.Println(len(body))
			// 解析quake返回的数据并写入数据库
			ParseQuekeGoDataAndWriteDb(db, reqjson, body, relatedapp)
		}
	}
}

func ParseQuekeGoDataAndWriteDb(db *gorm.DB, reqjson utils.Reqjson, body string, relatedapp string) {
	// 通过每次命令查询后解析quake_go返回的数据，解析后放入数据库
	// 此处如果有报错参考原项目的
	dataResult := utils.RespLoadJson[utils.SearchJson](body).Data
	// fmt.Println(dataResult)
	if reqjson.Field != "" && reqjson.Field != "ip,port" {
		for _, value := range dataResult {
			if value.Service.HTTP[reqjson.Field] == nil {
				fmt.Println(value.IP + ":" + "  " + strconv.Itoa(value.Port))
			} else {
				fmt.Println(value.IP + ":" + strconv.Itoa(value.Port) + "  " + value.Service.HTTP[reqjson.Field].(string))
			}
			WriteSerachEngineAsset(db, value.IP, strconv.Itoa(value.Port), relatedapp)
		}
	} else {
		for _, value := range dataResult {
			fmt.Println(value.IP + ":" + strconv.Itoa(value.Port))
			WriteSerachEngineAsset(db, value.IP, strconv.Itoa(value.Port), relatedapp)
		}
	}
}

func WriteSerachEngineAsset(db *gorm.DB, ip string, port string, relatedapp string) {
	// 解析每条数据到实例上
	var searchengineasset db_model.SeacrhEngineAsset
	searchengineasset.Ip = ip
	searchengineasset.Port = port
	searchengineasset.Relatedapp = relatedapp
	// 插入数据库
	res := db.Clauses(clause.OnConflict{
		// ID发生冲突时候，根据ip和port索引唯一，更新updatetime
		Columns:   []clause.Column{{Name: "ip"}, {Name: "port"}},
		DoUpdates: clause.AssignmentColumns([]string{"updatetime"}),
	}).Create(&searchengineasset)
	// 排除ID冲突错误错误后，将其他错误（字段冲突）打印出来
	if res.Error != nil && !strings.Contains(res.Error.Error(), "unique constraint") {
		log.Println("插入出错" + res.Error.Error())
		return
	}
}
