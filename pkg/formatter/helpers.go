package formatter

import "strings"

// isStdLibImport returns true if the import path is a standard library import
// for the given language. This is used to filter out stdlib imports when
// the --no-std-imports flag is set.
func isStdLibImport(language, importPath string) bool {
	switch language {
	case "go":
		return isGoStdLib(importPath)
	case "python":
		return isPythonStdLib(importPath)
	case "javascript", "typescript":
		return isJSStdLib(importPath)
	case "c", "cpp":
		return isCStdLib(importPath)
	case "java":
		return isJavaStdLib(importPath)
	case "rust":
		return isRustStdLib(importPath)
	default:
		return false
	}
}

// isGoStdLib checks if a Go import is from the standard library.
// Go stdlib imports have no dots in the path (e.g., "fmt", "os/exec").
// Handles formats: `import "fmt"`, `import _ "embed"`, `import alias "github.com/foo/bar"`
func isGoStdLib(importPath string) bool {
	// Extract the quoted path between the first and last double quotes
	first := strings.Index(importPath, `"`)
	last := strings.LastIndex(importPath, `"`)
	if first < 0 || first == last {
		return false
	}
	path := importPath[first+1 : last]
	// Go stdlib paths have no dots (no domain like github.com)
	return path != "" && !strings.Contains(path, ".")
}

// isPythonStdLib checks if a Python import is from the standard library.
// Format: `import os`, `from sys import argv`
func isPythonStdLib(importPath string) bool {
	path := importPath
	// Extract module name from "from X import Y" or "import X"
	if strings.HasPrefix(path, "from ") {
		path = strings.TrimPrefix(path, "from ")
		if idx := strings.Index(path, " import"); idx > 0 {
			path = path[:idx]
		}
	} else if strings.HasPrefix(path, "import ") {
		path = strings.TrimPrefix(path, "import ")
	}
	path = strings.TrimSpace(path)
	// Get top-level module name
	if idx := strings.Index(path, "."); idx > 0 {
		path = path[:idx]
	}
	_, ok := pythonStdLibModules[path]
	return ok
}

// isJSStdLib checks if a JS/TS import is a Node.js builtin module.
// Format: `import "node:fs"`, `import "path"`, or full import statements.
// Note: Browser-targeted projects may see false positives for names
// that collide with Node.js builtins (e.g., "path", "events").
func isJSStdLib(importPath string) bool {
	path := importPath
	// Extract module path from full import statements
	// e.g., `import { readFile } from 'fs'` or `import fs from "node:fs"`
	if strings.Contains(path, " from ") {
		parts := strings.Split(path, " from ")
		path = parts[len(parts)-1]
	} else if strings.HasPrefix(path, "import ") {
		path = strings.TrimPrefix(path, "import ")
	} else if strings.HasPrefix(path, "require(") {
		path = strings.TrimPrefix(path, "require(")
		path = strings.TrimSuffix(path, ")")
	}
	path = strings.Trim(strings.TrimSpace(path), `"';`)
	// node: prefix is always builtin
	if strings.HasPrefix(path, "node:") {
		return true
	}
	_, ok := nodeBuiltinModules[path]
	return ok
}

// isCStdLib checks if a C/C++ include is a system header.
// System headers use angle brackets: `#include <stdio.h>`
// User headers use quotes: `#include "myheader.h"`
func isCStdLib(importPath string) bool {
	path := strings.TrimSpace(importPath)
	if !strings.HasPrefix(path, "#include") {
		return false
	}
	path = strings.TrimPrefix(path, "#include")
	path = strings.TrimSpace(path)
	return strings.HasPrefix(path, "<") && strings.HasSuffix(path, ">")
}

// isJavaStdLib checks if a Java import is from the standard library.
// Java stdlib: java.* and javax.* packages
// Format: `import java.util.List`
func isJavaStdLib(importPath string) bool {
	path := strings.TrimPrefix(importPath, "import ")
	path = strings.TrimPrefix(path, "static ")
	path = strings.TrimSpace(path)
	path = strings.TrimSuffix(path, ";")
	return strings.HasPrefix(path, "java.") || strings.HasPrefix(path, "javax.")
}

// isRustStdLib checks if a Rust use statement is from the standard library.
// Rust stdlib: std::*, core::*, alloc::*
// Format: `use std::io`
func isRustStdLib(importPath string) bool {
	path := strings.TrimPrefix(importPath, "use ")
	path = strings.TrimSpace(path)
	path = strings.TrimSuffix(path, ";")
	return strings.HasPrefix(path, "std::") ||
		strings.HasPrefix(path, "core::") ||
		strings.HasPrefix(path, "alloc::")
}

