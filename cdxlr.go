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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CycloneDX/cyclonedx-go"
	"github.com/github/go-spdx/v2/spdxexp"
	"io"
	"net/http"
)

const (
	LicenseListSourceUrl = "https://raw.githubusercontent.com/spdx/license-list-data/main/json/licenses.json"
)

var (
	spdxLicenses LicenseList
)

// GenerateMapping extracts SPDX license identifiers from a CycloneDC license list.
// * If a SPDX expression is provided, it is parsed and the found licenses
// * Else if an ID is provided, it is verified and added to the results
// * Else if a name is passed, the best matching (at least 90% match) name, on equality, the latest is selected
// * Else if a URL is passed, the SPDX license containing the URL is selected
func GenerateMapping(licenses *cyclonedx.Licenses) ([]string, error) {
	if len(spdxLicenses.Licenses) < 1 {
		if err := loadLicenses(); err != nil {
			return nil, err
		}
	}
	licenseIds := make([]string, 0)
	for _, lin := range *licenses {
		if lin.License == nil && lin.Expression != "" {
			extracted, err := spdxexp.ExtractLicenses(lin.Expression)
			if err != nil {
				return nil, err
			}
			licenseIds = uappend(licenseIds, extracted...)
		} else if lin.License != nil {
			if lin.License.ID != "" {
				if spdxLicenses.ValidId(lin.License.ID) {
					licenseIds = uappend(licenseIds, lin.License.ID)
				}
			} else if lin.License.Name != "" {
				found, _ := spdxLicenses.FindByName(lin.License.Name)
				licenseIds = uappend(licenseIds, found.LicenseId)
			} else if lin.License.URL != "" {
				lic := spdxLicenses.FindByUrl(lin.License.URL)
				if lic != nil {
					licenseIds = uappend(licenseIds, lic.LicenseId)
				}
			}
		} else {
			return nil, fmt.Errorf("invalid license definition")
		}
	}
	return licenseIds, nil
}

// loadLicenses fetches the current list of licenses from the SPDX github page
func loadLicenses() error {
	resp, err := http.Get(LicenseListSourceUrl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, resp.Body); err != nil {
		return err
	}
	if err = json.Unmarshal(buf.Bytes(), &spdxLicenses); err != nil {
		return err
	}
	return nil
}

// uappend appends an element to a slice if not already exists.
func uappend[T comparable](slice []T, elems ...T) []T {
	var found bool
	for _, elem := range elems {
		// try to find element in slice
		found = false
		for _, f := range slice {
			if elem == f {
				found = true
				break
			}
		}
		if !found {
			slice = append(slice, elem)
		}
	}
	return slice
}
