// Code originally generated by pg-bindings generator.
// Regenerate with central/logimbue/store/postgres/gen.go.

//go:build sql_integration
// +build sql_integration

package n27ton28

import (
	"context"
	"sort"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_27_to_n_28_postgres_log_imbues/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_27_to_n_28_postgres_log_imbues/postgres"
	pghelper "github.com/stackrox/rox/migrator/migrations/postgreshelper"
	"github.com/stackrox/rox/pkg/bolthelper"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/protoutils"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
	bolt "go.etcd.io/bbolt"
)

func TestMigration(t *testing.T) {
	suite.Run(t, new(postgresMigrationSuite))
}

type postgresMigrationSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	ctx         context.Context

	legacyDB   *bolt.DB
	postgresDB *pghelper.TestPostgres
}

var _ suite.TearDownTestSuite = (*postgresMigrationSuite)(nil)

func (s *postgresMigrationSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(env.PostgresDatastoreEnabled.EnvVar(), "true")
	if !env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	var err error
	s.legacyDB, err = bolthelper.NewTemp(s.T().Name() + ".db")
	s.NoError(err)

	s.Require().NoError(err)

	s.ctx = sac.WithAllAccess(context.Background())
	s.postgresDB = pghelper.ForT(s.T(), true)
}

func (s *postgresMigrationSuite) TearDownTest() {
	testutils.TearDownDB(s.legacyDB)
	s.postgresDB.Teardown(s.T())
}

func (s *postgresMigrationSuite) TestLogImbueMigration() {
	newStore := pgStore.New(s.postgresDB.Pool)
	legacyStore := legacy.New(s.legacyDB)

	// Prepare data and write to legacy DB
	var logImbues []*storage.LogImbue

	for i := 0; i < 200; i++ {
		logImbue := &storage.LogImbue{}
		s.NoError(testutils.FullInit(logImbue, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		logImbues = append(logImbues, logImbue)
		s.NoError(legacyStore.Upsert(s.ctx, logImbue))
	}

	// Move
	s.NoError(move(s.postgresDB.GetGormDB(), s.postgresDB.Pool, legacyStore))

	// Verify
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(logImbues), count)

	// Log ids may change, sort by time and compare.
	sort.SliceStable(logImbues, func(i, j int) bool {
		return protoutils.After(logImbues[j].GetTimestamp(), logImbues[i].GetTimestamp())
	})

	fetched, err := newStore.GetAll(s.ctx)
	s.NoError(err)
	sort.SliceStable(fetched, func(i, j int) bool {
		return protoutils.After(fetched[j].GetTimestamp(), fetched[i].GetTimestamp())
	})

	for i, logImbue := range logImbues {
		s.Equal(logImbue.GetLog(), fetched[i].GetLog())
		s.Equal(logImbue.GetTimestamp().GetSeconds(), fetched[i].GetTimestamp().GetSeconds())
	}
}
