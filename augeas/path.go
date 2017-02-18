package augeas

import "regexp"

var beginWithSlashMatcher = regexp.MustCompile("^/")
var bracketedIndexMatcher = regexp.MustCompile("\[(\d+)\]")
var treeValueNodeMatcher = regexp.MustCompile("/\[value\]$")
var bracketedIndexSurroundedBySlashesMatcher = regexp.MustCompile("/\[(\d+)\](?:/|$)")


func GetRegularPath(path string) string {
	rewrittenPath := beginWithSlashMatcher.ReplaceAllLiteralString(path, "")
	rewrittenPath  = bracketedIndexMatcher.ReplaceAllString(rewrittenPath, "/$1")
	return rewrittenPath
}

func GetAugeasPath(path string) string {
	rewrittenPath := treeValueNodeMatcher.ReplaceAllLiteralString(path, "/")
	rewrittenPath  = bracketedIndexSurroundedBySlashesMatcher.ReplaceAllString(rewrittenPath, "[$1]")
	return rewrittenPath
}
