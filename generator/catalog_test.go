/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generator

import (
	"encoding/json"
	"testing"

	"k8s.io/klog/v2"
)

func TestGetCatalog(t *testing.T) {
	g := &Generator{
		Services: []Service{
			{
				Plans: []Plan{
					{
						FromPool: Pull{
							Tags:     3,
							Metadata: 4,
							Free:     1,
						},
					},
					{
						FromPool: Pull{
							Tags:     3,
							Metadata: 4,
						},
					},
				},
				FromPool: Pull{
					Tags:                 3,
					Metadata:             4,
					InstancesRetrievable: 1,
					BindingsRetrievable:  1,
					Bindable:             1,
					Requires:             2,
				},
			},
		},
	}
	AssignPoolGoT(g)

	catalog, err := g.GetCatalog()
	if err != nil {
		t.Errorf("Got error, %v", err)
	}

	catalogBytes, _ := json.MarshalIndent(catalog, "", "  ")

	catalogJson := string(catalogBytes)

	klog.Info(catalogJson)
}

func TestGetPlans(t *testing.T) {

	g := Generator{
		PlanPool: []string{"AAA", "BBB", "CCC", "DDD", "EEE"},
	}
	klog.Info(g.planNames(1, 5))
	klog.Info(g.planNames(2, 5))
}
