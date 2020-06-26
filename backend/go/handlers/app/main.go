package main

import (
	"context"
	"fmt"
	"os"
	"plugin"

	"github.com/acenteracms/acenteralib"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/myplugin/gofaas"
	models "github.com/myplugin/gofaas/models"
	uuid "github.com/satori/go.uuid"
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
	fmt.Println("INIT IN HERE MAIN")
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

	genHandlerObject := &models.GenericHandler{Awslambda: Awslambda}
	genHandlerObject.Initialize(r)

	/*
		genericHandler := models.GenericHandler{Awslambda: Awslambda, ElementType: ""}

		genericHandler := models.GenericHandler{Awslambda: Awslambda, ElementType: ""}
		genericHandler.InitializeRoutes(r)

		// projectHandler := ProjectHandler{gofaas.GenericHandler{Awslambda: Awslambda, ElementType: "PROJECT", ActionWord: "PROJECT"}}
		// projectHandler.InitializeRoutes(r)

		taskHandler := models.GenericHandler{ElementType: "TASKS"}
		taskHandler.InitializeRoutes(r)

		// postHandler := PostHandler{gofaas.GenericHandler{Awslambda: Awslambda, ElementType: "POSTS", ActionWord: "Posts"}}
		// postHandler.InitializeRoutes(r)
	*/
}

func RegisterRoutes(r *gin.Engine) {
	// REST API Queries here
	r.POST(fmt.Sprintf("/api/plugins/%s/setup/bootstrap", gofaas.PLUGIN_NAME), Awslambda.PluginNotifyAPIGateway(gofaas.AppPluginSiteBootstrap))
	r.GET(fmt.Sprintf("/api/plugins/%s/settings", gofaas.PLUGIN_NAME), Awslambda.PluginNotifyAPIGateway(gofaas.GetSettings))

	r.GET("/api/internal/test", Awslambda.PluginNotifyAPIGatewayJWTSecured(gofaas.ExecuteFast))
	r.POST("/api/internal/test", Awslambda.PluginNotifyAPIGatewayJWTSecured(gofaas.ExecuteFast))
	r.GET("/test", Awslambda.PluginNotifyAPIGatewayJWTSecured(gofaas.ExecuteFast))
	r.POST("/test", Awslambda.PluginNotifyAPIGatewayJWTSecured(gofaas.ExecuteFast))

	// Add any others
	// authenticated ones are, examle
	r.GET(fmt.Sprintf("/api/plugins/%s/me", gofaas.PLUGIN_NAME), Awslambda.PluginNotifyAPIGatewayJWTSecured(gofaas.GetMe))
}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if ginLambda == nil {
		// stdout and stderr are sent to AWS CloudWatch Logs
		r := gin.Default()
		RegisterRoutes(r)
		ginLambda = ginadapter.New(r)
	}

	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.Proxy(req)
}

func main() {
	// if (os.Getenv("PLUGIN_ACTION") == "Authorizer") {
	// lambda.Start(gofaas.NotifyAPIGatewayJWTAuth(gofaas.AuthorizerHandler))
	// } else {

	if os.Getenv("RUNWS") != "" {
		r := gin.Default()
		RegisterRoutes(r)
		fmt.Println("TEST PORT A\n")
		THEPORT := os.Getenv("CUSTOM_PORT")
		fmt.Println("TEST PORT A " + THEPORT + "\n")
		if THEPORT == "" {
			THEPORT = "3000"
		}
		fmt.Println("TEST PORT B " + THEPORT + "\n")
		r.Run(":" + THEPORT) // listen and serve on 0.0.0.0:8080
	} else {
		if os.Getenv("TYPE") == "WEBSITE" {
			lambda.Start(Awslambda.NotifyAPIGateway(WebsitePublic))
		} else if os.Getenv("TYPE") == "MODELSWS" {
			lambda.Start(Awslambda.NotifyAppSyncJWTSecure(ExecuteFast))
			//lambda.Start(Awslambda.NotifyAPIGateway(WebsitePublic))
		} else if os.Getenv("TYPE") == "MODELS" {
			lambda.Start(Awslambda.NotifyAppSyncJWTSecure(Execute))
		} else {
			lambda.Start(Handler)
		}
	}
}

func ExecuteFast(sharedlib acenteralib.SharedLib, ctx context.Context, reqObj acenteralib.RequestObject, e resolvers.Invocation) (interface{}, error) {
	fmt.Println("GOT REQUEST HERE OF\n")
	fmt.Println(ctx)
	fmt.Println(reqObj)
	fmt.Println(e)
	fmt.Println("======\n")
	return r.Handle(e, reqObj)
}

func Execute(sharedlib acenteralib.SharedLib, ctx context.Context, reqObj acenteralib.RequestObject, e resolvers.Invocation) (interface{}, error) {
	fmt.Println("GOT REQUEST HERE OF\n")
	fmt.Println(ctx)
	fmt.Println(reqObj)
	fmt.Println(e)
	fmt.Println("======\n")
	return r.Handle(e, reqObj)
}
