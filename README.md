# LiteLLM Terraform Provider

[![Build Status](https://github.com/bitop/terraform-provider-litellm/workflows/tests/badge.svg)](https://github.com/bitop/terraform-provider-litellm/actions)

This Terraform provider allows you to manage LiteLLM resources through Infrastructure as Code. It provides support for managing models, teams, and team members via the LiteLLM REST API.

## Documentation

Full documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/bitop/litellm/latest/docs).

## Features

- Manage LiteLLM model configurations
- Create and manage teams
- Configure team members and their permissions
- Set usage limits and budgets
- Control access to specific models

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.16 (for development)

## Using the Provider

```hcl
terraform {
  required_providers {
    litellm = {
      source  = "bitop/litellm"
      version = "1.0.0"
    }
  }
}

provider "litellm" {
  api_base = "http://your-litellm-instance:4000"
  api_key  = "your-api-key"
}

# Example model configuration
resource "litellm_model" "gpt4" {
  model_name          = "gpt-4-proxy"
  custom_llm_provider = "openai"
  model_api_key       = "sk-your-api-key"
  model_api_base      = "https://api.openai.com/v1"
  api_version         = "2023-05-15"
  base_model         = "gpt-4"
  tier               = "paid"
  tpm                = 100000
  rpm                = 1000
}
```

## Development

### Building the Provider

1. Clone the repository
```sh
git clone https://github.com/bitop/terraform-provider-litellm.git
```

2. Enter the repository directory
```sh
cd terraform-provider-litellm
```

3. Build the provider
```sh
make install
```

### Testing

To run the tests:

```sh
make test
```

### Contributing

Contributions are welcome! Please read our [contributing guidelines](CONTRIBUTING.md) first.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
