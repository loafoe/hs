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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/philips-software/go-hsdp-api/iron"
	"github.com/spf13/cobra"
)

// ironCmd represents the iron command
var ironCmd = &cobra.Command{
	Use:   "iron",
	Short: "Interaction with HSPD IronIO",
	Long:  `This is a replacement of the iron CLI with a focus on dockerized tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(ironCmd)
	ironCmd.PersistentFlags().StringP("cluster", "c", "", "Cluster to use")
}

func readIronConfig(path ...string) (*iron.Config, error) {
	var configFile string
	if len(path) == 0 {
		home, _ := os.UserHomeDir()
		configFile = filepath.Join(home, ".iron.json")
	} else {
		configFile = path[0]
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var config iron.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	if config.ProjectID == "" {
		return nil, fmt.Errorf("invalid config: %v", config)
	}
	return &config, nil
}
