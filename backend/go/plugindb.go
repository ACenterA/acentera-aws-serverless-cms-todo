package gofaas

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
	"time"

	"github.com/acenteracms/acenteralib"
	"github.com/pkg/errors"
)

type PluginInfo struct {
	Id      string `json:"id"`
	SortKey string `json:"sk"`
	Title   string `json:"title"`
	Created string `json:"created"`
	Type    string `json:"type"`
	Data    string `json:"json_data"` // this is a JSON Document for now
}

func UpdateAMIProjectStateById(ctx context.Context, objectId string, newData string) (res *PluginInfo, err error) {
	// Update AMI
	//fmt.Println("WILL UPDATE OF ", objectId, " and.. sk = SELF ")
	_, err = DynamoDB.UpdateItemWithContext(ctx, &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(objectId),
			},
			"sk": &dynamodb.AttributeValue{
				S: aws.String("SELF"),
			},
		},
		UpdateExpression: aws.String("SET #ATTR = :val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": &dynamodb.AttributeValue{
				S: aws.String(newData), // time.Now().UTC().Format(time.RFC3339)),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#ATTR": aws.String("json_data"),
		},
		// ConditionExpression: aws.String(":version = #VERSION"),
		ReturnValues: aws.String("NONE"),
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
	return nil, nil
}

func GetAMIProjectById(ctx context.Context, site string, uuid string) (res *PluginInfo, err error) {
	// For now we do not use an GSI table...
	// The GetPlugin isn't something that will occurs often
	//fmt.Println("WILL QUERY USING ", "gsk = ", uuid, " and gpk = SELF in table ", os.Getenv("APP_DATA_TABLE_NAME"))
	out, errTmp := DynamoDB.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(uuid),
			},
			":sk": {
				S: aws.String("SELF"),
			},
		},
		// IndexName: aws.String("gsi-data-index"),
		KeyConditionExpression: aws.String("id = :pk and sk = :sk"),
		// no filter we use plugin#name for now FilterExpression: aws.String("title = :title"),
		Limit: aws.Int64(1), // should only return 1
		// all ... ProjectionExpression:   aws.String("id, sk, jsondata"),
		// all ... for now ProjectionExpression:   aws.String("SongTitle"),
		TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	})
	err = errTmp

	// TODO: Make better implementation here
	if err != nil {
		//fmt.Println("ERR HERE A")
		//fmt.Println(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				//fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				//fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				//fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				//fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			//fmt.Println(err.Error())
		}
		return nil, ResponseError{"not found", 404}
	}
	//fmt.Println("ERR HERE B")
	if out == nil {
		return nil, ResponseError{"not found", 404}
	}

	//fmt.Println("1 ERR HERE C 1")
	//fmt.Println(*out)
	items := []PluginInfo{}
	//fmt.Println("1 ERR HERE C 2")
	if *out.Count == 1 {
		//fmt.Println("1 ERR HERE C 3")
		err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &items)
		return &items[0], err
	}
	return nil, ResponseError{"not found", 404}
}

func GetProjectElementsProjectById(ctx context.Context, uuid string) (res *[]PluginInfo, err error) {
	// The GetPlugin isn't something that will occurs often
	//fmt.Println("GET USING GPK OF :", uuid, " and gsk: active# with gsi..")
	out, errTmp := DynamoDB.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(uuid),
			},
			":gsk": {
				S: aws.String("active#"), // this hssould be like
			},
		},
		IndexName:              aws.String("gsi-data-index"),
		KeyConditionExpression: aws.String("gpk = :pk and begins_with( gsk, :gsk )"),
		// Limit: aws.Int64(1), // should only return 1
		// all ... ProjectionExpression:   aws.String("id, sk, jsondata"),
		// all ... for now ProjectionExpression:   aws.String("SongTitle"),
		TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	})
	/*
		//fmt.Println("WILL QUERY USING ", "gsk = ", site)
		out, errTmp := DynamoDB.Query(&dynamodb.QueryInput{
	    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
	        ":pk": {
	            S: aws.String(site),
					},
					":sk": {
						S: aws.String(fmt.Sprintf("%v#", k)),
					},
	    },
			IndexName: aws.String("gsi-data-index"),
	    KeyConditionExpression: aws.String("gpk = :pk and begins_with( gsk, :sk )"),
	*/
	err = errTmp

	// TODO: Make better implementation here
	if err != nil {
		//fmt.Println("ERR HERE A")
		//fmt.Println(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				//fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				//fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				//fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				//fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			//fmt.Println(err.Error())
		}
		return nil, ResponseError{"not found", 404}
	}
	//fmt.Println("ERR HERE B")
	if out == nil {
		return nil, ResponseError{"not found", 404}
	}

	//fmt.Println("1 ERR HERE C 1")
	//fmt.Println(*out)
	items := []PluginInfo{}
	//fmt.Println("1 ERR HERE C 2")
	if *out.Count >= 1 {
		//fmt.Println("1 ERR HERE C 3")
		err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &items)
		return &items, err
	}
	return nil, ResponseError{"not found", 404}
}

