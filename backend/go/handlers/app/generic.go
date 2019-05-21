package main

import (
	// "context"
	// "github.com/aws/aws-lambda-go/lambda"
	resolvers "github.com/sbstjn/appsync-resolvers"
	// 	 "reflect"
	"fmt"
	"os"
	// "plugin"
	"github.com/pkg/errors"
	"strings"
	"time"
	// "github.com/satori/go.uuid"
	"github.com/acenteracms/acenteralib"

	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// ID Is always a String
// https://docs.aws.amazon.com/appsync/latest/devguide/scalars.html
var (
	SELF = "SELF"
)

type createInputEvent struct {
	Input  map[string]interface{} `json:"input"`
	Parent map[string]interface{} `json:"parent,omitempty"`
	Type   string                 `json:"type,omitempty"`
}

type deleteInput struct {
	Id string `json:"id"`
}

type deleteInputEvent struct {
	Input deleteInput `json:"input"`
	Type  string      `json:"type,omitempty"`
}

type updateInputEvent struct {
	Input map[string]interface{} `json:"input"`
	Type  string                 `json:"type,omitempty"`
}

type listPaginatedGenericInput struct {
	Input       map[string]interface{} `json:"input"`
	ParentType  string                 `json:"parent_type,omitempty"`
	ParentValue string                 `json:"parent_value,omitempty"`
	Parent      string                 `json:"parent,omitempty"`
	Type        string                 `json:"type,omitempty"`
	// Type					string				`json:"type"`
	Limit     int    `json:"limit"`
	NextToken string `json:"nextToken"`
}

type PaginatedGeneric struct {
	Items     []*map[string]interface{} `json:"items"`
	Type      string                    `json:"type,omitempty"`
	NextToken string                    `json:"nextToken"`
}

/*
type Task struct {
				Id							string		`json:"id"`
        Title						string		`json:"title"`
        Status					string		`json:"status"`
				Description			string		`json:"description"`
				Completed				string		`json:"completed"`
        Created					string		`json:"created"`
}
*/

func (p GenericHandler) handleDeleteGeneric(reqObj *acenteralib.RequestObject, identity map[string]interface{}, event deleteInputEvent) (*map[string]interface{}, error) {
	// TODO: Add Security using the reqObj.User / Roles ... or the reqObj.Session

	elementType := p.ElementType
	if elementType == "" {
		elementType = event.Type
	}

	inactiveGSK := Awslambda.GenerateSortKeyWithTS(fmt.Sprintf("deleted#%s", elementType), "#")
	activeOldSortKey := fmt.Sprintf("active#%s#", elementType)

	output, err := acenteralib.DynamoDB.UpdateItem(&dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(event.Input.Id),
			},
			"sk": &dynamodb.AttributeValue{
				S: aws.String(elementType),
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

	var item map[string]interface{}
	if err == nil {
		// fmt.Println(out.Attributes)
		err = dynamodbattribute.UnmarshalMap(output.Attributes, &item)
	}
	// fmt.Println("Would returm item of ")
	// fmt.Println(item)
	return &item, err
}

func (p GenericHandler) handleUpdateGeneric(reqObj *acenteralib.RequestObject, identity map[string]interface{}, event updateInputEvent) (*map[string]interface{}, error) {

	// TODO: Add Security using the reqObj.User / Roles ... or the reqObj.Session
	timeUtc := time.Now().UTC()
	CreatedTime := timeUtc.Format(time.RFC3339)

	elementType := p.ElementType
	if event.Type != "" {
		elementType = event.Type
	}
	sortKey := elementType
	// Ok receieved update but also a specific SortKey
	if val, ok := event.Input["sk"]; ok {
		sortKey = val.(string)
	}

	// keyCond represents the Key Condition Expression
	// pKeyCond := expression.Key("id").Equal(expression.Value(event.Input["id"])) // "someValue"))
	keyCond := expression.KeyAnd(expression.Key("id").Equal(expression.Value(strings.ToLower(event.Input["id"].(string)))), expression.Key("sk").Equal(expression.Value(sortKey))) // event.Input["sk"])) // SELF or type?
	fmt.Println("UPDATE ITEM HERRE USING val id :", strings.ToLower(event.Input["id"].(string)))

	fmt.Println("UPDATE ITEM HERRE USING val sk :", sortKey)
	fmt.Println(event)
	fmt.Println(keyCond)

	// proj represents the Projection Expression
	proj := expression.ProjectionBuilder{} // NamesList() //expression.Name("aName"), expression.Name("anotherName"), expression.Name("oneOtherName"))
	update := expression.UpdateBuilder{}
	hasUpdated := 0

	hasSk := 0
	for k, v := range event.Input {
		lowerrKey := strings.ToLower(k)
		// Key cannot be part of the Update Statement
		if lowerrKey == "id" || lowerrKey == "sk" {
		} else {
			nameLower := expression.Name(lowerrKey)
			value := expression.Value(v)
			if lowerrKey == "updated" {
				value = expression.Value(CreatedTime)
				hasUpdated = 1
			}
			proj = expression.AddNames(proj, nameLower)
			update = update.Set(nameLower, value)
			fmt.Println("111 - Add of name of and val:", proj, nameLower, value)
		}
	}
	if hasSk <= -1 {
		fmt.Println("NO SK..")
		nameLower := expression.Name("sk")
		value := expression.Value(sortKey)
		proj = expression.AddNames(proj, nameLower)
		update = update.Set(nameLower, value)
	}
	if hasUpdated == 0 {
		k := "updated"
		nameLower := expression.Name(strings.ToLower(k))
		value := expression.Value(CreatedTime)
		proj = expression.AddNames(proj, nameLower)
		update = update.Set(nameLower, value)
	}
	// TODO: Add UserId update ?

	//  condUpdateCond := expression.And(  (expression.Name("id").Equal(expression.Value(strings.ToLower(event.Input["id"].(string))))), (expression.Name("sk").Equal(expression.Value(sortKey))))
	//  condUpdateCond := expression.And(expression.AttributeExists(expression.Name("id")), expression.AttributetExists(expression.Name("sk"))) // event.Input["sk"])) // sortKey or type?

	// builder := expression.NewBuilder().WithCondition(condUpdateCond).WithUpdate(update).WithProjection(proj) // ??

	builder := expression.NewBuilder().WithUpdate(update)
	exp, err := builder.Build()

	dynamoDdbUpdateInput := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(event.Input["id"].(string)),
			},
			"sk": &dynamodb.AttributeValue{
				S: aws.String(sortKey),
			},
		},
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		UpdateExpression:          exp.Update(),
		ConditionExpression:       exp.Condition(),
		ReturnValues:              aws.String("ALL_NEW"),
		TableName:                 aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	}

	output, err := acenteralib.DynamoDB.UpdateItem(dynamoDdbUpdateInput)
	if err != nil {
		fmt.Println(err)
		if ae, ok := err.(awserr.RequestFailure); ok && ae.Code() == "ConditionalCheckFailedException" {
			//fmt.Println("Update failed due to constraint...")
			return nil, errors.WithStack(err)
		}
		return nil, errors.WithStack(err)
	}

	var item map[string]interface{}
	if err == nil {
		err = dynamodbattribute.UnmarshalMap(output.Attributes, &item)
		fmt.Println("Would returm item of ")
		fmt.Println(item)
	}
	return &item, err

}

