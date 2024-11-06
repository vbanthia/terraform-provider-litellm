package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLiteLLMTeamMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiteLLMTeamMemberCreate,
		Read:   resourceLiteLLMTeamMemberRead,
		Update: resourceLiteLLMTeamMemberUpdate,
		Delete: resourceLiteLLMTeamMemberDelete,

		Schema: map[string]*schema.Schema{
			"member_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLiteLLMTeamMemberCreate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("not implemented")
}

func resourceLiteLLMTeamMemberRead(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("not implemented")
}

func resourceLiteLLMTeamMemberUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("not implemented")
}

func resourceLiteLLMTeamMemberDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("not implemented")
}
