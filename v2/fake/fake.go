package v2

import (
	"errors"

	"github.com/pmorie/go-open-service-broker-client/v2"
)

// NewFakeClientFunc generates a v2.CreateFunc that returns a FakeClient with
// the given FakeClientConfiguration.  It is useful for injecting the
// FakeClient in code that uses the v2.CreateFunc interface.
func NewFakeClientFunc(config FakeClientConfiguration) v2.CreateFunc {
	return func(_ *v2.ClientConfiguration) (v2.Client, error) {
		return &FakeClient{
			CatalogReaction:           config.CatalogReaction,
			ProvisionReaction:         config.ProvisionReaction,
			UpdateInstanceReaction:    config.UpdateInstanceReaction,
			DeprovisionReaction:       config.DeprovisionReaction,
			PollLastOperationReaction: config.PollLastOperationReaction,
			BindReaction:              config.BindReaction,
			UnbindReaction:            config.UnbindReaction,
		}, nil
	}
}

type FakeClientConfiguration struct {
	CatalogReaction           *CatalogReaction
	ProvisionReaction         *ProvisionReaction
	UpdateInstanceReaction    *UpdateInstanceReaction
	DeprovisionReaction       *DeprovisionReaction
	PollLastOperationReaction *PollLastOperationReaction
	BindReaction              *BindReaction
	UnbindReaction            *UnbindReaction
}

type FakeClient struct {
	CatalogReaction           *CatalogReaction
	ProvisionReaction         *ProvisionReaction
	UpdateInstanceReaction    *UpdateInstanceReaction
	DeprovisionReaction       *DeprovisionReaction
	PollLastOperationReaction *PollLastOperationReaction
	BindReaction              *BindReaction
	UnbindReaction            *UnbindReaction
}

var _ v2.Client = &FakeClient{}

func (c *FakeClient) GetCatalog() (*v2.CatalogResponse, error) {
	if c.CatalogReaction != nil {
		return c.CatalogReaction.Response, c.CatalogReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) ProvisionInstance(r *v2.ProvisionRequest) (*v2.ProvisionResponse, error) {
	if c.ProvisionReaction != nil {
		return c.ProvisionReaction.Response, c.ProvisionReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) UpdateInstance(r *v2.UpdateInstanceRequest) (*v2.UpdateInstanceResponse, error) {
	if c.UpdateInstanceReaction != nil {
		return c.UpdateInstanceReaction.Response, c.UpdateInstanceReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) DeprovisionInstance(r *v2.DeprovisionRequest) (*v2.DeprovisionResponse, error) {
	if c.DeprovisionReaction != nil {
		return c.DeprovisionReaction.Response, c.DeprovisionReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) PollLastOperation(r *v2.LastOperationRequest) (*v2.LastOperationResponse, error) {
	if c.PollLastOperationReaction != nil {
		return c.PollLastOperationReaction.Response, c.PollLastOperationReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) Bind(r *v2.BindRequest) (*v2.BindResponse, error) {
	if c.BindReaction != nil {
		return c.BindReaction.Response, c.BindReaction.Error
	}

	return nil, UnexpectedActionError()
}

func (c *FakeClient) Unbind(r *v2.UnbindRequest) (*v2.UnbindResponse, error) {
	if c.UnbindReaction != nil {
		return c.UnbindReaction.Response, c.UnbindReaction.Error
	}

	return nil, UnexpectedActionError()
}

func UnexpectedActionError() error {
	return errors.New("Unexpected action")
}

type CatalogReaction struct {
	Response *v2.CatalogResponse
	Error    error
}

type ProvisionReaction struct {
	Response *v2.ProvisionResponse
	Error    error
}

type UpdateInstanceReaction struct {
	Response *v2.UpdateInstanceResponse
	Error    error
}

type DeprovisionReaction struct {
	Response *v2.DeprovisionResponse
	Error    error
}

type PollLastOperationReaction struct {
	Response *v2.LastOperationResponse
	Error    error
}

type BindReaction struct {
	Response *v2.BindResponse
	Error    error
}

type UnbindReaction struct {
	Response *v2.UnbindResponse
	Error    error
}
