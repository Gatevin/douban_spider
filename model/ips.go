package model

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
)

type IP struct {
	Address string    `json:"address"`
	Port string       `json:"port"`
	Anonymous bool    `json:"anonymous"`
	Type string       `json:"type"`
	Location string   `jsong:"location"`
}

type IpPool struct {
	IpList *[]IP           `json:"ip_list"`
	SourceWebsite string   `json:"source_website"`
}

type IpPoolMgr struct {
	IpPool
}

func (ippm *IpPoolMgr) Prepare() error {
	ippm.SourceWebsite = "lab.crossincode.com"
	return nil
}

func (ippm *IpPoolMgr) FetchIpList(num int) error {
	if num <= 0 {
		return nil
	}
	resp, body, errs := gorequest.New().Get("http://lab.crossincode.com/proxy/get/?num=" + num).End()

	if errs !- nil {
		fmt.Println("get proxy ip error, %s", errs)
		return nil
	}

	if response.StatusCode != 200 {
		fmt.Println("Request response not 200, error code: ", response.StatusCode)
		return nil
	}
	
	fmt.Println("Returned body is ", body)

	return nil
}
