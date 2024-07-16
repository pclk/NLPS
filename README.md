# NLPS

Natural Language to Powershell Script

## Installation

Run `bin\mysetup.exe` on Windows.
This will:

- Install NLPS under Programs (x86)
- Add NLPS to PATH, so that you can use cmdlet `nlps` to run its functions
- Use C:\Users\{username}\AppData\Roaming\nlps to store configuration files
- Install this README file, accessible via the Start Menu

## Usage

`nlps [command] [flags]`

Commands:

- chat: Interactively chat, generate and run PowerShell commands.
- run: Generate and run a PowerShell command. The flags include modifiers for outputs.
  - -s, --silent: Run the command without asking for confirmation.
  - -n, --no-output: Do not display the command output and errors.
  - -e, --error-only: Only display the command error.
- history: Display and manage previously executed commands.
  - -c, --clear: Clear the history.
- config: Configure NLPS.

Global Flags:
-h, --help: help for nlps
-v, --verbose: Enable verbose output

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

```
NLPS/
├── cmd/
│ ├── root.go
│ ├── chat.go
│ ├── run.go
│ ├── history.go
│ ├── alias.go
│ ├── config.go
├── internal/
│ ├── ai/
│ │ └── openai.go
│ ├── powershell/
│ │ └── executor.go
│ ├── config/
│ │ └── config.go
│ ├── history/
│ │ └── history.go
│ └── ui/
│ └── styles.go
├── utils.go
├── main.go
├── README.md (You are here!)
```
