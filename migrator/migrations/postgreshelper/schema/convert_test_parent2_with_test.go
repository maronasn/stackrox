// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

// ConvertTestParent2FromProto converts a `*storage.TestParent2` to Gorm model
func ConvertTestParent2FromProto(obj *storage.TestParent2) (*schema.TestParent2, error) {
	serialized, err := obj.Marshal()
	if err != nil {
		return nil, err
	}
	model := &schema.TestParent2{
		Id:         obj.GetId(),
		ParentId:   obj.GetParentId(),
		Val:        obj.GetVal(),
		Serialized: serialized,
	}
	return model, nil
}

// ConvertTestParent2ToProto converts Gorm model `TestParent2` to its protobuf type object
func ConvertTestParent2ToProto(m *schema.TestParent2) (*storage.TestParent2, error) {
	var msg storage.TestParent2
	if err := msg.Unmarshal(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}

func TestTestParent2Serialization(t *testing.T) {
	obj := &storage.TestParent2{}
	assert.NoError(t, testutils.FullInit(obj, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	m, err := ConvertTestParent2FromProto(obj)
	assert.NoError(t, err)
	conv, err := ConvertTestParent2ToProto(m)
	assert.NoError(t, err)
	assert.Equal(t, obj, conv)
}
