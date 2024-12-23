package service

import (
	"fmt"

	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	ServiceName string `json:"serviceName"`
	Description string `json:"description"`
}

// ServiceList 服务列表响应
type ServiceList struct {
	Services []ServiceInfo `json:"services"`
}

// ListService 获取服务列表信息
func ListService(env, countryCode string) (*ServiceList, error) {
	client, err := GetFCClient(env, countryCode)
	if err != nil {
		return nil, err
	}

	// 加载配置
	config, err := LoadClientConfig(env, countryCode)
	if err != nil {
		return nil, err
	}

	listServicesHeaders := &fc_open20210406.ListServicesHeaders{
		XFcAccountId: tea.String(config.AliyunConfig.MainAccountID),
	}
	listServicesRequest := &fc_open20210406.ListServicesRequest{}
	runtime := &util.RuntimeOptions{}

	resp, err := client.ListServicesWithOptions(listServicesRequest, listServicesHeaders, runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}

	serviceList := &ServiceList{
		Services: make([]ServiceInfo, 0, len(resp.Body.Services)),
	}

	for _, service := range resp.Body.Services {
		serviceInfo := ServiceInfo{
			ServiceName: tea.StringValue(service.ServiceName),
			Description: tea.StringValue(service.Description),
		}
		serviceList.Services = append(serviceList.Services, serviceInfo)
	}

	return serviceList, nil
}
