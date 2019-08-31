package main


import (
	"fmt"
	"douban_spider/history_collector"
)

func main() {
	fmt.Println("Hello word！")
	history_collector.DoubanHandler.MainUrl="https://movie.douban.com/"
	history_collector.DoubanHandler.Test()
}
