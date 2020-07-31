package launch

import (
	"os"
)

// Checklist check list before launch service, panic if anythins wrong
//
//	launch.Checklist()
//
func Checklist() {

	name := os.Getenv("NAME")
	if name == "" {
		panic("env NAME=\"sample-jp\" not found")
	}

	region := os.Getenv("REGION")
	if region == "" {
		panic("env REGION=\"us\" not found")
	}

	branch := os.Getenv("BRANCH")
	if branch == "" {
		panic("env BRANCH=\"stable\" not found")
	}
}
