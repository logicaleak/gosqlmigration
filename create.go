package main

import (
	"github.com/codegangsta/cli"
	"log"
	"path/filepath"
	"os"
	"strconv"
	"fmt"
	"regexp"
)

func Create(ctx *cli.Context) {
	err := create(service)
	if err!= nil {
		log.Println(err)
	}
}

func getLatestAmongMigrationSlice(migrationSlice []*Migration) string {
	var latestName string = migrationSlice[0].Name
	for _, migrationPtr := range migrationSlice {
		if migrationPtr.Name > latestName {
			latestName = migrationPtr.Name
		}
	}
	return latestName
}

func setUpNewSqlFile(f *os.File) {
	f.WriteString("-- rambler up")
	f.WriteString("\n")
	f.WriteString("\n")
	f.WriteString("\n")
	f.WriteString("-- rambler down")
}

func create(servicer Servicer) error {
	directoryPath := servicer.GetDirectoryPath()
	migrations, err := servicer.Available()
	if err != nil {
		return err
	}

	var newMigrationFileName string
	if len(migrations) == 0 {
		newMigrationFileName = filepath.Join(directoryPath, "migration0.sql")
	} else {
		var latestName string = getLatestAmongMigrationSlice(migrations)
		r := regexp.MustCompile("[0-9]+")
		numericPartOfName := r.FindString(latestName)
		fmt.Println("lastName : " + latestName)
		fmt.Println(numericPartOfName)
		integerNumericPart, err := strconv.Atoi(numericPartOfName)
		if err != nil {
			panic(err)
		}

		nextSuffixOfMigration := integerNumericPart + 1

		newMigrationFileName = filepath.Join(directoryPath, "migration" + strconv.Itoa(nextSuffixOfMigration) + ".sql")
	}

	fmt.Println("filename : " + newMigrationFileName)
	f, err := os.Create(newMigrationFileName)
	setUpNewSqlFile(f)
	if err != nil {
		panic(err)
	}
	return nil
}