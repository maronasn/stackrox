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
	// CreateTableSecretsStmt holds the create statement for table `secrets`.
	CreateTableSecretsStmt = &postgres.CreateStmts{
		GormModel: (*Secrets)(nil),
		Children: []*postgres.CreateStmts{
			&postgres.CreateStmts{
				GormModel: (*SecretsFiles)(nil),
				Children: []*postgres.CreateStmts{
					&postgres.CreateStmts{
						GormModel: (*SecretsFilesRegistries)(nil),
						Children:  []*postgres.CreateStmts{},
					},
				},
			},
		},
	}

	// SecretsSchema is the go schema for table `secrets`.
	SecretsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("secrets")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.Secret)(nil)), "secrets")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_SECRETS, "secret", (*storage.Secret)(nil)))
		RegisterTable(schema, CreateTableSecretsStmt)
		return schema
	}()
)

const (
	SecretsTableName                = "secrets"
	SecretsFilesTableName           = "secrets_files"
	SecretsFilesRegistriesTableName = "secrets_files_registries"
)

// Secrets holds the Gorm model for Postgres table `secrets`.
type Secrets struct {
	Id          string     `gorm:"column:id;type:varchar;primaryKey"`
	Name        string     `gorm:"column:name;type:varchar"`
	ClusterId   string     `gorm:"column:clusterid;type:varchar"`
	ClusterName string     `gorm:"column:clustername;type:varchar"`
	Namespace   string     `gorm:"column:namespace;type:varchar"`
	CreatedAt   *time.Time `gorm:"column:createdat;type:timestamp"`
	Serialized  []byte     `gorm:"column:serialized;type:bytea"`
}

// SecretsFiles holds the Gorm model for Postgres table `secrets_files`.
type SecretsFiles struct {
	SecretsId   string             `gorm:"column:secrets_id;type:varchar;primaryKey"`
	Idx         int                `gorm:"column:idx;type:integer;primaryKey;index:secretsfiles_idx,type:btree"`
	Type        storage.SecretType `gorm:"column:type;type:integer"`
	CertEndDate *time.Time         `gorm:"column:cert_enddate;type:timestamp"`
	SecretsRef  Secrets            `gorm:"foreignKey:secrets_id;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}

// SecretsFilesRegistries holds the Gorm model for Postgres table `secrets_files_registries`.
type SecretsFilesRegistries struct {
	SecretsId       string       `gorm:"column:secrets_id;type:varchar;primaryKey"`
	SecretsFilesIdx int          `gorm:"column:secrets_files_idx;type:integer;primaryKey"`
	Idx             int          `gorm:"column:idx;type:integer;primaryKey;index:secretsfilesregistries_idx,type:btree"`
	Name            string       `gorm:"column:name;type:varchar"`
	SecretsFilesRef SecretsFiles `gorm:"foreignKey:secrets_id,secrets_files_idx;references:secrets_id,idx;belongsTo;constraint:OnDelete:CASCADE"`
}
