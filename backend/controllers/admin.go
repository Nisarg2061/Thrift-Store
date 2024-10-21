package controllers

import (
	"github.com/TMP-The-Major-Project/Thrift-Store/backend/database"
	"github.com/TMP-The-Major-Project/Thrift-Store/backend/models"
	"github.com/gofiber/fiber/v2"
)

func CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}
	db := database.Connect()
	db.Create(&product)
	return c.Status(fiber.StatusCreated).JSON(product)
}

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product

	db := database.Connect()

	db.Find(&products)
	return c.Status(fiber.StatusOK).JSON(products)
}

func AddToCart(c *fiber.Ctx) error {
	cartItem := new(models.CartItem)
	if err := c.BodyParser(cartItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	db := database.Connect()

	// Check if the product exists
	var product models.Product
	if err := db.First(&product, cartItem.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	// Add to cart
	cartItem.TotalPrice = product.NewPrice * float64(cartItem.Quantity)
	db.Create(&cartItem)

	return c.Status(fiber.StatusCreated).JSON(cartItem)
}

func RemoveFromCart(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.Connect()

	// Find and delete the cart item
	if err := db.Delete(&models.CartItem{}, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Item removed from cart"})
}

func GetCartTotal(c *fiber.Ctx) error {
	var cartItems []models.CartItem
	db := database.Connect()

	// Find all items in the cart
	db.Find(&cartItems)

	var total float64 = 0

	// Sum the total prices of each item in the cart
	for _, item := range cartItems {
		total += item.TotalPrice
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"cart_total": total})
}

func GetCartItems(c *fiber.Ctx) error {
	var cartItems []models.CartItem
	db := database.Connect()

	// Retrieve all items in the cart
	if err := db.Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve cart items"})
	}

	// Create a response object
	response := models.CartResponse{
		Cart: cartItems,
	}

	// Return the response with a 200 status
	return c.Status(fiber.StatusOK).JSON(response)
}
