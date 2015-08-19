package main

import (
	"encoding/json"
	"fmt"
	"github.com/fzzy/radix/redis"
	"sort"
	"strings"
)

func BasicSelection(_company, _tenent int, _sessionId string, ch chan []string) {
	requestKey := fmt.Sprintf("Request:%d:%d:%s", _company, _tenent, _sessionId)
	fmt.Println(requestKey)

	client, err := redis.Dial("tcp", redisIp)
	errHndlr(err)
	defer client.Close()

	// select database
	r := client.Cmd("select", redisDb)
	errHndlr(r.Err)
	strResObj, _ := client.Cmd("get", requestKey).Str()
	fmt.Println(strResObj)

	var reqObj Request
	json.Unmarshal([]byte(strResObj), &reqObj)

	var resourceConcInfo = make([]ConcurrencyInfo, 0)
	var matchingResources = make([]string, 0)
	if len(reqObj.AttributeInfo) > 0 {
		var tagArray = make([]string, 3)

		tagArray[0] = fmt.Sprintf("company_%d", reqObj.Company)
		tagArray[1] = fmt.Sprintf("tenant_%d", reqObj.Tenant)
		//tagArray[2] = fmt.Sprintf("class_%s", reqObj.Class)
		//tagArray[3] = fmt.Sprintf("type_%s", reqObj.Type)
		//tagArray[4] = fmt.Sprintf("category_%s", reqObj.Category)
		tagArray[2] = fmt.Sprintf("objtype_%s", "Resource")

		attInfo := make([]string, 0)

		for _, value := range reqObj.AttributeInfo {
			for _, att := range value.AttributeCode {
				attInfo = AppendIfMissing(attInfo, att)
			}
		}

		sort.Sort(ByStringValue(attInfo))
		for _, att := range attInfo {
			fmt.Println("attCode", att)
			tagArray = AppendIfMissing(tagArray, fmt.Sprintf("attribute_%s", att))
		}

		tags := fmt.Sprintf("tag:*%s*", strings.Join(tagArray, "*"))
		fmt.Println(tags)
		val, _ := client.Cmd("keys", tags).List()
		lenth := len(val)
		fmt.Println(lenth)

		for _, match := range val {
			strResKey, _ := client.Cmd("get", match).Str()
			splitVals := strings.Split(strResKey, ":")
			if len(splitVals) == 4 {
				concInfo := GetConcurrencyInfo(reqObj.Company, reqObj.Tenant, splitVals[3], reqObj.Class, reqObj.Type, reqObj.Category)
				resourceConcInfo = append(resourceConcInfo, concInfo)
				//matchingResources = AppendIfMissing(matchingResources, strResKey)
				//fmt.Println(strResKey)
			}
		}

		sort.Sort(timeSlice(resourceConcInfo))

		for _, res := range resourceConcInfo {
			resKey := fmt.Sprintf("Resource:%d:%d:%s", reqObj.Company, reqObj.Tenant, res.ResourceId)
			matchingResources = AppendIfMissing(matchingResources, resKey)
			fmt.Println(resKey)
		}

	}

	ch <- matchingResources

}
