package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/lancewf/money-report-go/eshelp"
)

type BillTypeResponse struct {
	Key         int    `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PurchaseResponse struct {
	Key         int     `json:"key"`
	Store       string  `json:"store"`
	Cost        float32 `json:"cost"`
	Notes       string  `json:"notes"`
	Billtypekey int     `json:"billtypekey"`
	Dayofmonth  int     `json:"dayofmonth"`
	Month       int     `json:"month"`
	Year        int     `json:"year"`
}

type PurchaseResponseEs struct {
	Store       string    `json:"store"`
	Cost        float32   `json:"cost"`
	Notes       string    `json:"notes"`
	Billtypekey int       `json:"billtypekey"`
	Date        time.Time `json:"date"`
}

type AllocatedAmountResponse struct {
	Key             int     `json:"key"`
	Amount          float32 `json:"amount"`
	Billtypekey     int     `json:"billtypekey"`
	StartDayofmonth int     `json:"startdayofmonth"`
	StartMonth      int     `json:"startmonth"`
	StartYear       int     `json:"startyear"`
	EndDayofmonth   int     `json:"enddayofmonth"`
	EndMonth        int     `json:"endmonth"`
	EndYear         int     `json:"endyear"`
}

type AllocatedAmountEs struct {
	Amount      float32   `json:"amount"`
	Billtypekey int       `json:"billtypekey"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}

func LoadData(esURL string) {
	esClient := eshelp.NewElasticSearchClient(esURL)
	loadAllocatedAmounts(esClient)
	loadBillTypes(esClient)
	loadPurchases(esClient)
}

func loadAllocatedAmounts(esClient *eshelp.ElasticSearchClient) {
	fmt.Println("Transfering AllocatedAmounts")
	var indexName = "allocated-amounts"
	var indexTypeName = "all"

	if !esClient.IndexExists(indexName) {
		fmt.Println("creating index")
		esClient.CreateIndex(indexName)
		fmt.Println("adding mapping")
		esClient.CreateMapping(indexName, indexTypeName, allocatedAmountsMapping())

		fmt.Println("Adding data from old site")
		for _, a := range getAllocatedAmountsFromOldSite() {
			start := time.Date(a.StartYear, time.Month(a.StartMonth), a.StartDayofmonth, 0, 0, 0, 0, time.UTC)
			end := time.Date(a.EndYear, time.Month(a.EndMonth), a.EndDayofmonth, 0, 0, 0, 0, time.UTC)
			aes := AllocatedAmountEs{a.Amount, a.Billtypekey, start, end}
			esClient.CreateDoc(indexName, indexTypeName, aes)
		}
	}
}

func getAllocatedAmountsFromOldSite() []AllocatedAmountResponse {
	var allocatedURL = "http://moneyreport.sjcmmsn.com/index.php/services/getAllocatedAmounts"
	resp, err := http.Get(allocatedURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	allocatedAmounts := new([]AllocatedAmountResponse)
	json.NewDecoder(bytes.NewBuffer(body)).Decode(allocatedAmounts)

	return *allocatedAmounts
}

func loadPurchases(esClient *eshelp.ElasticSearchClient) {
	fmt.Println("Transfering PurchasesFromOldSite")
	var indexName = "purchases"
	var indexTypeName = "all"

	if !esClient.IndexExists(indexName) {
		fmt.Println("creating index")
		esClient.CreateIndex(indexName)
		fmt.Println("adding mapping")
		esClient.CreateMapping(indexName, indexTypeName, purchaseMapping())

		fmt.Println("Adding data from old site")
		for _, p := range getPurchasesFromOldSite() {
			date := time.Date(p.Year, time.Month(p.Month), p.Dayofmonth, 0, 0, 0, 0, time.UTC)
			pes := PurchaseResponseEs{p.Store, p.Cost, p.Notes, p.Billtypekey, date}
			esClient.CreateDoc(indexName, indexTypeName, pes)
		}
	}
}

//{"key":13950,"store":"Friday Harbor Market Place","cost":57.76,"notes":"","billtypekey":7,"dayofmonth":1,"month":7,"year":2017}
func getPurchasesFromOldSite() []PurchaseResponse {
	var purchasURL = "http://moneyreport.sjcmmsn.com/index.php/services/getPurchases"
	resp, err := http.PostForm(purchasURL, url.Values{"startmonth": {"7"}, "startdaymonth": {"1"}, "startyear": {"2000"}, "endmonth": {"9"}, "enddaymonth": {"1"}, "endyear": {"2020"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	purchases := new([]PurchaseResponse)
	json.NewDecoder(bytes.NewBuffer(body)).Decode(purchases)

	return *purchases
}

func loadBillTypes(esClient *eshelp.ElasticSearchClient) {
	var oldBillTypesURL = "http://moneyreport.sjcmmsn.com/index.php/services/getBillTypes"
	fmt.Println("Transfering BillTypes")
	var indexName = "bill-type"
	var indexTypeName = "all"

	if !esClient.IndexExists(indexName) {
		fmt.Println("creating index")
		esClient.CreateIndex(indexName)
		fmt.Println("adding mapping")
		esClient.CreateMapping(indexName, indexTypeName, billTypeMapping())

		fmt.Println("Adding data from old site")
		for _, bt := range getBillTypesFromOldSite(oldBillTypesURL) {
			esClient.CreateDoc(indexName, indexTypeName, bt)
		}
	}
}

func getBillTypesFromOldSite(oldBillTypesURL string) []BillTypeResponse {
	resp, err := http.Get(oldBillTypesURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	billtypes := new([]BillTypeResponse)
	json.NewDecoder(bytes.NewBuffer(body)).Decode(billtypes)

	return *billtypes
}
