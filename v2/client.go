package v2

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
)

const (
	// XBrokerAPIVersion is the header for the Open Service Broker API
	// version.
	XBrokerAPIVersion = "X-Broker-Api-Version"

	catalogURL            = "%s/v2/catalog"
	serviceInstanceURLFmt = "%s/v2/service_instances/%s"
	lastOperationURLFmt   = "%s/v2/service_instances/%s/last_operation"
	bindingURLFmt         = "%s/v2/service_instances/%s/service_bindings/%s"
	queryParamFmt         = "%s=%s"
)

// NewClient is a CreateFunc for creating a new functional Client and
// implements the CreateFunc interface.
func NewClient(config *ClientConfiguration) (Client, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(config.TimeoutSeconds) * time.Second,
	}
	transport := &http.Transport{}
	if config.TLSConfig != nil {
		transport.TLSClientConfig = config.TLSConfig
	} else {
		transport.TLSClientConfig = &tls.Config{}
	}
	if config.Insecure {
		transport.TLSClientConfig.InsecureSkipVerify = true
	}
	httpClient.Transport = transport

	c := &client{
		Name:                config.Name,
		URL:                 strings.TrimRight(config.URL, "/"),
		APIVersion:          config.APIVersion,
		EnableAlphaFeatures: config.EnableAlphaFeatures,
		httpClient:          httpClient,
	}

	if config.AuthConfig != nil {
		if config.AuthConfig.BasicAuthConfig == nil {
			return nil, errors.New("BasicAuthConfig is required is AuthConfig is provided")
		}

		c.BasicAuthConfig = config.AuthConfig.BasicAuthConfig
	}

	return c, nil
}

var _ CreateFunc = NewClient

// client provides a functional implementation of the Client interface.
type client struct {
	Name                string
	URL                 string
	APIVersion          APIVersion
	BasicAuthConfig     *BasicAuthConfig
	EnableAlphaFeatures bool
	Verbose             bool

	httpClient *http.Client
}

var _ Client = &client{}

func (c *client) GetCatalog() (*CatalogResponse, error) {
	fullURL := fmt.Sprintf(catalogURL, c.URL)

	response, err := c.prepareAndDoRequest(http.MethodGet, fullURL, nil)
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

func (c *client) ProvisionInstance(r *ProvisionRequest) (*ProvisionResponse, error) {
	if err := validateProvisionRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(serviceInstanceURLFmt, c.URL, r.InstanceID)
	if r.AcceptsIncomplete {
		fullURL += "?accepts_incomplete=true"
	}

	requestBody := &provisionRequestBody{
		serviceID:        r.ServiceID,
		planID:           r.PlanID,
		organizationGUID: r.OrganizationGUID,
		spaceGUID:        r.SpaceGUID,
		parameters:       r.Parameters,
	}

	response, err := c.prepareAndDoRequest(http.MethodPut, fullURL, requestBody)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusCreated, http.StatusOK, http.StatusAccepted:
		responseBodyObj := &provisionSuccessResponseBody{}
		if err := c.unmarshalResponse(response, responseBodyObj); err != nil {
			return nil, err
		}

		var opPtr *OperationKey
		if responseBodyObj.operation != nil {
			opStr := *responseBodyObj.operation
			op := OperationKey(opStr)
			opPtr = &op
		}

		userResponse := &ProvisionResponse{
			DashboardURL: responseBodyObj.dashboardURL,
			OperationKey: opPtr,
		}
		if response.StatusCode == http.StatusAccepted {
			if c.Verbose {
				glog.Infof("broker %q: received asynchronous response", c.Name)
			}
			userResponse.Async = true
		}

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}
}

func (c *client) UpdateInstance(r *UpdateInstanceRequest) (*UpdateInstanceResponse, error) {
	if err := validateUpdateInstanceRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(serviceInstanceURLFmt, c.URL, r.InstanceID)
	if r.AcceptsIncomplete {
		fullURL += "?accepts_incomplete=true"
	}

	requestBody := &updateInstanceRequestBody{
		serviceID:  r.ServiceID,
		planID:     r.PlanID,
		parameters: r.Parameters,
	}

	response, err := c.prepareAndDoRequest(http.MethodPatch, fullURL, requestBody)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return &UpdateInstanceResponse{}, nil
	case http.StatusAccepted:
		responseBodyObj := &asyncSuccessResponseBody{}
		if err := c.unmarshalResponse(response, responseBodyObj); err != nil {
			return nil, err
		}

		var opPtr *OperationKey
		if responseBodyObj.operation != nil {
			opStr := *responseBodyObj.operation
			op := OperationKey(opStr)
			opPtr = &op
		}

		userResponse := &UpdateInstanceResponse{
			Async:        true,
			OperationKey: opPtr,
		}

		// TODO: fix op key handling

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}

