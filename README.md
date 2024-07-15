# PowerShell Command Generator

## Overview

This project consists of two main files: `main.go` and `powershell.go`. Together, they create a program that generates and executes PowerShell commands using OpenAI's GPT model. The program allows users to input a description of what they want to do, generates a corresponding PowerShell command, and then executes it if the user confirms.

## main.go

### Purpose

`main.go` serves as the entry point of the application and handles the main logic flow, user interaction, and integration with the OpenAI API.

### Functions

#### `main()`

- **Purpose**: Entry point of the program.
- **Functionality**:
  1. Initializes the OpenAI client and PowerShell session.
  2. Runs the main interaction loop.
  3. Handles user input, command generation, confirmation, and execution.
- **Philosophy**: Keeps the main function clean and high-level, delegating specific tasks to other functions.

#### `initOpenAI() *openai.Client`

- **Purpose**: Initializes the OpenAI client.
- **Functionality**:
  1. Loads environment variables from a `.env.local` file.
  2. Creates and returns a new OpenAI client using the API key.
- **Philosophy**: Separates the initialization of external services, making it easier to modify or replace in the future.

#### `getUserInput(reader *bufio.Reader) string`

- **Purpose**: Prompts the user for input and reads their response.
- **Functionality**:
  1. Prints a prompt asking the user to describe the desired PowerShell command.
  2. Reads the user's input, trimming any whitespace.
- **Philosophy**: Encapsulates user input logic, making it easy to modify how input is collected if needed.

#### `generateCommand(client *openai.Client, userInput string) string`

- **Purpose**: Generates a PowerShell command based on user input using the OpenAI API.
- **Functionality**:
  1. Sends a request to the OpenAI API with the user's input.
  2. Handles potential errors or empty responses.
  3. Extracts and returns the generated command.
- **Philosophy**: Isolates the AI interaction, making it easier to modify or replace the AI service if needed.

#### `confirmExecution(reader *bufio.Reader) bool`

- **Purpose**: Asks the user for confirmation before executing the generated command.
- **Functionality**:
  1. Prompts the user to confirm execution.
  2. Interprets the user's response, defaulting to "yes" if the input is empty.
- **Philosophy**: Adds a safety check, ensuring the user approves of the command before execution.

#### `executeCommand(ps *PowerShell, command string)`

- **Purpose**: Executes the generated PowerShell command and displays the output.
- **Functionality**:
  1. Sends the command to the PowerShell session.
  2. Formats and prints the output.
- **Philosophy**: Separates command execution from other logic, allowing for easier modification of how commands are run or output is displayed.

## powershell.go

### Purpose

`powershell.go` encapsulates all PowerShell-related functionality, providing a clean interface for interacting with PowerShell from the main program.

### Structures

#### `PowerShell`

- **Purpose**: Represents a PowerShell session.
- **Fields**:
  - `cmd`: The underlying PowerShell process.
  - `stdin`: Writer for sending commands to PowerShell.
  - `stdout`: Reader for receiving output from PowerShell.
- **Philosophy**: Encapsulates all PowerShell-related state, providing a clean abstraction for PowerShell interactions.

### Functions

#### `randomString(n int) string`

- **Purpose**: Generates a random string of specified length.
- **Functionality**: Creates a string of random letters (uppercase and lowercase).
- **Philosophy**: Provides a utility for creating unique markers, enhancing the reliability of command output parsing.

#### `initPowerShell() *PowerShell`

- **Purpose**: Initializes a new PowerShell session.
- **Functionality**:
  1. Starts a new PowerShell process with specific flags.
  2. Sets up input and output pipes.
  3. Initializes the PowerShell environment (e.g., setting encoding).
- **Philosophy**: Centralizes PowerShell initialization, ensuring consistent setup across the application.

#### `(ps *PowerShell) sendCommand(command string) string`

- **Purpose**: Sends a command to PowerShell and retrieves the output.
- **Functionality**:
  1. Wraps the command in a try-catch block for error handling.
  2. Uses a unique marker to identify the end of command output.
  3. Sets the output encoding to UTF-8 to avoid character encoding issues.
  4. Captures and returns the command output.
- **Philosophy**: Provides a robust method for command execution and output capture, handling potential errors and encoding issues.

#### `(ps *PowerShell) close()`

- **Purpose**: Closes the PowerShell session.
- **Functionality**: Terminates the PowerShell process.
- **Philosophy**: Ensures proper cleanup of resources when the program exits.
