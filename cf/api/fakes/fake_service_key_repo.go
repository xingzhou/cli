package fakes

import (
	"github.com/cloudfoundry/cli/cf/models"
)

type FakeServiceKeyRepo struct {
	CreateServiceKeyMethod CreateServiceKeyType
	ListServiceKeysMethod  ListServiceKeysType
	GetServiceKeyMethod    GetServiceKeyType
}

type CreateServiceKeyType struct {
	InstanceId string
	KeyName    string

	Error error
}

type ListServiceKeysType struct {
	InstanceId string

	ServiceKeys []models.ServiceKey
	Error       error
}

type GetServiceKeyType struct {
	InstanceId string
	KeyName    string

	ServiceKey models.ServiceKey
	Error      error
}

func NewFakeServiceKeyRepo() *FakeServiceKeyRepo {
	return &FakeServiceKeyRepo{
		CreateServiceKeyMethod: CreateServiceKeyType{},
		ListServiceKeysMethod:  ListServiceKeysType{},
		GetServiceKeyMethod:    GetServiceKeyType{},
	}
}

func (f *FakeServiceKeyRepo) CreateServiceKey(instanceId string, serviceKeyName string) (apiErr error) {
	f.CreateServiceKeyMethod.InstanceId = instanceId
	f.CreateServiceKeyMethod.KeyName = serviceKeyName

	return f.CreateServiceKeyMethod.Error
}

func (f *FakeServiceKeyRepo) ListServiceKeys(instanceId string) (serviceKeys []models.ServiceKey, apiErr error) {
	f.ListServiceKeysMethod.InstanceId = instanceId

	return f.ListServiceKeysMethod.ServiceKeys, f.ListServiceKeysMethod.Error
}

func (f *FakeServiceKeyRepo) GetServiceKey(instanceId string, serviceKeyName string) (serviceKey models.ServiceKey, apiErr error) {
	f.GetServiceKeyMethod.InstanceId = instanceId

	return f.GetServiceKeyMethod.ServiceKey, f.GetServiceKeyMethod.Error
}
