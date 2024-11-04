terraform {
  required_providers {
    litellm = {
      source = "local/litellm/litellm"
    }
  }
}

provider "litellm" {
  api_base = var.litellm_api_base
  api_key  = var.litellm_api_key
}

variable "litellm_api_base" {
  type        = string
  description = "LiteLLM API Base URL"
}

variable "litellm_api_key" {
  type        = string
  description = "LiteLLM API Key"
  sensitive   = true
}

resource "litellm_model" "test_model" {
  model_name          = "test-azure-gpt4"
  custom_llm_provider = "azure"
  base_model          = "gpt-4"
  model_api_key       = var.model_api_key
  model_api_base      = var.model_api_base
  api_version         = "2023-05-15"
  rpm                 = 50
  tpm                 = 80000
  tier                = "free"
}

variable "model_api_key" {
  type        = string
  description = "Model API Key"
  sensitive   = true
}

variable "model_api_base" {
  type        = string
  description = "Model API Base URL"
}
