package v2

import (
	"fmt"
	"net/http"
	"reflect"
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
		request            *UpdateInstanceRequest
		httpChecks         httpChecks
		httpReaction       httpReaction
		expectedResponse   *UpdateInstanceResponse
		expectedErrMessage string
		expectedErr        error
	}{
		{
			name:    "success - ok",
			request: defaultUpdateInstanceRequest(),
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successUpdateInstanceResponseBody,
			},
			expectedResponse: successUpdateInstanceResponse(),
		},
		{
			name:    "success - async",
			request: defaultAsyncUpdateInstanceRequest(),
			httpReaction: httpReaction{
				status: http.StatusAccepted,
				body:   successAsyncUpdateInstanceResponseBody,
			},
			expectedResponse: successUpdateInstanceResponseAsync(),
		},
		{
			name:    "accepted with malformed response",
			request: defaultAsyncUpdateInstanceRequest(),
			httpReaction: httpReaction{
				status: http.StatusAccepted,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:    "http error",
			request: defaultUpdateInstanceRequest(),
			httpReaction: httpReaction{
				err: fmt.Errorf("http error"),
			},
			expectedErrMessage: "http error",
		},
		{
			name:    "200 with malformed response",
			request: defaultUpdateInstanceRequest(),
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   malformedResponse,
			},
			expectedResponse: successUpdateInstanceResponse(),
		},
		{
			name:    "500 with malformed response",
			request: defaultUpdateInstanceRequest(),
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:    "500 with conventional failure response",
			request: defaultUpdateInstanceRequest(),
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   conventionalFailureResponseBody,
			},
			expectedErr: testHttpStatusCodeError(),
		},
	}

	for _, tc := range cases {
		if tc.httpChecks.URL == "" {
			tc.httpChecks.URL = "/v2/service_instances/test-instance-id"
		}

		if tc.httpChecks.body == "" {
			tc.httpChecks.body = "{}"
		}

		doUpdateInstanceInstanceTest(t, tc.name, tc.request, tc.httpChecks, tc.httpReaction, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func doUpdateInstanceInstanceTest(
	t *testing.T,
	name string,
	request *UpdateInstanceRequest,
	httpChecks httpChecks,
	httpReaction httpReaction,
	expectedResponse *UpdateInstanceResponse,
	expectedErrMessage string,
	expectedErr error,
) {
	klient := &client{
		Name:          "test client",
		Verbose:       true,
		URL:           "https://example.com",
		doRequestFunc: doHTTP(t, name, httpChecks, httpReaction),
	}

	response, err := klient.UpdateInstance(request)
	if err != nil && expectedErrMessage == "" && expectedErr == nil {
		t.Errorf("%v: error getting catalog: %v", name, err)
		return
	} else if err != nil && expectedErrMessage != "" && expectedErrMessage != err.Error() {
		t.Errorf("%v: unexpected error message: expected %v, got %v", name, expectedErrMessage, err)
		return
	} else if err != nil && expectedErr != nil && !reflect.DeepEqual(expectedErr, err) {
		t.Errorf("%v: unexpected error: expected %+v, got %v", name, expectedErr, err)
		return
	}

	if e, a := expectedResponse, response; !reflect.DeepEqual(e, a) {
		t.Errorf("%v: unexpected diff in catalog response; expected %+v, got %+v", name, e, a)
		return
	}

}
