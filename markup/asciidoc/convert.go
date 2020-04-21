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

// Package asciidoc converts Asciidoc to HTML using Asciidoc or Asciidoctor
// external binaries.
package asciidoc

import (
	"os/exec"
	"path/filepath"

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/markup/asciidoc/asciidoc_config"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/internal"
)

// Provider is the package entry point.
var Provider converter.ProviderProvider = provider{}

type provider struct {
}

func (p provider) New(cfg converter.ProviderConfig) (converter.Provider, error) {
	return converter.NewProvider("asciidoc", func(ctx converter.DocumentContext) (converter.Converter, error) {
		return &asciidocConverter{
			ctx: ctx,
			cfg: cfg,
		}, nil
	}), nil
}

type asciidocConverter struct {
	ctx converter.DocumentContext
	cfg converter.ProviderConfig
}

func (a *asciidocConverter) Convert(ctx converter.RenderContext) (converter.Result, error) {
	return converter.Bytes(a.getAsciidocContent(ctx.Src, a.ctx)), nil
}

func (c *asciidocConverter) Supports(feature identity.Identity) bool {
	return false
}

// getAsciidocContent calls asciidoctor or asciidoc as an external helper
// to convert AsciiDoc content to HTML.
func (a *asciidocConverter) getAsciidocContent(src []byte, ctx converter.DocumentContext) []byte {
	path := getAsciidoctorExecPath()
	if path == "" {
		path = getAsciidocExecPath()
		if path == "" {
			a.cfg.Logger.ERROR.Println("asciidoctor / asciidoc not found in $PATH: Please install.\n",
				"                 Leaving AsciiDoc content unrendered.")
			return src
		}
	}

	args := a.parseArgs()

	a.cfg.Logger.INFO.Println("Rendering", ctx.DocumentName, "with", path, "using asciidoc args", args, "...")

	return internal.ExternallyRenderContent(a.cfg, ctx, src, path, args)
}

func (a *asciidocConverter) parseArgs() []string {
	var cfg = a.cfg.MarkupConfig.AsciidocExt
	args := []string{}

	if contains(asciidoc_config.BackendWhitelist, cfg.Backend) {
		args = append(args, "-b", cfg.Backend)
	}

	for _, extension := range cfg.Extensions {
		if contains(asciidoc_config.ExtensionsWhitelist, extension) {
			args = append(args, "-r", extension)
			continue
		}
		a.cfg.Logger.ERROR.Println("Unsupported asciidoctor extension was passed in.")
	}

	if cfg.NoHeaderOrFooter == true {
		args = append(args, "--no-header-footer")
	}

	if cfg.SectionNumbers == true {
		args = append(args, "--section-numbers")
	}

	if cfg.Verbose == true {
		args = append(args, "-v")
	}

	if contains(asciidoc_config.SafeModeWhitelist, cfg.SafeMode) {
		args = append(args, "--safe-mode", cfg.SafeMode)
	}

	workingFolderCurrent := a.cfg.MarkupConfig.AsciidocExt.WorkingFolderCurrent

	if workingFolderCurrent {
		contentDir := filepath.Dir(a.ctx.Filename)
		destinationDir := a.cfg.Cfg.GetString("destination")

		a.cfg.Logger.INFO.Println("destinationDir", destinationDir)
		if destinationDir == "" {
			a.cfg.Logger.ERROR.Println("markup.asciidocext.workingFolderCurrent requires hugo command option --destination to be set")
		}

		outDir, err := filepath.Abs(filepath.Dir(filepath.Join(destinationDir, a.ctx.DocumentName)))
		if err != nil {
			a.cfg.Logger.ERROR.Println("asciidoctor outDir", err)
		}
		args = append(args, "--base-dir", contentDir, "-a", "outdir="+outDir)
	}

	args = append(args, "-")

	return args
}

func contains(s []string, elem interface{}) bool {
	for _, a := range s {
		if a == elem {
			return true
		}
	}

	return false
}

func (a *asciidocConverter) getAsciidoctorArgs(ctx converter.DocumentContext) []string {
	args := make([]string, 10)

	return args
}

func getAsciidocExecPath() string {
	path, err := exec.LookPath("asciidoc")
	if err != nil {
		return ""
	}
	return path
}

func getAsciidoctorExecPath() string {
	path, err := exec.LookPath("asciidoctor")
	if err != nil {
		return ""
	}
	return path
}

// Supports returns whether Asciidoc or Asciidoctor is installed on this computer.
func Supports() bool {
	return (getAsciidoctorExecPath() != "" ||
		getAsciidocExecPath() != "")
}
