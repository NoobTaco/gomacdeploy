package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func configureDefaultSettings(settings []string) {
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