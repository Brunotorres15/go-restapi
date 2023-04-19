package routes

import (
	"errors"
	"time"

	"github.com/Brunotorres15/go-restapi/database"
	"github.com/Brunotorres15/go-restapi/models"
	"github.com/gofiber/fiber/v2"
)

type OrderSerializer struct {
	ID        uint              `json:"id"`
	User      User              `json:"user"`
	Product   ProductSerializer `json:"product"`
	CreatedAt time.Time         `kspm:"order_date"`
}

func CreateResponseOrder(order models.Order, user User, product ProductSerializer) OrderSerializer {
	return OrderSerializer{ID: order.ID, User: user, Product: product, CreatedAt: order.CreatedAt}
}

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	// pega as informações e coloca na variável order (informação de id do usuário, id do produto, etc...)
	if err := c.BodyParser(&order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// verifica se o usuário existe
	var user models.User
	if err := findUser(order.UserRefer, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// verifica se o produto existe
	var product models.Product
	if err := findProduct(order.ProductRefer, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// cria a order no banco
	database.Database.Db.Create(&order)

	// serialização para usar no return
	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(200).JSON(responseOrder)
}

func findOrder(id int, order *models.Order) error {
	database.Database.Db.Find(&order, "id= ?", id)
	if order.ID == 0 {
		return errors.New("Order does not exist.")
	}
	return nil
}

func GetOrders(c *fiber.Ctx) error {

	orders := []models.Order{}

	database.Database.Db.Find(&orders)
	responseOrders := []OrderSerializer{}

	for _, order := range orders {
		var user models.User
		var product models.Product
		database.Database.Db.Find(&user, "id = ?", order.UserRefer)
		database.Database.Db.Find(&product, "id = ?", order.ProductRefer)
		responseOrder := CreateResponseOrder(order, CreateResponseUser(user), CreateResponseProduct(product))
		responseOrders = append(responseOrders, responseOrder)

	}

	return c.Status(200).JSON(responseOrders)

}

func GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var order models.Order

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer.")
	}

	if err := findOrder(id, &order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var user models.User
	var product models.Product

	database.Database.Db.First(&user, order.UserRefer)
	database.Database.Db.First(&product, order.ProductRefer)
	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)

	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(200).JSON(responseOrder)
}
