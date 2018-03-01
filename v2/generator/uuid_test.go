package generator

import (
	"testing"
)

func testForDuplicates(list []string, t *testing.T) {
	uuids := map[string]string{}
	for _, name := range list {
		uuid := IDFrom(name)
		if uuids[uuid] != "" {
			t.Fatalf("Found a collision with %s and %s", uuids[uuid], name)
		} else {
			uuids[uuid] = name
		}
	}
}

func TestClassUUID(t *testing.T) {
	testForDuplicates(ClassNames, t)
}

func TestPlanUUID(t *testing.T) {
	testForDuplicates(PlanNames, t)
}
