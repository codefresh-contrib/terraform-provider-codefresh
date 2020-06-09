package codefresh

import (
	cfClient "github.com/codefresh-io/terraform-provider-codefresh/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if url := os.Getenv("CODEFRESH_API_URL"); url != "" {
						return url, nil
					}
					return "https://g.codefresh.io/api", nil
				},
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CODEFRESH_API_KEY", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"codefresh_project":  resourceProject(),
			"codefresh_pipeline": resourcePipeline(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {

	apiURL := d.Get("api_url").(string)
	token := d.Get("token").(string)

	return cfClient.NewClient(apiURL, token), nil
}
