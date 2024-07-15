# NLPS

Natural Language to Powershell Script

## Installation

Run `bin\mysetup.exe` on Windows.
This will:

- Install NLPS under Programs (x86)
- Add NLPS to PATH, so that you can use cmdlet `nlps` to run its functions
- Use C:\Users\{username}\AppData\Roaming\nlps to store configuration files
- Install this README file, accessible via the Start Menu

## Accessing this README

After installation, you can access this README file via:

- The Start Menu shortcut in the NLPS group
- The installation directory (typically C:\Program Files (x86)\NLPS)

## Uninstallation

Go to Programs (x86)/NLPS, and click on the uninstaller.
This will undo all actions above, including:

- Removing NLPS from Programs (x86)
- Removing NLPS from PATH
- Deleting configuration files from C:\Users\{username}\AppData\Roaming\nlps
- Removing this README file

## Build

`go build main.go`

Install [Inno Setup](https://jrsoftware.org/isdl.php#stable)
Add C:\Program Files (x86)\Inno Setup 6 to PATH

`iscc nlps.iss`

## Testing

You can set your OPENAI_API_KEY in .env.local and run `go test github.com/pclk/NLPS/test`.

For manual testing, you can `go run main.go`. Note that it will create a config file at your AppData.
