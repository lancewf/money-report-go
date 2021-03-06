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

type billTypeResponse struct {
	Key         int    `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type purchaseResponse struct {
	Key         int     `json:"key"`
	Store       string  `json:"store"`
	Cost        float32 `json:"cost"`
	Notes       string  `json:"notes"`
	Billtypekey int     `json:"billtypekey"`
	Dayofmonth  int     `json:"dayofmonth"`
	Month       int     `json:"month"`
	Year        int     `json:"year"`
}

type purchaseResponseEs struct {
	Store       string    `json:"store"`
	Cost        float32   `json:"cost"`
	Notes       string    `json:"notes"`
	Billtypekey int       `json:"billtypekey"`
	Date        time.Time `json:"date"`
}

type allocatedAmountResponse struct {
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

type allocatedAmountEs struct {
	Amount      float32   `json:"amount"`
	Billtypekey int       `json:"billtypekey"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}

// LoadData all the data from the old site into ES
//
func LoadData(esURL string, oldMoneyReportBaseURL string) {
	esClient := eshelp.NewElasticSearchClient(esURL)
	loadAllocatedAmounts(esClient, oldMoneyReportBaseURL)
	loadBillTypes(esClient, oldMoneyReportBaseURL)
	loadPurchases(esClient, oldMoneyReportBaseURL)
}

func loadAllocatedAmounts(esClient *eshelp.ElasticSearchClient, oldMoneyReportBaseURL string) {
	fmt.Println("Transfering AllocatedAmounts")
	var indexName = "allocated-amounts"
	var indexTypeName = "all"

	if !esClient.IndexExists(indexName) {
		fmt.Println("creating index")
		esClient.CreateIndex(indexName)
		fmt.Println("adding mapping")
		esClient.CreateMapping(indexName, indexTypeName, allocatedAmountsMapping())

		fmt.Println("Adding data from old site")
		for _, a := range getAllocatedAmountsFromOldSite(oldMoneyReportBaseURL) {
			start := time.Date(a.StartYear, time.Month(a.StartMonth), a.StartDayofmonth, 0, 0, 0, 0, time.UTC)
			end := time.Date(a.EndYear, time.Month(a.EndMonth), a.EndDayofmonth, 0, 0, 0, 0, time.UTC)
			aes := allocatedAmountEs{a.Amount, a.Billtypekey, start, end}
			esClient.CreateDoc(indexName, indexTypeName, aes)
		}
	}
}

func getAllocatedAmountsFromOldSite(oldMoneyReportBaseURL string) []allocatedAmountResponse {
	var allocatedURL = oldMoneyReportBaseURL + "/index.php/services/getAllocatedAmounts"
	resp, err := http.Get(allocatedURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	allocatedAmounts := new([]allocatedAmountResponse)
	json.NewDecoder(bytes.NewBuffer(body)).Decode(allocatedAmounts)

	return *allocatedAmounts
}

func loadPurchases(esClient *eshelp.ElasticSearchClient, oldMoneyReportBaseURL string) {
	fmt.Println("Transfering PurchasesFromOldSite")
	var indexName = "purchases"
	var indexTypeName = "all"

	if !esClient.IndexExists(indexName) {
		fmt.Println("creating index")
		esClient.CreateIndex(indexName)
		fmt.Println("adding mapping")
		esClient.CreateMapping(indexName, indexTypeName, purchaseMapping())

		fmt.Println("Adding data from old site")
		for _, p := range getPurchasesFromOldSite(oldMoneyReportBaseURL) {
			date := time.Date(p.Year, time.Month(p.Month), p.Dayofmonth, 0, 0, 0, 0, time.UTC)
			pes := purchaseResponseEs{p.Store, p.Cost, p.Notes, p.Billtypekey, date}
			esClient.CreateDoc(indexName, indexTypeName, pes)
		}
	}
}

//{"key":13950,"store":"Friday Harbor Market Place","cost":57.76,"notes":"","billtypekey":7,"dayofmonth":1,"month":7,"year":2017}
func getPurchasesFromOldSite(oldMoneyReportBaseURL string) []purchaseResponse {
	var purchasURL = oldMoneyReportBaseURL + "/index.php/services/getPurchases"
	resp, err := http.PostForm(purchasURL, url.Values{"startmonth": {"5"}, "startdaymonth": {"1"}, "startyear": {"2017"}, "endmonth": {"9"}, "enddaymonth": {"1"}, "endyear": {"2020"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	purchases := new([]purchaseResponse)
	json.NewDecoder(bytes.NewBuffer(body)).Decode(purchases)

	return *purchases
}

func loadBillTypes(esClient *eshelp.ElasticSearchClient, oldMoneyReportBaseURL string) {
	var oldBillTypesURL = oldMoneyReportBaseURL + "/index.php/services/getBillTypes"
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

func getBillTypesFromOldSite(oldBillTypesURL string) []billTypeResponse {
	resp, err := http.Get(oldBillTypesURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	billtypes := new([]billTypeResponse)
	json.NewDecoder(bytes.NewBuffer(body)).Decode(billtypes)

	return *billtypes
}
