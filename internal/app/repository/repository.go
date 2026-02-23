package repository

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	orders       []Order
	db           *gorm.DB
	applications map[string]Experiment
	MinioClient  *MinioClient // Добавляем Minio клиент
	lastAppID    int          // для генерации ID
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db: db,
	}, nil
}

//func NewRepository() (*Repository, error) {
//	minioClient, err := NewMinioClient(
//		"localhost:9000", // эндпоинт Minio
//		"root",           // MINIO_ROOT_USER
//		"root12345",      // MINIO_ROOT_PASSWORD
//		"lab2-images",    // название bucket
//		false,            // useSSL = false для localhost
//	)
//
//	if err != nil {
//		log.Printf("Предупреждение: не удалось подключиться к Minio: %v", err)
//		minioClient = nil
//	}
//
//	repo := &Repository{
//		applications: make(map[string]Experiment),
//		MinioClient:  minioClient,
//		lastAppID:    0,
//	}
//
//	orders, err := repo.GetOrders()
//	if err != nil {
//		return nil, err
//	}
//	repo.orders = orders
//
//	repo.createTestApplication()
//
//	return repo, nil
//}

type Order struct {
	ID           int
	Title        string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
	Formula      string
	Image        string
	Video        string
	MolarMass    string
	Available    int
	Description  string
	Purity       float64
	DeliveryDate string
	BoilingPoint int
	MeltingPoint int
	Density      int
	Comment      string
	Coefficient  float64
}

type Experiment struct {
	ID          string
	CreatedAt   string
	Customer    string
	Status      string
	Items       []ExperimentItem
	TotalAmount int
	Result      string
}

// ExperimentItem - услуга в заявке
type ExperimentItem struct {
	ID          int
	OrderID     int
	Title       string
	Formula     string
	Image       string
	MolarMass   string
	Quantity    int
	Comment     string
	Coefficient float64
}

//func (r *Repository) GetOrders() ([]Order, error) {
//	orders := []Order{ // массив элементов из наших структур
//		{
//			ID:          1,
//			Title:       "Водород",
//			Formula:     "H2",
//			Coefficient: 4,
//			Image:       r.getImageURL("hydrogen.png"),
//			Video:       r.getVideoURL("sabatie_reaction.mp4"),
//			MolarMass:   "2",
//			Available:   100,
//			Description: "Водород - самый легкий химический элемент. " +
//				"В реакции Сабатье выступает в качестве восстановителя, " +
//				"реагирует с углекислым газом с образованием метана и воды. " +
//				"Хранится под высоким давлением в специальных баллонах.",
//			Purity:       99.999,
//			DeliveryDate: "15.04.2026",
//			BoilingPoint: -252,
//			MeltingPoint: -259,
//			Density:      89,
//			Comment:      "",
//		},
//		{
//			ID:          2,
//			Title:       "Углекислый газ",
//			Formula:     "CO2",
//			Coefficient: 1,
//			Image:       r.getImageURL("CO2.png"),
//			Video:       r.getVideoURL("sabatie_reaction.mp4"),
//			MolarMass:   "44",
//			Available:   500,
//			Description: "Углекислый газ - бесцветный газ со слегка кисловатым запахом. " +
//				"В реакции Сабатье является источником углерода для синтеза метана." +
//				"Поставляется в сжиженном состоянии в криогенных емкостях.",
//			Purity:       99.95,
//			DeliveryDate: "18.04.2025",
//			BoilingPoint: -78,
//			MeltingPoint: -56,
//			Density:      1977,
//			Comment:      "",
//		},
//		{
//			ID:          3,
//			Title:       "Никель",
//			Formula:     "Ni",
//			Coefficient: 0,
//			Image:       r.getImageURL("nickel.png"),
//			Video:       r.getVideoURL("sabatie_reaction.mp4"),
//			MolarMass:   "59",
//			Available:   50,
//			Description: "Никель - переходный металл серебристо-белого цвета. " +
//				"Используется в качестве катализатора в реакции Сабатье." +
//				"Поставляется в виде нанесенного катализатора на оксиде алюминия.",
//			Purity:       99.98,
//			DeliveryDate: "20.04.2025",
//			BoilingPoint: 2730,
//			MeltingPoint: 1455,
//			Density:      8908,
//			Comment:      "",
//		},
//		{
//			ID:          4,
//			Title:       "Энергия",
//			Formula:     "E",
//			Coefficient: 0,
//			Image:       r.getImageURL("energy.jpg"),
//			Video:       r.getVideoURL("sabatie_reaction.mp4"),
//			MolarMass:   "-",
//			Available:   10000,
//			Description: "Энергия для активации реакции Сабатье. " +
//				"Необходима для нагрева реактора до рабочей температуры 300-400°C. " +
//				"Может подаваться в виде электричества, пара или горячего теплоносителя. " +
//				"Расход энергии зависит от масштаба установки.",
//			Purity:       0,
//			DeliveryDate: "Постоянно",
//			BoilingPoint: 0,
//			MeltingPoint: 0,
//			Density:      0,
//			Comment:      "",
//		},
//	}
//	r.orders = orders
//	return orders, nil
//}

