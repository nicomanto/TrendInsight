package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type ParameterStore struct {
	client *ssm.Client
}

func NewParameterStore() (*ParameterStore, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := ssm.NewFromConfig(cfg)
	return &ParameterStore{
		client: client,
	}, nil
}

func (ps *ParameterStore) GetParametersByPath(path string, decrypt bool, maxResult int32) (map[string]string, error) {
	params, err := ps.client.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		WithDecryption: &decrypt,
		Path:           &path,
		MaxResults:     &maxResult,
	})
	if err != nil {
		return nil, err
	}
	returnedParams := map[string]string{}
	for _, p := range params.Parameters {
		if p.Name == nil || p.Value == nil {
			continue
		}
		n := strings.Split(*p.Name, "/")
		returnedParams[n[len(n)-1]] = *p.Value
	}
	return returnedParams, nil
}

func (ps *ParameterStore) GetParameterValue(name string, withDecryption bool) (string, error) {
	results, err := ps.client.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		return "", err
	}
	if results.Parameter.Value == nil {
		return "", fmt.Errorf("failed to find parameter %s", name)
	}
	return *results.Parameter.Value, nil
}
