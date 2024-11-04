# LiteLLM Model Test Example

This example demonstrates how to use the LiteLLM Terraform provider to manage LiteLLM models.

## Usage

1. Copy the example variables file:
```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:
```hcl
litellm_api_base = "https://your-litellm-api-base"
litellm_api_key  = "your-litellm-api-key"
model_api_key    = "your-model-api-key"
model_api_base   = "https://your-model-api-base"
```

3. Initialize Terraform:
```bash
terraform init
```

4. Apply the configuration:
```bash
terraform apply
```

## Variables

* `litellm_api_base` - The base URL of your LiteLLM API
* `litellm_api_key` - Your LiteLLM API key
* `model_api_key` - The API key for the model (e.g., Azure OpenAI API key)
* `model_api_base` - The base URL for the model (e.g., Azure OpenAI endpoint)

## Notes

* Never commit `terraform.tfvars` to version control as it contains sensitive information
* The example `.gitignore` is configured to prevent committing sensitive files
* Use environment variables or a secrets management solution in production environments
