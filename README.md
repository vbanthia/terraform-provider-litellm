# LiteLLM Terraform Provider

This Terraform provider allows you to manage LiteLLM resources through Infrastructure as Code. It provides support for managing models, teams, and team members via the LiteLLM REST API.

## Features

- Manage LiteLLM model configurations
- Create and manage teams
- Configure team members and their permissions
- Set usage limits and budgets
- Control access to specific models
- Specify model modes (e.g., completion, embeddings, image generation)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.16 (for development)

## Using the Provider

To use the LiteLLM provider in your Terraform configuration, you need to declare it in the `terraform` block:

```hcl
terraform {
  required_providers {
    litellm = {
      source  = "litellm/litellm"
      version = "~> 1.0.0"
    }
  }
}

provider "litellm" {
  api_base = var.litellm_api_base
  api_key  = var.litellm_api_key
}
```

Then, you can use the provider to manage LiteLLM resources. Here's an example of creating a model configuration:

```hcl
resource "litellm_model" "gpt4" {
  model_name          = "gpt-4-proxy"
  custom_llm_provider = "openai"
  model_api_key       = var.openai_api_key
  model_api_base      = "https://api.openai.com/v1"
  base_model          = "gpt-4"
  tier                = "paid"
  mode                = "completion"
  
  input_cost_per_million_tokens  = 30.0
  output_cost_per_million_tokens = 60.0
}
```

For full details on the `litellm_model` resource, see the [model resource documentation](docs/resources/model.md).

### Available Resources

- `litellm_model`: Manage model configurations. [Documentation](docs/resources/model.md)
- `litellm_team`: Manage teams. [Documentation](docs/resources/team.md)
- `litellm_team_member`: Manage team members. [Documentation](docs/resources/team_member.md)

## Development

### Project Structure

The project is organized as follows:

```
terraform-provider-litellm/
├── litellm/
│   ├── provider.go
│   ├── resource_model.go
│   ├── resource_model_crud.go
│   ├── types.go
│   └── utils.go
├── main.go
├── go.mod
├── go.sum
├── Makefile
└── ...
```

### Building the Provider

1. Clone the repository
```sh
git clone https://github.com/your-username/terraform-provider-litellm.git
```

2. Enter the repository directory
```sh
cd terraform-provider-litellm
```

3. Build and install the provider
```sh
make install
```

### Development Commands

The Makefile provides several useful commands for development:

- `make build`: Builds the provider
- `make install`: Builds and installs the provider
- `make test`: Runs the test suite
- `make fmt`: Formats the code
- `make vet`: Runs go vet
- `make lint`: Runs golangci-lint
- `make clean`: Removes build artifacts and installed provider

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
