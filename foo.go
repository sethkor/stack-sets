package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws/session"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

package main

import (
"gopkg.in/alecthomas/kingpin.v2"
"fmt"
"github.com/aws/aws-sdk-go/aws/session"
"github.com/aws/aws-sdk-go/aws/awserr"
"github.com/aws/aws-sdk-go/service/cloudformation"
"github.com/aws/aws-sdk-go/aws"
"reflect"
"github.com/aws/aws-sdk-go/service/connect"
)

var (

	profiles		= kingpin.Arg("profiles", "The AWS profiles to deploy the slave role too").Required().Strings()
	delete		= kingpin.Flag("delete", "Delete the slaves").Short('d').Bool()
	master		= kingpin.Flag("master", "The master account").Short('m').Required().Strings()
	region    	= kingpin.Flag("region", "Region to deploy the slave roles too").Short('r').Required().String()
	template    = kingpin.Flag("template", "The stack template").Short('r').ExistingFile()
	stackName 	= kingpin.Flag("stackname",  "The stack name to use in cloud formation").Required().Short('s').String()


)

func deleteStack() {
	input
	result, err := svc.DeleteStack(input)
}

func createStack() {

}

func main() {
	//Parse Command Line
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	var input

	var svc cloudformation.CloudFormation
	var theFunc uintptr

	//if(*delete) {
	//	input = &cloudformation.DeleteStackInput{
	//		StackName: aws.String(*stackName),
	//	}
	//	theFunc = reflect.ValueOf(svc.DeleteStack).Pointer()
	//} else {
	//	input = &cloudformation.CreateStackInput{
	//		Capabilities: []*string{aws.String("CAPABILITY_NAMED_IAM"),},
	//		StackName: aws.String(*stackName),
	//		TemplateBody: template,
	//	}
	//	theFunc = reflect.ValueOf(svc.CreateStack).Pointer()
	//}


	for _, profile := range *profiles {
		// Specify profile for config and region for requests
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: profile,
			SharedConfigState: session.SharedConfigEnable,
		}))


		svc := cloudformation.New(sess)

		svc.
		if(delete) {
			result, err := svc.DeleteStack(input)
		}


		if err != nil {
			fmt.Println("I got an error.")
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:

					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		fmt.Println(result)

	}

}
