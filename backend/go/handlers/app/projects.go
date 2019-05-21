package main

import (
	// "context"
	// "github.com/aws/aws-lambda-go/events"
	// "github.com/aws/aws-lambda-go/lambda"
	"fmt"
	resolvers "github.com/sbstjn/appsync-resolvers"
	"os"
	//    "plugin"
	// "strconv"
	"github.com/acenteracms/acenteralib"
	"github.com/pkg/errors"
	"time"

	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// UUIDGen is a UUID generator that can be mocked for testing
// ID Is always a String
// https://docs.aws.amazon.com/appsync/latest/devguide/scalars.html
type createProjectInputEvent struct {
	Input createProjectInput `json:"input"`
}

type createProjectInput struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

type deleteProjectInput struct {
	Id string `json:"id"`
}

type deleteProjectInputEvent struct {
	Input deleteProjectInput `json:"input"`
}

type updateProjectInputEvent struct {
	Input Project `json:"input"`
}

type listPaginatedInput struct {
	Limit     int    `json:"limit"`
	NextToken string `json:"nextToken"`
}

type PaginatedProject struct {
	Items     []*Project `json:"items"`
	NextToken string     `json:"nextToken"`
}

type Project struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Created string `json:"created"`
}

func (p ProjectHandler) handleDeleteProject(reqObj *acenteralib.RequestObject, identity map[string]interface{}, event deleteProjectInputEvent) (*Project, error) {
	// TODO: Security

	inactiveGSK := Awslambda.GenerateSortKeyWithTS(fmt.Sprintf("deleted#%s", p.ElementType), "#")
	activeOldSortKey := fmt.Sprintf("active#%s#", p.ElementType)

	output, err := acenteralib.DynamoDB.UpdateItem(&dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(event.Input.Id),
			},
			"sk": &dynamodb.AttributeValue{
				S: aws.String(p.ElementType),
			},
		},
		UpdateExpression: aws.String("SET #status = :status, #gsk = :gsk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": &dynamodb.AttributeValue{
				S: aws.String("deleted"),
			},
			":gsk": &dynamodb.AttributeValue{
				S: aws.String(inactiveGSK),
			},
			":activeInfo": &dynamodb.AttributeValue{
				S: aws.String(activeOldSortKey),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"),
			"#gsk":    aws.String("gsk"),
		},
		ConditionExpression: aws.String("begins_with(gsk, :activeInfo)"),
		ReturnValues:        aws.String("ALL_NEW"),
		TableName:           aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	})

	if err != nil {
		//fmt.Println("ADD PLUGINS IN DB ERR")
		//fmt.Println(err)
		if ae, ok := err.(awserr.RequestFailure); ok && ae.Code() == "ConditionalCheckFailedException" {
			//fmt.Println("Update failed due to constraint...")
			return nil, errors.WithStack(err)
		}
		return nil, errors.WithStack(err)
	}

	item := Project{}
	if err == nil {
		// fmt.Println(out.Attributes)
		err = dynamodbattribute.UnmarshalMap(output.Attributes, &item)
	}
	// fmt.Println("Would returm item of ")
	// fmt.Println(item)
	return &item, err

}

func (p ProjectHandler) handleUpdateProject(reqObj *acenteralib.RequestObject, identity map[string]interface{}, event updateProjectInputEvent) (*Project, error) {
	// TODO: Security

	timeUtc := time.Now().UTC()
	CreatedTime := timeUtc.Format(time.RFC3339)

	/*
		id: ID!
		title: String!
		status: String!
	*/
	output, err := acenteralib.DynamoDB.UpdateItem(&dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(event.Input.Id),
			},
			"sk": &dynamodb.AttributeValue{
				S: aws.String(p.ElementType),
			},
		},
		UpdateExpression: aws.String("SET #status = :status, #title = :title, #updated = :updated"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": &dynamodb.AttributeValue{
				S: aws.String(event.Input.Status),
			},
			":updated": &dynamodb.AttributeValue{
				S: aws.String(CreatedTime),
			},
			":title": &dynamodb.AttributeValue{
				S: aws.String(event.Input.Title),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#status":  aws.String("status"),
			"#title":   aws.String("title"),
			"#updated": aws.String("updated"),
		},
		// no conditions for now ... ConditionExpression: aws.String("begins_with(gsk, :activeInfo)"),
		ReturnValues: aws.String("ALL_NEW"),
		TableName:    aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	})

	if err != nil {
		//fmt.Println("ADD PLUGINS IN DB ERR")
		//fmt.Println(err)
		if ae, ok := err.(awserr.RequestFailure); ok && ae.Code() == "ConditionalCheckFailedException" {
			//fmt.Println("Update failed due to constraint...")
			return nil, errors.WithStack(err)
		}
		return nil, errors.WithStack(err)
	}

	item := Project{}
	if err == nil {
		// fmt.Println(out.Attributes)
		err = dynamodbattribute.UnmarshalMap(output.Attributes, &item)
	}
	// fmt.Println("Would returm item of ")
	// fmt.Println(item)
	return &item, err

}

