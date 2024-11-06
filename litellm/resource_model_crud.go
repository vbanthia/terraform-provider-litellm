package litellm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
