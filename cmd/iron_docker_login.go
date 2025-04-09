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
	"github.com/dip-software/go-dip-api/iron"

	"github.com/spf13/cobra"
)

// dockerLoginCmd represents the login command
var dockerLoginCmd = &cobra.Command{
	Use:   "login -u username -p password -e email -s server",
	Short: "register docker credentials with Iron",
	Long:  `register docker credentials with Iron`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := readIronConfig()
		if err != nil {
			fmt.Printf("error reading iron config: %v\n", err)
			return
		}
		if len(config.ClusterInfo) == 0 {
			fmt.Printf("missing required cluster info in config")
		}
		config.Debug = debug
		client, err := iron.NewClient(config)
		if err != nil {
			fmt.Printf("error configuring iron client: %v\n", err)
			return
		}

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		email, _ := cmd.Flags().GetString("email")
		server, _ := cmd.Flags().GetString("server")

		ok, _, err := client.Codes.DockerLogin(iron.DockerCredentials{
			Username:      username,
			Password:      password,
			Email:         email,
			ServerAddress: server,
		})
		if err != nil {
			fmt.Printf("error registering credentials: %v\n", err)
			fmt.Printf("\n")
			_ = cmd.Help()
			return
		}
		if !ok {
			fmt.Printf("credentials verification failed.\n")
			return
		}
		fmt.Printf("credentials stored successfully.\n")
	},
}

func init() {
	dockerCmd.AddCommand(dockerLoginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dockerLoginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	dockerLoginCmd.Flags().StringP("username", "u", "", "Docker registry username")
	dockerLoginCmd.Flags().StringP("password", "p", "", "Docker registry password")
	dockerLoginCmd.Flags().StringP("email", "e", "", "Docker registry email address")
	dockerLoginCmd.Flags().StringP("server", "s", "", "Docker registry server address")

}
