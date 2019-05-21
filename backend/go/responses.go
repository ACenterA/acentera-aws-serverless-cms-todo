package gofaas

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
)

var (
	responseEmpty = events.APIGatewayProxyResponse{}
	response404   = events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("404\n"),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		StatusCode: 404,
	}
	Response404 = events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("404\n"),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		StatusCode: 404,
	}
	response301 = events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("redirect\n"),
		Headers: map[string]string{
			"Location":     "/",
			"Content-Type": "text/html",
		},
		StatusCode: 301,
	}
	responseUnAuthenticated = events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("{%q: %q}\n", "status", "UnAuthorized"),
		Headers: map[string]string{
			"Content-Type": "application/javascript",
		},
		StatusCode: 401,
	}
	responseEmptyCreated = events.APIGatewayProxyResponse{
		Body: "",
		Headers: map[string]string{
			"Content-Type": "application/javascript",
		},
		StatusCode: 201,
	}
	responseError = events.APIGatewayProxyResponse{
		Body: "",
		Headers: map[string]string{
			"Content-Type": "application/javascript",
		},
		StatusCode: 403,
	}
	responseDeleted = events.APIGatewayProxyResponse{
		Body: "",
		Headers: map[string]string{
			"Content-Type": "application/javascript",
		},
		StatusCode: 204,
	}
)

type PluginSettings struct {
	AMIProjects []PluginInfo `json:"AMIProjects"`
	Clusters    []PluginInfo `json:"Clusters"`
}

// ResponseError is an error type that indicates a non-200 response
type ResponseError struct {
	Body       string
	StatusCode int
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("%s (%d)", e.Body, e.StatusCode)
}

// Response returns an API Gateway Response event
func (e ResponseError) Response() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("{%q: %q}\n", "error", e.Body),
		StatusCode: e.StatusCode,
	}, nil
}

func RestResponseError(err string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("{%q: %q}\n", "message", err),
		Headers: map[string]string{
			"Content-Type": "application/javascript",
		},
		StatusCode: 403,
	}, nil
}

func responseNeedMFA(u string, url string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body: fmt.Sprintf("{%q: %q, %q: %q, %q: %q}\n", "status", "NeedMFA", "code", u, "url", url),
		Headers: map[string]string{
			"Content-Type": "application/javascript",
		},
		StatusCode: 200,
	}
}

func RestResponseSmallCache(u interface{}) (events.APIGatewayProxyResponse, error) {
	//fmt.Println("Reveived:", u)
	b, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return responseEmpty, errors.WithStack(err)
	}

	return events.APIGatewayProxyResponse{
		Body: string(b) + "\n",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Cache-Control": "public, must-revalidate, proxy-revalidate, max-age=30",
		},
		StatusCode: 200,
	}, nil
}

func RestResponseNoCache(u interface{}) (events.APIGatewayProxyResponse, error) {
	b, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return responseEmpty, errors.WithStack(err)
	}

	return events.APIGatewayProxyResponse{
		Body: string(b) + "\n",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Cache-Control": "public, must-revalidate, proxy-revalidate, max-age=0",
		},
		StatusCode: 200,
	}, nil
}

func RestResponseMedCache(u interface{}) (events.APIGatewayProxyResponse, error) {
	b, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return responseEmpty, errors.WithStack(err)
	}

	return events.APIGatewayProxyResponse{
		Body: string(b) + "\n",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Cache-Control": "public, must-revalidate, proxy-revalidate, max-age=8600",
		},
		StatusCode: 200,
	}, nil
}
