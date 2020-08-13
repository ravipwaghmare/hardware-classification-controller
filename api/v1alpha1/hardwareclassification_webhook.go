/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var hardwareclassificationlog = logf.Log.WithName("hardwareclassification-resource")

func (r *HardwareClassification) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-metal3-sigs-k8s-io-v1alpha1-hardwareclassification,mutating=true,failurePolicy=fail,groups=metal3.sigs.k8s.io,resources=hardwareclassifications,verbs=create;update,versions=v1alpha1,name=mhardwareclassification.kb.io

var _ webhook.Defaulter = &HardwareClassification{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *HardwareClassification) Default() {
	hardwareclassificationlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-metal3-sigs-k8s-io-v1alpha1-hardwareclassification,mutating=false,failurePolicy=fail,groups=metal3.sigs.k8s.io,resources=hardwareclassifications,versions=v1alpha1,name=vhardwareclassification.kb.io

var _ webhook.Validator = &HardwareClassification{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *HardwareClassification) ValidateCreate() error {
	hardwareclassificationlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *HardwareClassification) ValidateUpdate(old runtime.Object) error {
	hardwareclassificationlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *HardwareClassification) ValidateDelete() error {
	hardwareclassificationlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
