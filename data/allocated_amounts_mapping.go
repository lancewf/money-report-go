package data

func allocatedAmountsMapping() string {
	return `
{
	"dynamic": false,
	"properties": {
		"amount": {
			"type": "double"
		},
		"billtypekey": {
			"type": "integer"
		},
		"start": {
			"type": "date"
		},
		"end": {
			"type": "date"
		}
	}
}
	`
}
