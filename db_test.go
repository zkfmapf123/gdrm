package dynamodbgo

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var orm = DynamoGORM(context.Background())

func Test_IsExistsTable(t *testing.T) {

	isOk, _ := orm.isExistTableName("not-found")
	assert.Equal(t, isOk, false)
}

func Test_isInsertTable(t *testing.T) {

	value := map[string]any{
		"user-id": "user-1",
		"bb":      true,
		"cc":      true,
	}

	for i := 0; i < 10; i++ {

		value["user-id"] = fmt.Sprintf("user-%d", i)
		err := orm.Insert(TableParmas{
			tableName:   "users",
			primarykey:  "user-id",
			billingMode: true,
		}, value)

		if err != nil {
			log.Fatalln(err)
		}
	}

}
