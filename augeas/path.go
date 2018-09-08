package augeas

import (
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const treeValueName = "[value]"

var beginWithSlashMatcher = regexp.MustCompile("^/")
var bracketedIndexMatcher = regexp.MustCompile("\\[(\\d+)\\]")
var treeValueNodeMatcher = regexp.MustCompile("/\\[value\\]$")
var bracketedIndexSurroundedBySlashesMatcher = regexp.MustCompile("/\\[(\\d+)\\](?:/|$)")

func isDirectory(path string) bool {
	return !strings.HasSuffix(path, treeValueName)
}

func getFilesystemPath(augeasPath string, isDir bool) string {
	//rewrittenPath := beginWithSlashMatcher.ReplaceAllLiteralString(augeasPath, "")
	rewrittenPath := bracketedIndexMatcher.ReplaceAllString(augeasPath, "/[$1]")

	if !isDir {
		rewrittenPath += "/" + treeValueName
	}

	return filepath.Clean(rewrittenPath)
}

func getAugeasPath(filesystemPath string, isDir bool) string {
	rewrittenPath := bracketedIndexSurroundedBySlashesMatcher.ReplaceAllString(filesystemPath, "[$1]")
	if !isDir {
		rewrittenPath = treeValueNodeMatcher.ReplaceAllLiteralString(rewrittenPath, "/")
	}
	//rewrittenPath = "/" + rewrittenPath
	return path.Clean(rewrittenPath)
}
