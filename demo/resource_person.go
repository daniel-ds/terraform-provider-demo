package demo

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
)

const index = "terraform"

type Person struct {
	Age   int
	Hobby string
}

func resourcePerson() *schema.Resource {
	return &schema.Resource{
		Create: resourcePersonCreate,
		Read:   resourcePersonRead,
		Update: resourcePersonUpdate,
		Delete: resourcePersonDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"age": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"hobby": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourcePersonCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	client := meta.(*elasticsearch.Client)
	reader := esutil.NewJSONReader(Person{
		Age:   d.Get("age").(int),
		Hobby: d.Get("hobby").(string),
	})

	response, e := client.Index(index, reader, client.Index.WithDocumentID(name), client.Index.WithOpType("create"))
	if e != nil {
		return e
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusConflict {
		return fmt.Errorf("person: %s already exists", name)
	}
	if response.IsError() {
		return fmt.Errorf("error in Create response for resource with name: %s, Status code: %v", name, response.StatusCode)
	}
	d.SetId(name)
	return resourcePersonRead(d, meta)
}

func resourcePersonRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*elasticsearch.Client)
	response, e := client.Get(index, d.Id())
	if e != nil {
		return e
	}
	defer response.Body.Close()
	if response.IsError() && response.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error in Read response for resource with name: %s, Status code: %v", d.Id(), response.StatusCode)
	}
	if response.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Removing person from state because no longer exists")
		d.SetId("")
		return nil
	}
	bytes, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return e
	}

	manyBytes := gjson.GetManyBytes(bytes, "_source.Age", "_source.Hobby")
	d.Set("age", manyBytes[0].Int())
	d.Set("hobby", manyBytes[1].Str)

	return nil
}

func resourcePersonUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Id()
	client := meta.(*elasticsearch.Client)
	reader := esutil.NewJSONReader(Person{
		Age:   d.Get("age").(int),
		Hobby: d.Get("hobby").(string),
	})

	response, e := client.Index(index, reader, client.Index.WithDocumentID(name))
	if e != nil {
		return e
	}
	defer response.Body.Close()
	if response.IsError() {
		return fmt.Errorf("error in Update response for resource with name: %s, Status code: %v", d.Id(), response.StatusCode)
	}
	return resourcePersonRead(d, meta)
}

func resourcePersonDelete(d *schema.ResourceData, meta interface{}) error {
	name := d.Id()
	client := meta.(*elasticsearch.Client)

	response, e := client.Delete(index, name)
	if e != nil {
		return e
	}
	defer response.Body.Close()
	if response.IsError() && response.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error deleting person with name: %s, StatusCode: %v", name, response.StatusCode)
	}
	return nil
}
