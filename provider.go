package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderConfig holds the provider configuration
type ProviderConfig struct {
	APIBase string
	APIKey  string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_base": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc("LITELLM_API_BASE", nil),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("LITELLM_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"litellm_model":       resourceLiteLLMModel(),
			"litellm_team":        resourceLiteLLMTeam(),
			"litellm_team_member": resourceLiteLLMTeamMember(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &ProviderConfig{
		APIBase: d.Get("api_base").(string),
		APIKey:  d.Get("api_key").(string),
	}
	return config, nil
}
