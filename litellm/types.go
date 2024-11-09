package litellm

// ProviderConfig holds the configuration for the LiteLLM provider.
type ProviderConfig struct {
	APIBase string
	APIKey  string
}

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	Error struct {
		Message interface{} `json:"message"`
	} `json:"error"`
}

// ModelResponse represents a response from the API containing model information.
type ModelResponse struct {
	ModelName     string                 `json:"model_name"`
	LiteLLMParams LiteLLMParams          `json:"litellm_params"`
	ModelInfo     ModelInfo              `json:"model_info"`
	Additional    map[string]interface{} `json:"additional"`
}

// TeamResponse represents a response from the API containing team information.
type TeamResponse struct {
	TeamID         string                 `json:"team_id"`
	TeamAlias      string                 `json:"team_alias"`
	OrganizationID string                 `json:"organization_id"`
	Metadata       map[string]interface{} `json:"metadata"`
	TPMLimit       int                    `json:"tpm_limit"`
	RPMLimit       int                    `json:"rpm_limit"`
	MaxBudget      float64                `json:"max_budget"`
	BudgetDuration string                 `json:"budget_duration"`
	Models         []string               `json:"models"`
	Blocked        bool                   `json:"blocked"`
}

// ModelRequest represents a request to create or update a model.
type ModelRequest struct {
	ModelName     string                 `json:"model_name"`
	LiteLLMParams LiteLLMParams          `json:"litellm_params"`
	ModelInfo     ModelInfo              `json:"model_info"`
	Additional    map[string]interface{} `json:"additional"`
}

// LiteLLMParams represents the parameters for LiteLLM.
type LiteLLMParams struct {
	CustomLLMProvider  string  `json:"custom_llm_provider"`
	TPM                int     `json:"tpm"`
	RPM                int     `json:"rpm"`
	APIKey             string  `json:"api_key"`
	APIBase            string  `json:"api_base"`
	APIVersion         string  `json:"api_version"`
	Model              string  `json:"model"`
	InputCostPerToken  float64 `json:"input_cost_per_token"`
	OutputCostPerToken float64 `json:"output_cost_per_token"`
	AWSAccessKeyID     string  `json:"aws_access_key_id"`
	AWSSecretAccessKey string  `json:"aws_secret_access_key"`
	AWSRegionName      string  `json:"aws_region_name"`
}

// ModelInfo represents information about a model.
type ModelInfo struct {
	ID        string `json:"id"`
	DBModel   bool   `json:"db_model"`
	BaseModel string `json:"base_model"`
	Tier      string `json:"tier"`
	Mode      string `json:"mode"`
}

// Add any other type definitions here as needed
