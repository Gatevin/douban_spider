package utils

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/gocolly/colly"
	"time"
	"os"
	//"reflect"
	"strings"
	"strconv"
	"bufio"
	//"math/rand"
)

const MAX_SIZE_IP_POOL = 1000

type IP struct {
	Address string    `json:"address"`
	Port string       `json:"port"`
	Anonymous string  `json:"anonymous"`
	//Location string   `jsong:"location"`
	HttpS string      `json:"https"`
}

type IpPool struct {
	IpList []*IP           `json:"ip_list"`
	SourceWebsite string   `json:"source_website"`
}

type IpPoolMgr struct {
	IpPool
	BannedIp map[string]bool
	UsingIp map[string]bool
	CheckedIp map[string]bool
}

var Ipmgr = &IpPoolMgr{}

func (ippm *IpPoolMgr) Prepare() error {
	ippm.SourceWebsite = "lab.crossincode.com"
	ippm.BannedIp = make(map[string]bool)
	ippm.UsingIp = make(map[string]bool)
	ippm.CheckedIp = make(map[string]bool)
	return nil
}

func (ippm *IpPoolMgr) PrintPoolInfo() error {
	for i, ip := range ippm.IpList {
		fmt.Println("IP: ", i, " Address:Port ", ip.Address, ":", ip.Port, " Annonymous: ", ip.Anonymous)
	}
	return nil
}

func (ippm *IpPoolMgr) CheckProxy(ip *IP) bool {
	if ip == nil {
		return false
	}
	if ip.Anonymous != "高匿" || ip.HttpS != "https" {
		return false
	}
	if checked, exist := ippm.CheckedIp[ip.Address+":"+ip.Port]; checked && exist {
		fmt.Println("IP:", ip, "checked, ignore")
		return false
	}
	ippm.CheckedIp[ip.Address+":"+ip.Port] = true
	//ip.Address = "221.178.232.130"
	//ip.Port = "8080"
	//ip.Address = "221.178.232.130"
	//ip.Port = "8080"
	// 47.110.130.152 8080

	//http
	//ip.Address = "111.29.3.221"
	//ip.Port = "8080"
	request := gorequest.New().Proxy("http://" + ip.Address + ":" + ip.Port).Timeout(time.Duration(5 * time.Second))
	resp, _, errs := request.Post("https://www.baidu.com").End()

	if errs == nil && resp.StatusCode == 200{
		fmt.Println("find usable high anonymous https proxy: ", ip)
		return true
	} else {
		fmt.Println("find " + ip.Address + ":" + ip.Port + " can't be used, ignored")
		//fmt.Println(errs)
		//fmt.Println(body)
		//fmt.Println(resp.StatusCode)
		return false
	}
}

func (ippm *IpPoolMgr) FetchIpList() error {
	ippm.IpList = make([]*IP, 0, MAX_SIZE_IP_POOL)
	for {
		nowLen := len(ippm.IpList)
		if nowLen < 2 {
			fmt.Println("高匿名代理数量不够...搜集中 (", nowLen, "/ 2 )")
			ippm.FetchIpListFromData5U()
			//ippm.FetchIpListFromCrossIncode()
			time.Sleep(time.Duration(20) * time.Second)
		} else {
			fileName := "ip_proxy.txt"
			f, err := os.OpenFile("./save_result/" + fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
			defer f.Close()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				file_content, err := Json.Marshal(ippm.IpList)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					newWriter := bufio.NewWriterSize(f, 1024)
					if _, err = newWriter.Write(file_content); err != nil {
						fmt.Println(err)
					}
					if err = newWriter.Flush(); err != nil {
						fmt.Println(err)
					}
					//fmt.Println(string(file_content))
				}
			}
			break
		}
	}
	return nil
}

