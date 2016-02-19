package main

import (
	"fmt"
	"github.com/bradfitz/slice"
	"github.com/codegangsta/cli"
	"log"
)

// Reverse available migrations based on the provided context.
func Reverse(ctx *cli.Context) {
	err := reverse(service, ctx.Bool("all"))
	if err != nil {
		log.Println(err)
	}
}

func ReverseTo(ctx *cli.Context) {
//	err := reverseTo(service, migrationName)
}

func reverse(service Servicer, all bool) error {
	initialized, err := service.Initialized()
	if err != nil {
		return fmt.Errorf("unable to check for database state: %s", err)
	}

	if !initialized {
		return fmt.Errorf("uninitialized database")
	}

	available, err := service.Available()
	if err != nil {
		return fmt.Errorf("unable to retrieve available migrations: %s", err)
	}

	applied, err := service.Applied()
	if err != nil {
		return fmt.Errorf("unable to retrieve applied migrations: %s", err)
	}

	if len(applied) == 0 {
		return fmt.Errorf("There are no applied migrations to reverse...")
	}

	//Check if available and applied migration consistency is unbroken
	err = checkAvailableAppliedSync(available, applied)
	if err != nil {
		return err
	}

	//Sort the applied migrations from last to first
	slice.Sort(applied, func(i, j int) bool {
		return applied[i].Name > applied[j].Name
	})

	for _, migration := range applied {
		err := service.Reverse(migration)
		if err != nil {
			return err
		}

		if !all {
			return nil
		}
	}

	return nil
}

func reverseTo(service Servicer, migrationName string)  error {
	return nil
}
