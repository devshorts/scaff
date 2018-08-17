package scaff

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"math"
	"github.com/sirupsen/logrus"
	"fmt"
	"io/ioutil"
)

type FileResolver struct {
}

func NewFileResolver() FileResolver {
	return FileResolver{}
}

func (f FileResolver) GetAllDirs(path string) []string {
	var dirs []string

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, path)
		}

		return nil
	})

	sort.Slice(dirs, func(i, j int) bool { return dirs[i] > dirs[j] })

	return dirs
}

func (f FileResolver) GetAllFiles(path string) []os.FileInfo {
	var files []os.FileInfo

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, info)
		}

		return nil
	})

	return files
}

func (f FileResolver) popSegment(path string) (string, string) {
	split := strings.Split(path, "/")

	idx := int(math.Max(float64(len(split)-1), 0))

	remaining := split[0:idx]

	return split[idx], strings.Join(remaining, "/")
}

func (f FileResolver) TemplatePath(path string, formatter RuleRunner, dryRun bool) {
	segment, remaining := f.popSegment(path)

	updated := formatter.Replace(segment)

	newPath := filepath.Join(remaining, updated)

	if newPath != path {
		logrus.Info(fmt.Sprintf("Updating %s to %s. DryRun %t", path, newPath, dryRun))

		if !dryRun {
			os.Rename(path, newPath)
		}
	}
}
func (f FileResolver) TemplateFile(info os.FileInfo, runner RuleRunner, b bool) {
	ioutil.ReadFile(info.Name())
}
