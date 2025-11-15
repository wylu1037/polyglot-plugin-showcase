package adapter

import (
	"fmt"

	"github.com/wylu1037/polyglot-plugin-showcase/plugins/desensitization/impl"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

// DesensitizationAdapter adapts the desensitization implementation to the common plugin interface
type DesensitizationAdapter struct {
	impl *impl.DesensitzerImpl
}

func NewDesensitizationAdapter() *DesensitizationAdapter {
	return &DesensitizationAdapter{
		impl: &impl.DesensitzerImpl{},
	}
}

func (a *DesensitizationAdapter) GetMetadata() (*common.MetadataResponse, error) {
	return &common.MetadataResponse{
		Name:        "desensitization",
		Version:     "1.0.0",
		Description: "Data desensitization plugin",
		Methods: []string{
			"DesensitizeName",
			"DesensitizeTelNo",
			"DesensitizeIDNumber",
			"DesensitizeEmail",
			"DesensitizeBankCard",
			"DesensitizeAddress",
		},
		Capabilities: map[string]string{
			"type": "desensitization",
		},
		ProtocolVersion: common.CurrentProtocolVersion,
	}, nil
}

func (a *DesensitizationAdapter) Execute(method string, params map[string]string) (*common.ExecuteResponse, error) {
	data, ok := params["data"]
	if !ok {
		errMsg := "missing 'data' parameter"
		return &common.ExecuteResponse{
			Success: false,
			Error:   &errMsg,
		}, nil
	}

	var result string
	var err error

	// Route to appropriate method
	switch method {
	case "DesensitizeName":
		result, err = a.impl.DesensitizeName(data)
	case "DesensitizeTelNo":
		result, err = a.impl.DesensitizeTelNo(data)
	case "DesensitizeIDNumber":
		result, err = a.impl.DesensitizeIDNumber(data)
	case "DesensitizeEmail":
		result, err = a.impl.DesensitizeEmail(data)
	case "DesensitizeBankCard":
		result, err = a.impl.DesensitizeBankCard(data)
	case "DesensitizeAddress":
		result, err = a.impl.DesensitizeAddress(data)
	default:
		errMsg := fmt.Sprintf("unknown method: %s", method)
		return &common.ExecuteResponse{
			Success: false,
			Error:   &errMsg,
		}, nil
	}

	if err != nil {
		errMsg := err.Error()
		return &common.ExecuteResponse{
			Success: false,
			Error:   &errMsg,
		}, nil
	}

	return &common.ExecuteResponse{
		Result:  &result,
		Success: true,
	}, nil
}
