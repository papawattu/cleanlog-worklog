package main

import (
	"context"
	"log"
	"net/http"
)

// Controller is an interface that defines the methods that a controller should implement
type Controller interface {
	Start() error

	Stop() error

	HandleRequest(w http.ResponseWriter, r *http.Request)
}

// Controllers is a struct that holds the controllers

type Controllers struct {
	controllers []Controller
}

// Start starts the controllers
func (c *Controllers) Start() error {
	for _, controller := range c.controllers {
		if err := controller.Start(); err != nil {
			log.Fatalf("Error starting controller: %v", err)
			return err
		}
	}
	return nil
}

// Stop stops the controllers
func (c *Controllers) Stop() error {
	for _, controller := range c.controllers {
		if err := controller.Stop(); err != nil {
			log.Fatalf("Error stopping controller: %v", err)
			return err
		}
	}
	return nil
}

// HandleRequest handles the request
func (c *Controllers) HandleRequest(w http.ResponseWriter, r *http.Request) {
	for _, controller := range c.controllers {
		log.Printf("Handling request: %v", r)
		controller.HandleRequest(w, r)
	}
}

// MakeControllers creates the controllers
func MakeControllers(ctx context.Context, controllers ...Controller) (*Controllers, error) {
	log.Printf("Creating controllers")
	return &Controllers{controllers: controllers}, nil
}
