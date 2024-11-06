package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLiteLLMTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiteLLMTeamCreate,
		Read:   resourceLiteLLMTeamRead,
		Update: resourceLiteLLMTeamUpdate,
		Delete: resourceLiteLLMTeamDelete,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLiteLLMTeamCreate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("not implemented")
}

func resourceLiteLLMTeamRead(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("not implemented")
}

func resourceLiteLLMTeamUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("not implemented")
}

func resourceLiteLLMTeamDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("not implemented")
}
