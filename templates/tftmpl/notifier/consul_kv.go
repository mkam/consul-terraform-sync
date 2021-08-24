package notifier

import (
	"log"

	"github.com/hashicorp/consul-terraform-sync/templates"
	"github.com/hashicorp/hcat/dep"
)

// ConsulKV is a custom notifier expected to be used
// for a template that contains consulKVNotifier template function
// (tmplfunc) and any other tmplfuncs e.g. services tmplfunc.
//
// This notifier only notifies on changes to Consul KV pairs and once-mode.
// It suppresses notifications for changes to other tmplfuncs.
type ConsulKV struct {
	templates.Template

	// count all dependencies needed to complete once-mode
	once     bool
	depTotal int
	counter  int
}

// NewConsulKV creates a new ConsulKVNotifier.
// serviceCount parameter: the number of services the task is configured with
func NewConsulKV(tmpl templates.Template, serviceCount int) *ConsulKV {
	return &ConsulKV{
		Template: tmpl,
		depTotal: serviceCount + 1, // for additional Consul KV dep
	}
}

// Notify notifies when a Consul KV pair or set of pairs changes.
//
// Notifications are sent when:
// A. There is a change in the Consul KV dependency for a single key pair in
//    which only the value of the key pair is returned (dep.KvValue)
// B. There is a change in the Consul KV dependency for a set of key pairs in
//    which a list of key pairs is returned ([]*dep.KeyPair)
// C. All the dependencies have been received for the first time. This is
//    regardless of the dependency type that "completes" having received all the
//    dependencies. Note: this is a special notification sent to handle a race
//    condition that causes hanging during once-mode
//
// Notification are suppressed when:
//  - Other types of dependencies that are not Consul KV. For example,
//    Services ([]*dep.HealthService).
func (n *ConsulKV) Notify(d interface{}) (notify bool) {
	log.Printf("[DEBUG] (notifier.cs) received dependency change type %T", d)
	notify = false

	if exists, ok := d.(dep.KVExists); ok {
		log.Printf("[DEBUG] (notifier.cs) notify Consul KV pair change")
		notify = true

		if !n.once && bool(exists) {
			// expect a KvValue change for once mode
			n.depTotal++
		}
	}

	if !n.once {
		n.counter++
		// after all dependencies are received, notify so once-mode can complete
		if n.counter >= n.depTotal {
			log.Printf("[DEBUG] (notifier.cs) notify once-mode complete")
			n.once = true
			notify = true
		}
	}

	if _, ok := d.(dep.KvValue); ok {
		log.Printf("[DEBUG] (notifier.cs) notify Consul KV pair change")
		notify = true
	}

	if _, ok := d.([]*dep.KeyPair); ok {
		log.Printf("[DEBUG] (notifier.cs) notify Consul KV pair change")
		notify = true
	}

	if notify {
		n.Template.Notify(d)
	}

	return notify
}
