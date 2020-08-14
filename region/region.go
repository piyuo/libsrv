package region

import "os"

// Current return current region
//
var Current = os.Getenv("REGION")
