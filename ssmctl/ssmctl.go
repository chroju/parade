package ssmctl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type SSMManager struct {
	svc ssmiface.SSMAPI
}

func New() (*SSMManager, error) {
	sess := session.Must(session.NewSession())
	svc := ssm.New(sess)

	return &SSMManager{
		svc: svc,
	}, nil
}

func (s *SSMManager) GetParameter(query string, withDecryption bool) (*ssm.Parameter, error) {
	params := &ssm.GetParameterInput{
		Name:           aws.String(query),
		WithDecryption: aws.Bool(withDecryption),
	}

	resp, err := s.svc.GetParameter(params)
	if err != nil {
		return nil, err
	}

	return resp.Parameter, nil
}

func (s *SSMManager) DescribeParameters() ([]*ssm.ParameterMetadata, error) {
	params := &ssm.DescribeParametersInput{
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

func (s *SSMManager) PutParameter(key string, value string, isEncryption bool, isForce bool) error {
	var paramType string
	if isEncryption {
		paramType = "SecureString"
	} else {
		paramType = "String"
	}

	param := &ssm.PutParameterInput{
		Name:      aws.String(key),
		Value:     aws.String(value),
		Type:      aws.String(paramType),
		Overwrite: aws.Bool(isForce),
	}

	if _, err := s.svc.PutParameter(param); err != nil {
		return err
	}

	return nil
}

func (s *SSMManager) DeleteParameter(key string) error {
	param := &ssm.DeleteParameterInput{
		Name: aws.String(key),
	}

	if _, err := s.svc.DeleteParameter(param); err != nil {
		return err
	}

	return nil
}
