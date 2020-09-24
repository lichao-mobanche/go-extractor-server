package command

import (
	"github.com/cfhamlet/os-rq-pod/pkg/runner"
	"github.com/cfhamlet/os-rq-pod/pkg/command"
	"github.com/lichao-mobanche/go-extractor-server/pkg/config"
	"github.com/lichao-mobanche/go-extractor-server/extractor-server/global"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func init() {
	Root.AddCommand(command.NewRunCommand("db-req-pod", run))
}
func run(conf *viper.Viper) {

	newConfig := func() (*viper.Viper, error) {
		err := config.LoadConfig(conf, global.EnvPrefix, global.DefaultConfig)
		return conf, err
	}

	var r *runner.Runner

	app := fx.New(
		fx.Provide(
			runner.New,
			newConfig,
		),
		fx.Invoke(),
		fx.Populate(&r),
	)
	r.Run(app)
}