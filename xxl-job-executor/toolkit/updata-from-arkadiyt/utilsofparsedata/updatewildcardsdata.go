package utilsofparsedata

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"updata-from-arkadiyt/db_model"

	"github.com/soapffz/common-go-functions/pkg"
	"github.com/thinkeridea/go-extend/exunicode/exutf8"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

func UpdateWildcardsData(db *gorm.DB, serverJkey string) {
	// 读取本地拉取的arkadiyt/bounty-targets-data文件
	wdtxtFile, _ := os.Open("bounty-targets-data/data/wildcards.txt")
	defer wdtxtFile.Close()

	r := bufio.NewReader(wdtxtFile)
	// 按行读取文件
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("读取文件失败：", err.Error())
			os.Exit(0)
		}
		line = strings.TrimSpace(line)
		// 解析通配符域名，若开头两位为*.并且剩余数据中没有*则继续进行解析
		bg_2 := SubStrRuneSubString(line, 0, 2)
		root := SubStrRuneSubString(line, 2, 0)
		if bg_2 == "*." && !strings.Contains(root, "*") {
			OperateWildCardDomain(db, root)
		}
		// 这里暂时留出来，后面再处理其他特殊的通配符域名
	}
}

func OperateWildCardDomain(db *gorm.DB, root string) {
	// 传入可用于递归的根部域名进行后续操作

	var wd db_model.ArkadiytWildcard
	// 进行字符处理
	newwildcard := GetCleanDomain(root)
	if newwildcard != "" {
		wd.Wildcarddomain = root
		rootdomainofwildcard := pkg.GetRootDomain(root)
		if rootdomainofwildcard != "" {
			wd.Rootdomainofwildcard = rootdomainofwildcard
		}
		WriteWildCardData(db, wd)
	}
}

func WriteWildCardData(db *gorm.DB, wddatastruct db_model.ArkadiytWildcard) {
	// 写入记录
	res := db.Clauses(clause.OnConflict{
		// 发生冲突时,根据domain，更新updatetime
		Columns:   []clause.Column{{Name: "wildcarddomain"}},
		DoUpdates: clause.AssignmentColumns([]string{"updatetime"}),
	}).Create(&wddatastruct)
	// 排除ID冲突错误错误后，将其他错误（字段冲突）打印出来
	if res.Error != nil && !strings.Contains(res.Error.Error(), "unique constraint") {
		log.Println("插入出错" + res.Error.Error())
		return
	}
}

func GetCleanDomain(domain string) string {
	// 传入域名进行清理
	// 去除前后空格
	domain = strings.TrimSpace(domain)
	// 去除www前缀
	domain = strings.TrimPrefix(domain, "www.")
	// 如果传入的是个ip，返回空字符串
	match, _ := regexp.MatchString("[A-Za-z]", domain)
	if !match {
		return ""
	}
	if domain == "" {
		return ""
	} else {
		return domain
	}
}
