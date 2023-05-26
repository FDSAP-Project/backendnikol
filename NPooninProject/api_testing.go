package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ---------- STRUCT ----------//

type Repository struct {
	DB *gorm.DB
}
type userlogin struct {
	Name      string `json:"name"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	ContactNo string `json:"contact"`
}
type Configuration struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

//---------- END ----------//

func main() {
	//getting env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	//config
	config := &Configuration{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSSLMODE"),
		DBName:   os.Getenv("DB_DBNAME"),
	}

	//connecting to db
	db, err := NewConnection(config)
	if err != nil {
		fmt.Println(err)
	}

	//migration to db
	err = MigrateClientDetails(db)
	if err != nil {
		fmt.Println(err)
	}

	repo := Repository{
		DB: db,
	}

	app := fiber.New()

	repo.SetupRoute(app)

	log.Fatal(app.Listen(":3000"))
}

// ---------- FUNCTIONS ----------//
func NewConnection(config *Configuration) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s password=%s user=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Password, config.User, config.DBName, config.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	return db, nil
}

func MigrateClientDetails(db *gorm.DB) error {
	err := db.AutoMigrate(&ClientDetails{})
	return err
}

// Set Routes
func (repo *Repository) SetupRoute(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/login", repo.InsertDataFromDB)

}

// INSERT API
func (repo *Repository) InsertDataFromDB(c *fiber.Ctx) error {
	fetchedData := ClientDetails{}

	repo.DB.Table("client_details").Select("*").Find(&fetchedData)

	return c.JSON(fetchedData)
}

//---------- END ----------//
