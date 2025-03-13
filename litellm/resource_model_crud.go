package litellm

import (
	"fmt"

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
	client, ok := m.(*Client)
	if !ok {
		return fmt.Errorf("invalid type assertion for client")
	}

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

	// Create thinking configuration if enabled
	var thinking map[string]interface{}
	if d.Get("thinking_enabled").(bool) {
		thinking = map[string]interface{}{
			"type":          "enabled",
			"budget_tokens": d.Get("thinking_budget_tokens").(int),
		}
	}

	modelReq := ModelRequest{
		ModelName: d.Get("model_name").(string),
		LiteLLMParams: LiteLLMParams{
			CustomLLMProvider:   customLLMProvider,
			TPM:                 d.Get("tpm").(int),
			RPM:                 d.Get("rpm").(int),
			APIKey:              d.Get("model_api_key").(string),
			APIBase:             d.Get("model_api_base").(string),
			APIVersion:          d.Get("api_version").(string),
			Model:               modelName,
			InputCostPerToken:   inputCostPerToken,
			OutputCostPerToken:  outputCostPerToken,
			InputCostPerPixel:   d.Get("input_cost_per_pixel").(float64),
			OutputCostPerPixel:  d.Get("output_cost_per_pixel").(float64),
			InputCostPerSecond:  d.Get("input_cost_per_second").(float64),
			OutputCostPerSecond: d.Get("output_cost_per_second").(float64),
			AWSAccessKeyID:      d.Get("aws_access_key_id").(string),
			AWSSecretAccessKey:  d.Get("aws_secret_access_key").(string),
			AWSRegionName:       d.Get("aws_region_name").(string),
			VertexProject:       d.Get("vertex_project").(string),
			VertexLocation:      d.Get("vertex_location").(string),
			VertexCredentials:   d.Get("vertex_credentials").(string),
			ReasoningEffort:     d.Get("reasoning_effort").(string),
			Thinking:            thinking,
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

	resp, err := MakeRequest(client, "POST", endpoint, modelReq)
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
	client, ok := m.(*Client)
	if !ok {
		return fmt.Errorf("invalid type assertion for client")
	}

	resp, err := MakeRequest(client, "GET", fmt.Sprintf("%s?litellm_model_id=%s", endpointModelInfo, d.Id()), nil)
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

	// Update the state with values from the response or fall back to the data passed in during creation
	d.Set("model_name", GetStringValue(modelResp.ModelName, d.Get("model_name").(string)))
	d.Set("custom_llm_provider", GetStringValue(modelResp.LiteLLMParams.CustomLLMProvider, d.Get("custom_llm_provider").(string)))
	d.Set("tpm", GetIntValue(modelResp.LiteLLMParams.TPM, d.Get("tpm").(int)))
	d.Set("rpm", GetIntValue(modelResp.LiteLLMParams.RPM, d.Get("rpm").(int)))
	d.Set("model_api_base", GetStringValue(modelResp.LiteLLMParams.APIBase, d.Get("model_api_base").(string)))
	d.Set("api_version", GetStringValue(modelResp.LiteLLMParams.APIVersion, d.Get("api_version").(string)))
	d.Set("base_model", GetStringValue(modelResp.ModelInfo.BaseModel, d.Get("base_model").(string)))
	d.Set("tier", GetStringValue(modelResp.ModelInfo.Tier, d.Get("tier").(string)))
	d.Set("mode", GetStringValue(modelResp.ModelInfo.Mode, d.Get("mode").(string)))

	// Store sensitive information
	d.Set("model_api_key", d.Get("model_api_key"))
	d.Set("aws_access_key_id", d.Get("aws_access_key_id"))
	d.Set("aws_secret_access_key", d.Get("aws_secret_access_key"))
	d.Set("aws_region_name", GetStringValue(modelResp.LiteLLMParams.AWSRegionName, d.Get("aws_region_name").(string)))

	// Store cost information
	d.Set("input_cost_per_million_tokens", d.Get("input_cost_per_million_tokens"))
	d.Set("output_cost_per_million_tokens", d.Get("output_cost_per_million_tokens"))

	// Handle thinking configuration
	if modelResp.LiteLLMParams.Thinking != nil {
		if thinkingType, ok := modelResp.LiteLLMParams.Thinking["type"].(string); ok && thinkingType == "enabled" {
			d.Set("thinking_enabled", true)
			if budgetTokens, ok := modelResp.LiteLLMParams.Thinking["budget_tokens"].(float64); ok {
				d.Set("thinking_budget_tokens", int(budgetTokens))
			}
		}
	} else {
		d.Set("thinking_enabled", false)
	}

	return nil
}

func resourceLiteLLMModelUpdate(d *schema.ResourceData, m interface{}) error {
	return createOrUpdateModel(d, m, true)
}

func resourceLiteLLMModelDelete(d *schema.ResourceData, m interface{}) error {
	client, ok := m.(*Client)
	if !ok {
		return fmt.Errorf("invalid type assertion for client")
	}

	deleteReq := struct {
		ID string `json:"id"`
	}{
		ID: d.Id(),
	}

	resp, err := MakeRequest(client, "POST", endpointModelDelete, deleteReq)
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
