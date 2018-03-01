package generator

import (
	"fmt"

	"math/rand"

	"sort"

	"github.com/pmorie/go-open-service-broker-client/v2"
)

// GetCatalog will produce a valid GetCatalog response based on the generator settings.
func (g *Generator) GetCatalog() (*v2.CatalogResponse, error) {
	if len(g.Services) == 0 {
		return nil, fmt.Errorf("no services defined")
	}

	services := make([]v2.Service, len(g.Services))

	for s, _ := range services {
		services[s].Plans = make([]v2.Plan, len(g.Services[s].Plans))
		service := &services[s]
		service.Name = ClassNames[s]
		service.ID = IDFrom(ClassNames[s])
		service.Tags = tags(s, g.Services[s].Tags)
		planNames := planNames(s, len(service.Plans))
		for p, _ := range service.Plans {
			service.Plans[p].Name = planNames[p]
			service.Plans[p].ID = IDFrom(planNames[p])
		}
	}

	return &v2.CatalogResponse{
		Services: services,
	}, nil
}

func getSubsetWithoutDuplicates(count int, seed int64, list []string) []string {
	rand.Seed(seed)

	plans := map[string]int32{}

	// Get count of plan names without duplicates
	for len(plans) < count {
		x := rand.Int31n(int32(len(list)))
		plans[list[x]] = x
	}

	keys := []string(nil)
	for k := range plans {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func planNames(service, count int) []string {
	return getSubsetWithoutDuplicates(count, int64(service), PlanNames)
}

func tags(service, count int) []string {
	return getSubsetWithoutDuplicates(count, int64(service*1000), TagNames)
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
