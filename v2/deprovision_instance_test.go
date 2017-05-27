package v2

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func defaultDeprovisionRequest() *DeprovisionRequest {
	return &DeprovisionRequest{
		InstanceID: testInstanceID,
		ServiceID:  testServiceID,
		PlanID:     testPlanID,
	}
}

func defaultAsyncDeprovisionRequest() *DeprovisionRequest {
	r := defaultDeprovisionRequest()
	r.AcceptsIncomplete = true
	return r
}

const successDeprovisionResponseBody = `{}`

func successDeprovisionResponse() *DeprovisionResponse {
	return &DeprovisionResponse{}
}

const successAsyncDeprovisionResponseBody = `{
  "operation": "test-operation-key"
}`

func successDeprovisionResponseAsync() *DeprovisionResponse {
	r := successDeprovisionResponse()
	r.Async = true
	r.OperationKey = &testOperation
	return r
}

func TestDeprovisionInstance(t *testing.T) {
	cases := []struct {
		name               string
		request            *DeprovisionRequest
		httpChecks         httpChecks
		httpReaction       httpReaction
		expectedResponse   *DeprovisionResponse
		expectedErrMessage string
		expectedErr        error
	}{
		{
			name:    "success - ok",
			request: defaultDeprovisionRequest(),
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successDeprovisionResponseBody,
			},
			expectedResponse: successDeprovisionResponse(),
		},
		{
			name:    "success - gone",
			request: defaultDeprovisionRequest(),
			httpReaction: httpReaction{
				status: http.StatusGone,
				body:   successDeprovisionResponseBody,
			},
			expectedResponse: successDeprovisionResponse(),
		},
		{
			name:    "success - async",
			request: defaultAsyncDeprovisionRequest(),
			httpReaction: httpReaction{
				status: http.StatusAccepted,
				body:   successAsyncDeprovisionResponseBody,
			},
			expectedResponse: successDeprovisionResponseAsync(),
		},
		{
			name:    "accepted with malformed response",
			request: defaultAsyncDeprovisionRequest(),
			httpReaction: httpReaction{
				status: http.StatusAccepted,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:    "http error",
			request: defaultDeprovisionRequest(),
			httpReaction: httpReaction{
				err: fmt.Errorf("http error"),
			},
			expectedErrMessage: "http error",
		},
		{
			name:    "200 with malformed response",
			request: defaultDeprovisionRequest(),
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   malformedResponse,
			},
			expectedResponse: successDeprovisionResponse(),
		},
		{
			name:    "500 with malformed response",
			request: defaultDeprovisionRequest(),
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:    "500 with conventional failure response",
			request: defaultDeprovisionRequest(),
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

		doDeprovisionInstanceTest(t, tc.name, tc.request, tc.httpChecks, tc.httpReaction, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func doDeprovisionInstanceTest(
	t *testing.T,
	name string,
	request *DeprovisionRequest,
	httpChecks httpChecks,
	httpReaction httpReaction,
	expectedResponse *DeprovisionResponse,
	expectedErrMessage string,
	expectedErr error,
) {
	klient := &client{
		Name:          "test client",
		Verbose:       true,
		URL:           "https://example.com",
		doRequestFunc: doHTTP(t, name, httpChecks, httpReaction),
	}

	response, err := klient.DeprovisionInstance(request)
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
