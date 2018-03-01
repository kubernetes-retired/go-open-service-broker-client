package generator

import (
	"encoding/json"
	"testing"
)

func TestCreateGenerator(t *testing.T) {
	g := CreateGenerator(3, Parameters{
		Services: ServiceRanges{
			Plans:               5,
			Tags:                6,
			Metadata:            4,
			Requires:            2,
			Bindable:            10,
			BindingsRetrievable: 1,
		},
		Plans: PlanRanges{
			Metadata: 4,
			Bindable: 10,
			Free:     4,
		},
	})
	AssignPoolGoT(g)

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