func (p GenericHandler) handleListAll(reqObj *acenteralib.RequestObject, identity map[string]interface{}) (*PaginatedGeneric, error) {
	// TODO ADD SECURITY

	// Element Type ?
	queryInput := listPaginatedGenericInput{
		Limit: 100,
	}

	output := &PaginatedGeneric{
		Items:     make([]*map[string]interface{}, 0),
		NextToken: "",
	}
	lastResult := &PaginatedGeneric{
		NextToken: "",
	}
	var err error
	for {
		queryInput.NextToken = lastResult.NextToken
		lastResult, err = p.handleList(reqObj, identity, queryInput)
		lastResult.NextToken = lastResult.NextToken
		output.Items = append(output.Items, lastResult.Items...)
		if lastResult.NextToken == "" {
			break
		}
	}
	return output, err
}

func (p GenericHandler) handleList(reqObj *acenteralib.RequestObject, identity map[string]interface{}, input listPaginatedGenericInput) (*PaginatedGeneric, error) {

	elementType := p.ElementType
	if elementType == "" {
		elementType = input.Type
	}

	// TODO: Add Security using the reqObj.User / Roles ... or the reqObj.Session
	if reqObj.User != nil {
		fmt.Println(reqObj.User.Username)
		fmt.Println(reqObj.User.Roles)
	}
	limits := int64(20)
	if input.Limit >= 1 && input.Limit <= 100 {
		limits = int64(input.Limit)
	}
	// no need : site := os.Getenv("SITE")

	// TODO: Get this from end-user session?
	parent := ""
	if val, ok := input.Input["parent"]; ok {
		//do something here
		parent = val.(string)
	}

	if parent == "" {
		parent = os.Getenv("SITE")
	}

	query := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(parent),
			},
			":sk": {
				S: aws.String(fmt.Sprintf("active#%s#", elementType)),
			},
		},
		IndexName:              aws.String("gsi-data-index"),
		KeyConditionExpression: aws.String("gpk = :pk and begins_with(gsk, :sk)"),
		// FilterExpression: aws.String("title = :title"),
		Limit: aws.Int64(limits),
		// TODO: last settings
		// TaskionExpression:   aws.String("id, sk, name, jsondata"),
		// all ... for now TaskionExpression:   aws.String("SongTitle"),
		TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	}
	fmt.Println("WILL Query using :", parent, " and active#", elementType, "# using gsi-data-index")
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
	output := PaginatedGeneric{}
	items := make([]*map[string]interface{}, 0)
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

