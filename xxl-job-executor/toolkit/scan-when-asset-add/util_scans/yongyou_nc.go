package util_scans

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"scan-when-asset-add/db_model"

	"github.com/soapffz/common-go-functions/pkg"
)

func Yongyou_nc(asset_l []db_model.BountyAsset, yongyou_nc_poc_path string) []string {
	// yongyou_nc资产进行扫描
	// 先格式化url，然后写入临时文件用于扫描
	var url_l []string
	for _, singledata := range asset_l {
		single_url := "http://" + singledata.Ip + ":" + singledata.Port
		url_l = append(url_l, single_url)
	}
	nc_url_filepath := pkg.WriteSliceReturnRandomFilename(url_l)

	// 根据文件行数是否有变化来判断文件是否有新增，因为重复执行文件行数不会变化
	before_line_num := pkg.CountFileLine("nuclei_results")

	// 执行对应nuclei脚本
	nc_command := exec.Command("nuclei", "-s", "medium,high,critical", "-t", yongyou_nc_poc_path, "-l", nc_url_filepath, "-o", "nuclei_results", "-silent", "-nts", "-irr")
	nc_command.Stdout = os.Stdout
	nc_command.Stderr = os.Stderr
	err3 := nc_command.Run()
	if err3 != nil {
		if strings.Contains(err3.Error(), "exit status 1") {
			os.Exit(0)
		}
	}

	// 删除从数据库查到的临时文件
	err2 := os.Remove(nc_url_filepath)
	if err2 != nil {
		// 删除失败
		fmt.Println("临时文件删除失败！")
	}

	after_line_num := pkg.CountFileLine("nuclei_results")
	if after_line_num > before_line_num {
		// nuclei扫描出了新的结果，读取新增行进行结果解析
		if before_line_num == 0 {
			before_line_num += 1
		}
		var vul_result_l []string
		new_nc_asset_line_l := pkg.ReadSpecifiedLineInFile("nuclei_results", before_line_num, after_line_num)
		for _, result_line := range new_nc_asset_line_l {
			// 先判断行的长度不为空
			if len(result_line) > 0 {
				split_l := strings.Split(result_line, "[high] ")
				rce_url := split_l[1]
				vul_result_l = append(vul_result_l, rce_url)
			}
		}
		return vul_result_l
	}
	return nil

}
