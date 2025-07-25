// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"fmt"

	"github.com/gardener/gardener/extensions/pkg/controller/worker"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	api "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws"
	"github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/helper"
)

// UpdateMachineImagesStatus implements genericactuator.WorkerDelegate.
func (w *WorkerDelegate) UpdateMachineImagesStatus(ctx context.Context) error {
	if w.machineImages == nil {
		if err := w.generateMachineConfig(ctx); err != nil {
			return fmt.Errorf("unable to generate the machine config: %w", err)
		}
	}

	// Decode the current worker provider status.
	workerStatus, err := w.decodeWorkerProviderStatus()
	if err != nil {
		return fmt.Errorf("unable to decode the worker provider status: %w", err)
	}

	workerStatus.MachineImages = w.machineImages
	if err := w.updateWorkerProviderStatus(ctx, workerStatus); err != nil {
		return fmt.Errorf("unable to update worker provider status: %w", err)
	}

	return nil
}

func (w *WorkerDelegate) findMachineImage(name, version string, region string, arch *string) (string, error) {
	ami, err := helper.FindAMIForRegionFromCloudProfile(w.cloudProfileConfig, name, version, region, arch)
	if err == nil {
		return ami, nil
	}

	// Try to look up machine image in worker provider status as it was not found in componentconfig.
	if providerStatus := w.worker.Status.ProviderStatus; providerStatus != nil {
		workerStatus := &api.WorkerStatus{}
		if _, _, err := w.decoder.Decode(providerStatus.Raw, nil, workerStatus); err != nil {
			return "", fmt.Errorf("could not decode worker status of worker '%s': %w", k8sclient.ObjectKeyFromObject(w.worker), err)
		}

		machineImage, err := helper.FindMachineImage(workerStatus.MachineImages, name, version, arch)
		if err != nil {
			return "", worker.ErrorMachineImageNotFound(name, version, *arch, region)
		}

		return machineImage.AMI, nil
	}

	return "", worker.ErrorMachineImageNotFound(name, version, *arch, region)
}

func appendMachineImage(machineImages []api.MachineImage, machineImage api.MachineImage) []api.MachineImage {
	if _, err := helper.FindMachineImage(machineImages, machineImage.Name, machineImage.Version, machineImage.Architecture); err != nil {
		return append(machineImages, machineImage)
	}
	return machineImages
}
