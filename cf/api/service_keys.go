package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cloudfoundry/cli/cf/api/resources"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/net"
)

type ServiceKeyRepository interface {
	CreateServiceKey(instanceId string, keyName string) (apiErr error)
	ListServiceKeys(instanceId string) (serviceKeys []models.ServiceKey, apiErr error)
	GetServiceKey(instanceId string, keyName string) (serviceKey models.ServiceKey, apiErr error)
}

type CloudControllerServiceKeyRepository struct {
	config  core_config.Reader
	gateway net.Gateway
}

func NewCloudControllerServiceKeyRepository(config core_config.Reader, gateway net.Gateway) (repo CloudControllerServiceKeyRepository) {
	return CloudControllerServiceKeyRepository{
		config:  config,
		gateway: gateway,
	}
}

func (c CloudControllerServiceKeyRepository) CreateServiceKey(instanceId string, keyName string) (apiErr error) {
	path := "/v2/service_keys"
	data := fmt.Sprintf(`{"service_instance_guid":"%s","name":"%s"}`, instanceId, keyName)

	err := c.gateway.CreateResource(c.config.ApiEndpoint(), path, strings.NewReader(data))

	if httpErr, ok := err.(errors.HttpError); ok && httpErr.ErrorCode() == errors.SERVICE_KEY_NAME_TAKEN {
		return errors.NewModelAlreadyExistsError("Service key", keyName)
	}

	return nil
}

func (c CloudControllerServiceKeyRepository) ListServiceKeys(instanceId string) (serviceKeys []models.ServiceKey, apiErr error) {
	path := fmt.Sprintf("/v2/service_keys?q=service_instance_guid:%s", instanceId)

	serviceKeys = []models.ServiceKey{}
	apiErr = c.gateway.ListPaginatedResources(
		c.config.ApiEndpoint(),
		path,
		resources.ServiceKeyResource{},
		func(resource interface{}) bool {
			serviceKey := resource.(resources.ServiceKeyResource).ToModel()
			serviceKeys = append(serviceKeys, serviceKey)
			return true
		})

	if apiErr != nil {
		return []models.ServiceKey{}, apiErr
	}

	return serviceKeys, nil
}

func (c CloudControllerServiceKeyRepository) GetServiceKey(instanceId string, keyName string) (serviceKey models.ServiceKey, apiErr error) {
	url := fmt.Sprintf("%s/v2/service_keys?q=service_instance_guid:%s;%s", c.config.ApiEndpoint(), instanceId, url.QueryEscape("name:"+keyName))

	serviceKeyResource := new(resources.ServiceKeyResource)
	apiErr = c.gateway.GetResource(url, serviceKeyResource)

	if apiErr != nil {
		return models.ServiceKey{}, apiErr
	}

	return serviceKeyResource.ToModel(), nil
}
