terraform {
  required_providers {
    litellm = {
      source  = "local/litellm"
      version = "1.0.0"
    }
  }
}

provider "litellm" {
  api_base = "http://localhost:4000"
  api_key  = "sk-your-api-key"
}

# Test model with token costs
resource "litellm_model" "cost_test" {
  model_name          = "cost-test-model"
  custom_llm_provider = "openai"
  model_api_key       = "sk-test"
  model_api_base      = "https://api.openai.com/v1"
  base_model          = "gpt-4"
  tier                = "paid"

  # Setting costs per million tokens
  # $30 per million tokens = $0.00003 per token
  input_cost_per_million_tokens = 45.0
  # $60 per million tokens = $0.00006 per token
  output_cost_per_million_tokens = 100.0
}
