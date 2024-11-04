# litellm_team_member Resource

Manages individual team member configurations in LiteLLM. This resource allows you to add, update, and remove team members with specific permissions and budget limits.

## Example Usage

```hcl
resource "litellm_team_member" "engineer" {
  team_id            = litellm_team.engineering.team_id
  user_id            = "user_3"
  user_email         = "engineer@example.com"
  role               = "user"
  max_budget_in_team = 200.0

  depends_on = [litellm_team.engineering]
}
```

## Argument Reference

The following arguments are supported:

* `team_id` - (Required) The ID of the team this member belongs to.

* `user_id` - (Required) Unique identifier for the user.

* `user_email` - (Required) Email address of the user.

* `role` - (Required) The role of the team member. Valid values are:
  * `admin`
  * `user`

* `max_budget_in_team` - (Optional) Maximum budget allocated to this team member within the team's budget.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The unique identifier for the team member configuration.

## Import

Team members can be imported using the format `team_id:user_id`:

```shell
terraform import litellm_team_member.engineer team_123456:user_3
