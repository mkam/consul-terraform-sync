package api

import (
	"net/http"

	"github.com/hashicorp/consul-terraform-sync/api/oapigen"
	"github.com/hashicorp/consul-terraform-sync/logging"
)

// DeleteTaskByName deletes an existing task and its events. Does not delete
// if the task is active.
func (h *TaskLifeCycleHandler) DeleteTaskByName(w http.ResponseWriter, r *http.Request, name string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	ctx := r.Context()
	requestID := requestIDFromContext(ctx)
	logger := logging.FromContext(r.Context()).Named(createTaskSubsystemName).With("task_name", name)
	logger.Trace("delete task request", "delete_task_request", name)

	// Check if task exists
	_, err := h.ctrl.Task(ctx, name)
	if err != nil {
		logger.Trace("task not found", "error", err)
		sendError(w, r, http.StatusNotFound, err)
		return
	}

	err = h.ctrl.TaskDelete(ctx, name)
	if err != nil {
		sendError(w, r, http.StatusInternalServerError, err)
		return
	}

	writeResponse(w, r, http.StatusAccepted, oapigen.TaskResponse{
		RequestId: requestID,
	})
}
