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
	firstShelfMap, squareMap, columnCatalogMap, columnCoverMap, shelfListMap, columnListMap, contentMap sync.Map

	shelfListTime, contentTime, columnListTime, squareListTime, columnCatalogTime, columnCoverTime time.Time
)

// RegisterRouter -
func RegisterRouter(r gin.IRouter) {
	r.GET("/content", getContent)

	r.GET("/shelf/list", entrance)

	r.GET("/square/list", getSquareList)

	r.GET("/column/list", getColumnList)
	r.GET("/column/catalog", getColumnCatalog)
	r.GET("/column/cover", getColumnCover)
}

/*func getShelfList(c *gin.Context, ch chan interface{}) {
	var (
		List  RepoResp
		Resp  RespShelfList
		Resps []RespShelfList

		//ch = make(chan interface{}, 1)
	)

	ListNow := time.Now()
	interval := ListNow.Sub(shelfListTime)
	timer, _ := time.ParseDuration("1h")

	url := fmt.Sprintf(RepoURL, GroupID)

	val, ok := shelfListMap.Load(GroupID)
	if ok {
		if interval > timer {
			go getShelfCatalog(c)
			vall, _ := firstShelfMap.Load(FirstShelfRepoID)
			err := callAPI(c, url, &List)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "lists": val, "first": vall})
				return
			}

			for _, v := range List.Repo.Data {
				if v.Description == Bookshelf {
					Resp.ID = v.ID
					Resp.Title = v.Name
					Resps = append(Resps, Resp)
				}
			}

			shelfListMap.Store(GroupID, Resps)
			shelfListTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "lists": Resps})
			return
		}
		vall, _ := firstShelfMap.Load(FirstShelfRepoID)
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "lists": val, "first": vall})
		return
	}

	go getShelfCatalog(c)
		s := <-ch
		if s == errGet {
			c.Error(errGet)
			c.JSON(http.StatusFailedDependency, gin.H{"status": http.StatusFailedDependency})
			return
		}

		firstShelfMap.Store(FirstShelfRepoID, s)

	err := callAPI(c, url, &List)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	for _, v := range List.Repo.Data {
		if v.Description == Bookshelf {
			Resp.ID = v.ID
			Resp.Title = v.Name
			Resps = append(Resps, Resp)
		}
	}

	shelfListMap.Store(GroupID, Resps)
	shelfListTime = time.Now()
	vall, _ := firstShelfMap.Load(FirstShelfRepoID)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "lists": Resps, "first": vall})
}*/

func getContent(c *gin.Context) {
	var (
		content struct {
			RepoID string `json:"repo_id" binding:"required"`
			ID     string `json:"id"      binding:"required"`
		}

		Resp ContentRespon
	)

	DetailNow := time.Now()
	interval := DetailNow.Sub(contentTime)
	timer, _ := time.ParseDuration("1h")

	err := c.ShouldBind(&content)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	url := fmt.Sprintf(ContentURL, content.RepoID, content.ID)

	val, ok := contentMap.Load(content.ID)
	if ok {
		if interval > timer {
			err := callAPI(c, url, &Resp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "content": val})
				return
			}

			contentMap.Store(content.ID, Resp)
			contentTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "content": Resp})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "content": val})
		return
	}

	err = callAPI(c, url, &Resp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	contentMap.Store(content.ID, Resp)
	contentTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "content": Resp})
}

func getColumnList(c *gin.Context) {
	var (
		Repo  RepoResp
		Resp  RespColumnList
		Resps []RespColumnList
	)

	GroupNow := time.Now()
	interval := GroupNow.Sub(columnListTime)
	timer, _ := time.ParseDuration("1h")

	url := fmt.Sprintf(RepoURL, GroupID)

	val, ok := columnListMap.Load(GroupID)
	if ok {
		if interval > timer {
			err := callAPI(c, url, &Repo)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "column_list": val})
				return
			}

			for _, v := range Repo.Repo.Data {
				if v.Description == Column {
					Resp.ID = v.ID
					Resp.Title = v.Name
					Resp.Update = v.UpdatedAt
					Resps = append(Resps, Resp)
				}
			}

			columnListMap.Store(GroupID, Resps)
			columnListTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "column_list": Resps})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "column_list": val})
		return
	}

	err := callAPI(c, url, &Repo)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	for _, v := range Repo.Repo.Data {
		if v.Description == Column {
			Resp.ID = v.ID
			Resp.Title = v.Name
			Resp.Update = v.UpdatedAt
			Resps = append(Resps, Resp)
		}
	}

	columnListMap.Store(GroupID, Resps)
	columnListTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "column_list": Resps})
}

func getSquareList(c *gin.Context) {
	var (
		RepoResp ListRespon
		Resp     RespSquareList
		Resps    []RespSquareList
	)

	ListNow := time.Now()
	interval := ListNow.Sub(squareListTime)
	timer, _ := time.ParseDuration("1h")

	url := fmt.Sprintf(ListURL, SquareRepoID)

	val, ok := squareMap.Load(SquareRepoID)
	if ok {
		if interval > timer {
			err := callAPI(c, url, &RepoResp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "square_list": val})
				return
			}

			for _, v := range RepoResp.List.Data {
				if v.Status > 0 {
					Resp.ID = v.ID
					Resp.Title = v.Title
					Resp.Cover = v.Cover
					Resp.LikesCount = v.LikesCount
					Resp.Update = v.UpdatedAt
					Resp.Description = v.Description
					Resps = append(Resps, Resp)
				}
			}

			squareMap.Store(SquareRepoID, Resps)
			squareListTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "square_list": Resps})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "square_list": val})
		return
	}

	err := callAPI(c, url, &RepoResp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	for _, v := range RepoResp.List.Data {
		if v.Status > 0 {
			Resp.ID = v.ID
			Resp.Title = v.Title
			Resp.Cover = v.Cover
			Resp.LikesCount = v.LikesCount
			Resp.Update = v.UpdatedAt
			Resp.Description = v.Description
			Resps = append(Resps, Resp)
		}
	}

	squareMap.Store(SquareRepoID, Resps)
	squareListTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "square_list": Resps})
}

