package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
)

// Apply available migrations based on the provided context.
func Apply(ctx *cli.Context) {
	err := apply(service, ctx.Bool("all"))
	if err != nil {
		log.Println(err)
	}
}

func apply(service Servicer, all bool) error {
	initialized, err := service.Initialized()
	if err != nil {
		return fmt.Errorf("unable to check for database state: %s", err)
	}

	if !initialized {
		err := service.Initialize()
		if err != nil {
			return fmt.Errorf("unable to initialize database: %s", err)
		}
	}

	available, err := service.Available()
	if err != nil {
		return fmt.Errorf("unable to retrieve available migrations: %s", err)
	}

	applied, err := service.Applied()
	if err != nil {
		return fmt.Errorf("unable to retrieve applied migrations: %s", err)
	}

	//Check if available and applied migration consistency is unbroken
	err = checkAvailableAppliedSync(available, applied)
	if err != nil {
		return err
	}

	appliedCount := len(applied)
	availableCount := len(available)

	//There are no inconsistencies between applied and available migrations
	//However every migration is already applied
	if availableCount == appliedCount {
		return fmt.Errorf("All migrations are already applied")
	}

	//Start from the first nonaplied migration
	for _, migration := range available[appliedCount:] {
		err := service.Apply(migration)
		if err != nil {
			return err
		}

		if !all {
			return nil
		}
	}

	return nil
}
