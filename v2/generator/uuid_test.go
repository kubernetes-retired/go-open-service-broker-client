package generator

import (
	"testing"
)

func testForDuplicates(list []string, t *testing.T) {
	uuids := map[string]string{}
	for _, name := range list {
		uuid := newUUID(name)
		if uuids[uuid] != "" {
			t.Fatalf("Found a collision with %s and %s", uuids[uuid], name)
		} else {
			uuids[uuid] = name
		}
	}
}

func TestServiceClassUUID(t *testing.T) {
	testForDuplicates(ServiceClassNames, t)
}

func TestServiceClassIDMatch(t *testing.T) {
	if ServiceClassID(0) != "" && ServiceClassID(0) != ServiceClassID(0) {
		t.Fatalf("ID did not match on second call.")
	}
}

func TestServicePlanUUID(t *testing.T) {
	testForDuplicates(ServicePlanNames, t)
}

func TestServicePlanIDMatch(t *testing.T) {
	if ServicePlanID(0) != "" && ServicePlanID(0) != ServicePlanID(0) {
		t.Fatalf("ID did not match on second call.")
	}
}
