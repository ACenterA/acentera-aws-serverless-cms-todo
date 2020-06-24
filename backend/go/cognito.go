package gofaas

import (
	"os"
	"fmt"
	"context"
	"reflect"
	"strings"

	cfn "github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

func init() {
}

//func(ctx context.Context, event Event) (reason string, err error) {
func CognitoMFASettings(ctx context.Context, event cfn.Event) (reason string, err error) { // (events.APIGatewayProxyResponse, error) { //} )(physicalResourceID string, data map[string]interface{}, err error) {
	//	cfn.CustomResourceFunction
	// fmt.Println("CognitoClientSettings")

	r := cfn.NewResponse(&event)
	r.Status = cfn.StatusSuccess
	if event.RequestType == "Delete" {
		// Always success
		r.Status = cfn.StatusSuccess
	} else if (event.RequestType == "Create" || event.RequestType == "Update") {
		mfa := "ON"
		//TODO: Use event.ResourceProperties["MFA"] ??
		if (os.Getenv("ENFORCE_MFA") == "true") {
		  mfa = "ON"
		} else {
		  mfa = "OPTIONAL"
		}

		csoftwareTokenMfa := &cognitoidentityprovider.SoftwareTokenMfaConfigType{
			Enabled: aws.Bool((mfa == "ON")),
		}
		_, errorTmp := CognitoIdentityServiceProvider.SetUserPoolMfaConfig(&cognitoidentityprovider.SetUserPoolMfaConfigInput{
			MfaConfiguration:              aws.String(mfa),
			SoftwareTokenMfaConfiguration: csoftwareTokenMfa,
			UserPoolId:                    aws.String(event.ResourceProperties["UserPoolId"].(string)),
		})
		if errorTmp != nil {
			// fmt.Println("Set MFA Config error:", errorTmp)
			r.Status = cfn.StatusFailed
		} else {
			r.Status = cfn.StatusSuccess
		}
	}

	// r.PhysicalResourceID, r.Data, err = lambdaFunction(ctx, event)
	if r.PhysicalResourceID == "" {
		r.PhysicalResourceID = "ACenterACognitoUserPoolDomain"
	}

	if err != nil {
		r.Status = cfn.StatusFailed
		r.Reason = err.Error()
		// fmt.Printf("sending status failed: %s", r.Reason)
	} else {
		r.Status = cfn.StatusSuccess
	}

	err = r.Send() // (client)
	if err != nil {
		reason = err.Error()
		// fmt.Println("Reason error: ", reason)
	} else {
		// fmt.Println("Successfully generated value sent response...")
		reason = "Successfully generated value..."
	}

	return reason, err

}

//func(ctx context.Context, event Event) (reason string, err error) {
func CognitoClientSettings(ctx context.Context, event cfn.Event) (reason string, err error) { // (events.APIGatewayProxyResponse, error) { //} )(physicalResourceID string, data map[string]interface{}, err error) {
	//	cfn.CustomResourceFunction
	// fmt.Println("CognitoClientSettings")

	r := cfn.NewResponse(&event)

	if event.RequestType == "Delete" {
		// Always success
		r.Status = cfn.StatusSuccess
	} else {
		callbackUrls := make([]*string, 0)
		if val, ok := event.ResourceProperties["CallbackURL"]; ok {
			rt := reflect.TypeOf(val)
			switch rt.Kind() {
			case reflect.Slice:
				// its a slice... // fmt.Println(k, "is a slice with element type", rt.Elem())
				// fmt.Println("GOT VAL HERE", val)
				arr := val.([]interface{})
				for i := 0; i < len(arr); i++ {
					// fmt.Printf("Parsing logoutUrls of : %s\n", arr[i])
					callbackUrls = append(callbackUrls, aws.String(arr[i].(string)))
				}
			case reflect.Array:
				// its an array
				// // fmt.Println(k, "is an array with element type", rt.Elem())
				arr := val.([]string)
				for i := 0; i < len(arr); i++ {
					// fmt.Printf("Parsing logoutUrls of : %s\n", arr[i])
					callbackUrls = append(callbackUrls, aws.String(arr[i]))
				}
			default:
				// Its a simple string
				arr := val.(string)
				// fmt.Printf("got callbackUrl provider of : %s\n", arr)
				callbackUrls = append(callbackUrls, aws.String(arr))
			}
		}

		logoutUrls := make([]*string, 0)
		if val, ok := event.ResourceProperties["LogoutURL"]; ok {

			rt := reflect.TypeOf(val)
			switch rt.Kind() {
			case reflect.Slice:
				// its a slice... // fmt.Println(k, "is a slice with element type", rt.Elem())
				// fmt.Println("GOT VAL HERE", val)
				arr := val.([]interface{})
				for i := 0; i < len(arr); i++ {
					// fmt.Printf("Parsing logoutUrls of : %s\n", arr[i])
					logoutUrls = append(logoutUrls, aws.String(arr[i].(string)))
				}
			case reflect.Array:
				// its an array
				// // fmt.Println(k, "is an array with element type", rt.Elem())
				arr := val.([]string)
				// fmt.Println("GOT VAL HERE A", val)
				for i := 0; i < len(arr); i++ {
					// fmt.Printf("Parsing logoutUrls of : %s\n", arr[i])
					logoutUrls = append(logoutUrls, aws.String(arr[i]))
				}
			default:
				// Its a simple string
				arr := val.(string)
				// fmt.Printf("got logoutUrl provider of : %s\n", arr)
				logoutUrls = append(logoutUrls, aws.String(arr))
			}
		}

		supportedIdentityProviders := make([]*string, 0)
		if val, ok := event.ResourceProperties["SupportedIdentityProviders"]; ok {
			rt := reflect.TypeOf(val)
			switch rt.Kind() {
			case reflect.Slice:
				// its a slice... // fmt.Println(k, "is a slice with element type", rt.Elem())
				// fmt.Println("Z GOT VAL HERE", val)
				arr := val.([]interface{})
				// fmt.Println(arr)
				for i := 0; i < len(arr); i++ {
					// fmt.Printf("Parsing supportedIdentityProvider of : %s\n", arr[i])
					supportedIdentityProviders = append(supportedIdentityProviders, aws.String(arr[i].(string)))
				}
			case reflect.Array:
				// its an array
				// // fmt.Println(k, "is an array with element type", rt.Elem())
				arr := val.([]string)
				for i := 0; i < len(arr); i++ {
					// fmt.Printf("Parsing supportedIdentityProvider of : %s\n", arr[i])
					supportedIdentityProviders = append(supportedIdentityProviders, aws.String(arr[i]))
				}
			default:
				// Its a simple string
				arr := val.(string)
				// fmt.Printf("got supported provider of : %s\n", arr)
				supportedIdentityProviders = append(supportedIdentityProviders, aws.String(arr))
			}
		}

		allowedOAuthFlows := make([]*string, 0)
		if val, ok := event.ResourceProperties["AllowedOAuthFlows"]; ok {
			rt := reflect.TypeOf(val)
			switch rt.Kind() {
			case reflect.Slice:
				// its a slice... // fmt.Println(k, "is a slice with element type", rt.Elem())
				// fmt.Println("W GOT VAL HERE", val)
				arr := val.([]interface{})
				for i := 0; i < len(arr); i++ {
					// fmt.Printf("Parsing oauthflows of : %s\n", arr[i])
					allowedOAuthFlows = append(allowedOAuthFlows, aws.String(arr[i].(string)))
				}
			case reflect.Array:
				// its an array
				// // fmt.Println(k, "is an array with element type", rt.Elem())
				arr := val.([]string)
				for i := 0; i < len(arr); i++ {
					// fmt.Printf("Parsing oauthflows of : %s\n", arr)
					allowedOAuthFlows = append(allowedOAuthFlows, aws.String(arr[i]))
				}
			default:
				// Its a simple string
				arr := val.(string)
				// fmt.Printf("got oauthflows of : %s\n", arr)
				allowedOAuthFlows = append(allowedOAuthFlows, aws.String(arr))
			}

		}

		allowedOAuthScopes := make([]*string, 0)
		if val, ok := event.ResourceProperties["AllowedOAuthScopes"]; ok {
			rt := reflect.TypeOf(val)
			switch rt.Kind() {
			case reflect.Slice:
				// its a slice... // fmt.Println(k, "is a slice with element type", rt.Elem())
				arr := val.([]interface{})
				for i := 0; i < len(arr); i++ {
					// fmt.Printf("Parsing scope of : %s\n", arr[i])
					allowedOAuthScopes = append(allowedOAuthScopes, aws.String(arr[i].(string)))
				}
			case reflect.Array:
				// its an array
				// // fmt.Println(k, "is an array with element type", rt.Elem())
				arr := val.([]string)
				for i := 0; i < len(arr); i++ {
					// fmt.Printf("Parsing scope of : %s\n", arr[i])
					allowedOAuthScopes = append(allowedOAuthScopes, aws.String(arr[i]))
				}
			default:
				// Its a simple string
				arr := val.(string)
				// fmt.Printf("got oauthFScopes  of : %s\n", arr)
				allowedOAuthScopes = append(allowedOAuthScopes, aws.String(arr))
			}
		}

		ExplicitAuthFlows := make([]*string, 0)
		if val, ok := event.ResourceProperties["ExplicitAuthFlows"]; ok {
			rt := reflect.TypeOf(val)
			switch rt.Kind() {
			case reflect.Slice:
				// its a slice... // fmt.Println(k, "is a slice with element type", rt.Elem())
				arr := val.([]interface{})
				for i := 0; i < len(arr); i++ {
					fmt.Printf("Parsing explicitAuthFlow of : %s\n", arr[i])
					ExplicitAuthFlows = append(ExplicitAuthFlows, aws.String(arr[i].(string)))
				}
			case reflect.Array:
				// its an array
				// // fmt.Println(k, "is an array with element type", rt.Elem())
				arr := val.([]string)
				for i := 0; i < len(arr); i++ {
					fmt.Printf("Parsing explicitAuthFlow of : %s\n", arr[i])
					ExplicitAuthFlows = append(ExplicitAuthFlows, aws.String(arr[i]))
				}
			default:
				// Its a simple string
				arr := val.(string)
				fmt.Printf("got ExplicitAuthFlows  of : %s\n", arr)
				ExplicitAuthFlows = append(ExplicitAuthFlows, aws.String(arr))
			}
		}


		allowedOAuthFlowClient := false
		if val, ok := event.ResourceProperties["AllowedOAuthFlowsUserPoolClient"]; ok {
			arr := val.(string)
			if arr == "true" {
				allowedOAuthFlowClient = true
			}
		}

		_, errz := CognitoIdentityServiceProvider.UpdateUserPoolClientWithContext(ctx, &cognitoidentityprovider.UpdateUserPoolClientInput{
			UserPoolId:                      aws.String(event.ResourceProperties["UserPoolId"].(string)),
			ClientId:                        aws.String(event.ResourceProperties["UserPoolClientId"].(string)),
			SupportedIdentityProviders:      supportedIdentityProviders,
			CallbackURLs:                    callbackUrls,
			LogoutURLs:                      logoutUrls,
			AllowedOAuthFlowsUserPoolClient: aws.Bool(allowedOAuthFlowClient),
			AllowedOAuthFlows:               allowedOAuthFlows,
			AllowedOAuthScopes:              allowedOAuthScopes,
			ExplicitAuthFlows:               ExplicitAuthFlows,
		})
		err = errz

		fmt.Println("SET EXPLICIT AUTH OF :", ExplicitAuthFlows)
		fmt.Println("GOT UPDATE USERPOOOL CLIENT ERROR OF :", err)

                mfa := "ON"
                //TODO: Use event.ResourceProperties["MFA"] ??
                if (os.Getenv("ENFORCE_MFA") == "true") {
                  mfa = "ON"
                } else {
                  mfa = "OPTIONAL"
                }

                csoftwareTokenMfa := &cognitoidentityprovider.SoftwareTokenMfaConfigType{
                        Enabled: aws.Bool((mfa == "ON")),
                }
                _, errz = CognitoIdentityServiceProvider.SetUserPoolMfaConfig(&cognitoidentityprovider.SetUserPoolMfaConfigInput{
                        MfaConfiguration:              aws.String(mfa),
                        SoftwareTokenMfaConfiguration: csoftwareTokenMfa,
                        UserPoolId:                    aws.String(event.ResourceProperties["UserPoolId"].(string)),
                })

		// fmt.Println(resp)
		/*
		 if r.PhysicalResourceID == "" {
		   r.PhysicalResourceID = "ACenterACognitoUserPoolClientSettings"
		 }
		*/
		if (errz != nil ){
		// r.PhysicalResourceID = *resp.UserPoolClient.ClientId
		  fmt.Println("UserPoolClient got error mfa?", errz)
		}
	}

	// r.PhysicalResourceID, r.Data, err = lambdaFunction(ctx, event)
	if r.PhysicalResourceID == "" {
		r.PhysicalResourceID = "ACenterACognitoUserPoolDomain"
	}

	if err != nil {
		r.Status = cfn.StatusFailed
		r.Reason = err.Error()
		// fmt.Printf("sending status failed: %s", r.Reason)
	} else {
		r.Status = cfn.StatusSuccess
	}

	err = r.Send() // (client)
	if err != nil {
		reason = err.Error()
		// fmt.Println("Reason error: ", reason)
	} else {
		// fmt.Println("Successfully generated value sent response...")
		reason = "Successfully generated value..."
	}

	return reason, err
}

/*
func (c *CognitoIdentityProvider) SetUserPoolMfaConfig(input *SetUserPoolMfaConfigInput) (*SetUserPoolMfaConfigOutput, error)
type SetUserPoolMfaConfigInput struct {

    // The MFA configuration.
    MfaConfiguration *string `type:"string" enum:"UserPoolMfaType"`

    // The SMS text message MFA configuration.
    SmsMfaConfiguration *SmsMfaConfigType `type:"structure"`

    // The software token MFA configuration.
    SoftwareTokenMfaConfiguration *SoftwareTokenMfaConfigType `type:"structure"`

    // The user pool ID.
    //
    // UserPoolId is a required field
    UserPoolId *string `min:"1" type:"string" required:"true"`
    // contains filtered or unexported fields
}
*/

//func(ctx context.Context, event Event) (reason string, err error) {
/*
func CognitoMFA(ctx context.Context, event cfn.Event)  (reason string, err error) {// (events.APIGatewayProxyResponse, error) { //} )(physicalResourceID string, data map[string]interface{}, err error) {

 }
*/

func deleteUserPoolDomain(ctx context.Context, domain string) (err error) {
	resp, errorTmp := CognitoIdentityServiceProvider.DescribeUserPoolDomainWithContext(ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
		Domain: aws.String(domain),
	})

	err = errorTmp
	// fmt.Println("Received response")
	// fmt.Println(resp)
	// fmt.Println(resp.DomainDescription)
	// fmt.Println("Received errr")
	// fmt.Println(err)

	if resp != nil {
		if resp.DomainDescription != nil {
			if resp.DomainDescription.Domain != nil {
				// it exists
				_, err = CognitoIdentityServiceProvider.DeleteUserPoolDomainWithContext(ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
					UserPoolId: resp.DomainDescription.UserPoolId,
					Domain:     aws.String(domain),
				})
			}
		}
	}
	return
}

