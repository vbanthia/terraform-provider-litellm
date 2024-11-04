# Token Cost Test Example

This example validates the token cost conversion functionality of the LiteLLM provider. It demonstrates setting costs per million tokens and verifies that they are correctly converted to per-token costs when sent to the API.

## Expected Behavior

The provider should convert:
- `input_cost_per_million_tokens = 30.0` to `input_cost_per_token = 0.00003`
- `output_cost_per_million_tokens = 60.0` to `output_cost_per_token = 0.00006`

## Usage

```bash
terraform init
terraform plan
terraform apply
```

After applying, you can verify the costs were set correctly by checking the model configuration in LiteLLM.
