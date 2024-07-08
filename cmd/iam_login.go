package cmd

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

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/philips-software/go-hsdp-api/iam"

	"github.com/labstack/echo/v4"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

type tokenOutput struct {
	AccessToken string `json:"access_token"`
}

// dockerLoginCmd represents the login command
var iamLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into HSDP IAM using browser authentication flow",
	Long:  `Log into HSDP IAM using browser authentication flow`,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		environment, _ := cmd.Flags().GetString("environment")
		if region == "" {
			region = currentWorkspace.DefaultRegion
		}
		if environment == "" {
			environment = currentWorkspace.DefaultEnvironment
		}
		serviceID, _ := cmd.Flags().GetString("service-id")
		if serviceID == "" { // Try reading from file
			serviceIdFile, _ := cmd.Flags().GetString("service-id-file")
			content, err := os.ReadFile(serviceIdFile)
			if err == nil {
				serviceID = string(content)
			}
		}

		if (clientID == "" || clientSecret == "") && serviceID == "" {
			fmt.Printf("this feature only works with official binaries.\n")
			return
		}
		// IAM
		iamClient, err := iam.NewClient(http.DefaultClient, &iam.Config{
			Region:         region,
			Environment:    environment,
			OAuth2ClientID: clientID,
			OAuth2Secret:   clientSecret,
		})
		if err != nil {
			fmt.Printf("error initializing IAM client: %v\n", err)
			os.Exit(1)
		}
		if serviceID != "" { // Service Identity flow
			privateKeyFile, _ := cmd.Flags().GetString("private-key-file")
			key, err := os.ReadFile(privateKeyFile)
			if err != nil {
				fmt.Printf("error reading private key: %v\n", err)
				os.Exit(1)
			}
			err = iamClient.ServiceLogin(iam.Service{
				ServiceID:  serviceID,
				PrivateKey: string(key),
			})
			if err != nil {
				fmt.Printf("error logging in: %v\n", err)
				os.Exit(1)
			}
			if clientID != "" {
				introspect, _, err := iamClient.Introspect()
				if err != nil {
					fmt.Printf("error performing introspect: %v\n", err)
					return
				}
				currentWorkspace.IAMUserUUID = introspect.Sub
				currentWorkspace.IAMAccessTokenExpires = introspect.Expires
			}
			token, _ := iamClient.Token()
			currentWorkspace.IAMAccessToken = token
			currentWorkspace.IAMIDToken = iamClient.IDToken()
			currentWorkspace.IAMRegion = region
			currentWorkspace.IAMEnvironment = environment
			if err := currentWorkspace.save(); err != nil {
				fmt.Printf("failed to save workspace: %v\n", err)
			}
			if jsonOut {
				data, _ := json.Marshal(tokenOutput{token})
				fmt.Printf("%s\n", string(data))
			}
			return
		}
		e := echo.New()
		e.HideBanner = true
		redirectURI := "http://localhost:35444/callback"
		loginSuccess := false
		e.GET("/callback", func(c echo.Context) error {
			code := c.QueryParam("code")
			err := iamClient.CodeLogin(code, redirectURI)
			if err != nil {
				_ = c.HTML(http.StatusForbidden, "<html><body>Login failed</body></html>")
				go func() {
					time.Sleep(1 * time.Second)
					_ = e.Shutdown(context.Background())
				}()
				return err
			}
			_ = c.HTML(http.StatusOK, "<html><body>You are now logged in! Feel free to close this window...</body></html>")
			loginSuccess = true
			go func() {
				time.Sleep(2 * time.Second)
				_ = e.Shutdown(context.Background())
			}()
			return nil
		})
		baseIAMURL := iamClient.BaseIAMURL().String()
		fmt.Printf("login using your browser to login...\n")
		err = browser.OpenURL(baseIAMURL + "/authorize/oauth2/authorize?response_type=code&client_id=hsappclient&redirect_uri=http://localhost:35444/callback")
		if err != nil {
			fmt.Printf("failed to open browser login: %v\n", err)
			return
		}
		done := make(chan bool)
		go func(done chan bool) {
			select {
			case <-done:
				return
			case <-time.After(5 * time.Minute):
				fmt.Printf("timed out waiting for login. Exiting ...\n")
				_ = e.Shutdown(context.Background())
			}
		}(done)
		_ = e.Start(":35444")
		if !loginSuccess {
			fmt.Printf("login failed. Please try again ...\n")
			return
		}

		introspect, _, err := iamClient.Introspect()
		if err != nil {
			fmt.Printf("error performing introspect: %v\n", err)
			return
		}
		fmt.Printf("logged in as: %s\n", introspect.Username)
		token, _ := iamClient.Token()
		currentWorkspace.IAMAccessToken = token
		currentWorkspace.IAMRefreshToken = iamClient.RefreshToken()
		currentWorkspace.IAMIDToken = iamClient.IDToken()
		currentWorkspace.IAMUserUUID = introspect.Sub
		currentWorkspace.IAMRegion = region
		currentWorkspace.IAMEnvironment = environment
		currentWorkspace.IAMAccessTokenExpires = introspect.Expires
		if jsonOut {
			data, _ := json.Marshal(tokenOutput{token})
			fmt.Printf("%s\n", string(data))
		}
		if err := currentWorkspace.save(); err != nil {
			fmt.Printf("failed to save workspace: %v\n", err)
		}
	},
}

func init() {
	iamCmd.AddCommand(iamLoginCmd)
	iamLoginCmd.Flags().String("service-id", "", "The service ID to use")
	iamLoginCmd.Flags().String("service-id-file", "", "A file containing the service id")
	iamLoginCmd.Flags().String("private-key-file", "", "A file containing the private key")
}
