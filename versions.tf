terraform {
  required_version = ">= 0.13.0"
  required_providers {
    litellm = {
      source  = "bitop/litellm"
      version = "1.0.0"
    }
  }
}
