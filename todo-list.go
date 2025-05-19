package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// Todo represents a single todo item
type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

var (
	todoList []Todo
	nextID   = 1
	mutex    sync.Mutex
)

func main() {
	r := gin.Default()

	// Routes
	r.GET("/todos", getTodos)
	r.POST("/todos", addTodo)
	r.DELETE("/todos/:id", deleteTodo)

	// Start server
	r.Run(":8080")
}

func getTodos(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()
	c.JSON(http.StatusOK, todoList)
}

func addTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil || newTodo.Task == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	newTodo.ID = nextID
	nextID++
	todoList = append(todoList, newTodo)
	c.JSON(http.StatusCreated, newTodo)
}

func deleteTodo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	for i, todo := range todoList {
		if todo.ID == id {
			todoList = append(todoList[:i], todoList[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
}

