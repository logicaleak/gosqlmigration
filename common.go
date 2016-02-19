package main
import "fmt"

func checkAvailableAppliedSync(available []*Migration, applied []*Migration) error {
	//If number of applied is bigger than available, there must be missing migrations at the end of the sequence
	if missingCount := len(applied) - len(available); missingCount > 0 {
		return fmt.Errorf("%s number of missing migrations", missingCount)
	}

	//If available is more than applied, then we will have an available slice
	//until all the existing number of applied migrations
	//If they are the same size, this just wont have any effect
	var slicedAvailable []*Migration
	slicedAvailable = available[0:len(applied)]

	i := 0
	j := 0
	for i <= len(slicedAvailable) && j <= len(applied) {
		if slicedAvailable[i].Name > applied[j].Name {
			return fmt.Errorf("missing migration: %s", applied[j].Name)
		}

		if slicedAvailable[i].Name < applied[j].Name {
			return fmt.Errorf("out of order migration: %s", slicedAvailable[i].Name)
		}
		i++
		j++
	}

	return nil

}
