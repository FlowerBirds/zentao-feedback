package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"log"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"time"
)

var sessionId = ""
var cookies, _ = cookiejar.New(nil)
var baseUrl = os.Getenv("ZENTAO_URL")

var allProduct = [...]Tuple{
	{name: "TempoAI", value: "2"},
	{name: "TempoBI", value: "1"},
	{name: "TempoDF", value: "3"},
	{name: "主数据", value: "15"},
}

var allStatus = [...]Tuple{
	{name: "已处理未关闭", value: "replied"},
	{name: "待处理", value: "wait"},
	{name: "处理中", value: "commenting"},
}

// 登录禅道，返回http.Client和zentaosid
func loginZentao(username, password string) error {
	profileUrl := baseUrl + "/biz/my-profile.html"
	_, content, err := DoGet(profileUrl, nil)
	if strings.Contains(content, "location='/biz/user-login-") {
		log.Println("session expired and try to login.")
	} else {
		log.Println("session is available.")
		return nil
	}

	loginURL := fmt.Sprintf("%s/biz/user-login-%s.json", baseUrl, sessionId)
	params := make(map[string]string)
	params["account"] = username
	params["password"] = password
	params["zentaosid"] = sessionId
	_, _, err = DoGet(loginURL, params)
	if err != nil {
		log.Println(err)
	}

	return nil
}

// 1. 触发后台查询，带 product 和 status 参数
func triggerSearch(product, status string, openedDate string, year string) (string, error) {
	url := baseUrl + "/biz/search-buildQuery.html" // 替换为实际地址
	// 构建表单参数
	params := map[string]string{
		"fieldtitle":          "",
		"fieldid":             "",
		"fieldmodule":         "",
		"fieldproduct":        "",
		"fieldstatus":         "",
		"fieldassignedTo":     "",
		"fieldmailto":         "",
		"fielddesc":           "",
		"fieldpublic":         "",
		"fieldopenedBy":       "",
		"fieldopenedDate":     "",
		"fieldprocessedBy":    "",
		"fieldprocessedDate":  "",
		"fieldclosedBy":       "",
		"fieldclosedDate":     "",
		"fieldclosedReason":   "0",
		"fieldproclassify":    "",
		"fieldprojectname":    "",
		"fieldproductversion": "",
		"fieldphonenum":       "",
		"fieldExpectedtime":   "",
		"fieldProblemImpact":  "",
		"fieldscore":          "",
		"fieldadvice":         "",
		"fieldbugreview":      "",
		"andOr1":              "AND",
		"field1":              "product",
		"operator1":           "=",
		"value1":              product,
		"andOr2":              "and",
		"field2":              "id",
		"operator2":           "=",
		"value2":              "",
		"andOr3":              "and",
		"field3":              "product",
		"operator3":           "=",
		"value3":              "",
		"groupAndOr":          "and",
		"andOr4":              "AND",
		"field4":              "status",
		"operator4":           "=",
		"value4":              status,
		"andOr5":              "and",
		"field5":              "status",
		"operator5":           "=",
		"value5":              "",
		"andOr6":              "and",
		"field6":              "desc",
		"operator6":           "include",
		"value6":              "",
		"module":              "feedback",
		"actionURL":           "/biz/feedback-admin-bysearch-myQueryID-editedDate_desc,id_desc.html",
		"groupItems":          "3",
		"formType":            "lite",
	}
	if openedDate != "" {
		params["field5"] = "openedDate"
		params["operator5"] = ">="
		params["value5"] = openedDate
	}

	if year != "" {
		params["field2"] = "openedDate"
		params["operator2"] = ">="
		params["value2"] = year + "-01-01"

		params["field3"] = "openedDate"
		params["operator3"] = "<="
		params["value3"] = year + "-12-31"
	}

	statusCode, content, err := DoPost(url, params)
	log.Println("Status Code:", statusCode)
	if err != nil {
		return "", err
	}
	return content, nil
}

// 2. 获取反馈列表页面
func fetchFeedbackList() (string, error) {
	url := baseUrl + "/biz/feedback-admin-bysearch-myQueryID-editedDate_desc,id_desc.html" // 替换为实际地址
	statusCode, content, err := DoPost(url, nil)
	log.Println("Status Code:", statusCode)
	if err != nil {
		return "", err
	}
	return content, nil
}

// 解析页面，获取总数、分页等信息
func parseFeedbackOverview(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "解析页面失败"
	}
	// 假设页面有 .pager 元素，包含总数和分页信息
	ul := doc.Find("ul.pager")
	dataRecTotal := ul.AttrOr("data-rec-total", "")
	log.Println(dataRecTotal)
	return dataRecTotal
}

func main() {
	r := gin.Default()
	sessionId = genZentaoSID()
	r.GET("/api/feedback/overview", func(c *gin.Context) {
		year := c.Query("year")
		username := os.Getenv("ZENTAO_USER")
		password := os.Getenv("ZENTAO_PASS")
		if username == "" || password == "" {
			username = "demo" // 可替换
			password = "demo" // 可替换
		}
		err := loginZentao(username, password)
		if err != nil {
			c.String(500, "登录失败: %v", err)
			return
		}
		// 先触发后台查询，传递 product 和 status 参数
		// product := "2"     // 可根据实际需求获取
		// status := "closed" // 可根据实际需求获取
		var content []string
		for _, product := range allProduct {
			var message []string
			var num = 0
			for _, status := range allStatus {
				_, err = triggerSearch(product.value, status.value, "", year)
				if err != nil {
					c.String(500, "查询触发失败: %v", err)
					return
				}
				// 再获取反馈列表页面
				html, err := fetchFeedbackList()
				if err != nil {
					c.String(500, "抓取失败: %v", err)
					return
				}
				total := parseFeedbackOverview(html)
				i, err := strconv.Atoi(total)
				num += i
				message = append(message, fmt.Sprintf("%s%d个", status.name, i))
				time.Sleep(time.Duration(500))
			}
			// 本周新增
			monday := getThisWeekMonday()
			log.Println(monday.Format("2006-01-02"))
			_, err = triggerSearch(product.value, "", monday.Format("2006-01-02"), "")
			if err != nil {
				c.String(500, "查询触发失败: %v", err)
				return
			} // 再获取反馈列表页面
			html, err := fetchFeedbackList()
			if err != nil {
				c.String(500, "抓取失败: %v", err)
				return
			}
			total := parseFeedbackOverview(html)
			i, err := strconv.Atoi(total)
			message = append(message, fmt.Sprintf("本周新增%d个", i))
			time.Sleep(time.Duration(200))
			content = append(content, fmt.Sprintf("%s: 未关闭%d个,其中%s", product.name, num, strings.Join(message, ",")))
		}

		c.String(200, strings.Join(content, "\n"))
	})
	r.Run(":8080")
}
