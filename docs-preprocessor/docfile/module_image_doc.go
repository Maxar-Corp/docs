package docfile

import (
	"fmt"
	"regexp"
	"github.com/gruntwork-io/docs/docs-preprocessor/errors"
	"github.com/gruntwork-io/docs/docs-preprocessor/logger"
	"github.com/gruntwork-io/docs/docs-preprocessor/file"
)

const IS_MODULE_IMAGE_DOC_REGEX = `^packages/([\w -]+)/modules/_images/([\w -]+\.(jpg|jpeg|gif|png|svg))$`
const IS_MODULE_IMAGE_DOC_REGEX_NUM_CAPTURE_GROUPS = 3

// Represents a non-overview document that's part of a specific module.
type ModuleImageDoc struct {
	relPath string
	absPath string
}

func NewModuleImageDoc(absPath string, relPath string) *ModuleImageDoc {
	return &ModuleImageDoc{ absPath: absPath, relPath: relPath }
}

func (d *ModuleImageDoc) IsMatch() bool {
	return checkRegex(d.relPath, IS_MODULE_IMAGE_DOC_REGEX)
}

func (d *ModuleImageDoc) Copy(outputPathRoot string) error {
	outRelPath, err := d.getRelOutputPath()
	if err != nil {
		return errors.WithStackTrace(err)
	}

	outAbsPath := fmt.Sprintf("%s/%s", outputPathRoot, outRelPath)

	logger.Logger.Printf("Copying MODULE-IMAGE-DOC file %s to %s\n", d.absPath, outAbsPath)
	err = file.CopyFile(d.absPath, outAbsPath)
	if err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}

func (d *ModuleImageDoc) getRelOutputPath() (string, error) {
	var outputPath string

	regex := regexp.MustCompile(IS_MODULE_IMAGE_DOC_REGEX)
	submatches := regex.FindAllStringSubmatch(d.relPath, -1)

	if len(submatches) == 0 || len(submatches[0]) != IS_MODULE_IMAGE_DOC_REGEX_NUM_CAPTURE_GROUPS + 1 {
		return outputPath, errors.WithStackTrace(&WrongNumberOfCaptureGroupsFoundInPathRegEx{ docTypeName: "ModuleImageDoc", path: d.relPath, regEx: IS_MODULE_IMAGE_DOC_REGEX })
	}

	// If we were parsing d.relPath = packages/package-vpc/modules/_images/sample.jpg...
	packageName := submatches[0][1] // = package-vpc
	imageName := submatches[0][2] // = sample.jpg

	outputPath = fmt.Sprintf("packages/%s/_images/%s", packageName, imageName)

	return outputPath, nil
}