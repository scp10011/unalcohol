package doc

import (
	"go/ast"
	"strings"
)

const (
	Description = "@DESCRIPTION"
	Tags        = "@TAGS"
	Url         = "@URL"
	Method      = "@METHOD"
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
			switch strings.ToUpper(words[1]) {
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
			case "@GET":
				doc.Method = []string{"GET"}
				doc.URL = words[2]
				break
			case "@POST":
				doc.Method = []string{"POST"}
				doc.URL = words[2]
				break
			case "@PUT":
				doc.Method = []string{"PUT"}
				doc.URL = words[2]
				break
			case "@DELETE":
				doc.Method = []string{"DELETE"}
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
