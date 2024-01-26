package cmd

/*
Copyright © 2022 Andy Lo-A-Foe

*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type versionRequest struct {
	Version string `json:"version"`
}

// pullCmd represents the pull command
var tfstatePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls a specific state version",
	Long:  `Pulls a specific state version from the backend`,
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		key, err := tfstateKey(key)
		if err != nil {
			fmt.Printf("error reading key: %v\n", err)
			os.Exit(1)
		}
		if len(args) < 1 {
			fmt.Printf("specify version to retrieve as the only argument\n")
			os.Exit(1)
		}
		version := args[0]

		tfstateEndpoint := currentWorkspace.TFStateInstanceURL + "/versions?ref=" + key
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
		jsonRequest := versionRequest{
			Version: version,
		}
		body, _ := json.Marshal(jsonRequest)

		req, _ := http.NewRequest("RETRIEVE", tfstateEndpoint, bytes.NewReader(body))
		req.Header.Set("Authorization", "Basic "+currentWorkspace.TFStateCreds)
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Printf("error fetch list: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		var state map[string]interface{}

		switch resp.StatusCode {
		case http.StatusNoContent:
			fmt.Printf("{}\n")
			os.Exit(1)
		case http.StatusOK:
			break
		default:
			fmt.Printf("Unexpected status: %d\n", resp.StatusCode)
			os.Exit(1)
		}
		if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
			fmt.Printf("error decoding versions body: %v\n", err)
			os.Exit(1)
		}

		data, _ := json.Marshal(state)
		fmt.Printf("%s\n", pretty(data))
		return
	},
}

func init() {
	tfstateCmd.AddCommand(tfstatePullCmd)
	tfstatePullCmd.Flags().StringP("key", "k", "", "The tfstate key")

}
