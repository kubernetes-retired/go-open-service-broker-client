package v2

import (
	"fmt"
	"net/http"
	"testing"
)

func defaultUpdateInstanceRequest() *UpdateInstanceRequest {
	return &UpdateInstanceRequest{
		InstanceID: testInstanceID,
		ServiceID:  testServiceID,
		PlanID:     strPtr(testPlanID),
	}
}

func defaultAsyncUpdateInstanceRequest() *UpdateInstanceRequest {
	r := defaultUpdateInstanceRequest()
	r.AcceptsIncomplete = true
	return r
}

const successUpdateInstanceResponseBody = `{}`

func successUpdateInstanceResponse() *UpdateInstanceResponse {
	return &UpdateInstanceResponse{}
}

const successAsyncUpdateInstanceResponseBody = `{
  "operation": "test-operation-key"
}`

func successUpdateInstanceResponseAsync() *UpdateInstanceResponse {
	r := successUpdateInstanceResponse()
	r.Async = true
	r.OperationKey = &testOperation
	return r
}

func TestUpdateInstanceInstance(t *testing.T) {
	cases := []struct {
		name               string
		enableAlpha        bool
		request            *UpdateInstanceRequest
		httpChecks         httpChecks
		httpReaction       httpReaction
		expectedResponse   *UpdateInstanceResponse
		expectedErrMessage string
		expectedErr        error
	}{
		{
			name: "invalid request",
			request: func() *UpdateInstanceRequest {
				r := defaultUpdateInstanceRequest()
				r.InstanceID = ""
				return r
			}(),
			expectedErrMessage: "instanceID is required",
		},
		{
			name: "success - ok",
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successUpdateInstanceResponseBody,
			},
			expectedResponse: successUpdateInstanceResponse(),
		},
		{
			name:    "success - async",
			request: defaultAsyncUpdateInstanceRequest(),
			httpChecks: httpChecks{
				params: map[string]string{
					asyncQueryParamKey: "true",
				},
			},
			httpReaction: httpReaction{
				status: http.StatusAccepted,
				body:   successAsyncUpdateInstanceResponseBody,
			},
			expectedResponse: successUpdateInstanceResponseAsync(),
		},
		{
			name:    "accepted with malformed response",
			request: defaultAsyncUpdateInstanceRequest(),
			httpChecks: httpChecks{
				params: map[string]string{
					asyncQueryParamKey: "true",
				},
			},
			httpReaction: httpReaction{
				status: http.StatusAccepted,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name: "http error",
			httpReaction: httpReaction{
				err: fmt.Errorf("http error"),
			},
			expectedErrMessage: "http error",
		},
		{
			name: "200 with malformed response",
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   malformedResponse,
			},
			expectedResponse: successUpdateInstanceResponse(),
		},
		{
			name: "500 with malformed response",
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name: "500 with conventional failure response",
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   conventionalFailureResponseBody,
			},
			expectedErr: testHttpStatusCodeError(),
		},
	}

	for _, tc := range cases {
		if tc.request == nil {
			tc.request = defaultUpdateInstanceRequest()
		}

		if tc.httpChecks.URL == "" {
			tc.httpChecks.URL = "/v2/service_instances/test-instance-id"
		}

		if tc.httpChecks.body == "" {
			tc.httpChecks.body = "{}"
		}

		klient := newTestClient(t, tc.name, tc.enableAlpha, tc.httpChecks, tc.httpReaction)

		response, err := klient.UpdateInstance(tc.request)

		doResponseChecks(t, tc.name, response, err, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func TestValidateUpdateInstanceRequest(t *testing.T) {
	cases := []struct {
		name    string
		request *UpdateInstanceRequest
		valid   bool
	}{
		{
			name:    "valid",
			request: defaultUpdateInstanceRequest(),
			valid:   true,
		},
		{
			name: "missing instance ID",
			request: func() *UpdateInstanceRequest {
				r := defaultUpdateInstanceRequest()
				r.InstanceID = ""
				return r
			}(),
			valid: false,
		},
	}

	for _, tc := range cases {
		err := validateUpdateInstanceRequest(tc.request)
		if err != nil {
			if tc.valid {
				t.Errorf("%v: expected valid, got error: %v", tc.name, err)
			}
		} else if !tc.valid {
			t.Errorf("%v: expected invalid, got valid", tc.name)
		}
	}
}
