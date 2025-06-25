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
	TableName string

	Primarykey  string // if already exists not used
	SortedKey   string // if already exists not used
	BillingMode bool   // true : on-demain , false : provisioned
}

func Insert(params TableParmas, data map[string]any) error {
	orm := dynamoGORM(context.Background())

	attempt := 1
	for attempt = 1; attempt < RETRY_COUNT; attempt++ {

		if !isExistTableName(params.TableName) {
			createTable(&orm, params)
		}

		if isActiveDynamoTable(params.TableName) {
			break
		}

		time.Sleep(time.Second * 2)
		fmt.Println("테이블 생성 중... ")
	}

	item, err := serializeToDynamoDB(data)
	if err != nil {
		return err
	}

	_, err = orm.db.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: &params.TableName,
		Item:      item,
	})

	return err
}

func createTable(orm *dynamoGORMParmas, params TableParmas) error {

	attrDefinition := []types.AttributeDefinition{
		{
			AttributeName: wrapString(params.Primarykey),
			AttributeType: types.ScalarAttributeTypeS,
		},
	}

	keySchema := []types.KeySchemaElement{
		{
			AttributeName: wrapString(params.Primarykey),
			KeyType:       types.KeyTypeHash,
		},
	}

	// TOBE. 정렬 키는 추후 구성 예정
	_, err := orm.db.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName:            wrapString(params.TableName),
		AttributeDefinitions: attrDefinition,
		KeySchema:            keySchema,
		BillingMode:          getBillingMode(params.BillingMode),
	})

	return err
}

func isActiveDynamoTable(tableName string) bool {
	ctx := context.Background()
	orm := dynamoGORM(ctx)

	res, err := orm.db.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: wrapString(tableName),
	})

	if err != nil {
		// log.Println(err)
		return false
	}

	return res.Table.TableStatus == types.TableStatusActive
}
