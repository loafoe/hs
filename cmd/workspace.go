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
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/philips-software/go-hsdp-api/has"
	"github.com/philips-software/go-hsdp-api/iron"
	"github.com/spf13/cobra"
)

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"ws"},
	Short:   "Manage workspaces",
	Long: `Manages workspaces

Workspaces are a core concept of hs. You can setup multiple workspaces
each having their own regional and environment based configuration. 
This is very convenient if you have global deployments or are working with 
mulitple customers and need to context switch frequently.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

type workspaceConfig struct {
	sync.Mutex
	Name                  string      `json:"-"`
	Version               int         `json:"Version"`
	DefaultRegion         string      `json:"DefaultRegion"`
	DefaultEnvironment    string      `json:"DefaultEnvironemnt"`
	IAMUserUUID           string      `json:"IAMUserUUID"`
	IAMAccessToken        string      `json:"IAMAccessToken"`
	IAMAccessTokenExpires int64       `json:"IAMAccessTokenExpires"`
	IAMRefreshToken       string      `json:"IAMRefreshToken"`
	IAMIDToken            string      `json:"IAMIDToken"`
	IAMRegion             string      `json:"IAMRegion"`
	IAMEnvironment        string      `json:"IAMEnvironment"`
	HASRegion             string      `json:"HASRegion"`
	HASConfig             has.Config  `json:"HASConfig"`
	IronConfig            iron.Config `json:"IronConfig"`
}

func (w *workspaceConfig) iamExpireTime() *time.Time {
	if w.IAMAccessTokenExpires == 0 {
		return nil
	}
	tm := time.Unix(w.IAMAccessTokenExpires, 0)
	return &tm
}

func (w *workspaceConfig) iamLoginExpired() bool {
	expired := w.iamExpireTime()
	if expired == nil {
		return true
	}
	return expired.Before(time.Now())
}

func (w *workspaceConfig) root() string {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	workspaceRoot := filepath.Join(home, ".hs", "workspaces")
	err = os.MkdirAll(workspaceRoot, 0700)
	if err != nil {
		fmt.Printf("workspace directory error: %v\n", err)
		os.Exit(1)
	}
	return workspaceRoot

}

func (w *workspaceConfig) configFile(name ...string) string {
	useName := w.Name
	if len(name) > 0 {
		useName = name[0]
	}
	return filepath.Join(w.root(), useName+".config.json")
}

func (w *workspaceConfig) list() ([]string, string, error) {
	workspaceList := make([]string, 0)
	glob := filepath.Join(w.root(), "*.config.json")
	list, err := filepath.Glob(glob)
	if err != nil {
		return list, "", err
	}
	for _, l := range list {
		workspaceList = append(workspaceList, workspaceName(l))
	}
	return workspaceList, w.current(), nil
}

func workspaceName(file string) string {
	baseName := filepath.Base(file)
	return strings.TrimSuffix(baseName, ".config.json")
}

func (w *workspaceConfig) current() string {
	current := filepath.Join(w.root(), "current")
	_, err := os.Lstat(current)
	if os.IsNotExist(err) {
		return ""
	}
	currentConfig, err := os.Readlink(current)
	if err != nil {
		return ""
	}
	return workspaceName(currentConfig)
}

func (w *workspaceConfig) delete(workspace string) error {
	if workspace == "default" {
		return fmt.Errorf("cannot remove default")
	}
	if workspace == w.current() {
		_ = w.setDefault("default")
	}
	return os.Remove(w.configFile(workspace))
}

func (w *workspaceConfig) save() error {
	w.Lock()
	defer w.Unlock()
	w.Version = 1
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(w.configFile(), data, 0600); err != nil {
		return err
	}
	return nil
}

func (w *workspaceConfig) load(workspace string) error {
	w.Lock()
	defer w.Unlock()
	newTarget := &workspaceConfig{}
	newTarget.Name = workspace
	data, err := ioutil.ReadFile(newTarget.configFile())
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, w)
	if err != nil {
		return err
	}
	w.Name = workspace
	return nil
}

func (w *workspaceConfig) setDefault(workspace string) error {
	w.Lock()
	defer w.Unlock()
	current := filepath.Join(w.root(), "current")
	if _, err := os.Lstat(current); err == nil {
		if err := os.Remove(current); err != nil {
			return err
		}
	}
	w.Name = workspace
	if err := os.Symlink(w.configFile(), current); err != nil {
		return err
	}
	return nil
}

func (w *workspaceConfig) ensureRoot() {
	// Check if we have a default namespace
	defaultWorkspaceConfigFile := filepath.Join(w.root(), "default.config.json")
	stat, err := os.Stat(defaultWorkspaceConfigFile)
	if os.IsNotExist(err) {
		// Create
		defaultWorkspace := &workspaceConfig{
			Name:               "default",
			Version:            1,
			DefaultRegion:      "us-east",
			DefaultEnvironment: "client-test",
		}
		if err := defaultWorkspace.save(); err != nil {
			fmt.Printf("failed to save default workspace config: %v\n", err)
			os.Exit(1)
		}
	} else {
		if stat.IsDir() {
			fmt.Printf("directory instead of file: %s\n", defaultWorkspaceConfigFile)
			os.Exit(1)
		}
	}
}

func (w *workspaceConfig) ensureDefault() {
	w.Lock()
	defer w.Unlock()
	// Check if we have a current workspace
	currentWorkspaceFile := filepath.Join(w.root(), "current")
	_, err := os.Stat(currentWorkspaceFile)
	if os.IsNotExist(err) {
		// Link to default
		if err := os.Symlink("default.config.json", currentWorkspaceFile); err != nil {
			fmt.Printf("Failed to set default workspace: %v\n", err)
			os.Exit(1)
		}
	}
}

func (w *workspaceConfig) init() error {
	w.ensureRoot()
	w.ensureDefault()
	return w.load(w.current())
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
}
