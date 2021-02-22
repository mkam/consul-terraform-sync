package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/hashicorp/consul-terraform-sync/config"
	"github.com/hashicorp/consul-terraform-sync/driver"
	"github.com/hashicorp/consul-terraform-sync/event"
	"github.com/mitchellh/mapstructure"
)

const taskPath = "tasks"

// taskHandler handles the tasks endpoint
type taskHandler struct {
	store   *event.Store
	drivers map[string]driver.Driver
	version string
}

// newTaskHandler returns a new taskHandler
func newTaskHandler(store *event.Store, drivers map[string]driver.Driver,
	version string) *taskHandler {

	return &taskHandler{
		store:   store,
		drivers: drivers,
		version: version,
	}
}

// ServeHTTP serves the tasks endpoint
func (h *taskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[TRACE] (api.task) requesting tasks '%s'", r.URL.Path)

	switch r.Method {
	case http.MethodPatch:
		h.updateTask(w, r)
	default:
		err := fmt.Errorf("'%s' in an unsupported method. The task API "+
			"currently supports the method(s): '%s'", r.Method, http.MethodPatch)
		log.Printf("[TRACE] (api.task) unsupported method: %s", err)
		jsonErrorResponse(w, http.StatusMethodNotAllowed, err)
	}
}

// UpdateTaskConfig contains the fields available for patch updating a task.
// Not all task configuration is available for update
type UpdateTaskConfig struct {
	Enabled *bool `mapstructure:"enabled"`
}

// updateTask does a patch update to an existing task
func (h *taskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	taskName, err := getTaskName(r.URL.Path, taskPath, h.version)
	if err != nil {
		log.Printf("[TRACE] (api.task) bad request: %s", err)
		jsonErrorResponse(w, http.StatusBadRequest, err)

		return
	}

	if taskName == "" {
		err := fmt.Errorf("No task name was included in the api request. " +
			"Updating a task requires the task name: '/v1/tasks/:task_name'")
		log.Printf("[TRACE] (api.task) bad request: %s", err)
		jsonErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	d, ok := h.drivers[taskName]
	if !ok {
		err := fmt.Errorf("A task with the name '%s' does not exist or has not "+
			"been initialized yet", taskName)
		log.Printf("[TRACE] (api.task) task not found: %s", err)
		jsonErrorResponse(w, http.StatusNotFound, err)
		return
	}

	runOp, err := runOption(r)
	if err != nil {
		log.Printf("[TRACE] (api.task) unsupported run option: %s", err)
		jsonErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[TRACE] (api.task) unable to read request body from update "+
			"request for task '%s': %s", taskName, err)
		jsonErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	conf, err := decodeBody(body)
	if err != nil {
		log.Printf("[TRACE] (api.task) problem decoding body from update request "+
			"for task '%s': %s", taskName, err)
		jsonErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	patch := driver.PatchTask{
		RunOption: runOp,
	}
	if conf.Enabled != nil {
		log.Printf("[INFO] (api.task) Updating task '%s' to be enabled=%t",
			taskName, config.BoolVal(conf.Enabled))
		patch.Enabled = config.BoolVal(conf.Enabled)
	}

	ctx := context.Background()

	plan, err := d.UpdateTask(ctx, patch)
	if err != nil {
		log.Printf("[TRACE] (api.task) error while updating task '%s': %s",
			taskName, err)
		jsonErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse(w, http.StatusOK, plan)
}

func decodeBody(body []byte) (UpdateTaskConfig, error) {
	var raw map[string]interface{}

	err := json.Unmarshal(body, &raw)
	if err != nil {
		return UpdateTaskConfig{}, err
	}

	var config UpdateTaskConfig
	var md mapstructure.Metadata
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		ErrorUnused:      false,
		Metadata:         &md,
		Result:           &config,
	})
	if err != nil {
		return UpdateTaskConfig{}, err
	}

	if err = decoder.Decode(raw); err != nil {
		return UpdateTaskConfig{}, err
	}

	if len(md.Unused) > 0 {
		sort.Strings(md.Unused)
		err := fmt.Errorf("request body's JSON contains unsupported keys: %s",
			strings.Join(md.Unused, ", "))
		return UpdateTaskConfig{}, err
	}

	return config, nil
}

// runOption returns a run option for updating the task
func runOption(r *http.Request) (string, error) {
	// `?run=<option>` parameter
	const runKey = "run"

	keys, ok := r.URL.Query()[runKey]
	if !ok {
		return "", nil
	}

	if len(keys) != 1 {
		return "", fmt.Errorf("cannot support more than one run query "+
			"parameter, got run values: %v", keys)
	}

	value := keys[0]
	value = strings.ToLower(value)
	switch value {
	case driver.RunOptionNow, driver.RunOptionInspect:
		return value, nil
	default:
		return "", fmt.Errorf("unsupported run parameter value. only "+
			"supporting run values %s and %s but got %s",
			driver.RunOptionNow, driver.RunOptionInspect, value)
	}
}
