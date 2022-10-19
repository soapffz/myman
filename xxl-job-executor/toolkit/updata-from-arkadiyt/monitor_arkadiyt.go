package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"updata-from-arkadiyt/db_model"
	"updata-from-arkadiyt/utilsofparsedata"

	"github.com/soapffz/common-go-functions/pkg"
	"github.com/spf13/viper"
	"github.com/thinkeridea/go-extend/exunicode/exutf8"
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
		log.Fatal("[ERR] 初始化失败：" + err.Error())
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

func main() {
	// 初始化后返回数据库链接
	username := viper.GetString("mysql.username") //账号
	password := viper.GetString("mysql.password") //密码
	host := viper.GetString("mysql.host")         //数据库地址，可以是ip或者域名
	port := viper.GetString("mysql.port")         //数据库端口
	dbname := viper.GetString("mysql.dbname")     //数据库名
	db := pkg.GetMysqlConnByGorm(username, password, host, port, dbname)

	// server酱推送的key
	serverJkey := viper.GetString("serverJ.serverJkey")

	// 判断本地是否已经拉取了arkadiyt的数据，没有的话则拉取并退出
	_, err := os.Stat("bounty-targets-data")
	if os.IsNotExist(err) {
		// 文件不存在
		log.Println("本地没有arkadiyt的数据，进行拉取")
		gitClonecmd := exec.Command("git", "clone", "https://github.com/arkadiyt/bounty-targets-data")
		gitClonecmd.Stdout = os.Stdout
		gitClonecmd.Stderr = os.Stderr
		gitCloneerr := gitClonecmd.Run()
		if gitCloneerr != nil {
			log.Fatal("[ERR] git clone失败：" + gitCloneerr.Error())
			pkg.PushMsgByServerJ(serverJkey, "bounty-targets-data下载失败", "bounty-targets-data下载失败："+gitCloneerr.Error())
			os.Exit(0)
		}
		os.Exit(0)
	}

	// 临时测试用
	// utilsofparsedata.UpdateWildcardsData(db, serverJkey)
	// os.Exit(0)

	// 监测拉取arkadiyt数据后获取最新一次变动的文件
	// 初始化查询时间，默认与xxl-job定时时间保持一致，查询此时间范围内变动的文件（与git日志比较）
	xxljob_crontab_second := 259200
	chagedsFileL := GetArkadiytDataChangeFiles(serverJkey, xxljob_crontab_second)

	// 根据获取到有更新的文件列表进行对应更新
	if len(chagedsFileL) >= 1 {
		for _, i := range chagedsFileL {
			switch i {
			case "hackerone":
				// 更新Hackerone Data
				log.Println("[Warn] 监测到Hackerone数据更新，正在解析到数据库中..")
				//把模型与数据库中的表对应起来
				if e := db.AutoMigrate(&db_model.ArkadiytHackerone{}); e != nil {
					log.Fatalln(e.Error())
				}
				utilsofparsedata.UpdateH1Data(db, serverJkey)
			case "bugcrowd":
				// 更新BugCrowd Data
				log.Println("[Warn] 监测到BugCrowd数据更新，正在解析到数据库中..")
				//把模型与数据库中的表对应起来
				if e := db.AutoMigrate(&db_model.ArkadiytBugcrowd{}); e != nil {
					log.Fatalln(e.Error())
				}
				utilsofparsedata.UpdateBugCrowdData(db, serverJkey)
			case "intigriti":
				// 更新Intigriti Data
				log.Println("[Warn] 监测到Intigriti数据更新，正在解析到数据库中..")
				//把模型与数据库中的表对应起来
				if e := db.AutoMigrate(&db_model.ArkadiytIntigriti{}); e != nil {
					log.Fatalln(e.Error())
				}
				utilsofparsedata.UpdateIntigritiData(db, serverJkey)
			case "wildcards":
				// 更新wildcard Data
				log.Println("[Warn] 监测到wildcard数据更新，正在解析到数据库中..")
				//把模型与数据库中的表对应起来
				if e := db.AutoMigrate(&db_model.ArkadiytWildcard{}); e != nil {
					log.Fatalln(e.Error())
				}
				utilsofparsedata.UpdateWildcardsData(db, serverJkey)
			default:
				log.Println("[Info] 暂时没有获取到新的文件变更")
			}
		}
	} else {
		log.Println("[Info] 暂时没有获取到新的文件变更")
	}
	log.Println("本次程序运行完毕，即将退出")
}

