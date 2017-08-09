package v2

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

type fakeOriginatingIdentity struct {
	platform   string
	value      string
	errMessage string
}

func (i *fakeOriginatingIdentity) Platform() string {
	return i.platform
}

func (i *fakeOriginatingIdentity) Value() (string, error) {
	if i.errMessage != "" {
		return "", errors.New(i.errMessage)
	}
	return i.value, nil
}

func TestCloudFoundryOriginatingIdentity(t *testing.T) {
	cases := []struct {
		name               string
		userId             string
		extra              map[string]interface{}
		expectedValue      map[string]interface{}
		expectedErrMessage string
	}{
		{
			name:          "only user_id",
			userId:        "user",
			expectedValue: map[string]interface{}{"user_id": "user"},
		},
		{
			name:               "missing user_id",
			userId:             "",
			expectedErrMessage: "UserId is required",
		},
		{
			name:   "additional properties",
			userId: "user",
			extra: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			expectedValue: map[string]interface{}{
				"user_id": "user",
				"key1":    "value1",
				"key2":    "value2",
			},
		},
	}

	for _, tc := range cases {
		oi := &CloudFoundryOriginatingIdentity{
			UserId: tc.userId,
			Extra:  tc.extra,
		}

		doOriginatingIdentityChecks(t, tc.name, oi, "cloudfoundry", tc.expectedValue, tc.expectedErrMessage)
	}
}

func TestKubernetesOriginatingIdentity(t *testing.T) {
	cases := []struct {
		name               string
		username           string
		uid                string
		groups             []string
		extra              map[string]interface{}
		expectedValue      map[string]interface{}
		expectedErrMessage string
	}{
		{
			name:     "only username and uid",
			username: "user",
			uid:      "1234",
			expectedValue: map[string]interface{}{
				"username": "user",
				"uid":      "1234",
			},
		},
		{
			name:               "missing username",
			username:           "",
			uid:                "1234",
			expectedErrMessage: "Username is required",
		},
		{
			name:               "missing username",
			username:           "user",
			uid:                "",
			expectedErrMessage: "Uid is required",
		},
		{
			name:     "groups",
			username: "user",
			uid:      "1234",
			groups:   []string{"group1", "group2"},
			expectedValue: map[string]interface{}{
				"username": "user",
				"uid":      "1234",
				"groups":   []interface{}{"group1", "group2"},
			},
		},
		{
			name:     "additional properties",
			username: "user",
			uid:      "1234",
			extra: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			expectedValue: map[string]interface{}{
				"username": "user",
				"uid":      "1234",
				"key1":     "value1",
				"key2":     "value2",
			},
		},
	}

	for _, tc := range cases {
		oi := &KubernetesOriginatingIdentity{
			Username: tc.username,
			Uid:      tc.uid,
			Groups:   tc.groups,
			Extra:    tc.extra,
		}

		doOriginatingIdentityChecks(t, tc.name, oi, "kubernetes", tc.expectedValue, tc.expectedErrMessage)
	}
}

func doOriginatingIdentityChecks(t *testing.T, name string, oi OriginatingIdentity, expectedPlatform string, expectedValue map[string]interface{}, expectedErrMessage string) {
	if e, a := expectedPlatform, oi.Platform(); e != a {
		t.Errorf("%v: unexpected platform: expected %v, actual %v", name, e, a)
		return
	}

	valueAsJson, err := oi.Value()
	if err == nil {
		if len(expectedErrMessage) != 0 {
			t.Errorf("%v: did not get expected error: %v", name, expectedErrMessage)
			return
		}
		value := map[string]interface{}{}
		err := json.Unmarshal([]byte(valueAsJson), &value)
		if err != nil {
			t.Errorf("%v: could not unmarshal value: value=%v, err=%v", name, valueAsJson, err)
			return
		}
		if e, a := expectedValue, value; !reflect.DeepEqual(e, a) {
			t.Errorf("%v: unexpected value: expected %v got %v", name, e, a)
			return
		}
	} else {
		if e, a := expectedErrMessage, err.Error(); e != a {
			t.Errorf("%v: unexpected error: expected %v got %v", name, e, a)
			return
		}
	}
}

func TestBuildOriginatingIdentityHeaderValue(t *testing.T) {
	cases := []struct {
		name                string
		originatingIdentity *fakeOriginatingIdentity
		platform            string
		value               string
		valueErrMessage     string
		expectedHeaderValue string
		expectedErrMessage  string
	}{
		{
			name: "valid originating identity",
			originatingIdentity: &fakeOriginatingIdentity{
				platform: "fakeplatform",
				value:    "{\"user\":\"name\"}",
			},
			expectedHeaderValue: "fakeplatform eyJ1c2VyIjoibmFtZSJ9",
		},
		{
			name: "value raising error",
			originatingIdentity: &fakeOriginatingIdentity{
				platform:   "fakeplatform",
				errMessage: "fakeerror",
			},
			expectedErrMessage: "fakeerror",
		},
	}
	for _, tc := range cases {
		headerValue, err := buildOriginatingIdentityHeaderValue(tc.originatingIdentity)
		if tc.expectedErrMessage != "" {
			if err == nil {
				t.Errorf("%v: expected error not thrown: expected %v", tc.name, tc.expectedErrMessage)
				return
			}
			if e, a := tc.expectedErrMessage, err.Error(); e != a {
				t.Errorf("%v: unexpected error message: expected %v, got %v", tc.name, e, a)
				return
			}
		} else if err != nil {
			t.Errorf("%v: unexpected error: %v", tc.name, err)
			return
		} else {
			if e, a := tc.expectedHeaderValue, headerValue; e != a {
				t.Errorf("%v: unexpected header value: expected %v, got %v", tc.name, e, a)
			}
		}
	}
}
