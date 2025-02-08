// +build windows

package main

import (
	"fmt"
	"net"
	"github.com/marko2276/wslroutesvc/runner"
	"github.com/marko2276/wslroutesvc/network"
)


func fixRoutes(wslIfaceName string, runner runner.Runner) int {
	wslIface := network.NewIface(wslIfaceName, runner)

	if (wslIface.ID == "") || (wslIface.IP.String() == "<nil>" ) {
		fmt.Printf("Could not find interface ID for WSL interface %s\n", wslIfaceName)
		return -1
	}

	fmt.Printf("%s interface ID: %s, IP: %s \n", wslIfaceName, wslIface.ID, wslIface.IP)

	routeList := network.NewRouteList(runner)

	for _, r := range routeList.Routes {
		if r.Network.Contains(wslIface.IP) && r.InterfaceID != wslIface.ID {
			maskSize, _ := r.Network.Mask.Size()

			// Prevent broad routes from qualifying
			if maskSize < 16 {
				continue
			}

			// Remove the route
			out, err := r.Remove(runner)

			if err != nil {
				fmt.Printf("Failed to remove route %s with interface ID %s!\n%s\n%v\n", r.Network, r.InterfaceID, out, err)
				continue
			}
			fmt.Printf("Route %s with interface ID %s removed!", r.Network, r.InterfaceID)
		}
	}

	wslNet := net.IPNet{wslIface.IP.Mask(net.CIDRMask(20, 32)), net.CIDRMask(20, 32)}
	fmt.Println("WSL network is ", wslNet.String())
	/*Add route for WSL*/
	wslRoute := network.NewRoute(wslNet, wslIface.ID)
	out, err := wslRoute.Add(runner)

	if err != nil {
		fmt.Printf("Failed to add route %s with interface ID %s!\n%s\n%v\n", wslRoute.Network.String(), wslRoute.InterfaceID, out, err)
		return -2
	}
	fmt.Printf("Route %s with interface ID %s added!", wslRoute.Network.String(), wslRoute.InterfaceID)
	return 0

}
