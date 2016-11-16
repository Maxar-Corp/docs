package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/docs/globs"
)

const GENERATOR_TESTS_FIXTURES_PATH = "test-fixtures/generator-tests"

func TestShouldSkipPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		path      string
		inputPath string
		excludes  []string
		expected  bool
	}{
		{"", "", []string{}, true},
		{".", ".", []string{}, true},
		{"foo/bar/baz", "foo/bar/baz", []string{}, true},
		{"foo/bar/baz/blah", "foo/bar/baz", []string{}, false},
		{"foo/bar/baz/blah", "foo/bar", []string{"*"}, true},
		{"foo/bar/baz/blah", "foo/bar", []string{"**"}, true},
		{"foo/bar/baz/blah", "foo/bar", []string{"*.*"}, false},
		{"foo/bar/baz/blah", "foo/bar", []string{"some/other/path"}, false},
		{"foo/bar/baz/blah", "foo/bar", []string{"foo/**/blah"}, true},
		{"foo/bar/baz/blah", "foo/bar", []string{"foo/**/abc"}, false},
	}

	for _, testCase := range testCases {
		globs, err := globs.ToGlobs(testCase.excludes)
		assert.Nil(t, err, "Failed to compile glob patterns: %s", testCase.excludes)

		actual := shouldSkipPath(testCase.path, &Opts{InputPath: testCase.inputPath, Excludes: globs})
		assert.Equal(t, testCase.expected, actual, "path = %s, inputPath = %s, excludes = %s", testCase.path, testCase.inputPath, testCase.excludes)
	}
}

// func TestProcessDocumentationFile(t *testing.T) {
// 	t.Parallel()

// 	testCases := []struct {
// 		inputFixturePath          string
// 		expectedOutputFixturePath string
// 	}{
// 		{"documentation.txt", "documentation-output.txt"},
// 		{"documentation-no-urls.md", "documentation-no-urls-output.md"},
// 		{"documentation-with-urls.md", "documentation-with-urls-output.md"},
// 		{"logo.png", "logo-output.png"},
// 	}

// 	for _, testCase := range testCases {
// 		actual, err := getContentsForDocumentationFile(testCase.inputFixturePath, &Opts{InputPath: GENERATOR_TESTS_FIXTURES_PATH})
// 		assert.Nil(t, err, "Error processing file %s: %v", testCase.inputFixturePath, err)

// 		expected := readProcessorTestsFixturesFile(t, testCase.expectedOutputFixturePath)
// 		assert.Equal(t, expected, actual, "inputFixturePath = %s, expectedOutputFixturePath = %s.\nExpected:\n%s\nActual:\n%s", testCase.inputFixturePath, testCase.expectedOutputFixturePath, string(expected), string(actual))
// 	}
// }

// func readProcessorTestsFixturesFile(t *testing.T, file string) []byte {
// 	bytes, err := ioutil.ReadFile(path.Join(GENERATOR_TESTS_FIXTURES_PATH, file))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	return bytes
// }
