/*
Copyright © 2020 Andy Lo-A-Foe <andy.lo-a-foe@philips.com>

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
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/philips-software/go-hsdp-api/console"

	"github.com/spf13/cobra"
)

// uaaLoginCmd represents the login command
var uaaLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to UAA",
	Long:  `Login to the regional UAA endpoint.`,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		environment, _ := cmd.Flags().GetString("environment")
		if region == "" {
			region = currentWorkspace.DefaultRegion
		}
		if environment == "" {
			environment = currentWorkspace.DefaultEnvironment
		}
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		username, password, err := credentials(username, password)
		fmt.Printf("\n")
		if err != nil {
			fmt.Printf("error logging in: %v\n", err)
			os.Exit(1)
		}
		consoleClient, err := console.NewClient(http.DefaultClient, &console.Config{
			Region:   region,
			DebugLog: "/tmp/console.log",
		})
		if err != nil {
			fmt.Printf("error initializing CONSOLE client: %v\n", err)
			os.Exit(1)
		}
		err = consoleClient.Login(username, password)
		if err != nil {
			fmt.Printf("error logging in: %v\n", err)
			os.Exit(1)
		}
		persistUAACredentials(consoleClient)
		fmt.Printf("%s\n", consoleClient.Token())
	},
}

func persistUAACredentials(consoleClient *console.Client) {
	currentWorkspace.UAAToken = consoleClient.Token()
	currentWorkspace.UAARefreshToken = consoleClient.RefreshToken()
	currentWorkspace.UAAIDToken = consoleClient.IDToken()
	currentWorkspace.UAAAccessTokenExpires = consoleClient.Expires()
	_ = currentWorkspace.save()
}

func credentials(username string, password string) (string, string, error) {
	var err error
	if username != "" && password != "" {
		return username, password, nil
	}
	reader := bufio.NewReader(os.Stdin)

	if username == "" {
		fmt.Print("Enter Username: ")
		username, err = reader.ReadString('\n')
		if err != nil {
			return "", "", err
		}
	}
	if password == "" {
		fmt.Print("Enter Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", "", err
		}
		password = string(bytePassword)
	}
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

func init() {
	uaaCmd.AddCommand(uaaLoginCmd)
	uaaLoginCmd.Flags().StringP("region", "t", "eu-west", "The region to login")
	uaaLoginCmd.Flags().StringP("username", "u", "", "UAA username")
	uaaLoginCmd.Flags().StringP("password", "p", "", "UAA password")
}
