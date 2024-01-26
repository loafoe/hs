/*
Copyright © 2021 Andy Lo-A-Foe <andy.lo-a-foe@philips.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
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

// pkiCertsListCmd represents the list command
var pkiCertsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li"},
	Short:   "Lists known certificates",
	Long:    `Lists the known certificates under a HSDP PKI policy CA`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("must specify PKI logical path\n")
			os.Exit(1)
		}
		logicalPath := args[0]
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
		certs, _, err := pkiClient.Services.GetCertificates(logicalPath, nil)
		if err != nil {
			fmt.Printf("error getting certificate list: %v\n", err)
			os.Exit(1)
		}
		data, err := json.Marshal(certs)
		if err != nil {
			fmt.Printf("error marshalling CerticateList result: %v\n", err)
			return
		}
		persistUAACredentials(consoleClient)
		fmt.Println(pretty(data))
	},
}

func init() {
	pkiCertsCmd.AddCommand(pkiCertsListCmd)
}
