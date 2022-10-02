package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"update-wildcard-domain/db_model"
	"update-wildcard-domain/tools"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/spf13/viper"
	"github.com/thinkeridea/go-extend/exunicode/exutf8"
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

// RuneSubString 提供符文数量截取字符串的方法，针对多字节字符安全高效的截取
// 如果 start 是非负数，返回的字符串将从 string 的 start 位置开始，从 0 开始计算。例如，在字符串 “abcdef” 中，在位置 0 的字符是 “a”，位置 2 的字符串是 “c” 等等。
// 如果 start 是负数，返回的字符串将从 string 结尾处向前数第 start 个字符开始。
// 如果 string 的长度小于 start，将返回空字符串。
//
// 如果提供了正数的 length，返回的字符串将从 start 处开始最多包括 length 个字符（取决于 string 的长度）。
// 如果提供了负数的 length，那么 string 末尾处的 length 个字符将会被省略（若 start 是负数则从字符串尾部算起）。如果 start 不在这段文本中，那么将返回空字符串。
// 如果提供了值为 0 的 length，返回的子字符串将从 start 位置开始直到字符串结尾。
func SubStrRuneSubString(s string, bg_index int, length int) string {
	return exutf8.RuneSubString(s, bg_index, length)
}

var wg sync.WaitGroup

func main() {
	// 初始化后返回数据库链接
	username := viper.GetString("mysql.username") //账号
	password := viper.GetString("mysql.password") //密码
	host := viper.GetString("mysql.host")         //数据库地址，可以是ip或者域名
	port := viper.GetString("mysql.port")         //数据库端口
	dbname := viper.GetString("mysql.dbname")     //数据库名
	db := tools.ConnectMysqlDb(username, password, host, port, dbname)
	//// 读取当前目录下文件扫描进行添加
	// path := "./"
	// files := loadFile(path)
	// for _, v := range files {
	// 	wg.Add(1)
	//  go writeRecordByReadFile(db, path+v)
	// }
	// 解析bounty-targets-datajson数据
	req, err := http.Get("https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/master/data/hackerone_data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)   // ioutil读取后就是[]byte类型
	str, err := jsonvalue.Unmarshal(body) //反序列化，数据类型必须是[]byte类型
	if err != nil {
		log.Println(err)
	}
	str.RangeArray(func(i int, v *jsonvalue.V) bool { // 最外层是数组所以遍历下数组
		scope_data, _ := v.Get("targets", "in_scope") // 每个Object对象，获取targets->in_scope之后又是数组，再遍历下数组
		scope_data.RangeArray(func(i int, v *jsonvalue.V) bool {
			assetype, _ := v.GetString("asset_type")
			if assetype == "URL" { // 获取类型为URL的数据
				url, _ := v.GetString("asset_identifier") //获取url
				bg_2 := SubStrRuneSubString(url, 0, 2)    // 解析URL，若开头两位为*.并且剩余数据中没有*则判断为通配符域名
				root := SubStrRuneSubString(url, 2, 0)
				if bg_2 == "*." && !strings.Contains(root, "*") {
					// fmt.Println(url)
					// fmt.Println(root)
					wg.Add(1)
					go writeParsedUrl(db, root)
				}
			}
			return true // true表示继续遍历，false表示停止遍历
		})
		return true // true表示继续遍历，false表示停止遍历
	})
	wg.Wait()
}

func writeParsedUrl(db *gorm.DB, wildcarddomain string) {
	defer wg.Done()
	var wildcard db_model.BountyAssetWildcarddomain
	wildcard.Rootdomain = wildcarddomain
	res := db.Clauses(clause.OnConflict{
		// ID发生冲突时候，根据id和rootdomain联合索引唯一，更新updatetime
		Columns:   []clause.Column{{Name: "id"}, {Name: "rootdomain"}},
		DoUpdates: clause.AssignmentColumns([]string{"updatetime"}),
	}).Create(&wildcard)
	// 排除ID冲突错误错误后，将其他错误（字段冲突）打印出来
	if res.Error != nil && !strings.Contains(res.Error.Error(), "unique constraint") {
		log.Println("插入出错" + res.Error.Error())
		return
	}
}
