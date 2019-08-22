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
	squareMap, columnCatalogMap, imageMap, lr, columnListMap, contentMap sync.Map

	listTime, contentTime, columnListTime, squareListTime, columnCatalogTime, imageListTime time.Time
)

// RegisterRouter -
func RegisterRouter(r gin.IRouter) {
	r.GET("/square/list", getSquareList)

	r.GET("/column/list", getColumnList)
	r.GET("/column/catalog", getColumnCatalog)
	r.GET("/column/content", getContent)

	r.GET("/image/list", getImageList)
}

func getList(c *gin.Context) {
	var (
		list struct {
			RepoID string `json:"repo_id" binding:"required"`
		}

		Resp ListRespon
	)

	ListNow := time.Now()
	interval := ListNow.Sub(listTime)
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
			err := callAPI(c, url, &Resp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "lists": val})
				return
			}

			lr.Store(list.RepoID, Resp)
			listTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "lists": Resp})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "lists": val})
		return
	}

	err = callAPI(c, url, &Resp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	lr.Store(list.RepoID, Resp)
	listTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "lists": Resp})
}

func getContent(c *gin.Context) {
	var (
		content struct {
			RepoID string `json:"repo_id" binding:"required"`
			ID     string `json:"id"      binding:"required"`
		}

		Resp DetailRespon
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

	url := fmt.Sprintf(DetailURL, content.RepoID, content.ID)

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
		Resp  RespRepo
		Resps []RespRepo
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
				if v.Description == "column" {
					Resp.ID = v.ID
					Resp.Name = v.Name
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
		if v.Description == "column" {
			Resp.ID = v.ID
			Resp.Name = v.Name
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

	url := fmt.Sprintf(ListURL, GroupID)

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
		Resp     RespColumn
		Resps    []RespColumn
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
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "last_lists": val})
				return
			}

			for _, v := range RepoResp.List.Data {
				if v.Status > 0 {
					Resp.Title = v.Title
					Resp.Cover = v.Cover
					Resp.Update = v.PublishedAt
					Resps = append(Resps, Resp)
				}
			}

			columnCatalogMap.Store(column.RepoID, Resps)
			columnCatalogTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "last_lists": Resps})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "last_lists": val})
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
			Resp.Title = v.Title
			Resp.Cover = v.Cover
			Resp.Update = v.PublishedAt
			Resps = append(Resps, Resp)
		}
	}

	columnCatalogMap.Store(column.RepoID, Resps)
	columnCatalogTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "last_lists": Resps})
}

func getImageList(c *gin.Context) {
	var (
		RepoResp ListRespon
		Resp     RespImage
		Resps    []RespImage
	)

	ListNow := time.Now()
	interval := ListNow.Sub(imageListTime)
	timer, _ := time.ParseDuration("1h")

	url := fmt.Sprintf(ListURL, ImageRepoID)

	val, ok := imageMap.Load(ImageRepoID)
	if ok {
		if interval > timer {
			err := callAPI(c, url, &RepoResp)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "image_lists": val})
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

			imageMap.Store(ImageRepoID, Resps)
			imageListTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "image_lists": Resps})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "image_lists": val})
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
			Resps = append(Resps, Resp)
		}
	}

	imageMap.Store(ImageRepoID, Resps)
	imageListTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "image_lists": Resps})
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
