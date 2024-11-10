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

// ModelRequest represents a request to create or update a model.
type ModelRequest struct {
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

// Key represents a LiteLLM API key.
type Key struct {
	Key                  string                 `json:"key,omitempty"`
	Models               []string               `json:"models"`
	Spend                float64                `json:"spend,omitempty"`
	MaxBudget            float64                `json:"max_budget,omitempty"`
	UserID               string                 `json:"user_id,omitempty"`
	TeamID               string                 `json:"team_id,omitempty"`
	MaxParallelRequests  int                    `json:"max_parallel_requests,omitempty"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	TPMLimit             int                    `json:"tpm_limit,omitempty"`
	RPMLimit             int                    `json:"rpm_limit,omitempty"`
	BudgetDuration       string                 `json:"budget_duration"`
	AllowedCacheControls []string               `json:"allowed_cache_controls,omitempty"`
	SoftBudget           float64                `json:"soft_budget,omitempty"`
	KeyAlias             string                 `json:"key_alias,omitempty"`
	Duration             string                 `json:"duration,omitempty"`
	Aliases              map[string]interface{} `json:"aliases,omitempty"`
	Config               map[string]interface{} `json:"config,omitempty"`
	Permissions          map[string]interface{} `json:"permissions,omitempty"`
	ModelMaxBudget       map[string]interface{} `json:"model_max_budget,omitempty"`
	ModelRPMLimit        map[string]interface{} `json:"model_rpm_limit,omitempty"`
	ModelTPMLimit        map[string]interface{} `json:"model_tpm_limit,omitempty"`
	Guardrails           []string               `json:"guardrails,omitempty"`
	Blocked              bool                   `json:"blocked"`
	Tags                 []string               `json:"tags,omitempty"`
}

// KeyResponse represents a response from the API containing key information.
type KeyResponse struct {
	Key string `json:"key"`
}
