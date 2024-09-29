package db

import (


	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	dsn := "host=aws-0-us-east-1.pooler.supabase.com user=postgres.ehngutoknbaqeoszckbe password=eZ1wx5Gtt3SEH7O3 dbname=postgres port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = database

	// Run migrations
	// if err := db.Debug().AutoMigrate(&models.User{}, &models.Problem{}, &models.BaseModel{}, &models.FunctionSignature{}, &models.Parameter{}, &models.Submission{}, &models.TestCase{}); err != nil {
	// log.Fatalf("Migration failed: %v", err)
	// } else {
	//     log.Println("Migrations ran successfully")
	// }

}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}
