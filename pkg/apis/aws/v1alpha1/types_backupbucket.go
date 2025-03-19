// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RetentionType defines the level at which immutability properties are applied on objects.
type RetentionType string

const (
	// BucketLevelImmutability sets the immutability feature on the bucket level.
	BucketLevelImmutability RetentionType = "bucket"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BackupBucketConfig represents the configuration for a backup bucket.
type BackupBucketConfig struct {
	metav1.TypeMeta `json:",inline"`

	// Immutability defines the immutability configuration for the backup bucket.
	// +optional
	Immutability *ImmutableConfig `json:"immutability,omitempty"`
}

// ImmutableConfig represents the immutability configuration for a backup bucket.
type ImmutableConfig struct {
	// RetentionType specifies the type of retention for the backup bucket.
	// Currently allowed value is:
	// - "bucket": retention policy applies on the entire bucket.
	RetentionType RetentionType `json:"retentionType"`

	// RetentionPeriod specifies the immutability retention period for the backup bucket.
	// S3 only supports immutability durations in days or years, therefore this field must be set as multiple of 24h.
	RetentionPeriod metav1.Duration `json:"retentionPeriod"`

	// S3 provides two retention modes that apply different levels of protection to objects:
	// allowed valus are: "Governance" or "Compliance" mode.
	Mode string `json:"mode"`
}
