package context

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// treeNode represents a node in the directory tree.
type treeNode struct {
	children map[string]*treeNode
	tokens   int  // token count for leaf nodes (files)
	isFile   bool // true if this node is a file (leaf)
}

// BuildTree generates a directory tree string from file paths.
// The root parameter is used to calculate relative paths.
func BuildTree(root string, paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	rootNode := &treeNode{}

	// Insert all paths into the tree
	for _, path := range paths {
		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}

		parts := strings.Split(rel, string(filepath.Separator))
		current := rootNode

		for _, part := range parts {
			if current.children == nil {
				current.children = make(map[string]*treeNode)
			}
			if _, exists := current.children[part]; !exists {
				current.children[part] = &treeNode{}
			}
			current = current.children[part]
		}
	}

	// Build the output string (root level has no prefix)
	var buf strings.Builder
	renderNode(&buf, rootNode, "", true)
	return strings.TrimSuffix(buf.String(), "\n")
}

// maxTreeDepth is the maximum recursion depth for rendering tree nodes.
// Prevents stack overflow on pathologically deep directory structures.
const maxTreeDepth = 100

// renderNode recursively renders the tree structure.
// isRoot indicates whether this is the root level (no connectors).
func renderNode(buf *strings.Builder, n *treeNode, prefix string, isRoot bool) {
	renderNodeDepth(buf, n, prefix, isRoot, 0)
}

func renderNodeDepth(buf *strings.Builder, n *treeNode, prefix string, isRoot bool, depth int) {
	if depth > maxTreeDepth {
		buf.WriteString(prefix + "... (truncated: depth > " + fmt.Sprintf("%d", maxTreeDepth) + ")\n")
		return
	}
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
		renderNodeDepth(buf, child, newPrefix, false, depth+1)
	}
}

// FileTokenCount holds a file path and its token count.
type FileTokenCount struct {
	Path   string
	Tokens int
}

// BuildTokenTree generates a directory tree string with per-file token counts.
// Each file node shows its token count, and directory nodes show the sum of their children.
func BuildTokenTree(root string, files []FileTokenCount) string {
	if len(files) == 0 {
		return ""
	}

	rootNode := &treeNode{}

	// Insert all paths into the tree with token counts
	for _, f := range files {
		rel, err := filepath.Rel(root, f.Path)
		if err != nil {
			rel = f.Path
		}

		parts := strings.Split(rel, string(filepath.Separator))
		current := rootNode

		for i, part := range parts {
			if current.children == nil {
				current.children = make(map[string]*treeNode)
			}
			if _, exists := current.children[part]; !exists {
				current.children[part] = &treeNode{}
			}
			current = current.children[part]

			// Mark leaf node with token count
			if i == len(parts)-1 {
				current.tokens = f.Tokens
				current.isFile = true
			}
		}
	}

	// Calculate directory totals bottom-up
	calcDirTokens(rootNode)

	// Build the output string
	var buf strings.Builder
	renderTokenNode(&buf, rootNode, "", true)
	buf.WriteString(fmt.Sprintf("\nTotal: %s tokens\n", formatNumber(rootNode.tokens)))
	return strings.TrimSuffix(buf.String(), "\n")
}

// calcDirTokens recursively sums token counts for directory nodes.
func calcDirTokens(n *treeNode) int {
	if n.isFile {
		return n.tokens
	}
	var sum int
	for _, child := range n.children {
		sum += calcDirTokens(child)
	}
	n.tokens = sum
	return sum
}

// renderTokenNode recursively renders the tree structure with token counts.
func renderTokenNode(buf *strings.Builder, n *treeNode, prefix string, isRoot bool) {
	keys := make([]string, 0, len(n.children))
	for k := range n.children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, key := range keys {
		child := n.children[key]
		isLast := i == len(keys)-1

		label := fmt.Sprintf("%s (%s tokens)", key, formatNumber(child.tokens))

		if isRoot {
			buf.WriteString(label + "\n")
		} else {
			connector := "├── "
			if isLast {
				connector = "└── "
			}
			buf.WriteString(prefix + connector + label + "\n")
		}

		var newPrefix string
		if isRoot {
			newPrefix = ""
		} else if isLast {
			newPrefix = prefix + "    "
		} else {
			newPrefix = prefix + "│   "
		}
		renderTokenNode(buf, child, newPrefix, false)
	}
}

// formatNumber formats an integer with comma separators (e.g., 1234 → "1,234").
func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	var buf strings.Builder
	remainder := len(s) % 3
	if remainder > 0 {
		buf.WriteString(s[:remainder])
	}
	for i := remainder; i < len(s); i += 3 {
		if buf.Len() > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(s[i : i+3])
	}
	return buf.String()
}
