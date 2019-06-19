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
	// listRecode
	lr sync.Map
	// detailRecode
	dr sync.Map
	// groupRecode
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

	r.GET("/getlist", getList)
	r.GET("/getdetails", getDetails)
	r.GET("/getrepo", getRepo)

}

func getList(c *gin.Context) {
	var list struct {
		RepoID string `json:"repoid"`
	}

	var RepoResp ListRespon

	ListNow := time.Now()
	interval := ListNow.Sub(ListTime)
	timer, _ := time.ParseDuration("24h")

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
			RepoID string `json:"repoid"`
			ID     string `json:"id"`
		}
	)

	var DeResp DetailRespon

	DetailNow := time.Now()
	interval := DetailNow.Sub(DetailTime)
	timer, _ := time.ParseDuration("10s")

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
			GroupID  string `json:"groupid"`
			RepoName string `json:"reponame"`
		}

		Repo RepoResp
	)

	GroupNow := time.Now()
	interval := GroupNow.Sub(GroupTime)
	timer, _ := time.ParseDuration("24h")

	err := c.ShouldBind(&Group)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	url := fmt.Sprintf(RepoURL, Group.GroupID)

	val, ok := gr.Load(Group.RepoName)
	if ok {
		if interval > timer {
			err := callAPI(c, url, &Repo)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway, "group_repos": val})
				return
			}
			gr.Store(Group.RepoName, Repo)
			GroupTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "group_repos": Repo})
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

	gr.Store(Group.RepoName, Repo)
	GroupTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "group_repos": Repo})
}

func callAPI(c *gin.Context, url string, obj interface{}) error {
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	Token := c.Request.Header
	request.Header.Add("X-Auth-Token", Token["X-Auth-Token"][0])

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
