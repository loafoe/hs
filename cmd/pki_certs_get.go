/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/philips-software/go-hsdp-api/console"
	"github.com/philips-software/go-hsdp-api/pki"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// pkiCertsGetCmd represents the get command
var pkiCertsGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	Short:   "Get information about a certifcate",
	Long:    `Retrieves and lists certificate info based on serial`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("must specify PKI logical path and certicate serial\n")
			os.Exit(1)
		}
		logicalPath := args[0]
		serial := args[1]
		region, _ := cmd.Flags().GetString("region")
		environment, _ := cmd.Flags().GetString("environment")
		if region == "" {
			region = currentWorkspace.DefaultRegion
		}
		if environment == "" {
			environment = currentWorkspace.DefaultEnvironment
		}
		if currentWorkspace.UAAToken == "" {
			fmt.Printf("Login to UAA first using 'uaa login'\n")
			os.Exit(1)
		}
		consoleClient, err := console.NewClient(http.DefaultClient, &console.Config{
			Region: region,
		})
		if err != nil {
			fmt.Printf("error initializing CONSOLE client: %v\n", err)
			os.Exit(1)
		}
		consoleClient.SetTokens(currentWorkspace.UAAToken,
			currentWorkspace.UAARefreshToken, currentWorkspace.UAAIDToken, currentWorkspace.UAAAccessTokenExpires)
		iamClient, err := getIAMClient(cmd)
		if err != nil {
			fmt.Printf("error initializing IAM client: %v\n", err)
			os.Exit(1)
		}
		pkiClient, err := pki.NewClient(consoleClient, iamClient, &pki.Config{
			Region:      region,
			Environment: environment,
		})
		if err != nil {
			fmt.Printf("error initializing PKI client: %v\n", err)
			os.Exit(1)
		}
		cert, _, err := pkiClient.Services.GetCertificateBySerial(logicalPath, serial, nil)
		if err != nil {
			fmt.Printf("error getting certificate: %v\n", err)
			os.Exit(1)
		}
		data, err := json.Marshal(cert)
		if err != nil {
			fmt.Printf("error marshalling certificate result: %v\n", err)
			return
		}
		persistUAACredentials(consoleClient)
		fmt.Println(pretty(data))
	},
}

func init() {
	pkiCertsCmd.AddCommand(pkiCertsGetCmd)
}
