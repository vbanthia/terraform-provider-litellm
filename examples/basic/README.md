# Basic LiteLLM Example

This example demonstrates a basic setup of the LiteLLM provider with a single model configuration.

## Usage

To run this example:

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your values:
```hcl
api_base = "http://your-litellm-instance:4000"
api_key  = "your-litellm-api-key"
openai_api_key = "your-openai-api-key"
```

2. Initialize Terraform:
```sh
terraform init
```

3. Review the plan:
```sh
terraform plan
```

4. Apply the configuration:
```sh
terraform apply
```

## Notes

- This example assumes you have a running LiteLLM instance
- Replace all API keys with your actual keys
