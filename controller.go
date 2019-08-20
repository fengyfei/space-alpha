package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	likeMap sync.Map
	lastMap sync.Map
	allMap  sync.Map

	lr sync.Map
	dr sync.Map
	gr sync.Map

	// ListTime -
	ListTime time.Time
	// DetailTime -
	DetailTime time.Time
	// GroupTime -
	GroupTime time.Time
)

// RegisterRouter -
func RegisterRouter(r gin.IRouter) {

	r.GET("/list/recommend", getRecommendList)
	r.GET("/list/lastest", getLastestList)
	r.GET("/list/all", getAllList)
	r.GET("/getdetails", getDetails)
	r.GET("/getrepo", getRepo)

}

func getList(c *gin.Context) {
	var list struct {
		RepoID string `json:"repo_id" binding:"required"`
	}

	var RepoResp ListRespon

	ListNow := time.Now()
	interval := ListNow.Sub(ListTime)
	timer, _ := time.ParseDuration("1h")

	err := c.ShouldBind(&list)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	url := fmt.Sprintf(ListURL, list.RepoID)

	val, ok := lr.Load(list.RepoID)
	if ok {
		if interval > timer {
			err := callAPI(c, url, &RepoResp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "lists": val})
				return
			}

			lr.Store(list.RepoID, RepoResp)
			ListTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "lists": RepoResp})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "lists": val})
		return
	}

	err = callAPI(c, url, &RepoResp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	lr.Store(list.RepoID, RepoResp)
	ListTime = time.Now()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "lists": RepoResp})
}

func getDetails(c *gin.Context) {
	var (
		detail struct {
			RepoID string `json:"repo_id" binding:"required"`
			ID     string `json:"id"      binding:"required"`
		}
	)

	var DeResp DetailRespon

	DetailNow := time.Now()
	interval := DetailNow.Sub(DetailTime)
	timer, _ := time.ParseDuration("1h")

	err := c.ShouldBind(&detail)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	url := fmt.Sprintf(DetailURL, detail.RepoID, detail.ID)

	val, ok := dr.Load(detail.ID)
	if ok {
		if interval > timer {
			err := callAPI(c, url, &DeResp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "details": val})
				return
			}

			dr.Store(detail.ID, DeResp)
			DetailTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "details": DeResp})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "details": val})
		return
	}

	err = callAPI(c, url, &DeResp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	dr.Store(detail.ID, DeResp)
	DetailTime = time.Now()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "details": DeResp})

}

func getRepo(c *gin.Context) {
	var (
		Group struct {
			GroupID string `json:"group_id" binding:"required"`
		}

		Repo  RepoResp
		Resp  RespRepo
		Resps []RespRepo
	)

	GroupNow := time.Now()
	interval := GroupNow.Sub(GroupTime)
	timer, _ := time.ParseDuration("1h")

	err := c.ShouldBind(&Group)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	url := fmt.Sprintf(RepoURL, Group.GroupID)

	val, ok := gr.Load(Group.GroupID)
	if ok {
		if interval > timer {
			err := callAPI(c, url, &Repo)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "group_repos": val})
				return
			}

			for _, v := range Repo.Repo.Data {
				Resp.ID = v.ID
				Resp.Name = v.Name
				Resps = append(Resps, Resp)
			}

			gr.Store(Group.GroupID, Resps)
			GroupTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "group_repos": Resps})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "group_repos": val})
		return
	}

	err = callAPI(c, url, &Repo)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	for _, v := range Repo.Repo.Data {
		Resp.ID = v.ID
		Resp.Name = v.Name
		Resps = append(Resps, Resp)
	}

	gr.Store(Group.GroupID, Resps)
	GroupTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "group_repos": Resps})
}

