/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v2

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (c *client) PollBindingLastOperation(r *BindingLastOperationRequest) (*LastOperationResponse, error) {
	if err := c.validateClientVersionIsAtLeast(Version2_14()); err != nil {
		return nil, AsyncBindingOperationsNotAllowedError{
			reason: err.Error(),
		}
	}

	if err := validateBindingLastOperationRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(bindingLastOperationURLFmt, c.URL, r.InstanceID, r.BindingID)
	params := map[string]string{}

	if r.ServiceID != nil {
		params[VarKeyServiceID] = *r.ServiceID
	}
	if r.PlanID != nil {
		params[VarKeyPlanID] = *r.PlanID
	}
	if r.OperationKey != nil {
		op := *r.OperationKey
		opStr := string(op)
		params[VarKeyOperation] = opStr
	}

	response, err := c.prepareAndDo(http.MethodGet, fullURL, params, nil /* request body */, r.OriginatingIdentity)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = drainReader(response.Body)
		response.Body.Close()
	}()

	switch response.StatusCode {
	case http.StatusOK:
		userResponse := &LastOperationResponse{}
		if err := c.unmarshalResponse(response, userResponse); err != nil {
			return nil, HTTPStatusCodeError{StatusCode: response.StatusCode, ResponseError: err}
		}

		if c.EnableAlphaFeatures {
			delayInSeconds := response.Header.Get(PollingDelayHeader)
			if delay, err := strconv.Atoi(delayInSeconds); err == nil {
				pollDelay := time.Duration(delay) * time.Second
				userResponse.PollDelay = &pollDelay
			}
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
