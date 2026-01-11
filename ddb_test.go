package goddb

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func Test_DDBStart(t *testing.T) {

	awsClient := aws.NewConfig()
	ddbClient := dynamodb.NewFromConfig(*awsClient)

	ddb := NewDDB(ddbClient)
	ddb.
		AddTable("user_logs", DDBAttributes{
			Name:          "id",
			AttributeType: types.ScalarAttributeTypeS,
		}, DDBAttributes{
			Name:          "id",
			AttributeType: types.ScalarAttributeTypeS,
		}, true, DDBBillingMode{
			isOnDemand: false,
		})

}
