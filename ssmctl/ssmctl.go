package ssmctl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

const (
	// DescribeOptionEquals is a describe parameters API option means equals
	DescribeOptionEquals = "Equals"
	// DescribeOptionContains is a describe parameters API option means contains
	DescribeOptionContains = "Contains"
	// DescribeOptionBeginsWith is a describe parameters API option means beginsWith
	DescribeOptionBeginsWith = "BeginsWith"
)

// SSMManager is the wrapper of SSM API.
type SSMManager struct {
	svc ssmiface.SSMAPI
}

// Parameter is the struct of SSM parameter.
type Parameter struct {
	Name  string
	Value string
	Type  string
}

// New returns a new SSMManager.
func New(profile, region string) (*SSMManager, error) {
	config := &aws.Config{}
	if profile != "" {
		config.Credentials = credentials.NewSharedCredentials("", profile)
	}

	if region != "" {
		config.Region = aws.String(region)
	}

	sess := session.Must(session.NewSession(config))
	svc := ssm.New(sess)
	_, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}

	return &SSMManager{
		svc: svc,
	}, nil
}

// GetParameter gets a SSM parameter.
func (s *SSMManager) GetParameter(query string, withDecryption bool) (*Parameter, error) {
	params := &ssm.GetParameterInput{
		Name:           aws.String(query),
		WithDecryption: aws.Bool(withDecryption),
	}

	resp, err := s.svc.GetParameter(params)
	if err != nil {
		return nil, err
	}
	value := *resp.Parameter.Value
	if !withDecryption && *resp.Parameter.Type == "SecureString" {
		value = "(encrypted)"
	}

	return &Parameter{
		Name:  *resp.Parameter.Name,
		Value: value,
		Type:  *resp.Parameter.Type,
	}, nil
}

// DescribeParameters describes SSM parameters.
func (s *SSMManager) DescribeParameters(query string, option string) ([]*Parameter, error) {
	params := &ssm.DescribeParametersInput{
		MaxResults: aws.Int64(50),
	}
	if query != "" {
		filter := &ssm.ParameterStringFilter{
			Key:    aws.String("Name"),
			Option: aws.String(option),
			Values: aws.StringSlice([]string{query}),
		}
		params.ParameterFilters = []*ssm.ParameterStringFilter{filter}
	}

	var metaDatas []*ssm.ParameterMetadata
	for {
		resp, err := s.svc.DescribeParameters(params)
		if err != nil {
			return nil, err
		}
		metaDatas = append(metaDatas, resp.Parameters...)
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	result := make([]*Parameter, len(metaDatas))
	for i, v := range metaDatas {
		result[i] = &Parameter{
			Name: *v.Name,
			Type: *v.Type,
		}
	}

	return result, nil
}

// PutParameter puts a SSM parameter.
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

// DeleteParameter deletes a SSM parameter.
func (s *SSMManager) DeleteParameter(key string) error {
	param := &ssm.DeleteParameterInput{
		Name: aws.String(key),
	}

	if _, err := s.svc.DeleteParameter(param); err != nil {
		return err
	}

	return nil
}
