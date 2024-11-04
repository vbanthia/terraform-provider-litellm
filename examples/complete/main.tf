terraform {
  required_providers {
    litellm = {
      source  = "bitop/litellm"
      version = "1.0.0"
    }
  }
}

provider "litellm" {
  api_base = "http://localhost:4000"
  api_key  = "sk-your-api-key"
}

# Example model configuration
resource "litellm_model" "gpt4_proxy" {
  model_name          = "gpt-4-proxy"
  custom_llm_provider = "openai"
  model_api_key       = "sk-your-api-key"
  model_api_base      = "https://api.openai.com/v1"
  api_version         = "2023-05-15"
  base_model          = "gpt-4"
  tier                = "paid"
  tpm                 = 100000
  rpm                 = 1000

  # Cost configuration (per million tokens)
  input_cost_per_million_tokens  = 30.0 # $0.03 per 1k tokens = $30 per million
  output_cost_per_million_tokens = 60.0 # $0.06 per 1k tokens = $60 per million
}

# Example team configuration
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

# Example team member configuration
resource "litellm_team_member" "engineer" {
  team_id            = litellm_team.engineering.team_id
  user_id            = "user_3"
  user_email         = "engineer@example.com"
  role               = "user"
  max_budget_in_team = 200.0

  depends_on = [litellm_team.engineering]
}
