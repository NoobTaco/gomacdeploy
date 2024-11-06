// Package main provides a Go application that automates the installation
// of Homebrew, various formulae, casks, and Mac App Store applications.
// It also includes optional package installations and system cleanup tasks.
//
// Project: Go Mac Deploy (gomacdeploy)
// Author: Mike Norton
// Date: 2024-11-06
// Description: This Go application automates the installation
//
//	of Homebrew, various formulae, casks, and Mac
//	App Store applications. It also includes
//	optional package installations and system
//	cleanup tasks.
//
// The application reads a configuration file (config.yaml) to determine
// which packages and settings to install and configure. It performs the
// following tasks:
// - Clears the terminal screen
// - Prints ASCII art
// - Prompts for the root password
// - Keeps sudo alive
// - Updates macOS
// - Installs Rosetta (if needed)
// - Installs Homebrew (if not already installed)
// - Sets up Homebrew environment
// - Checks and updates Homebrew
// - Installs specified formulae
// - Installs specified casks
// - Installs specified Mac App Store applications
// - Installs .NET (if desired)
// - Configures default system settings
// - Configures Dock settings
// - Sets up Git login
// - Cleans up Homebrew installations
// - Reboots the system
//
// Version: v0.1.0
//
// This project was inspired by https://github.com/donnybrilliant/install.sh
//
// This project uses the MIT license. See the LICENSE file for details.
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

// TODO: Add more error handling
// TODO: Add more comments

type Config struct {
	Casks           []string `yaml:"casks"`
	Formulae        []string `yaml:"formulae"`
	AppStore        []string `yaml:"appStore"`
	DefaultSettings []string `yaml:"defaultSettings"`
	DockReplace     []string `yaml:"dockReplace"`
	DockAdd         []string `yaml:"dockAdd"`
	DockRemove      []string `yaml:"dockRemove"`
}

func main() {
	config, err := readConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	clearScreen()
	printASCIIArt()
	promptForRootPassword()
	keepSudoAlive()
	updateMacOS()
	installRosetta()
	installHomebrew()
	setupHomebrew()
	checkAndUpdateHomebrew()
	installFormulae(config.Formulae)
	installCasks(config.Casks)
	installAppStoreApps(config.AppStore)
	installDotNet()
	configureDefaultSettings(config.DefaultSettings)
	configureDockSettings(config.DockReplace, config.DockAdd, config.DockRemove)
	setupGitLogin()
	cleanup()
	finishAndReboot()

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
	clearScreen()
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
	clearScreen()
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
	clearScreen()
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
	clearScreen()
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
	clearScreen()
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
	clearScreen()
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
	clearScreen()
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
	clearScreen()
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

// Cleanup
func cleanup() {
	clearScreen()
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

func setupGitLogin() {
	clearScreen()
	reader := bufio.NewReader(os.Stdin)

	// Check if Git username is already set
	cmd := exec.Command("git", "config", "--global", "user.name")
	existingName, err := cmd.Output()
	if err == nil && len(existingName) > 0 {
		fmt.Printf("Existing Git username: %s\n", strings.TrimSpace(string(existingName)))
		fmt.Print("Do you want to overwrite it? [y/N]: ")
		reply, _ := reader.ReadString('\n')
		reply = strings.TrimSpace(reply)
		if strings.ToLower(reply) != "y" {
			fmt.Println("Keeping existing Git username.")
			return
		}
	}

	// Check if Git email is already set
	cmd = exec.Command("git", "config", "--global", "user.email")
	existingEmail, err := cmd.Output()
	if err == nil && len(existingEmail) > 0 {
		fmt.Printf("Existing Git email: %s\n", strings.TrimSpace(string(existingEmail)))
		fmt.Print("Do you want to overwrite it? [y/N]: ")
		reply, _ := reader.ReadString('\n')
		reply = strings.TrimSpace(reply)
		if strings.ToLower(reply) != "y" {
			fmt.Println("Keeping existing Git email.")
			return
		}
	}

	fmt.Println("SET UP GIT")
	fmt.Print("Please enter your git username: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Please enter your git email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	cmd = exec.Command("git", "config", "--global", "user.name", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to set git username: %v\n", err)
		return
	}

	cmd = exec.Command("git", "config", "--global", "user.email", email)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to set git email: %v\n", err)
		return
	}

	cmd = exec.Command("git", "config", "--global", "color.ui", "true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to set git color.ui: %v\n", err)
		return
	}

	fmt.Println("Git is Setup")
}

func finishAndReboot() {
	clearScreen()
	fmt.Println("______ _____ _   _  _____ ")
	fmt.Println("|  _  \\  _  | \\ | ||  ___|")
	fmt.Println("| | | | | | |  \\| || |__  ")
	fmt.Println("| | | | | | | .   ||  __| ")
	fmt.Println("| |/ /\\ \\_/ / |\\  || |___ ")
	fmt.Println("|___/  \\___/\\_| \\_/\\____/ ")

	fmt.Println()
	fmt.Println()
	fmt.Print("Would you like to reboot now? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	reply, _ := reader.ReadString('\n')
	reply = strings.TrimSpace(reply)
	if strings.ToLower(reply) == "y" {
		cmd := exec.Command("sudo", "reboot")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error rebooting: %v\n", err)
		}
		os.Exit(0)
	} else {
		fmt.Println("Reboot canceled.")
	}
}
