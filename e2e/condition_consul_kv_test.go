// +build e2e

package e2e

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/consul-terraform-sync/api"
	"github.com/hashicorp/consul-terraform-sync/templates/tftmpl"
	"github.com/hashicorp/consul-terraform-sync/testutils"
	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type kvTaskOpts struct {
	path              string
	recurse           bool
	sourceIncludesVar bool
}

func kvTaskConfig(taskName string, opts kvTaskOpts) string {
	var module string
	if opts.sourceIncludesVar {
		module = "consul_kv_file"
	} else {
		module = "local_instances_file"
	}

	conditionTask := fmt.Sprintf(`task {
		name = "%s"
		services = ["web", "api"]
		source = "./test_modules/%s"
		condition "consul-kv" {
			path = "%s"
			source_includes_var = %t
			recurse = %t
		}
	}
	`, taskName, module, opts.path, opts.sourceIncludesVar, opts.recurse)
	return conditionTask
}

func ctsSetup(t *testing.T, srv *testutil.TestServer, tempDir string, taskConfig string) *api.Client {
	cleanup := testutils.MakeTempDir(t, tempDir)
	t.Cleanup(func() {
		cleanup()
	})

	config := baseConfig(tempDir).appendConsulBlock(srv).appendTerraformBlock().
		appendString(taskConfig)
	configPath := filepath.Join(tempDir, configFile)
	config.write(t, configPath)

	cts, stop := api.StartCTS(t, configPath)
	t.Cleanup(func() {
		stop(t)
	})

	err := cts.WaitForAPI(defaultWaitForAPI)
	require.NoError(t, err)

	return cts
}

/*
Source Existing New
Single
Recurse

add, update, remove

key exists initially
key doesn't exist initially

recurse, vars

false, false
false true
make sure that it doesn't take a sub key

true false
true true
make sure that it does take a sub key
make sure it doesn't take a rando key


Checking for none:
service change doesn't affect
bad key
*/

// Tests new key
func TestConditionConsulKV_NewKeyIncludesVar(t *testing.T) {
	t.Parallel()

	// Set up Consul server
	srv := newTestConsulServer(t)
	t.Cleanup(func() {
		srv.Stop()
	})

	// Configure and start CTS
	taskName := "consul_kv_condition_new"
	tempDir := fmt.Sprintf("%s%s", tempDirPrefix, taskName)
	path := "test-key"
	config := kvTaskConfig(taskName, kvTaskOpts{
		path:              path,
		recurse:           false,
		sourceIncludesVar: true,
	})
	cts := ctsSetup(t, srv, tempDir, config)

	// 0. Confirm only one event. Confirm empty var consul_kv
	eventCountBase := eventCount(t, taskName, cts.Port())
	require.Equal(t, 1, eventCountBase)

	workingDir := fmt.Sprintf("%s/%s", tempDir, taskName)
	content := testutils.CheckFile(t, true, workingDir, tftmpl.TFVarsFilename)
	assert.Contains(t, content, "consul_kv = {\n}")

	resourcesPath := filepath.Join(workingDir, resourcesDir)
	testutils.CheckFile(t, true, resourcesPath, "api.txt")
	testutils.CheckFile(t, true, resourcesPath, "web.txt")
	testutils.CheckFile(t, false, resourcesPath, "db.txt")

	// Deregister a service, confirm no event
	testutils.DeregisterConsulService(t, srv, "web")
	time.Sleep(defaultWaitForNoEvent)
	eventCountNow := eventCount(t, taskName, cts.Port())
	require.Equal(t, eventCountBase, eventCountNow,
		"change in event count. task was unexpectedly triggered")
	testutils.CheckFile(t, true, resourcesPath, "web.txt")

	// Create key that is monitored by task, check for event
	now := time.Now()
	value := "test-value"
	srv.SetKVString(t, path, value)
	api.WaitForEvent(t, cts, taskName, now, defaultWaitForEvent)
	eventCountNow = eventCount(t, taskName, cts.Port())
	eventCountBase++
	require.Equal(t, eventCountBase, eventCountNow,
		"event count did not increment once. task was not triggered as expected")
	content = testutils.CheckFile(t, true, resourcesPath, fmt.Sprintf("%s.txt", path))
	assert.Equal(t, value, content)
	testutils.CheckFile(t, false, resourcesPath, "web.txt")

	// Update the key value, check for event
	now = time.Now()
	value = "new-test-value"
	srv.SetKVString(t, path, value)
	api.WaitForEvent(t, cts, taskName, now, defaultWaitForEvent)
	eventCountNow = eventCount(t, taskName, cts.Port())
	eventCountBase++
	require.Equal(t, eventCountBase, eventCountNow,
		"event count did not increment once. task was not triggered as expected")
	content = testutils.CheckFile(t, true, resourcesPath, fmt.Sprintf("%s.txt", path))
	assert.Equal(t, value, content)

	// Remove the key, check for event
	now = time.Now()
	testutils.DeleteKV(t, srv, path)
	api.WaitForEvent(t, cts, taskName, now, defaultWaitForEvent)
	eventCountNow = eventCount(t, taskName, cts.Port())
	eventCountBase++
	require.Equal(t, eventCountBase, eventCountNow,
		"event count did not increment once. task was not triggered as expected")
	content = testutils.CheckFile(t, false, resourcesPath, fmt.Sprintf("%s.txt", path))

	// Add a key prefixed by the monitored path, check for no event
	now = time.Now()
	prefixedKey := fmt.Sprintf("%s/sub-key", path)
	srv.SetKVString(t, prefixedKey, "test")
	time.Sleep(defaultWaitForNoEvent)
	eventCountNow = eventCount(t, taskName, cts.Port())
	require.Equal(t, eventCountBase, eventCountNow,
		"change in event count. task was unexpectedly triggered")
	testutils.CheckFile(t, false, resourcesPath, fmt.Sprintf("%s.txt", prefixedKey))
}

