package v2

import (
	"crypto/tls"
)

type APIVersion string

const (
	// APIVersion2_11 represents the 2.11 version of the Open Service Broker
	// API.
	APIVersion2_11 APIVersion = "2.11"
)

// LatestAPIVersion returns the latest supported API version in the current
// release of this library.
func LatestAPIVersion() APIVersion {
	return APIVersion2_11
}

// AuthConfig is a union-type representing the possible auth configurations a
// client may use to authenticate to a broker.  Currently, only basic auth is
// supported.
type AuthConfig struct {
	BasicAuthConfig *BasicAuthConfig
}

// BasicAuthConfig represents a set of basic auth credentials.
type BasicAuthConfig struct {
	// Username is the basic auth username.
	Username string
	// Password is the basic auth password.
	Password string
}

// ClientConfiguration represents the configuration of a Client.
type ClientConfiguration struct {
	// Name is the name to use for this client in log messages.  Using the
	// logical name of the Broker this client is for is recommended.
	Name string
	// URL is the URL to use to contact the broker.
	URL string
	// APIVersion is the APIVersion to use for this client.  API features
	// adopted after the 2.11 version of the API will only be sent if
	// APIVersion is an API version that supports them.
	APIVersion APIVersion
	// AuthInfo is the auth configuration the client should use to authenticate
	// to the broker.
	AuthConfig *AuthConfig
	// TLSConfig is the TLS configuration to use when communicating with the
	// broker.
	TLSConfig *tls.Config
	// Insecure represents whether the 'InsecureSkipVerify' TLS configuration
	// field should be set.  If the TLSConfig field is set and this field is
	// set to true, it overrides the value in the TLSConfig field.
	Insecure bool
	// TimeoutSeconds is the length of the timeout of any request to the
	// broker, in seconds.
	TimeoutSeconds int
	// EnableAlphaFeatures controls whether alpha features in the Open Service
	// Broker API are enabled in a client.  Features are considered to be
	// alpha if they have been accepted into the Open Service Broker API but
	// not released in a version of the API specification.  Features are
	// indicated as being alpha when the client API fields they represent
	// begin with the 'Alpha' prefix.
	//
	// If alpha features are not enabled, the client will not send or return
	// any request parameters or request or response fields that correspond to
	// alpha features.
	EnableAlphaFeatures bool
}

// DefaultClientConfiguration returns a default ClientConfiguration:
//
// - latest API version
// - 60 second timeout (referenced as a typical timeout in the Open Service
//   Broker API spec)
// - alpha features disabled
func DefaultClientConfiguration() *ClientConfiguration {
	return &ClientConfiguration{
		APIVersion:          LatestAPIVersion(),
		TimeoutSeconds:      60,
		EnableAlphaFeatures: false,
	}
}

// Client defines the interface to the v2 Open Service Broker client.  The
// logical lifecycle of client operations is:
//
// 1.  Get the broker's catalog of services with the GetCatalog method
// 2.  Provision a new instance of a service with the ProvisionInstance method
// 3.  Update the parameters or plan of an instance with the UpdateInstance method
// 4.  Deprovision an instance with the DeprovisionInstance method
//
// Some services and plans support binding from an instance of the service to
// an application.  The logical lifecycle of a binding is:
//
// 1.  Create a new binding to an instance of a service with the Bind method
// 2.  Delete a binding to an instance with the Unbind method
type Client interface {
	// GetCatalog calls GET on the Broker's catalog endpoint (/v2/catalog) and
	// returns a response and an error.
	GetCatalog() (*CatalogResponse, error)
	// ProvisionInstance sends a provision request to the broker and returns
	// the response or an error.
	ProvisionInstance(r *ProvisionRequest) (*ProvisionResponse, error)
	// UpdateInstance sends an update instance request to the broker and
	// returns the response or an error.
	UpdateInstance(r *UpdateInstanceRequest) (*UpdateInstanceResponse, error)
	// DeprovisionInstance sends a deprovision request to the broker and
	// returns the response or an error.
	DeprovisionInstance(r *DeprovisionRequest) (*DeprovisionResponse, error)
	// PollLastOperation sends a request to query the last operation for a
	// service instance to the broker and returns the response or an error.
	PollLastOperation(r *LastOperationRequest) (*LastOperationResponse, error)
	// Bind sends a bind request to the broker and returns the response or an
	// error.
	Bind(r *BindRequest) (*BindResponse, error)
	// Unbind sends an unbind request to the broker and returns the response
	// or an error.
	Unbind(r *UnbindRequest) (*UnbindResponse, error)
}

// CreateFunc allows control over which implementation of a Client is
// returned.  Users of the Client interface may need to create clients for
// multiple brokers in a way that makes normal dependency injection
// prohibitive.  In order to make such code testable, users of the API can
// inject a CreateFunc, and use the CreateFunc from the fake package in tests.
type CreateFunc func(*ClientConfiguration) (Client, error)
