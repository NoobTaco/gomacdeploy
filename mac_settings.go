package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func configureDefaultSettings(settings []string) {
	clearScreen()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Configure default system settings? [Y/n]: ")
	reply, _ := reader.ReadString('\n')
	reply = strings.TrimSpace(reply)
	if reply == "" || strings.ToLower(reply) == "y" {
		fmt.Println("Configuring default settings...")
		for _, setting := range settings {
			cmd := exec.Command("bash", "-c", setting)
			// add a printout in terminal of the cmd prompt
			fmt.Printf("Applying setting: %s\n", setting)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Failed to apply setting: %s. Continuing...\n", setting)
			}
		}
	}
}

func configureDockSettings(replaceItems, addItems, removeItems []string) {
	clearScreen()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Apply Dock settings? [y/N]: ")
	reply, _ := reader.ReadString('\n')
	reply = strings.TrimSpace(reply)
	if strings.ToLower(reply) == "y" {
		fmt.Println("Installing dockutil...")
		cmd := exec.Command("brew", "install", "dockutil")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to install dockutil: %v\n", err)
			return
		}

		// Handle replacements
		for _, item := range replaceItems {
			parts := strings.Split(item, "|")
			if len(parts) == 2 {
				addApp := parts[0]
				replaceApp := parts[1]
				cmd := exec.Command("dockutil", "--add", addApp, "--replacing", replaceApp)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Printf("Failed to replace %s with %s: %v\n", replaceApp, addApp, err)
				}
			}
		}

		// Handle additions
		for _, app := range addItems {
			cmd := exec.Command("dockutil", "--add", app)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Failed to add %s: %v\n", app, err)
			}
		}

		// Handle removals
		for _, app := range removeItems {
			cmd := exec.Command("dockutil", "--remove", app)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Failed to remove %s: %v\n", app, err)
			}
		}
	}
}
