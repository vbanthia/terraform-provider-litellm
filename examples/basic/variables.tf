variable "api_base" {
  description = "The base URL of your LiteLLM instance"
  type        = string
}

variable "api_key" {
  description = "Your LiteLLM API key"
  type        = string
  sensitive   = true
}

variable "openai_api_key" {
  description = "Your OpenAI API key"
  type        = string
  sensitive   = true
}
