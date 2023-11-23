// routes.go
package main

import (
	"gorepository/model"
	"gorepository/repository"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, repos *repository.Repositories) {

	userRepo := repos.UserRepo
	postRepo := repos.PostRepo

	app.Get("/users", func(c *fiber.Ctx) error {
		// Parse pagination parameters
		page, err := strconv.Atoi(c.Query("page", "1"))
		if err != nil || page < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid page number"})
		}

		pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
		if err != nil || pageSize <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid page size"})
		}

		// This could be multiple query parameters or a single comma-separated parameter
		sortParams := c.Query("sort", "Name ASC")
		sortColumns := strings.Split(sortParams, ",")

		// Handle any potential conversion errors
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid pagination parameters"})
		}

		conditions := []func(*gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Where("is_deleted = ?", false)
			},
		}

		// Add sorting and pagination options
		opts := []repository.Option{
			repository.WithSorting(sortColumns),
			repository.WithPaging(page, pageSize),
		}

		var users []model.User
		err = repos.UserRepo.GetAllWithConditions(&users, conditions, opts...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch users"})
		}

		return c.JSON(users)
	})

	app.Get("/post/user/:id", func(c *fiber.Ctx) error {
		userID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
		}

		// Parse pagination parameters
		page, err := strconv.Atoi(c.Query("page", "1"))
		if err != nil || page < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid page number"})
		}

		pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
		if err != nil || pageSize <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid page size"})
		}

		// Parse sorting parameters
		sortParams := c.Query("sort", "posts.created_at ASC")
		sortColumns := strings.Split(sortParams, ",")

		conditions := []func(*gorm.DB) *gorm.DB{
			func(db *gorm.DB) *gorm.DB {
				return db.Joins("JOIN users ON users.id = posts.user_id").
					Select("posts.*, users.name as user_name").
					Where("posts.user_id = ?", userID)
			},
		}

		// Add sorting and pagination options
		opts := []repository.Option{
			repository.WithSorting(sortColumns),
			repository.WithPaging(page, pageSize),
		}

		var posts []model.PostWithUserName // Define a slice to store the result
		err = repos.PostRepo.GetAllWithConditions(&posts, conditions, opts...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch posts"})
		}

		return c.JSON(posts)
	})

	app.Post("/user", func(c *fiber.Ctx) error {
		var user model.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		newUser, err := userRepo.Create(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create user"})
		}

		return c.JSON(newUser)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}

		user, err := userRepo.FindByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user not found"})
		}

		return c.JSON(user)
	})

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}

		var user model.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		user.ID = uint(id)
		updatedUser, err := userRepo.Update(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot update user"})
		}

		return c.JSON(updatedUser)
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}

		err = userRepo.Delete(uint(id))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete user"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	// Routes for listing all posts
	app.Get("/posts", func(c *fiber.Ctx) error {
		posts, err := postRepo.GetAll()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch posts"})
		}
		return c.JSON(posts)
	})

	// Route for creating a new post
	app.Post("/post", func(c *fiber.Ctx) error {
		var post model.Post
		if err := c.BodyParser(&post); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		newPost, err := postRepo.Create(post)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create post"})
		}

		return c.JSON(newPost)
	})

	// Route for getting a single post by ID
	app.Get("/posts/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}

		post, err := postRepo.FindByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "post not found"})
		}

		return c.JSON(post)
	})

	// Route for updating a post
	app.Put("/posts/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}

		var post model.Post
		if err := c.BodyParser(&post); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		post.ID = uint(id)
		updatedPost, err := postRepo.Update(post)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot update post"})
		}

		return c.JSON(updatedPost)
	})

	// Route for deleting a post
	app.Delete("/posts/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
		}

		err = postRepo.Delete(uint(id))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete post"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
