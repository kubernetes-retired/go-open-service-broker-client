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

package v2

import (
	"fmt"
	"net/http"
)

func (c *client) GetCatalog() (*CatalogResponse, error) {
	fullURL := fmt.Sprintf(catalogURL, c.URL)

	response, err := c.prepareAndDo(http.MethodGet, fullURL, nil /* params */, nil /* request body */, nil /* originating identity */)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = drainReader(response.Body)
		response.Body.Close()
	}()

	switch response.StatusCode {
	case http.StatusOK:
		catalogResponse := &CatalogResponse{}
		if err := c.unmarshalResponse(response, catalogResponse); err != nil {
			return nil, HTTPStatusCodeError{StatusCode: response.StatusCode, ResponseError: err}
		}

		if c.APIVersion.IsLessThan(Version2_13()) || !c.EnableAlphaFeatures {
			c.pruneCatalogResponse(catalogResponse)
		}

		return catalogResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}
}

func (c *client) pruneCatalogResponse(catalogResponse *CatalogResponse) {
	for ii := range catalogResponse.Services {
		for jj := range catalogResponse.Services[ii].Plans {
			if c.APIVersion.IsLessThan(Version2_13()) {
				catalogResponse.Services[ii].Plans[jj].Schemas = nil
			}
			if !c.EnableAlphaFeatures {
				catalogResponse.Services[ii].Plans[jj].MaintenanceInfo = nil
				catalogResponse.Services[ii].Plans[jj].MaximumPollingDuration = nil
				catalogResponse.Services[ii].Plans[jj].PlanUpdateable = nil
			}
		}
	}
}
