package main

import "C"
import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

const DEBUG = false
const VERSION = "0.6.0"
const BASE_HELP = `[style.bold][style.bright.magenta]toml-alias[style.reset] - tool for easy TOML configurable alias creation.
Depending on current shell prompt localize configuration file in [style.italic][style.bright.yellow]$HOME/.config/toml-alias/config.toml[style.reset] .
Configuration file consists custom aliases activated by binary call (eg. by adding the binary into $PATH directory).
Binary name should be renamed to alias what is configured in TOML file.
For more info please visit [style.italic][style.underline][style.bright.blue]https://github.com/kkuba91/toml-alias[style.reset]

Default args:
    -h, --help      [style.italic]Prints out this help message[style.reset]
    -V, --version   [style.italic]Returns version of the tool[style.reset]

`

type AliasStageConfig struct {
	PrintStdout bool     `toml:"print-stdout"`
	PrintMatch  bool     `toml:"print-match"`
	AllowFail   bool     `toml:"allow-fail"`
	MatchStdout string   `toml:"match-stdout"`
	MatchMsg    string   `toml:"match-msg"`
	OnEnd       string   `toml:"print-on-end"`
	OnSuccess   string   `toml:"print-on-success"`
	OnFailure   string   `toml:"print-on-failure"`
	Cmd         []string `toml:"cmd"`
	PreCmd      []string `toml:"pre-cmd"`
	PostCmd     []string `toml:"post-cmd"`
}
type AliasConfig struct {
	HideBaseHelp    bool               `toml:"hide-base-help"`
	HideBaseVersion bool               `toml:"hide-base-version"`
	CustomHelp      string             `toml:"custom-help"`
	CustomVersion   string             `toml:"custom-version"`
	Stage           []AliasStageConfig `toml:"stage"`
}

type Config map[string]AliasConfig

func getAliasName() string {
	// getAliasName gets binary name - this binary, what is called.
	//
	// Returns:
	//   - A string containing binary name.

	execPath, _ := os.Executable()
	execName := strings.Split(filepath.Base(execPath), ".exe")[0]
	logPrint("binary: ", execName, "\n")
	return execName
}

func getHomeDir() string {
	// getHomeDir gets $HOME / %HOME% / windows User's path to search for configuration.
	//
	// Returns:
	//   - A string containing User's directory path.

	dirname := os.Getenv("HOME")
	if dirname == "" {
		dirnameWin, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		dirname = dirnameWin
	}
	logPrint("home directory: ", dirname, "\n")
	return dirname
}

func readAllAliasesConfigFromFile() Config {
	// readAllAliasesConfigFromFile reads and parses all aliases configuration from a TOML file.
	//
	// Returns:
	//   - A `Config` struct containing all aliases defined in the configuration file.
	//
	// Notes:
	//   - The function assumes the helper function `getHomeDir()` to check config file exists.
	//   - It uses the `github.com/BurntSushi/toml` packages for file handling and TOML parsing.
	//   - The function terminates execution on critical errors, such as missing or invalid configuration files.
	//
	// Example Configuration File Path:
	//   - `~/.config/toml-alias/config.toml` for Linux, Mac or cygwin (Windows)
	//   - `C:/Users/JohnnyBravo/.config/toml-alias/config.toml` for Windows and executing from CMD

	// Build the path to the configuration file
	configPath := filepath.Join(getHomeDir(), ".config", "toml-alias", "config.toml")
	logPrint("config path: ", configPath, "\n")

	// Check if the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err != nil {
			log.Fatalf("Failed to read configuration: %v", fmt.Errorf("configuration file not found: %s", configPath))
		}
		return nil
	}

	// Parse the TOML file
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		if err != nil {
			log.Fatalf("Failed to read configuration: %v", fmt.Errorf("failed to parse TOML file: %w", err))
		}
		return nil
	}
	logPrint("Configuration file read successfully!")
	return config
}

