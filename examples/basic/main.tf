terraform {
  required_providers {
    litellm = {
      source  = "bitop/litellm"
      version = "1.0.0"
    }
  }
}

provider "litellm" {
  api_base = var.api_base
  api_key  = var.api_key
}

# Simple model configuration with cost settings
resource "litellm_model" "gpt4" {
  model_name          = "gpt-4"
  custom_llm_provider = "openai"
  model_api_key       = var.openai_api_key
  model_api_base      = "https://api.openai.com/v1"
  base_model          = "gpt-4"
  tier                = "paid"

  # Cost configuration (per million tokens)
  input_cost_per_million_tokens  = 30.0 # $0.03 per 1k tokens = $30 per million
  output_cost_per_million_tokens = 60.0 # $0.06 per 1k tokens = $60 per million
}
