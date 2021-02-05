package service

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"tftest/helpers"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

// EC2Instance holds the data for an AWS EC2 instance.
type EC2Instance struct {
	InstanceID          string
	Name                string
	Region              string
	InstanceDescription types.Instance
}

// DescribeByID gets the instance details using the instance ID.
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

	if len(resp.Reservations) > 0 && len(resp.Reservations[0].Instances) > 0 {
		instance.InstanceDescription = resp.Reservations[0].Instances[0]
	} else {
		return fmt.Errorf("error: no instances were found with the ID [%v]", instance.InstanceID)
	}

	for _, tag := range instance.InstanceDescription.Tags {
		if *tag.Key == "Name" {
			instance.Name = *tag.Value
		}
	}

	return nil
}

// DescribeByName gets the instance details using the instance Name.
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

	if len(resp.Reservations) > 0 && len(resp.Reservations[0].Instances) > 0 {
		instance.InstanceDescription = resp.Reservations[0].Instances[0]
	} else {
		return fmt.Errorf("error: no instances were found with the name [%v]", instance.Name)
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

// ValidateProperties checks if the expected values and acutal values are the same.
func (instance *EC2Instance) ValidateProperties(props map[string]string) (validationResults []helpers.ValidationResult, err error) {
	for key, value := range props {
		validationResult := helpers.ValidationResult{
			ID:            uuid.NewString(),
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

// SSHCommand connects to the EC2 instance using ssh and runs the given command.
func (instance *EC2Instance) SSHCommand(username string, keyFile string, command string) (resp string, err error) {
	authMethod, err := publicKeyFile(keyFile)
	if err != nil {
		return resp, err
	}

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{authMethod},
	}

	connection, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", *instance.InstanceDescription.PrivateIpAddress), sshConfig)
	if err != nil {
		connection.Close()
		return resp, err
	}
	defer connection.Close()

	session, err := connection.NewSession()
	if err != nil {
		session.Close()
		return resp, err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		session.Close()
		return resp, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return resp, err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		return resp, err
	}

	go io.Copy(os.Stderr, stderr)

	err = session.Run(command)
	if err != nil {
		return resp, err
	}

	stdoutBuff := make([]byte, 1000)

	for {
		num, err := stdout.Read(stdoutBuff)
		if err != nil {
			break
		}

		current := string(stdoutBuff[:num])

		if checkPrompt(current) {
			resp += current
			break
		}

		resp += current
	}

	return resp, nil
}

// getFieldValue using reflection to get the value of the field specified.
func (instance *EC2Instance) getFieldValue(field string) (value string) {
	r := reflect.ValueOf(instance.InstanceDescription)
	f := reflect.Indirect(r).FieldByName(field)

	return f.String()
}

// checkPrompt checks if the current SSH command output is a prompt.
func checkPrompt(s string) bool {
	prompt := regexp.MustCompile(".*@?.*(#|>) $")
	m := prompt.FindStringSubmatch(s)
	return m != nil
}

// publicKeyFile reads the private key and returns an SSH auth method based on it.
func publicKeyFile(file string) (authMethod ssh.AuthMethod, err error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}

	authMethod = ssh.PublicKeys(key)

	return authMethod, nil
}
