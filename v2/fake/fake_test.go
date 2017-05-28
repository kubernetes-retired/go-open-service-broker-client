package fake_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pmorie/go-open-service-broker-client/v2"
	"github.com/pmorie/go-open-service-broker-client/v2/fake"
)

func catalogResponse() *v2.CatalogResponse {
	return &v2.CatalogResponse{
		Services: []v2.Service{
			{
				ID:          "acb56d7c-XXXX-XXXX-XXXX-feb140a59a66",
				Name:        "fake-service",
				Description: "fake service",
				Tags: []string{
					"tag1",
					"tag2",
				},
				Requires: []string{
					"route_forwarding",
				},
				Bindable:      true,
				PlanUpdatable: truePtr(),
				Plans: []v2.Plan{
					{
						ID:          "d3031751-XXXX-XXXX-XXXX-a42377d3320e",
						Name:        "fake-plan-1",
						Description: "description1",
						Metadata: map[string]interface{}{
							"b": "c",
							"d": "e",
						},
					},
				},
				DashboardClient: &v2.DashboardClient{
					ID:          "398e2f8e-XXXX-XXXX-XXXX-19a71ecbcf64",
					Secret:      "277cabb0-XXXX-XXXX-XXXX-7822c0a90e5d",
					RedirectURI: "http://localhost:1234",
				},
				Metadata: map[string]interface{}{
					"a": "b",
					"c": "d",
				},
			},
		},
	}
}

func truePtr() *bool {
	b := true
	return &b
}

func TestGetCatalog(t *testing.T) {
	cases := []struct {
		name     string
		reaction *fake.CatalogReaction
		response *v2.CatalogResponse
		err      error
	}{
		{
			name: "unexpected action",
			err:  fake.UnexpectedActionError(),
		},
		{
			name: "response",
			reaction: &fake.CatalogReaction{
				Response: catalogResponse(),
			},
			response: catalogResponse(),
		},
		{
			name: "error",
			reaction: &fake.CatalogReaction{
				Error: errors.New("oops"),
			},
			err: errors.New("oops"),
		},
	}

	for _, tc := range cases {
		fakeClient := &fake.FakeClient{
			CatalogReaction: tc.reaction,
		}

		response, err := fakeClient.GetCatalog()

		if !reflect.DeepEqual(tc.response, response) {
			t.Errorf("%v: unexpected response; expected %+v, got %+v", tc.name, tc.response, response)
		}

		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("%v: unexpected error; expected %+v, got %+v", tc.name, tc.err, err)
		}
	}
}

func provisionResponse() *v2.ProvisionResponse {
	return &v2.ProvisionResponse{
		Async: true,
	}
}

func TestProvisionInstance(t *testing.T) {
	cases := []struct {
		name     string
		reaction *fake.ProvisionReaction
		response *v2.ProvisionResponse
		err      error
	}{
		{
			name: "unexpected action",
			err:  fake.UnexpectedActionError(),
		},
		{
			name: "response",
			reaction: &fake.ProvisionReaction{
				Response: provisionResponse(),
			},
			response: provisionResponse(),
		},
		{
			name: "error",
			reaction: &fake.ProvisionReaction{
				Error: errors.New("oops"),
			},
			err: errors.New("oops"),
		},
	}

	for _, tc := range cases {
		fakeClient := &fake.FakeClient{
			ProvisionReaction: tc.reaction,
		}

		response, err := fakeClient.ProvisionInstance(&v2.ProvisionRequest{})

		if !reflect.DeepEqual(tc.response, response) {
			t.Errorf("%v: unexpected response; expected %+v, got %+v", tc.name, tc.response, response)
		}

		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("%v: unexpected error; expected %+v, got %+v", tc.name, tc.err, err)
		}
	}
}

func updateInstanceResponse() *v2.UpdateInstanceResponse {
	return &v2.UpdateInstanceResponse{
		Async: true,
	}
}

func TestUpdateInstance(t *testing.T) {
	cases := []struct {
		name     string
		reaction *fake.UpdateInstanceReaction
		response *v2.UpdateInstanceResponse
		err      error
	}{
		{
			name: "unexpected action",
			err:  fake.UnexpectedActionError(),
		},
		{
			name: "response",
			reaction: &fake.UpdateInstanceReaction{
				Response: updateInstanceResponse(),
			},
			response: updateInstanceResponse(),
		},
		{
			name: "error",
			reaction: &fake.UpdateInstanceReaction{
				Error: errors.New("oops"),
			},
			err: errors.New("oops"),
		},
	}

	for _, tc := range cases {
		fakeClient := &fake.FakeClient{
			UpdateInstanceReaction: tc.reaction,
		}

		response, err := fakeClient.UpdateInstance(&v2.UpdateInstanceRequest{})

		if !reflect.DeepEqual(tc.response, response) {
			t.Errorf("%v: unexpected response; expected %+v, got %+v", tc.name, tc.response, response)
		}

		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("%v: unexpected error; expected %+v, got %+v", tc.name, tc.err, err)
		}
	}
}

