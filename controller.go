package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	yuque = make(map[string]string)

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

	ListNow := time.Now()
	interval := ListNow.Sub(ListTime)
	timer, _ := time.ParseDuration("24h")

	err := c.ShouldBind(&list)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	url := "http://127.0.0.1:9090/list?RepoID=" + list.RepoID

	val, ok := lr.Load(list.RepoID)
	if ok {
		if interval > timer {
			newbody, err := callAPI(c, url)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
				return
			}

			lr.Store(list.RepoID, string(newbody))
			ListTime = time.Now()
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "List": string(newbody)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "List": val})
		return
	}

	body, err := callAPI(c, url)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden})
		return
	}

	lr.Store(list.RepoID, string(body))
	ListTime = time.Now()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "List": string(body)})
}

func getDetails(c *gin.Context) {
	var (
		detail struct {
			RepoID string `json:"repoid"`
			ID     string `json:"id"`
		}
	)

	DetailNow := time.Now()
	interval := DetailNow.Sub(DetailTime)
	timer, _ := time.ParseDuration("24h")

	err := c.ShouldBind(&detail)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	url := "http://127.0.0.1:9090/detail?RepoID=" + detail.RepoID + "&ID=" + detail.ID

	val, ok := dr.Load(detail.RepoID)
	fmt.Println(ok, "ppppp")
	if ok {
		if interval > timer {
			newbody, err := callAPI(c, url)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "List": val})
				return
			}

			handleDetail(newbody, detail.RepoID)
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "List": yuque})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "List": val})
		return
	}

	body, err := callAPI(c, url)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	handleDetail(body, detail.RepoID)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Detail": yuque})

}
func getRepo(c *gin.Context) {
	var Group struct {
		GroupID  string `json:"groupid"`
		RepoName string `json:"reponame"`
	}

	GroupNow := time.Now()
	interval := GroupNow.Sub(GroupTime)
	timer, _ := time.ParseDuration("24h")

	err := c.ShouldBind(&Group)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable})
		return
	}

	url := "http://127.0.0.1:9090/groups?GroupID=" + Group.GroupID + "&RepoName=" + Group.RepoName

	val, ok := gr.Load(Group.RepoName)
	if ok {
		if interval > timer {
			newbody, err := callAPI(c, url)
			if err != nil {
				c.Error(err)
				c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
				return
			}
			gr.Store(Group.RepoName, string(newbody))
			GroupTime = time.Now()

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "GroupRepo": string(newbody)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "GroupRepo": val})
		return
	}

	body, err := callAPI(c, url)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	gr.Store(Group.RepoName, string(body))
	GroupTime = time.Now()

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "RepoList": string(body)})
}

func callAPI(c *gin.Context, url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Error(err)
		return nil, err
	}

	Token := c.Request.Header
	request.Header.Add("X-Auth-Token", Token["X-Auth-Token"][0])

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		c.Error(err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.Error(err)
		return nil, err
	}

	check := bytes.Contains(body, []byte("ERROR"))
	if check {
		log.Println("Requst ERROR")
		return nil, err
	}

	return body, nil
}

func handleDetail(body []byte, id string) {

	var con = make([]string, 5)
	con[0] = "封面"
	con[1] = "作者"
	con[2] = "内容简介"
	con[3] = "捐助者"
	con[4] = "时间"
	info := bytes.Split(body, []byte("****"))

	yuque = map[string]string{
		con[0]: string(info[1]),
		con[1]: string(info[3]),
		con[2]: string(info[5]),
		con[3]: string(info[7]),
		con[4]: string(info[9]),
	}
	DetailTime = time.Now()
	dr.Store(id, yuque)

}
