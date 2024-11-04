package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	Additional map[string]interface{} `json:"additionalProp1,omitempty"`
}

type ModelRequest struct {
	ModelName     string                 `json:"model_name"`
	LiteLLMParams LiteLLMParams          `json:"litellm_params"`
	ModelInfo     ModelInfo              `json:"model_info"`
	Additional    map[string]interface{} `json:"additionalProp1,omitempty"`
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
				ForceNew: true,
			},
			"custom_llm_provider": {
				Type:     schema.TypeString,
				Optional: true,
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
				Optional: true,
			},
			"tier": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "free",
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
		},
	}
}

func isModelNotFoundError(errResp ErrorResponse) bool {
	// Check string message
	if msg, ok := errResp.Error.Message.(string); ok {
		if strings.Contains(msg, "model not found") {
			return true
		}
	}

	// Check map message
	if msgMap, ok := errResp.Error.Message.(map[string]interface{}); ok {
		if errStr, ok := msgMap["error"].(string); ok {
			if strings.Contains(errStr, "Model with id=") && strings.Contains(errStr, "not found in db") {
				return true
			}
		}
	}

	return false
}

func handleAPIResponse(resp *http.Response, reqBody interface{}) error {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
			if isModelNotFoundError(errResp) {
				return fmt.Errorf("model_not_found")
			}
		}
		reqBodyBytes, _ := json.Marshal(reqBody)
		return fmt.Errorf("API request failed: Status: %s, Response: %s, Request: %s",
			resp.Status, string(bodyBytes), string(reqBodyBytes))
	}

	return nil
}

func createOrUpdateModel(d *schema.ResourceData, m interface{}, isUpdate bool) error {
	config := m.(*ProviderConfig)

	// Convert cost per million tokens to cost per token
	inputCostPerToken := d.Get("input_cost_per_million_tokens").(float64) / 1000000.0
	outputCostPerToken := d.Get("output_cost_per_million_tokens").(float64) / 1000000.0

	modelReq := ModelRequest{
		ModelName: d.Get("model_name").(string),
		LiteLLMParams: LiteLLMParams{
			CustomLLMProvider:  d.Get("custom_llm_provider").(string),
			TPM:                d.Get("tpm").(int),
			RPM:                d.Get("rpm").(int),
			APIKey:             d.Get("model_api_key").(string),
			APIBase:            d.Get("model_api_base").(string),
			APIVersion:         d.Get("api_version").(string),
			Model:              d.Get("model_name").(string),
			InputCostPerToken:  inputCostPerToken,
			OutputCostPerToken: outputCostPerToken,
		},
		ModelInfo: ModelInfo{
			DBModel:   false,
			BaseModel: d.Get("base_model").(string),
			Tier:      d.Get("tier").(string),
		},
		Additional: make(map[string]interface{}),
	}

	if isUpdate {
		modelReq.ModelInfo.ID = d.Id()
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

	if err := handleAPIResponse(resp, modelReq); err != nil {
		if isUpdate && err.Error() == "model_not_found" {
			// If update fails because model doesn't exist, try to create it
			return createOrUpdateModel(d, m, false)
		}
		return fmt.Errorf("failed to %s model: %v", map[bool]string{true: "update", false: "create"}[isUpdate], err)
	}

	d.SetId(d.Get("model_name").(string))
	return nil
}

func resourceLiteLLMModelCreate(d *schema.ResourceData, m interface{}) error {
	return createOrUpdateModel(d, m, false)
}

func resourceLiteLLMModelRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*ProviderConfig)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/models", config.APIBase), nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-api-key", config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := handleAPIResponse(resp, nil); err != nil {
		if err.Error() == "model_not_found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read model: %v", err)
	}

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

	if err := handleAPIResponse(resp, deleteReq); err != nil {
		if err.Error() == "model_not_found" {
			// If the model is already gone, that's fine
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to delete model: %v", err)
	}

	d.SetId("")
	return nil
}
