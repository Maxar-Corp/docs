package file

import (
	"path/filepath"
	"os"
	"strings"
	"io"
	"os/exec"
	"fmt"
	"net/http"
	"github.com/gruntwork-io/docs/errors"
	"io/ioutil"
	"github.com/gruntwork-io/docs/logger"
)

// Convert the given path into a path relative to basePath
func GetPathRelativeTo(path string, basePath string) (string, error) {
	inputFolderAbs, err := filepath.Abs(basePath)
	if err != nil {
		return "", errors.WithStackTrace(err)
	}

	fileAbs, err := filepath.Abs(path)
	if err != nil {
		return "", errors.WithStackTrace(err)
	}

	relPath, err := filepath.Rel(inputFolderAbs, fileAbs)
	if err != nil {
		return "", errors.WithStackTrace(err)
	}

	return relPath, nil
}

// Return true if the path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Return the size, in bytes, of the given file
func GetFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, errors.WithStackTrace(err)
	}
	return fileInfo.Size(), nil
}

// Create a directory and all the parent directories at the given path
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}

// There is no way to know for sure if a file is text or binary. The best we can do is use various heuristics to guess.
// The best set of heuristics is in the Unix/Linux file command, so we use that if it's available. Otherwise, we turn
// to Go's HTTP package. For more info, see: http://stackoverflow.com/q/16760378/483528
//
// Note: for empty files, with no bytes to read to tell the type of content it contains, guessing is particularly hard.
// Since an empty file can be anything it wants to be, this function will return true for empty files!
func IsTextFile(path string) (bool, error) {
	if !PathExists(path) {
		return false, errors.WithStackTrace(NoSuchFile(path))
	}

	size, err := GetFileSize(path)
	if err != nil {
		return false, errors.WithStackTrace(err)
	}
	if size == 0 {
		return true, nil
	}

	mimeType, err := GuessMimeType(path)
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(mimeType, "text"), nil
}

// Copy the given file. If a file already exists at dstPath, return an error.
func CopyFile(srcPath, dstPath string) error {
	containingDir := getContainingDirectory(dstPath)

	err := CreateDir(containingDir)
	if err != nil {
		return errors.WithStackTrace(fmt.Errorf("Error while making directory %s", containingDir))
	}

	// TODO: Uncomment this since it's for dev only.
	//if isFileExist(dstPath) {
	//	return errors.WithStackTrace(fmt.Errorf("A file already exists at the path %s. Overwriting existing files is not permiitted. Most likely, another file with a conflicting name was already written to this location.\n", dstPath))
	//}

	bytes, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return errors.WithStackTrace(err)
	}

	err = ioutil.WriteFile(dstPath, bytes, os.ModePerm)
	if err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}

// Copy all files from the given srcPath to the given dstPath
func CopyFiles(srcPath, dstPath string) error {
	if ! IsDir(srcPath) {
		return errors.WithStackTrace(NotADirectory(srcPath))
	}

	// If the dstPath dir doesn't exist, create it
	err := CreateDir(dstPath)
	if err != nil {
		return errors.WithStackTrace(err)
	}

	// Copy all files in srcPath into dstPath
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, fileErr error) error {
		relPath, err := GetPathRelativeTo(path, srcPath)
		if err != nil {
			return errors.WithStackTrace(err)
		}

		fileSrcPath := filepath.Join(srcPath, relPath)
		fileDstPath := filepath.Join(dstPath, relPath)

		if ! IsDir(fileSrcPath) {
			err = CopyFile(fileSrcPath, fileDstPath)
			if err != nil {
				return errors.WithStackTrace(err)
			}
		}

		return nil
	})
}

// Read the body of the file at the given path
func ReadFile(srcPath string) (string, error) {
	bytes, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return "", errors.WithStackTrace(err)
	}

	return string(bytes[:]), nil
}

// Write the given file. If a file already exists at dstPath, return an error.
func WriteFile(body string, dstPath string) error {
	containingDir := getContainingDirectory(dstPath)

	err := CreateDir(containingDir)
	if err != nil {
		return errors.WithStackTrace(fmt.Errorf("Error while making directory %s", containingDir))
	}

	// TODO: Uncomment this since it's for dev only.
	//if isFileExist(dstPath) {
	//	return errors.WithStackTrace(fmt.Errorf("A file already exists at the path %s. Overwriting existing files is not permiitted. Most likely, another file with a conflicting name was already written to this location.\n", dstPath))
	//}

	bodyAsBytes := []byte(body)
	err = ioutil.WriteFile(dstPath, bodyAsBytes, os.ModePerm)
	if err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}

// Return true if the given path is a directory. We panic versus returning an error because this error is unlikely to occur
// and it's painful to take four lines just to call an IsDir() function.
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		logger.Logger.Panicf("ERROR: IsDir() failed while attempting to read %s\n", path)
	}

	return fileInfo.IsDir()
}

// Return true if the file at the given path exists
func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getContainingDirectory(path string) string {
	return filepath.Dir(path)
}

func mkDirRecursive(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// Guess the mime type for the given file using a variety of heuristics. Under the hood, uses the Unix/Linux file
// command, if available, and Go's HTTP package otherwise.
func GuessMimeType(path string) (string, error) {
	if CommandInstalled("file") {
		return guessMimeTypeUsingFileCommand(path)
	} else {
		return guessMimeTypeUsingGoHttpPackage(path)
	}
}

// Use the Unix/Linux "file" command to determine the mime type. This performs a number of checks and tends to do a
// good job with most files.
func guessMimeTypeUsingFileCommand(path string) (string, error) {
	return RunCommandAndGetOutput("file", "-b", "--mime", path)
}

// Use a package built into Go for detecting the mime type of arbitrary content. In my experience, it doesn't work
// very well.
func guessMimeTypeUsingGoHttpPackage(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errors.WithStackTrace(err)
	}
	defer file.Close()

	first512Bytes := make([]byte, 512)
	numBytesRead, err := file.Read(first512Bytes)
	if err != nil && err != io.EOF {
		return "", errors.WithStackTrace(err)
	}

	// If it's an empty file, there is no real distinction, so default to "false", as there is not much processing
	// you can do on an empty file anyway
	if numBytesRead == 0 {
		return "", nil
	}

	return http.DetectContentType(first512Bytes), nil
}

// Run the given command return its stdout and stderr as a string
func RunCommandAndGetOutput(command string, args ... string) (string, error) {
	cmd := exec.Command(command, args...)

	bytes, err := cmd.Output()
	if err != nil {
		return "", errors.WithStackTrace(err)
	}

	return string(bytes), nil
}

// Return true if the OS has the given command installed
func CommandInstalled(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// custom error types

type NoSuchFile string
func (path NoSuchFile) Error() string {
	return fmt.Sprintf("File %s does not exist", string(path))
}

type NotADirectory string
func (dir NotADirectory) Error() string {
	return fmt.Sprintf("The path %s is not a directory.\n", string(dir))
}