package data

func purchaseMapping() string {
	return `
{
	"dynamic": false,
	"properties": {
		"store": {
			"type": "keyword"
		},
		"cost": {
			"type": "double"
		},
		"notes": {
			"type": "text"
		},
		"billtypekey": {
			"type": "integer"
		},
		"date": {
			"type": "date"
		}
	}
}
	`
}
