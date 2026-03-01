package tui

import "github.com/manzanita-research/caspar/pkg/ghost"

type siteLoadedMsg struct {
	site        *ghost.SiteInfo
	postCount   int
	pageCount   int
	memberCount int
	tagCount    int
}

type postsLoadedMsg struct {
	posts      []ghost.Post
	pagination *ghost.Pagination
}

type postDetailMsg struct {
	post ghost.Post
}

type pagesLoadedMsg struct {
	pages      []ghost.Page
	pagination *ghost.Pagination
}

type errMsg struct {
	err error
}
