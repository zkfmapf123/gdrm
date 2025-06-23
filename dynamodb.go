package dynamodbgo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoGORMParmas struct {
	db dynamodb.Client
}

func DynamoGORM(c context.Context) DynamoGORMParmas {

	client, err := config.LoadDefaultConfig(c)
	if err != nil {
		panic(err)
	}

	db := dynamodb.NewFromConfig(client)
	return DynamoGORMParmas{
		db: *db,
	}
}

func (orm *DynamoGORMParmas) isExistTableName(tableName string) (bool, error) {

	_, err := orm.db.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: wrapString(tableName),
	})

	// not found
	if err != nil {
		return false, err
	}

	return true, nil
}

func wrapString(v string) *string {
	return aws.String(v)
}

func wrapInt(v int) *int {
	return aws.Int(v)
}

func wrapBool(v bool) *bool {
	return aws.Bool(v)
}

func getBillingMode(v bool) types.BillingMode {
	if v {
		return types.BillingModePayPerRequest
	}

	return types.BillingModeProvisioned

}

func convertToAttributeValueMap(data map[string]any) map[string]types.AttributeValue {
	item := make(map[string]types.AttributeValue)

	for key, value := range data {
		switch v := value.(type) {
		case string:
			item[key] = &types.AttributeValueMemberS{Value: v}
		case int:
			item[key] = &types.AttributeValueMemberN{Value: strconv.Itoa(v)}
		case int64:
			item[key] = &types.AttributeValueMemberN{Value: strconv.FormatInt(v, 10)}
		case float64:
			item[key] = &types.AttributeValueMemberN{Value: strconv.FormatFloat(v, 'f', -1, 64)}
		case bool:
			item[key] = &types.AttributeValueMemberBOOL{Value: v}
		case []byte:
			item[key] = &types.AttributeValueMemberB{Value: v}
		case []string:
			item[key] = &types.AttributeValueMemberSS{Value: v}
		case []int:
			numbers := make([]string, len(v))
			for i, num := range v {
				numbers[i] = strconv.Itoa(num)
			}
			item[key] = &types.AttributeValueMemberNS{Value: numbers}
		case map[string]any:
			nestedMap := convertToAttributeValueMap(v)
			item[key] = &types.AttributeValueMemberM{Value: nestedMap}
		case []any:
			// any 슬라이스 처리
			list := make([]types.AttributeValue, len(v))
			for i, val := range v {
				switch listVal := val.(type) {
				case string:
					list[i] = &types.AttributeValueMemberS{Value: listVal}
				case int:
					list[i] = &types.AttributeValueMemberN{Value: strconv.Itoa(listVal)}
				case bool:
					list[i] = &types.AttributeValueMemberBOOL{Value: listVal}
				default:
					list[i] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", val)}
				}
			}
			item[key] = &types.AttributeValueMemberL{Value: list}
		default:
			// 지원하지 않는 타입은 문자열로 변환
			item[key] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", value)}
		}
	}

	return item
}