func (ippm *IpPoolMgr) FetchIpListFromCrossIncode() error {
	var num = 20
	request := gorequest.New()
	resp, body, errs := request.Get("http://lab.crossincode.com/proxy/get/?num=" + strconv.Itoa(num) + "&head=https").End()

	if errs != nil {
		fmt.Println("get proxy ip error ", errs)
		return nil
	}

	if resp.StatusCode != 200 {
		fmt.Println("Request response not 200, error code: ", resp.StatusCode)
		return nil
	}

	var resMap map[string]interface{}

	if err := Json.Unmarshal([]byte(body), &resMap); err != nil {
		fmt.Println("Json unmarshal error in body")
		return nil
	}

	var proxiesI interface{}
	var okProxy bool
	if proxiesI, okProxy = resMap["proxies"]; okProxy != true {
		fmt.Println("proxies not found in resMap")
		return nil
	}
	var proxies = proxiesI.([]interface{})
	for _, each := range proxies{
		if ipInfoMapI, ok := each.(map[string]interface{}); ok {
			var newIp = &IP{}
			var ipInfoI interface{}
			var ipAnonymousI interface{}
			ipInfoI, _ = ipInfoMapI["https"]
			ipAnonymousI, _ = ipInfoMapI["类型"]
			var ipInfo string
			var ipAnonymous string
			ipInfo, _ = ipInfoI.(string)
			ipAnonymous, _ = ipAnonymousI.(string)
			var tmp []string 
			tmp = strings.Split(ipInfo, ":")
			if len(tmp) == 2 {
				newIp.Address = tmp[0]
				newIp.Port = tmp[1]
			}
			newIp.Anonymous = ipAnonymous
			newIp.HttpS = "https"
			if ippm.CheckProxy(newIp){
				ippm.IpList = append(ippm.IpList, newIp)
			}
		} else {
			fmt.Println("Nothing")
		}
	}
	return nil
}

func (ippm *IpPoolMgr) FetchIpListFromData5U() error {

	var ipColly *colly.Collector
	ipColly = colly.NewCollector(
	)
	
    ipColly.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting: ", r.URL)
        r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
    })
    ipColly.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited: ", r.Request.URL)
		//fmt.Println("Result is: ", string(r.Body))
        //fmt.Println("Result is: ", string(r.Body))
    })

    //这一部分先取影评
    ipColly.OnHTML("ul[style^=margin-top] li[style^=text-align] ul[class=l2]",func(e *colly.HTMLElement){
		ip := &IP{}
        address := e.ChildText("span:first-child>li")
        port := e.ChildText("span:nth-child(2)>li")
		anonymous := e.ChildText("span:nth-child(3)>li")
		httpS := e.ChildText("span:nth-child(4)>li")
		
		ip.Address = address
		ip.Port = port
		ip.Anonymous = anonymous
		ip.HttpS = httpS
		fmt.Println("Find ip: ", ip)
		ip.Port = "8080"
		if ippm.CheckProxy(ip) {
			ippm.IpList = append(ippm.IpList, ip)
		} else{
			ip.Port = "80"
			if ippm.CheckProxy(ip) {
				ippm.IpList = append(ippm.IpList, ip)
			}
		}
    })

    ipColly.OnError(func(r *colly.Response, e error) {
        fmt.Println("Request URL: ", r.Request.URL, " failed with error", e)
        //fmt.Println("Retrying url: ", r.Request.URL)
        //r.Request.Retry()
    })
	
	ipColly.Visit("http://www.data5u.com")

	ipColly.Wait()
	

	return nil
}

func (ippm *IpPoolMgr) BanAnonymousIp(ip *IP) error {
	if ip == nil {
		return nil
	}
	ippm.BannedIp[ip.Address] = true
	return nil
}

func (ippm *IpPoolMgr) GetAnonymousIp() *IP {
	for _, ip := range ippm.IpList {
		if banned, exist := ippm.BannedIp[ip.Address]; banned || exist {
			continue
		}
		if using, exist := ippm.UsingIp[ip.Address]; using || exist {
			continue
		}
		ippm.UsingIp[ip.Address] = true
		return ip
	}
	fmt.Println("No available ip found in ip list")
	return nil
}

func (ippm *IpPoolMgr) GetAnonymousIpWithIndex(index int) *IP {
	if ippm.IpList == nil {
		return nil
	}
	n := len(ippm.IpList)
	if n <= index {
		fmt.Println("index out of IpList bound")
		index = index % n
	}
	ip := ippm.IpList[index]
	if banned, exist := ippm.BannedIp[ip.Address]; banned && exist {
		fmt.Println("ip at index position is banned")
		return nil
	}
	if using, exist := ippm.UsingIp[ip.Address]; using && exist {
		fmt.Println("ip at index position is already using")
		return nil
	}
	ippm.UsingIp[ip.Address] = true
	return ip
}

func (ippm *IpPoolMgr) ReturnAnonymousIp(ip *IP) error {
	if ip == nil {
		return nil
	}
	if using, exist := ippm.UsingIp[ip.Address]; using && exist {
		ippm.UsingIp[ip.Address] = false
	}
	return nil
}