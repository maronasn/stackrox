// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableTestGrandparentsStmt holds the create statement for table `test_grandparents`.
	CreateTableTestGrandparentsStmt = &postgres.CreateStmts{
		GormModel: (*TestGrandparents)(nil),
		Children: []*postgres.CreateStmts{
			&postgres.CreateStmts{
				GormModel: (*TestGrandparentsEmbeddeds)(nil),
				Children: []*postgres.CreateStmts{
					&postgres.CreateStmts{
						GormModel: (*TestGrandparentsEmbeddedsEmbedded2)(nil),
						Children:  []*postgres.CreateStmts{},
					},
				},
			},
		},
	}

	// TestGrandparentsSchema is the go schema for table `test_grandparents`.
	TestGrandparentsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("test_grandparents")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.TestGrandparent)(nil)), "test_grandparents")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(61), "testgrandparent", (*storage.TestGrandparent)(nil)))
		RegisterTable(schema, CreateTableTestGrandparentsStmt)
		return schema
	}()
)

const (
	TestGrandparentsTableName                   = "test_grandparents"
	TestGrandparentsEmbeddedsTableName          = "test_grandparents_embeddeds"
	TestGrandparentsEmbeddedsEmbedded2TableName = "test_grandparents_embeddeds_embedded2"
)

// TestGrandparents holds the Gorm model for Postgres table `test_grandparents`.
type TestGrandparents struct {
	Id         string  `gorm:"column:id;type:varchar;primaryKey"`
	Val        string  `gorm:"column:val;type:varchar"`
	Priority   int64   `gorm:"column:priority;type:bigint"`
	RiskScore  float32 `gorm:"column:riskscore;type:numeric"`
	Serialized []byte  `gorm:"column:serialized;type:bytea"`
}

// TestGrandparentsEmbeddeds holds the Gorm model for Postgres table `test_grandparents_embeddeds`.
type TestGrandparentsEmbeddeds struct {
	TestGrandparentsId  string           `gorm:"column:test_grandparents_id;type:varchar;primaryKey"`
	Idx                 int              `gorm:"column:idx;type:integer;primaryKey;index:testgrandparentsembeddeds_idx,type:btree"`
	Val                 string           `gorm:"column:val;type:varchar"`
	TestGrandparentsRef TestGrandparents `gorm:"foreignKey:test_grandparents_id;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}

// TestGrandparentsEmbeddedsEmbedded2 holds the Gorm model for Postgres table `test_grandparents_embeddeds_embedded2`.
type TestGrandparentsEmbeddedsEmbedded2 struct {
	TestGrandparentsId           string                    `gorm:"column:test_grandparents_id;type:varchar;primaryKey"`
	TestGrandparentsEmbeddedsIdx int                       `gorm:"column:test_grandparents_embeddeds_idx;type:integer;primaryKey"`
	Idx                          int                       `gorm:"column:idx;type:integer;primaryKey;index:testgrandparentsembeddedsembedded2_idx,type:btree"`
	Val                          string                    `gorm:"column:val;type:varchar"`
	TestGrandparentsEmbeddedsRef TestGrandparentsEmbeddeds `gorm:"foreignKey:test_grandparents_id,test_grandparents_embeddeds_idx;references:test_grandparents_id,idx;belongsTo;constraint:OnDelete:CASCADE"`
}
