package main

import (
	"fmt"
	"os"

	//"github.com/davecgh/go-spew/spew"
	"github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"

	"github.com/appnexus/ankh/internal/ankh"
	"github.com/appnexus/ankh/internal/helm"
	"github.com/appnexus/ankh/internal/kubectl"
)

var log = logrus.New()

func main() {
	app := cli.App("ankh", "AppNexus Kubernetes Helper")
	app.Spec = "[-v]"

	var (
		verbose = app.BoolOpt("v verbose", false, "Verbose debug mode")
	)

	formatter := logrus.TextFormatter{
		DisableTimestamp: true,
	}
	log.Out = os.Stdout
	log.Formatter = &formatter
	if *verbose {
		log.Level = logrus.DebugLevel
	} else {
		log.Level = logrus.InfoLevel
	}

	app.Command("apply", "Deploy an ankh file to a kubernetes cluster", func(cmd *cli.Cmd) {

		cmd.Spec = "[-f] [--dry-run]"

		var (
			filename = cmd.StringOpt("f filename", "ankh.yaml", "Config file name")
			dryRun   = cmd.BoolOpt("dry-run", false, "Perform a dry-run and don't actually apply anything to a cluster")
		)

		cmd.Action = func() {
			ankhConfig, err := ankh.GetAnkhConfig()
			check(err)

			config, err := ankh.ProcessAnkhFile(filename)
			check(err)

			helmOutput, err := helm.Template(log, config, ankhConfig)
			check(err)

			action := kubectl.Apply
			log.Info("starting kubectl")
			kubectlOutput, err := kubectl.Execute(log, action, dryRun, helmOutput, config, ankhConfig)
			check(err)

			fmt.Println(kubectlOutput)

			log.Info(helmOutput)
			log.Info("complete")
			os.Exit(0)
		}
	})

	app.Command("template", "Output the results of templating an ankh file", func(cmd *cli.Cmd) {

		cmd.Spec = "[-f]"

		var (
			filename = cmd.StringOpt("f filename", "ankh.yaml", "Config file name")
		)

		cmd.Action = func() {
			ankhConfig, err := ankh.GetAnkhConfig()
			check(err)

			config, err := ankh.ProcessAnkhFile(filename)
			check(err)

			log.Info("starting helm template")
			helmOutput, err := helm.Template(log, config, ankhConfig)
			check(err)

			fmt.Println(helmOutput)
			log.Info("complete")
			os.Exit(0)
		}
	})

	app.Run(os.Args)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}