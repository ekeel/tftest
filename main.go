package main

import (
	"flag"
	"os"
	"strings"
	"tftest/model"
	"tftest/service"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "hcl", "", "")

	flag.Parse()

	// config, err := model.Parse("D:\\projects\\tftest\\example\\ec2.hcl")
	config, err := model.Parse(configFile)
	if err != nil {
		panic(err)
	}

	validate(config)
}

func validate(config model.Config) (err error) {
	for _, test := range config.Tests {
		switch strings.ToLower(test.Type) {
		case "ec2":
			var instance service.EC2Instance
			switch strings.ToLower(test.QueryBy) {
			case "name":
				instance = service.EC2Instance{Name: test.InstanceName}
				err = instance.DescribeByName()
				if err != nil {
					panic(err)
				}
			case "id":
				instance = service.EC2Instance{InstanceID: test.InstanceID}
				err = instance.DescribeByID()
				if err != nil {
					panic(err)
				}
			}

			fieldValidationResults, err := instance.ValidateFields(test.Fields)
			if err != nil {
				panic(err)
			}

			tagValidationResults, err := instance.ValidateTags(test.Tags)
			if err != nil {
				panic(err)
			}

			test.ValidationResults = append(fieldValidationResults, tagValidationResults...)
		case "security_group":
			var secGroup service.SecurityGroup
			switch strings.ToLower(test.QueryBy) {
			case "name":
				secGroup = service.SecurityGroup{Name: test.InstanceName}
				err = secGroup.DescribeByName()
				if err != nil {
					panic(err)
				}
			case "id":
				secGroup = service.SecurityGroup{GroupID: test.InstanceID}
				err = secGroup.DescribeByID()
				if err != nil {
					panic(err)
				}
			}

			fieldValidationResults, err := secGroup.ValidateFields(test.Fields)
			if err != nil {
				panic(err)
			}

			tagValidationResults, err := secGroup.ValidateTags(test.Tags)
			if err != nil {
				panic(err)
			}

			test.ValidationResults = append(fieldValidationResults, tagValidationResults...)
		}
	}

	printResults(config)

	return nil
}

func printResults(config model.Config) {
	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	tbl.AppendHeader(table.Row{"Test Name", "Type", "Field Name", "Passed", "Expected Value", "Actual Value"})

	for _, test := range config.Tests {
		for _, res := range test.ValidationResults {
			tbl.AppendRows([]table.Row{
				{test.Name, res.Type, res.Name, res.IsMatch, res.ExpectedValue, res.ActualValue},
			})
		}
		tbl.AppendSeparator()
	}

	tbl.SetRowPainter(table.RowPainter(func(row table.Row) text.Colors {
		switch row[3] {
		case true:
			return text.Colors{text.BgGreen, text.FgBlack}
		case false:
			return text.Colors{text.BgRed, text.FgBlack}
		default:
			return text.Colors{text.BgYellow, text.FgBlack}
		}
	}))

	tbl.Render()
}
