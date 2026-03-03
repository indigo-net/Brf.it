package formatter

import (
	"fmt"
	"testing"
)

func TestIsStdLibImport(t *testing.T) {
	tests := []struct {
		language   string
		importPath string
		want       bool
	}{
		// Go - stdlib
		{"go", `import "fmt"`, true},
		{"go", `import "os"`, true},
		{"go", `import "os/exec"`, true},
		{"go", `import "net/http"`, true},
		{"go", `import "strings"`, true},
		{"go", `import "testing"`, true},
		{"go", `import _ "embed"`, true},
		// Go - non-stdlib
		{"go", `import "github.com/foo/bar"`, false},
		{"go", `import "github.com/indigo-net/Brf.it/pkg/parser"`, false},
		{"go", `import pkgcontext "github.com/indigo-net/Brf.it/internal/context"`, false},
		{"go", `import alias "github.com/foo/bar"`, false},

		// Python - stdlib
		{"python", "import os", true},
		{"python", "import sys", true},
		{"python", "from os import path", true},
		{"python", "from collections import OrderedDict", true},
		{"python", "import json", true},
		{"python", "from pathlib import Path", true},
		{"python", "from zoneinfo import ZoneInfo", true},
		{"python", "import graphlib", true},
		// Python - non-stdlib
		{"python", "import requests", false},
		{"python", "from flask import Flask", false},
		{"python", "import numpy", false},

		// JavaScript - stdlib
		{"javascript", `import "fs"`, true},
		{"javascript", `import "path"`, true},
		{"javascript", `import "node:fs"`, true},
		{"javascript", `import { readFile } from 'fs'`, true},
		{"javascript", `import http from "node:http"`, true},
		// JavaScript - non-stdlib
		{"javascript", `import "react"`, false},
		{"javascript", `import { useState } from 'react'`, false},
		{"javascript", `import "./utils"`, false},

		// TypeScript - stdlib
		{"typescript", `import "fs"`, true},
		{"typescript", `import "node:path"`, true},
		// TypeScript - non-stdlib
		{"typescript", `import "express"`, false},
		{"typescript", `import "./types"`, false},

		// C/C++ - stdlib (system headers)
		{"c", `#include <stdio.h>`, true},
		{"c", `#include <stdlib.h>`, true},
		{"cpp", `#include <iostream>`, true},
		{"cpp", `#include <vector>`, true},
		// C/C++ - non-stdlib (user headers)
		{"c", `#include "myheader.h"`, false},
		{"cpp", `#include "myheader.h"`, false},

		// Java - stdlib
		{"java", "import java.util.List", true},
		{"java", "import java.io.File", true},
		{"java", "import javax.swing.JFrame", true},
		{"java", "import static java.lang.Math.PI", true},
		// Java - non-stdlib
		{"java", "import com.google.gson.Gson", false},
		{"java", "import org.apache.commons.lang3.StringUtils", false},

		// Rust - stdlib
		{"rust", "use std::collections::HashMap", true},
		{"rust", "use std::io::{self, Read, Write}", true},
		{"rust", "use core::fmt", true},
		{"rust", "use alloc::vec::Vec", true},
		// Rust - non-stdlib
		{"rust", "use serde::Deserialize", false},
		{"rust", "use crate::utils::*", false},
		{"rust", "use tokio::runtime::Runtime", false},

		// Unknown language - always false
		{"ruby", `require "json"`, false},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%s/%s", tt.language, tt.importPath)
		t.Run(name, func(t *testing.T) {
			got := isStdLibImport(tt.language, tt.importPath)
			if got != tt.want {
				t.Errorf("isStdLibImport(%q, %q) = %v, want %v", tt.language, tt.importPath, got, tt.want)
			}
		})
	}
}

func TestGetEmptyComment(t *testing.T) {
	tests := []struct {
		lang     string
		expected string
	}{
		{"go", "// (empty)"},
		{"typescript", "// (empty)"},
		{"javascript", "// (empty)"},
		{"java", "// (empty)"},
		{"c", "// (empty)"},
		{"python", "# (empty)"},
		{"ruby", "# (empty)"},
		{"html", "<!-- (empty) -->"},
		{"xml", "<!-- (empty) -->"},
		{"unknown", "// (empty)"},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			result := getEmptyComment(tt.lang)
			if result != tt.expected {
				t.Errorf("getEmptyComment(%q) = %q, want %q", tt.lang, result, tt.expected)
			}
		})
	}
}
