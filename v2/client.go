package v2

type APIVersion string

const (
	// XBrokerAPIVersion is the header for the Open Service Broker API.
	XBrokerAPIVersion = "X-Broker-Api-Version"
	// APIVersion2_11 represents the 2.11 version of the Open Service Broker
	// API.
	APIVersion2_11 APIVersion = "2.11"
)

// BasicAuthInfo represents a set of basic auth credentials.
type BasicAuthInfo struct {
	// Username is the basic auth username.
	Username string
	// Password is the basic auth password.
	Password string
}

// ClientConfiguration represents the configuration of a Client.
type ClientConfiguration struct {
	APIVersion          APIVersion
	BasicAuthInfo       *BasicAuthInfo
	EnableAlphaFeatures bool
}

// Client defines the interface to the v2 Open Service Broker client.  The
// logical lifecycle of client operations is:
//
// 1.  Get the broker's catalog of services with the GetCatalog method
// 2.  Provision a new instance of a service with the ProvisionInstance method
// 3.  Update the parameters or plan of an instance with the UpdateInstance method
// 4.  Deprovision an instance with the DeprovisionInstance method.
//
// Some services support binding from an instance of the service to an
// application.  The logical lifecycle of a binding is:
//
// 1.  Create a new binding to an instance of a service with the Bind method
// 2.  Delete a binding to an instance with the Unbind method
type Client interface {
	// GetCatalog calls GET on the Broker's catalog endpoint (/v2/catalog) and
	// returns a response and an error.
	GetCatalog() (*CatalogResponse, error)
	// ProvisionInstance sends a provision request to the broker and returns
	// the response or an error.
	ProvisionInstance(r ProvisionRequest) (*ProvisionResponse, error)
	// UpdateInstance sends an update instance request to the broker and
	// returns the response or an error.
	UpdateInstance(r UpdateInstanceRequest) (*UpdateInstanceResponse, error)
	// DeprovisionInstance sends a deprovision request to the broker and
	// returns the response or an error.
	DeprovisionInstance(r DeprovisionRequest) (*DeprovisionResponse, error)
	// PollLastOperation sends a request to query the last operation for a
	// service instance to the broker and returns the response or an error.
	PollLastOperation(r LastOperationRequest) (*LastOperationResponse, error)
	// Bind sends a bind request to the broker and returns the response or an
	// error.
	Bind(r BindRequest) (*BindResponse, error)
	// Unbind sends an unbind request to the broker and returns the response
	// or an error.
	Unbind(r UnbindRequest) (*UnbindResponse, error)
}

type client struct {
	URL           string
	BasicAuthInfo *BasicAuthInfo
}

func (c *client) GetCatalog() (*CatalogResponse, error) {
	return nil, nil
}

func (c *client) ProvisionInstance(r ProvisionRequest) (*ProvisionResponse, error) {
	return nil, nil
}

func (c *client) UpdateInstance(r UpdateInstanceRequest) (*UpdateInstanceResponse, error) {
	return nil, nil
}

func (c *client) DeprovisionInstance(r DeprovisionRequest) (*DeprovisionResponse, error) {
	return nil, nil
}

func (c *client) PollLastOperation(r LastOperationRequest) (*LastOperationResponse, error) {
	return nil, nil
}

func (c *client) Bind(r BindRequest) (*BindResponse, error) {
	return nil, nil
}

func (c *client) Unbind(r UnbindRequest) (*UnbindResponse, error) {
	return nil, nil
}
