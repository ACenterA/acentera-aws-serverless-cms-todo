package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"plugin"

	"github.com/acenteracms/acenteralib"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/myplugin/gofaas"
	"github.com/satori/go.uuid"
	resolvers "github.com/sbstjn/appsync-resolvers"
)

var ginLambda *ginadapter.GinLambda

var Awslambda acenteralib.SharedLib

var UUIDGen = func() uuid.UUID {
	uv4, _ := uuid.NewV4()
	return uv4
}

var r = resolvers.New()

func init() {
	// load modulez
	// 1. open the so file to load the symbols
	plug, err := plugin.Open("aws.so")
	if err != nil {
		plug, err = plugin.Open("/opt/aws.so")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	GetFilter, err := plug.Lookup("GetFilter")
	if err != nil {
		panic(err)
	}
	filter, err := GetFilter.(func() (acenteralib.SharedLib, error))()
	Awslambda = filter

	genericHandler := GenericHandler{ElementType: ""}
	genericHandler.InitializeRoutes(r)

	projectHandler := GenericHandler{ElementType: "PROJECT"}
	projectHandler.InitializeRoutes(r)

	taskHandler := GenericHandler{ElementType: "TASKS"}
	taskHandler.InitializeRoutes(r)

	postHandler := GenericHandler{ElementType: "POSTS"}
	postHandler.InitializeRoutes(r)
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if ginLambda == nil {
		log.Printf("Gin cold start")
		// stdout and stderr are sent to AWS CloudWatch Logs
		r := gin.Default()
		r.POST("/api/plugins/serverless-cms/setup/bootstrap", Awslambda.PluginNotifyAPIGateway(gofaas.AppPluginSiteBootstrap))
		r.GET("/api/plugins/serverless-cms/settings", Awslambda.PluginNotifyAPIGateway(gofaas.GetSettings))
		// r.GET("/pets/:id", getPet)
		// r.POST("/pets", createPet)

		ginLambda = ginadapter.New(r)
	}

	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.Proxy(req)
}

func main() {
	// if (os.Getenv("PLUGIN_ACTION") == "Authorizer") {
	// lambda.Start(gofaas.NotifyAPIGatewayJWTAuth(gofaas.AuthorizerHandler))
	// } else {

	if os.Getenv("TYPE") == "WEBSITE" {
		lambda.Start(Awslambda.NotifyAPIGateway(WebsitePublic))
	} else if os.Getenv("TYPE") == "MODELS" {
		lambda.Start(Awslambda.NotifyAppSyncJWTSecure(Execute))
	} else {
		lambda.Start(Handler)
	}
}

func Execute(sharedlib acenteralib.SharedLib, ctx context.Context, reqObj acenteralib.RequestObject, e resolvers.Invocation) (interface{}, error) {
	return r.Handle(e, reqObj)
}
