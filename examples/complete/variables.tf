variable "api_base" {
  description = "The base URL of your LiteLLM instance"
  type        = string
}

variable "api_key" {
  description = "Your LiteLLM API key"
  type        = string
  sensitive   = true
}
