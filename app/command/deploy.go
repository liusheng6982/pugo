package command

import (
	"errors"

	"github.com/codegangsta/cli"
	"github.com/go-xiaohei/pugo-static/app/builder"
	"github.com/go-xiaohei/pugo-static/app/deploy"
	"gopkg.in/inconshreveable/log15.v2"
)

func Deploy(opt *builder.BuildOption) cli.Command {
	return cli.Command{
		Name:     "deploy",
		Usage:    "deploy site to other platform",
		HideHelp: true,
		Flags: []cli.Flag{
			destFlag,
			themeFlag,
			debugFlag,
			watchFlag,
		},
		Action: deploySite(opt),
		Before: setDebugMode,
	}
}

func deploySite(opt *builder.BuildOption) func(ctx *cli.Context) {
	// build action
	return func(ctx *cli.Context) {

		if iniFile == nil {
			log15.Crit("Deploy.Fail", "error", errors.New("please add conf.ini to set deploy options"))
		}
		deployer, err := deploy.New(iniFile)
		if err != nil {
			log15.Crit("Deploy.Fail", "error", err.Error())
		}

		// real deploy action, in builder hook
		afterFunc := func(b *builder.Builder, c *builder.Context) error {
			if b.IsWatching() || isWatch || ctx.Bool("watch") {
				return deployer.RunAsync(b, c)
			}
			return deployer.Run(b, c)
		}

		// add hook to opt
		opt.After(afterFunc)

		// run build site
		buildSite(opt, false)(ctx)
	}
}
