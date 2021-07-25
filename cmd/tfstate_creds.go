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
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// uaaLoginCmd represents the login command
var tfstateCredsCmd = &cobra.Command{
	Use:   "creds",
	Short: "Set credentials for tfstate",
	Long:  `Sets the credentials to use for the tfstate instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		address, _ := cmd.Flags().GetString("address")
		address, err := tfstateAddress(address)
		if err != nil {
			fmt.Printf("error reading address: %v\n", err)
			os.Exit(1)
		}
		username, password, err = credentials(username, password)
		fmt.Printf("\n")
		if err != nil {
			fmt.Printf("error reading credentials: %v\n", err)
			os.Exit(1)
		}
		persistTFState(address, base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
	},
}

func persistTFState(address, credentials string) {
	currentWorkspace.TFStateInstanceURL = address
	currentWorkspace.TFStateCreds = credentials
	_ = currentWorkspace.save()
}

func tfstateAddress(address string) (string, error) {
	var err error
	if address != "" {
		return address, nil
	}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Address: ")
	address, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(address), nil
}

func init() {
	tfstateCmd.AddCommand(tfstateCredsCmd)
	tfstateCredsCmd.Flags().StringP("address", "a", "", "The tfstate address")
	tfstateCredsCmd.Flags().StringP("username", "u", "", "tfstate username")
	tfstateCredsCmd.Flags().StringP("password", "p", "", "tfstate password")
}
