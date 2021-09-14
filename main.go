package main

import (
	restful "github.com/emicklei/go-restful/v3"
	"log"
	"net/http"
)


func main() {
	//name := "griblje"
	//estimation := Estimate(name)
	//fmt.Println(estimation)

	//fmt.Println(EstimateAll())
	restful.Add(New())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