//func(ctx context.Context, event Event) (reason string, err error) {
func CognitoClientDomains(ctx context.Context, event cfn.Event) (reason string, err error) { // (events.APIGatewayProxyResponse, error) { //} )(physicalResourceID string, data map[string]interface{}, err error) {
	//	cfn.CustomResourceFunction
	fmt.Println("CognitoClientDomains starts ...")
	r := cfn.NewResponse(&event)

	fmt.Println("1 - CognitoClientDomains starts ...")
	if event.RequestType == "Delete" {
		// Always success
		// fmt.Println("CognitoClientDomains Got Delete Domain ...")
	        fmt.Println("2 - CognitoClientDomains starts ...")
		if val, ok := event.ResourceProperties["Domain"]; ok {
			//do something here
			// fmt.Println("CognitoClientDomains Got Domain Yes ...", strings.ToLower(strings.Trim(val.(string), " ")))
			err = deleteUserPoolDomain(ctx, strings.ToLower(strings.Trim(val.(string), " ")))
		}
		// Should always return true ?
		// r.Status = cfn.StatusSuccess
	} else if event.RequestType == "Update" {
		// fmt.Println("CognitoClientDomains Got Update ... ")
		// No checks here, we always have a Domain in OldResourceProperties
	        fmt.Println("3 - CognitoClientDomains starts ...")
		deleteUserPoolDomain(ctx, strings.ToLower(strings.Trim(event.OldResourceProperties["Domain"].(string), " ")))
		// fmt.Println("CognitoClientDomains Got Creating ... ")
		_, errz := CognitoIdentityServiceProvider.CreateUserPoolDomainWithContext(ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			UserPoolId: aws.String(event.ResourceProperties["UserPoolId"].(string)),
			Domain:     aws.String(strings.ToLower(strings.Trim(event.ResourceProperties["Domain"].(string), " "))),
		})
		err = errz
		// fmt.Println("COGNITO RESP")
		// fmt.Println(resp)
		/*
			 if (resp != nil) {
			 	r.PhysicalResourceID = *resp.CloudFrontDomain
			}
		*/
	} else if event.RequestType == "Create" {
	        fmt.Println("4 - CognitoClientDomains starts ...")
		// fmt.Println("CognitoClientDomains Got Create ... ")
		domainTmp := strings.ToLower(strings.Trim(event.ResourceProperties["Domain"].(string), " "))
		// fmt.Println("Recieved creation of domain : ", domainTmp)
		_, errz := CognitoIdentityServiceProvider.CreateUserPoolDomainWithContext(ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
			UserPoolId: aws.String(event.ResourceProperties["UserPoolId"].(string)),
			Domain:     aws.String(domainTmp),
		})
		fmt.Println("CognitoClientDomains Error: ")
		fmt.Println(err)
		err = errz

		// fmt.Println("COGNITO RESP")
		// fmt.Println(resp)
		/*if (resp != nil) {
		 	r.PhysicalResourceID = *resp.CloudFrontDomain
		}*/
		r.PhysicalResourceID = "ACenterACognitoUserPoolDomain"
	}
	fmt.Println("CognitoClientDomains Completed... using physical id of ")

	// r.PhysicalResourceID, r.Data, err = lambdaFunction(ctx, event)
	if r.PhysicalResourceID == "" {
		r.PhysicalResourceID = "ACenterACognitoUserPoolDomain"
	}
	// // fmt.Printf("PhysicalResourceID must exist, copying Log Stream name")
	// r.PhysicalResourceID = lambdacontext.LogStreamName
	//}

	if err != nil {
		r.Status = cfn.StatusFailed
		r.Reason = err.Error()
		fmt.Printf("sending status failed: %s", r.Reason)
	} else {
		r.Status = cfn.StatusSuccess
	}

	err = r.Send() // (client)
	if err != nil {
		reason = err.Error()
		fmt.Println("Reason error: ", reason)
	} else {
		fmt.Println("Successfully generated value sent response...")
		reason = "Successfully generated value..."
	}

	return reason, err
}
