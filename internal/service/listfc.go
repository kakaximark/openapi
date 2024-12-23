package service

import (
	"fmt"

	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// FcInfo 函数信息
type FcInfo struct {
	FunctionName string `json:"functionName"`
	Description  string `json:"description"`
}

// FcList 函数列表响应
type FcList struct {
	Fcs []FcInfo `json:"functions"`
}

// ListFc 获取函数列表信息
func ListFc(env, countryCode, serviceName string) (*FcList, error) {
	client, err := GetFCClient(env, countryCode)
	if err != nil {
		return nil, err
	}

	// 从数据库获取配置
	config, err := LoadClientConfig(env, countryCode)
	if err != nil {
		return nil, err
	}

	listFunctionsHeaders := &fc_open20210406.ListFunctionsHeaders{
		XFcAccountId: tea.String(config.AliyunConfig.MainAccountID),
	}
	listFunctionsRequest := &fc_open20210406.ListFunctionsRequest{}
	runtime := &util.RuntimeOptions{}

	resp, err := client.ListFunctionsWithOptions(tea.String(serviceName), listFunctionsRequest, listFunctionsHeaders, runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to get functions: %v", err)
	}

	fcList := &FcList{
		Fcs: make([]FcInfo, 0, len(resp.Body.Functions)),
	}

	for _, fc := range resp.Body.Functions {
		fcInfo := FcInfo{
			FunctionName: tea.StringValue(fc.FunctionName),
			Description:  tea.StringValue(fc.Description),
		}
		fcList.Fcs = append(fcList.Fcs, fcInfo)
	}

	return fcList, nil
}
