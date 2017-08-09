package v2

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

func buildOriginatingIdentityHeaderValue(i OriginatingIdentity) (string, error) {
	if i == nil {
		return "", nil
	}
	platform := i.Platform()
	value, err := i.Value()
	if err != nil {
		return "", err
	}
	encodedValue := base64.StdEncoding.EncodeToString([]byte(value))
	headerValue := fmt.Sprintf("%v %v", platform, encodedValue)
	return headerValue, nil
}

var _ OriginatingIdentity = &CloudFoundryOriginatingIdentity{}

func (i *CloudFoundryOriginatingIdentity) Platform() string {
	return "cloudfoundry"
}

func (i *CloudFoundryOriginatingIdentity) Value() (string, error) {
	if err := i.validate(); err != nil {
		return "", err
	}
	baseProperties := map[string]interface{}{
		"user_id": i.UserId,
	}
	return marshalValue(baseProperties, i.Extra)
}

func (i *CloudFoundryOriginatingIdentity) validate() error {
	if len(i.UserId) == 0 {
		return errors.New("UserId is required")
	}
	return nil
}

var _ OriginatingIdentity = &KubernetesOriginatingIdentity{}

func (i *KubernetesOriginatingIdentity) Platform() string {
	return "kubernetes"
}

func (i *KubernetesOriginatingIdentity) Value() (string, error) {
	if err := i.validate(); err != nil {
		return "", err
	}
	baseProperties := map[string]interface{}{
		"username": i.Username,
		"uid":      i.Uid,
	}
	if len(i.Groups) > 0 {
		baseProperties["groups"] = i.Groups
	}
	return marshalValue(baseProperties, i.Extra)
}

func (i *KubernetesOriginatingIdentity) validate() error {
	if len(i.Username) == 0 {
		return errors.New("Username is required")
	}
	if len(i.Uid) == 0 {
		return errors.New("Uid is required")
	}
	return nil
}

func marshalValue(baseProperties map[string]interface{}, extra map[string]interface{}) (string, error) {
	value := baseProperties
	for extraKey, extraValue := range extra {
		value[extraKey] = extraValue
	}
	valueAsJson, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(valueAsJson), nil
}
