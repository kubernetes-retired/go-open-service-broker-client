package v2

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

// internal message body types

type provisionRequestBody struct {
	ServiceID        string                 `json:"service_id"`
	PlanID           string                 `json:"plan_id"`
	OrganizationGUID string                 `json:"organization_guid"`
	SpaceGUID        string                 `json:"space_guid"`
	Parameters       map[string]interface{} `json:"parameters,omitempty"`
}

type provisionSuccessResponseBody struct {
	DashboardURL *string `json:"dashboard_url"`
	Operation    *string `json:"operation"`
}

func (c *client) ProvisionInstance(r *ProvisionRequest) (*ProvisionResponse, error) {
	if err := validateProvisionRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(serviceInstanceURLFmt, c.URL, r.InstanceID)
	if r.AcceptsIncomplete {
		fullURL += "?accepts_incomplete=true"
	}

	requestBody := &provisionRequestBody{
		ServiceID:        r.ServiceID,
		PlanID:           r.PlanID,
		OrganizationGUID: r.OrganizationGUID,
		SpaceGUID:        r.SpaceGUID,
		Parameters:       r.Parameters,
	}

	response, err := c.prepareAndDoFunc(http.MethodPut, fullURL, requestBody)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusCreated, http.StatusOK, http.StatusAccepted:
		responseBodyObj := &provisionSuccessResponseBody{}
		if err := c.unmarshalResponse(response, responseBodyObj); err != nil {
			return nil, err
		}

		var opPtr *OperationKey
		if responseBodyObj.Operation != nil {
			opStr := *responseBodyObj.Operation
			op := OperationKey(opStr)
			opPtr = &op
		}

		userResponse := &ProvisionResponse{
			DashboardURL: responseBodyObj.DashboardURL,
			OperationKey: opPtr,
		}
		if response.StatusCode == http.StatusAccepted {
			if c.Verbose {
				glog.Infof("broker %q: received asynchronous response", c.Name)
			}
			userResponse.Async = true
		}

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}
}
