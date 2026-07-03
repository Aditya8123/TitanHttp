package router

import (
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

// Static registers a GET route that serves static files from the provided root directory.
// The prefix is the URL path prefix (e.g., "/static/"), and root is the local directory path.
func (r *Router) Static(prefix, root string) {
	// Ensure prefix ends with a slash so the wildcard works predictably
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	// Register a wildcard route (e.g., "/static/*filepath")
	pattern := prefix + "*filepath"

	r.Get(pattern, func(req *http.Request) *http.Response {
		// Extract the requested file path from the route parameters
		fp, ok := req.Params["filepath"]
		if !ok || fp == "" {
			return http.NewResponse404()
		}

		// Prevent path traversal vulnerabilities (e.g., navigating outside the root directory)
		if strings.Contains(fp, "..") {
			return http.NewResponse400()
		}

		// Construct the absolute or relative path to the file on the filesystem
		fullPath := filepath.Join(root, fp)

		// Attempt to open the file
		file, err := os.Open(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				return http.NewResponse404()
			}
			// Other errors (e.g., permission denied)
			return http.NewResponse500()
		}
		// file is intentionally NOT closed here. It is assigned to resp.Stream and will be closed by resp.WriteTo()

		// Get file information to check if it's a directory
		info, err := file.Stat()
		if err != nil {
			file.Close()
			return http.NewResponse500()
		}

		// Handle directories
		if info.IsDir() {
			indexPath := filepath.Join(fullPath, "index.html")
			indexFile, err := os.Open(indexPath)
			if err != nil {
				file.Close()
				// If index.html doesn't exist or can't be opened, return 403 Forbidden
				return http.NewResponse403()
			}
			
			indexInfo, err := indexFile.Stat()
			if err != nil || indexInfo.IsDir() {
				file.Close()
				indexFile.Close()
				return http.NewResponse403()
			}

			// Replace file with indexFile
			file.Close()
			file = indexFile
			fullPath = indexPath
			info = indexInfo
		}

		// Detect MIME type
		contentType := mime.TypeByExtension(filepath.Ext(fullPath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// Construct a successful HTTP response
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Headers["Content-Type"] = contentType
		resp.Headers["Cache-Control"] = "public, max-age=3600"
		resp.Headers["Content-Length"] = strconv.FormatInt(info.Size(), 10)
		resp.Stream = file

		return resp
	})
}
