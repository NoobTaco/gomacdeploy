package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// ============================================================
//  Project: Homebrew and App Installer
//  Author: Mike Norton
//  Date: 2024-11-06
//  Description: This Go application automates the installation
//               of Homebrew, various formulae, casks, and Mac
//               App Store applications. It also includes
//               optional package installations and system
//               cleanup tasks.
// ============================================================

// TODO: Add more error handling
// TODO: Add more comments
// TODO: Break up into separate modules
// TODO: Find a better name
// TODO: Add Settings module
// TODO: Add dock module
// TODO: Setup Git Login information

type Config struct {
	Casks    []string `yaml:"casks"`
	Formulae []string `yaml:"formulae"`
	AppStore []string `yaml:"appStore"`
	DefaultSettings []string `yaml:"defaultSettings"`
}

// const (
//     RED   = "\033[0;31m"
//     GREEN = "\033[0;32m"
//     NC    = "\033[0m" // No Color
// )

func main() {
	config, err := readConfig("config.yaml")
	if err != nil {
			fmt.Printf("Error reading config: %v\n", err)
			os.Exit(1)
	}

    // clearScreen()
    // printASCIIArt()
    // promptForRootPassword()
    // keepSudoAlive()
    // updateMacOS()
		// installRosetta()
		// installHomebrew()
		// setupHomebrew()
		// checkAndUpdateHomebrew()
		// installFormulae(config.Formulae)
		// installCasks(config.Casks)
		// installAppStoreApps(config.AppStore)
		// installDotNet()
		// cleanup()
		configureDefaultSettings(config.DefaultSettings)

}

func readConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
			return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
			return nil, err
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
			return nil, err
	}

	return &config, nil
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error clearing screen: %v\n", err)
	}
}

func printASCIIArt() {
    fmt.Println(" _           _        _ _       _     ")
    fmt.Println("(_)         | |      | | |     | |    ")
    fmt.Println(" _ _ __  ___| |_ __ _| | |  ___| |__  ")
    fmt.Println("| | |_ \\/ __| __/ _  | | | / __| |_ \\ ")
    fmt.Println("| | | | \\__ \\ || (_| | | |_\\__ \\ | | |")
    fmt.Println("|_|_| |_|___/\\__\\__,_|_|_(_)___/_| |_|")
    fmt.Println()
    fmt.Println()
    fmt.Println("Enter root password")
}

func promptForRootPassword() {
	cmd := exec.Command("sudo", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error prompting for root password: %v\n", err)
		os.Exit(1)
	}
}

// TODO Fix sudo keep alive
func keepSudoAlive() {
	go func() {
		for {
			cmd := exec.Command("sudo", "-n", "true")
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error keeping sudo alive: %v\n", err)
			}
			time.Sleep(60 * time.Second)
		}
	}()
}

func updateMacOS() {
    fmt.Println("Updating macOS...")
    cmd := exec.Command("sudo", "softwareupdate", "-i", "-a")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Printf("Error updating macOS: %v\n", err)
    }
}

// TODO Check for better command line options
func installRosetta() {
	fmt.Println("Checking if Rosetta is installed...")
	cmd := exec.Command("arch", "-x86_64", "/usr/bin/true")
	err := cmd.Run()
	if err == nil {
			fmt.Println("Rosetta is already installed.")
			return
	}

	fmt.Println("Installing Rosetta...")
	cmd = exec.Command("sudo", "softwareupdate", "--install-rosetta", "--agree-to-license")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
			fmt.Printf("Error installing Rosetta: %v\n", err)
	}
}

func installHomebrew() {
	fmt.Println("Checking if Homebrew is installed...")
	cmd := exec.Command("brew", "--version")
	err := cmd.Run()
	if err == nil {
			fmt.Println("Homebrew is already installed.")
			return
	}

	fmt.Println("Installing Homebrew...")
	cmd = exec.Command("bash", "-c", "NONINTERACTIVE=1 /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
			fmt.Printf("Error installing Homebrew: %v\n", err)
	}
}

