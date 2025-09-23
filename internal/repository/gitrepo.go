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
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitRepo struct {
	Path string
	URL  string
	repo *git.Repository
}

func NewGitRepo(path, url string) *GitRepo {
	return &GitRepo{Path: path, URL: url}
}

func (g *GitRepo) EnsureCloned() error {
	if _, err := os.Stat(g.Path); os.IsNotExist(err) {
		fmt.Println("Cloning repo to", g.Path)
		r, err := git.PlainClone(g.Path, false, &git.CloneOptions{
			URL:      g.URL,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
		g.repo = r
		return nil
	}
	r, err := git.PlainOpen(g.Path)
	if err != nil {
		return err
	}
	g.repo = r
	return nil
}

func (g *GitRepo) DefaultBranch() (string, error) {
	rem, err := g.repo.Remote("origin")
	if err != nil {
		return "", err
	}
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}
	for _, r := range refs {
		if r.Type() == plumbing.SymbolicReference && r.Name().String() == "HEAD" {
			return strings.TrimPrefix(r.Target().String(), "refs/heads/"), nil
		}
	}

	return "main", nil
}

func (g *GitRepo) Pull() error {
	cmd := exec.Command("git", "-C", g.Path, "pull", "--rebase")
	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
	)
	return cmd.Run()
}

func (g *GitRepo) GetCommitHash(refName plumbing.ReferenceName) (string, error) {
	ref, err := g.repo.Reference(refName, true)
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}

func (g *GitRepo) GetFile(refName, filePath string) ([]byte, error) {

	commit, err := g.repo.ResolveRevision(plumbing.Revision(refName))
	if err != nil {
		return nil, err
	}

	// get tree
	commitObj, err := g.repo.CommitObject(*commit)
	if err != nil {
		return nil, err
	}
	tree, err := commitObj.Tree()
	if err != nil {
		return nil, err
	}

	// find file
	entry, err := tree.File(filePath)
	if err != nil {
		return nil, err
	}

	contents, err := entry.Contents()
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}
