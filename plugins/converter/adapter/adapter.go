package adapter

import (
	"fmt"

	"github.com/wylu1037/polyglot-plugin-showcase/plugins/converter/impl"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

// ConverterAdapter adapts the converter implementation to the common plugin interface
type ConverterAdapter struct {
	impl *impl.ConverterImpl
}

func NewConverterAdapter() *ConverterAdapter {
	return &ConverterAdapter{
		impl: &impl.ConverterImpl{},
	}
}

func (a *ConverterAdapter) GetMetadata() (*common.MetadataResponse, error) {
	return &common.MetadataResponse{
		Name:        "converter",
		Version:     "1.0.0",
		Description: "Data format converter plugin - converts JSON to CSV, TXT, HTML",
		Methods: []string{
			"ConvertToCSV",
			"ConvertToTXT",
			"ConvertToHTML",
		},
		Capabilities: map[string]string{
			"type":           "converter",
			"input_format":   "json",
			"output_formats": "csv,txt,html",
		},
		ProtocolVersion: common.CurrentProtocolVersion,
	}, nil
}

func (a *ConverterAdapter) Execute(method string, params map[string]string) (*common.ExecuteResponse, error) {
	data, ok := params["data"]
	if !ok {
		errMsg := "missing 'data' parameter"
		return &common.ExecuteResponse{
			Success: false,
			Error:   &errMsg,
		}, nil
	}

	// Extract options from params (all params except 'data' are considered options)
	options := make(map[string]string)
	for key, value := range params {
		if key != "data" {
			options[key] = value
		}
	}

	var result string
	var err error

	// Route to appropriate method
	switch method {
	case "ConvertToCSV":
		result, err = a.impl.ConvertToCSV(data, options)
	case "ConvertToTXT":
		result, err = a.impl.ConvertToTXT(data, options)
	case "ConvertToHTML":
		result, err = a.impl.ConvertToHTML(data, options)
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