func (c *client) DeprovisionInstance(r *DeprovisionRequest) (*DeprovisionResponse, error) {
	if err := validateDeprovisionRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(serviceInstanceURLFmt, c.URL, r.InstanceID)
	if r.AcceptsIncomplete {
		fullURL += "?accepts_incomplete=true"
	}

	requestServiceID := string(r.ServiceID)
	requestPlanID := string(r.PlanID)

	requestBody := &deprovisionInstanceRequestBody{
		serviceID: &requestServiceID,
		planID:    &requestPlanID,
	}

	response, err := c.prepareAndDoRequest(http.MethodDelete, fullURL, requestBody)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusGone:
		return &DeprovisionResponse{}, nil
	case http.StatusAccepted:
		responseBodyObj := &asyncSuccessResponseBody{}
		if err := c.unmarshalResponse(response, responseBodyObj); err != nil {
			return nil, err
		}

		var opPtr *OperationKey
		if responseBodyObj.operation != nil {
			opStr := *responseBodyObj.operation
			op := OperationKey(opStr)
			opPtr = &op
		}

		userResponse := &DeprovisionResponse{
			Async:        true,
			OperationKey: opPtr,
		}

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}

const (
	serviceIDKey = "service_id"
	planIDKey    = "plan_id"
	operationKey = "operation"
)

func (c *client) PollLastOperation(r *LastOperationRequest) (*LastOperationResponse, error) {
	// TODO: support special handling for delete responses

	if err := validateLastOperationRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(lastOperationURLFmt, c.URL, r.InstanceID)
	var queryParamBuffer bytes.Buffer
	switch {
	case r.ServiceID != nil:
		appendQueryParam(&queryParamBuffer, serviceIDKey, *r.ServiceID)
		fallthrough
	case r.PlanID != nil:
		appendQueryParam(&queryParamBuffer, planIDKey, *r.PlanID)
		fallthrough
	case r.OperationKey != nil:
		op := *r.OperationKey
		opStr := string(op)
		appendQueryParam(&queryParamBuffer, operationKey, opStr)
	}
	if queryParamBuffer.Len() > 0 {
		fullURL += "?" + queryParamBuffer.String()
	}

	response, err := c.prepareAndDoRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		userResponse := &LastOperationResponse{}
		if err := c.unmarshalResponse(response, userResponse); err != nil {
			return nil, err
		}

		return userResponse, nil
	case http.StatusGone:
		// TODO: async operations for deprovision have a special case to be
		// handled here
		fallthrough
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}

const (
	bindResourceAppGUIDKey = "app_guid"
	bindResourceRouteKey   = "route"
)

func (c *client) Bind(r *BindRequest) (*BindResponse, error) {
	if err := validateBindRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(bindingURLFmt, c.URL, r.InstanceID, r.BindingID)

	requestBody := &bindRequestBody{
		serviceID:  r.ServiceID,
		planID:     r.PlanID,
		parameters: r.Parameters,
	}

	if r.BindResource != nil {
		requestBody.bindResource = map[string]interface{}{}
		if r.BindResource.AppGUID != nil {
			requestBody.bindResource[bindResourceAppGUIDKey] = *r.BindResource.AppGUID
		}
		if r.BindResource.Route != nil {
			requestBody.bindResource[bindResourceRouteKey] = *r.BindResource.AppGUID
		}
	}

	response, err := c.prepareAndDoRequest(http.MethodPut, fullURL, requestBody)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		userResponse := &BindResponse{}
		if err := c.unmarshalResponse(response, userResponse); err != nil {
			return nil, err
		}

		return userResponse, nil
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}

