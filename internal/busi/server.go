package busi

import (
	"aggregate-task/pkg/models/evmmodel"
	"aggregate-task/pkg/utils"
	"context"

	"github.com/sirupsen/logrus"
)

type Server struct {
	Ctx context.Context
	Cf  utils.TomlConfig
}

func NewServer(ctx context.Context) *Server {
	return &Server{Ctx: ctx}
}

func (s *Server) initconfig() {
	if err := utils.InitConfFile(Flags.Config, &s.Cf); err != nil {
		logrus.Fatalf("Load configuration file err: %v", err)
	}

	if s.Cf.AggregateTask.FinalityEpoch > 900 {
		logrus.Fatal("the finality should be less than or equal to 900, https://docs.filecoin.io/reference/reference/glossary/#finality")
	}

	utils.EngineGroup = utils.NewEngineGroup(s.Ctx, &[]utils.EngineInfo{{utils.DBExtract, s.Cf.AggregateTask.NotifyDB, nil}, {utils.DBOBTask, s.Cf.AggregateTask.ObservatoryDB, evmmodel.Tables}})
}

func (s *Server) setLogTimeformat() {
	timeFormater := new(logrus.TextFormatter)
	timeFormater.FullTimestamp = true
	logrus.SetFormatter(timeFormater)
}

func (s *Server) Start() {
	s.initconfig()
	s.setLogTimeformat()

	go HttpServerStart(s.Cf.AggregateTask.Addr)

	AggregateTaskStart(s.Ctx, s.Cf.AggregateTask.FinalityEpoch, s.Cf.Task.Name, s.Cf.Task.Depend)
}