//func (r *Repository) GetOrder(id int) (Order, error) {
//	for _, order := range r.orders {
//		if order.ID == id {
//			return order, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
//		}
//	}
//	return Order{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
//}
//
//func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
//	orders, err := r.GetOrders()
//	if err != nil {
//		return []Order{}, err
//	}
//
//	var result []Order
//	for _, order := range orders {
//		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
//			result = append(result, order)
//		}
//	}
//
//	return result, nil
//}

func (r *Repository) getImageURL(filename string) string {
	if r.MinioClient != nil {
		return r.MinioClient.GetFileURL(filename)
	}
	// Если Minio недоступен, используем локальный путь
	return "/static/img/" + filename
}

func (r *Repository) getVideoURL(filename string) string {
	if r.MinioClient != nil {
		return r.MinioClient.GetFileURL(filename)
	}
	return ""
}

func (r *Repository) createTestApplication() {
	r.lastAppID++
	appID := fmt.Sprintf("EXP-%d", r.lastAppID)

	app := Experiment{
		ID:          appID,
		CreatedAt:   time.Now().Format("02.01.2006 15:04"),
		Customer:    "Иванов И.И.",
		Status:      "В обработке",
		Items:       []ExperimentItem{},
		TotalAmount: 0,
		Result:      "CH4 - 95% выход",
	}

	r.applications[appID] = app
}

// Добавление услуги в заявку
//func (r *Repository) AddToExperiment(appID string, orderID int) error {
//	app, exists := r.applications[appID]
//	if !exists {
//		return fmt.Errorf("заявка не найдена")
//	}
//
//	order, err := r.GetOrder(orderID)
//	if err != nil {
//		return err
//	}
//
//	for i, item := range app.Items {
//		if item.OrderID == orderID {
//			app.Items[i].Quantity++
//			app.TotalAmount++
//			r.applications[appID] = app
//			return nil
//		}
//	}
//
//	newItem := ExperimentItem{
//		ID:        len(app.Items) + 1,
//		OrderID:   order.ID,
//		Title:     order.Title,
//		Formula:   order.Formula,
//		Image:     order.Image,
//		MolarMass: order.MolarMass,
//		Quantity:  1,
//		Comment:   "",
//	}
//
//	app.Items = append(app.Items, newItem)
//	app.TotalAmount++
//	r.applications[appID] = app
//
//	return nil
//}

func (r *Repository) GetExperiment(appID string) (Experiment, error) {
	app, exists := r.applications[appID]
	if !exists {
		return Experiment{}, fmt.Errorf("заявка не найдена")
	}
	return app, nil
}

// Получение списка всех заявок
func (r *Repository) GetExperiments() ([]Experiment, error) {
	apps := make([]Experiment, 0, len(r.applications))
	for _, app := range r.applications {
		apps = append(apps, app)
	}
	return apps, nil
}

// Создание заявки
func (r *Repository) CreateExperiment(customer string) (string, error) {
	r.lastAppID++
	appID := fmt.Sprintf("APP-%d", r.lastAppID)

	app := Experiment{
		ID:          appID,
		CreatedAt:   time.Now().Format("02.01.2006 15:04"),
		Customer:    customer,
		Status:      "Новая",
		Items:       []ExperimentItem{},
		TotalAmount: 0,
		Result:      "Ожидание расчета",
	}

	r.applications[appID] = app
	return appID, nil
}

// Обновление поля м-м для услуги в заявке
func (r *Repository) UpdateExperimentItem(appID string, itemID int, quantity int, comment string) error {
	app, exists := r.applications[appID]
	if !exists {
		return fmt.Errorf("заявка не найдена")
	}

	for i, item := range app.Items {
		if item.ID == itemID {
			app.Items[i].Quantity = quantity
			app.Items[i].Comment = comment
			r.applications[appID] = app
			return nil
		}
	}

	return fmt.Errorf("элемент не найден")
}

func (r *Repository) RemoveFromExperiment(appID string, itemID int) error {
	app, exists := r.applications[appID]
	if !exists {
		return fmt.Errorf("заявка не найдена")
	}

	for i, item := range app.Items {
		if item.ID == itemID {
			// Удаляем элемент
			app.Items = append(app.Items[:i], app.Items[i+1:]...)
			app.TotalAmount--
			r.applications[appID] = app
			return nil
		}
	}

	return fmt.Errorf("элемент не найден в заявке")
}

func (r *Repository) ClearExperiment(appID string) error {
	app, exists := r.applications[appID]
	if !exists {
		return fmt.Errorf("заявка не найдена")
	}

	app.Items = []ExperimentItem{}
	app.TotalAmount = 0
	r.applications[appID] = app
	return nil
}