func (p GenericHandler) handleCreateGeneric(reqObj *acenteralib.RequestObject, identity map[string]interface{}, mutation createInputEvent) (*map[string]interface{}, error) {

	// TODO: Create if user is Admin

	elementType := p.ElementType
	if elementType == "" {
		elementType = mutation.Type
	}

	taskName := mutation.Input["title"].(string)
	site := os.Getenv("SITE")
	fmt.Println("Will be Creating ", elementType, " .. title :", taskName)

	parent := "" // os.Getenv("SITE")
	if mutation.Parent != nil {
		fmt.Println("1 - Will be Creating ", elementType, " .. title :", taskName)
		if val, ok := mutation.Parent["value"]; ok {
			parent = val.(string)
		}
	}
	fmt.Println("2- Will be Creating ", elementType, " .. title :", taskName)
	var dynamoPutInput dynamodb.PutItemInput
	if parent == "" {
		dynamoPutInput = *Awslambda.CreateAppItemParent(taskName, elementType, site, "") // leave default sitei and sort key of active#ELEMENTTYPE#TS
	} else {
		dynamoPutInput = *Awslambda.CreateAppItemParentAndPlugin(taskName, elementType, parent, "", site, "") // leave default sitei and sort key of active#ElementTypetTT#TS
	}
	fmt.Println("3- Will be Creating ", elementType, " .. title :", taskName)

	// keyCond represents the Key Condition Expression
	// pKeyCond := expression.Key("id").Equal(expression.Value(event.Input["id"])) // "someValue"))
	// keyCond := expression.KeyAnd(expression.Key("id").Equal(expression.Value(strings.ToLower(mutation.Input["id"].(string)))), expression.Key("sk").Equal(expression.Value(sortKey))) // event.Input["sk"])) // sortKey or type?

	// proj represents the Projection Expression
	fmt.Println("4- Will be Creating ", elementType, " .. title :", taskName)
	proj := expression.ProjectionBuilder{} // NamesList() //expression.Name("aName"), expression.Name("anotherName"), expression.Name("oneOtherName"))
	update := expression.UpdateBuilder{}

	fmt.Println("5- Will be Creating ", elementType, " .. title :", taskName)
	attributes := map[string]*dynamodb.AttributeValue{}
	/*
		  ":e_mail": &dynamodb.AttributeValue{
		        S: &user.Email,
		    },
			}
	*/

	/*
		for k, v := range m {
			rt := reflect.TypeOf(v)
			switch rt.Kind() {
			case reflect.Slice:
				fmt.Println(k, "is a slice with element type", rt.Elem())
			case reflect.Array:
				fmt.Println(k, "is an array with element type", rt.Elem())
			case reflect.String:
				fmt.Println(k, "is an strirng with element type")
			case reflect.Map:
				fmt.Println(k, "is an map with element type")
			case reflect.Bool:
				fmt.Println(k, "is an BOOL with element type")
			default:
				fmt.Println(k, "is something else entirely", rt.Kind())
			}
		}
	*/

	fmt.Println("6- Will be Creating ", elementType, " .. title :", taskName)

	// This include the Generated dUUID?
	/*uuid := ""
	sk := sortKey
	for k, attr := range dynamoPutInput.Item {
		if (k == "id") {
			uuid = *attr.S
		}
		if (k == "sk") {
			sk = *attr.S
		}
		attributes[k] = attr
	}
	*/

	// Override any values including sk if user wants it ...
	for k, v := range mutation.Input {
		lowerNameKey := strings.ToLower(k)
		nameLower := expression.Name(lowerNameKey)
		fmt.Println("7- Will be Creating ", elementType, " .. title :", taskName)
		value := expression.Value(v)

		// objType := 0
		// rt := reflect.TypeOf(v)

		// var attrType dynamodb.AttributeValue

		attrB, err := dynamodbattribute.Marshal(v)
		// fmt.Println("GOT ERR - 1")
		if err != nil {
			fmt.Println(err)
		}
		attributes[lowerNameKey] = attrB
		dynamoPutInput.Item[lowerNameKey] = attrB
		/*&dynamodb.AttributeValue{
			SS: v.([]string)value.Value,
		}*/

		/*
			switch rt.Kind() {
			case reflect.Slice:
				// SS
				objType = 1
				// fmt.Println(k, "is a slice with element type", rt.Elem())
				if (strings.HasPrefix(rt.Elem().String(), "string")) {
					attrB, err := dynamodbattribute.Marshal(v)
					fmt.Println("GOT ERR - 1")
					fmt.Prinrtln(err)
					attributes[lowerNameKey] = attributes
					&dynamodb.AttributeValue{
						SS:??
					}
					// attrType = dynamodb.AttributeValue.SS
				} else if (strings.HasPrefix(rt.Elem().String(), "byte")) {
					attrType = dynamodb.AttributeValue.BS
				} else if (strings.HasPrefix(rt.Elem().String(), "bool") || strings.HasPrefix(rt.Elem().String(), "BOOL")) {
					// ??? attrType =
					attrType = dynamodb.AttributeValue.L
				} else {
					attrType = dynamodb.AttributeValue.NS
				}
			case reflect.Array:
				objType = 1
				// fmt.Println(k, "is an array with element type", rt.Elem())
				if (strings.HasPrefix(rt.Elem().String(), "string")) {
					attrType = dynamodb.AttributeValue.SS
				} else if (strings.HasPrefix(rt.Elem().String(), "byte")) {
					attrType = BS
				} else if (strings.HasPrefix(rt.Elem().String(), "bool") || strings.HasPrefix(rt.Elem().String(), "BOOL")) {
					// ??? attrType = dynamodb.AttributeValue.NB
				} else {
					attrType = NS
				}
			case reflect.String:
				// fmt.Println(k, "is an strirng with element type")
				objType = 2
				attrType = S
			case reflect.Map:
				// fmt.Println(k, "is an map with element type")
				objType = 3
				attrType = M
			case reflect.Bool:
				objType = 4
				// fmt.Println(k, "is an BOOL with element type")
				attrType = BOOL
			case reflect.Byte:
				objType = 4
				// fmt.Println(k, "is an BOOL with element type")
				attrType = B
			default:
				fmt.Println(k, "is something else entirely", rt.Kind())
				if (strings.HasPrefix(rt.Kind(),"int") || strings.HasPrefix(rt.Kind(),"float")) {
					objType = 5
				} else {
					objType = 2
				}
			}
		*/
		proj = expression.AddNames(proj, nameLower)
		update = update.Set(nameLower, value)
	}

	fmt.Println("8- Will be Creating ", elementType, " .. title :", taskName)

	condUpdateCond := expression.And(expression.AttributeNotExists(expression.Name("id")), expression.AttributeNotExists(expression.Name("sk"))) // event.Input["sk"])) // sortKey or type?
	// builder := expression.NewBuilder().WithKeyCondition(keyCond).WithProjection(proj).WithCondition(condUpdateCond).WithUpdate(update)
	// builder := expression.NewBuilder().WithKeyCondition(keyCond).WithCondition(condUpdateCond).WithUpdate(update)
	fmt.Println("9- Will be Creating ", elementType, " .. title :", taskName)
	fmt.Println("Condupdate is ")
	fmt.Println(condUpdateCond)
	fmt.Println("Update is ")
	fmt.Println(update)
	builder := expression.NewBuilder().WithCondition(condUpdateCond).WithUpdate(update)
	exp, err := builder.Build()
	// update := expression.UpdateBuilder{}
	fmt.Println(err)
	fmt.Println(exp)

	/*
		if (uuid == "") {
			uuid = UUIDGen().String()
		}
	*/
	/*
		dynamoDdbUpdateInput := &dynamodb.UpdateItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"id": &dynamodb.AttributeValue{
					S: aws.String(uuid), // mutation.Input["id"].(string)),
				},
				"sk": &dynamodb.AttributeValue{
					S: aws.String(sk),
				},
			},
			// Key:  										 exp.KeyCondition(),
			ExpressionAttributeNames:  exp.Names(),
			ExpressionAttributeValues: exp.Values(),
			UpdateExpression:          exp.Update(),
			ConditionExpression:  		 exp.Condition(),
			ReturnValues: aws.String("ALL_NEW"),
			TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
		}
	*/
	// output, errPut := acenteralib.DynamoDB.UpdateItem(dynamoDdbUpdateInput)

	fmt.Println("Put item using.....")
	// fmt.Println(dynamoDdbUpdateInput)
	fmt.Println(dynamoPutInput)

	output, errPut := acenteralib.DynamoDB.PutItem(&dynamoPutInput)

	fmt.Println("Outpout is :")
	fmt.Println(output)

	var item map[string]interface{}
	fmt.Println("ERRR IS")
	fmt.Println(errPut)
	if errPut == nil {
		// fmt.Println(out.Attributes)
		errPut = dynamodbattribute.UnmarshalMap(dynamoPutInput.Item, &item) // output.Attributes, &item) // dynamoPutInput.Item, &item)
	}
	// fmt.Println("Would returm item of ")
	fmt.Println(item)
	fmt.Println("ERR put ..")
	fmt.Println(errPut)
	return &item, errPut
}

