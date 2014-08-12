// +build integration

package tests

import (
	"os"
	"testing"
	"time"

	"github.com/deis/deis/tests/utils"
)

var (
	appsCreateCmd  = "apps:create {{.AppName}}"
	appsListCmd    = "apps:list"
	appsRunCmd     = "apps:run echo hello"
	appsOpenCmd    = "apps:open --app={{.AppName}}"
	appsLogsCmd    = "apps:logs --app={{.AppName}}"
	appsInfoCmd    = "apps:info --app={{.AppName}}"
	appsDestroyCmd = "apps:destroy --app={{.AppName}} --confirm={{.AppName}}"
)

func TestApps(t *testing.T) {
	params := appsSetup(t)
	appsCreateTest(t, params)
	appsListTest(t, params, false)
	appsLogsTest(t, params)
	appsInfoTest(t, params)
	appsRunTest(t, params)
	appsOpenTest(t, params)
	limitsSetTest(t, params, 3)
	appsOpenTest(t, params)
	limitsUnsetTest(t, params, 5)
	appsDestroyTest(t, params)
	appsListTest(t, params, true)
}

func appsSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "appssample"
	utils.Execute(t, authLoginCmd, cfg, false, "")
	utils.Execute(t, gitCloneCmd, cfg, false, "")
	return cfg
}

func appsCreateTest(t *testing.T, params *utils.DeisTestConfig) {
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	cmd := appsCreateCmd
	utils.Execute(t, cmd, params, false, "")
	utils.Execute(t, cmd, params, true, "App with this Id already exists")
}

func appsDestroyTest(t *testing.T, params *utils.DeisTestConfig) {
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, appsDestroyCmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	if err := utils.Rmdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
}

func appsInfoTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Execute(t, appsInfoCmd, params, false, "")
}

func appsListTest(t *testing.T, params *utils.DeisTestConfig, notflag bool) {
	utils.CheckList(t, appsListCmd, params, params.AppName, notflag)
}

func appsLogsTest(t *testing.T, params *utils.DeisTestConfig) {
	cmd := appsLogsCmd
	utils.Execute(t, cmd, params, true, "204 NO CONTENT")
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, gitPushCmd, params, false, "")
	// TODO: nginx needs a few seconds to wake up here--fixme!
	time.Sleep(5000 * time.Millisecond)
	utils.Curl(t, params)
	utils.Execute(t, cmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
}

func appsOpenTest(t *testing.T, params *utils.DeisTestConfig) {
	utils.Curl(t, params)
}

func appsRunTest(t *testing.T, params *utils.DeisTestConfig) {
	cmd := appsRunCmd
	if err := utils.Chdir(params.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, cmd, params, false, "")
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, cmd, params, true, "Not found")
}
