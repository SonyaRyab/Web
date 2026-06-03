package main

import (
	// импорты для подключения к БД и другие
	"fmt"
    "log"
    "math/rand"
    "os"
    "time"

    "github.com/go-faker/faker/v4"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "lab4/internal/app/ds"
)

func main(){
	// читаем значения из окружения (или можешь захардкодить)
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "postgres")
	password := getenv("DB_PASSWORD", "postgres1234")
	dbname := getenv("DB_NAME", "postgres")
	
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	// чтобы rand работал нормально
	rand.Seed(time.Now().UnixNano())

	// заранее получаем пользователей из БД (user10, user11, user12)
	var users []ds.User
	if err := db.Where("login IN ?", []string{"user10", "user11", "user12"}).Find(&users).Error; err != nil {
		log.Fatalf("cannot load users: %v", err)
	}
	if len(users) == 0 {
		log.Fatalf("no users found for logins user10/user11/user12")
	}

	numEntries :=100000
	
	for i := 0; i < numEntries; i++ {
			newEntry := ds.Methane{
			Status: getRandomStatus(),
			AdminID: users[rand.Intn(len(users))].ID,
			DateCreate: getRandomDate(),
		}
		
		// генерируем рандомное имя
		if err := faker.FakeData(&newEntry.Name); err != nil {
			log.Printf("faker error: %v", err)
		}
		
		if err := db.Create(&newEntry).Error; err != nil {
			log.Printf("could not create application: %v", err)
		}
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getRandomStatus() string {
    statuses := []string{"draft", "formed", "completed", "rejected"}
    return statuses[rand.Intn(len(statuses))]
}

func getRandomDate() time.Time {
    // Случайная дата за последний год
    now := time.Now()
    delta := rand.Int63n(int64(365 * 24 * time.Hour))
    return now.Add(-time.Duration(delta))
}