package v2

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

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

func truePtr() *bool {
	b := true
	return &b
}

func falsePtr() *bool {
	b := false
	return &b
}

func TestGetCatalogSuccess(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/v2/catalog", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		bodyBytes := []byte(okCatalogBytes)
		_, err := w.Write(bodyBytes)
		if err != nil {
			t.Fatalf("error writing response bytes: %v", err)
		}
	})

	server := httptest.NewServer(router)
	URL := server.URL
	defer server.Close()

	config := DefaultClientConfiguration()
	config.URL = URL

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}

	catalog, err := client.GetCatalog()
	if err != nil {
		t.Fatalf("error getting catalog: %v", err)
	}

	if e, a := okCatalogResponse(), catalog; !reflect.DeepEqual(e, a) {
		t.Fatalf("unexpected diff in catalog response; expected %+v, got %+v", e, a)
	}
}
