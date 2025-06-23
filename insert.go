package dynamodbgo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	RETRY_COUNT = 10
)

type TableParmas struct {
	tableName string

	primarykey  string // if already exists not used
	sortedKey   string // if already exists not used
	billingMode bool   // true : on-demain , false : provisioned
}

func (orm *DynamoGORMParmas) Insert(params TableParmas, data map[string]any) error {

	attempt := 1
	for attempt = 1; attempt < RETRY_COUNT; attempt++ {

		isOk, _ := orm.isExistTableName(params.tableName)
		if isOk {
			break
		}

		createTable(orm, params)
		time.Sleep(time.Second * 1)
		fmt.Println("테이블 생성 중... ")
	}

	_, err := orm.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: &params.tableName,
		Item:      convertToAttributeValueMap(data),
	})

	return err
}

func createTable(orm *DynamoGORMParmas, params TableParmas) error {

	attrDefinition := []types.AttributeDefinition{
		{
			AttributeName: wrapString(params.primarykey),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}

	keySchema := []types.KeySchemaElement{
		{
			AttributeName: wrapString(params.primarykey),
			KeyType:       types.KeyTypeHash,
		},
	}

	// TOBE. 정렬 키는 추후 구성 예정

	_, err := orm.db.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName:            wrapString(params.tableName),
		AttributeDefinitions: attrDefinition,
		KeySchema:            keySchema,
		BillingMode:          getBillingMode(params.billingMode),
	})

	return err
}