func getColumnCatalog(c *gin.Context) {
	var (
		column struct {
			RepoID string `json:"repo_id" binding:"required"`
		}

		RepoResp ListRespon
		Resp     RespColumnCatalog
		Resps    []RespColumnCatalog
	)

	ListNow := time.Now()
	interval := ListNow.Sub(columnCatalogTime)
	timer, _ := time.ParseDuration("1h")

	err := c.ShouldBind(&column)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	url := fmt.Sprintf(ListURL, column.RepoID)

	val, ok := columnCatalogMap.Load(column.RepoID)
	if ok {
		if interval > timer {
			err := callAPI(c, url, &RepoResp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "column_catalog": val})
				return
			}

			for _, v := range RepoResp.List.Data {
				if v.Status > 0 {
					Resp.ID = v.ID
					Resp.Title = v.Title
					Resp.Cover = v.Cover
					Resp.Update = v.PublishedAt
					Resps = append(Resps, Resp)
				}
			}

			columnCatalogMap.Store(column.RepoID, Resps)
			columnCatalogTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "column_catalog": Resps})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "column_catalog": val})
		return
	}

	err = callAPI(c, url, &RepoResp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	for _, v := range RepoResp.List.Data {
		if v.Status > 0 {
			Resp.ID = v.ID
			Resp.Title = v.Title
			Resp.Cover = v.Cover
			Resp.Update = v.PublishedAt
			Resps = append(Resps, Resp)
		}
	}

	columnCatalogMap.Store(column.RepoID, Resps)
	columnCatalogTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "column_catalog": Resps})
}

func getColumnCover(c *gin.Context) {
	var (
		RepoResp ListRespon
		Resp     RespColumnCover
		Resps    []RespColumnCover
	)

	ListNow := time.Now()
	interval := ListNow.Sub(columnCoverTime)
	timer, _ := time.ParseDuration("1h")

	url := fmt.Sprintf(ListURL, CoverRepoID)

	val, ok := columnCoverMap.Load(CoverRepoID)
	if ok {
		if interval > timer {
			err := callAPI(c, url, &RepoResp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "column_cover": val})
				return
			}

			for _, v := range RepoResp.List.Data {
				if v.Status > 0 {
					Resp.ID = v.ID
					Resp.Title = v.Title
					Resp.Cover = v.Cover
					Resps = append(Resps, Resp)
				}
			}

			columnCoverMap.Store(CoverRepoID, Resps)
			columnCoverTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "column_cover": Resps})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "column_cover": val})
		return
	}

	err := callAPI(c, url, &RepoResp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	for _, v := range RepoResp.List.Data {
		if v.Status > 0 {
			Resp.ID = v.ID
			Resp.Title = v.Title
			Resp.Cover = v.Cover
			Resps = append(Resps, Resp)
		}
	}

	columnCoverMap.Store(CoverRepoID, Resps)
	columnCoverTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "column_cover": Resps})
}

func getShelfList(c *gin.Context, ch chan interface{}) {
	var (
		List  RepoResp
		Resp  RespShelfList
		Resps []RespShelfList
	)

	url := fmt.Sprintf(RepoURL, GroupID)

	err := callAPI(c, url, &List)
	if err != nil {
		c.Error(err)
		ch <- errGetList
	}

	for _, v := range List.Repo.Data {
		if v.Description == Bookshelf {
			Resp.ID = v.ID
			Resp.Title = v.Name
			Resps = append(Resps, Resp)
		}
	}

	ch <- Resps
}

func getFirstShelfRepo(c *gin.Context, ch chan interface{}) {
	var (
		RepoResp ListRespon
		Resp     RespShelfCatalog
		Resps    []RespShelfCatalog
	)

	url := fmt.Sprintf(ListURL, FirstShelfRepoID)

	err := callAPI(c, url, &RepoResp)
	if err != nil {
		c.Error(err)
		ch <- errGetFirst
	}

	for _, v := range RepoResp.List.Data {
		if v.Status > 0 {
			Resp.ID = v.ID
			Resp.Title = v.Title
			Resp.Cover = v.Cover
			Resps = append(Resps, Resp)
		}
	}
	//firstShelfMap.Store(FirstShelfRepoID, Resps)

	ch <- Resps
}

//entrance
func entrance(c *gin.Context) {
	var chFirst = make(chan interface{}, 1)
	var chList = make(chan interface{}, 1)

	go getFirstShelfRepo(c, chFirst)
	firstShelfRepo := <-chFirst
	if firstShelfRepo == errGetFirst {
		c.Error(errGetFirst)
		c.JSON(http.StatusFailedDependency, gin.H{"status": http.StatusFailedDependency})
		return
	}
	go getShelfList(c, chList)
	shelfList := <-chList
	if shelfList == errGetList {
		c.Error(errGetList)
		c.JSON(http.StatusFailedDependency, gin.H{"status": http.StatusFailedDependency})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "column_cover": firstShelfRepo, "list": shelfList})
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
