/*
gocu is a curl copycat, a CLI HTTP client focused on simplicity and ease of use
Copyright (C) 2025  Andrés C.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"github.com/0x0ACF/gocu/internal/cache"

	"github.com/spf13/cobra"
)

func init() {
	varsCmd.AddCommand(varsLsCmd, varsGetCmd, varsAddCmd, varsModCmd, varsRmCmd, varsClearCmd)
	rootCmd.AddCommand(varsCmd)
}

var varsCmd = &cobra.Command{
	Use:   "vars",
	Args:  cobra.MinimumNArgs(1),
	Short: "Manage the variables used as placeholders",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var varsLsCmd = &cobra.Command{
	Use:   "ls",
	Args:  cobra.NoArgs,
	Short: "Lists all saved variables",
	Run: func(cmd *cobra.Command, args []string) {
		listSavedVariables()
	},
}

var varsGetCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.ExactArgs(1),
	Short: "Gets a variable value",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		val, err := cache.GetVariable(name)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(val)
		}
	},
}

var varsAddCmd = &cobra.Command{
	Use:   "add",
	Args:  cobra.ExactArgs(2),
	Short: "Adds a new variable",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		val := args[1]

		err := cache.AddVariable(name, val)

		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

var varsModCmd = &cobra.Command{
	Use:   "mod",
	Args:  cobra.ExactArgs(2),
	Short: "Modifies a variable value",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		val := args[1]

		err := cache.ModifyVariable(name, val)

		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

var varsRmCmd = &cobra.Command{
	Use:   "rm",
	Args:  cobra.ExactArgs(1),
	Short: "Removes a variable",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		err := cache.RemoveVariable(name)

		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

var varsClearCmd = &cobra.Command{
	Use:   "clear",
	Args:  cobra.NoArgs,
	Short: "Removes all saved variables",
	Run: func(cmd *cobra.Command, args []string) {
		err := cache.RemoveAllVariables()

		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func listSavedVariables() {
	vars := cache.Variables()

	if len(vars) == 0 {
		fmt.Println("No variables saved")
	} else {
		for name, val := range vars {
			fmt.Printf("%s=%s\n", name, val)
		}
	}
}
