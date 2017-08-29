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

1.  Getting a broker's 'catalog' of services: `Client.GetCatalog`
2.  Provisioning a new instance of a service: `Client.ProvisionInstance`
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

To provision a new instance of a service, call the `Client.Provision` method:

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

	return client.GetCatalog()
}
```


### Updating properties of an instance

### Deprovisioning an instance

### Checking the status of an asynchronous operation

### Binding to an instance

### Unbinding from an instance