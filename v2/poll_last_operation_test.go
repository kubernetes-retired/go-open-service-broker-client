package v2

import (
	"fmt"
	"net/http"
	"reflect"
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
		name               string
		request            *LastOperationRequest
		httpChecks         httpChecks
		httpReaction       httpReaction
		expectedResponse   *LastOperationResponse
		expectedErrMessage string
		expectedErr        error
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
			name: "500 with malformed response",
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   conventionalFailureResponseBody,
			},
			expectedErr: testHttpStatusCodeError(),
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

		doPollLastOperationTest(t, tc.name, tc.request, tc.httpChecks, tc.httpReaction, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func doPollLastOperationTest(
	t *testing.T,
	name string,
	request *LastOperationRequest,
	httpChecks httpChecks,
	httpReaction httpReaction,
	expectedResponse *LastOperationResponse,
	expectedErrMessage string,
	expectedErr error,
) {
	klient := &client{
		Name:          "test client",
		Verbose:       true,
		URL:           "https://example.com",
		doRequestFunc: doHTTP(t, name, httpChecks, httpReaction),
	}

	response, err := klient.PollLastOperation(request)
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