func (c *client) Unbind(r *UnbindRequest) (*UnbindResponse, error) {
	if err := validateUnbindRequest(r); err != nil {
		return nil, err
	}

	fullURL := fmt.Sprintf(bindingURLFmt, c.URL, r.InstanceID, r.BindingID)
	var queryParamBuffer bytes.Buffer
	appendQueryParam(&queryParamBuffer, serviceIDKey, r.ServiceID)
	appendQueryParam(&queryParamBuffer, planIDKey, r.PlanID)
	fullURL += "?" + queryParamBuffer.String()

	response, err := c.prepareAndDoRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusGone:
		// TODO: should we establish that the response body ('{}') is correct?
		return &UnbindResponse{}, nil
	default:
		return nil, c.handleFailureResponse(response)
	}

	return nil, nil
}

const (
	contentType = "Content-Type"
	jsonType    = "application/json"
)

// prepareAndDoRequest prepares a request for the given method, URL, and
// message body, and executes the request, returning an http.Response or an
// error.  Errors returned from this function represent http-layer errors and
// not errors in the Open Service Broker API.
func (c *client) prepareAndDoRequest(method, URL string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		bodyReader = bytes.NewReader(bodyBytes)
	}

	request, err := http.NewRequest(method, URL, bodyReader)
	if err != nil {
		return nil, err
	}

	request.Header.Set(XBrokerAPIVersion, string(c.APIVersion))
	if bodyReader != nil {
		request.Header.Set(contentType, jsonType)
	}
	if c.BasicAuthConfig != nil {
		request.SetBasicAuth(c.BasicAuthConfig.Username, c.BasicAuthConfig.Password)
	}

	if c.Verbose {
		glog.Infof("broker %q: doing request to %q", c.Name, URL)
	}

	return c.httpClient.Do(request)
}

// appendQueryParam appends key=value to buffer if value is non-null.
// If buffer is non-empty appends &key=value
func appendQueryParam(buffer *bytes.Buffer, key, value string) error {
	if value == "" {
		return nil
	}
	if buffer.Len() > 0 {
		_, err := buffer.WriteString("&")
		if err != nil {
			return err
		}
	}
	_, err := buffer.WriteString(fmt.Sprintf(queryParamFmt, key, value))
	return err
}

// unmarshalResponse unmartials the response body of the given response into
// the given object or returns an error.
func (c *client) unmarshalResponse(response *http.Response, obj interface{}) error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if c.Verbose {
		glog.Info("broker %q: response body: %v", c.Name, string(body))
	}

	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}

	return nil
}

// handleFailureResponse returns an HTTPStatusCodeError for the given
// response.
func (c *client) handleFailureResponse(response *http.Response) error {
	brokerResponse := &failureResponseBody{}
	if err := c.unmarshalResponse(response, brokerResponse); err != nil {
		return err
	}

	return HTTPStatusCodeError{
		StatusCode:   response.StatusCode,
		ErrorMessage: brokerResponse.err,
		Description:  brokerResponse.description,
	}
}

// internal message body types

type provisionRequestBody struct {
	serviceID        string                 `json:"service_id"`
	planID           string                 `json:"plan_id"`
	organizationGUID string                 `json:"organization_guid"`
	spaceGUID        string                 `json:"space_guid"`
	parameters       map[string]interface{} `json:"parameters,omitempty"`
}

type provisionSuccessResponseBody struct {
	dashboardURL *string `json:"dashboard_url"`
	operation    *string `json:"operation"`
}

type deprovisionInstanceRequestBody struct {
	serviceID *string `json:"service_id"`
	planID    *string `json:"plan_id,omitempty"`
}

type updateInstanceRequestBody struct {
	serviceID  string                 `json:"service_id"`
	planID     *string                `json:"plan_id,omitempty"`
	parameters map[string]interface{} `json:"parameters,omitempty"`

	// Note: this client does not currently support the 'previous_values'
	// field of this request body.
}

type bindRequestBody struct {
	serviceID    string                 `json:"service_id"`
	planID       string                 `json:"plan_id"`
	parameters   map[string]interface{} `json:"parameters,omitempty"`
	bindResource map[string]interface{} `json:"bind_resource,omitempty"`
}

type asyncSuccessResponseBody struct {
	operation *string `json:"operation"`
}

type failureResponseBody struct {
	err         *string `json:"error,omitempty"`
	description *string `json:"description,omitempty"`
}
