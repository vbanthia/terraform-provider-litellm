package litellm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	endpointModelNew    = "/model/new"
	endpointModelUpdate = "/model/update"
	endpointModelInfo   = "/model/info"
	endpointModelDelete = "/model/delete"
)

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
	modelID := d.Id()
	if !isUpdate {
		modelID = uuid.New().String()
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

	endpoint := endpointModelNew
	if isUpdate {
		endpoint = endpointModelUpdate
	}

	resp, err := makeRequest(config, "POST", endpoint, modelReq)
	if err != nil {
		return fmt.Errorf("failed to %s model: %w", map[bool]string{true: "update", false: "create"}[isUpdate], err)
	}
	defer resp.Body.Close()

	_, err = handleAPIResponse(resp, modelReq)
	if err != nil {
		if isUpdate && err.Error() == "model_not_found" {
			return createOrUpdateModel(d, m, false)
		}
		return fmt.Errorf("failed to %s model: %w", map[bool]string{true: "update", false: "create"}[isUpdate], err)
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

	resp, err := makeRequest(config, "GET", fmt.Sprintf("%s?litellm_model_id=%s", endpointModelInfo, d.Id()), nil)
	if err != nil {
		return fmt.Errorf("failed to read model: %w", err)
	}
	defer resp.Body.Close()

	modelResp, err := handleAPIResponse(resp, nil)
	if err != nil {
		if err.Error() == "model_not_found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read model: %w", err)
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

	resp, err := makeRequest(config, "POST", endpointModelDelete, deleteReq)
	if err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}
	defer resp.Body.Close()

	_, err = handleAPIResponse(resp, deleteReq)
	if err != nil {
		if err.Error() == "model_not_found" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to delete model: %w", err)
	}

	d.SetId("")
	return nil
}

func makeRequest(config *ProviderConfig, method, endpoint string, body interface{}) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		req, err = http.NewRequest(method, fmt.Sprintf("%s%s", config.APIBase, endpoint), bytes.NewBuffer(jsonData))
	} else {
		req, err = http.NewRequest(method, fmt.Sprintf("%s%s", config.APIBase, endpoint), nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.APIKey)

	client := &http.Client{}
	return client.Do(req)
}
