// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package asciidoc_config holds asciidoc related configuration.
package asciidoc_config

// DefaultConfig holds the default asciidoc configuration.
var (
	Default = Config{
		Backend:              "html5",
		DocType:              "article",
		Extensions:           []string{""},
		NoHeaderOrFooter:     true,
		SafeMode:             "safe",
		SectionNumbers:       false,
		Verbose:              false,
		WorkingFolderCurrent: false,
	}

	ExtensionsWhitelist = []string{
		"asciidoctor-html5s",
		"asciidoctor-diagram",
		"asciidoctor-interdoc-reftext",
		"asciidoctor-katex",
		"asciidoctor-latex",
		"asciidoctor-question",
		"asciidoctor-rouge",
	}

	SafeModeWhitelist = []string{
		"unsafe", //default
		"safe",
		"server",
		"secure",
	}

	BackendWhitelist = []string{
		"html5",
		"html5s",
		"xhtml5",
		"docbook5",
		"docbook45",
		"manpage",
	}
)

// Config configures asciidoc.
type Config struct {
	Backend              string
	DocType              string
	Extensions           []string
	NoHeaderOrFooter     bool
	SafeMode             string
	SectionNumbers       bool
	Verbose              bool
	WorkingFolderCurrent bool
}
