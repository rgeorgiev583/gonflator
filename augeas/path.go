package augeas

import (
	"path"
	"path/filepath"
	"regexp"
)

var beginWithSlashMatcher = regexp.MustCompile("^/")
var bracketedIndexMatcher = regexp.MustCompile("\\[(\\d+)\\]")
var treeValueNodeMatcher = regexp.MustCompile("/\\[value\\]$")
var bracketedIndexSurroundedBySlashesMatcher = regexp.MustCompile("/\\[(\\d+)\\](?:/|$)")

func getFilesystemPath(augeasPath string) string {
	rewrittenPath := beginWithSlashMatcher.ReplaceAllLiteralString(augeasPath, "")
	rewrittenPath = bracketedIndexMatcher.ReplaceAllString(rewrittenPath, "/[$1]")
	return filepath.Clean(rewrittenPath)
}

func getAugeasPath(filesystemPath string) string {
	rewrittenPath := treeValueNodeMatcher.ReplaceAllLiteralString(filesystemPath, "/")
	rewrittenPath = bracketedIndexSurroundedBySlashesMatcher.ReplaceAllString(rewrittenPath, "[$1]")
	return path.Clean(rewrittenPath)
}
