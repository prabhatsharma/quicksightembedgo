package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/quicksight"
)

func main() {

	//Create a session
	sess := session.Must(session.NewSession())

	// This is the IAM role that will be assumed for createing users. This should have the permission to call GetDashboardEmbedUrl
	// Also the user/role(could be an ec2 role, ecs task role or EKS pod role or lambda execution role) that is executing this go program must has the stsAssumeRole permission to assume the roleName
	roleName := "quicksight-embedded"

	awsAccountID := "107995894928"
	iamRoleARN := "arn:aws:iam::" + awsAccountID + ":role/" + roleName
	userEmail := "someone@gmail.com"
	identityType := "IAM"
	userRegistrationRegion := "us-east-1"
	namespace := "default"
	userRole := "READER"

	// Step 1 - AssumeRole
	creds := stscreds.NewCredentials(sess, iamRoleARN)

	//Step 2: Register User. This might fail if user already exist, but no harm or foul if this fails due to UserAlready exists. Just continue
	client := quicksight.New(sess, &aws.Config{Credentials: creds, Region: &userRegistrationRegion})

	ruInput := quicksight.RegisterUserInput{
		AwsAccountId: &awsAccountID,
		Email:        &userEmail,
		IamArn:       &iamRoleARN,
		Namespace:    &namespace,
		IdentityType: &identityType,
		SessionName:  &userEmail,
		UserRole:     &userRole,
	}

	ruOutput, ruOutputError := client.RegisterUser(&ruInput)
	if ruOutputError != nil {
		fmt.Println(ruOutputError.Error())
	} else {
		fmt.Println(ruOutput.String())
	}

	// Step 3: Get the embeddedURL
	dashboardID := "81d2ae9f-57bf-42b1-ad9e-9703718f36f6"
	userDashboardRegion := "us-west-2"

	// Need to create separate client since dashboard region could be different from us-east-1 which is the user region
	client2 := quicksight.New(sess, &aws.Config{Credentials: creds, Region: &userDashboardRegion})
	userARN := "arn:aws:quicksight:us-east-1:" + awsAccountID + ":user/" + namespace + "/" + roleName + "/" + userEmail

	dashboardIdentityType := "QUICKSIGHT"

	eURLInput := quicksight.GetDashboardEmbedUrlInput{
		AwsAccountId: &awsAccountID,
		DashboardId:  &dashboardID,
		IdentityType: &dashboardIdentityType, //Needs to be QUICKSIGHT here and not IAM even  though an IAM role is being used that assumes the role
		UserArn:      &userARN,
	}

	eURLOutput, errEmbed := client2.GetDashboardEmbedUrl(&eURLInput)

	if errEmbed != nil {
		fmt.Println("\nStep 3.2 - ", errEmbed.Error())
	} else {
		fmt.Println(eURLOutput)
	}
}
