package util_scans

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"scan-when-asset-add/db_model"

	"github.com/soapffz/common-go-functions/pkg"
	"gorm.io/gorm"
)

func Yongyou_nc(db *gorm.DB, args_relatedapp_type string, asset_l []db_model.BountyAsset, yongyou_nc_poc_path string, genrepoer_flag bool, serverJkey string) {
	// 使用用友nc的nuclei模版进行漏扫，扫完后有漏洞url的进行域名解析及权重查询，写进数据库

	// 先格式化url，然后写入临时文件用于扫描
	var url_l []string
	for _, singledata := range asset_l {
		single_url := "http://" + singledata.Ip + ":" + singledata.Port
		url_l = append(url_l, single_url)
	}
	nc_url_filepath := pkg.WriteSliceReturnRandomFilename(url_l)

	// 根据文件行数是否有变化来判断文件是否有新增，因为重复执行文件行数不会变化
	before_line_num := pkg.CountFileLine("nuclei_results")
	// fmt.Println(before_line_num)

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
	// fmt.Println(after_line_num)
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
		// fmt.Println(vul_result_l)
		if len(vul_result_l) > 0 {
			// nuclei扫出了新的结果，则根据ip进行域名解析及权重查询，写进数据库
			fmt.Println("nuclei扫描出的新的漏洞，正在解析域名及根域名网站权重...")
			for _, line_vul_l := range vul_result_l {
				vul_l := strings.Replace(line_vul_l, "\n", "", -1)
				// 解析漏洞url得到ip及port
				u, err := url.Parse(vul_l)
				if err != nil {
					log.Fatal(err)
				}
				host := u.Host
				ip_port_l := strings.Split(host, ":")
				ip := ip_port_l[0]
				port := ip_port_l[1]

				// 根据ip进行域名解析及网站权重查询
				domain, root_domain, root_domain_web_weight := pkg.Ip2DomainAndWebWeight(ip)
				if domain != "" && root_domain_web_weight >= 2 {
					// 有新的漏洞可以提交了，则推送通知
					content := root_domain + "\n" + domain + "\n" + strconv.Itoa(root_domain_web_weight) + "\n" + vul_l
					pkg.PushMsgByServerJ(serverJkey, "用友nc新漏洞链接推送", "有新的用友nc漏洞可以提交了：\n"+content)
				}
				UpdateRecordWithVulnInfo(db, args_relatedapp_type, ip, port, vul_l, domain, root_domain, root_domain_web_weight, genrepoer_flag, serverJkey)
			}
		}
	}
}
