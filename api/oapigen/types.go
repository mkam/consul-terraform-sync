// Package oapigen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.3 DO NOT EDIT.
package oapigen

import (
	"encoding/json"
	"fmt"
)

// BufferPeriod defines model for BufferPeriod.
type BufferPeriod struct {
	Enabled *bool   `json:"enabled,omitempty"`
	Max     *string `json:"max,omitempty"`
	Min     *string `json:"min,omitempty"`
}

// CatalogServicesCondition defines model for CatalogServicesCondition.
type CatalogServicesCondition struct {
	Datacenter       *string                            `json:"datacenter,omitempty"`
	Namespace        *string                            `json:"namespace,omitempty"`
	NodeMeta         *CatalogServicesCondition_NodeMeta `json:"node_meta,omitempty"`
	Regexp           string                             `json:"regexp"`
	UseAsModuleInput *bool                              `json:"use_as_module_input,omitempty"`
}

// CatalogServicesCondition_NodeMeta defines model for CatalogServicesCondition.NodeMeta.
type CatalogServicesCondition_NodeMeta struct {
	AdditionalProperties map[string]string `json:"-"`
}

// Condition defines model for Condition.
type Condition struct {
	CatalogServices *CatalogServicesCondition `json:"catalog_services,omitempty"`
	ConsulKv        *ConsulKVCondition        `json:"consul_kv,omitempty"`
	Schedule        *ScheduleCondition        `json:"schedule,omitempty"`
	Services        *ServicesCondition        `json:"services,omitempty"`
}

// ConsulKVCondition defines model for ConsulKVCondition.
type ConsulKVCondition struct {
	Datacenter       *string `json:"datacenter,omitempty"`
	Namespace        *string `json:"namespace,omitempty"`
	Path             string  `json:"path"`
	Recurse          *bool   `json:"recurse,omitempty"`
	UseAsModuleInput *bool   `json:"use_as_module_input,omitempty"`
}

// ConsulKVModuleInput defines model for ConsulKVModuleInput.
type ConsulKVModuleInput struct {
	Datacenter *string `json:"datacenter,omitempty"`
	Namespace  *string `json:"namespace,omitempty"`
	Path       string  `json:"path"`
	Recurse    *bool   `json:"recurse,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Error     Error     `json:"error"`
	RequestId RequestID `json:"request_id"`
}

// ModuleInput defines model for ModuleInput.
type ModuleInput struct {
	ConsulKv *ConsulKVModuleInput `json:"consul_kv,omitempty"`
	Services *ServicesModuleInput `json:"services,omitempty"`
}

// RequestID defines model for RequestID.
type RequestID string

// Run defines model for Run.
type Run struct {
	// Whether or not infrastructure changes were detected during task inspection.
	ChangesPresent *bool   `json:"changes_present,omitempty"`
	Plan           *string `json:"plan,omitempty"`

	// Enterprise only. URL of Terraform Cloud run that corresponds to the task run.
	TfcRunUrl *string `json:"tfc_run_url,omitempty"`
}

// ScheduleCondition defines model for ScheduleCondition.
type ScheduleCondition struct {
	Cron string `json:"cron"`
}

// ServicesCondition defines model for ServicesCondition.
type ServicesCondition struct {
	CtsUserDefinedMeta *ServicesCondition_CtsUserDefinedMeta `json:"cts_user_defined_meta,omitempty"`
	Datacenter         *string                               `json:"datacenter,omitempty"`
	Filter             *string                               `json:"filter,omitempty"`
	Names              *[]string                             `json:"names,omitempty"`
	Namespace          *string                               `json:"namespace,omitempty"`
	Regexp             *string                               `json:"regexp,omitempty"`
	UseAsModuleInput   *bool                                 `json:"use_as_module_input,omitempty"`
}

// ServicesCondition_CtsUserDefinedMeta defines model for ServicesCondition.CtsUserDefinedMeta.
type ServicesCondition_CtsUserDefinedMeta struct {
	AdditionalProperties map[string]string `json:"-"`
}

// ServicesModuleInput defines model for ServicesModuleInput.
type ServicesModuleInput struct {
	Names  *[]string `json:"names,omitempty"`
	Regexp *string   `json:"regexp,omitempty"`
}

// Task defines model for Task.
type Task struct {
	BufferPeriod *BufferPeriod `json:"buffer_period,omitempty"`
	Condition    *Condition    `json:"condition,omitempty"`
	Description  *string       `json:"description,omitempty"`
	Enabled      *bool         `json:"enabled,omitempty"`
	Module       string        `json:"module"`
	ModuleInput  *ModuleInput  `json:"module_input,omitempty"`
	Name         string        `json:"name"`
	Providers    *[]string     `json:"providers,omitempty"`
	Services     *[]string     `json:"services,omitempty"`
	Variables    *VariableMap  `json:"variables,omitempty"`
	Version      *string       `json:"version,omitempty"`
}

// TaskDeleteResponse defines model for TaskDeleteResponse.
type TaskDeleteResponse struct {
	RequestId RequestID `json:"request_id"`
}

// TaskRequest defines model for TaskRequest.
type TaskRequest Task

// TaskResponse defines model for TaskResponse.
type TaskResponse struct {
	RequestId RequestID `json:"request_id"`
	Run       *Run      `json:"run,omitempty"`
	Task      *Task     `json:"task,omitempty"`
}

// VariableMap defines model for VariableMap.
type VariableMap struct {
	AdditionalProperties map[string]string `json:"-"`
}

// CreateTaskJSONBody defines parameters for CreateTask.
type CreateTaskJSONBody TaskRequest

// CreateTaskParams defines parameters for CreateTask.
type CreateTaskParams struct {
	// Different modes for running. Supports run now which runs the task immediately
	// and run inspect which creates a dry run task that is inspected and discarded
	// at the end of the inspection.
	Run *CreateTaskParamsRun `json:"run,omitempty"`
}

// CreateTaskParamsRun defines parameters for CreateTask.
type CreateTaskParamsRun string

// CreateTaskJSONRequestBody defines body for CreateTask for application/json ContentType.
type CreateTaskJSONRequestBody CreateTaskJSONBody

// Getter for additional properties for CatalogServicesCondition_NodeMeta. Returns the specified
// element and whether it was found
func (a CatalogServicesCondition_NodeMeta) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for CatalogServicesCondition_NodeMeta
func (a *CatalogServicesCondition_NodeMeta) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for CatalogServicesCondition_NodeMeta to handle AdditionalProperties
func (a *CatalogServicesCondition_NodeMeta) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for CatalogServicesCondition_NodeMeta to handle AdditionalProperties
func (a CatalogServicesCondition_NodeMeta) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for ServicesCondition_CtsUserDefinedMeta. Returns the specified
// element and whether it was found
func (a ServicesCondition_CtsUserDefinedMeta) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ServicesCondition_CtsUserDefinedMeta
func (a *ServicesCondition_CtsUserDefinedMeta) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ServicesCondition_CtsUserDefinedMeta to handle AdditionalProperties
func (a *ServicesCondition_CtsUserDefinedMeta) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ServicesCondition_CtsUserDefinedMeta to handle AdditionalProperties
func (a ServicesCondition_CtsUserDefinedMeta) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for VariableMap. Returns the specified
// element and whether it was found
func (a VariableMap) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for VariableMap
func (a *VariableMap) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for VariableMap to handle AdditionalProperties
func (a *VariableMap) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for VariableMap to handle AdditionalProperties
func (a VariableMap) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}