// source includes var false
func TestConditionConsulKV_NewKey(t *testing.T) {
	t.Parallel()

	// Set up Consul server
	srv := newTestConsulServer(t)
	t.Cleanup(func() {
		srv.Stop()
	})

	// Configure and start CTS
	taskName := "consul_kv_condition_new"
	tempDir := fmt.Sprintf("%s%s", tempDirPrefix, taskName)
	path := "test-key"
	config := kvTaskConfig(taskName, kvTaskOpts{
		path:              path,
		recurse:           false,
		sourceIncludesVar: false,
	})
	cts := ctsSetup(t, srv, tempDir, config)

	// 0. Confirm only one event. Confirm empty var consul_kv
	eventCountBase := eventCount(t, taskName, cts.Port())
	require.Equal(t, 1, eventCountBase)

	workingDir := fmt.Sprintf("%s/%s", tempDir, taskName)
	content := testutils.CheckFile(t, true, workingDir, tftmpl.TFVarsFilename)
	assert.NotContains(t, content, "consul_kv = {")

	resourcesPath := filepath.Join(workingDir, resourcesDir)
	testutils.CheckFile(t, true, resourcesPath, "api.txt")
	testutils.CheckFile(t, true, resourcesPath, "web.txt")
	testutils.CheckFile(t, false, resourcesPath, "db.txt")

	// Deregister a service, confirm no event
	testutils.DeregisterConsulService(t, srv, "web")
	time.Sleep(defaultWaitForNoEvent)
	eventCountNow := eventCount(t, taskName, cts.Port())
	require.Equal(t, eventCountBase, eventCountNow,
		"change in event count. task was unexpectedly triggered")
	testutils.CheckFile(t, true, resourcesPath, "web.txt")

	// Create key that is monitored by task, check for event
	now := time.Now()
	value := "test-value"
	srv.SetKVString(t, path, value)
	api.WaitForEvent(t, cts, taskName, now, defaultWaitForEvent)
	eventCountNow = eventCount(t, taskName, cts.Port())
	eventCountBase++
	require.Equal(t, eventCountBase, eventCountNow,
		"event count did not increment once. task was not triggered as expected")

	// Update the key value, check for event
	now = time.Now()
	value = "new-test-value"
	srv.SetKVString(t, path, value)
	api.WaitForEvent(t, cts, taskName, now, defaultWaitForEvent)
	eventCountNow = eventCount(t, taskName, cts.Port())
	eventCountBase++
	require.Equal(t, eventCountBase, eventCountNow,
		"event count did not increment once. task was not triggered as expected")

	// Remove the key, check for event
	now = time.Now()
	testutils.DeleteKV(t, srv, path)
	api.WaitForEvent(t, cts, taskName, now, defaultWaitForEvent)
	eventCountNow = eventCount(t, taskName, cts.Port())
	eventCountBase++
	require.Equal(t, eventCountBase, eventCountNow,
		"event count did not increment once. task was not triggered as expected")

	// Add a key prefixed by the monitored path, check for no event
	now = time.Now()
	prefixedKey := fmt.Sprintf("%s/sub-key", path)
	srv.SetKVString(t, prefixedKey, "test")
	time.Sleep(defaultWaitForNoEvent)
	eventCountNow = eventCount(t, taskName, cts.Port())
	require.Equal(t, eventCountBase, eventCountNow,
		"change in event count. task was unexpectedly triggered")
}

