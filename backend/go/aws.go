package gofaas

import (
	"context"
	"os"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	jwt "github.com/dgrijalva/jwt-go"

	// "github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
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

type COGNITOIDENTITYAPI interface {
	UpdateUserPoolClientWithContext(ctx aws.Context, input *cognitoidentityprovider.UpdateUserPoolClientInput, opts ...request.Option) (*cognitoidentityprovider.UpdateUserPoolClientOutput, error)
	UpdateUserPoolClient(input *cognitoidentityprovider.UpdateUserPoolClientInput) (*cognitoidentityprovider.UpdateUserPoolClientOutput, error)
	CreateUserPoolDomainWithContext(ctx aws.Context, input *cognitoidentityprovider.CreateUserPoolDomainInput, opts ...request.Option) (*cognitoidentityprovider.CreateUserPoolDomainOutput, error)
	CreateUserPoolDomain(input *cognitoidentityprovider.CreateUserPoolDomainInput) (*cognitoidentityprovider.CreateUserPoolDomainOutput, error)
	DeleteUserPoolDomainWithContext(ctx aws.Context, input *cognitoidentityprovider.DeleteUserPoolDomainInput, opts ...request.Option) (*cognitoidentityprovider.DeleteUserPoolDomainOutput, error)
	DeleteUserPoolDomain(input *cognitoidentityprovider.DeleteUserPoolDomainInput) (*cognitoidentityprovider.DeleteUserPoolDomainOutput, error)
	DescribeUserPoolDomainWithContext(ctx aws.Context, input *cognitoidentityprovider.DescribeUserPoolDomainInput, opts ...request.Option) (*cognitoidentityprovider.DescribeUserPoolDomainOutput, error)
	DescribeUserPoolDomain(input *cognitoidentityprovider.DescribeUserPoolDomainInput) (*cognitoidentityprovider.DescribeUserPoolDomainOutput, error)
	CreateGroupWithContext(ctx aws.Context, input *cognitoidentityprovider.CreateGroupInput, opts ...request.Option) (*cognitoidentityprovider.CreateGroupOutput, error)
	CreateGroup(input *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error)
	GetGroupWithContext(ctx aws.Context, input *cognitoidentityprovider.GetGroupInput, opts ...request.Option) (*cognitoidentityprovider.GetGroupOutput, error)
	GetGroup(input *cognitoidentityprovider.GetGroupInput) (*cognitoidentityprovider.GetGroupOutput, error)
	AdminCreateUserWithContext(ctx aws.Context, input *cognitoidentityprovider.AdminCreateUserInput, opts ...request.Option) (*cognitoidentityprovider.AdminCreateUserOutput, error)
	AdminCreateUser(input *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error)
	AdminAddUserToGroupWithContext(ctx aws.Context, input *cognitoidentityprovider.AdminAddUserToGroupInput, opts ...request.Option) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
	AdminAddUserToGroup(input *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
	GetUserWithContext(ctx aws.Context, input *cognitoidentityprovider.GetUserInput, opts ...request.Option) (*cognitoidentityprovider.GetUserOutput, error)
	GetUser(input *cognitoidentityprovider.GetUserInput) (*cognitoidentityprovider.GetUserOutput, error)
	AdminGetUserWithContext(ctx aws.Context, input *cognitoidentityprovider.AdminGetUserInput, opts ...request.Option) (*cognitoidentityprovider.AdminGetUserOutput, error)
	AdminGetUser(input *cognitoidentityprovider.AdminGetUserInput) (*cognitoidentityprovider.AdminGetUserOutput, error)
	AdminInitiateAuthRequest(input *cognitoidentityprovider.AdminInitiateAuthInput) (req *request.Request, output *cognitoidentityprovider.AdminInitiateAuthOutput)
	InitiateAuthWithContext(ctx aws.Context, input *cognitoidentityprovider.InitiateAuthInput, opts ...request.Option) (*cognitoidentityprovider.InitiateAuthOutput, error)
	InitiateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, error)
	AssociateSoftwareTokenWithContext(ctx aws.Context, input *cognitoidentityprovider.AssociateSoftwareTokenInput, opts ...request.Option) (*cognitoidentityprovider.AssociateSoftwareTokenOutput, error)
	AssociateSoftwareToken(input *cognitoidentityprovider.AssociateSoftwareTokenInput) (*cognitoidentityprovider.AssociateSoftwareTokenOutput, error)
	VerifySoftwareTokenWithContext(ctx aws.Context, input *cognitoidentityprovider.VerifySoftwareTokenInput, opts ...request.Option) (*cognitoidentityprovider.VerifySoftwareTokenOutput, error)
	VerifySoftwareToken(input *cognitoidentityprovider.VerifySoftwareTokenInput) (*cognitoidentityprovider.VerifySoftwareTokenOutput, error)
	AdminSetUserMFAPreferenceWithContext(ctx aws.Context, input *cognitoidentityprovider.AdminSetUserMFAPreferenceInput, opts ...request.Option) (*cognitoidentityprovider.AdminSetUserMFAPreferenceOutput, error)
	AdminSetUserMFAPreference(input *cognitoidentityprovider.AdminSetUserMFAPreferenceInput) (*cognitoidentityprovider.AdminSetUserMFAPreferenceOutput, error)
	SetUserMFAPreference(input *cognitoidentityprovider.SetUserMFAPreferenceInput) (*cognitoidentityprovider.SetUserMFAPreferenceOutput, error)
	SetUserPoolMfaConfig(input *cognitoidentityprovider.SetUserPoolMfaConfigInput) (*cognitoidentityprovider.SetUserPoolMfaConfigOutput, error)
	AdminConfirmSignUp(input *cognitoidentityprovider.AdminConfirmSignUpInput) (*cognitoidentityprovider.AdminConfirmSignUpOutput, error)
	AdminSetUserPassword(input *cognitoidentityprovider.AdminSetUserPasswordInput) (*cognitoidentityprovider.AdminSetUserPasswordOutput, error)
}
type COGNITOIDENTITY interface {
	GetOpenIdTokenWithContext(ctx aws.Context, input *cognitoidentity.GetOpenIdTokenInput, opts ...request.Option) (*cognitoidentity.GetOpenIdTokenOutput, error)
	GetOpenIdToken(input *cognitoidentity.GetOpenIdTokenInput) (*cognitoidentity.GetOpenIdTokenOutput, error)
	GetIdWithContext(ctx aws.Context, input *cognitoidentity.GetIdInput, opts ...request.Option) (*cognitoidentity.GetIdOutput, error)
	GetId(input *cognitoidentity.GetIdInput) (*cognitoidentity.GetIdOutput, error)
	GetCredentialsForIdentityWithContext(ctx aws.Context, input *cognitoidentity.GetCredentialsForIdentityInput, opts ...request.Option) (*cognitoidentity.GetCredentialsForIdentityOutput, error)
	GetCredentialsForIdentity(input *cognitoidentity.GetCredentialsForIdentityInput) (*cognitoidentity.GetCredentialsForIdentityOutput, error)
	GetIdentityPoolRolesWithContext(ctx aws.Context, input *cognitoidentity.GetIdentityPoolRolesInput, opts ...request.Option) (*cognitoidentity.GetIdentityPoolRolesOutput, error)
	GetIdentityPoolRoles(input *cognitoidentity.GetIdentityPoolRolesInput) (*cognitoidentity.GetIdentityPoolRolesOutput, error)
	SetIdentityPoolRoles(input *cognitoidentity.SetIdentityPoolRolesInput) (*cognitoidentity.SetIdentityPoolRolesOutput, error)
	// GetIdentityPoolRolesWithContext(ctx aws.Context, input *cognitoidentity.GetIdentityPoolRolesInput) (*cognitoidentity.GetIdentityPoolRolesOutput, error)
}

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
		Endpoint:   aws.String(IfThenElse(os.Getenv("AWS_SAM_LOCAL") == "true", "http://dynamodb:8000", "")),
		MaxRetries: aws.Int(5),
	},
	))

	cognitoSess = session.Must(session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_DEFAULT_REGION"))}))

	DynamoDB = NewDynamoDB()
	S3       = NewS3()
	EC2      = NewEC2()
	IAM      = NewIAM()
	ASG      = NewASG(sess)
	CognitoIdentityServiceProvider = NewCognitoIdentityServiceProvider()
	CognitoIdentity                = NewCognitoIdentity()
	ECS      = NewECS(sess)
)

