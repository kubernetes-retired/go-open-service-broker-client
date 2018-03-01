package generator

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/pmorie/go-open-service-broker-client/v2"
)

// GetCatalog will produce a valid GetCatalog response based on the generator settings.
func (g *Generator) GetCatalog() (*v2.CatalogResponse, error) {
	if len(g.Services) == 0 {
		return nil, fmt.Errorf("no services defined")
	}

	services := make([]v2.Service, len(g.Services))

	for i, _ := range services {
		populateService(i, &services[i])
	}

	return &v2.CatalogResponse{
		Services: services,
	}, nil
}

func populateService(i int, s *v2.Service) {
	if len(ServiceClassNames) < i {
		glog.Error("out of range for generated class name")
		return
	}
	s.Name = ServiceClassNames[i]
	s.ID = ServiceClassID(i)
}

//
//const okCatalogBytes = `{
//  "services": [{
//    "name": "fake-service",
//    "id": "acb56d7c-XXXX-XXXX-XXXX-feb140a59a66",
//    "description": "fake service",
//    "tags": ["tag1", "tag2"],
//    "requires": ["route_forwarding"],
//    "bindable": true,
//    "bindings_retrievable": true,
//    "metadata": {
//    	"a": "b",
//    	"c": "d"
//    },
//    "dashboard_client": {
//      "id": "398e2f8e-XXXX-XXXX-XXXX-19a71ecbcf64",
//      "secret": "277cabb0-XXXX-XXXX-XXXX-7822c0a90e5d",
//      "redirect_uri": "http://localhost:1234"
//    },
//    "plan_updateable": true,
//    "plans": [{
//      "name": "fake-plan-1",
//      "id": "d3031751-XXXX-XXXX-XXXX-a42377d3320e",
//      "description": "description1",
//      "metadata": {
//      	"b": "c",
//      	"d": "e"
//      }
//    }]
//  }]
//}`