func (p ProjectHandler) handleListAllProjects(reqObj *acenteralib.RequestObject, identity map[string]interface{}) (*PaginatedProject, error) {
	// TODO ADD SECURITY

	queryInput := listPaginatedInput{
		Limit: 100,
	}

	output := &PaginatedProject{
		Items:     make([]*Project, 0),
		NextToken: "",
	}
	lastResult := &PaginatedProject{
		NextToken: "",
	}
	var err error
	for {
		queryInput.NextToken = lastResult.NextToken
		lastResult, err = p.handleListProjects(reqObj, identity, queryInput)
		lastResult.NextToken = lastResult.NextToken
		output.Items = append(output.Items, lastResult.Items...)
		if lastResult.NextToken == "" {
			break
		}
	}

	return output, err
}

func (p ProjectHandler) handleListProjects(reqObj *acenteralib.RequestObject, identity map[string]interface{}, input listPaginatedInput) (*PaginatedProject, error) {
	// fmt.Println(reqObj)
	fmt.Println("USER IS :")
	fmt.Println(reqObj.User)
	if reqObj.User != nil {
		fmt.Println("NOT NUIL USER IS :")
		fmt.Println(reqObj.User.Username)
		fmt.Println(reqObj.User.Roles)
	}
	limits := int64(20)
	if input.Limit >= 1 && input.Limit <= 100 {
		limits = int64(input.Limit)
	}
	site := os.Getenv("SITE")

	query := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(site),
			},
			":sk": {
				S: aws.String(fmt.Sprintf("active#%s#", p.ElementType)),
			},
		},
		IndexName:              aws.String("gsi-data-index"),
		KeyConditionExpression: aws.String("gpk = :pk and begins_with(gsk, :sk)"),
		// FilterExpression: aws.String("title = :title"),
		Limit: aws.Int64(limits),
		// TODO: last settings
		// ProjectionExpression:   aws.String("id, sk, name, jsondata"),
		// all ... for now ProjectionExpression:   aws.String("SongTitle"),
		TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	}

	var lastK map[string]interface{}
	if input.NextToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(input.NextToken)
		if err != nil {
			// fmt.Println("[ERROR] - Cannot decrypt start Key")
			return nil, err
		}
		// fmt.Println("Decoded to :", string(decoded))

		errT := json.Unmarshal(decoded, &lastK)
		if errT != nil {
			fmt.Println("[ERROR] - DynamoDB decode error:", errT)
			return nil, errT
		}
		newStartKey, errTF := dynamodbattribute.MarshalMap(lastK)
		if errTF != nil {
			fmt.Println("[ERROR] - DynamoDB Get EVAL START KEY", errTF)
			return nil, errTF
		}

		query.ExclusiveStartKey = newStartKey
	}

	out, err := acenteralib.DynamoDB.Query(query)
	output := PaginatedProject{}
	items := make([]*Project, 0)
	if err == nil && *out.Count >= 1 {
		err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &items)
	}
	output.Items = items
	lastEvaluatedKey := ""
	if out.LastEvaluatedKey != nil {
		//map[string]*dynamodb.AttributeValue
		var lastKey map[string]interface{}
		err = dynamodbattribute.UnmarshalMap(out.LastEvaluatedKey, &lastKey) // , lastKey)
		jsonString, _ := json.Marshal(lastKey)
		jsonStringBase64 := base64.StdEncoding.EncodeToString([]byte(jsonString))
		lastEvaluatedKey = jsonStringBase64
	}
	output.NextToken = lastEvaluatedKey

	return &output, err
}

func (p ProjectHandler) handleCreateProject(reqObj *acenteralib.RequestObject, identity map[string]interface{}, mutation createProjectInputEvent) (*Project, error) {

	// TODO: Create if user is Admin
	projectName := mutation.Input.Title
	site := os.Getenv("SITE")

	fmt.Println("SITE IS :", site)
	fmt.Println("BadProject Element Type is :", p.ElementType)
	dynamoPutInput := Awslambda.CreateAppItemParent(projectName, p.ElementType, site, "") // leave default sitei and sort key of active#ELEMENTTYPE#TS

	statusInfo := &dynamodb.AttributeValue{
		S: aws.String(mutation.Input.Status),
	}
	dynamoPutInput.Item["status"] = statusInfo

	_, err := acenteralib.DynamoDB.PutItem(dynamoPutInput)

	item := Project{}
	if err == nil {
		// fmt.Println(out.Attributes)
		err = dynamodbattribute.UnmarshalMap(dynamoPutInput.Item, &item)
	}
	// fmt.Println("Would returm item of ")
	// fmt.Println(item)
	return &item, err
}

type ProjectHandler struct {
	ElementType string
}

func (p ProjectHandler) InitializeRoutes(r resolvers.Repository) error {
	var err error

	err = r.Add("query.listProjects", p.handleListProjects)
	err = r.Add("query.listAllProjects", p.handleListAllProjects)
	// err = r.Add("query.getTask", handleTaskById)
	// err = r.Add("field.task.notes", handleNotes)
	err = r.Add("mutation.createProject", p.handleCreateProject)
	err = r.Add("mutation.updateProject", p.handleUpdateProject)
	err = r.Add("mutation.deleteProject", p.handleDeleteProject)
	return err
}
