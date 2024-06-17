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
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"strings"
)

// SimilarityThreshold defines that names must meet a threshold of at least 90%
const SimilarityThreshold = 0.9

// License represents a single license of the SPDX license list in JSON format
type License struct {
	Reference             string   `json:"reference"`
	IsDeprecatedLicenseId bool     `json:"isDeprecatedLicenseId"`
	DetailsUrl            string   `json:"detailsUrl"`
	ReferenceNumber       int      `json:"referenceNumber"`
	Name                  string   `json:"name"`
	LicenseId             string   `json:"licenseId"`
	SeeAlso               []string `json:"seeAlso"`
	IsOsiApproved         bool     `json:"isOsiApproved"`
}

// LicenseList represents the structure of the SPDX license list in JSON format
type LicenseList struct {
	LicenseListVersion string    `json:"licenseListVersion"`
	Licenses           []License `json:"licenses"`
}

// FindByName searches for licenses ba similar names
// Returns the best-matching license and the similarity percentage
func (ll *LicenseList) FindByName(name string) (License, float64) {
	var bestMatch License
	var bestRating float64 = 0

	compareMetric := metrics.NewJaroWinkler()
	for _, l := range ll.Licenses {
		if name == l.LicenseId {
			// Wrong handling of CDX: ID as name
			return l, 1
		}
		similarity := strutil.Similarity(name, l.Name, compareMetric)
		if similarity >= SimilarityThreshold && similarity >= bestRating {
			bestMatch = l
			bestRating = similarity
		}
	}
	return bestMatch, bestRating
}

// FindByUrl searches an url in all SPDX licenses in the fields reference, detailsUrl and all seeAlso additional links
func (ll *LicenseList) FindByUrl(url string) *License {
	for _, l := range ll.Licenses {
		switch url {
		case l.Reference:
			return &l
		case l.DetailsUrl:
			return &l
		default:
			for _, additionalLink := range l.SeeAlso {
				if strings.HasPrefix(additionalLink, url) {
					return &l
				}
			}
		}
	}
	return nil
}

// ValidId checks if an ID is a valid SPDX license IO
func (ll *LicenseList) ValidId(id string) bool {
	for _, l := range ll.Licenses {
		if l.LicenseId == id {
			return true
		}
	}
	return false
}
