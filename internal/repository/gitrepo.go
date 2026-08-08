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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

type GitRepo struct {
	Path string
	URL  string
}

func NewGitRepo(path, url string) *GitRepo {
	return &GitRepo{Path: path, URL: url}
}

func (g *GitRepo) InitAllBranches() error {
	if strings.TrimSpace(g.URL) == "" {
		return errors.New("repository URL is empty (please set REPO_URL environment variable or REPO_URL_FILE)")
	}

	log.Printf("🔍 Fetching remote branch list from URL: %s", g.URL)
	branches, err := g.listRemoteBranches()
	if err != nil {
		return err
	}

	log.Printf("📋 Found %d remote branch(es): %v", len(branches), branches)

	for _, branch := range branches {
		if _, err := g.EnsureBranch(branch); err != nil {
			return err
		}
	}

	return nil
}

func (g *GitRepo) listRemoteBranches() ([]string, error) {
	if strings.TrimSpace(g.URL) == "" {
		return nil, errors.New("repository URL is empty (please set REPO_URL environment variable or REPO_URL_FILE)")
	}

	// Create a remote to list references without cloning
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{g.URL},
	})

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list remote branches from %s: %w", g.URL, err)
	}

	var branches []string
	for _, ref := range refs {
		refName := ref.Name().String()
		if strings.HasPrefix(refName, "refs/heads/") {
			branch := strings.TrimPrefix(refName, "refs/heads/")
			branches = append(branches, branch)
		}
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
		log.Printf("📦 Cloning branch %q from URL %s into %s...", branch, g.URL, targetPath)

		cloneOpts := &git.CloneOptions{
			URL:          g.URL,
			SingleBranch: true,
			Depth:        1,
		}

		if branch != "" {
			cloneOpts.ReferenceName = plumbing.NewBranchReferenceName(branch)
		}

		_, err := git.PlainClone(targetPath, false, cloneOpts)
		if err != nil {
			log.Printf("❌ Failed to clone branch %q into %s: %v", branch, targetPath, err)
			return "", fmt.Errorf("failed to clone branch %s: %w", branch, err)
		}
		log.Printf("✅ Successfully cloned branch %q into %s", branch, targetPath)
	} else {
		log.Printf("ℹ️ Branch %q directory already exists at %s", branch, targetPath)
	}
	return targetPath, nil
}

func (g *GitRepo) Pull(branch string) error {
	branchPath := filepath.Join(g.Path, branch)

	repo, err := git.PlainOpen(branchPath)
	if err != nil {
		return fmt.Errorf("failed to open repo at %s: %w", branchPath, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Force:         true,
	})

	// Already up-to-date is not an error
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("failed to pull branch %s: %w", branch, err)
	}

	return nil
}

func (g *GitRepo) GetCommitHashFromBranch(branch string) (string, error) {
	branchPath := filepath.Join(g.Path, branch)

	repo, err := git.PlainOpen(branchPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repo at %s: %w", branchPath, err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD for branch %s: %w", branch, err)
	}

	return head.Hash().String(), nil
}
