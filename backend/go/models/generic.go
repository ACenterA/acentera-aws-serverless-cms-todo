package models

import (
	// "context"
	// "github.com/aws/aws-lambda-go/lambda"
	pluralize "github.com/gertd/go-pluralize"

	"crypto/md5"
	"log"

	resolvers "github.com/sbstjn/appsync-resolvers"
	// 	 "reflect"
	"fmt"
	"os"

	"strings"
	"time"

	"github.com/pkg/errors"

	// "github.com/satori/go.uuid"
	"github.com/acenteracms/acenteralib"

	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"gopkg.in/yaml.v2"
)

// ID Is always a String
// https://docs.aws.amazon.com/appsync/latest/devguide/scalars.html
var (
	SELF         = "SELF"
	CustomModels = make(map[string]Models, 0)
)

type Models struct {
	Version               int64                  `yaml:"version"`
	Admin                 int64                  `yaml:"admin"`
	OneToMany             string                 `yaml:"one_to_many"`
	CustomId              string                 `yaml:"id"`
	OneToManyUpdateFields string                 `yaml:"one_to_many_update_parent_fields"`
	ModelID               string                 `yaml:"model"`
	Generic               int64                  `yaml:"generic"`
	Parent                string                 `yaml:"parent"`
	Plurial               string                 `yaml:"plurial"`
	Singular              string                 `yaml:"singular"`
	ClassName             string                 `yaml:"class"`
	Resolvers             map[string]interface{} `yaml:"resolvers"`
	Class                 *GenericHandler        `yaml:"-"`
}

