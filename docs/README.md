# Documentation

The [Open Service Broker API](https://github.com/openservicebrokerapi/servicebroker)
 describes an entity (service broker) that provides some set of capabilities
(services).  Services have different *plans* that describe different tiers of
the service.  New instances of the services are *provisioned* in order to be
used.  Some services can be *bound* to applications for programmatic use.

Example:

- Service: "database as a service"
- Instance: "My database"
- Binding: "Credentials to use my database in app 'guestbook'"

## Background Reading

Reading the
[API specification](https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md) is 
recommended before reading this documentation.

## API Fundamentals

There are 7 operations in the API:

1.  Getting a broker's 'catalog' of services: [`Client.GetCatalog`](#getting-a-brokers-catalog)
2.  Provisioning a new instance of a service: [`Client.ProvisionInstance`](#provisioning-a-new-instance-of-a-service)
3.  Updating properties of an instance: `Client.UpdateInstance`
4.  Deprovisioning an instance: `Client.DeprovisionInstance`
5.  Checking the status of an asynchronous operation (provision, update, or deprovision) on an instance: `Client.PollLastOperation`
6.  Binding to an instance: `Client.Bind`
7.  Unbinding from an instance: `Client.Unbind`

### Getting a broker's catalog

A broker's catalog holds information about the services a broker provides and
their plans.  A platform implementing the OSB API must first get the broker's
catalog.

```go
import (
	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

func GetBrokerCatalog(URL string) (*osb.CatalogResponse, error) {
	config := osb.DefaultClientConfiguration()
	config.URL = URL

	client, err := osb.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client.GetCatalog()
}
```

### Provisioning a new instance of a service

To provision a new instance of a service, call the `Client.Provision` method.

Key points:

1. `ProvisionInstance` returns a response from the broker for successful
   operations, or an error if the broker returned an error response or
   there was a problem communicating with the broker
2. Use the `IsHTTPError` method to test and convert errors from Brokers
   into the standard broker error type, allowing access to conventional
   broker-provided fields
3. The `response.Async` field indicates whether the broker is performing the
   provision concurrently.  See the `LastOperation` method for information
   about handling asynchronous operations.

```go
import (
	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

func ProvisionService(client osb.Client, request osb.ProvisionRequest) (*osb.CatalogResponse, error) {
	config := osb.DefaultClientConfiguration()
	config.URL = URL

	client, err := osb.NewClient(config)
	if err != nil {
		return nil, err
	}

	request := &ProvisionRequest{
		InstanceID: "my-dbaas-service-instance",

		// Made up parameters for a hypothetical service
		ServiceID: "dbaas-service",
		PlanID:    "dbaas-gold-plan",
		Parameters: map[string]interface{}{
			"tablespace-page-cost":      100,
			"tablespace-io-concurrency": 5,
		},

		// Set the AcceptsIncomplete field to indicate that this client can
		// support asynchronous operations (provision, update, deprovision).
		AcceptsIncomplete: true,
	}

	// ProvisionInstance returns a response from the broker for successful
	// operations, or an error if the broker returned an error response or
	// there was a problem communicating with the broker.
	resp, err := client.ProvisionInstance(request)
	if err != nil {
		// Use the IsHTTPError method to test and convert errors from Brokers
		// into the standard broker error type, allowing access to conventional
		// broker-provided fields.
		errHttp, isError := osb.IsHTTPError(err)
		if isError {
			// handle error response from broker
		} else {
			// handle errors communicating with the broker
		}
	}

	// The response.Async field indicates whether the broker is performing the
	// provision concurrently.  See the LastOperation method for information
	// about handling asynchronous operations.
	if response.Async {
		// handle asynchronous operation
	}
}
```




### Updating properties of an instance

### Deprovisioning an instance

### Checking the status of an asynchronous operation

### Binding to an instance

### Unbinding from an instance