package output

import (
	"os"

	"github.com/major1201/kubetrack/config"
	"github.com/major1201/kubetrack/gormutils"
	"github.com/major1201/kubetrack/log"
)

type MysqlOutput struct {
	conf *config.OutputMysql
}

func NewMysqlOutput(conf *config.OutputMysql) *MysqlOutput {
	if conf == nil {
		return nil
	}
	res := &MysqlOutput{
		conf: conf,
	}
	res.initDB()
	res.migrate()

	return res
}

func (lo *MysqlOutput) Name() string {
	return "mysql"
}

func (lo *MysqlOutput) Write(out OutputStruct) error {
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
