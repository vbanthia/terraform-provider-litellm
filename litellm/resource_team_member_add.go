package litellm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceLiteLLMTeamMemberAdd() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiteLLMTeamMemberAddCreate,
		Read:   resourceLiteLLMTeamMemberAddRead,
		Update: resourceLiteLLMTeamMemberAddUpdate,
		Delete: resourceLiteLLMTeamMemberAddDelete,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"member": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_email": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"admin",
								"user",
							}, false),
						},
					},
				},
			},
			"max_budget_in_team": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
		},
	}
}

func resourceLiteLLMTeamMemberAddCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	teamID := d.Get("team_id").(string)
	members := d.Get("member").(*schema.Set)
	maxBudget := d.Get("max_budget_in_team").(float64)

	// Convert members to the expected format
	membersList := make([]map[string]interface{}, 0, members.Len())
	for _, member := range members.List() {
		m := member.(map[string]interface{})
		memberData := map[string]interface{}{
			"role": m["role"].(string),
		}
		if userID, ok := m["user_id"].(string); ok && userID != "" {
			memberData["user_id"] = userID
		}
		if userEmail, ok := m["user_email"].(string); ok && userEmail != "" {
			memberData["user_email"] = userEmail
		}
		membersList = append(membersList, memberData)
	}

	memberData := map[string]interface{}{
		"member":             membersList,
		"team_id":            teamID,
		"max_budget_in_team": maxBudget,
	}

	log.Printf("[DEBUG] Create team members request payload: %+v", memberData)

	resp, err := MakeRequest(client, "POST", "/team/member_add", memberData)
	if err != nil {
		return fmt.Errorf("error adding team members: %v", err)
	}
	defer resp.Body.Close()

	if err := handleResponse(resp, "adding team members"); err != nil {
		return err
	}

	// Set ID as team_id since this resource manages all members for a team
	d.SetId(teamID)

	return resourceLiteLLMTeamMemberAddRead(d, m)
}

func resourceLiteLLMTeamMemberAddRead(d *schema.ResourceData, m interface{}) error {
	// The API doesn't provide a way to read specific team members
	// We'll maintain the state as is
	return nil
}

func resourceLiteLLMTeamMemberAddUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	teamID := d.Get("team_id").(string)
	maxBudget := d.Get("max_budget_in_team").(float64)

	o, n := d.GetChange("member")
	oldMembers := o.(*schema.Set)
	newMembers := n.(*schema.Set)

	// Find members to remove (in old but not in new)
	for _, member := range oldMembers.Difference(newMembers).List() {
		m := member.(map[string]interface{})
		deleteData := map[string]interface{}{
			"team_id": teamID,
		}
		if userID, ok := m["user_id"].(string); ok && userID != "" {
			deleteData["user_id"] = userID
		}
		if userEmail, ok := m["user_email"].(string); ok && userEmail != "" {
			deleteData["user_email"] = userEmail
		}

		resp, err := MakeRequest(client, "POST", "/team/member_delete", deleteData)
		if err != nil {
			return fmt.Errorf("error deleting team member: %v", err)
		}
		defer resp.Body.Close()

		if err := handleResponse(resp, "deleting team member"); err != nil {
			return err
		}
	}

	// Find members to add (in new but not in old)
	membersToAdd := newMembers.Difference(oldMembers).List()
	if len(membersToAdd) > 0 {
		membersList := make([]map[string]interface{}, 0, len(membersToAdd))
		for _, member := range membersToAdd {
			m := member.(map[string]interface{})
			memberData := map[string]interface{}{
				"role": m["role"].(string),
			}
			if userID, ok := m["user_id"].(string); ok && userID != "" {
				memberData["user_id"] = userID
			}
			if userEmail, ok := m["user_email"].(string); ok && userEmail != "" {
				memberData["user_email"] = userEmail
			}
			membersList = append(membersList, memberData)
		}

		memberData := map[string]interface{}{
			"member":             membersList,
			"team_id":            teamID,
			"max_budget_in_team": maxBudget,
		}

		log.Printf("[DEBUG] Adding new team members request payload: %+v", memberData)

		resp, err := MakeRequest(client, "POST", "/team/member_add", memberData)
		if err != nil {
			return fmt.Errorf("error adding team members: %v", err)
		}
		defer resp.Body.Close()

		if err := handleResponse(resp, "adding team members"); err != nil {
			return err
		}
	}

	return resourceLiteLLMTeamMemberAddRead(d, m)
}

func resourceLiteLLMTeamMemberAddDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	teamID := d.Get("team_id").(string)
	members := d.Get("member").(*schema.Set)

	// Delete each member
	for _, member := range members.List() {
		m := member.(map[string]interface{})
		deleteData := map[string]interface{}{
			"team_id": teamID,
		}
		if userID, ok := m["user_id"].(string); ok && userID != "" {
			deleteData["user_id"] = userID
		}
		if userEmail, ok := m["user_email"].(string); ok && userEmail != "" {
			deleteData["user_email"] = userEmail
		}

		resp, err := MakeRequest(client, "POST", "/team/member_delete", deleteData)
		if err != nil {
			return fmt.Errorf("error deleting team member: %v", err)
		}
		defer resp.Body.Close()

		if err := handleResponse(resp, "deleting team member"); err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}
