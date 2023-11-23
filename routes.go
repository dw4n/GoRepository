// routes.go
package main

import (
	"gorepository/model"
	"gorepository/repository"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userRepo repository.UserRepository, postRepo repository.PostRepository) {

	app.Get("/users", func(c *fiber.Ctx) error {
		users, err := userRepo.GetAll()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch users"})
		}
		return c.JSON(users)
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
