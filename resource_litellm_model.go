package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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

func resourceLiteLLMModel() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiteLLMModelCreate,
		Read:   resourceLiteLLMModelRead,
		Update: resourceLiteLLMModelUpdate,
		Delete: resourceLiteLLMModelDelete,

		Schema: map[string]*schema.Schema{
			"model_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_llm_provider": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tpm": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rpm": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"model_api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"model_api_base": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"base_model": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tier": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "free",
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"completion",
					"embeddings",
					"image_generation",
					"moderation",
					"audio_transcription",
				}, false),
			},
			"input_cost_per_million_tokens": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0.0,
			},
			"output_cost_per_million_tokens": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0.0,
			},
			"aws_access_key_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"aws_secret_access_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"aws_region_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func isModelNotFoundError(errResp ErrorResponse) bool {
	if msg, ok := errResp.Error.Message.(string); ok {
		if strings.Contains(msg, "model not found") {
			return true
		}
	}

	if msgMap, ok := errResp.Error.Message.(map[string]interface{}); ok {
		if errStr, ok := msgMap["error"].(string); ok {
			if strings.Contains(errStr, "Model with id=") && strings.Contains(errStr, "not found in db") {
				return true
			}
		}
	}

	return false
}

func handleAPIResponse(resp *http.Response, reqBody interface{}) (*ModelResponse, error) {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
			if isModelNotFoundError(errResp) {
				return nil, fmt.Errorf("model_not_found")
			}
		}
		reqBodyBytes, _ := json.Marshal(reqBody)
		return nil, fmt.Errorf("API request failed: Status: %s, Response: %s, Request: %s",
			resp.Status, string(bodyBytes), string(reqBodyBytes))
	}

	var modelResp ModelResponse
	if err := json.Unmarshal(bodyBytes, &modelResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &modelResp, nil
}

func createOrUpdateModel(d *schema.ResourceData, m interface{}, isUpdate bool) error {
	config := m.(*ProviderConfig)

	// Convert cost per million tokens to cost per token
	inputCostPerToken := d.Get("input_cost_per_million_tokens").(float64) / 1000000.0
	outputCostPerToken := d.Get("output_cost_per_million_tokens").(float64) / 1000000.0

	// Construct the model name in the format "custom_llm_provider/base_model"
	customLLMProvider := d.Get("custom_llm_provider").(string)
	baseModel := d.Get("base_model").(string)
	modelName := fmt.Sprintf("%s/%s", customLLMProvider, baseModel)

	// Generate a UUID for new models
	var modelID string
	if !isUpdate {
		modelID = uuid.New().String()
	} else {
		modelID = d.Id()
	}

	modelReq := ModelRequest{
		ModelName: d.Get("model_name").(string),
		LiteLLMParams: LiteLLMParams{
			CustomLLMProvider:  customLLMProvider,
			TPM:                d.Get("tpm").(int),
			RPM:                d.Get("rpm").(int),
			APIKey:             d.Get("model_api_key").(string),
			APIBase:            d.Get("model_api_base").(string),
			APIVersion:         d.Get("api_version").(string),
			Model:              modelName,
			InputCostPerToken:  inputCostPerToken,
			OutputCostPerToken: outputCostPerToken,
			AWSAccessKeyID:     d.Get("aws_access_key_id").(string),
			AWSSecretAccessKey: d.Get("aws_secret_access_key").(string),
			AWSRegionName:      d.Get("aws_region_name").(string),
		},
		ModelInfo: ModelInfo{
			ID:        modelID,
			DBModel:   true,
			BaseModel: baseModel,
			Tier:      d.Get("tier").(string),
			Mode:      d.Get("mode").(string),
		},
		Additional: make(map[string]interface{}),
	}

	jsonData, err := json.Marshal(modelReq)
	if err != nil {
		return err
	}

	endpoint := "/model/new"
	if isUpdate {
		endpoint = "/model/update"
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", config.APIBase, endpoint), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = handleAPIResponse(resp, modelReq)
	if err != nil {
		if isUpdate && err.Error() == "model_not_found" {
			return createOrUpdateModel(d, m, false)
		}
		return fmt.Errorf("failed to %s model: %v", map[bool]string{true: "update", false: "create"}[isUpdate], err)
	}

	d.SetId(modelID)

	// Read back the resource to ensure the state is consistent
	return resourceLiteLLMModelRead(d, m)
}

func resourceLiteLLMModelCreate(d *schema.ResourceData, m interface{}) error {
	return createOrUpdateModel(d, m, false)
}

func resourceLiteLLMModelRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/model/info?litellm_model_id=%s", config.APIBase, d.Id()), nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	modelResp, err := handleAPIResponse(resp, nil)
	if err != nil {
		if err.Error() == "model_not_found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read model: %v", err)
	}

	// Update the state with values from the response
	d.Set("model_name", modelResp.ModelName)
	d.Set("custom_llm_provider", modelResp.LiteLLMParams.CustomLLMProvider)
	d.Set("tpm", modelResp.LiteLLMParams.TPM)
	d.Set("rpm", modelResp.LiteLLMParams.RPM)
	d.Set("model_api_base", modelResp.LiteLLMParams.APIBase)
	d.Set("api_version", modelResp.LiteLLMParams.APIVersion)
	d.Set("base_model", modelResp.ModelInfo.BaseModel)
	d.Set("tier", modelResp.ModelInfo.Tier)
	d.Set("mode", modelResp.ModelInfo.Mode)
	d.Set("aws_access_key_id", modelResp.LiteLLMParams.AWSAccessKeyID)
	d.Set("aws_secret_access_key", modelResp.LiteLLMParams.AWSSecretAccessKey)
	d.Set("aws_region_name", modelResp.LiteLLMParams.AWSRegionName)

	return nil
}

func resourceLiteLLMModelUpdate(d *schema.ResourceData, m interface{}) error {
	return createOrUpdateModel(d, m, true)
}

func resourceLiteLLMModelDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	deleteReq := struct {
		ID string `json:"id"`
	}{
		ID: d.Id(),
	}

	jsonData, err := json.Marshal(deleteReq)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/model/delete", config.APIBase), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = handleAPIResponse(resp, deleteReq)
	if err != nil {
		if err.Error() == "model_not_found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to delete model: %v", err)
	}

	d.SetId("")
	return nil
}
