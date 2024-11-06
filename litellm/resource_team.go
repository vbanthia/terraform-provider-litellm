package litellm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLiteLLMTeam() *schema.Resource {
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

	// Generate a UUID for the team
	teamID := uuid.New().String()

	teamData := map[string]interface{}{
		"team_id":    teamID,
		"team_alias": d.Get("team_alias").(string),
		"metadata":   d.Get("metadata").(map[string]interface{}),
		"tpm_limit":  d.Get("tpm_limit").(int),
		"rpm_limit":  d.Get("rpm_limit").(int),
		"max_budget": d.Get("max_budget").(float64),
		"models":     d.Get("models").([]interface{}),
		"blocked":    d.Get("blocked").(bool),
	}

	// Only include organization_id if it's set
	if v, ok := d.GetOk("organization_id"); ok {
		teamData["organization_id"] = v.(string)
	}

	// Only include budget_duration if it's set
	if v, ok := d.GetOk("budget_duration"); ok {
		teamData["budget_duration"] = v.(string)
	}

	jsonData, err := json.Marshal(teamData)
	if err != nil {
		return fmt.Errorf("error marshalling team data: %v", err)
	}

	log.Printf("[DEBUG] Create team request payload: %s", string(jsonData))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/team/new", config.APIBase), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Create team response: Status: %s, Body: %s", resp.Status, string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating team: %s - %s", resp.Status, string(body))
	}

	d.SetId(teamID)

	log.Printf("[INFO] Team created with ID: %s", teamID)

	return resourceLiteLLMTeamRead(d, m)
}

func resourceLiteLLMTeamRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	log.Printf("[INFO] Reading team with ID: %s", d.Id())

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/team/info?team_id=%s", config.APIBase, d.Id()), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("x-api-key", config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Read team response: Status: %s, Body: %s", resp.Status, string(body))

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[WARN] Team with ID %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error reading team: %s - %s", resp.Status, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	d.Set("team_alias", result["team_alias"])
	d.Set("organization_id", result["organization_id"])
	d.Set("metadata", result["metadata"])
	d.Set("tpm_limit", result["tpm_limit"])
	d.Set("rpm_limit", result["rpm_limit"])
	d.Set("max_budget", result["max_budget"])
	d.Set("budget_duration", result["budget_duration"])
	d.Set("models", result["models"])
	d.Set("blocked", result["blocked"])

	log.Printf("[INFO] Successfully read team with ID: %s", d.Id())

	return nil
}

func resourceLiteLLMTeamUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	teamData := map[string]interface{}{
		"team_id":    d.Id(),
		"team_alias": d.Get("team_alias").(string),
		"metadata":   d.Get("metadata").(map[string]interface{}),
		"tpm_limit":  d.Get("tpm_limit").(int),
		"rpm_limit":  d.Get("rpm_limit").(int),
		"max_budget": d.Get("max_budget").(float64),
		"models":     d.Get("models").([]interface{}),
		"blocked":    d.Get("blocked").(bool),
	}

	// Only include organization_id if it's set
	if v, ok := d.GetOk("organization_id"); ok {
		teamData["organization_id"] = v.(string)
	}

	// Only include budget_duration if it's set
	if v, ok := d.GetOk("budget_duration"); ok {
		teamData["budget_duration"] = v.(string)
	}

	jsonData, err := json.Marshal(teamData)
	if err != nil {
		return fmt.Errorf("error marshalling team data: %v", err)
	}

	log.Printf("[DEBUG] Update team request payload: %s", string(jsonData))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/team/update", config.APIBase), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Update team response: Status: %s, Body: %s", resp.Status, string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error updating team: %s - %s", resp.Status, string(body))
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

	jsonData, err := json.Marshal(deleteData)
	if err != nil {
		return fmt.Errorf("error marshalling delete data: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/team/delete", config.APIBase), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Delete team response: Status: %s, Body: %s", resp.Status, string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting team: %s - %s", resp.Status, string(body))
	}

	log.Printf("[INFO] Successfully deleted team with ID: %s", d.Id())

	d.SetId("")
	return nil
}
