/*
 * Copyright (c) 2025 KAnggara75
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * See <https://www.gnu.org/licenses/gpl-3.0.html>.
 *
 * @author KAnggara75 on Mon 22/09/25 07.40
 * @project conflect repository
 * https://github.com/KAnggara75/conflect/tree/main/internal/repository
 */

package repository

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GitRepo struct {
	Path string
	URL  string
}

func NewGitRepo(path, url string) *GitRepo {
	return &GitRepo{Path: path, URL: url}
}

func (g *GitRepo) EnsureCloned() error {
	// cek apakah folder sudah ada
	if _, err := os.Stat(g.Path); os.IsNotExist(err) {
		log.Println("Cloning repo...")
		cmd := exec.Command("git", "clone", g.URL, g.Path)
		cmd.Env = append(os.Environ(),
			"GIT_TERMINAL_PROMPT=0",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git clone failed: %v, output: %s", err, string(out))
		}
		log.Println("Clone success")
		return nil
	}

	// kalau folder ada, pastikan itu repo git
	if _, err := os.Stat(filepath.Join(g.Path, ".git")); err == nil {
		log.Println("Repo already cloned")
		return nil
	}

	// folder ada tapi bukan repo git → hapus & clone ulang
	log.Println("Directory exists but not a git repo, removing and recloning...")
	if err := os.RemoveAll(g.Path); err != nil {
		return fmt.Errorf("failed to remove existing dir: %v", err)
	}

	cmd := exec.Command("git", "clone", g.URL, g.Path)
	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone retry failed: %v, output: %s", err, string(out))
	}
	log.Println("Clone success after cleanup")
	return nil
}

func (g *GitRepo) Pull() error {
	cmd := exec.Command("git", "-C", g.Path, "pull", "--rebase")
	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
	)
	return cmd.Run()
}

func (g *GitRepo) GetCommitHash() (string, error) {
	cmd := exec.Command("git", "-C", g.Path, "rev-parse", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git rev-parse failed: %v, output: %s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}
