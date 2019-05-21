package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	// "github.com/myplugin/gofaas/static"
	// "github.com/myplugin/gofaas"
	"github.com/acenterastatic/static"
	"github.com/aws/aws-lambda-go/events"
	// "gopkg.in/h2non/filetype.v1"
	"encoding/base64"
	"github.com/acenteracms/acenteralib"
)

func init() {
}

func WebsitePublic(sharedlib acenteralib.SharedLib, ctx context.Context, reqObj acenteralib.RequestObject, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// func WebsitePublic(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//fmt.Println("TEST A")
	//fmt.Println("IN GOPLUGIN RECEIEVED OF ", e.Path)
	//fmt.Println(e)
	// TODO: Find if it has /static/xxx ?
	fileNamePath := strings.TrimPrefix(e.Path, "/static/")
	data, err := static.Asset(fmt.Sprintf("dist%v", e.Path))
	fmt.Println("Finding ", fmt.Sprintf("dist%v", e.Path), "")
	if err != nil {
		// Asset was not found.
		// fileNamePath = strings.TrimPrefix(e.Path,"/static/")
		fmt.Println("Finding ", fmt.Sprintf("dist/%v", fileNamePath))
		data, err = static.Asset(fmt.Sprintf("dist/%v", fileNamePath))
		// data, err = static.Asset(fmt.Sprintf("dist/%v", strings.TrimPrefix(e.Path,"/static/")))
		// data, err = static.Asset(e.Path)
	}
	/*
		if (strings.Contains(e.Path,"/plugins/")) {
			return events.APIGatewayProxyResponse{
				Body: string(data),
				Headers: map[string]string{
					"Content-Type": "application/javascript",
					"Cache-Control": "public, max-age=600",
				},
				StatusCode: 404,
			}, nil
		}
	*/
	if err != nil {

		if strings.HasSuffix(fileNamePath, ".js") || strings.HasSuffix(fileNamePath, ".css") {
			return acenteralib.Response404, nil
		} else {
			fmt.Println("Finding as last resort ", fmt.Sprintf("dist/%v", "index.html"))
			data, err = static.Asset(fmt.Sprintf("dist/%v", "index.html"))
			// return response301, nil
		}

	}
	extension := filepath.Ext(e.Path)
	//fmt.Println("Extension 2:", extension)

	if extension == ".css" {
		return events.APIGatewayProxyResponse{
			Body: string(data),
			Headers: map[string]string{
				"Content-Type":  "text/css",
				"Cache-Control": "public, max-age=8600",
			},
			StatusCode: 200,
		}, nil
	} else if extension == ".js" {
		return events.APIGatewayProxyResponse{
			Body: string(data),
			Headers: map[string]string{
				"Content-Type":  "application/javascript",
				"Cache-Control": "public, max-age=8600",
			},
			StatusCode: 200,
		}, nil
	} else if extension == ".gif" {
		output := base64.StdEncoding.EncodeToString(data)
		return events.APIGatewayProxyResponse{
			Body: output,
			Headers: map[string]string{
				"Content-Type":  "image/gif",
				"Cache-Control": "public, max-age=8600",
			},
			IsBase64Encoded: true,
			StatusCode:      200,
		}, nil
	} else if extension == ".png" {
		//fmt.Println("ffa Encode to string ofs data...")

		// read file content into buffer
		/*
		 				fReader := bufio.NewReader(data)
						buf := make([]byte, len(data))
		 			  fReader.Read(buf)
		*/
		// imgBase64Str := base64.StdEncoding.EncodeToString(buf)
		output := base64.StdEncoding.EncodeToString(data)
		// //fmt.Println("z Encode to string of set return as Base64 Encoded...", output)
		return events.APIGatewayProxyResponse{
			// Body: fmt.Sprintf("data:image/png;base64,%s", output),
			Body: fmt.Sprintf("%v", output),
			Headers: map[string]string{
				"Content-Type":  "image/png",
				"Cache-Control": "public, max-age=8600",
			},
			IsBase64Encoded: true,
			StatusCode:      200,
		}, nil
	} else if extension == ".jpg" || extension == ".jpeg" {
		output := base64.StdEncoding.EncodeToString(data)
		return events.APIGatewayProxyResponse{
			Body: output,
			Headers: map[string]string{
				"Content-Type":  "image/jpeg",
				"Cache-Control": "public, max-age=8600",
			},
			IsBase64Encoded: true,
			StatusCode:      200,
		}, nil
	} else if extension == ".svg" {
		output := base64.StdEncoding.EncodeToString(data)
		return events.APIGatewayProxyResponse{
			Body:            output,
			IsBase64Encoded: true,
			Headers: map[string]string{
				"Content-Type":  "image/svg+xml",
				"Cache-Control": "public, max-age=8600",
			},
			StatusCode: 200,
		}, nil
	} else if extension == ".woff2" {
		output := base64.StdEncoding.EncodeToString(data)
		return events.APIGatewayProxyResponse{
			Body:            output,
			IsBase64Encoded: true,
			Headers: map[string]string{
				"Content-Type":  "font/woff2",
				"Cache-Control": "public, max-age=8600",
			},
			StatusCode: 200,
		}, nil
	} else if extension == ".html" || extension == "" {
		// TODO: Detect language ??
		// %%TITLE%%
		strTitle := os.Getenv("TITLE")
		dataWithTitle := strings.Replace(string(data), "%%TITLE%%", strTitle, -1) //TODO: Customize name using CloudFormation Title
		if strTitle == "" {
			dataWithTitle = strings.Replace(string(data), "%%TITLE%%", "Serverless Portal", -1) //TODO: Customize name using CloudFormation Title
		}
		return events.APIGatewayProxyResponse{
			Body: dataWithTitle,
			Headers: map[string]string{
				"Content-Type":  "text/html",
				"Cache-Control": "public, must-revalidate, proxy-revalidate, max-age=0",
			},
			StatusCode: 200,
		}, nil
		// //fmt.Printf("File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
		return acenteralib.Response404, nil
	}
	return acenteralib.Response404, nil
}