func deprovisionResponse() *v2.DeprovisionResponse {
	return &v2.DeprovisionResponse{
		Async: true,
	}
}

func TestDeprovisionInstance(t *testing.T) {
	cases := []struct {
		name     string
		reaction *fake.DeprovisionReaction
		response *v2.DeprovisionResponse
		err      error
	}{
		{
			name: "unexpected action",
			err:  fake.UnexpectedActionError(),
		},
		{
			name: "response",
			reaction: &fake.DeprovisionReaction{
				Response: deprovisionResponse(),
			},
			response: deprovisionResponse(),
		},
		{
			name: "error",
			reaction: &fake.DeprovisionReaction{
				Error: errors.New("oops"),
			},
			err: errors.New("oops"),
		},
	}

	for _, tc := range cases {
		fakeClient := &fake.FakeClient{
			DeprovisionReaction: tc.reaction,
		}

		response, err := fakeClient.DeprovisionInstance(&v2.DeprovisionRequest{})

		if !reflect.DeepEqual(tc.response, response) {
			t.Errorf("%v: unexpected response; expected %+v, got %+v", tc.name, tc.response, response)
		}

		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("%v: unexpected error; expected %+v, got %+v", tc.name, tc.err, err)
		}
	}
}

func lastOperationResponse() *v2.LastOperationResponse {
	return &v2.LastOperationResponse{
		State: v2.StateSucceeded,
	}
}

func TestPollLastOperation(t *testing.T) {
	cases := []struct {
		name     string
		reaction *fake.PollLastOperationReaction
		response *v2.LastOperationResponse
		err      error
	}{
		{
			name: "unexpected action",
			err:  fake.UnexpectedActionError(),
		},
		{
			name: "response",
			reaction: &fake.PollLastOperationReaction{
				Response: lastOperationResponse(),
			},
			response: lastOperationResponse(),
		},
		{
			name: "error",
			reaction: &fake.PollLastOperationReaction{
				Error: errors.New("oops"),
			},
			err: errors.New("oops"),
		},
	}

	for _, tc := range cases {
		fakeClient := &fake.FakeClient{
			PollLastOperationReaction: tc.reaction,
		}

		response, err := fakeClient.PollLastOperation(&v2.LastOperationRequest{})

		if !reflect.DeepEqual(tc.response, response) {
			t.Errorf("%v: unexpected response; expected %+v, got %+v", tc.name, tc.response, response)
		}

		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("%v: unexpected error; expected %+v, got %+v", tc.name, tc.err, err)
		}
	}
}

func bindResponse() *v2.BindResponse {
	return &v2.BindResponse{
		Credentials: map[string]interface{}{
			"foo": "bar",
		},
	}
}

func TestBind(t *testing.T) {
	cases := []struct {
		name     string
		reaction *fake.BindReaction
		response *v2.BindResponse
		err      error
	}{
		{
			name: "unexpected action",
			err:  fake.UnexpectedActionError(),
		},
		{
			name: "response",
			reaction: &fake.BindReaction{
				Response: bindResponse(),
			},
			response: bindResponse(),
		},
		{
			name: "error",
			reaction: &fake.BindReaction{
				Error: errors.New("oops"),
			},
			err: errors.New("oops"),
		},
	}

	for _, tc := range cases {
		fakeClient := &fake.FakeClient{
			BindReaction: tc.reaction,
		}

		response, err := fakeClient.Bind(&v2.BindRequest{})

		if !reflect.DeepEqual(tc.response, response) {
			t.Errorf("%v: unexpected response; expected %+v, got %+v", tc.name, tc.response, response)
		}

		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("%v: unexpected error; expected %+v, got %+v", tc.name, tc.err, err)
		}
	}
}

func unbindResponse() *v2.UnbindResponse {
	return &v2.UnbindResponse{}
}

func TestUnbind(t *testing.T) {
	cases := []struct {
		name     string
		reaction *fake.UnbindReaction
		response *v2.UnbindResponse
		err      error
	}{
		{
			name: "unexpected action",
			err:  fake.UnexpectedActionError(),
		},
		{
			name: "response",
			reaction: &fake.UnbindReaction{
				Response: unbindResponse(),
			},
			response: unbindResponse(),
		},
		{
			name: "error",
			reaction: &fake.UnbindReaction{
				Error: errors.New("oops"),
			},
			err: errors.New("oops"),
		},
	}

	for _, tc := range cases {
		fakeClient := &fake.FakeClient{
			UnbindReaction: tc.reaction,
		}

		response, err := fakeClient.Unbind(&v2.UnbindRequest{})

		if !reflect.DeepEqual(tc.response, response) {
			t.Errorf("%v: unexpected response; expected %+v, got %+v", tc.name, tc.response, response)
		}

		if !reflect.DeepEqual(tc.err, err) {
			t.Errorf("%v: unexpected error; expected %+v, got %+v", tc.name, tc.err, err)
		}
	}
}
