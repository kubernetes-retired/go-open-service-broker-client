package generator

import (
	"encoding/json"
	"testing"

	"github.com/golang/glog"
)

func TestGetCatalog(t *testing.T) {
	g := Generator{
		Services: []Service{
			{

				Plans: []Plan{
					{},
					{},
				},
				Tags: 3,
			},
		},
	}

	catalog, err := g.GetCatalog()
	if err != nil {
		t.Errorf("Got error, %v", err)
	}

	catalogBytes, err := json.MarshalIndent(catalog, "", "  ")

	catalogJson := string(catalogBytes)

	if catalogJson != okCatalogBytes {
		t.Errorf("Catalog does not match. \n%s\n!=\n%s", catalogJson, okCatalogBytes)
	}
}

func TestGetPlans(t *testing.T) {
	glog.Info(planNames(1, 5))
	glog.Info(planNames(2, 5))
}

const okCatalogBytes = `{
  "services": [{
    "name": "fake-service",
    "id": "acb56d7c-XXXX-XXXX-XXXX-feb140a59a66",
    "description": "fake service",
    "tags": ["tag1", "tag2"],
    "requires": ["route_forwarding"],
    "bindable": true,
    "bindings_retrievable": true,
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
