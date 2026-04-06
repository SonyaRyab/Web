package main

import (
	"lab3/internal/app/ds"
	"lab3/internal/app/dsn"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Methane{},
		&ds.MethaneReagent{},
		&ds.Reagent{},
		&ds.User{},
	)
	if err != nil {
		panic("cant migrate db")
	}

	var reagentCount int64
	db.Model(&ds.Reagent{}).Count(&reagentCount)
	if reagentCount == 0 {
		reagents := []ds.Reagent{
			{
				Name:        "Водород",
				Formula:     "H₂",
				MolarMass:   2.02,
				Description: "Восстановитель в реакции Сабатье",
			},
			{
				Name:        "Углекислый газ",
				Formula:     "CO₂",
				MolarMass:   44.01,
				Description: "Исходное вещество",
			},
			{
				Name:        "Никель",
				Formula:     "Ni",
				MolarMass:   58.69,
				Description: "Катализатор реакции",
			},
		}
		for _, r := range reagents {
			db.Create(&r)
		}
	}

	var methaneCount int64
	db.Model(&ds.Methane{}).Count(&methaneCount)
	if methaneCount == 0 {
		defaultMethane := ds.Methane{
			Name:        "Тестовый эксперимент",
			Status:      "черновик",
			DateCreate:  time.Now(),
			AdminID:     1,
			Temperature: 350.0,
		}
		db.Create(&defaultMethane)
	}

	// Check if there is a default user row and add one if it's missing
	var count int64

	db.Model(&ds.User{}).Where("login = ?", "test").Count(&count)
	if count == 0 {
		defaultUser := ds.User{Username: "test", Login: "test", Password: "test123", IsModerator: false}

		err := db.Create(&defaultUser).Error

		if err != nil {
			panic("error creating default user")
		}
	}
}