// NewKMS is an xray instrumented KMS client
func NewCognitoIdentityServiceProvider() COGNITOIDENTITYAPI {
	c := cognitoidentityprovider.New(cognitoSess)
	// xray.AWS(c.Client)
	return c
}

func NewCognitoIdentity() COGNITOIDENTITY {
	c := cognitoidentity.New(cognitoSess)
	// xray.AWS(c.Client)
	return c
}

func NewDynamoDB() DynamoDBAPI {
	fmt.Println("1a OK TEST DYNAMO OF")
	fmt.Println(IfThenElse(os.Getenv("AWS_SAM_LOCAL") == "true", "http://dynamodb:8000", ""))
	fmt.Println("1a OK TEST DYNAMO OF - 1")

	sessTmp := sess
	if (os.Getenv("AWS_SAM_LOCAL") == "true") {
		fmt.Println("1a OK TEST DYNAMO OF IN SAM LOCAL")
		sessTmp = session.Must(session.NewSession(&aws.Config{
			Endpoint:   aws.String(IfThenElse(os.Getenv("AWS_SAM_LOCAL") == "true", "http://dynamodb:8000", "")),
			Credentials: credentials.NewStaticCredentials("a", "a", ""),
			MaxRetries: aws.Int(5),
		},
		))
	}

	c := dynamodb.New(sessTmp)
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
