package handler

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/Rakhimgaliev/tech-db-forum/project/db"
	"github.com/Rakhimgaliev/tech-db-forum/project/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func BindJSON(context *gin.Context, obj interface{}) error {
	if err := binding.JSON.Bind(context.Request, obj); err != nil {
		context.Error(err).SetType(gin.ErrorTypeBind)
		return err
	}
	return nil
}

func (h handler) CreatePosts(context *gin.Context) {
	posts := models.Posts{}

	err := BindJSON(context, &posts)
	if err != nil {
		context.JSON(400, err)
		return
	}

	err = db.CreatePosts(h.conn, context.Param("slug_or_id"), &posts)
	log.Println(err, posts)
	if err != nil {
		switch err {
		case db.ErrorPostCreateBadRequest:
			postsJSON, _ := json.Marshal(posts)
			context.Data(201, "application/json", postsJSON)
		case db.ErrorThreadNotFound, db.ErrorForumNotFound:
			context.JSON(404, err)
		case db.ErrorPostCreateConflict:
			context.JSON(409, err)
		default:
			context.JSON(500, err)
		}
		return
	}

	postsJSON, _ := json.Marshal(posts)
	context.Data(201, "application/json", postsJSON)
}

func (h handler) GetPosts(context *gin.Context) {
	posts := models.Posts{}
	queryArgs := context.Request.URL.Query()

	limit := 0
	if len(queryArgs["limit"]) > 0 {
		limit, _ = strconv.Atoi(queryArgs["limit"][0])
	}

	since := 0
	if len(queryArgs["since"]) > 0 {
		since, _ = strconv.Atoi(queryArgs["since"][0])
	}

	desc := false
	if len(queryArgs["desc"]) > 0 {
		if queryArgs["desc"][0] == "true" {
			desc = true
		}
	}

	sort := ""
	if len(queryArgs["sort"]) > 0 {
		sort = queryArgs["sort"][0]
	}

	err := db.GetPosts(h.conn, context.Param("slug_or_id"),
		limit, desc,
		since, sort, &posts)
	log.Println("-------------------err: ", err)

	if err != nil {
		if err == db.ErrorThreadNotFound {
			context.JSON(404, err)
			return
		}
		context.JSON(500, err)
		return
	}
	log.Println("posts count: ", len(posts))
	postsJSON, _ := json.Marshal(posts)
	context.Data(200, "application/json", postsJSON)
	return
}