type conf struct {
	Models map[string]Models `yaml:"models"`
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
	        yamlFile, err = ioutil.ReadFile("./handlers/app/conf.yaml")
	        if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func init() {
}

func (p GenericHandler) Initialize(r resolvers.Repository) error {

	// genericHandler := GenericHandler{Awslambda: p.Awslambda, ElementType: ""}
	// genericHandler.InitializeRoutes(r)

	// projectHandler := ProjectHandler{gofaas.GenericHandler{Awslambda: Awslambda, ElementType: "PROJECT", ActionWord: "PROJECT"}}
	// projectHandler.InitializeRoutes(r)

	var c conf
	yamlObjList := c.getConf()
	fmt.Println(c)
	pluralize := pluralize.NewClient()
	for k, v := range yamlObjList.Models {
		fmt.Println(" K :", k)
		fmt.Println(" V :", v)

		modelType := strings.ToLower(strings.ToLower(pluralize.Singular(k)))
		if v.ModelID != "" {
			modelType = v.ModelID
		}

		ActionWord := v.Plurial
		if v.Plurial == "" {
			ActionWord = strings.Title(strings.ToLower(modelType))
		}
		ActionWordSingular := strings.Title(strings.ToLower(v.Singular))
		if ActionWordSingular == "" {
			ActionWordSingular = strings.Title(strings.ToLower(pluralize.Singular(modelType)))
		}

		ActionWordPlurial := strings.Title(strings.ToLower(v.Plurial))
		if ActionWordPlurial == "" {
			ActionWordPlurial = strings.Title(strings.ToLower(pluralize.Plural(modelType)))
		}

		fmt.Println("INIT OF ROUTES HERE for ", modelType)
		customHandler := GenericHandler{Awslambda: p.Awslambda, Models: v, ElementType: modelType, ActionWord: ActionWord, ActionWordSingular: ActionWordSingular, ActionWordPlurial: ActionWordPlurial}
		customHandler.InitializeRoutes(r)
	}
	/*
		projectHandler := GenericHandler{Awslambda: p.Awslambda, ElementType: "PROJECT"}
		projectHandler.InitializeRoutes(r)

		taskHandler := GenericHandler{Awslambda: p.Awslambda, ElementType: "TASKS"}
		taskHandler.InitializeRoutes(r)

		postsHandler := GenericHandler{Awslambda: p.Awslambda, ElementType: "POSTS", ActionWord: "Posts"}
		postsHandler.InitializeRoutes(r)
	*/
	fmt.Println("INITIALIZEED ROUTES S")

	// postHandler := PostHandler{gofaas.GenericHandler{Awslambda: Awslambda, ElementType: "POSTS", ActionWord: "Posts"}}
	// postHandler.InitializeRoutes(r)
	return nil
}

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
	Input map[string]interface{} `json:"input"`
	//bad	ParentType  string                 `json:"parent_type,omitempty"`
	//bad	ParentValue string                 `json:"parent_value,omitempty"`
	//bad	Parent      string                 `json:"parent,omitempty"`
	Type string `json:"type,omitempty"`
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

func (p GenericHandler) HandleDeleteGeneric(reqObj *acenteralib.RequestObject, identity map[string]interface{}, event deleteInputEvent) (*map[string]interface{}, error) {
	// TODO: Add Security using the reqObj.User / Roles ... or the reqObj.Session

	elementType := p.ElementType
	if elementType == "" {
		elementType = event.Type
	}

	inactiveGSK := p.Awslambda.GenerateSortKeyWithTS(fmt.Sprintf("deleted#%s", elementType), "#")
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
			":activeGskInfo": &dynamodb.AttributeValue{
				S: aws.String(activeOldSortKey),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#status": aws.String("status"),
			"#gsk":    aws.String("gsk"),
		},
		ConditionExpression: aws.String("begins_with(gsk, :activeGskInfo)"),
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

func (p GenericHandler) HandleUpdateGeneric(reqObj *acenteralib.RequestObject, identity map[string]interface{}, event updateInputEvent) (*map[string]interface{}, error) {

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

func (p GenericHandler) HandleListAll(reqObj *acenteralib.RequestObject, identity map[string]interface{}) (*PaginatedGeneric, error) {
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
		lastResult, err = p.HandleList(reqObj, identity, queryInput)
		lastResult.NextToken = lastResult.NextToken
		output.Items = append(output.Items, lastResult.Items...)
		if lastResult.NextToken == "" {
			break
		}
	}
	return output, err
}

func (p GenericHandler) HandleListAdmin(reqObj *acenteralib.RequestObject, identity map[string]interface{}, input listPaginatedGenericInput) (*PaginatedGeneric, error) {
	return p.handleList(false, reqObj, identity, input)
}
func (p GenericHandler) HandleList(reqObj *acenteralib.RequestObject, identity map[string]interface{}, input listPaginatedGenericInput) (*PaginatedGeneric, error) {
	return p.handleList(false, reqObj, identity, input)
}
func (p GenericHandler) handleList(admin bool, reqObj *acenteralib.RequestObject, identity map[string]interface{}, input listPaginatedGenericInput) (*PaginatedGeneric, error) {

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

	fmt.Println("GOT PARENT TST", parent)
	if parent == "" {
		fmt.Println("GOT MODEL ", p.Models.Parent)
		if p.Models.Parent != "" {
			fmt.Println("test if input as it .. in ", input.Input)
			if val, ok := input.Input[p.Models.Parent]; ok {
				parent = val.(string)
			}
		}
	}
	fmt.Println("GOT PARENT 1 - TST", parent, " and elementType is :", elementType)

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
	fmt.Println("1 - WILL Query using :", parent, " and active#", elementType, "# using gsi-data-index")
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

func (p GenericHandler) HandleCreateGeneric(reqObj *acenteralib.RequestObject, identity map[string]interface{}, mutation createInputEvent) (*map[string]interface{}, error) {
	return p.handleCreateGeneric(reqObj, identity, mutation, false, nil, nil)
}

func (p GenericHandler) handleCreateGeneric(reqObj *acenteralib.RequestObject, identity map[string]interface{}, mutation createInputEvent, isNextMutation bool, nextItem *dynamodb.PutItemInput, parentItem *map[string]interface{}) (*map[string]interface{}, error) {
	// TODO: Create if user is Admin
	fmt.Println("Create Generic ... (A)")
	elementType := p.ElementType
	if elementType == "" {
		elementType = mutation.Type
	}
	hasChildMutation := false // isNextMutation // false

	// OneToMany ?????
	// if we receive an existing id that means we already have the metadata created ...
	existingID := ""
	if val, ok := mutation.Input["id"]; ok {
		existingID = val.(string)
	}

	// Ok, if we want an deterministic id value such as an md5 of the email address as primary key
	// compute the fct + field value
	// one_to_many_update_parent_fields
	newIDFormat := ""
	isIdFormated := false
	customIdFunction := strings.Split(p.Models.CustomId, ",")
	if len(customIdFunction) >= 2 {
		fct := customIdFunction[0]
		fieldfc := customIdFunction[1]
		// ie:  id: md5,email
		if fct == "md5" {
			if val, ok := mutation.Input[fieldfc]; ok {
				hasher := md5.New()
				hasher.Write([]byte(fmt.Sprintf("%v", val)))
				newIDFormat = hex.EncodeToString(hasher.Sum(nil))
				isIdFormated = true
			} else {
				// return missing create field
				return nil, errors.New(fmt.Sprintf("Mising field '%s'", val))
			}
		} else {
			// Not yet defined
			return nil, errors.New(fmt.Sprintf("Function not implemented '%s'", fct))
		}
	}

	/*
		id: md5,email
		parent: plugin_key
	*/
	var childDynamoPutInput dynamodb.PutItemInput
	var lstParentFields map[string]string
	parentSKPrefix := ""
	childSKSuffixField := ""
	// if !isNextMutation {
	if p.Models.OneToMany != "" {
		hasChildMutation = true
		lstParentFieldsTmp := strings.Split(strings.ToLower(p.Models.OneToManyUpdateFields), ",")
		fmt.Println("GOT SPLIT FILEDSD OF ", lstParentFieldsTmp)
		lstParentFields = make(map[string]string, 0)
		for i, v := range lstParentFieldsTmp {
			lstParentFields[v] = lstParentFieldsTmp[i]
		}

		// id: md5,email
		if p.Models.OneToMany == "plugin_key" {
			// Ok THIS IS AN METADATA ... we would create metadata and keep only  the proper fields ...
			// SK will be {PLUGIN_KEY}#{MODELID}  (for metadata ... )
			// SK Will be overided with thte PLUGIN_KEY prefix

			// ok fill up
			// childDynamoPutInput
			parentSKPrefix = os.Getenv("SITE_KEY")
		} else {
			// Ok we might need to creeate a parent object first
			childSKSuffixField = strings.ToLower(p.Models.OneToMany)
			if existingID != "" {
				// Ok we will only update the parent object to either add the field as part of the aray list, and other such as authors etc...
				// based on the list arrayss
				// fmt.Println(existingID)
			} else {
				// OK must fist create a new Parent item without the "current sk" ...
				// we will use the OneToMany "value" as "SK" (in a lower string format )
			}
		}
	} else {
		// ?? reg create
	}
	// }

	// Regular Create( from here )

	taskName := mutation.Input["title"].(string)
	site := os.Getenv("SITE")

	parent := "" // os.Getenv("SITE")
	if mutation.Parent != nil {
		if val, ok := mutation.Parent["value"]; ok {
			parent = val.(string)
		}
	}
	var dynamoPutInput dynamodb.PutItemInput
	// if no parents defined in request ...
	if parent == "" {
		if p.Models.Parent != "" {
			if val, ok := mutation.Input[p.Models.Parent]; ok {
				parent = val.(string)
			}
		}
	}

	if parent == "" {
		dynamoPutInput = *p.Awslambda.CreateAppItemParent(taskName, elementType, site, "") // leave default sitei and sort key of active#ELEMENTTYPE#TS
	} else {
		dynamoPutInput = *p.Awslambda.CreateAppItemParentAndPlugin(taskName, elementType, parent, "", site, "") // leave default sitei and sort key of active#ElementTypetTT#TS
	}

	if isIdFormated {
		dynamoPutInput.Item["id"].S = aws.String(newIDFormat)
	}

	// proj represents the Projection Expression
	proj := expression.ProjectionBuilder{} // NamesList() //expression.Name("aName"), expression.Name("anotherName"), expression.Name("oneOtherName"))
	update := expression.UpdateBuilder{}

	attributes := map[string]*dynamodb.AttributeValue{}

	// will add USERID / `type`#last-modified-date so we can query liike
	// Get latest created `type` (ie: tasks) for a user...
	fmt.Println(reqObj.Session)
	mutation.Input["upk"] = reqObj.Session.Userid

	// childSK := dynamoPutInput.Item["sk"].S
	if hasChildMutation {
		fmt.Println("PARENT WILL ONLY KEEP THE FOLLOWING FIELDS", lstParentFields)
		origSK := *dynamoPutInput.Item["sk"].S
		if childSKSuffixField != "" {
			if v, ok := mutation.Input[childSKSuffixField]; ok {
				origSK = origSK + "#" + v.(string)
			} else {
				// error
				return nil, errors.New(fmt.Sprintf("Field '%s' missing", childSKSuffixField))
			}
		}
		if !isNextMutation {
			// Ok this is the parent
			dynamoPutInput.Item["sk"].S = aws.String(parentSKPrefix + "#" + *dynamoPutInput.Item["sk"].S)
		} else {
			// Ok this is the child
			dynamoPutInput.Item["sk"].S = aws.String(origSK)
		}
	}

	// Override any values including sk if user wants it ...
	exisingKeys := make(map[string]*dynamodb.AttributeValue, 0)
	for lowerNameKey, attr := range dynamoPutInput.Item {
		exisingKeys[lowerNameKey] = attr
		if lowerNameKey == "id" || lowerNameKey == "sk" {
			// do not ad it to the update expressions
		} else {
			isFieldOk := false
			if hasChildMutation {
				if _, ok := lstParentFields[lowerNameKey]; ok {
					if !isNextMutation {
						isFieldOk = true
					}
				} else {
					if isNextMutation {
						isFieldOk = true
					}
				}
			} else {
				isFieldOk = true
			}

			if isFieldOk {
				nameLower := expression.Name(lowerNameKey)
				proj = expression.AddNames(proj, nameLower)

				value := expression.Value("")
				if attr.S != nil {
					value = expression.Value(attr.S)
				} else if attr.N != nil {
					value = expression.Value(attr.N)
				} else if attr.B != nil {
					value = expression.Value(attr.B)
				} else if attr.BOOL != nil {
					value = expression.Value(attr.BOOL)
				} else if attr.BS != nil {
					value = expression.Value(attr.BS)
				} else if attr.L != nil {
					value = expression.Value(attr.L)
				} else if attr.M != nil {
					value = expression.Value(attr.M)
				} else if attr.NS != nil {
					value = expression.Value(attr.NS)
				} else if attr.NULL != nil {
					value = expression.Value(attr.NULL)
				} else if attr.SS != nil {
					value = expression.Value(attr.SS)
				}
				update = update.Set(nameLower, value)
			} else {
				/*if (hasChildMutation && !isNextMutation) {
					childDynamoPutInput
				}
				*/
			}
		}
	}
	for k, v := range mutation.Input {
		lowerNameKey := strings.ToLower(k)
		fmt.Println("Check of  :", lowerNameKey)
		if _, ok := exisingKeys[lowerNameKey]; ok {
			// Already exists do not re-add
		} else {
			isFieldOk := false
			if hasChildMutation {
				if _, ok := lstParentFields[lowerNameKey]; ok {
					if !isNextMutation {
						isFieldOk = true
					}
				} else {
					if isNextMutation {
						isFieldOk = true
					}
				}
			} else {
				isFieldOk = true
			}

			if isFieldOk {
				nameLower := expression.Name(lowerNameKey)
				fmt.Println("22 - Adding of ", nameLower)
				value := expression.Value(v)
				attrB, _ := dynamodbattribute.Marshal(v)
				attributes[lowerNameKey] = attrB
				dynamoPutInput.Item[lowerNameKey] = attrB
				proj = expression.AddNames(proj, nameLower)
				update = update.Set(nameLower, value)
				fmt.Println("11 - Adding done..")
			} else {
				fmt.Println("Ignoring field: ", lowerNameKey)
			}
		}
	}
	fmt.Println("Ok Done..")
	condUpdateCond := expression.And(expression.AttributeNotExists(expression.Name("id")), expression.AttributeNotExists(expression.Name("sk"))) // event.Input["sk"])) // sortKey or type?
	builder := expression.NewBuilder().WithCondition(condUpdateCond).WithUpdate(update)
	exp, err := builder.Build()

	dynamoDdbUpdateInput := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": dynamoPutInput.Item["id"],
			"sk": dynamoPutInput.Item["sk"],
		},
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		UpdateExpression:          exp.Update(),
		ConditionExpression:       exp.Condition(),
		ReturnValues:              aws.String("ALL_NEW"),
		TableName:                 aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	}
	// fmt.Println(exp.Update())

	// Ok need to
	if hasChildMutation {
		if !isNextMutation {
			mutation.Input["ref_id"] = *dynamoPutInput.Item["id"].S
			mutation.Input["ref_sk"] = *dynamoPutInput.Item["sk"].S
		}
	}

	output, err := acenteralib.DynamoDB.UpdateItem(dynamoDdbUpdateInput)
	fmt.Println(output)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)
		if ae, ok := err.(awserr.RequestFailure); ok && ae.Code() == "ConditionalCheckFailedException" {
			//fmt.Println("Update failed due to constraint...")
			return nil, errors.WithStack(err)
		}
		return nil, errors.WithStack(err)
	}

	// do not use put items..
	var item map[string]interface{}
	if err == nil {
		err = dynamodbattribute.UnmarshalMap(output.Attributes, &item)
	}

	// TODO: Ned to return a UserPost Connection...

	if hasChildMutation {
		if !isNextMutation {
			// Ok create the child element without the parent fields
			return p.handleCreateGeneric(reqObj, identity, mutation, true, &childDynamoPutInput, &item)
		} else {
			// Ok we are the child, return the item
			// TODO: Return the merged parent fields ???
			return &item, err
		}
	} else {
		return &item, err
	}

	/*
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
	*/
}

