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

type tagsLoadedMsg struct {
	tags       []ghost.Tag
	pagination *ghost.Pagination
}

type membersLoadedMsg struct {
	members    []ghost.Member
	pagination *ghost.Pagination
}

type errMsg struct {
	err error
}

type postToggledMsg struct {
	post ghost.Post
}

type postToggleErrMsg struct {
	err error
}