func parseAliasConfig(config Config) AliasConfig {
	// parseAliasConfig retrieves single AliasConfig data from the provided configuration of all of them.
	//
	// Parameters:
	//   - config: A map where the keys are alias names (strings) and the values are `Config` structs.
	//
	// Returns:
	//   - An `AliasConfig` struct corresponding to the current alias name, or an empty `AliasConfig` if no match is found.

	thisAliasName := getAliasName()
	for name, alias := range config {
		if name == thisAliasName {
			logPrint("Configuration file parsed successfully for alias:", thisAliasName)
			return alias
		}
	}
	return AliasConfig{}
}

func checkValidConfig(aliasConfig *AliasStageConfig) {
	// checkValidConfig validates the provided AliasConfig has necessary data to call command.
	//
	// Parameters:
	//   - aliasConfig: A pointer to the AliasConfig struct containing the configuration of single alias.
	//
	// Behavior:
	//   - If validation fails, the function terminates the program using `log.Fatal` and provides an appropriate error message.
	//
	// Note:
	//   - This function does not return a value and is designed to enforce configuration validity before proceeding with further operations.

	if aliasConfig == nil {
		log.Fatal(fmt.Errorf("no alias configuration provided"))
	}

	if len(aliasConfig.Cmd) == 0 {
		log.Fatal(fmt.Errorf("no commands to execute, %s", aliasConfig.OnFailure))
	}
}

func executeCommands(aliasConfig *AliasStageConfig) string {
	// executeCommands executes a list of commands defined in the AliasStageConfig structure (pre-cmd, main-cmd and post-cmd).
	//
	// Parameters:
	//   - aliasConfig: A pointer to the AliasStageConfig struct containing the configuration of single alias.
	//
	// Returns:
	//   - A string containing the standard output (stdout) from the main command execution.

	// [Pre-Command]
	if len(aliasConfig.PreCmd) > 0 {
		exec.Command(aliasConfig.PreCmd[0], aliasConfig.PreCmd[1:]...).Output()
	}
	// [Main-Command]
	// The first element in Cmd is the command, and the rest are arguments.
	stdout, _ := exec.Command(aliasConfig.Cmd[0], aliasConfig.Cmd[1:]...).Output()
	// [Post-Command]
	if len(aliasConfig.PostCmd) > 0 {
		exec.Command(aliasConfig.PostCmd[0], aliasConfig.PostCmd[1:]...).Output()
	}

	// Print stdout to the console
	if aliasConfig.PrintStdout {
		print(string(stdout), "\n")
	}
	return string(stdout)
}

func printPostData(aliasConfig *AliasStageConfig, matchResult string) {
	// printPostData prints post-processing data based on the provided alias configuration and match result.
	//
	// If the `OnEnd` field in `aliasConfig` is not empty, it prints the value of `OnEnd` followed by a space.
	// If `matchResult` is not empty, it prints the `matchResult` followed by a newline.
	//
	// Parameters:
	//   - aliasConfig: A pointer to the AliasStageConfig struct containing the configuration of single alias.
	//   - matchResult: A string representing the result of an optional matching operation, typically derived from stdout.

	// Print ending string (defined in config)
	if aliasConfig.OnEnd != "" {
		print(aliasConfig.OnEnd, " ")
	}

	// Print match result from optional matching of stdout
	if matchResult != "" {
		print(matchResult, "\n")
	}
}

