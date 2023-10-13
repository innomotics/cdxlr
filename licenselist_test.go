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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLicenseList_FilterByName(t *testing.T) {
	tests := []struct {
		name        string
		argName     string
		wantRating  float64
		wantLicense License
	}{
		{
			"MIT License",
			"MIT License",
			1,
			License{LicenseId: "MIT"},
		},
		{
			"MIT License plural",
			"MIT Licenses",
			.9833333333333333,
			License{LicenseId: "MIT"},
		},
		{
			"MIT License short",
			"MIT Lic.",
			.9022727272727272,
			License{LicenseId: "MIT"},
		},
		{
			"Apache Licenses",
			"Apache License",
			.9555555555555556,
			License{LicenseId: "Apache-2.0"},
		},
		{
			"Apache Licenses 1",
			"Apache License 1",
			.9777777777777777,
			License{LicenseId: "Apache-1.1"},
		},
		{
			"Apache License too new",
			"Apache License 3.0",
			.9777777777777777,
			License{LicenseId: "Apache-2.0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, loadLicenses())
			l, rate := spdxLicenses.FindByName(tt.argName)
			assert.Equal(t, rate, tt.wantRating)
			assert.Equal(t, l.LicenseId, tt.wantLicense.LicenseId)
			fmt.Printf("%s (%f)\n", l.Name, rate)
		})
	}
}

func TestLicenseList_FindByUrl(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantLicense *License
	}{
		{
			"Find by 'referenceId'",
			"https://spdx.org/licenses/CC-BY-SA-4.0.html",
			&License{LicenseId: "CC-BY-SA-4.0"},
		},
		{
			"Find by 'detailsUrl'",
			"https://spdx.org/licenses/BSD-4.3TAHOE.json",
			&License{LicenseId: "BSD-4.3TAHOE"},
		},
		{
			"Find by first 'seeAlso'",
			"https://github.com/microsoft/Computational-Use-of-Data-Agreement/blob/master/C-UDA-1.0.md",
			&License{LicenseId: "C-UDA-1.0"},
		},
		{
			"Find by base of 'seeAlso'",
			"https://git.savannah.gnu.org/cgit/indent.git/tree/doc/indent.texi",
			&License{LicenseId: "BSD-4.3TAHOE"},
		},
		{
			"No valid URL found",
			"https://git.savannah.gnu.org/cgit/indent.git/tree/doc/extent.texi",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, loadLicenses())
			result := spdxLicenses.FindByUrl(tt.url)
			if tt.wantLicense == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.wantLicense.LicenseId, result.LicenseId)
			}
		})
	}
}

func TestLicenseList_ValidId(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{
			"Incomplete Apache license",
			"Apache",
			false,
		},
		{
			"Incomplete Apache 1 license",
			"Apache-1",
			false,
		},
		{
			"Incomplete Apache 1 license",
			"Apache-1.1",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, loadLicenses())
			assert.Equalf(t, tt.want, spdxLicenses.ValidId(tt.id), "ValidId(%v)", tt.id)
		})
	}
}
