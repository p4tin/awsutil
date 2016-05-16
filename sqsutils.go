package main

import (
	"os"
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/joho/godotenv"
	"strings"
)

var queueName 		string
var sqs_svc             *sqs.SQS
var sns_svc 		*sns.SNS

var aws_region		string
var sqs_url 		string
var sns_url 		string

var attrib 		string
var topicArn		string

func init() {
	err := godotenv.Load()
	check(err)

	sqs_url = os.Getenv("SQS_URL")
	sns_url = os.Getenv("SNS_URL")
	aws_region = os.Getenv("AWS_REGION")
}

func main() {
	app := cli.NewApp()

	app.Version = "1.1"
	app.Name = "sqsutil"
	app.Usage = "Utility to work with SQS/ElasticMq Queues"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "action, a",
			Value: "depth",
			Usage: "action to perform on the queue (depth, create, purge, send, receive)",
		},
		//cli.StringFlag{
		//	Name: "server, s",
		//	Value: "localhost:9324",
		//	Usage: "server and port number for the request",
		//},
		//cli.StringFlag{
		//	Name: "amazonId, i",
		//	Value: "",
		//	Usage: "Amazon ID for the request",
		//},
		cli.StringFlag{
			Name: "stringBody, t",
			Value: "",
			Usage: "String body for send message",
		},
		cli.StringFlag{
			Name: "fileBody, f",
			Value: "",
			Usage: "File body for send message",
		},
	}

	app.Action = func(c *cli.Context) {
		queue := ""
		if len(c.Args()) > 0 {
			queue = c.Args()[0]
		}
		topicName := ""
		if len(c.Args()) > 1 {
			topicName = c.Args()[1]
		}
		action := c.String("action")
		server := c.String("server")
		amazonId := c.String("amazonId")
		fileBody := c.String("fileBody")
		strBody := c.String("stringBody")
		createEndpoint(queue, server, amazonId)
		switch action {
		case "list-queues":
			ListQueues()
		case "create-queue":
			createSQSQueue()
		case "depth":
			getSQSQueueDepth()
		case "purge":
			purgeQueue()
		case "send":
			sendMessage(fileBody, strBody)
		case "receive":
			receiveMessage()
		case "delete-queue":
			DeleteQueue()
		case "list-topics":
			ListTopics()
		case "create-topic":
			CreateTopic(queue)
		case "create-subscription":
			CreateSqsSubscription(topicName)
		case "send-topic":
			createTopicArn(queue, amazonId)
			SendSnsMsg(fileBody, strBody)
		case "delete-topic":
			DeleteTopic(topicName)
		default:
			fmt.Println("Unrecognized action - try `sqsutil -h` for help.")
		}
	}
	app.Run(os.Args)
}

func createEndpoint(queue string, server string, amazonId string) {
	if amazonId == "" {
		sqs_svc = sqs.New(session.New(), &aws.Config{Endpoint: aws.String(sqs_url), Region: aws.String(aws_region)})
		sqs_url = sqs_url + "/queue/" + queue
		sns_svc = sns.New(session.New(), &aws.Config{Endpoint: aws.String(sns_url), Region: aws.String(aws_region)})
	} else {
		sqs_url = "https://" + server + "/" + amazonId + "/" + queue
		sqs_svc = sqs.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
		sns_svc = sns.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	}
	queueName = queue
}

func ListQueues() {
	params := &sqs.ListQueuesInput{}
	resp, err := sqs_svc.ListQueues(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	for _, queue := range resp.QueueUrls {
		fmt.Println(*queue)
	}
}

func createSQSQueue() {
	params := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName), // Required
	}
	resp, err := sqs_svc.CreateQueue(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)
}

func GetQueueArn() string {
	params := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(sqs_url), // Required
		AttributeNames: []*string{
			aws.String("QueueArn"), // Required
		},
	}
	resp, err := sqs_svc.GetQueueAttributes(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return ""
	}

	return *resp.Attributes["QueueArn"]
}

func getSQSQueueDepth() {

	attrib = "ApproximateNumberOfMessages"
	sendParams := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(sqs_url), // Required
		AttributeNames: []*string{
			&attrib, // Required
		},
	}
	resp2, sendErr := sqs_svc.GetQueueAttributes(sendParams)
	if sendErr != nil {
		fmt.Println("Depth: " + sendErr.Error())
		return
	}
	fmt.Println(resp2)
}

