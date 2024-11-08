// main_test.go
package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestReadConfig(t *testing.T) {
	content := `
casks:
	- google-chrome
formulae:
	- git
appStore:
	- 1234567890
defaultSettings:
	- "defaults write com.apple.finder AppleShowAllFiles YES"
dockReplace:
	- "/Applications/Safari.app|/Applications/Firefox.app"
dockAdd:
	- "/Applications/Slack.app"
dockRemove:
	- "/Applications/Mail.app"
`
	tmpfile, err := ioutil.TempFile("", "example.*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	config, err := readConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(config.Casks) != 1 || config.Casks[0] != "google-chrome" {
		t.Errorf("Expected google-chrome, got %v", config.Casks)
	}
}

func TestClearScreen(t *testing.T) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestPrintASCIIArt(t *testing.T) {
	printASCIIArt()
}

func TestPromptForRootPassword(t *testing.T) {
	cmd := exec.Command("sudo", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestKeepSudoAlive(t *testing.T) {
	go keepSudoAlive()
	time.Sleep(2 * time.Second)
}

func TestUpdateMacOS(t *testing.T) {
	cmd := exec.Command("sudo", "softwareupdate", "-i", "-a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestInstallRosetta(t *testing.T) {
	cmd := exec.Command("arch", "-x86_64", "/usr/bin/true")
	err := cmd.Run()
	if err == nil {
		t.Log("Rosetta is already installed.")
		return
	}

	cmd = exec.Command("sudo", "softwareupdate", "--install-rosetta", "--agree-to-license")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestInstallHomebrew(t *testing.T) {
	cmd := exec.Command("brew", "--version")
	err := cmd.Run()
	if err == nil {
		t.Log("Homebrew is already installed.")
		return
	}

	cmd = exec.Command("bash", "-c", "NONINTERACTIVE=1 /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSetupHomebrew(t *testing.T) {
	zprofilePath := os.Getenv("HOME") + "/.zprofile"
	homebrewInit := `eval "$(/opt/homebrew/bin/brew shellenv)"`

	file, err := os.Open(zprofilePath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(zprofilePath)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		} else {
			t.Fatalf("Expected no error, got %v", err)
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), homebrewInit) {
			t.Log("Homebrew initialization is already in .zprofile.")
			return
		}
	}

	file, err = os.OpenFile(zprofilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(homebrewInit + "\n")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	cmd := exec.Command("bash", "-c", homebrewInit)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCheckAndUpdateHomebrew(t *testing.T) {
	cmd := exec.Command("brew", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	cmd = exec.Command("brew", "doctor")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestInstallFormulae(t *testing.T) {
	formulae := []string{"git", "wget"}
	for _, formula := range formulae {
		cmd := exec.Command("brew", "install", formula)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}
}

func TestInstallCasks(t *testing.T) {
	casks := []string{"google-chrome", "visual-studio-code"}
	for _, cask := range casks {
		cmd := exec.Command("brew", "install", "--cask", cask)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}
}

func TestInstallAppStoreApps(t *testing.T) {
	cmd := exec.Command("mas", "--version")
	err := cmd.Run()
	if err != nil {
		cmd = exec.Command("brew", "install", "mas")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}

	apps := []string{"1234567890"}
	for _, app := range apps {
		cmd := exec.Command("mas", "install", app)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}
}

func TestInstallDotNet(t *testing.T) {
	cmd := exec.Command("dotnet", "--version")
	err := cmd.Run()
	if err == nil {
		t.Log(".NET is already installed.")
		return
	}

	cmd = exec.Command("brew", "install", "dotnet")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	zprofilePath := os.Getenv("HOME") + "/.zprofile"
	dotnetExport := `export DOTNET_ROOT="/opt/homebrew/opt/dotnet/libexec"`

	file, err := os.OpenFile(zprofilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dotnetExport + "\n")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	cmd = exec.Command("bash", "-c", dotnetExport)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestConfigureDefaultSettings(t *testing.T) {
	settings := []string{"defaults write com.apple.finder AppleShowAllFiles YES"}
	for _, setting := range settings {
		cmd := exec.Command("bash", "-c", setting)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}
}

func TestConfigureDockSettings(t *testing.T) {
	cmd := exec.Command("brew", "install", "dockutil")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	replaceItems := []string{"/Applications/Safari.app|/Applications/Firefox.app"}
	addItems := []string{"/Applications/Slack.app"}
	removeItems := []string{"/Applications/Mail.app"}

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
				t.Errorf("Expected no error, got %v", err)
			}
		}
	}

	for _, app := range addItems {
		cmd := exec.Command("dockutil", "--add", app)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	for _, app := range removeItems {
		cmd := exec.Command("dockutil", "--remove", app)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}
}

func TestCleanup(t *testing.T) {
	cmd := exec.Command("brew", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	cmd = exec.Command("brew", "upgrade")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	cmd = exec.Command("brew", "cleanup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	cmd = exec.Command("brew", "doctor")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSetupGitLogin(t *testing.T) {
	cmd := exec.Command("git", "config", "--global", "user.name", "testuser")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	cmd = exec.Command("git", "config", "--global", "user.email", "testuser@example.com")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	cmd = exec.Command("git", "config", "--global", "color.ui", "true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestFinishAndReboot(t *testing.T) {
	cmd := exec.Command("sudo", "reboot")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
