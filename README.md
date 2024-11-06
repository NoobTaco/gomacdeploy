# Homebrew and App Installer

This Go application automates the installation of Homebrew, various formulae, casks, and Mac App Store applications. It also includes optional package installations and system cleanup tasks.

## Features

- Clears the terminal screen
- Prints ASCII art
- Prompts for the root password
- Keeps sudo alive
- Updates macOS
- Installs Rosetta (if needed)
- Installs Homebrew (if not already installed)
- Sets up Homebrew environment
- Checks and updates Homebrew
- Installs specified formulae
- Installs specified casks
- Installs specified Mac App Store applications
- Installs .NET (if desired)
- Configures default system settings
- Configures Dock settings
- Sets up Git login
- Cleans up Homebrew installations

## Configuration

The application reads a configuration file (`config.yaml`) to determine which packages and settings to install and configure. Here is an example configuration:

```yaml
casks:
  - google-chrome
  - visual-studio-code
formulae:
  - git
  - wget
appStore:
  - 409201541  # Pages
  - 409203825  # Numbers
defaultSettings:
  - defaults write -g AppleShowAllExtensions -bool true
dockReplace:
  - /Applications/Google Chrome.app|Safari
dockAdd:
  - [WezTerm.app](http://_vscodecontentref_/0)
dockRemove:
  - FaceTime
dotfilesRepo: 'https://github.com/NoobTaco/dotfiles'
```
