// Copyright 2020 Paul Vollmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	GitCommit = "?"
	BuildDate = "?"
	BuildOS   = "?"
	BuildArch = "?"
)

func PrintInfo() {
	fmt.Printf("robotstxt-service\n")
	fmt.Printf("  Version    : %s\n", Version)
	fmt.Printf("  Git commit : %s\n", GitCommit)
	fmt.Printf("  Go Version : %s\n", runtime.Version())
	fmt.Printf("  Build      : %s\n", BuildDate)
	fmt.Printf("  OS/Arch    : %s/%s\n", BuildOS, BuildArch)
}
