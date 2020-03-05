package main

import (
	"os"
	"sort"

	"github.com/disaster37/check-rancher/v1/checkrancher"
	"github.com/disaster37/go-nagios"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func run(args []string) error {

	// Logger setting
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)

	// CLI settings
	app := cli.NewApp()
	app.Usage = "Check some usefull state about your Rancher project"
	app.Version = "develop"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: "Load configuration from `FILE`",
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "url",
			Usage:   "The Rancher base URL",
			EnvVars: []string{"RANCHER_URL"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "rancher-key",
			Usage:   "The Rancher API key",
			EnvVars: []string{"RANCHER_KEY"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "rancher-secret",
			Usage:   "The Rancher API secret",
			EnvVars: []string{"RANCHER_SECRET"},
		}),
		&cli.StringFlag{
			Name:  "project-name",
			Usage: "The project name you should to check certificates",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Display debug output",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:  "check-stack",
			Usage: "Check the stack state",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "stack-name",
					Usage: "The stack name you should to check",
				},
			},
			Action: checkrancher.CheckStack,
		},
		{
			Name:   "check-hosts",
			Usage:  "Check the hosts state in given project",
			Action: checkrancher.CheckHosts,
		},
		{
			Name:  "check-certificates",
			Usage: "Check all certificates validity on project",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "warning-days",
					Usage: "The number of days before certificate expire to fire warning",
				},
			},
			Action: checkrancher.CheckCertificates,
		},
	}

	app.Before = func(c *cli.Context) error {

		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		if c.String("config") != "" {
			before := altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("config"))
			return before(c)
		}
		return nil
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(args)
	return err

}

func main() {
	err := run(os.Args)
	if err != nil {
		monitoringData := nagiosPlugin.NewMonitoring()
		monitoringData.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		monitoringData.AddMessage("Error appear during check: %s", err)
		monitoringData.ToSdtOut()
	}
}

/*
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
	if c.String("project-name") == "" {
		return cli.NewExitError("You must set --project-name parameter", modelMonitoring.STATUS_UNKNOWN)
	}
	if c.String("stack-name") == "" {
		return cli.NewExitError("You must set --stack-name parameter", modelMonitoring.STATUS_UNKNOWN)
	}

	// Get Rancher connection
	rancherClient, err := rancherService.NewClient(rancherUrl, rancherKey, rancherSecret)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Error appear when try to connect on Rancher API: %v", err))
		monitoringData.ToSdtOut()
	}

	// Load Rancher project
	project, err := rancherClient.FindProjectByName(c.String("project-name"))
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to load environment %s: %v", c.String("project-name"), err))
		monitoringData.ToSdtOut()
	}
	if project == nil {
		monitoringData.AddMessage(fmt.Sprintf("Project %s not found (project not exist, or you are no access to this project)", c.String("project-name")))
		monitoringData.ToSdtOut()
	}

	// Load Rancher stack and all data associated
	stack, err := rancherClient.LoadStackByNameOnProject(c.String("stack-name"), project)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to load stack %s on project %s: %v", c.String("stack-name"), c.String("project-name"), err))
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
	if c.String("project-name") == "" {
		return cli.NewExitError("You must set --project-name parameter", modelMonitoring.STATUS_UNKNOWN)
	}

	// Get Rancher connection
	rancherClient, err := rancherService.NewClient(rancherUrl, rancherKey, rancherSecret)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Error appear when try to connect on Rancher API: %v", err))
		monitoringData.ToSdtOut()
	}

	// Load Rancher project
	project, err := rancherClient.FindProjectByName(c.String("project-name"))
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to load environment %s: %v", c.String("project-name"), err))
		monitoringData.ToSdtOut()
	}
	if project == nil {
		monitoringData.AddMessage(fmt.Sprintf("Project %s not found (project not exist, or you are no access to this project)", c.String("project-name")))
		monitoringData.ToSdtOut()
	}

	// Load Rancher hosts associated to project
	hosts, err := rancherClient.FindHostsByProjectLink(project)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to load hosts on project %s: %v", c.String("project-name"), err))
		monitoringData.ToSdtOut()
	}

	// Check the hosts state
	monitoringDataFinal, err := monitoringService.CheckHosts(hosts)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to check hosts state: %v", err))
		monitoringData.ToSdtOut()
	}

	monitoringDataFinal.ToSdtOut()

	return nil

}

// Perform the check certificates validity
func checkCertificatesProject(c *cli.Context) error {

	monitoringData := modelMonitoring.NewMonitoring()
	monitoringData.SetStatus(modelMonitoring.STATUS_UNKNOWN)

	// Check global parameters
	err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), modelMonitoring.STATUS_UNKNOWN)
	}

	// Check current parameters
	if c.String("project-name") == "" {
		return cli.NewExitError("You must set --project-name parameter", modelMonitoring.STATUS_UNKNOWN)
	}

	// Get Rancher connection
	rancherClient, err := rancherService.NewClient(rancherUrl, rancherKey, rancherSecret)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Error appear when try to connect on Rancher API: %v", err))
		monitoringData.ToSdtOut()
	}

	// Load Rancher project
	project, err := rancherClient.FindProjectByName(c.String("project-name"))
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to load project %s: %v", c.String("project-name"), err))
		monitoringData.ToSdtOut()
	}
	if project == nil {
		monitoringData.AddMessage(fmt.Sprintf("Project %s not found (project not exist, or you are no access to this project)", c.String("project-name")))
		monitoringData.ToSdtOut()
	}

	// Load Rancher certificates associated to project
	certificates, err := rancherClient.FindCertificatesByProjectLink(project)
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to load certificates on project %s: %v", c.String("project-name"), err))
		monitoringData.ToSdtOut()
	}

	// Check certificates validities
	monitoringDataFinal, err := monitoringService.CheckCertificates(certificates, c.Int("warning-days"))
	if err != nil {
		monitoringData.AddMessage(fmt.Sprintf("Somethink wrong when try to check certificates: %v", err))
		monitoringData.ToSdtOut()
	}

	monitoringDataFinal.ToSdtOut()

	return nil

}

*/
