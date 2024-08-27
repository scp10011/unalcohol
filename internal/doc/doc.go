package doc

import (
	"go/ast"
	"strings"
)

const (
	Description = "@Description"
	Tags        = "@Tags"
	Url         = "@URL"
	Method      = "@Method"
)

type Doc struct {
	Summary     string
	Description string
	Tags        []string
	URL         string
	Method      []string
	Name        string
}

func ParseDoc(group *ast.CommentGroup) *Doc {
	doc := Doc{}
	for _, comment := range group.List {
		words := strings.Fields(comment.Text)
		if len(words) > 2 {
			switch words[1] {
			case Tags:
				doc.Tags = strings.Split(words[2], ",")
				break
			case Description:
				doc.Description = words[2]
				break
			case Method:
				doc.Method = strings.Split(strings.ToUpper(words[2]), ",")
				break
			case Url:
				doc.URL = words[2]
				break
			default:
				if doc.Name == "" {
					doc.Name = words[1]
					doc.Summary = words[2]
				}
			}
		}
	}
	if doc.URL != "" {
		return &doc
	} else {
		return nil
	}
}
