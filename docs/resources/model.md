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
  base_model         = "gpt-4"
  tier               = "paid"
  tpm                = 100000
  rpm                = 1000
  
  # Cost configuration (per million tokens)
  input_cost_per_million_tokens  = 30.0    # $0.03 per 1k tokens = $30 per million
  output_cost_per_million_tokens = 60.0    # $0.06 per 1k tokens = $60 per million
}
```

## Argument Reference

The following arguments are supported:

* `model_name` - (Required) The name of the model configuration. This will be used to identify the model in API calls.

* `custom_llm_provider` - (Required) The LLM provider for this model. Examples include:
  * `openai`
  * `anthropic`
  * `azure`
  * `cohere`

* `model_api_key` - (Required) The API key for the underlying model provider.

* `model_api_base` - (Required) The base URL for the model provider's API.

* `api_version` - (Optional) The API version to use for the model provider.

* `base_model` - (Required) The actual model identifier from the provider (e.g., "gpt-4", "claude-2").

* `tier` - (Optional) The usage tier for this model. Valid values are:
  * `free`
  * `paid`
  Default is `free`.

* `tpm` - (Optional) Tokens per minute limit for this model.

* `rpm` - (Optional) Requests per minute limit for this model.

* `input_cost_per_million_tokens` - (Optional) Cost per million input tokens. This will be automatically converted to the per-token cost required by the API. For example:
  * Set to `30.0` for a cost of $0.03 per 1k tokens ($30 per million)
  * Set to `3.0` for a cost of $0.003 per 1k tokens ($3 per million)

* `output_cost_per_million_tokens` - (Optional) Cost per million output tokens. This will be automatically converted to the per-token cost required by the API. For example:
  * Set to `60.0` for a cost of $0.06 per 1k tokens ($60 per million)
  * Set to `6.0` for a cost of $0.006 per 1k tokens ($6 per million)

## Additional Configuration Options

The following optional parameters are also supported:

* `timeout` - (Optional) Request timeout in seconds.
* `stream_timeout` - (Optional) Streaming request timeout in seconds.
* `max_retries` - (Optional) Maximum number of retry attempts.
* `organization` - (Optional) Organization identifier for the model provider.
* `region_name` - (Optional) Region name for region-specific providers.

### Provider-Specific Options

#### Google Vertex AI
* `vertex_project` - Project ID for Google Vertex AI.
* `vertex_location` - Location for Google Vertex AI resources.
* `vertex_credentials` - Credentials for Google Vertex AI authentication.

#### AWS
* `aws_access_key_id` - AWS access key ID.
* `aws_secret_access_key` - AWS secret access key.
* `aws_region_name` - AWS region name.

#### IBM WatsonX
* `watsonx_region_name` - Region name for WatsonX services.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the model configuration.

## Import

Model configurations can be imported using the model name:

```shell
terraform import litellm_model.gpt4 gpt-4-proxy
