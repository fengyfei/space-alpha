package main

import (
	"errors"
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
	// Bookshelf -
	Bookshelf = "书架"
	// Column -
	Column = "专栏"
)

var (
	errGetFirst = errors.New("Getting first shelf is wrong")
	errGetList  = errors.New("Getting shelf list is wrong")

	// Timer -
	Timer, _ = time.ParseDuration("1h")

	// ContentURL -
	ContentURL = "https://yuque.yangchengkai.now.sh/api/v1/yuque/details?RepoID=%s&ID=%s"
	// ListURL -
	ListURL = "https://yuque.yangchengkai.now.sh/api/v1/yuque/list?RepoID=%s"
	// RepoURL -
	RepoURL = "https://yuque.yangchengkai.now.sh/api/v1/yuque/repo?GroupID=%s"
)
