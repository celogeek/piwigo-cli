package piwigotools

import (
	"path/filepath"
	"sort"
	"strings"
)

type Tree interface {
	Add(string) Tree
	AddPath(string) Tree
	FlatView() chan string
	TreeView() chan string
}

type node struct {
	Name     string
	children map[string]*node
}

func NewTree() Tree {
	return &node{
		Name: ".",
	}
}

func (t *node) Add(name string) Tree {
	if t.children == nil {
		t.children = map[string]*node{}
	}
	n, ok := t.children[name]
	if !ok {
		n = &node{Name: name}
		t.children[name] = n
	}
	return n
}
func (t *node) AddPath(path string) Tree {
	n := Tree(t)
	for _, name := range strings.Split(path, "/") {
		n = n.Add(name)
	}
	return n
}

func (t *node) Children() []*node {
	childs := make([]*node, len(t.children))
	i := 0
	for _, n := range t.children {
		childs[i] = n
		i++
	}
	sort.Slice(childs, func(i, j int) bool {
		return childs[i].Name < childs[j].Name
	})
	return childs
}

func (t *node) HasChildren() bool {
	return t.children != nil
}

func (t *node) FlatView() (out chan string) {
	out = make(chan string)
	go func() {
		defer close(out)
		var flatten func(string, *node)

		flatten = func(path string, t *node) {
			switch t.HasChildren() {
			case false:
				out <- path
			case true:
				for _, child := range t.Children() {
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
			children := t.Children()
			for i, st := range children {
				switch i {
				case len(children) - 1:
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
