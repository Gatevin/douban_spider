package main


import (
	"fmt"
	"douban_spider/collector"
)

func main() {
	fmt.Println("Hello word！")
	collector.DoubanUserHistoryHandler.Uid = "summermonica"
	collector.DoubanUserHistoryHandler.Uname= "井颯"
	collector.DoubanUserHistoryHandler.FetchHistoryWithUser()
}
