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
	"encoding/json"
	"fmt"

	"github.com/cheynewallace/tabby"

	"github.com/spf13/cobra"
)

// hasSessionsListCmd represents the list command
var hasSessionsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List HAS sessions",
	Long:    `Lists HAS sessions.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getHASClient(cmd, args)
		if err != nil {
			fmt.Printf("error initializing HAS client: %v\n", err)
			return
		}
		sessions, _, err := client.Sessions.GetSessions()
		if err != nil {
			fmt.Printf("error retrieving image list: %v\n", err)
			return
		}
		if sessions == nil || len(sessions.Sessions) == 0 {
			if jsonOut {
				fmt.Printf("[]\n")
				return
			}
			fmt.Printf("no sessions found\n")
			return
		}
		if jsonOut {
			data, _ := json.Marshal(sessions)
			fmt.Printf("%s\n", data)
			return
		}
		t := tabby.New()
		t.AddHeader("session id", "user", "region", "state", "url")
		for _, i := range sessions.Sessions {
			t.AddLine(i.SessionID,
				i.UserID,
				i.Region,
				i.State,
				i.SessionURL)
		}
		t.Print()

		fmt.Printf("\n")
	},
}

func init() {
	hasSessionsCmd.AddCommand(hasSessionsListCmd)
}
