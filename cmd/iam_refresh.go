package cmd

/*
Copyright © 2024 Andy Lo-A-Foe <andy.lo-a-foe@philips.com>

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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/philips-software/go-hsdp-api/iam"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// TokenResponse is a struct that matches the JSON response structure
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

// iamTokenCmd represents the token command
var iamRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Continuously refreshes a service identity token",
	Long:  `Refreshes access token, useful for sidecar processes.`,
	Run: func(cmd *cobra.Command, args []string) {
		region, _ := cmd.Flags().GetString("region")
		environment, _ := cmd.Flags().GetString("environment")
		every, _ := cmd.Flags().GetInt64("every")
		clientId, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")

		if every == 0 {
			slog.Error("every must be > 0", "every", every)
			return
		}
		if region == "" {
			region = os.Getenv("HSP_IAM_REGION")
		}
		if environment == "" {
			environment = os.Getenv("HSP_IAM_ENVIRONMENT")
		}

		iamClient, err := iam.NewClient(http.DefaultClient, &iam.Config{
			Region:         region,
			Environment:    environment,
			OAuth2ClientID: clientId,
			OAuth2Secret:   clientSecret,
		})
		if err != nil {
			slog.Error("error initializing IAM client", "error", err)
			return
		}

		tokenExchangeIssuer, _ := cmd.Flags().GetString("token-exchange-issuer")
		tokenFile, _ := cmd.Flags().GetString("token-file")
		// Loop here
		retries := 3

		for {
			serviceID, _ := cmd.Flags().GetString("service-id")
			if serviceID == "" { // Try reading from file
				serviceIdFile, _ := cmd.Flags().GetString("service-id-file")
				content, err := os.ReadFile(serviceIdFile)
				if err == nil {
					serviceID = string(content)
				}
			}
			connectorId, _ := cmd.Flags().GetString("connector-id")
			privateKeyFile, _ := cmd.Flags().GetString("private-key-file")

			if privateKeyFile == "" {
				slog.Error("private-key-file is required")
				return
			}
			if serviceID == "" {
				slog.Error("service-id is required")
				return
			}

			err = retry.Do(func() error {
				key, err := os.ReadFile(privateKeyFile)
				if err != nil {
					slog.Error("error reading private key", "error", err)
					return err
				}
				slog.Info("logging in", "serviceID", serviceID)
				err = iamClient.ServiceLogin(iam.Service{
					ServiceID:  serviceID,
					PrivateKey: string(key),
				})
				if err != nil {
					slog.Error("error logging in", "error", err)
					return err
				}

				token, _ := iamClient.Token()
				if tokenExchangeIssuer != "" {
					// Prepare the data to be sent in the request body
					slog.Info("exchanging token", "issuer", tokenExchangeIssuer)
					data := url.Values{}
					data.Set("connector_id", connectorId)
					data.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
					data.Set("scope", "openid groups federated:id")
					data.Set("requested_token_type", "urn:ietf:params:oauth:token-type:access_token")
					data.Set("subject_token", token)
					data.Set("subject_token_type", "urn:ietf:params:oauth:token-type:access_token")
					// Create a new request
					req, err := http.NewRequest("POST", tokenExchangeIssuer+"/token", bytes.NewBufferString(data.Encode()))
					if err != nil {
						return err
					}

					// Set the Authorization header
					auth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
					req.Header.Add("Authorization", "Basic "+auth)

					// Set the Content-Type header
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

					// Create an HTTP client and send the request
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						return err
					}
					defer resp.Body.Close()
					// Read and print the response body
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return err
					}
					// Parse the JSON response
					var tokenResponse TokenResponse
					if err := json.Unmarshal(body, &tokenResponse); err != nil {
						return err
					}
					token = tokenResponse.AccessToken
				}
				if tokenFile != "" {
					err = os.WriteFile(tokenFile, []byte(token), 0644)
					if err != nil {
						retries = retries - 1
						slog.Error("error writing token file", "error", err)
					} else {
						slog.Info("token written", "file", tokenFile)
					}
				}
				if jsonOut {
					data, _ := json.Marshal(tokenOutput{token})
					fmt.Printf("%s\n", string(data))
				}
				return nil
			}, retry.Attempts(uint(retries)), retry.Delay(5*time.Second))
			if err != nil {
				slog.Error("failed to get token", "error", err)
				return
			}
			// Wait for next cycle
			slog.Info("sleeping", "seconds", every)
			time.Sleep(time.Duration(every) * time.Second)
		}
	},
}

func init() {
	iamCmd.AddCommand(iamRefreshCmd)

	iamRefreshCmd.Flags().String("service-id", "", "The service ID to use")
	iamRefreshCmd.Flags().String("service-id-file", "", "A file containing the service id")
	iamRefreshCmd.Flags().String("private-key-file", "", "A file containing the private key")
	iamRefreshCmd.Flags().Int64("every", 900, "Refresh every n seconds")
	iamRefreshCmd.Flags().String("token-file", "token.txt", "The file to write the token to")
	iamRefreshCmd.Flags().String("token-exchange-issuer", "", "Exchanges the token with the specified issuer")
	iamRefreshCmd.Flags().String("connector-id", "hsdp", "The connector ID to use")
	iamRefreshCmd.Flags().String("client-id", "alloy", "The client ID to use")
	iamRefreshCmd.Flags().String("client-secret", "observability", "The client secret to use")
}