func sendMessage(file string, str string) {
	msg := "Testing 1,2,3,..."
	if str != "" {
		msg = str
	} else if file != "" {
		dat, err := ioutil.ReadFile(file)
		check(err)
		msg = string(dat)
	}
	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(msg), // Required
		QueueUrl:     aws.String(sqs_url), // Required
	}
	resp, err := sqs_svc.SendMessage(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)

}


func receiveMessage() {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(sqs_url), // Required
	}
	resp, err := sqs_svc.ReceiveMessage(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(resp.Messages) > 0 {
		for _, msg := range resp.Messages {
			fmt.Println(msg)
			delParams := &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(sqs_url),                        // Required
				ReceiptHandle: aws.String(*msg.ReceiptHandle), // Required
			}
			sqs_svc.DeleteMessage(delParams)
		}
	}
}

func purgeQueue() {
	params := &sqs.PurgeQueueInput{
		QueueUrl: aws.String(sqs_url), // Required
	}
	resp, err := sqs_svc.PurgeQueue(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)
}

func DeleteQueue() {
	params := &sqs.DeleteQueueInput{
		QueueUrl: aws.String(sqs_url), // Required
	}
	resp, err := sqs_svc.DeleteQueue(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func createTopicArn(topic string, amazonId string) {
	if amazonId == "" {
		topicArn = "arn:aws:sns:yopa-local:000000000000:" + topic
	} else {
		topicArn = "arn:aws:sns:us-east-1:" + amazonId + ":" + topic
	}
}


func GetTopicArn(topicName string) string {
	params := &sns.ListTopicsInput{
		NextToken: aws.String(""),
	}
	resp, err := sns_svc.ListTopics(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return ""
	}
	// Pretty-print the response data.
	for _, arn := range resp.Topics {
		if strings.HasSuffix(*arn.TopicArn, topicName){
			return *arn.TopicArn
		}
	}
	return ""
}

func ListTopics() {
	params := &sns.ListTopicsInput{
		NextToken: aws.String(""),
	}
	resp, err := sns_svc.ListTopics(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}
	// Pretty-print the response data.
	for _, arn := range resp.Topics {
		fmt.Println(*arn.TopicArn)
		params := &sns.ListSubscriptionsByTopicInput{
			TopicArn:  aws.String(*arn.TopicArn), // Required
		}
		resp, err := sns_svc.ListSubscriptionsByTopic(params)

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return
		}

		if len(resp.Subscriptions) > 0 {
			// Pretty-print the response data.
			fmt.Println("Subscriptions:")
			for _, subs := range resp.Subscriptions {
				if *subs.Protocol == "sqs" {
					fmt.Println("\t", *subs.Endpoint)
				}
			}
		}
	}

}

func CreateTopic(name string) {
	params := &sns.CreateTopicInput{
		Name: aws.String(name), // Required
	}
	resp, err := sns_svc.CreateTopic(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func CreateSqsSubscription(topicName string) {
	params := &sns.SubscribeInput{
		Protocol: aws.String("sqs"),
		TopicArn: aws.String(GetTopicArn(topicName)),
		Endpoint: aws.String(GetQueueArn()),
	}
	resp, err := sns_svc.Subscribe(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func SendSnsMsg(file string, str string) {
	msg := "Testing 1,2,3,..."
	if str != "" {
		msg = str
	} else if file != "" {
		dat, err := ioutil.ReadFile(file)
		check(err)
		msg = string(dat)
	}

	//Create a session object to talk to SNS (also make sure you have your key and secret setup in your .aws/credentials file)
	//svc := sns.New(session.New())
	// params will be sent to the publish call included here is the bare minimum params to send a message.
	params := &sns.PublishInput{
		Message: aws.String(msg), // This is the message itself (can be XML / JSON / Text - anything you want)
		TopicArn: aws.String(topicArn), //Get this from the Topic in the AWS console.

	}

	resp, err := sns_svc.Publish(params)   //Call to puclish the message

	if err != nil {
		//Check for errors
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func DeleteTopic(topicName string) {
	params := &sns.DeleteTopicInput{
		TopicArn: aws.String(GetTopicArn(topicName)), // Required
	}
	resp, err := sns_svc.DeleteTopic(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
