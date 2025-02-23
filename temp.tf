# Get gatorlinks from AD group
data "ldap_group" "ufai-ict-lps" {
  ou    = "OU=groups,OU=toolkit,OU=NaviGator,OU=Services,OU=UFIT,OU=Departments,OU=UF,DC=ad,DC=ufl,DC=edu"
  name  = "ufai-ict-lps"
  scope = 2
}

# Create Team in LiteLLM
resource "litellm_team" "ufai-ict-lps" {
  team_alias = "ufai-ict-lps"
  models     = [
    "mixtral-8x7b-instruct",
    "mistral-7b-instruct",
    "gpt-3.5-turbo",
    "gpt-4-turbo"
  ]


  metadata = {
    customer_name    = "ufai-ict-lps",
    customer_number  = "00001234",
    customer_contact = "ufit@mail.ufl.edu"
    ict              = true
  }
  max_budget      = 150.0
  blocked         = false
  budget_duration = "30d"
}

locals {
  team_members = {
    for member in data.ldap_group.ufai-ict-lps.members_names : member => {
      user_id    = "${member}@ufl.edu"
      user_email = "${member}@ufl.edu"
      role       = "user"
    }
  }
}

# Bang
resource "litellm_team_member_add" "ufai-ict-lps" {
  team_id = litellm_team.ufai-ict-lps.id
  
  dynamic "member" {
    for_each = local.team_members
    content {
      user_id    = lookup(member.value, "user_id", null)
      user_email = lookup(member.value, "user_email", null)
      role       = lookup(member.value, "role", null)
    }
  }
  max_budget_in_team = 100
}