package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLiteLLMTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLiteLLMTeamCreate,
		ReadContext:   resourceLiteLLMTeamRead,
		UpdateContext: resourceLiteLLMTeamUpdate,
		DeleteContext: resourceLiteLLMTeamDelete,
		Schema: map[string]*schema.Schema{
			"team_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"models": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"members_with_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user_email": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"blocked": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"spend": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"tpm_limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rpm_limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_budget": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"budget_duration": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceLiteLLMTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderConfig)
	client := &Client{
		APIBase: c.APIBase,
		APIKey:  c.APIKey,
	}

	// Create the request body with required fields
	team := map[string]interface{}{
		"team_alias": d.Get("team_alias").(string),
		"blocked":    d.Get("blocked").(bool),
	}

	// Handle models with default value
	if v, ok := d.GetOk("models"); ok {
		team["models"] = v.([]interface{})
	} else {
		team["models"] = []interface{}{"all-proxy-models"}
	}

	// Handle members_with_roles
	if v, ok := d.GetOk("members_with_roles"); ok {
		members := v.(*schema.Set).List()
		team["members_with_roles"] = members
	} else {
		team["members_with_roles"] = []interface{}{
			map[string]interface{}{
				"role":    "admin",
				"user_id": "default_user_id",
			},
		}
	}

	resp, err := client.CreateTeam(team)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the ID from the response
	if teamID, ok := resp["team_id"].(string); ok {
		d.SetId(teamID)
		d.Set("team_id", teamID)
	} else {
		return diag.Errorf("team_id not found in response: %v", resp)
	}

	// Set other fields from the response
	if err := setTeamData(d, resp); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceLiteLLMTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderConfig)
	client := &Client{
		APIBase: c.APIBase,
		APIKey:  c.APIKey,
	}

	team, err := client.GetTeam(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if team == nil {
		d.SetId("")
		return nil
	}

	if err := setTeamData(d, team); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setTeamData(d *schema.ResourceData, team map[string]interface{}) error {
	if v, ok := team["team_alias"].(string); ok {
		if err := d.Set("team_alias", v); err != nil {
			return err
		}
	}
	if v, ok := team["team_id"].(string); ok {
		if err := d.Set("team_id", v); err != nil {
			return err
		}
	}
	if v, ok := team["models"].([]interface{}); ok {
		if err := d.Set("models", v); err != nil {
			return err
		}
	}
	if v, ok := team["members_with_roles"].([]interface{}); ok {
		if err := d.Set("members_with_roles", v); err != nil {
			return err
		}
	}
	if v, ok := team["blocked"].(bool); ok {
		if err := d.Set("blocked", v); err != nil {
			return err
		}
	}
	if v, ok := team["spend"].(float64); ok {
		if err := d.Set("spend", v); err != nil {
			return err
		}
	}
	return nil
}

func resourceLiteLLMTeamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderConfig)
	client := &Client{
		APIBase: c.APIBase,
		APIKey:  c.APIKey,
	}

	// Only include fields that have changed
	team := map[string]interface{}{
		"team_id": d.Id(),
	}

	if d.HasChange("team_alias") {
		team["team_alias"] = d.Get("team_alias").(string)
	}

	// Send update request with only changed fields
	_, err := client.UpdateTeam(team)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLiteLLMTeamRead(ctx, d, m)
}

func resourceLiteLLMTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderConfig)
	client := &Client{
		APIBase: c.APIBase,
		APIKey:  c.APIKey,
	}

	err := client.DeleteTeam(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