func setupHomebrew() {
	zprofilePath := os.Getenv("HOME") + "/.zprofile"
	homebrewInit := `eval "$(/opt/homebrew/bin/brew shellenv)"`

	// Check if the initialization line is already in .zprofile
	file, err := os.Open(zprofilePath)
	if err != nil {
			if os.IsNotExist(err) {
					// Create the file if it does not exist
					file, err = os.Create(zprofilePath)
					if err != nil {
							fmt.Printf("Error creating .zprofile: %v\n", err)
							return
					}
			} else {
					fmt.Printf("Error opening .zprofile: %v\n", err)
					return
			}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
			if strings.Contains(scanner.Text(), homebrewInit) {
					fmt.Println("Homebrew initialization is already in .zprofile.")
					return
			}
	}

	// Append the initialization line to .zprofile
	file, err = os.OpenFile(zprofilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
			fmt.Printf("Error opening .zprofile for writing: %v\n", err)
			return
	}
	defer file.Close()

	_, err = file.WriteString(homebrewInit + "\n")
	if err != nil {
			fmt.Printf("Error writing to .zprofile: %v\n", err)
			return
	}

	// Immediately evaluate the Homebrew environment settings for the current session
	cmd := exec.Command("bash", "-c", homebrewInit)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
			fmt.Printf("Error evaluating Homebrew environment settings: %v\n", err)
	}
}

func checkAndUpdateHomebrew() {
	fmt.Println("Checking Homebrew installation and updating...")

	cmd := exec.Command("brew", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
			fmt.Printf("Error updating Homebrew: %v\n", err)
			return
	}

	cmd = exec.Command("brew", "doctor")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
			fmt.Printf("Error running brew doctor: %v\n", err)
			return
	}

	// Set the HOMEBREW_NO_INSTALL_CLEANUP environment variable
	os.Setenv("HOMEBREW_NO_INSTALL_CLEANUP", "1")
}

// Install Formulae
func installFormulae(formulae []string) {
	fmt.Println("Installing formulae...")
	for _, formula := range formulae {
			cmd := exec.Command("brew", "install", formula)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
					fmt.Printf("Failed to install %s. Continuing...\n", formula)
			}
	}
}

// Install Casks
func installCasks(casks []string) {
	fmt.Println("Installing casks...")
	for _, cask := range casks {
			cmd := exec.Command("brew", "install", "--cask", cask)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
					fmt.Printf("Failed to install %s. Continuing...\n", cask)
			}
	}
}

// Install App Store Apps
func installAppStoreApps(apps []string) {
	fmt.Println("Checking if mas is installed...")
	cmd := exec.Command("mas", "--version")
	err := cmd.Run()
	if err != nil {
			fmt.Println("mas is not installed. Installing mas...")
			cmd = exec.Command("brew", "install", "mas")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
					fmt.Printf("Error installing mas: %v\n", err)
					return
			}
	} else {
			fmt.Println("mas is already installed.")
	}

	fmt.Println("Installing Mac App Store applications...")
	for _, app := range apps {
			cmd := exec.Command("mas", "install", app)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
					fmt.Printf("Failed to install app %s. Continuing...\n", app)
			}
	}
}

// Install .NET
func installDotNet() {
	fmt.Println("Checking if .NET is installed...")
	cmd := exec.Command("dotnet", "--version")
	err := cmd.Run()
	if err == nil {
			fmt.Println(".NET is already installed.")
			return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Install .NET? [y/N]: ")
	reply, _ := reader.ReadString('\n')
	reply = strings.TrimSpace(reply)
	if strings.ToLower(reply) == "y" {
			cmd := exec.Command("brew", "install", "dotnet")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
					fmt.Printf("Failed to install .NET: %v\n", err)
					return
			}

			// Export DOTNET_ROOT to zsh
			zprofilePath := os.Getenv("HOME") + "/.zprofile"
			dotnetExport := `export DOTNET_ROOT="/opt/homebrew/opt/dotnet/libexec"`

			file, err := os.OpenFile(zprofilePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
					fmt.Printf("Error opening .zprofile for writing: %v\n", err)
					return
			}
			defer file.Close()

			_, err = file.WriteString(dotnetExport + "\n")
			if err != nil {
					fmt.Printf("Error writing to .zprofile: %v\n", err)
					return
			}

			// Immediately evaluate the DOTNET_ROOT environment setting for the current session
			cmd = exec.Command("bash", "-c", dotnetExport)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
					fmt.Printf("Error evaluating DOTNET_ROOT environment setting: %v\n", err)
			}
	}
}

// Cleanup
func cleanup() {
	fmt.Println("Cleaning up...")
	cmd := exec.Command("brew", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
			fmt.Printf("Error updating Homebrew: %v\n", err)
			return
	}

	cmd = exec.Command("brew", "upgrade")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
			fmt.Printf("Error upgrading Homebrew: %v\n", err)
			return
	}

	cmd = exec.Command("brew", "cleanup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
			fmt.Printf("Error cleaning up Homebrew: %v\n", err)
			return
	}

	cmd = exec.Command("brew", "doctor")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
			fmt.Printf("Error running brew doctor: %v\n", err)
			return
	}

	
}