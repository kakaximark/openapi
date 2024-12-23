package service

import (
	"fmt"

	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// PublicServiceInfo 发布服务信息
type PublicServiceInfo struct {
	VersionID    string `json:"versionId"`
	Description  string `json:"description"`
	CreatedTime  string `json:"createdTime"`
	LastModified string `json:"lastModifiedTime"`
}

// PublicService 发布服务
func PublicService(env, region, serviceName, description string) (*PublicServiceInfo, error) {
	client, err := GetFCClient(env, region)
	if err != nil {
		return nil, err
	}

	// 从数据库获取配置
	config, err := LoadClientConfig(env, region)
	if err != nil {
		return nil, err
	}

	publishServiceVersionHeaders := &fc_open20210406.PublishServiceVersionHeaders{
		XFcAccountId: tea.String(config.AliyunConfig.MainAccountID),
	}
	publishServiceVersionRequest := &fc_open20210406.PublishServiceVersionRequest{
		Description: tea.String(description),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.PublishServiceVersionWithOptions(tea.String(serviceName), publishServiceVersionRequest, publishServiceVersionHeaders, runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to publish service: %v", err)
	}

	settings := &PublicServiceInfo{
		VersionID:    tea.StringValue(resp.Body.VersionId),
		Description:  tea.StringValue(resp.Body.Description),
		CreatedTime:  tea.StringValue(resp.Body.CreatedTime),
		LastModified: tea.StringValue(resp.Body.LastModifiedTime),
	}

	return settings, nil
}
