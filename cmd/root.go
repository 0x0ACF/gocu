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
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/0x0ACF/gocu/internal/cache"
	"github.com/0x0ACF/gocu/pkg/http"

	"github.com/spf13/cobra"
)

var (
	method,
	data string
	headers []string
)

var rootCmd = &cobra.Command{
	Use:   "gocu",
	Short: "curl reimagined",
	Long:  `Gocu is a curl copycat, a CLI http client focused on simplicity and ease of use.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req := http.RequestInfo{
			Url:     replaceVars(args[0]),
			Method:  strings.ToUpper(method),
			Data:    replaceVars(data),
			Headers: extractHeaders(headers),
		}

		printRequestInfo(&req)

		resp, err := http.SendRequest(&req)

		if err != nil {
			log.Fatal(err)
		}

		printResponse(resp)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&method, "request", "X", "GET", "HTTP method to use")
	rootCmd.Flags().StringVarP(&data, "data", "d", "", "Data to send in the body of the request")
	rootCmd.Flags().StringArrayVarP(&headers, "header", "H", []string{}, "Header to add to the request")
}

func replaceVars(strToReplace string) string {
	if strToReplace != "" {
		ps := extractVarsPlaceholders(strToReplace)

		for _, p := range ps {
			varName := strings.Replace(p, "{{", "", 1)
			varName = strings.Replace(varName, "}}", "", 1)

			varVal, err := cache.GetVariable(varName)

			if err != nil {
				log.Fatal(err)
			}

			strToReplace = strings.Replace(strToReplace, p, varVal, -1)
		}
	}

	return strToReplace
}

func extractVarsPlaceholders(s string) []string {
	re, _ := regexp.Compile("{{.+?}}")

	return re.FindAllString(s, -1)
}

func extractHeaders(headersFlag []string) map[string]string {
	hs := map[string]string{
		"Content-Type": "application/json",
	}

	for _, h := range headersFlag {
		aux := strings.Split(h, ":")

		name := aux[0]
		val := strings.TrimSpace(strings.Join(aux[1:], ":"))

		hs[name] = replaceVars(val)
	}

	return hs
}

func printRequestInfo(r *http.RequestInfo) {
	fmt.Printf("%s %s\n", r.Method, r.Url)

	for name, val := range r.Headers {
		fmt.Printf("%s: %s\n", name, val)
	}

	j := http.PrettifyJson([]byte(r.Data))
	fmt.Printf("%s\n", j)
}

func printResponse(r *http.Response) {
	fmt.Println()
	fmt.Println(r.Status)
	fmt.Println(r.Data)
}
