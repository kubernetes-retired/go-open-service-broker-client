package v2

import (
	"fmt"
	"net/http"
)

func (c *client) Unbind(r *UnbindRequest) (*UnbindResponse, error) {
	if err := validateUnbindRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(bindingURLFmt, c.URL, r.InstanceID, r.BindingID)
	params := map[string]string{}
	params[serviceIDKey] = r.ServiceID
	params[planIDKey] = r.PlanID

	response, err := c.prepareAndDo(http.MethodDelete, fullURL, params, nil)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusGone:
		userResponse := &UnbindResponse{}
		if err := c.unmarshalResponse(response, userResponse); err != nil {
			return nil, err
		}

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}
