package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lancewf/money-report-go/data"
	elastic "gopkg.in/olivere/elastic.v5"
)

type BillTypeResponse struct {
	Key         int    `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PurchaseResponseEs struct {
	Store       string    `json:"store"`
	Cost        float32   `json:"cost"`
	Notes       string    `json:"notes"`
	Billtypekey int       `json:"billtypekey"`
	Date        time.Time `json:"date"`
}

func main() {
	//esURL := os.Getenv("ELASTICSEARCH_URL")
    esURL := "http://elasticsearch:9200/"
	fmt.Printf("esURL: %s\n", esURL)
	fmt.Println("loading ES...")

	data.LoadData(esURL)

	fmt.Println("starting server")
	r := gin.Default()
	r.GET("/bill_types", test)
	r.GET("/purchases", purchaseSearch)
	r.Run("0.0.0.0:1234")
}

func purchaseSearch(c *gin.Context) {
	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")
	billTypeKey := c.DefaultQuery("bill_type", "-1")
	costGte := c.DefaultQuery("cost_gte", "none")
	costLt := c.DefaultQuery("cost_lt", "none")

	esURL := "http://elasticsearch:9200/"

	esClient, err := elastic.NewClient(elastic.SetURL(esURL), elastic.SetSniff(false), elastic.SetHealthcheck(false))
	if err != nil {
		fmt.Printf("Could not create elasticsearch client %s\n", err)
	}

	bq := elastic.NewBoolQuery()

	if costGte != "none" || costLt != "none" {
		rangeQuery := elastic.NewRangeQuery("cost")

		if costGte != "none" {
			rangeQuery = rangeQuery.Gte(costGte)
		}

		if costLt != "none" {
			rangeQuery = rangeQuery.Lt(costLt)
		}

		bq.Must(rangeQuery)
	}

	if start != "" || end != "" {
		rangeQuery := elastic.NewRangeQuery("date")

		if start != "" {
			rangeQuery = rangeQuery.Gte(start)
		}

		if end != "" {
			rangeQuery = rangeQuery.Lt(end)
		}

		bq.Must(rangeQuery)
	}

	if billTypeKey != "-1" {
		bq.Must(elastic.NewTermQuery("billtypekey", billTypeKey))
	}

	searchResult, err := esClient.Search().
		Index("purchases").
		Query(bq).
		Sort("date", false).
		From(0).Size(1000).
		Do(c) // execute

	if err != nil {
		panic(err)
	}

	var ttyp PurchaseResponseEs
	var response []PurchaseResponseEs
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if purchase, ok := item.(PurchaseResponseEs); ok {
			response = append(response, purchase)
		}
	}

	c.JSON(200, response)
}

func test(c *gin.Context) {
	esURL := "http://elasticsearch:9200/"
	esClient, err := elastic.NewClient(elastic.SetURL(esURL), elastic.SetSniff(false), elastic.SetHealthcheck(false))

	if err != nil {
		fmt.Printf("Could not create elasticsearch client %s\n", err)
	}
	searchResult2, err := esClient.Search().
		Index("bill-type").
		Do(c) // execute

	if err != nil {
		panic(err)
	}

	var ttyp BillTypeResponse
	var response []BillTypeResponse
	for _, item := range searchResult2.Each(reflect.TypeOf(ttyp)) {
		if billType, ok := item.(BillTypeResponse); ok {
			response = append(response, billType)
		}
	}

	c.JSON(200, response)
}
