package piwigotools

import (
	"path/filepath"
	"sort"
)

type Tree interface {
	AddNode(string) Tree
	FlatView() chan string
	TreeView() chan string
}

type node struct {
	Name     string
	Children []*node
}

func NewTree() Tree {
	return &node{
		Name: ".",
	}
}

func (t *node) AddNode(name string) Tree {
	n := &node{Name: name}
	t.Children = append(t.Children, n)
	return n
}

func (t *node) FlatView() (out chan string) {
	out = make(chan string)
	go func() {
		defer close(out)
		var flatten func(string, *node)

		flatten = func(path string, t *node) {
			switch t.Children {
			case nil:
				out <- path
			default:
				sort.Slice(t.Children, func(i, j int) bool {
					return t.Children[i].Name < t.Children[j].Name
				})
				for _, child := range t.Children {
					flatten(filepath.Join(path, child.Name), child)
				}
			}
		}

		flatten("", t)
	}()
	return out
}

func (t *node) TreeView() (out chan string) {
	out = make(chan string)
	treeLinkChar := "│   "
	treeMidChar := "├── "
	treeEndChar := "└── "
	treeAfterEndChar := "    "

	go func() {
		defer close(out)

		var tree func(string, *node)

		tree = func(prefix string, t *node) {
			for i, st := range t.Children {
				switch i {
				case len(t.Children) - 1:
					out <- prefix + treeEndChar + st.Name
					tree(prefix+treeAfterEndChar, st)
				case 0:
					out <- prefix + treeMidChar + st.Name
					tree(prefix+treeLinkChar, st)
				default:
					out <- prefix + treeMidChar + st.Name
					tree(prefix+treeLinkChar, st)
				}
			}
		}

		out <- t.Name
		tree("", t)

	}()
	return out
}
