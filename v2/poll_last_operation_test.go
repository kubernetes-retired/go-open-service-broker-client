package v2

import (
	"fmt"
	"net/http"
	"testing"
)

func defaultLastOperationRequest() *LastOperationRequest {
	return &LastOperationRequest{
		InstanceID: testInstanceID,
		ServiceID:  strPtr(testServiceID),
		PlanID:     strPtr(testPlanID),
	}
}

const successLastOperationRequestBody = `{"service_id":"test-service-id","plan_id":"test-plan-id"}`

func successLastOperationResponse() *LastOperationResponse {
	return &LastOperationResponse{
		State:       StateSucceeded,
		Description: strPtr("test description"),
	}
}

const successLastOperationResponseBody = `{"state":"succeeded","description":"test description"}`

func inProgressLastOperationResponse() *LastOperationResponse {
	return &LastOperationResponse{
		State:       StateInProgress,
		Description: strPtr("test description"),
	}
}

const inProgressLastOperationResponseBody = `{"state":"in progress","description":"test description"}`

func failedLastOperationResponse() *LastOperationResponse {
	return &LastOperationResponse{
		State:       StateFailed,
		Description: strPtr("test description"),
	}
}

const failedLastOperationResponseBody = `{"state":"failed","description":"test description"}`

func TestPollLastOperation(t *testing.T) {
	cases := []struct {
		name                string
		enableAlpha         bool
		originatingIdentity string
		request             *LastOperationRequest
		httpChecks          httpChecks
		httpReaction        httpReaction
		expectedResponse    *LastOperationResponse
		expectedErrMessage  string
		expectedErr         error
	}{
		{
			name: "op succeeded",
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successLastOperationResponseBody,
			},
			expectedResponse: successLastOperationResponse(),
		},
		{
			name: "op in progress",
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   inProgressLastOperationResponseBody,
			},
			expectedResponse: inProgressLastOperationResponse(),
		},
		{
			name: "op failed",
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   failedLastOperationResponseBody,
			},
			expectedResponse: failedLastOperationResponse(),
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
			expectedErrMessage: "Status: 200; ErrorMessage: <nil>; Description: <nil>; ResponseError: unexpected end of JSON input",
		},
		{
			name: "500 with malformed response",
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   malformedResponse,
			},
			expectedErrMessage: "Status: 500; ErrorMessage: <nil>; Description: <nil>; ResponseError: unexpected end of JSON input",
		},
		{
			name: "500 with convential response",
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   conventionalFailureResponseBody,
			},
			expectedErr: testHttpStatusCodeError(),
		},
		{
			name: "op succeeded",
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successLastOperationResponseBody,
			},
			expectedResponse: successLastOperationResponse(),
		},
		{
			name:                "originating identity included",
			originatingIdentity: "fakeOI",
			httpChecks:          httpChecks{headers: map[string]string{XBrokerAPIOriginatingIdentity: "fakeOI"}},
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successLastOperationResponseBody,
			},
			expectedResponse: successLastOperationResponse(),
		},
		{
			name:                "originating identity excluded",
			originatingIdentity: "",
			httpChecks:          httpChecks{headers: map[string]string{XBrokerAPIOriginatingIdentity: ""}},
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successLastOperationResponseBody,
			},
			expectedResponse: successLastOperationResponse(),
		},
	}

	for _, tc := range cases {
		if tc.request == nil {
			tc.request = defaultLastOperationRequest()
		}

		if tc.httpChecks.URL == "" {
			tc.httpChecks.URL = "/v2/service_instances/test-instance-id/last_operation"
		}

		if len(tc.httpChecks.params) == 0 {
			tc.httpChecks.params = map[string]string{}
			tc.httpChecks.params[serviceIDKey] = testServiceID
			tc.httpChecks.params[planIDKey] = testPlanID
		}

		version := Version2_11()
		klient := newTestClient(t, tc.name, version, tc.enableAlpha, tc.originatingIdentity, tc.httpChecks, tc.httpReaction)

		response, err := klient.PollLastOperation(tc.request)

		doResponseChecks(t, tc.name, response, err, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func TestValidateLastOperationRequest(t *testing.T) {
	cases := []struct {
		name    string
		request *LastOperationRequest
		valid   bool
	}{
		{
			name:    "valid",
			request: defaultLastOperationRequest(),
			valid:   true,
		},
		{
			name: "missing instance ID",
			request: func() *LastOperationRequest {
				r := defaultLastOperationRequest()
				r.InstanceID = ""
				return r
			}(),
			valid: false,
		},
	}

	for _, tc := range cases {
		err := validateLastOperationRequest(tc.request)
		if err != nil {
			if tc.valid {
				t.Errorf("%v: expected valid, got error: %v", tc.name, err)
			}
		} else if !tc.valid {
			t.Errorf("%v: expected invalid, got valid", tc.name)
		}
	}
}
