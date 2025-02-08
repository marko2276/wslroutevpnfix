// +build windows

// Checks for any conflicting routes inserted by VPN
// that are required for WSL and overwrites them
// in order to have networking support for WSL available
// also when connected to VPN
//
package main

import (
	"fmt"
	"os"
	"strings"
	"bufio"
	"github.com/marko2276/wslroutesvc/runner"
)

func main() {
	const appName = "wslroutevpnfix"

	runner := runner.ExecRunner{}

	if len(os.Args) > 1 {
		cmd := strings.ToLower(os.Args[1])
		switch cmd {
		case "debug":
			fmt.Printf("Does nothing.....\n")
			return
		}
	}

	//list of known names for wsl eth addapters
	wslEthApapterNames := []string{"vEthernet (WSL)", "vEthernet (WSL (Hyper-V firewall))"}

	for _, wslEthApapter := range wslEthApapterNames {
        fmt.Printf("Trying adapter %s ...\n", wslEthApapter)
		if fixRoutes(wslEthApapter, &runner) == 0 {
			break
		}
    }

	buf := bufio.NewReader(os.Stdin)
    fmt.Print("Press Enter to continue ..... ")
    sentence, err := buf.ReadBytes('\n')
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(string(sentence))
    }

	return
}