type GenericHandler struct {
	ElementType string
}

func (p GenericHandler) InitializeRoutes(r resolvers.Repository) error {
	var err error

	if p.ElementType == "" {
		err = r.Add(fmt.Sprintf("query.list"), p.handleList)
		err = r.Add(fmt.Sprintf("query.listAll"), p.handleListAll)
		err = r.Add(fmt.Sprintf("mutation.create"), p.handleCreateGeneric)
		err = r.Add(fmt.Sprintf("mutation.update"), p.handleUpdateGeneric)
		err = r.Add(fmt.Sprintf("mutation.delete"), p.handleDeleteGeneric)
	} else {
		ActionWord := strings.Title(strings.ToLower(p.ElementType))
		err = r.Add(fmt.Sprintf("query.list%s", ActionWord), p.handleList)
		err = r.Add(fmt.Sprintf("query.list%ss", ActionWord), p.handleList) // plurial?
		err = r.Add(fmt.Sprintf("query.listAll%s", ActionWord), p.handleListAll)
		err = r.Add(fmt.Sprintf("query.listAll%ss", ActionWord), p.handleListAll) // plurial
		err = r.Add(fmt.Sprintf("mutation.create%s", ActionWord), p.handleCreateGeneric)
		err = r.Add(fmt.Sprintf("mutation.update%s", ActionWord), p.handleUpdateGeneric)
		err = r.Add(fmt.Sprintf("mutation.delete%s", ActionWord), p.handleDeleteGeneric)
	}
	return err
}
