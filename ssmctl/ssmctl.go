package ssmctl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMManager struct {
	svc *ssm.SSM
}

func New() (*SSMManager, error) {
	sess := session.Must(session.NewSession())
	svc := ssm.New(sess)

	return &SSMManager{
		svc: svc,
	}, nil
}

func (s *SSMManager) GetParameters(query string) ([]*ssm.Parameter, error) {
	params := &ssm.GetParametersInput {
		Names:  []*string{aws.String(query)},
		WithDecryption: aws.Bool(false),
	}

	resp, err := s.svc.GetParameters(params)
	if err != nil {
		return nil, err
	}

	return resp.Parameters, nil
}

func (s *SSMManager) DescribeParameters() ([]*ssm.ParameterMetadata, error) {
	params := &ssm.DescribeParametersInput {
		MaxResults: aws.Int64(50),
	}

	var result []*ssm.ParameterMetadata
	for {
		resp, err := s.svc.DescribeParameters(params)
		if err != nil {
			return nil, err
		}
		result = append(result, resp.Parameters...)
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return result, nil
}
