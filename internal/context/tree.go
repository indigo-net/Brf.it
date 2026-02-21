package context

import (
	"path/filepath"
	"sort"
	"strings"
)

// treeNode represents a node in the directory tree.
type treeNode struct {
	children map[string]*treeNode
}

// BuildTree generates a directory tree string from file paths.
// The root parameter is used to calculate relative paths.
func BuildTree(root string, paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	rootNode := &treeNode{children: make(map[string]*treeNode)}

	// Insert all paths into the tree
	for _, path := range paths {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}

		parts := strings.Split(rel, string(filepath.Separator))
		current := rootNode

		for _, part := range parts {
			if _, exists := current.children[part]; !exists {
				current.children[part] = &treeNode{
					children: make(map[string]*treeNode),
				}
			}
			current = current.children[part]
		}
	}

	// Build the output string (root level has no prefix)
	var buf strings.Builder
	renderNode(&buf, rootNode, "", true)
	return strings.TrimSuffix(buf.String(), "\n")
}

// renderNode recursively renders the tree structure.
// isRoot indicates whether this is the root level (no connectors).
func renderNode(buf *strings.Builder, n *treeNode, prefix string, isRoot bool) {
	// Get sorted keys for consistent output
	keys := make([]string, 0, len(n.children))
	for k := range n.children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, key := range keys {
		child := n.children[key]
		isLast := i == len(keys)-1

		if isRoot {
			// Root level: no prefix, no connector
			buf.WriteString(key + "\n")
		} else {
			// Child levels: with prefix and connector
			connector := "├── "
			if isLast {
				connector = "└── "
			}
			buf.WriteString(prefix + connector + key + "\n")
		}

		// Render children with updated prefix
		var newPrefix string
		if isRoot {
			newPrefix = ""
		} else if isLast {
			newPrefix = prefix + "    "
		} else {
			newPrefix = prefix + "│   "
		}
		renderNode(buf, child, newPrefix, false)
	}
}
