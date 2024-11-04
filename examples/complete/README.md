# Complete LiteLLM Example

This example demonstrates a complete setup of the LiteLLM provider, including:
- Model configuration
- Team setup
- Team member management

## Usage

To run this example:

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your values:
```hcl
api_base = "http://your-litellm-instance:4000"
api_key  = "your-api-key"
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
- Adjust the model configurations based on your needs
