package litellm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

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
