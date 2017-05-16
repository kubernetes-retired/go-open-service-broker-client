package v2

import (
	"bytes"
	"fmt"
	"net/http"
)

const (
	serviceIDKey = "service_id"
	planIDKey    = "plan_id"
	operationKey = "operation"
)

func (c *client) PollLastOperation(r *LastOperationRequest) (*LastOperationResponse, error) {
	// TODO: support special handling for delete responses

	if err := validateLastOperationRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(lastOperationURLFmt, c.URL, r.InstanceID)
	var queryParamBuffer bytes.Buffer
	switch {
	case r.ServiceID != nil:
		appendQueryParam(&queryParamBuffer, serviceIDKey, *r.ServiceID)
		fallthrough
	case r.PlanID != nil:
		appendQueryParam(&queryParamBuffer, planIDKey, *r.PlanID)
		fallthrough
	case r.OperationKey != nil:
		op := *r.OperationKey
		opStr := string(op)
		appendQueryParam(&queryParamBuffer, operationKey, opStr)
	}
	if queryParamBuffer.Len() > 0 {
		fullURL += "?" + queryParamBuffer.String()
	}

	response, err := c.prepareAndDoRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		userResponse := &LastOperationResponse{}
		if err := c.unmarshalResponse(response, userResponse); err != nil {
			return nil, err
		}

		return userResponse, nil
	case http.StatusGone:
		// TODO: async operations for deprovision have a special case to be
		// handled here
		fallthrough
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}
