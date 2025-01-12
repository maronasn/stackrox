// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

// ConvertTestG3GrandChild1FromProto converts a `*storage.TestG3GrandChild1` to Gorm model
func ConvertTestG3GrandChild1FromProto(obj *storage.TestG3GrandChild1) (*schema.TestG3GrandChild1, error) {
	serialized, err := obj.Marshal()
	if err != nil {
		return nil, err
	}
	model := &schema.TestG3GrandChild1{
		Id:         obj.GetId(),
		Val:        obj.GetVal(),
		Serialized: serialized,
	}
	return model, nil
}

// ConvertTestG3GrandChild1ToProto converts Gorm model `TestG3GrandChild1` to its protobuf type object
func ConvertTestG3GrandChild1ToProto(m *schema.TestG3GrandChild1) (*storage.TestG3GrandChild1, error) {
	var msg storage.TestG3GrandChild1
	if err := msg.Unmarshal(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}

func TestTestG3GrandChild1Serialization(t *testing.T) {
	obj := &storage.TestG3GrandChild1{}
	assert.NoError(t, testutils.FullInit(obj, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	m, err := ConvertTestG3GrandChild1FromProto(obj)
	assert.NoError(t, err)
	conv, err := ConvertTestG3GrandChild1ToProto(m)
	assert.NoError(t, err)
	assert.Equal(t, obj, conv)
}