func AMIGetByProject(ctx context.Context, site string) (res []PluginInfo, err error) {
	// For now we do not use an GSI table...
	// The GetPlugin isn't something that will occurs often
	fmt.Println("WILL QUERY USING ", "gsk = ", site, " and sk = amiprj#")
	out, errTmp := DynamoDB.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(site),
			},
			":sk": {
				S: aws.String("amiprj#"),
			},
		},
		IndexName:              aws.String("gsi-data-index"),
		KeyConditionExpression: aws.String("gpk = :pk and begins_with( gsk, :sk )"),
		// no filter we use plugin#name for now FilterExpression: aws.String("title = :title"),
		// Limit: aws.Int64(1), // should only return 1
		TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	})
	err = errTmp

	// TODO: Make better implementation here
	if err != nil {
		//fmt.Println("ERR HERE A")
		fmt.Println(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				//fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				//fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				//fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				//fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			//fmt.Println(err.Error())
		}
		return nil, ResponseError{"not found", 404}
	}
	//fmt.Println("ERR HERE B")
	if out == nil {
		return nil, ResponseError{"not found", 404}
	}

	//fmt.Println("ERR HERE C")
	//fmt.Println(*out)
	items := []PluginInfo{}
	if *out.Count >= 1 {
		err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &items)
		return items, err
	}
	return items, nil // none nil, ResponseError{"not found", 404}
}

func GetByProject(ctx context.Context, site string, k string) (res []PluginInfo, err error) {
	// For now we do not use an GSI table...
	// The GetPlugin isn't something that will occurs often
	//fmt.Println("WILL QUERY USING ", "gsk = ", site)
	out, errTmp := DynamoDB.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(site),
			},
			":sk": {
				S: aws.String(fmt.Sprintf("%v#", k)),
			},
		},
		IndexName:              aws.String("gsi-data-index"),
		KeyConditionExpression: aws.String("gpk = :pk and begins_with( gsk, :sk )"),
		// no filter we use plugin#name for now FilterExpression: aws.String("title = :title"),
		// Limit: aws.Int64(1), // should only return 1
		TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	})
	err = errTmp

	// TODO: Make better implementation here
	if err != nil {
		//fmt.Println("ERR HERE A")
		//fmt.Println(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				//fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				//fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				//fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				//fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			//fmt.Println(err.Error())
		}
		return nil, ResponseError{"not found", 404}
	}
	//fmt.Println("ERR HERE B")
	if out == nil {
		return nil, ResponseError{"not found", 404}
	}

	//fmt.Println("ERR HERE C")
	//fmt.Println(*out)
	items := []PluginInfo{}
	if *out.Count >= 1 {
		err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &items)
		return items, err
	}
	return items, nil // none nil, ResponseError{"not found", 404}
}

func GetByProjectAndFilter(ctx context.Context, site string, k string, filterKey string, filterValue string) (res []PluginInfo, err error) {
	// For now we do not use an GSI table...
	// The GetPlugin isn't something that will occurs often
	//fmt.Println("WILL QUERY USING ", "gsk = ", site)
	out, errTmp := DynamoDB.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(site),
			},
			":sk": {
				S: aws.String(fmt.Sprintf("%v#", k)),
			},
			":filterval": {
				S: aws.String(fmt.Sprintf("%v", filterValue)),
			},
		},
		IndexName:              aws.String("gsi-data-index"),
		KeyConditionExpression: aws.String("gpk = :pk and begins_with( gsk, :sk )"),
		FilterExpression:       aws.String(fmt.Sprintf("%v = :filterval", filterKey)),
		// Limit: aws.Int64(1),
		// TODO: last settings
		// ProjectionExpression:   aws.String("id, sk, title, jsondata"),
		// no filter we use plugin#name for now FilterExpression: aws.String("title = :title"),
		Limit:     aws.Int64(1), // should only return 1
		TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	})
	err = errTmp

	// TODO: Make better implementation here
	if err != nil {
		//fmt.Println("ERR HERE A")
		//fmt.Println(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				//fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				//fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				//fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				//fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			//fmt.Println(err.Error())
		}
		return nil, ResponseError{"not found", 404}
	}
	//fmt.Println("ERR HERE B")
	if out == nil {
		return nil, ResponseError{"not found", 404}
	}

	//fmt.Println("ERR HERE C")
	//fmt.Println(*out)
	items := []PluginInfo{}
	if *out.Count >= 1 {
		err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &items)
		return items, err
	}
	return items, nil // none nil, ResponseError{"not found", 404}
}

