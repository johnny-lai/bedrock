package service

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/johnny-lai/bedrock"
)

type Config struct {
	SvcHost string
}

type Service struct {
	airbrake bedrock.AirbrakeService
	newrelic bedrock.NewRelicService
	dbsvc    bedrock.GormService
	config   Config
}

func (s *Service) Configure(app *bedrock.Application) error {
	if err := s.airbrake.Configure(app); err != nil {
		return err
	}

	if err := s.newrelic.Configure(app); err != nil {
		return err
	}

	if err := s.dbsvc.Configure(app); err != nil {
		return err
	}

	if err := app.BindConfig(&s.config); err != nil {
		return err
	}

	return nil
}

func (s *Service) Migrate(app *bedrock.Application) error {
	db, err := s.dbsvc.Db()
	if err != nil {
		return err
	}
	db.SingularTable(true)
	return nil
}

func (s *Service) Build(app *bedrock.Application) error {
	if err := s.airbrake.Build(app); err != nil {
		return err
	}

	if err := s.newrelic.Build(app); err != nil {
		return err
	}

	if err := s.dbsvc.Build(app); err != nil {
		return err
	}

	db, err := s.dbsvc.Db()
	if err != nil {
		return err
	}
	db.SingularTable(true)

	r := app.Engine
	r.GET("/health", s.dbsvc.HealthHandler(app))
	r.GET("/panic", s.airbrake.PanicHandler(app))

	return nil
}

func (s *Service) Run(app *bedrock.Application) error {
	app.Engine.Run(s.config.SvcHost)

	return nil
}
