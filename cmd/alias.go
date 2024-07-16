package cmd

import (
	"fmt"

	"github.com/pclk/NLPS/internal/config"
	"github.com/pclk/NLPS/internal/ui"
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage command aliases",
	Long:  `This command allows you to add, remove, and list command aliases.`,
	Run:   runAlias,
}

func init() {
	rootCmd.AddCommand(aliasCmd)
	aliasCmd.AddCommand(addAliasCmd)
	aliasCmd.AddCommand(removeAliasCmd)
	aliasCmd.AddCommand(listAliasCmd)
}

var addAliasCmd = &cobra.Command{
	Use:   "add <name> <command>",
	Short: "Add a new alias",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name, command := args[0], args[1]
		if !isValidAlias(name) {
			fmt.Println(ui.Error("Invalid alias name. Alias names cannot contain spaces or special characters."))
			return
		}
		err := config.AddAlias(name, command)
		if err != nil {
			fmt.Println(ui.Error(fmt.Sprintf("Failed to add alias: %v", err)))
			return
		}
		fmt.Printf("%s %s %s %s\n", ui.Success("Alias"), ui.Command(name), ui.Success("added for command:"), ui.Value(command))
	},
}

var removeAliasCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove an alias",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if config.RemoveAlias(name) {
			fmt.Printf("%s %s %s\n", ui.Success("Alias"), ui.Command(name), ui.Success("removed."))
		} else {
			fmt.Printf("%s %s %s\n", ui.Error("Alias"), ui.Command(name), ui.Error("not found."))
		}
	},
}

var listAliasCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	Run: func(cmd *cobra.Command, args []string) {
		aliases := config.GetAliases()
		if len(aliases) == 0 {
			fmt.Println(ui.Info("No aliases found."))
			return
		}
		fmt.Println(ui.Info("Aliases:"))
		for name, command := range aliases {
			fmt.Printf("%s: %s\n", ui.Command(name), ui.Value(command))
		}
	},
}

func runAlias(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func isValidAlias(alias string) bool {
	for _, char := range alias {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return len(alias) > 0
}