func processStdout(aliasConfig *AliasStageConfig, stdout string) string {
	// processStdout processes the standard output (stdout) of a command based on the configuration provided in aliasConfig.
	//
	// Parameters:
	//   - aliasConfig: A pointer to the AliasStageConfig struct containing the configuration of single alias.
	//   - stdout: A string representing the standard output to be processed.
	//
	// Returns:
	//   - A string indicating either the success message (`OnSuccess`) or failure message (`OnFailure`) based on the match result.

	customMatchResult := ""
	if aliasConfig.MatchStdout != "" {
		matchRe, err := regexp.Compile(aliasConfig.MatchStdout)
		if err != nil {
			log.Fatalf("Error compiling regex: %s \n%s", err, aliasConfig.OnFailure)
		}
		matches := matchRe.FindStringSubmatch(stdout)
		if len(matches) > 0 {
			matchResult := ""
			groupNames := matchRe.SubexpNames()
			for i, name := range groupNames {
				if name == "match" && matchResult == "" {
					matchResult = matches[i]
				}
			}
			if matchResult == "" {
				matchResult = matches[0]
			}

			if aliasConfig.PrintMatch {
				matchMsg := "Result matched:"
				if aliasConfig.MatchMsg != "" {
					matchMsg = aliasConfig.MatchMsg
				}
				print(fmt.Sprintf("%s %s\n", matchMsg, matchResult))
			}
			customMatchResult = aliasConfig.OnSuccess
		} else if aliasConfig.AllowFail {
			customMatchResult = ""
		} else {
			customMatchResult = aliasConfig.OnFailure
		}
	}
	return customMatchResult
}

func executeAlias(aliasConfig *AliasStageConfig) {
	// executeAlias summary procedure to call whole proces and its parts one by one.
	//
	// Parameters:
	//   - aliasConfig: A pointer to the AliasStageConfig struct containing the configuration of single alias.
	//
	// Returns:
	//   - A string indicating either the success message (`OnSuccess`) or failure message (`OnFailure`) based on the match result.

	checkValidConfig(aliasConfig)
	stdout := executeCommands(aliasConfig)
	matchResult := processStdout(aliasConfig, stdout)
	printPostData(aliasConfig, matchResult)

}

func checkVersionAndPrint(thisAliasConfig AliasConfig) {
	// checkVersionAndPrint verify `-V` or `--version` be used and prints version data.
	//
	// Parameters:
	//   - thisAliasConfig: A pointer to the AliasConfig struct containing the configuration of single alias.

	logPrint("HideBaseVersion: ", thisAliasConfig.HideBaseVersion, "\n")
	logPrint("CustomVersion: ", thisAliasConfig.CustomVersion, "\n")
	firstArg := ""
	if len(os.Args) > 1 {
		firstArg = os.Args[1]
	}
	if firstArg == "-V" || firstArg == "--version" {
		if !thisAliasConfig.HideBaseVersion {
			print("Version=", VERSION, "\n")
		}
		if thisAliasConfig.CustomVersion != "" {
			print(thisAliasConfig.CustomVersion, "\n")
		}
		os.Exit(0)
	}
}

func checkHelpAndPrint(thisAliasConfig AliasConfig) {
	// checkHelpAndPrint verify `-h` or `--help` be used and prints help message.
	//
	// Parameters:
	//   - thisAliasConfig: A pointer to the AliasConfig struct containing the configuration of single alias.

	logPrint("HideBaseHelp: ", thisAliasConfig.HideBaseHelp, "\n")
	logPrint("CustomHelp: ", thisAliasConfig.CustomHelp, "\n")
	firstArg := ""
	if len(os.Args) > 1 {
		firstArg = os.Args[1]
	}
	if firstArg == "-h" || firstArg == "--help" {
		if !thisAliasConfig.HideBaseHelp {
			print(BASE_HELP)
		}
		if thisAliasConfig.CustomHelp != "" {
			print(thisAliasConfig.CustomHelp, "\n")
		}
		os.Exit(0)
	}
}

//export processAll
func processAll() {
	// processAll summary ALL procedures in one function.

	config := readAllAliasesConfigFromFile()
	thisAliasConfig := parseAliasConfig(config)

	checkVersionAndPrint(thisAliasConfig)
	checkHelpAndPrint(thisAliasConfig)

	for i, thisStageConfig := range thisAliasConfig.Stage {
		logPrint("Execute stage ", i+1, "\n")
		executeAlias(&thisStageConfig)
	}
	print("\n")
}

func main() {
}
