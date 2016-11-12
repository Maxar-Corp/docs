package docfile

import (
	"fmt"
	"regexp"
	"github.com/gruntwork-io/docs/docs-preprocessor/errors"
	"github.com/gruntwork-io/docs/docs-preprocessor/logger"
	"github.com/gruntwork-io/docs/docs-preprocessor/file"
)

const IS_PACKAGE_OVERVIEW_DOC_REGEX = `^packages/([\s\w -]+)/README.md$`
const IS_PACKAGE_OVERVIEW_DOC_REGEX_NUM_CAPTURE_GROUPS = 1

// Represents a non-overview document that's part of a specific module.
type PackageOverviewDoc struct {
	relPath string
	absPath string
}

func NewPackageOverviewDoc(absPath string, relPath string) *PackageOverviewDoc {
	return &PackageOverviewDoc{ absPath: absPath, relPath: relPath }
}

func (d *PackageOverviewDoc) IsMatch() bool {
	return checkRegex(d.relPath, IS_PACKAGE_OVERVIEW_DOC_REGEX)
}

func (d *PackageOverviewDoc) Copy(outputPathRoot string) error {
	outRelPath, err := d.getRelOutputPath()
	if err != nil {
		return errors.WithStackTrace(err)
	}

	outAbsPath := fmt.Sprintf("%s/%s", outputPathRoot, outRelPath)

	logger.Logger.Printf("Copying PACKAGE-OVERVIEW-DOC file %s to %s\n", d.absPath, outAbsPath)
	err = file.CopyFile(d.absPath, outAbsPath)
	if err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}

func (d *PackageOverviewDoc) getRelOutputPath() (string, error) {
	var outputPath string

	regex := regexp.MustCompile(IS_PACKAGE_OVERVIEW_DOC_REGEX)
	submatches := regex.FindAllStringSubmatch(d.relPath, -1)

	if len(submatches) == 0 || len(submatches[0]) != IS_PACKAGE_OVERVIEW_DOC_REGEX_NUM_CAPTURE_GROUPS + 1 {
		return outputPath, errors.WithStackTrace(&WrongNumberOfCaptureGroupsFoundInPathRegEx{ docTypeName: "PackageDocOverview", path: d.relPath, regEx: IS_PACKAGE_OVERVIEW_DOC_REGEX })
	}

	// If we were parsing d.relPath = packages/package-vpc/README.md...
	packageName := submatches[0][1] // = package-vpc

	return fmt.Sprintf("packages/%s/overview.md", packageName), nil
}