package demo

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
	return resourcePersonRead(d, meta)
}

func resourcePersonRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePersonUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourcePersonRead(d, meta)
}

func resourcePersonDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
