package v2

import (
	"fmt"
	"net/http"
	"testing"
)

func defaultUnbindRequest() *UnbindRequest {
	return &UnbindRequest{
		BindingID:  testBindingID,
		InstanceID: testInstanceID,
		ServiceID:  testServiceID,
		PlanID:     testPlanID,
	}
}

const successUnbindResponseBody = `{}`

func successUnbindResponse() *UnbindResponse {
	return &UnbindResponse{}
}

func TestUnbind(t *testing.T) {
	cases := []struct {
		name                string
		version             APIVersion
		enableAlpha         bool
		originatingIdentity *OriginatingIdentity
		request             *UnbindRequest
		httpChecks          httpChecks
		httpReaction        httpReaction
		expectedResponse    *UnbindResponse
		expectedErrMessage  string
		expectedErr         error
	}{
		{
			name: "invalid request",
			request: func() *UnbindRequest {
				r := defaultUnbindRequest()
				r.InstanceID = ""
				return r
			}(),
			expectedErrMessage: "instanceID is required",
		},
		{
			name: "success - ok",
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successUnbindResponseBody,
			},
			expectedResponse: successUnbindResponse(),
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
			name: "500 with conventional failure response",
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   conventionalFailureResponseBody,
			},
			expectedErr: testHTTPStatusCodeError(),
		},
		{
			name:                "originating identity included",
			version:             Version2_13(),
			originatingIdentity: testOriginatingIdentity,
			httpChecks:          httpChecks{headers: map[string]string{OriginatingIdentityHeader: testOriginatingIdentityHeaderValue}},
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successUnbindResponseBody,
			},
			expectedResponse: successUnbindResponse(),
		},
		{
			name:                "originating identity excluded",
			version:             Version2_13(),
			originatingIdentity: nil,
			httpChecks:          httpChecks{headers: map[string]string{OriginatingIdentityHeader: ""}},
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successUnbindResponseBody,
			},
			expectedResponse: successUnbindResponse(),
		},
		{
			name:                "originating identity not sent unless API version >= 2.13",
			version:             Version2_12(),
			originatingIdentity: testOriginatingIdentity,
			httpChecks:          httpChecks{headers: map[string]string{OriginatingIdentityHeader: ""}},
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successUnbindResponseBody,
			},
			expectedResponse: successUnbindResponse(),
		},
	}

	for _, tc := range cases {
		if tc.request == nil {
			tc.request = defaultUnbindRequest()
		}

		tc.request.OriginatingIdentity = tc.originatingIdentity

		if tc.httpChecks.URL == "" {
			tc.httpChecks.URL = "/v2/service_instances/test-instance-id/service_bindings/test-binding-id"
		}

		if len(tc.httpChecks.params) == 0 {
			tc.httpChecks.params = map[string]string{}
			tc.httpChecks.params[serviceIDKey] = testServiceID
			tc.httpChecks.params[planIDKey] = testPlanID
		}

		if tc.version.label == "" {
			tc.version = Version2_11()
		}

		klient := newTestClient(t, tc.name, tc.version, tc.enableAlpha, tc.httpChecks, tc.httpReaction)

		response, err := klient.Unbind(tc.request)

		doResponseChecks(t, tc.name, response, err, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func TestValidateUnbindRequest(t *testing.T) {
	cases := []struct {
		name    string
		request *UnbindRequest
		valid   bool
	}{
		{
			name:    "valid",
			request: defaultUnbindRequest(),
			valid:   true,
		},
		{
			name: "missing binding ID",
			request: func() *UnbindRequest {
				r := defaultUnbindRequest()
				r.BindingID = ""
				return r
			}(),
			valid: false,
		},
		{
			name: "missing instance ID",
			request: func() *UnbindRequest {
				r := defaultUnbindRequest()
				r.InstanceID = ""
				return r
			}(),
			valid: false,
		},
		{
			name: "missing service ID",
			request: func() *UnbindRequest {
				r := defaultUnbindRequest()
				r.ServiceID = ""
				return r
			}(),
			valid: false,
		},
		{
			name: "missing plan ID",
			request: func() *UnbindRequest {
				r := defaultUnbindRequest()
				r.PlanID = ""
				return r
			}(),
			valid: false,
		},
	}

	for _, tc := range cases {
		err := validateUnbindRequest(tc.request)
		if err != nil {
			if tc.valid {
				t.Errorf("%v: expected valid, got error: %v", tc.name, err)
			}
		} else if !tc.valid {
			t.Errorf("%v: expected invalid, got valid", tc.name)
		}
	}
}
