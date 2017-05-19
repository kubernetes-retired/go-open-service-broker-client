package v2

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
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
		prepareAndDo       prepareAndDoFunc
		responseBody       string
		expectedResponse   *ProvisionResponse
		expectedErrMessage string
		expectedErr        error
	}{
		{
			name:             "success - synchronous",
			request:          defaultProvisionRequest(),
			prepareAndDo:     returnHttpResponseFunc(http.StatusOK, successProvisionResponseBody),
			expectedResponse: successProvisionResponse(),
		},
		{
			name:             "success - asynchronous",
			request:          defaultProvisionRequest(),
			prepareAndDo:     returnHttpResponseFunc(http.StatusAccepted, successAsyncProvisionResponseBody),
			expectedResponse: successProvisionResponseAsync(),
		},
		{
			name:               "malformed response",
			request:            defaultProvisionRequest(),
			responseBody:       malformedResponse,
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:               "http error",
			request:            defaultProvisionRequest(),
			prepareAndDo:       returnErrFunc("http error"),
			expectedErrMessage: "http error",
		},
		{
			name:               "200 with malformed response",
			request:            defaultProvisionRequest(),
			prepareAndDo:       returnHttpResponseFunc(http.StatusOK, malformedResponse),
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:               "500 with malformed response",
			request:            defaultProvisionRequest(),
			prepareAndDo:       returnHttpResponseFunc(http.StatusInternalServerError, malformedResponse),
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:         "500 with conventional failure response",
			request:      defaultProvisionRequest(),
			prepareAndDo: returnHttpResponseFunc(http.StatusInternalServerError, conventionalFailureResponseBody),
			expectedErr:  testHttpStatusCodeError(),
		},
	}

	for _, tc := range cases {
		doProvisionInstanceTest(t, tc.name, tc.request, tc.responseBody, tc.prepareAndDo, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func doProvisionInstanceTest(
	t *testing.T,
	name string,
	request *ProvisionRequest,
	responseBody string,
	prepareAndDo prepareAndDoFunc,
	expectedResponse *ProvisionResponse,
	expectedErrMessage string,
	expectedErr error,
) {
	router := mux.NewRouter()
	router.HandleFunc("/v2/service_instances/test-instance-id", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		bodyBytes := []byte(responseBody)
		_, err := w.Write(bodyBytes)
		if err != nil {
			t.Errorf("%v: error writing response bytes: %v", name, err)
		}
	})

	server := httptest.NewServer(router)
	URL := server.URL
	defer server.Close()

	var klient Client
	if prepareAndDo != nil {
		klient = &client{
			Name:             "test client",
			Verbose:          true,
			URL:              URL,
			prepareAndDoFunc: prepareAndDo,
		}
	} else {
		config := DefaultClientConfiguration()
		config.URL = URL

		var err error
		klient, err = NewClient(config)
		if err != nil {
			t.Errorf("%v: error creating client: %v", name, err)
			return
		}
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
