package litellm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	endpointTeamNew    = "/team/new"
	endpointTeamInfo   = "/team/info"
	endpointTeamUpdate = "/team/update"
	endpointTeamDelete = "/team/delete"
)

func ResourceLiteLLMTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiteLLMTeamCreate,
		Read:   resourceLiteLLMTeamRead,
		Update: resourceLiteLLMTeamUpdate,
		Delete: resourceLiteLLMTeamDelete,

		Schema: map[string]*schema.Schema{
			"team_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"models": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"blocked": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceLiteLLMTeamCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	teamID := uuid.New().String()
	teamData := buildTeamData(d, teamID)

	log.Printf("[DEBUG] Create team request payload: %+v", teamData)

	resp, err := MakeRequest(config, "POST", endpointTeamNew, teamData)
	if err != nil {
		return fmt.Errorf("error creating team: %w", err)
	}
	defer resp.Body.Close()

	if err := handleResponse(resp, "creating team"); err != nil {
		return err
	}

	d.SetId(teamID)
	log.Printf("[INFO] Team created with ID: %s", teamID)

	return resourceLiteLLMTeamRead(d, m)
}

func resourceLiteLLMTeamRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	log.Printf("[INFO] Reading team with ID: %s", d.Id())

	resp, err := MakeRequest(config, "GET", fmt.Sprintf("%s?team_id=%s", endpointTeamInfo, d.Id()), nil)
	if err != nil {
		return fmt.Errorf("error reading team: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Team with ID %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	var teamResp TeamResponse
	if err := json.NewDecoder(resp.Body).Decode(&teamResp); err != nil {
		return fmt.Errorf("error decoding team info response: %w", err)
	}

	d.Set("team_alias", teamResp.TeamAlias)
	d.Set("organization_id", teamResp.OrganizationID)
	d.Set("metadata", teamResp.Metadata)
	d.Set("tpm_limit", teamResp.TPMLimit)
	d.Set("rpm_limit", teamResp.RPMLimit)
	d.Set("max_budget", teamResp.MaxBudget)
	d.Set("budget_duration", teamResp.BudgetDuration)
	d.Set("models", teamResp.Models)
	d.Set("blocked", teamResp.Blocked)

	log.Printf("[INFO] Successfully read team with ID: %s", d.Id())
	return nil
}

func resourceLiteLLMTeamUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	teamData := buildTeamData(d, d.Id())
	log.Printf("[DEBUG] Update team request payload: %+v", teamData)

	resp, err := MakeRequest(config, "POST", endpointTeamUpdate, teamData)
	if err != nil {
		return fmt.Errorf("error updating team: %w", err)
	}
	defer resp.Body.Close()

	if err := handleResponse(resp, "updating team"); err != nil {
		return err
	}

	log.Printf("[INFO] Successfully updated team with ID: %s", d.Id())
	return resourceLiteLLMTeamRead(d, m)
}

func resourceLiteLLMTeamDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	log.Printf("[INFO] Deleting team with ID: %s", d.Id())

	deleteData := map[string]interface{}{
		"team_ids": []string{d.Id()},
	}

	resp, err := MakeRequest(config, "POST", endpointTeamDelete, deleteData)
	if err != nil {
		return fmt.Errorf("error deleting team: %w", err)
	}
	defer resp.Body.Close()

	if err := handleResponse(resp, "deleting team"); err != nil {
		return err
	}

	log.Printf("[INFO] Successfully deleted team with ID: %s", d.Id())
	d.SetId("")
	return nil
}

func buildTeamData(d *schema.ResourceData, teamID string) map[string]interface{} {
	teamData := map[string]interface{}{
		"team_id":    teamID,
		"team_alias": d.Get("team_alias").(string),
	}

	for _, key := range []string{"organization_id", "metadata", "tpm_limit", "rpm_limit", "max_budget", "budget_duration", "models", "blocked"} {
		if v, ok := d.GetOk(key); ok {
			teamData[key] = v
		}
	}

	return teamData
}

func handleResponse(resp *http.Response, action string) error {
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error %s: %s - %s", action, resp.Status, string(body))
	}
	return nil
}
