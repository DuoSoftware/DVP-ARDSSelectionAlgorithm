// SelectionAlgo project main.go
package main

import (
	"code.google.com/p/gorest"
	"fmt"
	"net/http"
)

type Configuration struct {
	RedisIp   string
	RedisDb   int
	RedisPort string
	Port      string
}

type EnvConfiguration struct {
	RedisIp   string
	RedisDb   string
	RedisPort string
	Port      string
}

type AttributeData struct {
	AttributeCode     []string
	AttributeClass    string
	AttributeType     string
	AttributeCategory string
	WeightPrecentage  string
}

type Request struct {
	Company       int
	Tenant        int
	Class         string
	Type          string
	Category      string
	SessionId     string
	AttributeInfo []AttributeData
}

type ConcurrencyInfo struct {
	ResourceId        string
	LastConnectedTime string
}

func main() {
	fmt.Println("Initializting Main")
	InitiateRedis()
	gorest.RegisterService(new(SelectionAlgo))
	http.Handle("/", gorest.Handle())
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}