type GenericHandler struct {
	Awslambda          acenteralib.SharedLib
	Models             Models
	InitHandler        func(r resolvers.Repository) error
	ElementType        string
	ActionWord         string
	ActionWordSingular string
	ActionWordPlurial  string
}

func (p GenericHandler) InitializeRoutes(r resolvers.Repository) error {
	var err error
	fmt.Println("TEST HERE  ", p.InitHandler)
	if p.InitHandler != nil {
		err = p.InitHandler(r)
	} else {
		fmt.Println("TEST HERE  B ", p.ElementType)
		if p.ElementType == "" {
			// this is impossible... elementytpe is requirred
			err = r.Add(fmt.Sprintf("query.list"), p.HandleList)
			err = r.Add(fmt.Sprintf("query.listAll"), p.HandleListAll)
			err = r.Add(fmt.Sprintf("mutation.create"), p.HandleCreateGeneric)
			err = r.Add(fmt.Sprintf("mutation.update"), p.HandleUpdateGeneric)
			err = r.Add(fmt.Sprintf("mutation.delete"), p.HandleDeleteGeneric)
		} else {
			if p.Models.Version == 1 {
				fmt.Println("GraphQL Adding :", fmt.Sprintf("query.list%s", p.ActionWordPlurial))
				err = r.Add(fmt.Sprintf("query.list%s", p.ActionWordPlurial), p.HandleList) // query.listPosts(listInput) with nextTokens

				fmt.Println("GraphQL Adding :", fmt.Sprintf("query.%s", strings.ToLower(p.ActionWordPlurial)))
				err = r.Add(fmt.Sprintf("query.%s", strings.ToLower(p.ActionWordPlurial)), p.HandleList) // query.posts(listInput)

				// err = r.Add(fmt.Sprintf("query.batchget%s", p.ActionWordPlurial), p.HandleBatchGet)
				// err = r.Add(fmt.Sprintf("query.%s", strings.ToLower(p.ActionWordSingular)), p.HandleGetGeneric)			 // query.post(id=x)

				fmt.Println("GraphQL Adding :", fmt.Sprintf("query.listAll%s", p.ActionWordPlurial))
				err = r.Add(fmt.Sprintf("query.listAll%s", p.ActionWordPlurial), p.HandleListAll) // listAllPosts (no paging)
				// err = r.Add(fmt.Sprintf("query.listAll%ss", ActionWord), p.HandleListAll) // plurial

				fmt.Println("GraphQL Adding CRUDS for :", fmt.Sprintf("mutation.[create|update|delete|%s", p.ActionWordSingular))
				err = r.Add(fmt.Sprintf("mutation.create%s", p.ActionWordSingular), p.HandleCreateGeneric)
				err = r.Add(fmt.Sprintf("mutation.update%s", p.ActionWordSingular), p.HandleUpdateGeneric)
				err = r.Add(fmt.Sprintf("mutation.delete%s", p.ActionWordSingular), p.HandleDeleteGeneric)

				if p.Models.Admin == 1 {
					err = r.Add(fmt.Sprintf("query.list%sAdmin", p.ActionWordPlurial), p.HandleListAdmin) // query.listPosts(listInput) with nextTokens
				}

				fmt.Println("DID WE GOT RESEOLVERS?", p.Models.Resolvers)
				for v, k := range p.Models.Resolvers {
					fmt.Println("Adding custom resolvers : ", fmt.Sprintf("mutation.%s", v), " for :", k.(string))
					if k.(string) == "create" {
						fmt.Println("Added custom resolvers : ", fmt.Sprintf("mutation.%s", v), " for : creaete generic")
						err = r.Add(fmt.Sprintf("mutation.%s", v), p.HandleCreateGeneric)
						fmt.Println("ERR IS :", err)
					}
				}
			}
		}
	}
	return err
}
