package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func GetRegistryAuthorizationToken() (string, error) {
	sess, err := session.NewSession()
	svc := ecr.New(sess,&aws.Config{
		Region: aws.String("ap-south-1"),
	})
	input := &ecr.GetAuthorizationTokenInput{}
	result, err := svc.GetAuthorizationToken(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return "", nil
	}
	return *(result.AuthorizationData[0].AuthorizationToken), nil
}
