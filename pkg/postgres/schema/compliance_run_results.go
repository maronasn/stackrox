// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"
	"time"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableComplianceRunResultsStmt holds the create statement for table `compliance_run_results`.
	CreateTableComplianceRunResultsStmt = &postgres.CreateStmts{
		GormModel: (*ComplianceRunResults)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ComplianceRunResultsSchema is the go schema for table `compliance_run_results`.
	ComplianceRunResultsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("compliance_run_results")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ComplianceRunResults)(nil)), "compliance_run_results")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_COMPLIANCE_RESULTS, "compliancerunresults", (*storage.ComplianceRunResults)(nil)))
		RegisterTable(schema, CreateTableComplianceRunResultsStmt)
		return schema
	}()
)

const (
	ComplianceRunResultsTableName = "compliance_run_results"
)

// ComplianceRunResults holds the Gorm model for Postgres table `compliance_run_results`.
type ComplianceRunResults struct {
	RunMetadataRunId           string     `gorm:"column:runmetadata_runid;type:varchar;primaryKey"`
	RunMetadataStandardId      string     `gorm:"column:runmetadata_standardid;type:varchar"`
	RunMetadataClusterId       string     `gorm:"column:runmetadata_clusterid;type:varchar"`
	RunMetadataFinishTimestamp *time.Time `gorm:"column:runmetadata_finishtimestamp;type:timestamp"`
	Serialized                 []byte     `gorm:"column:serialized;type:bytea"`
}
