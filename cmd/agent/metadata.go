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
	fmt.Printf("Build version: %s", buildVersion)
	fmt.Printf("Build date: %s", buildDate)
	fmt.Printf("Build commit: %s", buildCommit)
}

func naIfEmpty(v *string) {
	if *v == "" {
		*v = "N/A"
	}
}
