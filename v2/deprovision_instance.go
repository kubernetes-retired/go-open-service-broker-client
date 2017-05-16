package v2

import (
	"fmt"
	"net/http"
)

// internal message body types

type deprovisionInstanceRequestBody struct {
	serviceID *string `json:"service_id"`
	planID    *string `json:"plan_id,omitempty"`
}

func (c *client) DeprovisionInstance(r *DeprovisionRequest) (*DeprovisionResponse, error) {
	if err := validateDeprovisionRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(serviceInstanceURLFmt, c.URL, r.InstanceID)
	if r.AcceptsIncomplete {
		fullURL += "?accepts_incomplete=true"
	}

	requestServiceID := string(r.ServiceID)
	requestPlanID := string(r.PlanID)

	requestBody := &deprovisionInstanceRequestBody{
		serviceID: &requestServiceID,
		planID:    &requestPlanID,
	}

	response, err := c.prepareAndDoRequest(http.MethodDelete, fullURL, requestBody)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusGone:
		return &DeprovisionResponse{}, nil
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

		userResponse := &DeprovisionResponse{
			Async:        true,
			OperationKey: opPtr,
		}

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}
