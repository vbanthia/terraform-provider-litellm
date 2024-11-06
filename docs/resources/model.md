# litellm_model Resource

Manages a LiteLLM model configuration. This resource allows you to create, update, and delete model configurations in your LiteLLM instance.

## Example Usage

```hcl
resource "litellm_model" "gpt4" {
  model_name          = "gpt-4-proxy"
  custom_llm_provider = "openai"
  model_api_key       = "sk-your-api-key"
  model_api_base      = "https://api.openai.com/v1"
  api_version         = "2023-05-15"
  base_model          = "gpt-4"
  tier                = "paid"
  mode                = "completion"
  tpm                 = 100000
  rpm                 = 1000
  
  # Cost configuration (per million tokens)
  input_cost_per_million_tokens  = 30.0    # $0.03 per 1k tokens = $30 per million
  output_cost_per_million_tokens = 60.0    # $0.06 per 1k tokens = $60 per million

  # AWS-specific configuration (if applicable)
  aws_access_key_id     = "your-aws-access-key-id"
  aws_secret_access_key = "your-aws-secret-access-key"
  aws_region_name       = "us-west-2"
}
```

## Argument Reference

The following arguments are supported:

* `model_name` - (Required) The name of the model configuration. This will be used to identify the model in API calls.

* `custom_llm_provider` - (Required) The LLM provider for this model (e.g., "openai", "anthropic", "azure", "bedrock").

* `model_api_key` - (Optional) The API key for the underlying model provider.

* `model_api_base` - (Optional) The base URL for the model provider's API.

* `api_version` - (Optional) The API version to use for the model provider.

* `base_model` - (Required) The actual model identifier from the provider (e.g., "gpt-4", "claude-2").

* `tier` - (Optional) The usage tier for this model. Valid values are "free" or "paid". Default is "free".

* `mode` - (Optional) The intended use of the model. Valid values are:
  * `completion`
  * `embeddings`
  * `image_generation`
  * `moderation`
  * `audio_transcription`

* `tpm` - (Optional) Tokens per minute limit for this model.

* `rpm` - (Optional) Requests per minute limit for this model.

* `input_cost_per_million_tokens` - (Optional) Cost per million input tokens. This will be automatically converted to the per-token cost required by the API.

* `output_cost_per_million_tokens` - (Optional) Cost per million output tokens. This will be automatically converted to the per-token cost required by the API.

### AWS-specific Configuration

* `aws_access_key_id` - (Optional) AWS access key ID for AWS-based models.

* `aws_secret_access_key` - (Optional) AWS secret access key for AWS-based models.

* `aws_region_name` - (Optional) AWS region name for AWS-based models.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the model configuration.

## Import

Model configurations can be imported using the model ID:

```shell
terraform import litellm_model.gpt4 <model-id>
```

Note: The model ID is generated when the model is created and is different from the `model_name`.
