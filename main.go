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

var (
	Version = "No version specified"
)

const (
	exitCodeOk             int = 0
	exitCodeError          int = 1
	exitCodeFlagParseError     = 10 + iota
	exitCodeAWSError
)

const helpString = `Usage:
  dkv [-hiv] [--table=dynamo_table] [--delete=key] [--region=region] [key=value]

Flags:
  -d, --delete  Delete a key
  -h, --help    Print this help message
  -r, --region  The AWS region the table is in
  -t, --table   The name of the DynamoDB table
  -v, --version Print the version number
`

var (
	f = flag.NewFlagSet("flags", flag.ContinueOnError)

	// options
	deleteFlag  = f.StringP("delete", "d", "", "Delete a key")
	helpFlag    = f.BoolP("help", "h", false, "Show help")
	regionFlag  = f.StringP("region", "r", "us-east-1", "The AWS region")
	tableFlag   = f.StringP("table", "t", "", "The Dynamo table")
	versionFlag = f.BoolP("version", "v", false, "Print the version")
)

type Item struct {
	Name  string
	Value string
}

func main() {
	if err := f.Parse(os.Args[1:]); err != nil {
		fmt.Println("hmm")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(exitCodeFlagParseError)
	}

	if *helpFlag == true {
		fmt.Print(helpString)
		os.Exit(exitCodeOk)
	}

	if *versionFlag == true {
		fmt.Println(Version)
		os.Exit(exitCodeOk)
	}

	if *tableFlag == "" {
		fmt.Fprintln(os.Stderr, "Error: Missing table name")
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
