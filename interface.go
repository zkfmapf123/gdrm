package goddb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DDBBillingMode struct {
	isOnDemand    bool
	isProvisioned struct {
		ReadCapacityUnits  int
		WriteCapacityUnits int
	}
}

type DDBAttributes struct {
	Name          string
	AttributeType types.ScalarAttributeType
}

type TableParms struct {
	Pk              string // primary key
	PkAttributeType types.ScalarAttributeType

	Sk              string // sort key
	SkAttributeType types.ScalarAttributeType

	IsCreate    bool
	BillingMode DDBBillingMode
}

type DDBClient struct {
	client *dynamodb.Client
	tables map[string]TableParms
}
