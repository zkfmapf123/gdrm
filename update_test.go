package dynamodbgo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Update(t *testing.T) {

	data, _ := FindByUsePK[ResponseParams]("users", "user_id", "user_14")

	assert.Equal(t, data.UserId, "user_14")
	assert.Equal(t, data.BB, true)
	assert.Equal(t, data.CC, true)

	UpdatePartial("users", "user_id", "user_14", map[string]any{
		"bb": false,
		"cc": false,
	})

	updatedData, _ := FindByUsePK[ResponseParams]("users", "user_id", "user_14")
	fmt.Println("updated ", updatedData)

}
