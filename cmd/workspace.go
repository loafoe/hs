package cmd

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

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/dip-software/go-dip-api/iam"

	"github.com/dip-software/go-dip-api/iron"
	"github.com/mitchellh/go-homedir"
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
	IAMSelectedOrg        string      `json:"IAMSelectedOrg"`
	IAMSelectedOrgName    string      `json:"IAMSelectedOrgName"`
	IronConfig            iron.Config `json:"IronConfig"`
	S3CredsProductKey     string      `json:"S3CredsProductKey"`
	UAAToken              string      `json:"UAAAccessToken"`
	UAARefreshToken       string      `json:"UAARefreshToken"`
	UAAAccessTokenExpires int64       `json:"UAAAccessTokenExpires"`
	UAAIDToken            string      `json:"UAAIDToken"`
	PKILogicalPath        string      `json:"PKILogicalPath"`
	TFStateCreds          string      `json:"TFStateCreds"`
	TFStateInstanceURL    string      `json:"TFStateInstanceURL"`
}

func (w *workspaceConfig) iamExpireTime() *time.Time {
	if w.IAMAccessTokenExpires == 0 {
		return nil
	}
	tm := time.Unix(w.IAMAccessTokenExpires, 0)
	return &tm
}

func (w *workspaceConfig) uaaExpireTime() *time.Time {
	if w.UAAAccessTokenExpires == 0 {
		return nil
	}
	tm := time.Unix(w.UAAAccessTokenExpires, 0)
	return &tm
}

func (w *workspaceConfig) iamLoginExpired() bool {
	expired := w.iamExpireTime()
	if expired == nil {
		return true
	}
	return expired.Before(time.Now())
}

func (w *workspaceConfig) uaaLoginExpired() bool {
	expired := w.uaaExpireTime()
	if expired == nil {
		return true
	}
	return expired.Before(time.Now())
}

func workspaceRoot() string {
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
	return filepath.Join(workspaceRoot(), useName+".config.json")
}

func (w *workspaceConfig) list() ([]string, string, error) {
	workspaceList := make([]string, 0)
	glob := filepath.Join(workspaceRoot(), "*.config.json")
	list, err := filepath.Glob(glob)
	if err != nil {
		return list, "", err
	}
	for _, l := range list {
		workspaceList = append(workspaceList, workspaceName(l))
	}
	return workspaceList, currentWorkspaceName(), nil
}

func workspaceName(file string) string {
	baseName := filepath.Base(file)
	return strings.TrimSuffix(baseName, ".config.json")
}

func currentWorkspaceName() string {
	current := filepath.Join(workspaceRoot(), "current")
	_, err := os.Lstat(current)
	if os.IsNotExist(err) {
		return ""
	}
	var currentConfig string
	if runtime.GOOS != "window" {
		currentConfig, err = os.Readlink(current)
		if err != nil {
			return ""
		}
	} else {
		// Windows
		data, err := os.ReadFile(current)
		if err != nil {
			return ""
		}
		currentConfig = string(data)
	}
	return workspaceName(currentConfig)
}

func (w *workspaceConfig) delete(workspace string) error {
	if workspace == "default" {
		return fmt.Errorf("cannot remove default")
	}
	if workspace == currentWorkspaceName() {
		_ = w.setDefault("default")
	}
	return os.Remove(w.configFile(workspace))
}

func (w *workspaceConfig) saveWithIAM(client *iam.Client) error {
	token, _ := client.Token()
	w.IAMAccessToken = token
	w.IAMRefreshToken = client.RefreshToken()
	w.IAMAccessTokenExpires = client.Expires()
	return w.save()
}
func (w *workspaceConfig) save() error {
	w.Lock()
	defer w.Unlock()
	w.Version = 1
	data, err := json.Marshal(w)
	if err != nil {
		return err
	}
	if err := os.WriteFile(w.configFile(), data, 0600); err != nil {
		return err
	}
	return nil
}

func loadWorkspaceConfig(workspace string) (*workspaceConfig, error) {
	newTarget := &workspaceConfig{}
	newTarget.Name = workspace
	data, err := os.ReadFile(newTarget.configFile())
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, newTarget)
	if err != nil {
		return nil, err
	}
	newTarget.Name = workspace
	return newTarget, nil
}

func (w *workspaceConfig) setDefault(workspace string) error {
	w.Lock()
	defer w.Unlock()
	current := filepath.Join(workspaceRoot(), "current")
	if _, err := os.Lstat(current); err == nil {
		if err := os.Remove(current); err != nil {
			return err
		}
	}
	w.Name = workspace
	if runtime.GOOS != "windows" {
		if err := os.Symlink(w.configFile(), current); err != nil {
			return err
		}
	} else {
		// Windows
		if err := os.WriteFile(current, []byte(w.configFile()), 0600); err != nil {
			return err
		}
	}
	return nil
}

func ensureRoot() {
	// Check if we have a default namespace
	defaultWorkspaceConfigFile := filepath.Join(workspaceRoot(), "default.config.json")
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

func ensureDefault() {
	// Check if we have a current workspace
	currentWorkspaceFile := filepath.Join(workspaceRoot(), "current")
	_, err := os.Stat(currentWorkspaceFile)
	if os.IsNotExist(err) {
		if runtime.GOOS != "windows" {
			// Link to default
			if err := os.Symlink("default.config.json", currentWorkspaceFile); err != nil {
				fmt.Printf("Failed to set default workspace: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Windows
			if err := os.WriteFile(currentWorkspaceFile, []byte("default.config.json"), 0600); err != nil {
				fmt.Printf("Failed to set default workspace: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
}
