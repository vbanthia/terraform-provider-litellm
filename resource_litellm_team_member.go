package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLiteLLMTeamMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLiteLLMTeamMemberCreate,
		ReadContext:   resourceLiteLLMTeamMemberRead,
		UpdateContext: resourceLiteLLMTeamMemberUpdate,
		DeleteContext: resourceLiteLLMTeamMemberDelete,
		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "member",
			},
			"max_budget_in_team": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceLiteLLMTeamMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderConfig)
	client := &Client{
		APIBase: c.APIBase,
		APIKey:  c.APIKey,
	}

	member := map[string]interface{}{
		"member": []map[string]interface{}{
			{
				"role":       d.Get("role").(string),
				"user_id":    d.Get("user_id").(string),
				"user_email": d.Get("user_email").(string),
			},
		},
		"team_id":            d.Get("team_id").(string),
		"max_budget_in_team": d.Get("max_budget_in_team").(float64),
	}

	_, err := client.sendRequest("POST", "/team/member_add", member)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("team_id").(string) + ":" + d.Get("user_id").(string))

	return resourceLiteLLMTeamMemberRead(ctx, d, m)
}

func resourceLiteLLMTeamMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// The API doesn't provide a direct way to read a single team member
	// We'll assume if the team exists, the member exists
	return nil
}

func resourceLiteLLMTeamMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderConfig)
	client := &Client{
		APIBase: c.APIBase,
		APIKey:  c.APIKey,
	}

	member := map[string]interface{}{
		"user_id":            d.Get("user_id").(string),
		"user_email":         d.Get("user_email").(string),
		"team_id":            d.Get("team_id").(string),
		"max_budget_in_team": d.Get("max_budget_in_team").(float64),
	}

	_, err := client.sendRequest("POST", "/team/member_update", member)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLiteLLMTeamMemberRead(ctx, d, m)
}

func resourceLiteLLMTeamMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderConfig)
	client := &Client{
		APIBase: c.APIBase,
		APIKey:  c.APIKey,
	}

	member := map[string]interface{}{
		"user_id":    d.Get("user_id").(string),
		"user_email": d.Get("user_email").(string),
		"team_id":    d.Get("team_id").(string),
	}

	_, err := client.sendRequest("POST", "/team/member_delete", member)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
