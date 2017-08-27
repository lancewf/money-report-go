package config

type Service struct {
	Name             string `json:"name"`
	Version          string `json:"version"`
	Host             string `json:"host"`
	Port             int    `json:"port"`
	OldURL           string `json:"old_url"`
	ElasticsearchURL string `json:"es_url"`
}

func GetDefault() Service {
	return Service{
		Version:          "0.0.1",
		Name:             "Money Report",
		Host:             "0.0.0.0",
		Port:             1234,
		OldURL:           "",
		ElasticsearchURL: "http://elasticsearch:9200/",
	}
}
