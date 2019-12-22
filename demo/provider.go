package demo

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DEMO_URL", "http://localhost:9200"),
			},
		},
		ResourcesMap:  map[string]*schema.Resource{},
		ConfigureFunc: configureFunction(),
	}
}

func configureFunction() schema.ConfigureFunc {
	return func(data *schema.ResourceData) (i interface{}, e error) {
		config := elasticsearch.Config{Addresses: []string{data.Get("url").(string)}}
		client, e := elasticsearch.NewClient(config)
		if e != nil {
			return nil, e
		}
		// Validate Client
		_, e = client.Info()
		if e != nil {
			return nil, e
		}
		return client, nil
	}
}
