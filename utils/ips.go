package utils

import (
	"fmt"
	"strconv"
	"github.com/parnurzeal/gorequest"
	"reflect"
	"strings"
	//"math/rand"
)

const MAX_SIZE_IP_POOL = 100

type IP struct {
	Address string    `json:"address"`
	Port string       `json:"port"`
	Anonymous string  `json:"anonymous"`
	//Location string   `jsong:"location"`
}

type IpPool struct {
	IpList []*IP           `json:"ip_list"`
	SourceWebsite string   `json:"source_website"`
}

type IpPoolMgr struct {
	IpPool
	BannedIp map[string]bool
	UsingIp map[string]bool
}

var Ipmgr = &IpPoolMgr{}

func (ippm *IpPoolMgr) Prepare() error {
	ippm.SourceWebsite = "lab.crossincode.com"
	ippm.BannedIp = make(map[string]bool)
	ippm.UsingIp = make(map[string]bool)
	return nil
}

func (ippm *IpPoolMgr) PrintPoolInfo() error {
	for i, ip := range ippm.IpList {
		fmt.Println("IP: ", i, " Address:Port ", ip.Address, ":", ip.Port, " Annonymous: ", ip.Anonymous)
	}
	return nil
}

func (ippm *IpPoolMgr) FetchIpList(num int) error {
	if num <= 0 {
		return nil
	}
	if num > MAX_SIZE_IP_POOL {
		num = MAX_SIZE_IP_POOL
	}

	ippm.IpList = make([]*IP, 0, MAX_SIZE_IP_POOL)

	resp, body, errs := gorequest.New().Get("http://lab.crossincode.com/proxy/get/?num=" + strconv.Itoa(num)).End()

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
	fmt.Println(proxiesI)
	fmt.Println(reflect.TypeOf(proxiesI))
	var proxies = proxiesI.([]interface{})
	for _, each := range proxies{
		fmt.Println(reflect.TypeOf(each))
		if ipInfoMapI, ok := each.(map[string]interface{}); ok {
			var newIp = &IP{}
			var ipInfoI interface{}
			var ipAnonymousI interface{}
			ipInfoI, _ = ipInfoMapI["http"]
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
			ippm.IpList = append(ippm.IpList, newIp)
		} else {
			fmt.Println("Nothing")
		}
	}

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
		ippm.BannedIp[ip.Address] = true
		ippm.UsingIp[ip.Address] = true
		return ip
	}
	fmt.Println("No availabel ip found in ip list")
	return nil
}
