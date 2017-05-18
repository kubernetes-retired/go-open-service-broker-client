package v2

import (
	"bytes"
	"fmt"
	"net/http"
)

func (c *client) Unbind(r *UnbindRequest) (*UnbindResponse, error) {
	if err := validateUnbindRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(bindingURLFmt, c.URL, r.InstanceID, r.BindingID)
	var queryParamBuffer bytes.Buffer
	appendQueryParam(&queryParamBuffer, serviceIDKey, r.ServiceID)
	appendQueryParam(&queryParamBuffer, planIDKey, r.PlanID)
	fullURL += "?" + queryParamBuffer.String()

	response, err := c.prepareAndDoFunc(http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusGone:
		// TODO: should we establish that the response body ('{}') is correct?
		return &UnbindResponse{}, nil
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}