// a, errT := CreateSiteData(reqObj, "A", createAmiObj.Name)

func CreateSiteData(ctx context.Context, reqObj acenteralib.RequestObject, strType string, pluginName string, jsonData string, pluginId string) (uuid string, err error) {
	// In this case all App Data are

	// Create a specific AMI Project UUID
	uuid = UUIDGen().String()
	CreatedTime := time.Now().UTC().Format(time.RFC3339)
	_, err = DynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(uuid),
			},
			"sk": &dynamodb.AttributeValue{
				S: aws.String("SELF"), // fmt.Sprintf("amiprj#%s",CreatedTime)),
			},
			"uk": &dynamodb.AttributeValue{
				S: aws.String(uuid),
			},
			"title": &dynamodb.AttributeValue{
				S: aws.String(pluginName),
			},
			"json_data": &dynamodb.AttributeValue{
				S: aws.String(jsonData),
			},
			"created": &dynamodb.AttributeValue{ // This Created Time is being used for access of the DATA_PERM_TABLE if we need to update title or something else ???
				S: aws.String(CreatedTime),
			},
			"gpk": &dynamodb.AttributeValue{ // in case we createa GSI on this table....
				S: aws.String(reqObj.Site),
			},
			"gsk": &dynamodb.AttributeValue{ // in case we createa GSI on this table....
				S: aws.String(fmt.Sprintf("%s#%s", strType, CreatedTime)),
			},
			"ppk": &dynamodb.AttributeValue{ // in case we createa GSI on this table....
				S: aws.String(pluginId),
			},
			"psk": &dynamodb.AttributeValue{ // in case we create a new GSI on this table... to list all iste for a plugin?
				S: aws.String(fmt.Sprintf("%s#%s", strType, CreatedTime)),
			},
		},
		TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	})

	if err != nil {
		return "", err
	}
	/*
	  // For a Site list top X AMIProject  (dashboard)
	  pk := reqObj.Site
		sk := fmt.Sprintf("amiprj#%s", CreatedTime)
		// OK at this point we successfully created a new Data Object
		// Lets make sure this current Site have access to this data...
		// Allow to list the latest AMI Projects "AMIPRJ#{CreationDate}" search by "Site ID's"
		_, err = DynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"id": &dynamodb.AttributeValue{
					S: aws.String(pk),
				},
				"sk": &dynamodb.AttributeValue{
					S: aws.String(sk),
				},
				"title": &dynamodb.AttributeValue{
					S: aws.String(pluginName),
				},
				"amiprj": &dynamodb.AttributeValue{                // Reference to the real Project if needed or for future queries ... or for GSI
					S: aws.String(uuid),
				},
				"created": &dynamodb.AttributeValue{
					S: aws.String(CreatedTime),
				},
			},
			TableName: aws.String(os.Getenv("DATA_PERM_TABLE_NAME")),
		})

		if (err != nil) {
			// TODO: Send an SNS to Delete the orphaned PutItem we just created ?
			return "", err
		}
	*/

	return uuid, nil
}

