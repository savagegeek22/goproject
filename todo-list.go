package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

var (
	todos  = []Todo{}
	nextID = 1
	mu     sync.Mutex
)

func main() {
	r := gin.Default()

	// Serve the frontend at "/"
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlPage)
	})

	// API routes
	r.GET("/todos", getTodos)
	r.POST("/todos", addTodo)
	r.DELETE("/todos/:id", deleteTodo)

	r.Run(":8080")
}

func getTodos(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, todos)
}

func addTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil || newTodo.Task == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo"})
		return
	}

	mu.Lock()
	newTodo.ID = nextID
	nextID++
	todos = append(todos, newTodo)
	mu.Unlock()

	c.JSON(http.StatusCreated, newTodo)
}

func deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

const htmlPage = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8" />
	<title>Todo List</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 600px; margin: 20px auto; }
		form { margin-bottom: 20px; }
		ul { padding-left: 20px; }
		li { margin-bottom: 8px; }
		button { margin-left: 10px; }
	</style>
</head>
<body>
	<h1>Todo List</h1>
	<form id="todo-form">
		<input type="text" id="task" placeholder="Enter todo" required />
		<button type="submit">Add Todo</button>
	</form>
	<ul id="todo-list"></ul>

	<script>
		async function fetchTodos() {
			const res = await fetch('/todos');
			const todos = await res.json();
			const list = document.getElementById('todo-list');
			list.innerHTML = '';
			todos.forEach(todo => {
				const li = document.createElement('li');
				li.textContent = todo.task;

				const btn = document.createElement('button');
				btn.textContent = 'Delete';
				btn.onclick = async () => {
					await fetch('/todos/' + todo.id, { method: 'DELETE' });
					fetchTodos();
				};

				li.appendChild(btn);
				list.appendChild(li);
			});
		}

		document.getElementById('todo-form').addEventListener('submit', async (e) => {
			e.preventDefault();
			const taskInput = document.getElementById('task');
			const task = taskInput.value.trim();
			if (!task) return;

			await fetch('/todos', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ task })
			});
			taskInput.value = '';
			fetchTodos();
		});

		// Load todos on page load
		fetchTodos();
	</script>
</body>
</html>
`

