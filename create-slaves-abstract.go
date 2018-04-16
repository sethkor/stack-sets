package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	pProfiles  = kingpin.Arg("profiles", "The AWS profiles to deploy the slave role too").Required().Strings()
	pDelete    = kingpin.Flag("delete", "Delete the slaves").Short('d').Bool()
	pMaster    = kingpin.Flag("master", "The master account").Short('m').Required().String()
	pRegion    = kingpin.Flag("region", "Region to deploy the slave roles too").Short('r').Required().String()
	pTemplate  = kingpin.Flag("template", "The stack template").Short('t').Required().ExistingFile()
	pStackName = kingpin.Flag("stackname", "The stack name to use in cloud formation").Required().Short('s').String()
)

// I didn't know how to use this
// func errHandler(err error, desc string) {
// 	if err != nil {
// 		fmt.Printf("I got an error %s: %v", desc, err)
// 		if aerr, ok := err.(awserr.Error); ok {
// 			switch aerr.Code() {
// 			default:
// 				fmt.Println(aerr.Error())
// 			}
// 		} else {
// 			// Print the error, cast err to awserr.Error to get the Code and
// 			// Message from an error.
// 			fmt.Println(err.Error())
// 		}
// 		return
// 	}
//
// }

func deleteProfiles(profiles []string) error {

	input := &cloudformation.DeleteStackInput{
		StackName: aws.String(*pStackName),
	}

	return doForEachProfile(profiles, "deleting a profile", func(svc *cloudformation.CloudFormation) (interface{}, error) {
		return svc.DeleteStack(input)
	})
}

func createProfiles(profiles []string, templatePath string) error {

	var awsParameters []*cloudformation.Parameter
	awsParameters = append(awsParameters, &cloudformation.Parameter{
		ParameterKey:   aws.String("MasterAccountId"),
		ParameterValue: aws.String(*pMaster),
	})

	awsParameters = append(awsParameters, &cloudformation.Parameter{
		ParameterKey:   aws.String("MasterRegion"),
		ParameterValue: aws.String(*pRegion),
	})

	var templateBody string
	if data, err := ioutil.ReadFile(templatePath); err != nil {
		if err != nil {
			return err
		}
	} else {
		templateBody = string(data)
	}

	//and run it for each profile
	return doForEachProfile(profiles, "creating a profile", func(svc *cloudformation.CloudFormation) (interface{}, error) {
		input := &cloudformation.CreateStackInput{
			Capabilities: []*string{aws.String("CAPABILITY_NAMED_IAM")},
			StackName:    aws.String(*pStackName),
			TemplateBody: &templateBody,
			Parameters:   awsParameters,
		}
		return svc.CreateStack(input)
	})
}

func doForEachProfile(profiles []string, operationDesc string, op func(*cloudformation.CloudFormation) (interface{}, error)) error {
	for _, profile := range profiles {
		// Specify profile for config and region for requests
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           profile,
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := cloudformation.New(sess)
		result, err := op(svc)
		if err != nil {
			// This is a choice, you may not want to bail out at the first error
			return err
		}
		fmt.Printf("%s: %v", operationDesc, result)
	}
	return nil
}

func main() {
	//Parse Command Line
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	profiles := *(pProfiles)

	if *pDelete {
		err := deleteProfiles(profiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An Error occured deleting profiles: %v", err)
			os.Exit(1)
		}
	} else {
		err := createProfiles(profiles, *pTemplate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "An Error occured creating profiles: %v", err)
			os.Exit(1)
		}
	}
}