func getRecommendList(c *gin.Context) {
	var (
		RepoResp ListRespon
		Resp     RespRecommendList
		Resps    []RespRecommendList
	)

	ListNow := time.Now()
	interval := ListNow.Sub(ListTime)
	timer, _ := time.ParseDuration("1h")

	val, ok := likeMap.Load(ArticleRepoID)
	if ok {
		if interval > timer {
			err := callAPI(c, RecommendListURL, &RepoResp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "like_lists": val})
				return
			}

			for _, v := range RepoResp.List.Data {
				if v.LikesCount > 0 && v.Status > 0 {
					Resp.Title = v.Title
					Resp.Cover = v.Cover
					Resp.LikesCount = v.LikesCount
					Resp.Description = v.CustomDescription
					Resps = append(Resps, Resp)
				}
			}

			likeMap.Store(ArticleRepoID, Resps)
			ListTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "like_lists": Resps})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "like_lists": val})
		return
	}

	err := callAPI(c, RecommendListURL, &RepoResp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	for _, v := range RepoResp.List.Data {
		if v.LikesCount > 0 && v.Status > 0 {
			Resp.Title = v.Title
			Resp.Cover = v.Cover
			Resp.LikesCount = v.LikesCount
			Resp.Description = v.Description
			Resps = append(Resps, Resp)
		}
	}

	likeMap.Store(ArticleRepoID, Resps)
	ListTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "like_lists": Resps})
}

func getLastestList(c *gin.Context) {
	var (
		RepoResp ListRespon
		Resp     RespLastestList
		Resps    []RespLastestList
	)

	ListNow := time.Now()
	interval := ListNow.Sub(ListTime)
	timer, _ := time.ParseDuration("1h")

	val, ok := lastMap.Load(ArticleRepoID)
	if ok {
		if interval > timer {
			err := callAPI(c, RecommendListURL, &RepoResp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "last_lists": val})
				return
			}

			for _, v := range RepoResp.List.Data {
				t := ListNow.Sub(v.PublishedAt)
				if t < Lastest && v.Status > 0 {
					Resp.Title = v.Title
					Resp.Cover = v.Cover
					Resp.Date = v.PublishedAt
					Resp.Description = v.Description
					Resps = append(Resps, Resp)
				}
			}

			lastMap.Store(ArticleRepoID, Resps)
			ListTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "last_lists": Resps})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "last_lists": val})
		return
	}

	err := callAPI(c, RecommendListURL, &RepoResp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	for _, v := range RepoResp.List.Data {
		t := ListNow.Sub(v.PublishedAt)
		if t < Lastest && v.Status > 0 {
			Resp.Title = v.Title
			Resp.Cover = v.Cover
			Resp.Date = v.PublishedAt
			Resp.Description = v.Description
			Resps = append(Resps, Resp)
		}
	}

	lastMap.Store(ArticleRepoID, Resps)
	ListTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "last_lists": Resps})
}

func getAllList(c *gin.Context) {
	var (
		RepoResp ListRespon
		Resp     RespRepo
		Resps    []RespRepo
	)

	ListNow := time.Now()
	interval := ListNow.Sub(ListTime)
	timer, _ := time.ParseDuration("1h")

	val, ok := allMap.Load(ArticleRepoID)
	if ok {
		if interval > timer {
			err := callAPI(c, RecommendListURL, &RepoResp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "all_lists": val})
				return
			}

			for _, v := range RepoResp.List.Data {
				if v.Status > 0 {
					Resp.ID = v.ID
					Resp.Name = v.Title
					Resps = append(Resps, Resp)
				}
			}

			allMap.Store(ArticleRepoID, Resps)
			ListTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "all_lists": Resps})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "all_lists": val})
		return
	}

	err := callAPI(c, RecommendListURL, &RepoResp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	for _, v := range RepoResp.List.Data {
		if v.Status > 0 {
			Resp.ID = v.ID
			Resp.Name = v.Title
			Resps = append(Resps, Resp)
		}
	}

	allMap.Store(ArticleRepoID, Resps)
	ListTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "all_lists": Resps})
}

func callAPI(c *gin.Context, url string, obj interface{}) error {
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	t := c.Request.Header.Get("X-Auth-Token")
	request.Header.Add("X-Auth-Token", t)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, obj)
}
