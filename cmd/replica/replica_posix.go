// +build !windows

package main

func run(bckFilters, distFilters []string) {
	runReplica(
		bckFilters,
		distFilters,
	)
}
