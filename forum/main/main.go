package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"

	_ "github.com/cockroachdb/apd"
	_ "github.com/hashicorp/go-version"
	_ "github.com/jackc/fake"
	_ "github.com/lib/pq"
	_ "github.com/pkg/errors"
	_ "github.com/satori/go.uuid"
	_ "github.com/shopspring/decimal"
	_ "github.com/sirupsen/logrus"
	_ "go.uber.org/zap"
	_ "gopkg.in/inconshreveable/log15.v2"
)

var pgxConfig = pgx.ConnConfig{
	Host:     "localhost",
	Port:     5432,
	Database: "docker",
	User:     "docker",
	Password: "docker",
}

func main() {
	log.Println("start")

	ConnPoolConfig := pgx.ConnPoolConfig{pgxConfig, 3, nil, 0}

	conn, err := pgx.NewConnPool(ConnPoolConfig)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("YEEH")
	}

	r := gin.Default()
	r.POST("/forum/:slug", func(c *gin.Context) {
		if c.Param("slug") == "create" {
			c.JSON(200, gin.H{
				"description": `Создание нового форума.`,
			})
		}
		c.JSON(200, gin.H{
			"description": `hz`,
		})
	})

	r.GET("/forum/:slug/details", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Получение информации о форуме по его идентификаторе.`,
		})
	})

	r.POST("/forum/:slug/create", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Добавление новой ветки обсуждения на форум.`,
		})
	})

	r.GET("/forum/:slug/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Получение списка пользователей, у которых есть пост или ветка обсуждения в данном форуме.
								Пользователи выводятся отсортированные по nickname в порядке возрастания.
								Порядок сотрировки должен соответсвовать побайтовому сравнение в нижнем регистре.`,
		})
	})

	r.GET("/forum/:slug/threads", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Получение списка ветвей обсужления данного форума.
								Ветви обсуждения выводятся отсортированные по дате создания.`,
		})
	})

	r.GET("/post/:id/details", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Получение информации о ветке обсуждения по его имени.`,
		})
	})

	r.POST("/post/:id/details", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Изменение сообщения на форуме.
								Если сообщение поменяло текст, то оно должно получить отметку "isEdited".`,
		})
	})

	r.POST("/service/clear", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Очистка всех данных в базе`,
		})
	})

	r.GET("/service/status", func(c *gin.Context) {
		Status(c, conn)
	})

	r.POST("/thread/:slug_or_id/create", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Добавление новых постов в ветку обсуждения на форум.
								Все посты, созданные в рамках одного вызова данного метода должны иметь одинаковую дату создания (Post.Created).`,
		})
	})

	r.GET("/thread/:slug_or_id/details", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Получение информации о ветке обсуждения по его имени.`,
		})
	})

	r.POST("/thread/:slug_or_id/details", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Обновление ветки обсуждения на форуме.`,
		})
	})

	r.GET("/thread/:slug_or_id/posts", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Получение информации о ветке обсуждения по его имени.`,
		})
	})

	r.POST("/thread/:slug_or_id/vote", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Изменение голоса за ветвь обсуждения.
								Один пользователь учитывается только один раз и может изменить своё
								мнение.`,
		})
	})

	r.POST("/user/:nickname/create", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Создание нового пользователя в базе данных.`,
		})
	})

	r.GET("/user/:nickname/profile", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Получение информации о пользователе форума по его имени.`,
		})
	})

	r.POST("/user/:nickname/profile", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"description": `Изменение информации в профиле пользователя.`,
		})
	})

	r.Run("127.0.0.1:5000")
}
