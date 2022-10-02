package scripts

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"update-wildcard-domain/db_model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var wg sync.WaitGroup

func writeRecordByReadFile(db *gorm.DB, filename string) {
	defer wg.Done()
	f, err := os.Open(filename)
	if err != nil {
		log.Println(filename + " error")
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	iter := 0 // 记录出错的行数
	for scanner.Scan() {
		var wildcard db_model.BountyAssetWildcarddomain
		iter++
		////  调用json.Unmarshal()将文本转换为结构体
		// if err = json.Unmarshal([]byte(scanner.Text()), &wildcard.Rootdomain); err != nil {
		// 	log.Println("转换错误--->" + scanner.Text())
		// 	return
		// }
		wildcard.Rootdomain = scanner.Text()
		res := db.Clauses(clause.OnConflict{
			// ID发生冲突时候，根据id和rootdomain联合索引唯一，更新updatetime
			Columns:   []clause.Column{{Name: "id"}, {Name: "rootdomain"}},
			DoUpdates: clause.AssignmentColumns([]string{"updatetime"}),
		}).Create(&wildcard)
		// 排除ID冲突错误错误后，将其他错误（字段冲突）打印出来
		if res.Error != nil && !strings.Contains(res.Error.Error(), "unique constraint") {
			log.Println("插入出错--->" + " 在" + filename + "第" + strconv.Itoa(iter) + "行")
			return
		}
	}
}
