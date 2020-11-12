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

package hcmanager

import (
	bmh "github.com/metal3-io/baremetal-operator/apis/metal3.io/v1alpha1"
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"
)

// MinMaxFilter it will perform the minimum and maximum comparison based on the value provided by the user and check for the valid host
func (mgr HardwareClassificationManager) MinMaxFilter(ProfileName string, HostList []bmh.HardwareDetails, expectedHardwareprofile hwcc.HardwareCharacteristics) []string {
	var validHosts []string
	for _, hardwareDetail := range HostList {
		if !checkCPUCount(mgr, hardwareDetail.CPU, expectedHardwareprofile.Cpu, hardwareDetail.Hostname) ||
			!checkRAM(mgr, hardwareDetail.RAMMebibytes, expectedHardwareprofile.Ram, hardwareDetail.Hostname) ||
			!checkNICS(mgr, len(hardwareDetail.NIC), expectedHardwareprofile.Nic, hardwareDetail.Hostname) ||
			!checkFirmware(mgr, hardwareDetail.Firmware, expectedHardwareprofile.Firmware, hardwareDetail.Hostname) ||
			!checkSystemVendor(mgr, hardwareDetail.SystemVendor, expectedHardwareprofile.SystemVendor, hardwareDetail.Hostname) ||
			!checkDiskDetails(mgr, hardwareDetail.Storage, expectedHardwareprofile.Disk, hardwareDetail.Hostname) {
			continue
		}
		validHosts = append(validHosts, hardwareDetail.Hostname)
	}
	return validHosts
}

//checkCPUCount this function checks the CPU details for both min and max parameters
func checkCPUCount(mgr HardwareClassificationManager, cpu bmh.CPU, expectedCPU *hwcc.Cpu, bmhName string) bool {
	if expectedCPU == nil {
		return true
	}
	if expectedCPU.MaximumCount > 0 {
		expectedMaxCPUCount := expectedCPU.MaximumCount
		mgr.Log.Info("Maximum CPU Count", "BareMetalHost", bmhName, "Expected", expectedMaxCPUCount, "Actual", cpu.Count)
		if expectedMaxCPUCount < cpu.Count {
			mgr.Log.Info("CPU Count Mismatched", "BareMetalHost", bmhName)
			return false
		}
	}
	if expectedCPU.MinimumCount > 0 {
		expectedMinCPUCount := expectedCPU.MinimumCount
		mgr.Log.Info("Minimum CPU Count", "BareMetalHost", bmhName, "Expected", expectedMinCPUCount, "Actual", cpu.Count)
		if expectedMinCPUCount > cpu.Count {
			mgr.Log.Info("CPU Count Mismatched", "BareMetalHost", bmhName)
			return false
		}
	}
	if expectedCPU.MaximumSpeedMHz > 0 {
		expectedMaxSpeedHz := bmh.ClockSpeed(expectedCPU.MaximumSpeedMHz)
		mgr.Log.Info("Maximum CPU ClockSpeed", "BareMetalHost", bmhName, "Expected", expectedMaxSpeedHz, "Actual", cpu.ClockMegahertz)
		if expectedMaxSpeedHz < cpu.ClockMegahertz {
			mgr.Log.Info("CPU ClockSpeed Mismatched", "BareMetalHost", bmhName)
			return false
		}
	}
	if expectedCPU.MinimumSpeedMHz > 0 {
		expectedMinSpeedHz := bmh.ClockSpeed(expectedCPU.MinimumSpeedMHz)
		mgr.Log.Info("Minimum CPU ClockSpeed", "BareMetalHost", bmhName, "Expected", expectedMinSpeedHz, "Actual", cpu.ClockMegahertz)
		if expectedMinSpeedHz > cpu.ClockMegahertz {
			mgr.Log.Info("CPU ClockSpeed Mismatched", "BareMetalHost", bmhName)
			return false
		}
	}
	return true
}

//checkNICS this function checks the nics details for both min and max parameters
func checkNICS(mgr HardwareClassificationManager, nics int, expectedNIC *hwcc.Nic, bmhName string) bool {
	if expectedNIC == nil {
		return true
	}
	if expectedNIC.MaximumCount > 0 {
		expectedMaxNicCount := expectedNIC.MaximumCount
		mgr.Log.Info("Maximum NIC Count", "BareMetalHost", bmhName, "Expected", expectedMaxNicCount, "Actual", nics)
		if expectedMaxNicCount < nics {
			mgr.Log.Info("NIC Count Mismatched", "BareMetalHost", bmhName)
			return false
		}
	}
	if expectedNIC.MinimumCount > 0 {
		expectedMinNicCount := expectedNIC.MinimumCount
		mgr.Log.Info("Minimum NIC Count", "BareMetalHost", bmhName, "Expected", expectedMinNicCount, "Actual", nics)
		if expectedMinNicCount > nics {
			mgr.Log.Info("NIC Count Mismatched", "BareMetalHost", bmhName)
			return false
		}
	}
	return true
}

