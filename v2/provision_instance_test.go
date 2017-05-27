package v2

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const (
	testInstanceID       = "test-instance-id"
	testServiceID        = "test-service-id"
	testPlanID           = "test-plan-id"
	testOrganizationGUID = "test-organization-guid"
	testSpaceGUID        = "test-space-guid"
)

func defaultProvisionRequest() *ProvisionRequest {
	return &ProvisionRequest{
		InstanceID:       testInstanceID,
		ServiceID:        testServiceID,
		PlanID:           testPlanID,
		OrganizationGUID: testOrganizationGUID,
		SpaceGUID:        testSpaceGUID,
	}
}

func defaultAsyncProvisionRequest() *ProvisionRequest {
	r := defaultProvisionRequest()
	r.AcceptsIncomplete = true
	return r
}

const successProvisionRequestBody = `{"service_id":"test-service-id","plan_id":"test-plan-id","organization_guid":"test-organization-guid","space_guid":"test-space-guid"}`

const successProvisionResponseBody = `{
  "dashboard_url": "https://example.com/dashboard"
}`

var testDashboardURL = "https://example.com/dashboard"
var testOperation OperationKey = "test-operation-key"

func successProvisionResponse() *ProvisionResponse {
	return &ProvisionResponse{
		DashboardURL: &testDashboardURL,
	}
}

const successAsyncProvisionResponseBody = `{
  "dashboard_url": "https://example.com/dashboard",
  "operation": "test-operation-key"
}`

func successProvisionResponseAsync() *ProvisionResponse {
	r := successProvisionResponse()
	r.Async = true
	r.OperationKey = &testOperation
	return r
}

func TestProvisionInstance(t *testing.T) {
	cases := []struct {
		name               string
		request            *ProvisionRequest
		httpChecks         httpChecks
		httpReaction       httpReaction
		expectedResponse   *ProvisionResponse
		expectedErrMessage string
		expectedErr        error
	}{
		{
			name:    "success - created",
			request: defaultProvisionRequest(),
			httpReaction: httpReaction{
				status: http.StatusCreated,
				body:   successProvisionResponseBody,
			},
			expectedResponse: successProvisionResponse(),
		},
		{
			name:    "success - ok",
			request: defaultProvisionRequest(),
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   successProvisionResponseBody,
			},
			expectedResponse: successProvisionResponse(),
		},
		{
			name:    "success - asynchronous",
			request: defaultAsyncProvisionRequest(),
			httpChecks: httpChecks{
				params: map[string]string{
					"accepts_incomplete": "true",
				},
			},
			httpReaction: httpReaction{
				status: http.StatusAccepted,
				body:   successAsyncProvisionResponseBody,
			},
			expectedResponse: successProvisionResponseAsync(),
		},
		{
			name:    "http error",
			request: defaultProvisionRequest(),
			httpReaction: httpReaction{
				err: fmt.Errorf("http error"),
			},
			expectedErrMessage: "http error",
		},
		{
			name:    "200 with malformed response",
			request: defaultProvisionRequest(),
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:    "500 with malformed response",
			request: defaultProvisionRequest(),
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:    "500 with conventional failure response",
			request: defaultProvisionRequest(),
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
			tc.httpChecks.body = successProvisionRequestBody
		}

		doProvisionInstanceTest(t, tc.name, tc.request, tc.httpChecks, tc.httpReaction, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func doProvisionInstanceTest(
	t *testing.T,
	name string,
	request *ProvisionRequest,
	httpChecks httpChecks,
	httpReaction httpReaction,
	expectedResponse *ProvisionResponse,
	expectedErrMessage string,
	expectedErr error,
) {
	klient := &client{
		Name:          "test client",
		Verbose:       true,
		URL:           "https://example.com",
		doRequestFunc: doHTTP(t, name, httpChecks, httpReaction),
	}

	response, err := klient.ProvisionInstance(request)
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
