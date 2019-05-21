package gofaas

import (
	// "context"
	"fmt"
	// "time"
	"os"
	"strings"

	// "math/rand"
	// "strings"
	"crypto/md5"
	"encoding/hex"

	// "crypto/sha1"
	// "encoding/base64"

	// "github.com/pkg/errors"
	"github.com/aws/aws-lambda-go/events"
	// cfn "github.com/aws/aws-lambda-go/cfn"
	// "github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"

	// "github.com/aws/aws-sdk-go/service/cognitoidentity"
	// "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	// "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	// "github.com/aws/aws-sdk-go/service/kms"
	// "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/acenteracms/acenteralib"
)

/*
var Awslambda acenteralib.SharedLib

func init() {
// load modulez
// 1. open the so file to load the symbols
fmt.Println("Trying to open aws.so ...")
plug, err := plugin.Open("aws.so")
if err != nil {
	plug, err = plugin.Open("/opt/aws.so")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
fmt.Println("Done .. Trying to open aws.so ...")
GetFilter, err := plug.Lookup("GetFilter")
if err != nil {
		panic(err)
}
fmt.Println("Get Filter ...")
filter, err := GetFilter.(func() (acenteralib.SharedLib, error))()
Awslambda = filter
}
*/

func GetSettings(lib acenteralib.SharedLib, c *gin.Context, reqObj acenteralib.RequestObject) (events.APIGatewayProxyResponse, error) {
	// We need to know what is the current API Gateway ... ???
	// reqObj, Type, Name
	// Returns an Object with an auto generated UUID ...
	// fmt.Println("QUERY USING UUID of : ", reqObj.Site, " and default site")
	/*
						fmt.Println("GOT SHARED LIB OF ", sharedlib)
		        fmt.Println("GOT REQOBJ LIB OF ", reqObj)
		        fmt.Println("GOT CONTEXT OF ")
		        fmt.Println(ctx)
		        fmt.Println("GOT EVENT PROXY REQUEST OF ")
		        fmt.Println(e)
		        fmt.Println("AUTHORIZER OF ")
		        fmt.Println(e.RequestContext.Authorizer)
		        fmt.Println("DONE....")
	*/
	/*
	   a, errT := AMIGetByProject(ctx, reqObj.Site)
	   //fmt.Println(a)
	   //fmt.Println(errT)
	   if (errT != nil) {
	           fmt.Println(errT.Error())
	           return RestResponseError(errT.Error())
	   }

	   z, errF := GetByProject(ctx, reqObj.Site, "clusterprj")
	   //fmt.Println(z)
	   //fmt.Println(errF)
	   if (errF != nil) {
	           fmt.Println(errF)
	           return RestResponseError(errF.Error())
	   }

	   respObj := &PluginSettings{
	           AMIProjects: a,
	           Clusters: z,
	   }
	*/
	respObj := &PluginSettings{}
	// a
	fmt.Println("Returning response...")
	return acenteralib.RestResponseNoCache(respObj)
}

type AppBootstrapReq struct {
	Accountid int  `json:"accountid"` // Used to validate account id vs AWS::AccountID ...
	Eula      bool `json:"eula"`
}

