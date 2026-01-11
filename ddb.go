package goddb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func NewDDB(dynamoDBClient *dynamodb.Client) *DDBClient {
	return &DDBClient{
		client: dynamoDBClient,
	}
}

func (c *DDBClient) AddTable(tableName string, pk, sk DDBAttributes, isCreate bool, billingMode DDBBillingMode) DDBClient {

	c.tables[tableName] = TableParms{
		Pk:              pk.Name,
		PkAttributeType: pk.AttributeType,

		Sk:              sk.Name,
		SkAttributeType: sk.AttributeType,

		IsCreate:    isCreate,
		BillingMode: billingMode,
	}

	return *c
}

func (c *DDBClient) Start() error {

	for tableName, parmas := range c.tables {

		if parmas.IsCreate {
			_, err := c.client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
				TableName: aws.String(tableName),

				// KeySchema
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String(parmas.Pk),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String(parmas.Sk),
						KeyType:       types.KeyTypeRange,
					},
				},

				// AttributeDefinitions
				AttributeDefinitions: []types.AttributeDefinition{
					{
						AttributeName: aws.String(parmas.Pk),
						AttributeType: types.ScalarAttributeType(parmas.PkAttributeType),
					},
					{
						AttributeName: aws.String(parmas.Sk),
						AttributeType: types.ScalarAttributeType(parmas.SkAttributeType),
					},
				},

				// BillingMode
				BillingMode: types.BillingMode(getBillingMode(parmas.BillingMode)),
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(int64(parmas.BillingMode.isProvisioned.ReadCapacityUnits)),
					WriteCapacityUnits: aws.Int64(int64(parmas.BillingMode.isProvisioned.WriteCapacityUnits)),
				},
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getBillingMode(billingMode DDBBillingMode) types.BillingMode {
	if billingMode.isOnDemand {
		return types.BillingModePayPerRequest
	}

	return types.BillingModeProvisioned
}
