# LiteLLM Terraform Provider

This Terraform provider allows you to manage LiteLLM resources through Infrastructure as Code. It provides support for managing models, teams, and team members via the LiteLLM REST API.

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
      source  = "local/litellm/litellm"
      version = "1.0.0"
    }
  }
}

provider "litellm" {
  api_base = var.litellm_api_base
  api_key  = var.litellm_api_key
}

# Example model configuration for AWS Bedrock
resource "litellm_model" "claude_aws_bedrock" {
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
```

## Development

### Building the Provider

1. Clone the repository
```sh
git clone https://github.com/your-username/terraform-provider-litellm.git
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

## Notes

- Always use environment variables or secure secret management solutions to handle sensitive information like API keys and AWS credentials.
- Refer to the `examples/` directory for more detailed usage examples.
- Make sure to keep your provider version updated for the latest features and bug fixes.
