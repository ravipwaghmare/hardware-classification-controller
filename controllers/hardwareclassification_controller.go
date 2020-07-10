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

package controllers

import (
	"context"
	hwcc "hardware-classification-controller/api/v1alpha1"
	"hardware-classification-controller/hcmanager"
	"strings"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	//HWControllerName Name to show in the logs
	HWControllerName = "HardwareClassification-Controller"
)

// HardwareClassificationReconciler reconciles a HardwareClassification object
type HardwareClassificationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile reconcile function
// +kubebuilder:rbac:groups=metal3.io,resources=hardwareclassifications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io,resources=hardwareclassifications/status,verbs=get;update;patch
// Add RBAC rules to access baremetalhost resources
// +kubebuilder:rbac:groups=metal3.io,resources=baremetalhosts,verbs=get;list;update
// +kubebuilder:rbac:groups=metal3.io,resources=baremetalhosts/status,verbs=get
func (hcReconciler *HardwareClassificationReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx := context.Background()

	// Get HardwareClassificationController to get values for Namespace and ExpectedHardwareConfiguration
	hardwareClassification := &hwcc.HardwareClassification{}

	if err := hcReconciler.Client.Get(ctx, req.NamespacedName, hardwareClassification); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(hardwareClassification, hcReconciler.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	defer func() {
		// Always attempt to Patch the hardwareClassification object and status after each reconciliation.
		if err := patchHelper.Patch(ctx, hardwareClassification); err != nil {
			hcReconciler.Log.Error(err, "Failed to Patch HardwareClassification")
		}
	}()

	// Get Expected Hardware Configuration from hardwareClassification
	extractedProfile := hardwareClassification.Spec.HardwareCharacteristics
	hcReconciler.Log.Info("Expected Hardware Configuration", "Profile", extractedProfile)

	// Get the new hardware classification manager
	hcManager := hcmanager.NewHardwareClassificationManager(hcReconciler.Client, hcReconciler.Log)

	ErrValidation := hcManager.ValidateExtractedHardwareProfile(extractedProfile)
	if ErrValidation != nil {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusEmpty, hwcc.ProfileMisConfigured, ErrValidation.Error())
		hcReconciler.Log.Error(ErrValidation, ErrValidation.Error())
		return ctrl.Result{}, nil
	}

	//Fetch baremetal host list for the given namespace
	hostList, bmhList, err := hcManager.FetchBmhHostList(hardwareClassification.ObjectMeta.Namespace)
	if err != nil {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusEmpty, hwcc.FetchBMHListFailure, err.Error())
		hcReconciler.Log.Error(err, err.Error())
		return ctrl.Result{}, nil
	}

	if len(hostList) == 0 {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusEmpty, hwcc.Empty, hwcc.NoBaremetalHost)
		hcReconciler.Log.Info(hwcc.NoBaremetalHost)
		return ctrl.Result{}, nil
	}

	//Extract the hardware details from the baremetal host list
	validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(extractedProfile, hostList)
	hcReconciler.Log.Info("Validated Hardware Details", "HardwareDetails", validatedHardwareDetails)

	//Compare the host list with extracted profile and fetch the valid host names
	validHost := hcManager.MinMaxFilter(hardwareClassification.ObjectMeta.Name, validatedHardwareDetails, extractedProfile)
	hcReconciler.Log.Info("Filtered Bare metal hosts", "ValidHost", validHost)

	if len(validHost) == 0 {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusUnMatched, hwcc.Empty, hwcc.NoValidHostFound)
		hcReconciler.Log.Info("Updated profile match status", "ProfileMatchStatus", hwcc.ProfileMatchStatusUnMatched)
		deleteLabelError := hcManager.DeleteHWCCLabel(ctx, hardwareClassification.ObjectMeta, bmhList)
		if len(deleteLabelError) > 0 {
			hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusEmpty, hwcc.LabelUpdateFailure, strings.Join(deleteLabelError, ","))
		}
		return ctrl.Result{}, nil
	}

	//Update BMHost Labels
	validHost = append(validHost, validHost...)
	setLabelError := hcManager.SetLabel(ctx, hardwareClassification.ObjectMeta, validHost, bmhList)
	if len(setLabelError) > 0 {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusEmpty, hwcc.LabelUpdateFailure, strings.Join(setLabelError, ","))
	} else {
		hcmanager.SetStatus(hardwareClassification, hwcc.ProfileMatchStatusMatched, hwcc.Empty, hwcc.NOError)
		hcReconciler.Log.Info(hwcc.LabelUpdated)
		hcReconciler.Log.Info("Updated profile match status", "ProfileMatchStatus", hwcc.ProfileMatchStatusMatched)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager will add watches for this controller
func (hcReconciler *HardwareClassificationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassification{}).
		Complete(hcReconciler)
}
