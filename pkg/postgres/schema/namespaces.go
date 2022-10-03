// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"fmt"
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableNamespacesStmt holds the create statement for table `namespaces`.
	CreateTableNamespacesStmt = &postgres.CreateStmts{
		GormModel: (*Namespaces)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// NamespacesSchema is the go schema for table `namespaces`.
	NamespacesSchema = func() *walker.Schema {
		schema := GetSchemaForTable("namespaces")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.NamespaceMetadata)(nil)), "namespaces")
		referencedSchemas := map[string]*walker.Schema{
			"storage.Cluster": ClustersSchema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_NAMESPACES, "namespacemetadata", (*storage.NamespaceMetadata)(nil)))
		schema.SetSearchScope([]v1.SearchCategory{
			v1.SearchCategory_IMAGE_VULNERABILITIES,
			v1.SearchCategory_COMPONENT_VULN_EDGE,
			v1.SearchCategory_IMAGE_COMPONENTS,
			v1.SearchCategory_IMAGE_COMPONENT_EDGE,
			v1.SearchCategory_IMAGE_VULN_EDGE,
			v1.SearchCategory_IMAGES,
			v1.SearchCategory_DEPLOYMENTS,
			v1.SearchCategory_NAMESPACES,
			v1.SearchCategory_CLUSTERS,
		}...)
		RegisterTable(schema, CreateTableNamespacesStmt)
		return schema
	}()
)

const (
	NamespacesTableName = "namespaces"
)

// Namespaces holds the Gorm model for Postgres table `namespaces`.
type Namespaces struct {
	Id          string            `gorm:"column:id;type:varchar;primaryKey"`
	Name        string            `gorm:"column:name;type:varchar"`
	ClusterId   string            `gorm:"column:clusterid;type:varchar"`
	ClusterName string            `gorm:"column:clustername;type:varchar"`
	Labels      map[string]string `gorm:"column:labels;type:jsonb"`
	Annotations map[string]string `gorm:"column:annotations;type:jsonb"`
	Serialized  []byte            `gorm:"column:serialized;type:bytea"`
}
