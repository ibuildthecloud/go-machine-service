package main

import (
	"os"

	"github.com/rancher/go-machine-service/controller"
	"github.com/rancher/go-machine-service/logging"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	GITCOMMIT = "HEAD"
)

var logger = logging.Logger()

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "Kube config for accessing kubernetes cluster",
			EnvVar: "KUBECONFIG",
		},
	}

	app.Action = func(c *cli.Context) error {
		return run(c.String("config"))
	}

	app.ExitErrHandler = func(c *cli.Context, err error) {
		logrus.Fatal(err)
	}

	app.Run(os.Args)
}

func run(kubeConfigFile string) error {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
	if err != nil {
		return err
	}

	management, err := config.NewManagementContext(*kubeConfig)
	if err != nil {
		return err
	}

	controller.Register(management)

	return management.StartAndWait()
}

//func main() {
//	processCmdLineFlags()
//	logrus.SetLevel(logrus.DebugLevel)
//
//	logger.WithField("gitcommit", GITCOMMIT).Info("Starting go-machine-service...")
//
//	apiURL := os.Getenv("CATTLE_URL")
//	accessKey := os.Getenv("CATTLE_ACCESS_KEY")
//	secretKey := os.Getenv("CATTLE_SECRET_KEY")
//
//	ready := make(chan bool, 2)
//	done := make(chan error)
//
//	go func() {
//		eventHandlers := map[string]events.EventHandler{
//			"machinedriver.reactivate": handlers.ActivateDriver,
//			"machinedriver.activate":   handlers.ActivateDriver,
//			"machinedriver.update":     handlers.ActivateDriver,
//			"machinedriver.error":      handlers.ErrorDriver,
//			"machinedriver.deactivate": handlers.DeactivateDriver,
//			"machinedriver.remove":     handlers.RemoveDriver,
//			"ping":                     handlers.PingNoOp,
//		}
//
//		router, err := events.NewEventRouter("machine-service", 2000, apiURL, accessKey, secretKey,
//			nil, eventHandlers, "machineDriver", 250, events.DefaultPingConfig)
//		if err == nil {
//			err = router.Start(ready)
//		}
//		done <- err
//	}()
//
//	go func() {
//		eventHandlers := map[string]events.EventHandler{
//			"host.provision": handlers.CreateMachineAndActivateMachine,
//			"host.remove":    handlers.PurgeMachine,
//			"ping":           handlers.PingNoOp,
//		}
//
//		router, err := events.NewEventRouter("machine-service", 2000, apiURL, accessKey, secretKey,
//			nil, eventHandlers, "host", 250, events.DefaultPingConfig)
//		if err == nil {
//			err = router.Start(ready)
//		}
//		done <- err
//	}()
//
//	go func() {
//		// Can not remove this as nothing will delete the handler entries
//		eventHandlers := map[string]events.EventHandler{
//			"ping": handlers.PingNoOp,
//		}
//
//		router, err := events.NewEventRouter("machine-service", 2000, apiURL, accessKey, secretKey,
//			nil, eventHandlers, "agent", 5, events.DefaultPingConfig)
//		if err == nil {
//			err = router.Start(ready)
//		}
//		done <- err
//	}()
//
//	go func() {
//		logger.Infof("Waiting for handler registration (1/2)")
//		<-ready
//		logger.Infof("Waiting for handler registration (2/2)")
//		<-ready
//		if err := dynamic.ReactivateOldDrivers(); err != nil {
//			logger.Fatalf("Error reactivating old drivers: %v", err)
//		}
//		if err := dynamic.DownloadAllDrivers(); err != nil {
//			logger.Fatalf("Error updating drivers: %v", err)
//		}
//	}()
//
//	err := <-done
//	if err == nil {
//		logger.Infof("Exiting go-machine-service")
//	} else {
//		logger.Fatalf("Exiting go-machine-service: %v", err)
//	}
//}
//
//func processCmdLineFlags() {
//	// Define command line flags
//	version := flag.Bool("v", false, "read the version of the go-machine-service")
//	flag.Parse()
//	if *version {
//		fmt.Printf("go-machine-service\t gitcommit=%s\n", GITCOMMIT)
//		os.Exit(0)
//	}
//}
