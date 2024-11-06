package litellm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceLiteLLMTeamMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiteLLMTeamMemberCreate,
		Read:   resourceLiteLLMTeamMemberRead,
		Update: resourceLiteLLMTeamMemberUpdate,
		Delete: resourceLiteLLMTeamMemberDelete,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"org_admin",
					"internal_user",
					"internal_user_viewer",
					"admin",
					"user",
				}, false),
			},
			"max_budget_in_team": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
		},
	}
}

func resourceLiteLLMTeamMemberCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	memberData := map[string]interface{}{
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

	jsonData, err := json.Marshal(memberData)
	if err != nil {
		return fmt.Errorf("error marshalling team member data: %v", err)
	}

	log.Printf("[DEBUG] Create team member request payload: %s", string(jsonData))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/team/member_add", config.APIBase), bytes.NewBuffer(jsonData))
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
	log.Printf("[DEBUG] Create team member response: Status: %s, Body: %s", resp.Status, string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating team member: %s - %s", resp.Status, string(body))
	}

	// Set a composite ID since there's no specific member ID returned
	d.SetId(fmt.Sprintf("%s:%s", d.Get("team_id").(string), d.Get("user_id").(string)))

	log.Printf("[INFO] Team member created with ID: %s", d.Id())

	return resourceLiteLLMTeamMemberRead(d, m)
}

func resourceLiteLLMTeamMemberRead(d *schema.ResourceData, m interface{}) error {
	// There's no specific endpoint to read a single team member
	// We might need to read the entire team and find the member
	// For now, we'll just return the data we have in the state
	log.Printf("[INFO] Reading team member with ID: %s", d.Id())
	return nil
}

func resourceLiteLLMTeamMemberUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	updateData := map[string]interface{}{
		"user_id":            d.Get("user_id").(string),
		"user_email":         d.Get("user_email").(string),
		"team_id":            d.Get("team_id").(string),
		"max_budget_in_team": d.Get("max_budget_in_team").(float64),
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return fmt.Errorf("error marshalling team member update data: %v", err)
	}

	log.Printf("[DEBUG] Update team member request payload: %s", string(jsonData))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/team/member_update", config.APIBase), bytes.NewBuffer(jsonData))
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
	log.Printf("[DEBUG] Update team member response: Status: %s, Body: %s", resp.Status, string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error updating team member: %s - %s", resp.Status, string(body))
	}

	log.Printf("[INFO] Successfully updated team member with ID: %s", d.Id())

	return resourceLiteLLMTeamMemberRead(d, m)
}

func resourceLiteLLMTeamMemberDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	deleteData := map[string]interface{}{
		"user_id":    d.Get("user_id").(string),
		"user_email": d.Get("user_email").(string),
		"team_id":    d.Get("team_id").(string),
	}

	jsonData, err := json.Marshal(deleteData)
	if err != nil {
		return fmt.Errorf("error marshalling team member delete data: %v", err)
	}

	log.Printf("[DEBUG] Delete team member request payload: %s", string(jsonData))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/team/member_delete", config.APIBase), bytes.NewBuffer(jsonData))
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
	log.Printf("[DEBUG] Delete team member response: Status: %s, Body: %s", resp.Status, string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting team member: %s - %s", resp.Status, string(body))
	}

	log.Printf("[INFO] Successfully deleted team member with ID: %s", d.Id())

	d.SetId("")
	return nil
}