// PluginNotifyAPIGateway(h HandlerPluginAPIGateway) (gin.HandlerFunc)
// PluginNotifyAPIGatewayJWTSecured(h HandlerPluginAPIGatewayWithJWT) (gin.HandlerFunc)
func AppPluginSiteBootstrap(lib acenteralib.SharedLib, c *gin.Context, reqObj acenteralib.RequestObject) (events.APIGatewayProxyResponse, error) {
	var appBootstrapReq AppBootstrapReq
	c.BindJSON(&appBootstrapReq)
	// fmt.Println("POST BOOTSTRAP TES 1T a")

	lambdadAccountId := "-1"
	stackName := os.Getenv("STACK_NAME")
	// fmt.Println("IS IT SAM", os.Getenv("AWS_SAM_LOCAL"))
	siteKey := os.Getenv("SITE_KEY")
	if siteKey == "" {
		siteKey = c.GetHeader("x-site")
	}

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%v", "admin")))
	adminTableHash := hex.EncodeToString(hasher.Sum(nil))

	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		lambdadAccountId = IfThenElse(os.Getenv("DEV_ACCOUNTID") == "", "0", os.Getenv("DEV_ACCOUNTID"))
	} else {
		// Running in AWS
		s, err := acenteralib.STS.GetCallerIdentity(&sts.GetCallerIdentityInput{}) //if no sort key
		if err == nil {
			lambdadAccountId = fmt.Sprintf("%s", string(*s.Account))
		} else {
			lambdadAccountId = ""
		}
	}

	fmt.Println("a - Comparing of ", lambdadAccountId, " with ", appBootstrapReq.Accountid)
	if lambdadAccountId == fmt.Sprintf("%v", appBootstrapReq.Accountid) {
		if lambdadAccountId != "" {
			// nPluginInfo, errorPlugin := lib.GetPluginById(siteKey)
			nPluginInfo, errorPlugin := lib.GetSite(siteKey, adminTableHash)
			fmt.Println("SITE KEY ", siteKey)
			if errorPlugin != nil || nPluginInfo != nil {
				if errorPlugin != nil {
					fmt.Println(errorPlugin)
					r := events.APIGatewayProxyResponse{
						Body: string(fmt.Sprintf("{ \"message\": \"Plugin bootstrap error\" }")),
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
						StatusCode: 501,
					}
					return r, nil
				} else {
					r := events.APIGatewayProxyResponse{
						Body: string(fmt.Sprintf("{ \"message\": \"Plugin bootstrap already completed\" }")),
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
						StatusCode: 201,
					}
					return r, nil
				}
			}

			fmt.Println("SITE HASH APP SETTINGS...")
			// Ok Get Administrator Site by Hash so we can add ourself to the plugin list ...
			siteHasher := md5.New()
			siteHasher.Write([]byte(fmt.Sprintf("%v", "admin")))
			siteHash := hex.EncodeToString(siteHasher.Sum(nil))
			appSettings, err := lib.GetAppSettings(siteHash)
			fmt.Println("SITE HASH APP SETTINGS... is")
			fmt.Println(appSettings)
			if err != nil {
				r := events.APIGatewayProxyResponse{
					Body: string(fmt.Sprintf("{ \"message\": \"Could not access DB\" }")),
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					StatusCode: 500,
				}
				return r, nil
			}
			appSettings.Version = PluginVersion // this plugin version
			appSettings.Id = siteKey            //

			fmt.Println("1 - SITE HASH APP SETTINGS... is")
			stackUrl := fmt.Sprintf("https://console.aws.amazon.com/cloudformation/home?region=%s#/stacks?stackId=%s&filter=active&tab=outputs", os.Getenv("REGION"), os.Getenv("STACK_ID"))
			appSettings.StackUrl = stackUrl

			appSettings.Title = os.Getenv("SITE_TITLE") // this plugin site key
			if appSettings.Title == "" {
				appSettings.Title = "Todo CMS Portal" // Couldd also be an Env Pararmeter?
			}

			nPlugin := &acenteralib.AppSettings{}
			nPlugin.Id = os.Getenv("SITE_KEY") // this plugin site key // if defined?
			nPlugin.Title = os.Getenv("SITE_TITLE")
			if nPlugin.Title == "" {
				nPlugin.Title = os.Getenv("PLUGIN_NAME")
			}
			if nPlugin.Title == "" {
				nPlugin.Title = nPlugin.Id
			}
			//TODO: Add plugin.Version from Properties ... ?
			nPlugin.Stage = os.Getenv("STAGE")
			//TODO: Add plugin.Version from Properties ...

			// appSiteRef, errTmp := GetSite(refSiteHash, "UnUsedName")

			nPlugin.Aws = appSettings.Aws

			cognitoMap := nPlugin.Aws["cognito"]
			if (cognitoMap == nil) {
				nPlugin.Aws["cognito"] = map[string]string{}
				cognitoMap = nPlugin.Aws["cognito"]
			}

			// In case we receive different userepool as parameters ...
			cognitoMap["APP_CLIENT_ID"] = IfThenElse(os.Getenv("UserPoolClientId") == "", cognitoMap["APP_CLIENT_ID"], os.Getenv("UserPoolClientId"))
			cognitoMap["IDENTITY_POOL_ID"] = IfThenElse(os.Getenv("IdentityId") == "", cognitoMap["IDENTITY_POOL_ID"], os.Getenv("IdentityId"))
			cognitoMap["USER_POOL_ID"] = IfThenElse(os.Getenv("IdentityId") == "", cognitoMap["USER_POOL_ID"], os.Getenv("UserPoolId"))
			cognitoMap["REGION"] = IfThenElse(os.Getenv("IdentityId") == "", cognitoMap["REGION"], os.Getenv("REGION"))

			nPlugin.Graphql = appSettings.Graphql
			GraphApiEndpoint := os.Getenv("GraphApiEndpoint")
			if (GraphApiEndpoint != "") {
				nPlugin.Graphql = map[string]string{}
				nPlugin.Graphql["URL"] = GraphApiEndpoint
				nPlugin.Graphql["REGION"] = os.Getenv("REGION")
			} else {
				if nPlugin.Graphql == nil {
					// The Parent does not have a graphql endpoint ?
					// TODO: Get any GraphQL
					req, resp := acenteralib.Cloudformation.DescribeStacksRequest(&cloudformation.DescribeStacksInput{
						StackName: aws.String(stackName),
					})
					fmt.Println("3 - SITE HASH APP SETTINGS... is")
					errSend := req.Send()
					if errSend != nil { // resp is now filled
						fmt.Println(errSend)
						r := events.APIGatewayProxyResponse{
							Body: string("{ \"message\": \"Stack not found\" }\n"),
							Headers: map[string]string{
								"Content-Type": "application/json",
							},
							StatusCode: 404,
						}
						return r, nil
					}
					stacks := resp.Stacks
					stackObject := stacks[0]
					GraphQLEndpoint := ""
					for _, outputIface := range stackObject.Outputs {
						// tmpOutput := outputIface.(map[string]string)
						if *outputIface.OutputKey == "ApiEndpoint" {
							GraphQLEndpoint = *outputIface.OutputValue
						}
					}
					if GraphQLEndpoint != "" {
						nPlugin.Graphql = map[string]string{}
						nPlugin.Graphql["URL"] = GraphQLEndpoint
						nPlugin.Graphql["REGION"] = os.Getenv("REGION")
					}
				}
			}

			pluginTmp := strings.ToLower(strings.Trim(nPlugin.Id, " "))

			// fmt.Println("Recieved creation of domain : ", pluginTmp)
			// fmt.Println("[DEBUG] Create Plugin ... ", pluginTmp)

			// validadte i not eaists_, err = GetSite(tableHash, pluginTmp)

			errTmp := lib.CreateSitePlugin(nPlugin.Id, pluginTmp, nPlugin, adminTableHash)
			// pluginUuid, errTmp := lib.CreateAppPlugin(nPlugin.Id, nPlugin.Title, nPlugin)
			fmt.Println(errTmp)
			if errTmp == nil {
				// fmt.Println("PluginUUID: ", pluginUuid)
				r := events.APIGatewayProxyResponse{
					Body: string(fmt.Sprintf("{ \"message\": \"Plugin bootstrap completed\" }")),
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					StatusCode: 201,
				}
				return r, nil
			}
			// Create the site informations
			// errCreate := AppSettingsCreation(appSettings)

			// Ok nice we have proper accountId defined ...
			/*
				      req, resp := Cloudformation.DescribeStacksRequest(&cloudformation.DescribeStacksInput{
				              StackName: aws.String(stackName),
				      })
							..
				      stacks  := resp.Stacks
				      stackObject := stacks[0]
				      //fmt.Println("STATE IS :", *stackObject.StackStatus)
				      // fmt.Println(*stackObject.StackStatus)
				      // fmt.Println(stackObject.Outputs)

							userPoolDomain := ""
					 		authenticatedRole := ""
					 		unAuthenticatedRole := ""

							 for _, outputIface := range stackObject.Outputs {
															 // tmpOutput := outputIface.(map[string]string)
															 if (*outputIface.OutputKey == "AppS3Assets") {
																			 // s3Assets = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "IdentityPoolId") {
																			 identityPoolId = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "UserPoolId") {
																			 userPoolId = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "UserPoolClientId") {
																			 userPoolClientId = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "AdminUsername") {
																			 adminUsername = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "AdminPhoneNumber") {
																			 adminPhoneNumber = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "CognitoAdminGroup") {
																			 cognitoAdminPool = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "CognitoUserGroup") {
																			 // cognitoUserPool = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "Stage") {
																			 stage = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "CognitoRegion") {
																			 region = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "UserPoolDomain") {
																			 userPoolDomain = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "CognitoUnAuthenticatedRoleArn") {
																			 unAuthenticatedRole = *outputIface.OutputValue
															 } else if (*outputIface.OutputKey == "CognitoAuthenticatedRoleArn") {
																			 authenticatedRole = *outputIface.OutputValue
															 }
							 }
			*/

		}
	}

	r := events.APIGatewayProxyResponse{
		Body: string(fmt.Sprintf("{ \"message\": \"Could not proceed with bootstrap...\" }")),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: 500,
	}
	return r, nil
}