func TestConditionConsulKV_ExistingKey(t *testing.T) {
	t.Parallel()

	// Set up Consul server, add a key
	srv := newTestConsulServer(t)
	t.Cleanup(func() {
		srv.Stop()
	})
	path := "test-key"
	value := "test-value"
	srv.SetKVString(t, path, value)

	// Configure and start CTS
	taskName := "consul_kv_condition_existing"
	tempDir := fmt.Sprintf("%s%s", tempDirPrefix, taskName)
	config := kvTaskConfig(taskName, kvTaskOpts{
		path:              path,
		recurse:           false,
		sourceIncludesVar: true,
	})
	cts := ctsSetup(t, srv, tempDir, config)

	// 0. Confirm only one event
	eventCountBase := eventCount(t, taskName, cts.Port())
	require.Equal(t, 1, eventCountBase)

	workingDir := fmt.Sprintf("%s/%s", tempDir, taskName)
	resourcesPath := filepath.Join(workingDir, resourcesDir)

	content := testutils.CheckFile(t, true, resourcesPath, fmt.Sprintf("%s.txt", path))
	assert.Equal(t, value, content)

	testutils.CheckFile(t, true, resourcesPath, "api.txt")
	testutils.CheckFile(t, true, resourcesPath, "web.txt")
	testutils.CheckFile(t, false, resourcesPath, "db.txt")

	// Update key that is monitored by task, check for event
	now := time.Now()
	value = "new-test-value"
	srv.SetKVString(t, path, value)
	api.WaitForEvent(t, cts, taskName, now, defaultWaitForEvent)
	eventCountNow := eventCount(t, taskName, cts.Port())
	eventCountBase++
	require.Equal(t, eventCountBase, eventCountNow,
		"event count did not increment once. task was not triggered as expected")
	content = testutils.CheckFile(t, true, resourcesPath, fmt.Sprintf("%s.txt", path))
	assert.Equal(t, value, content)

}

func TestConditionConsulKV_Register(t *testing.T) {
	t.Parallel()

	// Set up Consul server, add a key
	srv := newTestConsulServer(t)
	t.Cleanup(func() {
		srv.Stop()
	})
	path := "test-key"
	value := "test-value"
	srv.SetKVString(t, path, value)

	// Configure and start CTS
	taskName := "consul_kv_condition_existing"
	tempDir := fmt.Sprintf("%s%s", tempDirPrefix, taskName)
	config := kvTaskConfig(taskName, kvTaskOpts{
		path:              path,
		recurse:           false,
		sourceIncludesVar: true,
	})
	cts := ctsSetup(t, srv, tempDir, config)

	// 0. Confirm only one event
	eventCountBase := eventCount(t, taskName, cts.Port())
	require.Equal(t, 1, eventCountBase)

	workingDir := fmt.Sprintf("%s/%s", tempDir, taskName)
	resourcesPath := filepath.Join(workingDir, resourcesDir)

	content := testutils.CheckFile(t, true, resourcesPath, fmt.Sprintf("%s.txt", path))
	assert.Equal(t, value, content)

	testutils.CheckFile(t, true, resourcesPath, "api.txt")
	testutils.CheckFile(t, true, resourcesPath, "web.txt")
	testutils.CheckFile(t, false, resourcesPath, "db.txt")

	// Update key that is monitored by task, check for event
	now := time.Now()
	value = "new-test-value"
	srv.SetKVString(t, path, value)
	api.WaitForEvent(t, cts, taskName, now, defaultWaitForEvent)
	eventCountNow := eventCount(t, taskName, cts.Port())
	eventCountBase++
	require.Equal(t, eventCountBase, eventCountNow,
		"event count did not increment once. task was not triggered as expected")
	content = testutils.CheckFile(t, true, resourcesPath, fmt.Sprintf("%s.txt", path))
	assert.Equal(t, value, content)

}
