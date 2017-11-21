package main

import (
	"github.com/disaster37/check-rancher/model/monitoring"
	"github.com/disaster37/check-rancher/service/monitoring"
	"github.com/disaster37/check-rancher/service/rancher"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/urfave/cli.v1"
	"os"
)

var debug bool
var rancherUrl string
var rancherKey string
var rancherSecret string

version := "develop"

func main() {

	// Logger setting
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)

	// CLI settings
	app := cli.NewApp()
	app.Usage = "Check some usefull state about your Rancher environment"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "rancher-url",
			Usage:       "The Rancher base URL",
			Destination: &rancherUrl,
		},
		cli.StringFlag{
			Name:        "rancher-key",
			Usage:       "The Rancher API key",
			Destination: &rancherKey,
		},
		cli.StringFlag{
			Name:        "rancher-secret",
			Usage:       "The Rancher API secret",
			Destination: &rancherSecret,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Display debug output",
			Destination: &debug,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "check-stack",
			Usage: "Check the stack state",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "stack-name",
					Usage: "The stack name you should to check",
				},
			},
			Action: checkStack,
		},
		{
			Name:  "check-hosts-project",
			Usage: "Check the hosts state in given project",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "project-name",
					Usage: "The project name you should to check hosts",
				},
			},
			Action: checkHostsProject,
		},
		{
			Name:  "check-certificates",
			Usage: "Check all certificates validity",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "warning-days",
					Usage: "The number of days before certificate expire to fire warning",
				},
			},
			Action: checkCertificates,
		},
	}

	app.Run(os.Args)

}

// Check the global parameter
func manageGlobalParameters() error {
	if debug == true {
		log.SetLevel(log.DebugLevel)
	}

	if rancherUrl == "" {
		return errors.New("You must set --rancher-url parameter")
	}

	if rancherKey == "" {
		return errors.New("You must set --rancher-key parameter")
	}
	if rancherSecret == "" {
		return errors.New("You must set --rancher-secret parameter")
	}

	return nil
}

// Perform the check stack's state on given stack
func checkStack(c *cli.Context) error {

	monitoringData := modelMonitoring.NewMonitoring()
	monitoringData.SetStatus(modelMonitoring.STATUS_UNKNOWN)

	// Check global parameters
	err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), modelMonitoring.STATUS_UNKNOWN)
	}

	// Check current parameters
	if c.String("stack-name") == "" {
		return cli.NewExitError("You must set --stack-name parameter", modelMonitoring.STATUS_UNKNOWN)
	}

	// Get Rancher connection
	rancherClient, err := rancherService.NewClient(rancherUrl, rancherKey, rancherSecret)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Error appear when try to connect on Rancher API: %v", err))
		monitoringData.ToSdtOut()
	}

	// Load Rancher stack and all data associated
	stack, err := rancherClient.LoadStackByName(c.String("stack-name"))
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to load stack %s: %v", c.String("stack-name"), err))
		monitoringData.ToSdtOut()
	}

	// Check the stack state
	monitoringDataFinal, err := monitoringService.CheckStack(stack)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to check stack state %s: %v", c.String("stack-name"), err))
		monitoringData.ToSdtOut()
	}

	monitoringDataFinal.ToSdtOut()

	return nil

}

// Perform the check hosts'state on given project
func checkHostsProject(c *cli.Context) error {

	monitoringData := modelMonitoring.NewMonitoring()
	monitoringData.SetStatus(modelMonitoring.STATUS_UNKNOWN)

	// Check global parameters
	err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), modelMonitoring.STATUS_UNKNOWN)
	}

	// Check current parameters
	if c.String("environment-name") == "" {
		return cli.NewExitError("You must set --environment-name parameter", modelMonitoring.STATUS_UNKNOWN)
	}

	// Get Rancher connection
	rancherClient, err := rancherService.NewClient(rancherUrl, rancherKey, rancherSecret)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Error appear when try to connect on Rancher API: %v", err))
		monitoringData.ToSdtOut()
	}

	// Load Rancher hosts associated to environment
	project, err := rancherClient.LoadProjectByName(c.String("environment-name"))
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to load environment %s: %v", c.String("environment-name"), err))
		monitoringData.ToSdtOut()
	}

	// Check the hosts state
	monitoringDataFinal, err := monitoringService.CheckHostsProject(project)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to check environment state %s: %v", c.String("environment-name"), err))
		monitoringData.ToSdtOut()
	}

	monitoringDataFinal.ToSdtOut()

	return nil

}

// Perform the check certificates validity
func checkCertificates(c *cli.Context) error {

	monitoringData := modelMonitoring.NewMonitoring()
	monitoringData.SetStatus(modelMonitoring.STATUS_UNKNOWN)

	// Check global parameters
	err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), modelMonitoring.STATUS_UNKNOWN)
	}

	// Get Rancher connection
	rancherClient, err := rancherService.NewClient(rancherUrl, rancherKey, rancherSecret)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Error appear when try to connect on Rancher API: %v", err))
		monitoringData.ToSdtOut()
	}

	// Load Rancher certificates
	certificates, err := rancherClient.GetCertificates()
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to load certificates: %v", err))
		monitoringData.ToSdtOut()
	}

	monitoringDataFinal, err := monitoringService.CheckCertificates(certificates, c.Int("warning-days"))
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to check certificates: %v", err))
		monitoringData.ToSdtOut()
	}

	monitoringDataFinal.ToSdtOut()

	return nil

}
