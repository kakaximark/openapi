package service

import (
	"fmt"

	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// AliasInfo 函数别名信息
type AliasInfo struct {
	AliasName   string `json:"aliasName"`
	VersionId   string `json:"versionId"`
	Description string `json:"description"`
}

// UpdateAlias 更新函数别名
func UpdateAlias(env, countryCode, serviceName, aliasName, versionId string) (*AliasInfo, error) {
	client, err := GetFCClient(env, countryCode)
	if err != nil {
		return nil, err
	}

	// 从数据库获取配置
	config, err := LoadClientConfig(env, countryCode)
	if err != nil {
		return nil, err
	}

	updateAliasHeaders := &fc_open20210406.UpdateAliasHeaders{
		XFcAccountId: tea.String(config.AliyunConfig.MainAccountID),
	}
	updateAliasRequest := &fc_open20210406.UpdateAliasRequest{
		VersionId: tea.String(versionId),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.UpdateAliasWithOptions(tea.String(serviceName), tea.String(aliasName), updateAliasRequest, updateAliasHeaders, runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to update alias: %v", err)
	}

	settings := &AliasInfo{
		AliasName:   tea.StringValue(resp.Body.AliasName),
		VersionId:   tea.StringValue(resp.Body.VersionId),
		Description: tea.StringValue(resp.Body.Description),
	}

	return settings, nil
}
