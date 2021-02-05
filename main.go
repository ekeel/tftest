package main

import (
	"fmt"
	"tftest/service"
)

func main() {
	ec2Instance := service.EC2Instance{
		// InstanceID: "i-0585edc5229c84f02",
		Name:   "fxs-bam-np-instance-a",
		Region: "us-west-2",
	}

	err := ec2Instance.DescribeByName()
	// err := ec2Instance.DescribeByID()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", *ec2Instance.InstanceDescription.InstanceId)
	fmt.Printf("%#v\n", ec2Instance.Name)
}
