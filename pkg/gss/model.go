package gss

import "io"

type GitAuth struct {
	Username string
	Password string
}

type CloneRepoOptions struct {
	*GitAuth
	RepoUrl string
	Output  io.Writer
}

type SearchHit map[string][]string

type SearchResult map[string]*SearchHit
