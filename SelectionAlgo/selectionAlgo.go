package main

import (
	"code.google.com/p/gorest"
	"encoding/json"
	"fmt"
	"github.com/fzzy/radix/redis"
	"time"
)

type SelectionAlgo struct {
	gorest.RestService `root:"/SelectionAlgo/" consumes:"application/json" produces:"application/json"`
	basicSelection     gorest.EndPoint `method:"GET" path:"/Select/BasicSelection/{Company:int}/{Tenant:int}/{SessionId:string}" output:"[]string"`
}

func (selectionAlgo SelectionAlgo) BasicSelection(Company, Tenant int, SessionId string) []string {

	const longForm = "Jan 2, 2006 at 3:04pm (MST)"

	fmt.Println(Company)
	fmt.Println(Tenant)
	fmt.Println(SessionId)

	ch := make(chan []string)

	go BasicSelection(Company, Tenant, SessionId, ch)
	var result = <-ch
	close(ch)
	return result

}

func AppendIfMissing(dataList []string, i string) []string {
	for _, ele := range dataList {
		if ele == i {
			return dataList
		}
	}
	return append(dataList, i)
}

func GetConcurrencyInfo(_company, _tenant int, _resId, _class, _type, _category string) ConcurrencyInfo {
	client, err := redis.Dial("tcp", redisIp)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)
	key := fmt.Sprintf("ConcurrencyInfo:%d:%d:%s:%s:%s:%s", _company, _tenant, _resId, _class, _type, _category)
	fmt.Println(key)
	strCiObj, _ := client.Cmd("get", key).Str()
	fmt.Println(strCiObj)

	var ciObj ConcurrencyInfo
	json.Unmarshal([]byte(strCiObj), &ciObj)

	return ciObj
}

type ByStringValue []string

func (a ByStringValue) Len() int           { return len(a) }
func (a ByStringValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStringValue) Less(i, j int) bool { return a[i] < a[j] }

type timeSlice []ConcurrencyInfo

func (p timeSlice) Len() int {
	return len(p)
}
func (p timeSlice) Less(i, j int) bool {
	t1, _ := time.Parse(layout, p[i].LastConnectedTime)
	t2, _ := time.Parse(layout, p[j].LastConnectedTime)
	return t1.After(t2)
}
func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
