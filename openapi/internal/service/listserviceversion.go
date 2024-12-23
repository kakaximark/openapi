package service

import (
	"fmt"

	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// ServiceVersionInfo 服务版本信息
type ServiceVersionInfo struct {
	VersionID   string `json:"versionId"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdTime"`
	UpdatedAt   string `json:"lastModifiedTime"`
}

// ServiceList 服务列表响应
type ServiceVersionList struct {
	Versions []ServiceVersionInfo `json:"versions"`
}

// ListServiceVersion 获取服务版本列表信息
func ListServiceVersion(env, countryCode, serviceName string) (*ServiceVersionList, error) {
	client, err := GetFCClient(env, countryCode)
	if err != nil {
		return nil, err
	}

	// 从数据库获取配置
	config, err := LoadClientConfig(env, countryCode)
	if err != nil {
		return nil, err
	}

	listServiceVersionsHeaders := &fc_open20210406.ListServiceVersionsHeaders{
		XFcAccountId: tea.String(config.AliyunConfig.MainAccountID),
	}
	listServiceVersionsRequest := &fc_open20210406.ListServiceVersionsRequest{}
	runtime := &util.RuntimeOptions{}

	resp, err := client.ListServiceVersionsWithOptions(tea.String(serviceName), listServiceVersionsRequest, listServiceVersionsHeaders, runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}

	serviceList := &ServiceVersionList{
		Versions: make([]ServiceVersionInfo, 0, len(resp.Body.Versions)),
	}

	// 遍历响应中的服务列表
	for _, service := range resp.Body.Versions {
		serviceInfo := ServiceVersionInfo{
			VersionID:   tea.StringValue(service.VersionId),
			Description: tea.StringValue(service.Description),
			CreatedAt:   tea.StringValue(service.CreatedTime),
			UpdatedAt:   tea.StringValue(service.LastModifiedTime),
		}
		serviceList.Versions = append(serviceList.Versions, serviceInfo)
	}

	return serviceList, nil
}
