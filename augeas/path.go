package augeas

import "regexp"

type matchers struct {
	BeginWithSlash,
	BracketedIndex,
	TreeValueNode,
	BracketedIndexSurroundedBySlashes *regexp.Regexp
}

var matcherInstances matchers

func initMatcherInstances() {
	if matcherInstances != nil {
		return
	}

	matcherInstances = &matchers{
		BeginWithSlash: regexp.MustCompile("^/"),
		BracketedIndex: regexp.MustCompile("\[(\d+)\]"),
		TreeValueNode: regexp.MustCompile("/\[\]$"),
		BracketedIndexSurroundedBySlashes: regexp.MustCompile("/\[(\d+)\](?:/|$)"),
	}
}

func GetRegularPath(path string) string {
	initMatcherInstances()
	rewrittenPath := matcherInstances.BeginWithSlash.ReplaceAllLiteralString(path, "")
	rewrittenPath = matcherInstances.BracketedIndex.ReplaceAllString(rewrittenPath, "/$1")
	return rewrittenPath
}

func GetAugeasPath(path string) string {
	initMatcherInstances()
	rewrittenPath := matcherInstances.TreeValueNode.ReplaceAllLiteralString(path, "/")
	rewrittenPath = matcherInstances.BracketedIndexSurroundedBySlashes.ReplaceAllString(rewrittenPath, "[$1]")
	return rewrittenPath
}
