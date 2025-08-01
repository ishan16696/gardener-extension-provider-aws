// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package client

func copySlice[T any](array []T) []T {
	var cp []T
	return append(cp, array...)
}