// pythonStdLibModules is a set of Python standard library top-level modules.
var pythonStdLibModules = map[string]struct{}{
	"abc": {}, "aifc": {}, "argparse": {}, "array": {}, "ast": {},
	"asynchat": {}, "asyncio": {}, "asyncore": {}, "atexit": {},
	"base64": {}, "bdb": {}, "binascii": {}, "binhex": {}, "bisect": {},
	"builtins": {}, "bz2": {},
	"calendar": {}, "cgi": {}, "cgitb": {}, "chunk": {}, "cmath": {},
	"cmd": {}, "code": {}, "codecs": {}, "codeop": {}, "collections": {},
	"colorsys": {}, "compileall": {}, "concurrent": {}, "configparser": {},
	"contextlib": {}, "contextvars": {}, "copy": {}, "copyreg": {},
	"cProfile": {}, "crypt": {}, "csv": {}, "ctypes": {}, "curses": {},
	"dataclasses": {}, "datetime": {}, "dbm": {}, "decimal": {}, "difflib": {},
	"dis": {}, "distutils": {},
	"email": {}, "encodings": {}, "enum": {}, "errno": {},
	"faulthandler": {}, "fcntl": {}, "filecmp": {}, "fileinput": {},
	"fnmatch": {}, "fractions": {}, "ftplib": {}, "functools": {},
	"gc": {}, "getopt": {}, "getpass": {}, "gettext": {}, "glob": {},
	"graphlib": {}, "grp": {}, "gzip": {},
	"hashlib": {}, "heapq": {}, "hmac": {}, "html": {}, "http": {},
	"idlelib": {}, "imaplib": {}, "imghdr": {}, "imp": {}, "importlib": {},
	"inspect": {}, "io": {}, "ipaddress": {}, "itertools": {},
	"json": {},
	"keyword": {},
	"lib2to3": {}, "linecache": {}, "locale": {}, "logging": {}, "lzma": {},
	"mailbox": {}, "mailcap": {}, "marshal": {}, "math": {}, "mimetypes": {},
	"mmap": {}, "modulefinder": {}, "multiprocessing": {},
	"netrc": {}, "nis": {}, "nntplib": {}, "numbers": {},
	"operator": {}, "optparse": {}, "os": {}, "ossaudiodev": {},
	"parser": {}, "pathlib": {}, "pdb": {}, "pickle": {}, "pickletools": {},
	"pipes": {}, "pkgutil": {}, "platform": {}, "plistlib": {}, "poplib": {},
	"posix": {}, "posixpath": {}, "pprint": {}, "profile": {}, "pstats": {},
	"pty": {}, "pwd": {}, "py_compile": {}, "pyclbr": {}, "pydoc": {},
	"queue": {}, "quopri": {},
	"random": {}, "re": {}, "readline": {}, "reprlib": {}, "resource": {},
	"rlcompleter": {}, "runpy": {},
	"sched": {}, "secrets": {}, "select": {}, "selectors": {}, "shelve": {},
	"shlex": {}, "shutil": {}, "signal": {}, "site": {}, "smtpd": {},
	"smtplib": {}, "sndhdr": {}, "socket": {}, "socketserver": {},
	"sqlite3": {}, "ssl": {}, "stat": {}, "statistics": {}, "string": {},
	"stringprep": {}, "struct": {}, "subprocess": {}, "sunau": {},
	"symtable": {}, "sys": {}, "sysconfig": {}, "syslog": {},
	"tabnanny": {}, "tarfile": {}, "telnetlib": {}, "tempfile": {},
	"termios": {}, "test": {}, "textwrap": {}, "threading": {}, "time": {},
	"timeit": {}, "tkinter": {}, "token": {}, "tokenize": {}, "tomllib": {},
	"trace": {}, "traceback": {}, "tracemalloc": {}, "tty": {}, "turtle": {},
	"turtledemo": {}, "types": {}, "typing": {},
	"unicodedata": {}, "unittest": {}, "urllib": {}, "uu": {}, "uuid": {},
	"venv": {},
	"warnings": {}, "wave": {}, "weakref": {}, "webbrowser": {},
	"winreg": {}, "winsound": {}, "wsgiref": {},
	"xdrlib": {}, "xml": {}, "xmlrpc": {},
	"zipapp": {}, "zipfile": {}, "zipimport": {}, "zlib": {},
	"zoneinfo": {},
	"_thread": {},
}

// nodeBuiltinModules is a set of Node.js builtin module names.
var nodeBuiltinModules = map[string]struct{}{
	"assert": {}, "buffer": {}, "child_process": {}, "cluster": {},
	"console": {}, "constants": {}, "crypto": {}, "dgram": {},
	"diagnostics_channel": {}, "dns": {}, "domain": {}, "events": {},
	"fs": {}, "http": {}, "http2": {}, "https": {}, "inspector": {},
	"module": {}, "net": {}, "os": {}, "path": {}, "perf_hooks": {},
	"process": {}, "punycode": {}, "querystring": {}, "readline": {},
	"repl": {}, "stream": {}, "string_decoder": {}, "timers": {},
	"tls": {}, "trace_events": {}, "tty": {}, "url": {}, "util": {},
	"v8": {}, "vm": {}, "wasi": {}, "worker_threads": {}, "zlib": {},
}

// getEmptyComment returns the appropriate empty file comment for a language.
func getEmptyComment(lang string) string {
	switch lang {
	case "python", "ruby":
		return "# (empty)"
	case "html", "xml":
		return "<!-- (empty) -->"
	case "go", "c", "cpp", "java", "javascript", "typescript":
		return "// (empty)"
	default:
		return "// (empty)"
	}
}
