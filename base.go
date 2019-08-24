package main

import (
	"time"
)

const (
	// SquareRepoID -
	SquareRepoID = "305104"
	// CoverRepoID -
	CoverRepoID = "302601"
	// GroupID -
	GroupID = "334719"
	// FirstShelfRepoID -
	FirstShelfRepoID = "302601"
	// AboutusRepoID -
	AboutusRepoID = "305107"
	// Bookshelf -
	Bookshelf = "书架"
	// Column -
	Column = "专栏"
)

var (
	// Timer -
	Timer = time.Duration(7200000000000)

	// ContentURL -
	ContentURL = "https://yuque.yangchengkai.now.sh/api/v1/yuque/details?RepoID=%s&ID=%s"
	// ListURL -
	ListURL = "https://yuque.yangchengkai.now.sh/api/v1/yuque/list?RepoID=%s"
	// RepoURL -
	RepoURL = "https://yuque.yangchengkai.now.sh/api/v1/yuque/repo?GroupID=%s"
)
