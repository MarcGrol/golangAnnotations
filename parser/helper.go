package parser

import "go/ast"

func extractComments(commentGroup *ast.CommentGroup) []string {
	lines := make([]string, 0)
	if commentGroup != nil {
		for _, comment := range commentGroup.List {
			lines = append(lines, comment.Text)
		}
	}
	return lines
}

func extractTag(basicLit *ast.BasicLit) string {
	if basicLit != nil {
		return basicLit.Value
	}
	return ""
}