func GetPlugin(reqObj acenteralib.RequestObject, pluginName string) (u *PluginInfo, err error) {

	// For now we do not use an GSI table...
	// The GetPlugin isn't something that will occurs often
	//fmt.Println("WILL QUERY USING ", reqObj.Site, " and sk = plugin#")
	out, errTmp := DynamoDB.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(reqObj.Site),
			},
			":sk_value": {
				S: aws.String(fmt.Sprintf("plugin#%s", pluginName)), // # may contain date or a file date?
			},
		},
		// IndexName: aws.String("gsi-data-index"),
		// KeyConditionExpression: aws.String("gpk = :gpk and gsk = :gsk_value"),
		KeyConditionExpression: aws.String("id = :id and begins_with(sk, :sk_value)"),
		// no filter we use plugin#name for now FilterExpression: aws.String("title = :title"),
		// Limit: aws.Int64(1),
		// TODO: last settings
		ProjectionExpression: aws.String("id, sk, title, jsondata"),
		// all ... for now ProjectionExpression:   aws.String("SongTitle"),
		TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	})
	err = errTmp

	// TODO: Make better implementation here
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				//fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				//fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				//fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				//fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			//fmt.Println(err.Error())
		}
		return nil, ResponseError{"not found", 404}
	}

	if out == nil {
		return nil, ResponseError{"not found", 404}
	}

	items := []PluginInfo{}
	if *out.Count == 1 {
		err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &items)
		return &items[0], err
	}
	return nil, ResponseError{"not found", 404}
}

func CreateProjectData(ctx context.Context, reqObj acenteralib.RequestObject, strType string, currType string, pluginName string, jsonData string, pluginId string, refPK string) (uuid string, err error) {
	// In this case all App Data are

	// Create a specific AMI Project UUID
	uuid = UUIDGen().String()
	CreatedTime := time.Now().UTC().Format(time.RFC3339)
	_, err = DynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(uuid),
			},
			"sk": &dynamodb.AttributeValue{
				S: aws.String("SELF"), // fmt.Sprintf("amiprj#%s",CreatedTime)),
			},
			"uk": &dynamodb.AttributeValue{
				S: aws.String(uuid),
			},
			"title": &dynamodb.AttributeValue{
				S: aws.String(pluginName),
			},
			"json_data": &dynamodb.AttributeValue{
				S: aws.String(jsonData),
			},
			"created": &dynamodb.AttributeValue{ // This Created Time is being used for access of the DATA_PERM_TABLE if we need to update title or something else ???
				S: aws.String(CreatedTime),
			},
			"gpk": &dynamodb.AttributeValue{ // in case we createa GSI on this table....
				S: aws.String(refPK),
			},
			"gsk": &dynamodb.AttributeValue{ // in case we createa GSI on this table....
				S: aws.String(fmt.Sprintf("%s#%s", strType, CreatedTime)),
			},
			"ssk": &dynamodb.AttributeValue{ // in case we createa GSI on this table....
				S: aws.String(fmt.Sprintf("%s#%s", strType, CreatedTime)),
			},
			"spk": &dynamodb.AttributeValue{ // in case we createa GSI on this table....
				S: aws.String(reqObj.Site),
			},
			"ppk": &dynamodb.AttributeValue{ // in case we createa GSI on this table....
				S: aws.String(pluginId),
			},
			"psk": &dynamodb.AttributeValue{ // in case we create a new GSI on this table... to list all iste for a plugin?
				S: aws.String(fmt.Sprintf("%s#%s", strType, CreatedTime)),
			},
			"type": &dynamodb.AttributeValue{ // in case we create a new GSI on this table... to list all iste for a plugin?
				S: aws.String(fmt.Sprintf("%s", currType)),
			},
		},
		TableName: aws.String(os.Getenv("APP_DATA_TABLE_NAME")),
	})

	if err != nil {
		return "", err
	}
	/*
	  // For a Site list top X AMIProject  (dashboard)
	  pk := reqObj.Site
		sk := fmt.Sprintf("amiprj#%s", CreatedTime)
		// OK at this point we successfully created a new Data Object
		// Lets make sure this current Site have access to this data...
		// Allow to list the latest AMI Projects "AMIPRJ#{CreationDate}" search by "Site ID's"
		_, err = DynamoDB.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"id": &dynamodb.AttributeValue{
					S: aws.String(pk),
				},
				"sk": &dynamodb.AttributeValue{
					S: aws.String(sk),
				},
				"title": &dynamodb.AttributeValue{
					S: aws.String(pluginName),
				},
				"amiprj": &dynamodb.AttributeValue{                // Reference to the real Project if needed or for future queries ... or for GSI
					S: aws.String(uuid),
				},
				"created": &dynamodb.AttributeValue{
					S: aws.String(CreatedTime),
				},
			},
			TableName: aws.String(os.Getenv("DATA_PERM_TABLE_NAME")),
		})

		if (err != nil) {
			// TODO: Send an SNS to Delete the orphaned PutItem we just created ?
			return "", err
		}
	*/

	return uuid, nil
}
