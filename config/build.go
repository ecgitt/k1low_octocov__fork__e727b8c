package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/k1LoW/octocov/internal"
)

func (c *Config) Build() {
	// Repository
	if c.Repository == "" {
		c.Repository = os.Getenv("GITHUB_REPOSITORY")
	}

	// Coverage
	if c.Coverage == nil {
		c.Coverage = &Coverage{}
	}
	if c.Coverage.Path != "" {
		_, _ = fmt.Fprintln(os.Stderr, "Deprecated: coverage.path: has been deprecated. please use coverage.paths: instead.") //nostyle:handlerrors
		c.Coverage.Paths = append(c.Coverage.Paths, c.Coverage.Path)
	}
	if len(c.Coverage.Paths) == 0 {
		c.Coverage.Paths = append(c.Coverage.Paths, filepath.Dir(c.path))
	} else {
		var paths []string
		for _, p := range c.Coverage.Paths {
			p = filepath.FromSlash(p)
			globPath := p
			if !filepath.IsAbs(p) {
				globPath = filepath.Join(filepath.Dir(c.path), p)
			}
			matches, err := filepath.Glob(globPath)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "warn: %s\n", err) //nostyle:handlerrors
				continue
			}
			if len(matches) > 0 {
				paths = append(paths, matches...)
			} else if !strings.ContainsAny(p, "*?[") {
				paths = append(paths, globPath)
			}
		}
		c.Coverage.Paths = paths
	}

	// TestExecutionTime
	if c.TestExecutionTime == nil {
		c.TestExecutionTime = &TestExecutionTime{}
	}

	// Report

	// Central
	if c.Central != nil {
		if c.Central.Root == "" {
			c.Central.Root = "."
		}
		if !filepath.IsAbs(c.Central.Root) {
			c.Central.Root = filepath.Clean(filepath.Join(c.Root(), c.Central.Root))
		}
		if len(c.Central.Reports.Datastores) == 0 {
			c.Central.Reports.Datastores = append(c.Central.Reports.Datastores, defaultReportsDatastore)
		}
		if len(c.Central.Badges.Datastores) == 0 {
			c.Central.Badges.Datastores = append(c.Central.Badges.Datastores, defaultBadgesDatastore)
		}
	}

	// Push

	// Comment

	// Diff

	// GitRoot
	gitRoot, _ := internal.RootPath(c.Root()) //nostyle:handlerrors
	c.GitRoot = gitRoot
}
