package litellm

type ProviderConfig struct {
	APIBase string
	APIKey  string
}

type LiteLLMParams struct {
	CustomLLMProvider  string                 `json:"custom_llm_provider,omitempty"`
	TPM                int                    `json:"tpm,omitempty"`
	RPM                int                    `json:"rpm,omitempty"`
	APIKey             string                 `json:"api_key,omitempty"`
	APIBase            string                 `json:"api_base,omitempty"`
	APIVersion         string                 `json:"api_version,omitempty"`
	Timeout            int                    `json:"timeout,omitempty"`
	StreamTimeout      int                    `json:"stream_timeout,omitempty"`
	MaxRetries         int                    `json:"max_retries,omitempty"`
	Model              string                 `json:"model,omitempty"`
	InputCostPerToken  float64                `json:"input_cost_per_token,omitempty"`
	OutputCostPerToken float64                `json:"output_cost_per_token,omitempty"`
	AWSAccessKeyID     string                 `json:"aws_access_key_id,omitempty"`
	AWSSecretAccessKey string                 `json:"aws_secret_access_key,omitempty"`
	AWSRegionName      string                 `json:"aws_region_name,omitempty"`
	Additional         map[string]interface{} `json:"additionalProp1,omitempty"`
}

type ModelInfo struct {
	ID         string                 `json:"id,omitempty"`
	DBModel    bool                   `json:"db_model"`
	UpdatedAt  string                 `json:"updated_at,omitempty"`
	UpdatedBy  string                 `json:"updated_by,omitempty"`
	CreatedAt  string                 `json:"created_at,omitempty"`
	CreatedBy  string                 `json:"created_by,omitempty"`
	BaseModel  string                 `json:"base_model,omitempty"`
	Tier       string                 `json:"tier,omitempty"`
	Mode       string                 `json:"mode,omitempty"`
	Additional map[string]interface{} `json:"additionalProp1,omitempty"`
}

type ModelRequest struct {
	ModelName     string                 `json:"model_name"`
	LiteLLMParams LiteLLMParams          `json:"litellm_params"`
	ModelInfo     ModelInfo              `json:"model_info"`
	Additional    map[string]interface{} `json:"additionalProp1,omitempty"`
}

type ModelResponse struct {
	ID            string        `json:"id"`
	ModelName     string        `json:"model_name"`
	LiteLLMParams LiteLLMParams `json:"litellm_params"`
	ModelInfo     ModelInfo     `json:"model_info"`
}

type ErrorResponse struct {
	Error struct {
		Message interface{} `json:"message"` // Can be string or map[string]interface{}
		Type    string      `json:"type"`
		Param   string      `json:"param"`
		Code    string      `json:"code"`
	} `json:"error"`
}
