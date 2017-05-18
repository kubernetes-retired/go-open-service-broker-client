package v2

import (
	"fmt"
	"net/http"
)

// internal message body types

type updateInstanceRequestBody struct {
	serviceID  string                 `json:"service_id"`
	planID     *string                `json:"plan_id,omitempty"`
	parameters map[string]interface{} `json:"parameters,omitempty"`

	// Note: this client does not currently support the 'previous_values'
	// field of this request body.
}

func (c *client) UpdateInstance(r *UpdateInstanceRequest) (*UpdateInstanceResponse, error) {
	if err := validateUpdateInstanceRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(serviceInstanceURLFmt, c.URL, r.InstanceID)
	if r.AcceptsIncomplete {
		fullURL += "?accepts_incomplete=true"
	}

	requestBody := &updateInstanceRequestBody{
		serviceID:  r.ServiceID,
		planID:     r.PlanID,
		parameters: r.Parameters,
	}

	response, err := c.prepareAndDoFunc(http.MethodPatch, fullURL, requestBody)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return &UpdateInstanceResponse{}, nil
	case http.StatusAccepted:
		responseBodyObj := &asyncSuccessResponseBody{}
		if err := c.unmarshalResponse(response, responseBodyObj); err != nil {
			return nil, err
		}

		var opPtr *OperationKey
		if responseBodyObj.operation != nil {
			opStr := *responseBodyObj.operation
			op := OperationKey(opStr)
			opPtr = &op
		}

		userResponse := &UpdateInstanceResponse{
			Async:        true,
			OperationKey: opPtr,
		}

		// TODO: fix op key handling

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}
