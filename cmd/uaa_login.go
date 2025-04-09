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
	"net/http"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/dip-software/go-dip-api/console"

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
			Region: region,
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
		token, _ := consoleClient.Token()
		fmt.Printf("%v\n", token)
	},
}

func persistUAACredentials(consoleClient *console.Client) {
	token, _ := consoleClient.Token()
	currentWorkspace.UAAToken = token.AccessToken
	currentWorkspace.UAARefreshToken = token.RefreshToken
	currentWorkspace.UAAIDToken = consoleClient.IDToken()
	currentWorkspace.UAAAccessTokenExpires = consoleClient.Expires()
	_ = currentWorkspace.save()
}

func credentials(username string, password string) (string, string, error) {
	var err error
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter Username [%s]: ", username)
	entered, err := reader.ReadString('\n')
	if err != nil {
		return username, password, err
	}
	trimmedUsername := strings.TrimSpace(entered)
	if trimmedUsername != "" {
		username = trimmedUsername
	}

	masked := ""
	if len(password) > 0 {
		masked = "****"
	}
	fmt.Printf("Enter Password [%s]: ", masked)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return username, password, err
	}
	trimmedPassword := strings.TrimSpace(string(bytePassword))
	if trimmedPassword != "" {
		password = trimmedPassword
	}
	return username, password, nil
}

func init() {
	uaaCmd.AddCommand(uaaLoginCmd)
	uaaLoginCmd.Flags().StringP("region", "t", "eu-west", "The region to login")
	uaaLoginCmd.Flags().StringP("username", "u", "", "UAA username")
	uaaLoginCmd.Flags().StringP("password", "p", "", "UAA password")
}
