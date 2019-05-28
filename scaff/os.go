package scaff

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devshorts/scaff/scaff/config"
	"github.com/devshorts/scaff/scaff/file"
	"github.com/devshorts/scaff/scaff/lang"
	"github.com/sirupsen/logrus"
)

type FileResolver struct {
	config config.FileConfig
}

const DEFAULT_DELIM = "__"

func NewTemplator(config config.FileConfig) FileResolver {
	return FileResolver{config}
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

type FileData struct {
	FileInfo os.FileInfo
	Path     string
}

func (f FileResolver) GetAllFiles(path string) []FileData {
	var files []FileData

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, FileData{
				FileInfo: info,
				Path:     path,
			})
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

	pathDelimiter := DEFAULT_DELIM

	updated := formatter.Replace(segment, pathDelimiter)

	newPath := filepath.Join(remaining, updated)

	if newPath != path {
		logrus.Info(fmt.Sprintf("Updating %s to %s. DryRun %t", path, newPath, dryRun))

		if !dryRun {
			_, remaining := f.popSegment(newPath)

			if !file.Exists(remaining) {
				originalPermissions, _ := os.Stat(path)

				os.MkdirAll(remaining, originalPermissions.Mode())
			}

			os.Rename(path, newPath)
		}
	}
}

func (f FileResolver) TemplateFile(info FileData, runner RuleRunner, dryRun bool) {
	fileBytes, _ := ioutil.ReadFile(info.Path)

	tokenDelimiter := DEFAULT_DELIM

	fileExtension := filepath.Ext(info.Path)

	// if we have a specific delimiter for rules in this file type use it
	if delim, ok := f.config.FileDelims[fileExtension]; ok {
		tokenDelimiter = delim
	}

	contents := string(fileBytes)

	result := runner.Replace(contents, tokenDelimiter)

	switch fileExtension {
	case ".go":
		if f.config.LanguageRules.Go != nil {
			rules := *f.config.LanguageRules.Go
			result = lang.NewGoProcessor(rules, runner.ctx).Process(result)
		}
	}

	if contents != result {
		logrus.Info(fmt.Sprintf("Updating %s. DryRun %t", info.Path, dryRun))

		if !dryRun {
			ioutil.WriteFile(info.Path, []byte(result), info.FileInfo.Mode())
		}
	}

	// apply rules to the filename after processing file contents
	f.TemplatePath(info.Path, runner, dryRun)
}
