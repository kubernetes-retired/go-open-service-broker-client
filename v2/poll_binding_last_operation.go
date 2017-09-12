package v2

import (
	"errors"
	"fmt"
	"net/http"
)

func (c *client) PollBindingLastOperation(r *BindingLastOperationRequest) (*LastOperationResponse, error) {
	if !c.EnableAlphaFeatures {
		return nil, errors.New("alpha features must be enabled")
	}
	if err := validateBindingLastOperationRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(bindingLastOperationURLFmt, c.URL, r.InstanceID, r.BindingID)
	params := map[string]string{}

	if r.ServiceID != nil {
		params[serviceIDKey] = *r.ServiceID
	}
	if r.PlanID != nil {
		params[planIDKey] = *r.PlanID
	}
	if r.OperationKey != nil {
		op := *r.OperationKey
		opStr := string(op)
		params[operationKey] = opStr
	}

	response, err := c.prepareAndDo(http.MethodGet, fullURL, params, nil /* request body */, r.OriginatingIdentity)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		userResponse := &LastOperationResponse{}
		if err := c.unmarshalResponse(response, userResponse); err != nil {
			return nil, HTTPStatusCodeError{StatusCode: response.StatusCode, ResponseError: err}
		}

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}
}

func validateBindingLastOperationRequest(request *BindingLastOperationRequest) error {
	if request.InstanceID == "" {
		return required("instanceID")
	}

	if request.BindingID == "" {
		return required("bindingID")
	}

	return nil
}
