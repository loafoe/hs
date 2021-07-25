/*
Copyright © 2021 Andy Lo-A-Foe <andy.lo-a-foe@philips.com>

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
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// uaaLoginCmd represents the login command
var tfstateStatesCmd = &cobra.Command{
	Use:   "states",
	Short: "Lists all known state keys",
	Long:  `Lists all known state keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		tfstatesEndpoint := currentWorkspace.TFStateInstanceURL + "/states"
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
		req, _ := http.NewRequest(http.MethodGet, tfstatesEndpoint, nil)
		req.Header.Set("Authorization", "Basic "+currentWorkspace.TFStateCreds)
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Printf("error fetch list: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		states := make([]string, 0)

		if err := json.NewDecoder(resp.Body).Decode(&states); err != nil {
			fmt.Printf("error decoding states body: %v\n", err)
			os.Exit(1)
		}
		if jsonOut {
			data, _ := json.Marshal(states)
			fmt.Printf("%s\n", string(data))
			return
		}

		fmt.Printf("Found %d state keys\n", len(states))
		for _, v := range states {
			fmt.Printf("%s\n", v)
		}
	},
}

func init() {
	tfstateCmd.AddCommand(tfstateStatesCmd)
}
