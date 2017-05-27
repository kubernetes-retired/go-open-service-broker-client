package v2

import (
	"fmt"
	"net/http"
	"testing"
)

// const testBindingID = "test-binding-id"

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
		name               string
		request            *UnbindRequest
		httpChecks         httpChecks
		httpReaction       httpReaction
		expectedResponse   *UnbindResponse
		expectedErrMessage string
		expectedErr        error
	}{
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
			expectedErrMessage: "unexpected end of JSON input",
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
			tc.request = defaultUnbindRequest()
		}

		if tc.httpChecks.URL == "" {
			tc.httpChecks.URL = "/v2/service_instances/test-instance-id/service_bindings/test-binding-id"
		}

		if len(tc.httpChecks.params) == 0 {
			tc.httpChecks.params = map[string]string{}
			tc.httpChecks.params[serviceIDKey] = testServiceID
			tc.httpChecks.params[planIDKey] = testPlanID
		}

		klient := newTestClient(t, tc.name, tc.httpChecks, tc.httpReaction)

		response, err := klient.Unbind(tc.request)

		doResponseChecks(t, tc.name, response, err, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}
