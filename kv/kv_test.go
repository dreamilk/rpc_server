package kv

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Skip()

	ctx := context.Background()
	v, err := Get(ctx, "redis")
	assert.NoError(t, err)
	fmt.Println(v)
}
