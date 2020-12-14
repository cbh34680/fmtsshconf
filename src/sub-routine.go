package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// KeyValueMap ... [[Key=Value], ...]
type KeyValueMap map[string]string

// HostConfMap ... [Host] = [[Key=Value], ...]
type HostConfMap map[string]KeyValueMap

func mustReadSSHConf(sshConf string) (hostConf HostConfMap, blockOrder []string, maxKeyLen int) {
	data, err := ioutil.ReadFile(sshConf)
	if err != nil {
		panic(fmt.Sprintf("ReadFile: error=%v\n", err))
	}

	hostConf = HostConfMap{}
	blockOrder = []string{}
	maxKeyLen = 0

	curBlock := ""
	rLine := regexp.MustCompile(`^\s*([^\s=]+)\s*=?\s*(.*)$`)
	rHostVal := regexp.MustCompile(`\s+`)

	for _, line := range strings.Split(string(data), "\n") {

		matches := rLine.FindAllStringSubmatch(line, -1)

		if len(matches) != 1 {
			continue
		}
		if len(matches[0]) != 3 {
			continue
		}

		curKey := matches[0][1]

		if curKey[0:1] == "#" {
			continue
		}

		curVal := strings.TrimSpace(matches[0][2])
		noConf := false

		if curKey == "Host" || curKey == "Match" {

			if curKey == "Host" {

				arr := rHostVal.Split(curVal, -1)
				sort.Strings(arr)
				curVal = strings.Join(arr, " ")
			}

			curBlock = curKey + " " + curVal
			noConf = true
		}

		if _, ok := hostConf[curBlock]; !ok {

			blockOrder = append(blockOrder, curBlock)
			hostConf[curBlock] = KeyValueMap{}
		}

		if noConf {
			continue
		}

		conf := hostConf[curBlock]
		conf[curKey] = curVal

		curKeyLen := len(curKey)
		if curKeyLen > maxKeyLen {
			maxKeyLen = curKeyLen
		}

		fmt.Fprintln(os.Stderr, conf)
	}

	return
}

// MustWriteSSHConfIfNeed ... func
func MustWriteSSHConfIfNeed(delHostBlock, sshConf string) {

	fmt.Fprintf(os.Stderr, "load file=[%s]\n", sshConf)
	hostConf, blockOrder, maxKeyLen := mustReadSSHConf(sshConf)
	fmt.Fprintf(os.Stderr, "done.\n")

	if len(blockOrder) > 0 {

		fmt.Fprintf(os.Stderr, "save file=[%s]\n", sshConf)
		file, err := os.Create(sshConf)
		if err != nil {
			panic(fmt.Sprintf("File Create %v", err))
		}

		defer file.Close()

		for _, block := range blockOrder {

			if delHostBlock != "" {

				if block == delHostBlock {
					fmt.Fprintf(os.Stderr, "skip host=[%s]", strings.Split(delHostBlock, " ")[1])
					continue
				}
			}

			indent := ""

			if block != "" {
				fmt.Fprintf(os.Stderr, "block [%s]\n", block)
				fmt.Fprintln(file, block)

				indent = "  "
			}

			for key, value := range hostConf[block] {

				pFormat := "%s%-" + strconv.Itoa(maxKeyLen) + "s %s\n"
				fmt.Fprintf(file, pFormat, indent, key, value)
				//fmt.Fprintf(file, "%s%s %s\n", indent, key, value)
			}

			fmt.Fprintln(file)
		}
		fmt.Fprintf(os.Stderr, "done.\n")
	}
}
