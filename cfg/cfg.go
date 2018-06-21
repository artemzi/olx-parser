package cfg

type RequestCfg struct {
	BASEURL     string
	CATEGORY    string
	SUBCATEGORY string
	REGION      string
}

func NewRequestCfg() *RequestCfg {
	return &RequestCfg{
		"https://www.olx.ua",
		"transport",
		"legkovye-avtomobili",
		"donetsk",
	}
}