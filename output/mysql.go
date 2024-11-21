package output

import (
	"os"

	"github.com/major1201/kubetrack/config"
	"github.com/major1201/kubetrack/gormutils"
	"github.com/major1201/kubetrack/log"
	"github.com/robfig/cron/v3"
)

type MysqlOutput struct {
	ktconfig *config.KubeTrackConfiguration
	conf     *config.OutputMysql
}

func NewMysqlOutput(ktconfig *config.KubeTrackConfiguration, conf *config.OutputMysql) *MysqlOutput {
	if conf == nil {
		return nil
	}
	res := &MysqlOutput{
		ktconfig: ktconfig,
		conf:     conf,
	}
	res.initDB()
	res.migrate()
	res.initCleanupJob()

	return res
}

func (lo *MysqlOutput) Name() string {
	return "mysql"
}

func (lo *MysqlOutput) Write(out OutputStruct) error {
	ev := &Events{
		Cluster:   lo.ktconfig.Cluster,
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

func (lo *MysqlOutput) initDB() {
	os.Setenv(gormutils.EnvDBDriver, "mysql")
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

func (lo *MysqlOutput) migrate() {
	log.L.Info("migrating mysql")

	if err := gormutils.GetDB().AutoMigrate(&Events{}); err != nil {
		log.L.Error(err, "migrate error")
		os.Exit(1)
	}
}

func (lo *MysqlOutput) initCleanupJob() {
	if lo.conf.TTLDays <= 0 {
		return
	}

	cj := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger), cron.Recover(cron.DefaultLogger)))
	if _, err := cj.AddFunc("0 * * * *", lo.doCleanupJob); err != nil {
		panic(err)
	}
	cj.Start()
	log.L.Info("mysql cleanup job started, will run every hour", "ttl", lo.conf.TTLDays)
}

func (lo *MysqlOutput) doCleanupJob() {
	log.L.Info("running cleanup job", "ttlDays", lo.conf.TTLDays)
	if err := gormutils.Delete(&Events{}, "created_at < now() - interval ? day", lo.conf.TTLDays); err != nil {
		log.L.Error(err, "cron: delete data failed")
	}
}
