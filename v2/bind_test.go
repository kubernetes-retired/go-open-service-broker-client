package v2

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const testBindingID = "test-binding-id"

func defaultBindRequest() *BindRequest {
	return &BindRequest{
		BindingID:  testBindingID,
		InstanceID: testInstanceID,
		ServiceID:  testServiceID,
		PlanID:     testPlanID,
	}
}

const defaultBindRequestBody = `{"service_id":"test-service-id","plan_id":"test-plan-id"}`

const successBindResponseBody = `{
  "credentials": {
    "uri": "mysql://mysqluser:pass@mysqlhost:3306/dbname",
    "username": "mysqluser",
    "password": "pass",
    "host": "mysqlhost",
    "port": 3306,
    "database": "dbname"
  }
}`

func successBindResponse() *BindResponse {
	return &BindResponse{
		Credentials: map[string]interface{}{
			"uri":      "mysql://mysqluser:pass@mysqlhost:3306/dbname",
			"username": "mysqluser",
			"password": "pass",
			"host":     "mysqlhost",
			"port":     float64(3306),
			"database": "dbname",
		},
	}
}

func optionalFieldsBindRequest() *BindRequest {
	r := defaultBindRequest()
	r.Parameters = map[string]interface{}{
		"foo": "bar",
		"blu": 2,
	}
	r.BindResource = &BindResource{
		AppGUID: strPtr("test-app-guid"),
		Route:   strPtr("test-route"),
	}
	return r
}

const optionalFieldsBindRequestBody = `{"service_id":"test-service-id","plan_id":"test-plan-id","parameters":{"blu":2,"foo":"bar"},"bind_resource":{"app_guid":"test-app-guid","route":"test-app-guid"}}`

func TestBind(t *testing.T) {
	cases := []struct {
		name               string
		request            *BindRequest
		httpChecks         httpChecks
		httpReaction       httpReaction
		expectedResponse   *BindResponse
		expectedErrMessage string
		expectedErr        error
	}{
		{
			name:    "success - created",
			request: defaultBindRequest(),
			httpReaction: httpReaction{
				status: http.StatusCreated,
				body:   successBindResponseBody,
			},
			expectedResponse: successBindResponse(),
		},
		{
			name:    "success - ok",
			request: defaultBindRequest(),
			httpReaction: httpReaction{
				status: http.StatusCreated,
				body:   successBindResponseBody,
			},
			expectedResponse: successBindResponse(),
		},
		{
			name:    "success - optional fields",
			request: optionalFieldsBindRequest(),
			httpChecks: httpChecks{
				body: optionalFieldsBindRequestBody,
			},
			httpReaction: httpReaction{
				status: http.StatusCreated,
				body:   successBindResponseBody,
			},
			expectedResponse: successBindResponse(),
		},
		{
			name:    "http error",
			request: defaultBindRequest(),
			httpReaction: httpReaction{
				err: fmt.Errorf("http error"),
			},
			expectedErrMessage: "http error",
		},
		{
			name:    "200 with malformed response",
			request: defaultBindRequest(),
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:    "500 with malformed response",
			request: defaultBindRequest(),
			httpReaction: httpReaction{
				status: http.StatusInternalServerError,
				body:   malformedResponse,
			},
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:    "500 with conventional failure response",
			request: defaultBindRequest(),
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

		if tc.httpChecks.body == "" {
			tc.httpChecks.body = defaultBindRequestBody
		}

		doBindInstanceTest(t, tc.name, tc.request, tc.httpChecks, tc.httpReaction, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func doBindInstanceTest(
	t *testing.T,
	name string,
	request *BindRequest,
	httpChecks httpChecks,
	httpReaction httpReaction,
	expectedResponse *BindResponse,
	expectedErrMessage string,
	expectedErr error,
) {
	klient := &client{
		Name:          "test client",
		Verbose:       true,
		URL:           "https://example.com",
		doRequestFunc: doHTTP(t, name, httpChecks, httpReaction),
	}

	response, err := klient.Bind(request)
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
