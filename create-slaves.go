package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"sync"
	"io/ioutil"
)

var (

	pTargets		= kingpin.Arg("targets", "The AWS slaves to deploy the slave role too").Required().Strings()
	pDelete		= kingpin.Flag("delete", "Delete the slaves").Short('d').Bool()
	pAdmin		= kingpin.Flag("admin", "The admin account to be used by stack setst").Short('a').Required().String()
	pMaster		= kingpin.Flag("master", "The master account for the organization").Short('m').Required().String()
	pRegion    	= kingpin.Flag("region", "Region to deploy the slave roles too").Short('r').Required().String()
	pAccRole		= kingpin.Flag("role", "The cross account role").Short('x').Required().String()
	pSlaveStack 	= kingpin.Flag("slavestack",  "The slave stack cfn file").Short('s').Default("slave.yaml").String()
	pMasterStack 	= kingpin.Flag("masterstack",  "The master stack yaml file").Default("master.yml").String()


)

const kSlaveStackSetName = "cloudformation-stack-sets-slave-role"
const kMAsterStackSetName = "cloudformation-stack-sets-master-role"

func errHandler(err error) {

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

}

func main() {
	//Parse Command Line
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	//Call Function here depending on flag that was passed without do an if in the for loop
	if(*pDelete) {
		//Delete flag was passed so delete things


		//Delete from the master first

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Profile: *pMaster,
			SharedConfigState: session.SharedConfigEnable,
		}))

		cfSvc := cloudformation.New(sess)

		input := &cloudformation.DeleteStackInput{
			StackName:    aws.String(kSlaveStackSetName),
		}

		result, err := cfSvc.DeleteStack(input)

		errHandler(err)

		fmt.Println(result)

		//for _, profile := range *pTargets {
		//	// Specify profile for config and region for requests
		//	sess := session.Must(session.NewSessionWithOptions(session.Options{
		//		Profile: profile,
		//		SharedConfigState: session.SharedConfigEnable,
		//	}))
		//
		//	svc := cloudformation.New(sess)
		//
		//	input:= &cloudformation.DeleteStackInput{
		//		StackName: aws.String(kSlaveStackSetName),
		//	}
		//
		//	result, err := svc.DeleteStack(input)
		//
		//	errHandler(err)
		//
		//	fmt.Println(result)
		//}

	} else {
		//Create things

		//The master account profile is for an IAM user created in the master account.  We assume that during creation of
		//child/slave accounts from the master that a cross account role was created

		//Load the slave yaml file.  Thread the read as iops is slow nd we can set the params first
		var slaveTemplateBody string;
		var cfSvc *cloudformation.CloudFormation
		var sess *session.Session

		prepThreads :=2;
		var prepWg sync.WaitGroup

		prepWg.Add(prepThreads)

		go func() {
			defer prepWg.Done()
			if data, err := ioutil.ReadFile(*pSlaveStack); err != nil {
				errHandler(err)
			} else {
				slaveTemplateBody = string(data)
			}
		}()

		go func() {
			defer prepWg.Done()
			sess = session.Must(session.NewSessionWithOptions(session.Options{
				Profile: *pMaster,
				SharedConfigState: session.SharedConfigEnable,
			}))

			cfSvc = cloudformation.New(sess)
		}()


		var awsParameters []*cloudformation.Parameter
		awsParameters = append(awsParameters, &cloudformation.Parameter{
			ParameterKey: aws.String("MasterAccountId"),
			ParameterValue: aws.String(*pAdmin),
		})

		awsParameters = append(awsParameters, &cloudformation.Parameter{
			ParameterKey: aws.String("MasterRegion"),
			ParameterValue: aws.String(*pRegion),
		})

		input := &cloudformation.CreateStackInput{
			Capabilities: []*string{aws.String("CAPABILITY_NAMED_IAM")},
			StackName:    aws.String(kSlaveStackSetName),
			TemplateBody: &slaveTemplateBody,
			Parameters: awsParameters,
		}

		//Wait for the above threads to complete
		prepWg.Wait()

		//First we create a stack set target role in the master account as any stack sets run within or org won't be run from the master account
		//We can do this in a thread

		var wg sync.WaitGroup

		wg.Add(len(*pTargets)+1)

		go func () {
			defer wg.Done()
			result, err := cfSvc.CreateStack(input)

			errHandler(err)

			fmt.Println(result)
		}()


		// Specify profile for config and region for requests
		//sess := session.Must(session.NewSessionWithOptions(session.Options{
		//	Profile: *pMaster,
		//}))

		svc := sts.New(sess)

		//For the master account we must create the stack set slave role



		for _, account := range *pTargets {

			go func(account string) {
				defer wg.Done()

				roleArn := aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", account, *pAccRole))

				roleInput := &sts.AssumeRoleInput{
					RoleArn:         roleArn,
					RoleSessionName: aws.String(fmt.Sprintf("master-%s", account)),
				}



				result, err := svc.AssumeRole(roleInput)

				errHandler(err)

				fmt.Println(result)
			}(account)//go func
		}//for

		wg.Wait()
	}//else
}