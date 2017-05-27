package v2

import (
	"fmt"
	"net/http"
	"reflect"
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
			name:    "success - ok",
			request: defaultUnbindRequest(),
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successUnbindResponseBody,
			},
			expectedResponse: successUnbindResponse(),
		},
		{
			name:    "http error",
			request: defaultUnbindRequest(),
			httpReaction: httpReaction{
				err: fmt.Errorf("http error"),
			},
			expectedErrMessage: "http error",
		},
		{
			name:    "200 with malformed response",
			request: defaultUnbindRequest(),
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:    "500 with malformed response",
			request: defaultUnbindRequest(),
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:    "500 with conventional failure response",
			request: defaultUnbindRequest(),
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   conventionalFailureResponseBody,
			},
			expectedErr: testHttpStatusCodeError(),
		},
	}

	for _, tc := range cases {
		if tc.httpChecks.URL == "" {
			tc.httpChecks.URL = "/v2/service_instances/test-instance-id/service_bindings/test-binding-id"
		}

		if len(tc.httpChecks.params) != 0 {
			tc.httpChecks.params["service_id"] = testServiceID
			tc.httpChecks.params["plan_id"] = testPlanID
		}

		doUnbindInstanceTest(t, tc.name, tc.request, tc.httpChecks, tc.httpReaction, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func doUnbindInstanceTest(
	t *testing.T,
	name string,
	request *UnbindRequest,
	httpChecks httpChecks,
	httpReaction httpReaction,
	expectedResponse *UnbindResponse,
	expectedErrMessage string,
	expectedErr error,
) {
	klient := &client{
		Name:          "test client",
		Verbose:       true,
		URL:           "https://example.com",
		doRequestFunc: doHTTP(t, name, httpChecks, httpReaction),
	}

	response, err := klient.Unbind(request)
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
		t.Errorf("%v: unexpected diff in bind response; expected %+v, got %+v", name, e, a)
		return
	}

}
