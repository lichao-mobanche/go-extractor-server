package command

import (
	"github.com/cfhamlet/os-rq-pod/pkg/command"
	"github.com/cfhamlet/os-rq-pod/pkg/ginserv"
	"github.com/cfhamlet/os-rq-pod/pkg/runner"
	"github.com/lichao-mobanche/go-extractor-server/app/controllers"
	"github.com/lichao-mobanche/go-extractor-server/app/routers"
	"github.com/lichao-mobanche/go-extractor-server/pkg/config"
	"github.com/lichao-mobanche/go-extractor-server/server/exqueue"
	"github.com/lichao-mobanche/go-extractor-server/server/global"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func init() {
	Root.AddCommand(command.NewRunCommand("extractor-server", run))
}
func run(conf *viper.Viper) {

	newConfig := func() (*viper.Viper, error) {
		err := config.LoadConfig(conf, global.EnvPrefix, global.DefaultConfig)
		return conf, err
	}

	var r *runner.Runner

	exqueueGo := func(lc fx.Lifecycle, e *exqueue.ExQueue, r *runner.Runner) {
		runner.ServWait(lc, e, r)
	}

	app := fx.New(
		fx.Provide(
			runner.New,
			newConfig,
			exqueue.New,
			ginserv.NewEngine,
			ginserv.NewServer,
			ginserv.NewAPIGroup,
			controllers.NewExQueueController,
		),
		fx.Invoke(
			ginserv.LoadGlobalMiddlewares,
			exqueueGo,
			routers.RouteExQueueCtrl,
			runner.HTTPServerLifecycle,
		),
		fx.Populate(&r),
	)
	r.Run(app)
}
