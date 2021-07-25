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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/philips-labs/terraform-backend-hsdp/backend/types"
	"github.com/spf13/cobra"
)

// uaaLoginCmd represents the login command
var tfstateForceUnlockCmd = &cobra.Command{
	Use:   "force-unlock",
	Short: "Force unlock a locked state",
	Long:  `Force unlocks a locked state.`,
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		lockID, _ := cmd.Flags().GetString("id")
		key, err := tfstateKey(key)
		if err != nil {
			fmt.Printf("error reading key: %v\n", err)
			os.Exit(1)
		}
		lockID, err = tfstateLockID(lockID)
		if err != nil {
			fmt.Printf("error reading lock ID: %v\n", err)
			os.Exit(1)
		}
		tfstateEndpoint := currentWorkspace.TFStateInstanceURL + "/" + key
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
		lock := types.Lock{
			ID: lockID,
		}
		data, _ := json.Marshal(&lock)
		body := bytes.NewReader(data)

		req, _ := http.NewRequest("UNLOCK", tfstateEndpoint, body)
		req.Header.Set("Authorization", "Basic "+currentWorkspace.TFStateCreds)
		req.Header.Set("Content-Type", "application/json")
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Printf("error fetch list: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		response, _ := ioutil.ReadAll(resp.Body)
		if jsonOut {
			fmt.Printf("%s\n", string(response))
			return
		}
		fmt.Printf("STATUS %d\n", resp.StatusCode)
	},
}

func tfstateLockID(id string) (string, error) {
	var err error
	if id != "" {
		return id, nil
	}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Lock ID: ")
	id, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(id), nil
}

func init() {
	tfstateCmd.AddCommand(tfstateForceUnlockCmd)
	tfstateForceUnlockCmd.Flags().StringP("key", "k", "", "The tfstate key")
	tfstateForceUnlockCmd.Flags().StringP("id", "l", "", "The tfstate lock ID")
}
