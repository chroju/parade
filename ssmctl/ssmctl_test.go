package ssmctl

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

const (
	dummyEncryptedValue = "(encrypted)"
)

type mockSSMClient struct {
	ssmiface.SSMAPI
}

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

func NewMock() *SSMManager {
	svc := &mockSSMClient{}

	return &SSMManager{
		svc: svc,
	}
}

func (m *mockSSMClient) GetParameter(i *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	for _, v := range mockParameters {
		if *i.Name == *v.Name {
			if *v.Type == "SecureString" && *i.WithDecryption == false {
				v.Value = aws.String(dummyEncryptedValue)
			}
			return &ssm.GetParameterOutput{Parameter: v}, nil
		}
	}

	return nil, errors.New(ssm.ErrCodeParameterNotFound)
}

func TestGetParameter(t *testing.T) {
	m := NewMock()
	cases := []struct {
		query           string
		withDescryption bool
		expected        string
	}{
		{
			"/service1/dev/key1",
			false,
			"dev_value1",
		},
		{
			"/service1/dev/key2",
			true,
			"dev_value2",
		},
		{
			"/service1/prod/key2",
			false,
			dummyEncryptedValue,
		},
	}

	for _, c := range cases {
		result, err := m.GetParameter(c.query, c.withDescryption)
		if err != nil {
			t.Fatalf("Failed: error = %s", err)
		}

		if result.Value != c.expected {
			t.Errorf("want: %s\nget : %s", c.expected, result)
		}
	}

}
