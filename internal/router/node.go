package router

import (
	"strings"
)

// node represents a single segment in the routing tree.
type node struct {
	path     string  // the path segment this node represents (e.g., "users")
	children []*node // child nodes
	handler  Handler // non-nil if this node is a terminal route
	isParam  bool    // true if this segment is a parameter (starts with ':')
	paramKey string  // the name of the parameter (e.g., "id" for ":id")
}

// insert adds a route to the tree.
func (n *node) insert(pattern string, handler Handler) {
	segments := splitPath(pattern)
	curr := n

	for _, segment := range segments {
		child := curr.matchChild(segment)
		if child == nil {
			isParam := strings.HasPrefix(segment, ":")
			paramKey := ""
			if isParam {
				paramKey = segment[1:]
			}

			child = &node{
				path:     segment,
				isParam:  isParam,
				paramKey: paramKey,
			}
			curr.children = append(curr.children, child)
		}
		curr = child
	}

	curr.handler = handler
}

// search finds a handler for the given path, capturing any parameters along the way.
func (n *node) search(path string) (Handler, map[string]string) {
	segments := splitPath(path)
	params := make(map[string]string)

	curr := n
	for _, segment := range segments {
		// First try to find an exact match
		child := curr.matchExactChild(segment)
		if child == nil {
			// If no exact match, look for a parameter match
			child = curr.matchParamChild()
			if child == nil {
				// No match found
				return nil, nil
			}
			// Capture the parameter value
			params[child.paramKey] = segment
		}
		curr = child
	}

	return curr.handler, params
}

// matchExactChild looks for a child node with the exact path segment.
func (n *node) matchExactChild(segment string) *node {
	for _, child := range n.children {
		if !child.isParam && child.path == segment {
			return child
		}
	}
	return nil
}

// matchParamChild looks for a child node that acts as a parameter.
func (n *node) matchParamChild() *node {
	for _, child := range n.children {
		if child.isParam {
			return child
		}
	}
	return nil
}

// matchChild is used during insertion to find an existing segment (exact or param).
func (n *node) matchChild(segment string) *node {
	for _, child := range n.children {
		if child.path == segment {
			return child
		}
	}
	return nil
}

// splitPath splits a URL path into segments, ignoring empty segments (e.g., consecutive slashes).
func splitPath(path string) []string {
	parts := strings.Split(path, "/")
	var segments []string
	for _, part := range parts {
		if part != "" {
			segments = append(segments, part)
		}
	}
	// If the original path was exactly "/", we need an explicit root representation,
	// but in a segment-based trie, the root node handles it inherently when segments is empty.
	return segments
}