func GetArkadiytDataChangeFiles(serverJkey string, xxljob_crontab_second int) []string {
	// 拉取arkadiyt数据后获取最新一次变动的文件，返回对应资产变动的字符串列表

	// 执行git pull语句
	log.Println("[Info] 正在拉取arkadiyt数据.")
	gitpull_cmd := "cd bounty-targets-data/ && git pull"
	cmd1 := exec.Command("/bin/bash", "-c", gitpull_cmd)
	_, err1 := cmd1.CombinedOutput()
	if err1 != nil {
		log.Fatalln("[ERR] 拉取失败：" + err1.Error())
		pkg.PushMsgByServerJ(serverJkey, "bounty-targets-data拉取失败", "请检查网络，报错信息为："+err1.Error())
		os.Exit(0)
	}

	// 使用git获取最近一次文件变动时间，与xxl-job定时时间进行比较
	// 获取最近一次提交的commit
	getCommitId_cmd := "cd bounty-targets-data/ && git rev-parse HEAD"
	cmdxxx := exec.Command("/bin/bash", "-c", getCommitId_cmd)
	getCommitId_output, _ := cmdxxx.CombinedOutput()
	commitId := strings.Replace(string(getCommitId_output), "\n", "", -1)

	// 使用获取到的commit获取这次提交的时间
	getCommIdTime_cmd := "cd bounty-targets-data/ && echo $(git log --pretty=format:'%cd' --date=local --date=format:'%Y-%m-%d %H:%M:%S' " + commitId + " -1)"
	cmdyyy := exec.Command("/bin/bash", "-c", getCommIdTime_cmd)
	getCommitTime_output, _ := cmdyyy.CombinedOutput()

	// 处理提交时间，获取到的时间进行时区转换
	// 提交时间先转换为时间格式
	transCommitTime1, _ := time.Parse("2006-01-02 15:04:05", strings.Replace(string(getCommitTime_output), "\n", "", -1))
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	// 加上时区，会变成字符串，所以还需要再转换一次时间格式得到上海时区的提交时间
	transCommitTime2 := transCommitTime1.In(cstSh).Format("2006-01-02 15:04:05")
	commitTime, _ := time.Parse("2006-01-02 15:04:05", transCommitTime2)
	// 最近一次提交时间加上xxl-job定时时间
	m, _ := time.ParseDuration("1s")
	addTime := commitTime.Add(time.Duration(xxljob_crontab_second) * m)
	// fmt.Println(commitTime)
	// fmt.Println(addTime)

	// 当前时间也进行转换，保证在同一个时区
	// 当前时间加上时区，会变成字符串，在转换一次未时间格式
	parsedNowTime := time.Now().In(cstSh).Format("2006-01-02 15:04:05")
	nowTime, _ := time.Parse("2006-01-02 15:04:05", parsedNowTime)
	// fmt.Println(nowTime)

	var changedFileL []string
	// 比较时间，如果当前时间大于最近一次提交时间加上xxl-job定时时间，说明有新的文件变动
	// fmt.Println(addTime.After(nowTime))
	if addTime.After(nowTime) {
		// 有更新，获取更新的文件列表
		//使用git获取最近一次变动涉及变动的文件
		chafil_cmd := "cd bounty-targets-data/ && echo $(git diff --name-only HEAD~ HEAD)"
		get_diff_filesname_cmd := exec.Command("/bin/bash", "-c", chafil_cmd)
		getdiff_output, err2 := get_diff_filesname_cmd.CombinedOutput()
		if err2 != nil {
			log.Fatalln(err2.Error())
		}
		diff_result_l := strings.Split(string(getdiff_output), "README.md")
		if len(diff_result_l) >= 1 {
			// 如果检测到了文件变动
			for _, i := range diff_result_l {
				i = strings.TrimSpace(i)
				if len(i) > 1 {
					aa := strings.Split(i, "data/")[1]
					if strings.Contains(aa, ".json") {
						// 如果是json类格式，获取网站名称
						bb := strings.Split(aa, "_data.json")[0]
						changedFileL = append(changedFileL, bb)
					} else if strings.Contains(aa, ".txt") {
						bb := strings.Split(aa, ".txt")[0]
						changedFileL = append(changedFileL, bb)
					} else {
						continue
					}
				}
			}
			return changedFileL
		}
		return changedFileL
	} else {
		log.Println("[Info] 定时时间内未获取到新的提交，程序退出，等待下一次查询.")
	}
	return changedFileL
}
