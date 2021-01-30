package tools

import (
	"os"
	"path"
	"path/filepath"
)

// ErrorCheckWithPanic : panics on error
func ErrorCheckWithPanic(err error, message string) {
	if err != nil {
		panic(err.Error() + " : " + message)
	}
}

// GetResourcePath : Gets the path to a resource file
func GetResourcePath(directory, file string) string {
	dir := path.Join(GetExecutablePath(), directory, file)
	if checkFileExists(dir) {
		return dir
	}
	dir = path.Join(GetExecutablePath(),"..", directory, file)
	if checkFileExists(dir) {
		return dir
	}

	return ""
}

// GetExecutablePath : Returns the path of the executable
func GetExecutablePath() string {
	executable, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(executable)
}

func checkFileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}