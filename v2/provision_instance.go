package v2

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

// internal message body types

type provisionRequestBody struct {
	serviceID        string                 `json:"service_id"`
	planID           string                 `json:"plan_id"`
	organizationGUID string                 `json:"organization_guid"`
	spaceGUID        string                 `json:"space_guid"`
	parameters       map[string]interface{} `json:"parameters,omitempty"`
}

type provisionSuccessResponseBody struct {
	dashboardURL *string `json:"dashboard_url"`
	operation    *string `json:"operation"`
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
		serviceID:        r.ServiceID,
		planID:           r.PlanID,
		organizationGUID: r.OrganizationGUID,
		spaceGUID:        r.SpaceGUID,
		parameters:       r.Parameters,
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
		if responseBodyObj.operation != nil {
			opStr := *responseBodyObj.operation
			op := OperationKey(opStr)
			opPtr = &op
		}

		userResponse := &ProvisionResponse{
			DashboardURL: responseBodyObj.dashboardURL,
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
