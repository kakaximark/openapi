package service

import (
	"fmt"

	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// AliasInfo 别名信息
type AliasNewInfo struct {
	AliasName        string `json:"aliasName"`
	VersionID        string `json:"versionId"`
	Description      string `json:"description"`
	CreatedTime      string `json:"createdTime"`
	LastModifiedTime string `json:"lastModifiedTime"`
}

// AliasList 别名列表响应
type AliasList struct {
	Aliases []AliasNewInfo `json:"aliases"`
}

// ListAlias 获取别名列表信息
func ListAlias(env, countryCode, serviceName string) (*AliasList, error) {
	client, err := GetFCClient(env, countryCode)
	if err != nil {
		return nil, err
	}

	// 从数据库获取配置
	config, err := LoadClientConfig(env, countryCode)
	if err != nil {
		return nil, err
	}

	listAliasesHeaders := &fc_open20210406.ListAliasesHeaders{
		XFcAccountId: tea.String(config.AliyunConfig.MainAccountID),
	}
	listAliasesRequest := &fc_open20210406.ListAliasesRequest{}
	runtime := &util.RuntimeOptions{}

	resp, err := client.ListAliasesWithOptions(tea.String(serviceName), listAliasesRequest, listAliasesHeaders, runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to get aliases: %v", err)
	}

	aliasList := &AliasList{
		Aliases: make([]AliasNewInfo, 0, len(resp.Body.Aliases)),
	}

	for _, alias := range resp.Body.Aliases {
		aliasInfo := AliasNewInfo{
			AliasName:        tea.StringValue(alias.AliasName),
			VersionID:        tea.StringValue(alias.VersionId),
			Description:      tea.StringValue(alias.Description),
			CreatedTime:      tea.StringValue(alias.CreatedTime),
			LastModifiedTime: tea.StringValue(alias.LastModifiedTime),
		}
		aliasList.Aliases = append(aliasList.Aliases, aliasInfo)
	}

	return aliasList, nil
}
