# cdxlr - CycloneDX License Resolver

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This library supports mapping different types of valid CycloneDX licenses to valid SPDX License-IDs.

The initial problem is that the [CycloneDX Components Licenses](https://cyclonedx.org/docs/1.4/json/#metadata_component_licenses)
format allows for a various definition of licenses.
On the other hand, for automated tooling, the licenses must be provided in a standardized format.
The most widely used one is [SPDX Reference Format](https://spdx.github.io/spdx-spec/v2-draft/SPDX-license-list/) with
the SPDX License identifiers.

This library uses the official [License List Data](https://github.com/spdx/license-list-data) from github and tries to
match the CycloneDX License to an SPDX ID by different potential matchers. The matching tries first to match the more
qualified matchers and if no match was found, uses the next-well-qualifying one.

The implemented Matchers are:
* If no `license` field, but an `expression` field is set, the contained licenses are extracted using the [spdxexp](https://github.com/github/go-spdx) library.
* If a `license` field is set
  * if the `id` field is set, it is validated as valid SPDX identifier and used.
  * if the `name` field is set, it is [compared](https://github.com/adrg/strutil) with the SPDX licenses name and the best-matching result is used.
  * if the `url` field is set, the url is being tried to match the SPDX licenses' `reference`, `detailsUrl`, and `seeAlso` fields. 

This way, most generated CycloneDX files can be used in automated CI/CD pipelines

## Basic operation

Install the package by installing the module:
```shell
go get github.com/innomotics/cyclonedx-license-resolver
```

In your code, you can then use the library through the one-stop-shop method `GenerateMapping`:
```go
var cdxLicenses *cyclonedx.Licenses
...
spdxList, err := cdxlr.GenerateMapping(cdxLicenses)
if err != nil {
	log.Fatalf("error parsing CycloneDX licenses: %v\n", err)
}
fmt.Prinln(spdxList)
```

## License
This library is released under the Apache License Version 2.0 (see [LICENSE](./LICENSE.txt)).

## Contribution
Any contributions welcome. Please suggest any enhancements as issues.