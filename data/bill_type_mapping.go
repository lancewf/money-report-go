package data

func billTypeMapping() string {
	return `
{
	"dynamic": false,
	"properties": {
		"key": {
			"type": "integer"
		},
		"name": {
			"type": "keyword"
		},
		"description": {
			"type": "text"
		}
	}
}
	`
}
