package v2

import (
	"bytes"
	// "errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

// func init() {
// 	flag.Set("alsologtostderr", "true")
// 	flag.Set("v", "5")
// }

const okCatalogBytes = `{
  "services": [{
    "name": "fake-service",
    "id": "acb56d7c-XXXX-XXXX-XXXX-feb140a59a66",
    "description": "fake service",
    "tags": ["tag1", "tag2"],
    "requires": ["route_forwarding"],
    "bindable": true,
    "metadata": {
    	"a": "b",
    	"c": "d"
    },
    "dashboard_client": {
      "id": "398e2f8e-XXXX-XXXX-XXXX-19a71ecbcf64",
      "secret": "277cabb0-XXXX-XXXX-XXXX-7822c0a90e5d",
      "redirect_uri": "http://localhost:1234"
    },
    "plan_updateable": true,
    "plans": [{
      "name": "fake-plan-1",
      "id": "d3031751-XXXX-XXXX-XXXX-a42377d3320e",
      "description": "description1",
      "metadata": {
      	"b": "c",
      	"d": "e"
      }
    }]
  }]
}`

func okCatalogResponse() *CatalogResponse {
	return &CatalogResponse{
		Services: []Service{
			{
				ID:          "acb56d7c-XXXX-XXXX-XXXX-feb140a59a66",
				Name:        "fake-service",
				Description: "fake service",
				Tags: []string{
					"tag1",
					"tag2",
				},
				Requires: []string{
					"route_forwarding",
				},
				Bindable:      true,
				PlanUpdatable: truePtr(),
				Plans: []Plan{
					{
						ID:          "d3031751-XXXX-XXXX-XXXX-a42377d3320e",
						Name:        "fake-plan-1",
						Description: "description1",
						Metadata: map[string]interface{}{
							"b": "c",
							"d": "e",
						},
					},
				},
				DashboardClient: &DashboardClient{
					ID:          "398e2f8e-XXXX-XXXX-XXXX-19a71ecbcf64",
					Secret:      "277cabb0-XXXX-XXXX-XXXX-7822c0a90e5d",
					RedirectURI: "http://localhost:1234",
				},
				Metadata: map[string]interface{}{
					"a": "b",
					"c": "d",
				},
			},
		},
	}
}

const okCatalog2Bytes = `{
  "services": [{
    "name": "fake-service-2",
    "id": "fake-service-2-id",
    "description": "service-description-2",
    "bindable": false,
    "plans": [{
      "name": "fake-plan-2",
      "id": "fake-plan-2-id",
      "description": "description-2",
      "bindable": true
    }]
  }]
}`

func okCatalog2Response() *CatalogResponse {
	return &CatalogResponse{
		Services: []Service{
			{
				ID:          "fake-service-2-id",
				Name:        "fake-service-2",
				Description: "service-description-2",
				Bindable:    false,
				Plans: []Plan{
					{
						ID:          "fake-plan-2-id",
						Name:        "fake-plan-2",
						Description: "description-2",
						Bindable:    truePtr(),
					},
				},
			},
		},
	}
}

const malformedResponse = `{`

const conventionalFailureResponseBody = `{
	"error": "TestError",
	"description": "test error description"
}`

func testHttpStatusCodeError() error {
	errorMessage := "TestError"
	description := "test error description"
	return HTTPStatusCodeError{http.StatusInternalServerError, &errorMessage, &description}
}

func truePtr() *bool {
	b := true
	return &b
}

func falsePtr() *bool {
	b := false
	return &b
}

func closer(s string) io.ReadCloser {
	return nopCloser{bytes.NewBufferString(s)}
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type httpChecks struct {
	URL    string
	body   string
	params map[string]string
}

type httpReaction struct {
	status int
	body   string
	err    error
}

var walkingGhostErr = fmt.Errorf("test has already failed")

func doHTTP(t *testing.T, name string, checks httpChecks, reaction httpReaction) func(*http.Request) (*http.Response, error) {
	return func(request *http.Request) (*http.Response, error) {
		if len(checks.URL) > 0 && checks.URL != request.URL.Path {
			t.Errorf("%v: unexpected URL; expected %v, got %v", name, checks.URL, request.URL.Path)
			return nil, walkingGhostErr
		}

		if len(checks.params) > 0 {
			for k, v := range checks.params {
				actualValue := request.URL.Query().Get(k)
				if e, a := v, actualValue; e != a {
					t.Errorf("%v: unexpected parameter value for key %v; expected %v, got %v", name, k, e, a)
					return nil, walkingGhostErr
				}
			}
		}

		var bodyBytes []byte
		if request.Body != nil {
			var err error
			bodyBytes, err = ioutil.ReadAll(request.Body)
			if err != nil {
				t.Errorf("%v: error reading request body bytes: %v", name, err)
				return nil, walkingGhostErr
			}
		}

		if e, a := checks.body, string(bodyBytes); e != a {
			t.Errorf("%v: unexpected request body: expected %v, got %v", name, e, a)
			return nil, walkingGhostErr
		}

		return &http.Response{
			StatusCode: reaction.status,
			Body:       closer(reaction.body),
		}, reaction.err
	}
}

func TestGetCatalog(t *testing.T) {
	cases := []struct {
		name               string
		httpReaction       httpReaction
		expectedResponse   *CatalogResponse
		expectedErrMessage string
		expectedErr        error
	}{
		{
			name: "success 1",
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   okCatalogBytes,
			},
			expectedResponse: okCatalogResponse(),
		},
		{
			name: "success 2",
			httpReaction: httpReaction{
				status: http.StatusOK,
				body:   okCatalog2Bytes,
			},
			expectedResponse: okCatalog2Response(),
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
		doGetCatalogTest(t, tc.name, tc.httpReaction, tc.expectedResponse, tc.expectedErrMessage, tc.expectedErr)
	}
}

func doGetCatalogTest(
	t *testing.T,
	name string,
	httpReaction httpReaction,
	expectedResponse *CatalogResponse,
	expectedErrMessage string,
	expectedErr error,
) {

	klient := &client{
		Name:          "test client",
		Verbose:       true,
		URL:           "https://example.com",
		doRequestFunc: doHTTP(t, name, httpChecks{}, httpReaction),
	}

	response, err := klient.GetCatalog()
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
