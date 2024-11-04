# litellm_team Resource

Manages a team configuration in LiteLLM. Teams allow you to group users and manage their access to models and usage limits.

## Example Usage

```hcl
resource "litellm_team" "engineering" {
  team_alias = "engineering-team"
  models     = ["gpt-4-proxy", "claude-2"]

  members_with_roles {
    role       = "admin"
    user_id    = "user_1"
    user_email = "admin@example.com"
  }

  members_with_roles {
    role       = "user"
    user_id    = "user_2"
    user_email = "member@example.com"
  }

  blocked         = false
  tpm_limit       = 500000
  rpm_limit       = 5000
  max_budget      = 1000.0
  budget_duration = "monthly"
}
```

## Argument Reference

The following arguments are supported:

* `team_alias` - (Required) A human-readable identifier for the team.

* `models` - (Required) List of model names that this team can access.

* `members_with_roles` - (Required) List of team members and their roles. Each block supports:
  * `role` - (Required) The role of the team member. Valid values are:
    * `admin`
    * `user`
  * `user_id` - (Required) Unique identifier for the user.
  * `user_email` - (Required) Email address of the user.

* `blocked` - (Optional) Whether the team is blocked from making requests. Default is `false`.

* `tpm_limit` - (Optional) Team-wide tokens per minute limit.

* `rpm_limit` - (Optional) Team-wide requests per minute limit.

* `max_budget` - (Optional) Maximum budget allocated to the team.

* `budget_duration` - (Optional) Duration for the budget cycle. Valid values are:
  * `daily`
  * `weekly`
  * `monthly`
  * `yearly`

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `team_id` - The unique identifier for the team.

## Import

Teams can be imported using the team ID:

```shell
terraform import litellm_team.engineering team_123456
