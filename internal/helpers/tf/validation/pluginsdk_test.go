// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Modifications made on 2025-08-14

package validation

import (
	"testing"

	"github.com/valiparsa/terraform-provider-azuread/internal/helpers/tf/pluginsdk"
)

func TestValidateFloatInSlice(t *testing.T) {
	cases := map[string]struct {
		Value                  interface{}
		ValidateFunc           pluginsdk.SchemaValidateFunc
		ExpectValidationErrors bool
	}{
		"accept valid value": {
			Value:                  1.5,
			ValidateFunc:           FloatInSlice([]float64{1.0, 1.5, 2.0}),
			ExpectValidationErrors: false,
		},
		"accept valid negative value ": {
			Value:                  -1.0,
			ValidateFunc:           FloatInSlice([]float64{-1.0, 2.0}),
			ExpectValidationErrors: false,
		},
		"accept zero": {
			Value:                  0.0,
			ValidateFunc:           FloatInSlice([]float64{0.0, 2.0}),
			ExpectValidationErrors: false,
		},
		"reject out of range value": {
			Value:                  -1.0,
			ValidateFunc:           FloatInSlice([]float64{0.0, 2.0}),
			ExpectValidationErrors: true,
		},
		"reject incorrectly typed value": {
			Value:                  1,
			ValidateFunc:           FloatInSlice([]float64{0, 1, 2}),
			ExpectValidationErrors: true,
		},
	}

	for tn, tc := range cases {
		_, errors := tc.ValidateFunc(tc.Value, tn)
		if len(errors) > 0 && !tc.ExpectValidationErrors {
			t.Errorf("%s: unexpected errors %s", tn, errors)
		} else if len(errors) == 0 && tc.ExpectValidationErrors {
			t.Errorf("%s: expected errors but got none", tn)
		}
	}
}
