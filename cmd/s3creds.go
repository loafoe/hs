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
	"fmt"

	"github.com/philips-software/go-hsdp-api/s3creds"
	"github.com/spf13/cobra"
)

// s3credsCmd represents the s3creds command
var s3credsCmd = &cobra.Command{
	Use:     "s3creds",
	Aliases: []string{"s3c", "s3cr"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(s3credsCmd)
}

func getCredentialsClient(cmd *cobra.Command, _ []string) (*s3creds.Client, error) {
	productKey := currentWorkspace.S3CredsProductKey
	region := currentWorkspace.DefaultRegion
	environment := currentWorkspace.DefaultEnvironment
	iamClient, err := getIAMClient(cmd)
	if err != nil {
		return nil, fmt.Errorf("iam client: %w", err)
	}
	if productKey == "" {
		return nil, fmt.Errorf("no S3 Credentials productKey configured")
	}
	return s3creds.NewClient(iamClient, &s3creds.Config{
		Region:      region,
		Environment: environment,
	})
}
