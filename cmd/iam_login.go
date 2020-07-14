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
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/philips-software/go-hsdp-api/iam"

	"github.com/labstack/echo/v4"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

// dockerLoginCmd represents the login command
var iamLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into HSDP IAM using browser authentication flow",
	Long:  `Log into HSDP IAM using browser authentication flow`,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		environment, _ := cmd.Flags().GetString("environment")
		debug, _ := cmd.Flags().GetBool("debug")

		if clientID == "" || clientSecret == "" {
			fmt.Printf("this feature only works with official binaries.\n")
			return
		}
		// IAM
		iamClient, err := iam.NewClient(http.DefaultClient, &iam.Config{
			Region:         region,
			Environment:    environment,
			OAuth2ClientID: clientID,
			OAuth2Secret:   clientSecret,
			Debug:          debug,
		})
		if err != nil {
			fmt.Printf("error initializing IAM client: %v\n", err)
			os.Exit(1)
		}
		e := echo.New()
		e.HideBanner = true
		redirectURI := "http://localhost:35444/callback"
		loginSuccess := false
		e.GET("/callback", func(c echo.Context) error {
			code := c.QueryParam("code")
			err := iamClient.CodeLogin(code, redirectURI)
			if err != nil {
				c.HTML(http.StatusForbidden, "<html><body>Login failed</body></html>")
				go func() {
					time.Sleep(1 * time.Second)
					_ = e.Shutdown(context.Background())
				}()
				return err
			}
			c.HTML(http.StatusOK, "<html><body>You are now logged in! Feel free to close this window...</body></html>")
			loginSuccess = true
			go func() {
				time.Sleep(2 * time.Second)
				_ = e.Shutdown(context.Background())
			}()
			return nil
		})
		fmt.Printf("login using your browser ...\n")
		err = browser.OpenURL("https://iam-client-test.us-east.philips-healthsuite.com/authorize/oauth2/authorize?response_type=code&client_id=hsappclient&redirect_uri=http://localhost:35444/callback")
		if err != nil {
			fmt.Printf("failed to open browser login: %v\n", err)
			os.Exit(1)
		}
		_ = e.Start(":35444")
		if !loginSuccess {
			fmt.Printf("login failed. Please try again ...\n")
			os.Exit(1)
		}

		introspect, _, err := iamClient.Introspect()
		if err != nil {
			fmt.Printf("error performing introspect: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("logged in as: %s\n", introspect.Username)
		currentWorkspace.IAMAccessToken = iamClient.Token()
		currentWorkspace.IAMRefreshToken = iamClient.RefreshToken()
		currentWorkspace.save()
	},
}

func init() {
	iamCmd.AddCommand(iamLoginCmd)
}
