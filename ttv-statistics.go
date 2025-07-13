package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"ttv-statistics/api"
	"ttv-statistics/helixclient"
)

const (
	hostFlagName         string = "host"
	hostHelpText         string = "the host address where the API should be hosted"
	clientIDFlagName     string = "client-id"
	clientSecretFlagName string = "client-secret"
	clientIDHelpText     string = "the client ID used to access Twitch helix API"
	clientSecretHelpText string = "the client secret used to access Twitch helix API"
	helixHostFlagName    string = "helix-host"
	helixHostHelpText    string = "the host address of the twitch helix API "
)

var (
	stringFlags = []stringFlag{
		{
			ptr:          &api.Host,
			flagName:     hostFlagName,
			defaultValue: "",
			helpText:     hostHelpText,
		},
		{
			ptr:          &helixclient.ClientID,
			flagName:     clientIDFlagName,
			defaultValue: "",
			helpText:     clientIDHelpText,
		},
		{
			ptr:          &helixclient.ClientSecret,
			flagName:     clientSecretFlagName,
			defaultValue: "",
			helpText:     clientSecretHelpText,
		},
		{
			ptr:          &helixclient.HelixHost,
			flagName:     helixHostFlagName,
			defaultValue: "",
			helpText:     helixHostHelpText,
		},
	}
)

type stringFlag struct {
	ptr          *string
	flagName     string
	defaultValue string
	helpText     string
}

func parseFlags() error {

	for _, stringFlag := range stringFlags {
		flag.StringVar(stringFlag.ptr, stringFlag.flagName, stringFlag.defaultValue, stringFlag.helpText)
	}

	flag.Parse()

	missingFlags := []string{}

	for _, stringFlag := range stringFlags {
		if *stringFlag.ptr == "" {
			missingFlags = append(missingFlags, fmt.Sprintf("--%s", stringFlag.flagName))
		}
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("missing required flags: %s", strings.Join(missingFlags, ", "))
	}

	return nil

}

func runServerAndAwaitShutdown() error {

	server := api.NewTTVStatisticsServer()
	server.Run()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	return server.ShutDownServer(context.Background())

}

func main() {

	parseFlagsError := parseFlags()
	if parseFlagsError != nil {
		log.SetFlags(0)
		log.Printf("Startup Error: %v", parseFlagsError)
		flag.Usage()
		os.Exit(1)
	}

	clientAuthError := helixclient.InitHelixClientAuth(context.Background())
	if clientAuthError != nil {
		log.Printf("Failed to authenticate with TwithTV API. Error: %v", clientAuthError)
	}

	serverShutdownError := runServerAndAwaitShutdown()
	if serverShutdownError != nil {
		log.Printf("Server failed to shutdown gracefully. Error: %v", serverShutdownError)
		os.Exit(1)
	}

}
