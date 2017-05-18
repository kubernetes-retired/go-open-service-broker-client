package v2

import (
	"fmt"
	"net/http"
)

func (c *client) GetCatalog() (*CatalogResponse, error) {
	fullURL := fmt.Sprintf(catalogURL, c.URL)

	response, err := c.prepareAndDoFunc(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		catalogResponse := &CatalogResponse{}
		if err := c.unmarshalResponse(response, catalogResponse); err != nil {
			return nil, err
		}
		return catalogResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}
}
