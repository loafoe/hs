/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// iamLogoutCmd represents the logout command
var iamLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		iamClient, err := getIAMClient(cmd)
		if err != nil {
			fmt.Printf("error initalizing IAM client: %v\n", err)
			return
		}

		// Check the error return values of RevokeRefreshAccessToken and RevokeAccessToken
		err = iamClient.RevokeRefreshAccessToken()
		if err != nil {
			fmt.Printf("Error revoking refresh access token: %v\n", err)
		}

		err = iamClient.RevokeAccessToken()
		if err != nil {
			fmt.Printf("Error revoking access token: %v\n", err)
		}

		err = iamClient.EndSession()
		if err == io.EOF {
			fmt.Printf("Session ended\n")
		} else {
			fmt.Printf("Eror: %v\n", err)
		}
	},
}

func init() {
	iamCmd.AddCommand(iamLogoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// iamLogoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// iamLogoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
