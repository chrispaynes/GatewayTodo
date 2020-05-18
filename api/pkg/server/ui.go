package server

import (
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ServeUI serves the Angular static files from the directory
func ServeUI() {
	r := gin.Default()

	log.Info("serving UI files")
	handler := func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)

		if file == "" || ext == "" {
			c.File("./dist/todo-app/index.html")
		} else {
			c.File("./dist/todo-app/" + path.Join(dir, file))
		}
	}

	r.Group("/")

	log.Info("serving UI files")
	r.NoRoute(handler)

	if err := r.Run(":" + conf.UIPort); err != nil {
		log.Info("fatal")
		log.Fatal(err)
	}

	log.Info("serving UI files")
}
