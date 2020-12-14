package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	argSSHConf := flag.String("config", "", "config path")
	argConfirm := flag.Bool("confirm", true, "confirm overwrite")
	argDelHost := flag.String("delhost", "", "delete host")

	flag.Parse()
	fmt.Fprintf(os.Stderr, "arg config=[%s] confirm=[%v] delhost=[%s]\n", *argSSHConf, *argConfirm, *argDelHost)

	delHostBlock := ""

	if *argDelHost != "" {
		delHostBlock = strings.Join([]string{"Host", *argDelHost}, " ")
	}

	usr, _ := user.Current()

	sshConf := *argSSHConf
	if sshConf == "" {
		usrHomeDir := usr.HomeDir

		//
		// Windows では通常 HOME は設定されていないので、これが set されていた
		// 場合は HOME を優先する
		//
		if envHomeDir := os.Getenv("HOME"); envHomeDir != "" {
			usrHomeDir = envHomeDir
		}
		sshConf = filepath.Join(usrHomeDir, ".ssh", "config")
	}

	if !PathExists(sshConf) {
		fmt.Fprintf(os.Stderr, "%s: file not found.\n", sshConf)
		os.Exit(0)
	}

	if *argConfirm {
		fmt.Printf("update %s ok ? [yes/No] ", sshConf)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if strings.ToLower(scanner.Text()) != "yes" {
			fmt.Fprintf(os.Stderr, "user cancelled.\n")
			os.Exit(0)
		}
	}

	confDir, confFile := filepath.Split(sshConf)
	sshConfBack := filepath.Join(confDir, confFile+".back-fmtsshconf")

	if !PathExists(sshConfBack) {

		if err := CopyFile(sshConf, sshConfBack); err != nil {

			panic(fmt.Sprintf("CopyFile: error=%v\n", err))
		}

		fmt.Fprintf(os.Stderr, "%s: backup-file created.\n", sshConf)
	}

	MustWriteSSHConfIfNeed(delHostBlock, sshConf)

	fmt.Fprintln(os.Stderr, "all done.")
}
