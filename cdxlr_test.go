package cdxlr

/*
 * cdxlr - CycloneDX License Resolver
 *
 * Copyright (c) Innomotics GmbH, 2023
 *
 * Authors:
 *  Mathias Haimerl <mathias.haimerl@innomotics.com>
 *
 * This work is licensed under the terms of the Apache 2.0 license.
 * See the LICENSE.txt file in the top-level directory.
 *
 * SPDX-License-Identifier:	Apache-2.0
 */
import (
	"fmt"
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadLicenses(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"Standard success",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadLicenses(); (err != nil) != tt.wantErr {
				t.Errorf("loadLicenses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateMapping(t *testing.T) {
	tests := []struct {
		name     string
		licenses *cyclonedx.Licenses
		want     []string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "All valid license types contained",
			licenses: &cyclonedx.Licenses{
				{License: &cyclonedx.License{
					ID: "MIT",
				}},
				{License: &cyclonedx.License{
					Name: "Apache License 1.1",
				}},
				{License: &cyclonedx.License{
					URL: "https://spdx.org/licenses/BSD-4-Clause.html",
				}},
				{License: &cyclonedx.License{
					URL: "http://ecos.sourceware.org/old-license.html",
				}},
				{License: &cyclonedx.License{
					URL: "https://github.com/chromium/octane/blob/master/crypto.js",
				}},
				{Expression: "Ruby AND (SAX-PD OR LGPL-2.0-only WITH FLTK-Exception)"},
			},
			want:    []string{"MIT", "Apache-1.1", "BSD-4-Clause", "RHeCos-1.1", "MIT-Wu", "LGPL-2.0-only WITH FLTK-exception", "Ruby", "SAX-PD"},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid expression: small connector",
			licenses: &cyclonedx.Licenses{
				{Expression: "Apache-2.0 and MIT"},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Invalid expression: plus connector",
			licenses: &cyclonedx.Licenses{
				{Expression: "Apache-2.0+MIT"},
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Valid expression: all the same",
			licenses: &cyclonedx.Licenses{
				{Expression: "MIT AND MIT OR (MIT OR MIT)"},
			},
			want:    []string{"MIT"},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid License: Invalid License ID",
			licenses: &cyclonedx.Licenses{
				{License: &cyclonedx.License{ID: "Fnord"}},
			},
			want:    []string{},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid License: Only License Text provided",
			licenses: &cyclonedx.Licenses{
				{
					License: &cyclonedx.License{
						Text: &cyclonedx.AttachedText{
							Content: "This could either represent an MIT or an Apache-2.0 License",
						},
					},
				},
			},
			want:    []string{},
			wantErr: assert.NoError,
		},
		{
			name:     "Invalid License: Empty License Information",
			licenses: &cyclonedx.Licenses{{}},
			want:     nil,
			wantErr:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateMapping(tt.licenses)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateMapping(%v)", tt.licenses)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GenerateMapping(%v)", tt.licenses)
		})
	}
}
