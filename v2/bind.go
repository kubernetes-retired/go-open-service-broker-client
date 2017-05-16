package v2

import (
	"fmt"
	"net/http"
)

// internal message body types

type bindRequestBody struct {
	serviceID    string                 `json:"service_id"`
	planID       string                 `json:"plan_id"`
	parameters   map[string]interface{} `json:"parameters,omitempty"`
	bindResource map[string]interface{} `json:"bind_resource,omitempty"`
}

const (
	bindResourceAppGUIDKey = "app_guid"
	bindResourceRouteKey   = "route"
)

func (c *client) Bind(r *BindRequest) (*BindResponse, error) {
	if err := validateBindRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(bindingURLFmt, c.URL, r.InstanceID, r.BindingID)

	requestBody := &bindRequestBody{
		serviceID:  r.ServiceID,
		planID:     r.PlanID,
		parameters: r.Parameters,
	}

	if r.BindResource != nil {
		requestBody.bindResource = map[string]interface{}{}
		if r.BindResource.AppGUID != nil {
			requestBody.bindResource[bindResourceAppGUIDKey] = *r.BindResource.AppGUID
		}
		if r.BindResource.Route != nil {
			requestBody.bindResource[bindResourceRouteKey] = *r.BindResource.AppGUID
		}
	}

	response, err := c.prepareAndDoRequest(http.MethodPut, fullURL, requestBody)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		userResponse := &BindResponse{}
		if err := c.unmarshalResponse(response, userResponse); err != nil {
			return nil, err
		}

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}
