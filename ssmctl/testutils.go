package ssmctl

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

const (
	dummyEncryptedValue = "(encrypted)"
)

var mockParameters = []*ssm.Parameter{
	{
		Name:  aws.String("/service1/dev/key1"),
		Value: aws.String("dev_value1"),
		Type:  aws.String("String"),
	},
	{
		Name:  aws.String("/service1/dev/key2"),
		Value: aws.String("dev_value2"),
		Type:  aws.String("SecureString"),
	},
	{
		Name:  aws.String("/service1/prod/key1"),
		Value: aws.String("prod_value1"),
		Type:  aws.String("String"),
	},
	{
		Name:  aws.String("/service1/prod/key2"),
		Value: aws.String("prod_value2"),
		Type:  aws.String("SecureString"),
	},
}

type mockSSMClient struct {
	ssmiface.SSMAPI
}

func NewMockSSMManager() SSMManager {
	return &ssmManager{
		svc: &mockSSMClient{},
	}
}

func (m *mockSSMClient) GetParameter(i *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	for _, v := range mockParameters {
		if *i.Name == *v.Name {
			value := v.Value
			if *v.Type == "SecureString" && *i.WithDecryption == false {
				value = aws.String(dummyEncryptedValue)
			}
			return &ssm.GetParameterOutput{Parameter: &ssm.Parameter{
				Name:  v.Name,
				Type:  v.Type,
				Value: value,
			}}, nil
		}
	}

	return nil, errors.New(ssm.ErrCodeParameterNotFound)
}

func (m *mockSSMClient) DescribeParameters(i *ssm.DescribeParametersInput) (*ssm.DescribeParametersOutput, error) {
	var result []*ssm.ParameterMetadata
	var filter *ssm.ParameterStringFilter

	if len(i.ParameterFilters) != 0 {
		filter = i.ParameterFilters[0]
	}

	for _, v := range mockParameters {
		if filter == nil {
			result = append(result, &ssm.ParameterMetadata{
				Name: v.Name,
				Type: v.Type,
			})
			continue
		}

		switch *filter.Option {
		case DescribeOptionBeginsWith:
			if !strings.HasPrefix(*v.Name, *filter.Values[0]) && !strings.HasPrefix(*v.Name, "/"+*filter.Values[0]) {
				continue
			}
		case DescribeOptionEquals:
			if *v.Name != *filter.Values[0] {
				continue
			}
		case DescribeOptionContains:
			if !strings.Contains(*v.Name, *filter.Values[0]) {
				continue
			}
		}

		result = append(result, &ssm.ParameterMetadata{
			Name: v.Name,
			Type: v.Type,
		})
	}

	if len(result) == 0 {
		return nil, errors.New(ssm.ErrCodeParameterNotFound)
	}
	return &ssm.DescribeParametersOutput{
		Parameters: result,
	}, nil

}

func (m *mockSSMClient) PutParameter(*ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	return nil, nil
}

func (m *mockSSMClient) DeleteParameter(*ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error) {
	return nil, nil
}
