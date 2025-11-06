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
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

type GitRepo struct {
	Path string
	URL  string
	Repo *git.Repository
}

func NewGitRepo(path, url string) *GitRepo {
	return &GitRepo{Path: path, URL: url}
}

func (g *GitRepo) InitAllBranches() error {
	branches, err := g.listRemoteBranches()
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if _, err := g.EnsureBranch(branch); err != nil {
			return err
		}
	}

	return nil
}

func (g *GitRepo) listRemoteBranches() ([]string, error) {
	fmt.Printf("Listing remote branches %s", g.URL)

	cmd := exec.Command("git", "-c", "credential.helper=", "ls-remote", "-v", "--heads", g.URL)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list remote branches: %w", err)
	}

	var branches []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) != 2 {
			continue
		}
		ref := parts[1]
		if strings.HasPrefix(ref, "refs/heads/") {
			branch := strings.TrimPrefix(ref, "refs/heads/")
			branches = append(branches, branch)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return branches, nil
}

func (g *GitRepo) EnsureBranch(branch string) (string, error) {
	path := "origin"
	if branch != "" {
		path = branch
	}

	targetPath := filepath.Join(g.Path, path)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		var cmd *exec.Cmd

		if branch != "" {
			fmt.Printf("Cloning branch %s into %s\n", branch, targetPath)
			cmd = exec.Command(
				"git", "clone",
				"--branch", branch,
				"--single-branch",
				"--depth=1",
				g.URL,
				targetPath,
			)
		} else {
			fmt.Printf("Cloning default branch into %s\n", targetPath)
			cmd = exec.Command(
				"git", "clone",
				"--single-branch",
				"--depth=1",
				g.URL,
				targetPath,
			)

		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to clone branch %s: %w", branch, err)
		}
	}
	return targetPath, nil
}

func (g *GitRepo) Pull(branch string) error {
	branchPath := filepath.Join(g.Path, branch)

	cmd := exec.Command("git", "-C", branchPath, "pull", "--rebase")
	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
	)
	return cmd.Run()
}

func (g *GitRepo) GetCommitHashFromBranch(branch string) (string, error) {
	branchPath := filepath.Join(g.Path, branch)
	cmd := exec.Command("git", "rev-parse", branch)

	cmd.Dir = branchPath

	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git rev-parse error for branch %s: %v (%s)", branchPath, err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}
