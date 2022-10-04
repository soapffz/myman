package util_scans

import (
	"fmt"
	"net/url"
	"os/exec"
	"strconv"
	"strings"

	"scan-when-asset-add/db_model"

	"github.com/soapffz/common-go-functions/pkg"
	"gorm.io/gorm"
)

func ScanByNuclei(db *gorm.DB, args_relatedapp_type string, data db_model.BountyAsset, poc_path string, genrepoer_flag bool, serverJkey string) {
	// 传入的是每个bountyasset结构体
	// 使用对应漏洞类型的nuclei模版进行漏扫，扫完后有漏洞url的进行域名解析及权重查询，写进数据库

	// 取出ip和端口组合成链接
	data_url := "http://" + data.Ip + ":" + data.Port

	// 执行对应nuclei脚本扫描链接，每条执行完后都获取执行结果（CombinedOutput方法不能实时）
	nc_command := exec.Command("nuclei", "-s", "medium,high,critical", "-t", poc_path, "-u", data_url, "-silent", "-nts", "-irr")
	output, _ := nc_command.CombinedOutput()
	scan_result := string(output)

	// 如果执行有输出则表示扫描成功,提取漏洞url、ip、端口
	if len(scan_result) > 0 {
		// 获取漏洞链接
		no_head_vul_l := strings.Split(scan_result, "http://")
		if len(no_head_vul_l) == 2 {
			tmp_vul_l := "http://" + no_head_vul_l[1]
			vul_l := strings.Replace(tmp_vul_l, "\n", "", -1)
			// 解析漏洞url得到ip及port
			u, _ := url.Parse(vul_l)
			host := u.Host
			ip_port_l := strings.Split(host, ":")
			ip := ip_port_l[0]
			port := ip_port_l[1]

			// 根据ip进行域名解析及网站权重查询
			domain, root_domain, root_domain_web_weight := pkg.Ip2DomainAndWebWeight(ip)
			if domain != "" && root_domain_web_weight >= 2 {
				// 权重大于指定值，则推送消息通知
				content := root_domain + "\n" + domain + "\n" + strconv.Itoa(root_domain_web_weight) + "\n" + vul_l
				pkg.PushMsgByServerJ(serverJkey, "xxl-job监测到了新漏洞", "有新的"+args_relatedapp_type+"漏洞可以提交了：\n"+content)
			}
			fmt.Println("[" + args_relatedapp_type + "] Writing: " + vul_l)
			// 根据解析结果写入数据库
			UpdateRecordWithVulnInfo(db, args_relatedapp_type, ip, port, vul_l, domain, root_domain, root_domain_web_weight, genrepoer_flag, serverJkey)
		}
	}
}
