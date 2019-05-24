package gofaas

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	jwt "github.com/dgrijalva/jwt-go"

	// "github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	// "github.com/aws/aws-xray-sdk-go/xray"
)

func IfThenElse(condition bool, a string, b string) string {
	if condition {
		return a
	}
	return b
}

func getStageFromEnv() string {
	s := os.Getenv("STAGE")
	if s == "" {
		return "prod"
	}
	return s
}

func getPluginNameFromEnv() string {
	s := os.Getenv("PLUGIN_NAME")
	if s == "" {
		return "undefined"
	}
	return s
}

var (
	PluginVersion = "0.0.1"
	PLUGIN_NAME   = getPluginNameFromEnv()
	Stage         = getStageFromEnv()
)

// HandlerAPIGateway is an API Gateway Proxy Request handler function
type HandlerAPIGateway func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type RequestObject struct {
	User    *User    `json:"user"`
	Session *Session `json:"session"`
}

type UuidResponse struct {
	Uuid string `json:"uuid"`
}

// User represents a user
type User struct {
	ID               string   `json:"id"`
	Username         string   `json:"username"`
	Password         string   `json:"-"`
	SessionHash      string   `json:"session_hash"`
	Token            string   `json:"token"` //used to send back the JWT Token for the end-user session // we should move this
	Roles            []string `json:"roles"`
	Lock             string   `json:"lock"`
	PasswordResetJWT string   `json:"password_reset"`
	Sub              string   `json:"sub"`
	Aud              string   `json:"aud"`
	// UserObject 				CognitoUser `json:"-"` // Used for password reset
}

type SessionJWTClaim struct {
	RealSessionId string `json:"real_session_id"`
	AnonSessionId string `json:"anon_session_id"`
	CreatedAt     int64  `json:"at"`
	SessionId     string `json:"session_id"`
	Salt          string `json:"salt"`
	Username      string `json:"username"`
	jwt.StandardClaims
}

type Session struct {
	ID      string `json:"sessionid"`
	Ttl     int64  `json:"ttl"`
	Userid  string `json:"userid"`
	Sshkeys string `json:"sshkeys"`
	Token   string `json:"token"`
}

// DynamoDBAPI is a subset of dynamodbiface.DynamoDBAPI
type DynamoDBAPI interface {
	DeleteItemWithContext(ctx aws.Context, input *dynamodb.DeleteItemInput, opts ...request.Option) (*dynamodb.DeleteItemOutput, error)
	GetItemWithContext(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error)
	PutItemWithContext(ctx aws.Context, input *dynamodb.PutItemInput, opts ...request.Option) (*dynamodb.PutItemOutput, error)
	UpdateItemWithContext(ctx aws.Context, input *dynamodb.UpdateItemInput, opts ...request.Option) (*dynamodb.UpdateItemOutput, error)
	QueryWithContext(ctx aws.Context, input *dynamodb.QueryInput, opts ...request.Option) (*dynamodb.QueryOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}
type S3API interface {
	PutObjectWithContext(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error)
}

type EC2API interface {
	DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
	DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	DescribeNatGateways(input *ec2.DescribeNatGatewaysInput) (*ec2.DescribeNatGatewaysOutput, error)
	DescribeInternetGateways(input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error)
	DescribeSpotPriceHistory(input *ec2.DescribeSpotPriceHistoryInput) (*ec2.DescribeSpotPriceHistoryOutput, error)
}

type ECSAPI interface {
	ListServices(*ecs.ListServicesInput) (*ecs.ListServicesOutput, error)
	DescribeClusters(input *ecs.DescribeClustersInput) (*ecs.DescribeClustersOutput, error)
	DescribeTasks(input *ecs.DescribeTasksInput) (*ecs.DescribeTasksOutput, error)
	DescribeServices(input *ecs.DescribeServicesInput) (*ecs.DescribeServicesOutput, error)
	DescribeTaskDefinition(input *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error)
	DescribeContainerInstances(input *ecs.DescribeContainerInstancesInput) (*ecs.DescribeContainerInstancesOutput, error)
}

type AUTOSCALINGAPI interface {
	DescribeAutoScalingGroups(input *autoscaling.DescribeAutoScalingGroupsInput) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
	DescribeAutoScalingInstances(input *autoscaling.DescribeAutoScalingInstancesInput) (*autoscaling.DescribeAutoScalingInstancesOutput, error)
}

type IAMAPI interface {
	ListRolesWithContext(ctx aws.Context, input *iam.ListRolesInput, opts ...request.Option) (*iam.ListRolesOutput, error)
	ListInstanceProfilesWithContext(ctx aws.Context, input *iam.ListInstanceProfilesInput, opts ...request.Option) (*iam.ListInstanceProfilesOutput, error)
}
type HandlerAPIGatewayWithJWT func(context.Context, jwt.Claims, RequestObject, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// type NotifyAPIGatewayJWTSecured func(context.Context, jwt.Claims, RequestObject, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

/*
type AWSLambda interface {
	NotifyAPIGatewayJWTSecured(h HandlerAPIGatewayWithJWT) HandlerAPIGateway
}

type LambdaHelper interface {
	NotifyAPIGatewayJWTSecured(h HandlerAPIGatewayWithJWT) HandlerAPIGateway
}
*/
type PlgFilter interface {
	NotifyAPIGatewayJWTSecured(func()) interface{}
	Name() string
	Age() int
}

// type HandlerAPIGateway func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func NotifyAPIGateway(h HandlerAPIGateway) HandlerAPIGateway {
	//TODO: Implement Authorizer?
	return func(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		r, err := h(ctx, e)
		// notify(ctx, err)
		return r, err
	}
}

var (
	sess = session.Must(session.NewSession(&aws.Config{
		Endpoint:   aws.String(IfThenElse(os.Getenv("AWS_SAM_LOCAL") == "true", "http://172.17.0.1:8000", "")),
		MaxRetries: aws.Int(5),
	},
	))
	DynamoDB = NewDynamoDB()
	S3       = NewS3()
	EC2      = NewEC2()
	IAM      = NewIAM()
	ASG      = NewASG(sess)
	ECS      = NewECS(sess)
)

func NewDynamoDB() DynamoDBAPI {
	c := dynamodb.New(sess)
	// xray.AWS(c.Client)
	return c
}

func NewS3() S3API {
	c := s3.New(sess)
	// xray.AWS(c.Client)
	return c
}

func NewEC2() EC2API {
	c := ec2.New(sess)
	// xray.AWS(c.Client)
	return c
}

func NewECS(ses *session.Session) ECSAPI {
	c := ecs.New(ses)
	// xray.AWS(c.Client)
	return c
}

func NewASG(ses *session.Session) AUTOSCALINGAPI {
	c := autoscaling.New(ses)
	// xray.AWS(c.Client)
	return c
}

func NewIAM() IAMAPI {
	c := iam.New(sess)
	// xray.AWS(c.Client)
	return c
}
