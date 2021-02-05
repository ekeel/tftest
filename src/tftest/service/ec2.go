package service

import (
	"context"
	"reflect"
	"tftest/helpers"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Instance struct {
	InstanceID          string
	Name                string
	Region              string
	InstanceDescription types.Instance
}

func (instance *EC2Instance) DescribeByID() (err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(instance.Region))
	if err != nil {
		return err
	}

	svc := ec2.NewFromConfig(cfg)

	resp, err := svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: []string{instance.InstanceID},
	})
	if err != nil {
		return err
	}

	instance.InstanceDescription = resp.Reservations[0].Instances[0]

	for _, tag := range instance.InstanceDescription.Tags {
		if *tag.Key == "Name" {
			instance.Name = *tag.Value
		}
	}

	return nil
}

func (instance *EC2Instance) DescribeByName() (err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(instance.Region))
	if err != nil {
		return err
	}

	svc := ec2.NewFromConfig(cfg)

	resp, err := svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return err
	}

	for _, res := range resp.Reservations {
		for _, inst := range res.Instances {
			if *inst.Tags[0].Value == instance.Name {
				instance.InstanceDescription = inst
			}
		}
	}

	return nil
}

func (instance *EC2Instance) ValidateProperties(props map[string]string) (validationResults []helpers.ValidationResult, err error) {
	for key, value := range props {
		validationResult := helpers.ValidationResult{
			Name:          key,
			ExpectedValue: value,
			ActualValue:   instance.getFieldValue(key),
		}

		if validationResult.ActualValue == validationResult.ExpectedValue {
			validationResult.IsMatch = true
		} else {
			validationResult.IsMatch = false
		}

		validationResults = append(validationResults, validationResult)
	}

	return validationResults, nil
}

func (instance *EC2Instance) getFieldValue(field string) (value string) {
	r := reflect.ValueOf(instance.InstanceDescription)
	f := reflect.Indirect(r).FieldByName(field)

	return f.String()
}
