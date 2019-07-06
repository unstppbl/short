// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package dep

import (
	"tinyURL/app"
	"tinyURL/fw"
	"tinyURL/modern"
)

// Injectors from wire.go:

func InitGraphQlService(name string) fw.Service {
	logger := fw.NewLocalLogger()
	graphQlApi := app.NewGraphQlApi()
	server := modern.NewGraphGophers(logger, graphQlApi)
	service := fw.NewService(name, server, logger)
	return service
}
