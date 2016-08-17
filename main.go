package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	flag "github.com/ogier/pflag"
)

const (
	exitCodeOk             int = 0
	exitCodeError          int = 1
	exitCodeFlagParseError     = 10 + iota
	exitCodeAWSError
)

var (
	f = flag.NewFlagSet("flags", flag.ContinueOnError)

	// options
	deleteFlag = f.StringP("delete", "d", "", "Delete a key")
	regionFlag = f.StringP("region", "r", "us-east-1", "The AWS region")
	tableFlag  = f.StringP("table", "t", "", "The Dynamo table")
)

type Item struct {
	Name  string
	Value string
}

// type DynamoClient interface {
// 	Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
// }

func main() {
	if err := f.Parse(os.Args[1:]); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(exitCodeFlagParseError)
	}

	if *tableFlag == "" {
		fmt.Fprint(os.Stderr, "Error: Missing table name")
		os.Exit(exitCodeFlagParseError)
	}

	// setup dynamo client
	sess, err := session.NewSession(&aws.Config{Region: regionFlag})
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(exitCodeError)
	}
	svc := dynamodb.New(sess)

	if *deleteFlag != "" {
		params := &dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"Name": &dynamodb.AttributeValue{
					S: deleteFlag,
				},
			},
			TableName: tableFlag,
		}
		_, err := svc.DeleteItem(params)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(exitCodeError)
		}
		os.Exit(0)
	}

	args := f.Args()
	if len(args) == 0 {
		params := &dynamodb.ScanInput{
			TableName: tableFlag,
		}
		resp, err := svc.Scan(params)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(exitCodeError)
		}
		for _, item := range resp.Items {
			fmt.Printf("%s=%s\n", *item["Name"].S, *item["Value"].S)
		}
		os.Exit(exitCodeOk)
	}

	// set
	for _, pair := range args {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) < 2 {
			fmt.Printf("Error: \"%s\" is not a valid key-value pair\n", pair)
			os.Exit(exitCodeError)
		}
		params := &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"Name": &dynamodb.AttributeValue{
					S: aws.String(parts[0]),
				},
				"Value": &dynamodb.AttributeValue{
					S: aws.String(parts[1]),
				},
			},
			TableName: tableFlag,
		}
		_, err := svc.PutItem(params)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(exitCodeError)
		}
	}

	os.Exit(exitCodeOk)
}
