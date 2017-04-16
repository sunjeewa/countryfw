/*
 * Project: countryfw
 * File: countryfw.go
 *
 * Copyright (c) 2016 Sanjeewa Wijesundara
 *
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	url    = "http://www.ipdeny.com/ipblocks/data/countries/"
	update = true
	rule   = "iptables -m tcp -p tcp -s %s -j ACCEPT\n"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getZoneFile(zone string) {
	fmt.Println("Downloading file")

	output, err := os.Create(zone)
	check(err)
	defer output.Close()

	response, err := http.Get(url + zone)
	check(err)
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	check(err)
	fmt.Println(n, " bytes downloaded")
}

func genRules(zone string, outputFile string) {
	fmt.Println("Generating iptables rules from ", zone)

	file, err := os.Open(zone)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Printf(rule, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func main() {
	// flag
	var update bool
	flag.BoolVar(&update, "update", false, "update source")
	flag.Parse()

	if update {
		getZoneFile("au.zone")
	}
	genRules("au.zone", "firewall")
}
