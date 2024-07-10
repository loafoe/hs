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
	"encoding/base64"
	"encoding/json"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

// Key is a struct that matches the JSON response structure
type Key struct {
	Version     string `json:"v"`
	PrivateKey  string `json:"pk"`
	ID          string `json:"id"`
	Region      string `json:"r"`
	Environment string `json:"e"`
}

// iamKeygenCmd represents the token command
var iamKeygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Packages service credentials as a key file",
	Long:  `Use this command to package service credentials into a key file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure log output
		logLevel := &slog.LevelVar{}
		logLevel.Set(slog.LevelInfo)
		if debug {
			logLevel.Set(slog.LevelDebug)
		}
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})))

		region, _ := cmd.Flags().GetString("region")
		environment, _ := cmd.Flags().GetString("environment")

		if region == "" {
			region = os.Getenv("HSP_IAM_REGION")
		}
		if environment == "" {
			environment = os.Getenv("HSP_IAM_ENVIRONMENT")
		}
		if region == "" || environment == "" {
			slog.Error("region and environment must be set")
			return
		}
		tokenFile, _ := cmd.Flags().GetString("token-file")

		serviceID, _ := cmd.Flags().GetString("service-id")
		if serviceID == "" { // Try reading from file
			serviceIdFile, _ := cmd.Flags().GetString("service-id-file")
			content, err := os.ReadFile(serviceIdFile)
			if err == nil {
				serviceID = string(content)
			}
		}
		privateKeyFile, _ := cmd.Flags().GetString("private-key-file")

		if privateKeyFile == "" {
			slog.Error("private-key-file is required")
			return
		}
		if serviceID == "" {
			slog.Error("service-id is required")
			return
		}

		key, err := os.ReadFile(privateKeyFile)
		if err != nil {
			slog.Error("error reading private key", "error", err)
			return
		}
		// Serialize it
		keyData := Key{
			Version:     "1",
			PrivateKey:  string(key),
			ID:          serviceID,
			Region:      region,
			Environment: environment,
		}
		token, err := json.Marshal(keyData)
		if err != nil {
			slog.Error("error marshalling key", "error", err)
			return
		}

		base64Token := base64.StdEncoding.EncodeToString(token)

		// Write the key file
		err = os.WriteFile(tokenFile, []byte(base64Token), 0644)
		if err != nil {
			slog.Error("error writing token file", "error", err)
		} else {
			slog.Info("token written", "file", tokenFile)
		}
	},
}

func init() {
	iamCmd.AddCommand(iamKeygenCmd)

	iamKeygenCmd.Flags().String("service-id", "", "The service ID to use")
	iamKeygenCmd.Flags().String("service-id-file", "", "A file containing the service id")
	iamKeygenCmd.Flags().String("private-key-file", "", "A file containing the private key")
	iamKeygenCmd.Flags().String("token-file", "token.txt", "The file to write the token to")
}
