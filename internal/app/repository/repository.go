package repository

import (
	"fmt"
	"lab2/internal/app/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db           *gorm.DB
	applications map[string]Experiment
	MinioClient  *MinioClient // Добавляем Minio клиент
	lastAppID    int          // для генерации ID
}

func New(dsn string, cfg *config.Config) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	// Подключаем MinIO
	minioClient, err := NewMinioClient(
		cfg.Minio.Endpoint,
		cfg.Minio.AccessKey,
		cfg.Minio.SecretKey,
		cfg.Minio.Bucket,
		cfg.Minio.UseSSL,
	)

	if err != nil {
		fmt.Printf("Предупреждение: не удалось подключиться к Minio: %v\n", err)
		minioClient = nil
	}

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db:          db,
		MinioClient: minioClient,
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

func (r *Repository) GetImageURL(filename string) string {
	if r.MinioClient != nil && filename != "" {
		return r.MinioClient.GetFileURL(filename)
	}
	// Если Minio недоступен, используем локальный путь
	return "/img/" + filename
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
