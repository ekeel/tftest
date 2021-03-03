package service

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"tftest/model"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/google/uuid"
)

// SecurityGroup holds the data for an AWS security group.
type SecurityGroup struct {
	GroupID                  string
	Name                     string
	Region                   string
	Tags                     map[string]string
	SecurityGroupDescription types.SecurityGroup
}

func (secGroup *SecurityGroup) DescribeByID() (err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(secGroup.Region))
	if err != nil {
		return err
	}

	svc := ec2.NewFromConfig(cfg)

	resp, err := svc.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{secGroup.GroupID},
	})
	if err != nil {
		return err
	}

	if len(resp.SecurityGroups) == 1 {
		secGroup.SecurityGroupDescription = resp.SecurityGroups[0]
	} else {
		return fmt.Errorf("error: a security group with the given id [%v] was not found", secGroup.GroupID)
	}

	secGroup.Tags = getTags(secGroup.SecurityGroupDescription.Tags)

	for key, value := range secGroup.Tags {
		if key == "Name" {
			secGroup.Name = value
		}
	}

	return nil
}

func (secGroup *SecurityGroup) DescribeByName() (err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(secGroup.Region))
	if err != nil {
		return err
	}

	svc := ec2.NewFromConfig(cfg)

	resp, err := svc.DescribeSecurityGroups(context.TODO(), &ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		return err
	}

	if len(resp.SecurityGroups) <= 0 {
		return fmt.Errorf("error: a security group with the given name [%v] was not found", secGroup.Name)
	}

	for _, res := range resp.SecurityGroups {
		tags := getTags(res.Tags)

		if tags["Name"] == secGroup.Name {
			secGroup.Tags = tags
			secGroup.SecurityGroupDescription = res
		}
	}

	return nil
}

func (secGroup *SecurityGroup) ValidateFields(props map[string]string) (validationResults []model.ValidationResult, err error) {
	for key, value := range props {
		validationResult := model.ValidationResult{
			ID:            uuid.NewString(),
			Type:          "security_group:field",
			Name:          key,
			ExpectedValue: value,
			ActualValue:   secGroup.getFieldValue(key),
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

func (secGroup *SecurityGroup) ValidateTags(props map[string]string) (validationResults []model.ValidationResult, err error) {
	for key, value := range props {
		validationResult := model.ValidationResult{
			ID:            uuid.NewString(),
			Type:          "security_group:tag",
			Name:          key,
			ExpectedValue: value,
			ActualValue:   secGroup.Tags[key],
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

func (secGroup *SecurityGroup) getFieldValue(field string) (value string) {
	//switch strings.ToLower(field) {
	//case "ingress":
	//
	//}
	//
	//fmt.Printf("%#v\n", secGroup.SecurityGroupDescription.IpPermissions[0])

	obj := reflect.ValueOf(secGroup.SecurityGroupDescription)
	for _, p := range strings.Split(field, ".") {
		obj = reflect.Indirect(obj).FieldByName(p)
	}

	if obj.Kind() == reflect.Ptr {
		return obj.Elem().String()
	}

	return fmt.Sprintf("%v", obj)
}
