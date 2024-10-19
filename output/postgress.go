package output

import (
	"fmt"
	"os"

	"github.com/major1201/kubetrack/config"
	"github.com/major1201/kubetrack/gormutils"
	"github.com/major1201/kubetrack/log"
	"github.com/robfig/cron/v3"
)

type PostgresOutput struct {
	conf *config.OutputPostgres
}

func NewPostgresOutput(conf *config.OutputPostgres) *PostgresOutput {
	if conf == nil {
		return nil
	}
	out := &PostgresOutput{
		conf: conf,
	}
	out.initDB()
	out.migrate()
	out.initCleanupJob()

	return out
}

func (lo *PostgresOutput) Name() string {
	return "postgres"
}

func (lo *PostgresOutput) Write(out OutputStruct) error {
	ev := &Events{
		EventTime: out.EventTime,
		Source:    string(out.Source),
		EventType: string(out.EventType),

		APIVersion: out.ObjectRef.APIVersion,
		Kind:       out.ObjectRef.Kind,
		Namespace:  out.ObjectRef.Namespace,
		Name:       out.ObjectRef.Name,
		UID:        string(out.ObjectRef.UID),

		Fields:  gormutils.MustToJsonb(out.Fields),
		Message: out.Message,

		Object:    gormutils.MustToJsonb(out.Object),
		Diff:      out.Diff,
		JsonPatch: gormutils.MustToJsonb(out.JsonPatch),
	}
	return gormutils.Save(ev)
}

func (lo *PostgresOutput) initDB() {
	os.Setenv(gormutils.EnvDBDriver, "postgres")
	os.Setenv(gormutils.EnvDBConnection, lo.conf.DSN)
	// try db connection
	sqlDB, err := gormutils.GetDB().DB()
	if err != nil {
		log.L.Error(err, "get sql db failed")
		os.Exit(1)
	}
	if err = sqlDB.Ping(); err != nil {
		log.L.Error(err, "db connect failed")
		os.Exit(1)
	}
}

func (lo *PostgresOutput) migrate() {
	log.L.Info("migrating postgres")

	if err := gormutils.GetDB().AutoMigrate(&Events{}); err != nil {
		log.L.Error(err, "migrate error")
		os.Exit(1)
	}
}

func (lo *PostgresOutput) initCleanupJob() {
	if lo.conf.TTLDays <= 0 {
		return
	}

	cj := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	if _, err := cj.AddFunc("* * * * *", lo.doCleanupJob); err != nil {
		panic(err)
	}
	cj.Start()
	log.L.Info("postgres cleanup job started, will run every hour", "ttl", lo.conf.TTLDays)
}

func (lo *PostgresOutput) doCleanupJob() {
	log.L.Info("running cleanup job", "ttlDays", lo.conf.TTLDays)
	if err := gormutils.Delete(&Events{}, fmt.Sprintf("created_at < now() - INTERVAL '%d days'", lo.conf.TTLDays)); err != nil {
		log.L.Error(err, "cron: delete data failed")
	}
}