//checkRAM this function checks the ram details for both min and max parameters
func checkRAM(mgr HardwareClassificationManager, ram int, expectedRAM *hwcc.Ram, bmhName string) bool {
	if expectedRAM == nil {
		return true
	}
	if expectedRAM.MaximumSizeGB > 0 {
		expectedMaxRAM := expectedRAM.MaximumSizeGB
		mgr.Log.Info("Maximum RAM Size", "BareMetalHost", bmhName, "Expected", expectedMaxRAM, "Actual", ram)
		if expectedMaxRAM < ram {
			mgr.Log.Info("RAM Size Mismatched", "BareMetalHost", bmhName)
			return false
		}
	}
	if expectedRAM.MinimumSizeGB > 0 {
		expectedMinRAM := expectedRAM.MinimumSizeGB
		mgr.Log.Info("Minimum RAM Size", "BareMetalHost", bmhName, "Expected", expectedMinRAM, "Actual", ram)
		if expectedMinRAM > ram {
			mgr.Log.Info("RAM Size Mismatched", "BareMetalHost", bmhName)
			return false
		}
	}
	return true
}

//checkDiskDetails this function checks the Disk details for both min and max parameters
func checkDiskDetails(mgr HardwareClassificationManager, disks []bmh.Storage, expectedDisk *hwcc.Disk, bmhName string) bool {
	if expectedDisk == nil {
		return true
	}
	if expectedDisk.MaximumCount > 0 {
		expectedMaxDiskCount := expectedDisk.MaximumCount
		mgr.Log.Info("Maximum Disk Count", "BareMetalHost", bmhName, "Expected", expectedMaxDiskCount, "Actual", len(disks))
		if expectedMaxDiskCount < len(disks) {
			mgr.Log.Info("Disk Count Mismatched", "BareMetalHost", bmhName)
			return false
		}
	}
	if expectedDisk.MinimumCount > 0 {
		expectedMinDiskCount := expectedDisk.MinimumCount
		mgr.Log.Info("Minimum Disk Count", "BareMetalHost", bmhName, "Expected", expectedMinDiskCount, "Actual", len(disks))
		if expectedMinDiskCount > len(disks) {
			mgr.Log.Info("Disk Count Mismatched", "BareMetalHost", bmhName)
			return false
		}
	}
	for _, disk := range disks {
		if expectedDisk.MaximumIndividualSizeGB > 0 {
			expectedMaxDiskSize := bmh.Capacity(expectedDisk.MaximumIndividualSizeGB)
			mgr.Log.Info("Maximum Disk Size", "BareMetalHost", bmhName, "Expected", expectedMaxDiskSize, "Actual", disk.SizeBytes)
			if expectedMaxDiskSize < disk.SizeBytes {
				mgr.Log.Info("Disk Size Mismatched", "BareMetalHost", bmhName)
				return false
			}
		}
		if expectedDisk.MinimumIndividualSizeGB > 0 {
			expectedMinDiskSize := bmh.Capacity(expectedDisk.MinimumIndividualSizeGB)
			mgr.Log.Info("Minimum Disk Size", "BareMetalHost", bmhName, "Expected", expectedMinDiskSize, "Actual", disk.SizeBytes)
			if expectedMinDiskSize > disk.SizeBytes {
				mgr.Log.Info("Disk Size Mismatched", "BareMetalHost", bmhName)
				return false
			}
		}
	}
	return true
}

//checkFirmware this function checks the Firmware details
func checkFirmware(mgr HardwareClassificationManager, firmware bmh.Firmware, expectedFirmware *hwcc.Firmware, bmhName string) bool {
	if expectedFirmware == nil {
		return true
	}
	if expectedFirmware.Bios.Vendor != firmware.BIOS.Vendor && expectedFirmware.Bios.Version != firmware.BIOS.Version {
		mgr.Log.Info("Firmware Vendor", "BareMetalHost", bmhName, "Expected", expectedFirmware.Bios.Vendor, "Actual", firmware.BIOS.Vendor)
		mgr.Log.Info("Firmware Version", "BareMetalHost", bmhName, "Expected", expectedFirmware.Bios.Version, "Actual", firmware.BIOS.Version)
		return false
	}
	return true
}

//checkSystemvendor this function checks the SystemVendor details
func checkSystemVendor(mgr HardwareClassificationManager, systemVendor bmh.HardwareSystemVendor, expectedSystemVendor *hwcc.SystemVendor, bmhName string) bool {
	if expectedSystemVendor == nil {
		return true
	}
	if expectedSystemVendor.Manufacturer != systemVendor.Manufacturer && expectedSystemVendor.ProductName != systemVendor.ProductName {
		mgr.Log.Info("System Vendor Manufacturer", "BareMetalHost", bmhName, "Expected", expectedSystemVendor.Manufacturer, "Actual", systemVendor.Manufacturer)
		mgr.Log.Info("System Vendor Product Name", "BareMetalHost", bmhName, "Expected", expectedSystemVendor.ProductName, "Actual", systemVendor.ProductName)
		return false
	}
	return true
}
