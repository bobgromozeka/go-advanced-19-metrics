package main

import (
	"fmt"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func init() {
	naIfEmpty(&buildVersion)
	naIfEmpty(&buildDate)
	naIfEmpty(&buildCommit)
}

func printMetadata() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func naIfEmpty(v *string) {
	if *v == "" {
		*v = "N/A"
	}
}
