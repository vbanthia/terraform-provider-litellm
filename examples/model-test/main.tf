terraform {
  required_providers {
    litellm = {
      source  = "registry.terraform.io/local/litellm"
      version = "1.0.0"
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

variable "aws_access_key_id" {
  type        = string
  description = "AWS Access Key ID"
  sensitive   = true
}

variable "aws_secret_access_key" {
  type        = string
  description = "AWS Secret Access Key"
  sensitive   = true
}

variable "aws_region" {
  type        = string
  description = "AWS Region"
  default     = "us-west-2"
}

resource "litellm_model" "test_aws_bedrock_model" {
  model_name            = "claude-3.5-sonnet-v2"
  custom_llm_provider   = "bedrock"
  base_model            = "anthropic.claude-3-5-sonnet-20241022-v2:0"
  tier                  = "paid"
  aws_access_key_id     = var.aws_access_key_id
  aws_secret_access_key = var.aws_secret_access_key
  aws_region_name       = var.aws_region

  input_cost_per_million_tokens  = 4.0
  output_cost_per_million_tokens = 16.0
}

output "model_id" {
  value = litellm_model.test_aws_bedrock_model.id
}
