// 共获取40368条评论
package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func ConnecttoDataBase() (databaseptr *sqlx.DB, err error) {
	databaseptr, err = sqlx.Open("mysql", "root:********@tcp(127.0.0.1:3306)/BiliBiliComments")
	if err != nil {
		fmt.Println("Failed to connect database:")
		fmt.Println(err)
		return
	} else {
		fmt.Println("Database init successfully")
	}
	return
}

func PageGet(PageUrl string) (PageString string, err error) {
	PageResp, err1 := http.Get(PageUrl)
	if err1 != nil {
		err = err1
		return
	}
	defer PageResp.Body.Close()
	buf := make([]byte, 4096)
	for {
		n, err2 := PageResp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		PageString += string(buf[:n])
	}
	return
}

func filter(page string, rule string, i int) [][]string { //用于进行任意匹配
	regexRULE := regexp.MustCompile(rule)
	matched := regexRULE.FindAllStringSubmatch(page, i)

	return matched
}

func main() {
	database, err1 := ConnecttoDataBase()
	fmt.Println(database)
	if err1 != nil {
		return
	}

	var PackageUrl string
	//var limit int //用于测试时设置范围的变量
	//fmt.Println("请输入获取评论的页数：")
	fmt.Println("按任意键开始")
	fmt.Scanln()
	//fmt.Scanf("%d", &limit)
	//for i := 1; i <= limit; i++ {
	for i := 1; ; i++ {
		if i == 1 {
			PackageUrl = `https://api.bilibili.com/x/v2/reply/main?jsonp=jsonp&next=0&type=1&oid=21071819&mode=3&plat=1&_=1670411770802`
		} else {
			PackageUrl = `https://api.bilibili.com/x/v2/reply/main?jsonp=jsonp&next=` + strconv.Itoa(i) + `&type=1&oid=21071819&mode=3&plat=1&_=1670411770802`
		}

		PageString, err := PageGet(PackageUrl)
		if err != nil {
			continue
		}
		Comments := filter(PageString, `"content":{"message":"(?s:(.*?))",`, -1)

		if len(Comments) == 0 {
			fmt.Println("爬取结束")
			break
		}

		for _, CMT := range Comments {
			fmt.Println(CMT[1])
			FeedBack, err2 := database.Exec("INSERT INTO comments(Passage)VALUES(?)", CMT[1])
			if err2 != nil {
				fmt.Println("Exec failed, skip")
				continue
			} else {
				fmt.Println("Done!")
				fmt.Println(FeedBack)
			}
		}
		time.Sleep(2 * 10000 * time.Millisecond)
	} //for 循环结束

}
