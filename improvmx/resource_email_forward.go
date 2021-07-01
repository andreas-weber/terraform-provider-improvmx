package improvmx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	improvmxApi "github.com/issyl0/go-improvmx"
)

func resourceEmailForward() *schema.Resource {
	return &schema.Resource{
		Create: resourceEmailForwardCreate,
		Read:   resourceEmailForwardRead,
		Update: resourceEmailForwardUpdate,
		Delete: resourceEmailForwardDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEmailForwardImport,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},

			"alias_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"destination_email": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceEmailForwardCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*improvmxApi.Client)
	client.CreateEmailForward(d.Get("domain").(string), d.Get("alias_name").(string), d.Get("destination_email").(string))

	return resourceEmailForwardRead(d, meta)
}

func resourceEmailForwardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*improvmxApi.Client)
	resp := client.GetEmailForward(d.Get("domain").(string), d.Get("alias_name").(string))

	d.SetId(strconv.FormatInt(resp.Alias.Id, 10))
	d.Set("alias_name", resp.Alias.Alias)
	d.Set("destination_email", resp.Alias.Forward)

	return nil
}

func resourceEmailForwardUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*improvmxApi.Client)
	client.UpdateEmailForward(d.Get("domain").(string), d.Get("alias_name").(string), d.Get("destination_email").(string))

	return resourceEmailForwardRead(d, meta)
}

func resourceEmailForwardDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*improvmxApi.Client)
	client.DeleteEmailForward(d.Get("domain").(string), d.Get("alias_name").(string))

	return nil
}

func resourceEmailForwardImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "_")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Error Importing email forward. Please make sure the email forward ID is in the form DOMAIN_EMAILFORWARDNAME (i.e. example.com_hi)")
	}

	d.SetId(parts[1])
	d.Set("domain", parts[0])

	resourceEmailForwardRead(d, meta)

	return []*schema.ResourceData{d}, nil
}
