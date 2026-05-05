package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"

	"lab4/internal/app/ds"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (r *Repository) GetMethaneByID(id uint) (*ds.Methane, error) {
	var methane ds.Methane
	err := r.db.Preload("Professor").Preload("Researcher").Where("id = ?", id).First(&methane).Error
	if err != nil {
		return nil, err
	}
	return &methane, nil
}

func (r *Repository) FormMethaneByOwner(id uint, userID uint, updates map[string]interface{}) error {
	var methane ds.Methane
	err := r.db.Where("id = ?", id).First(&methane).Error
	if err != nil {
		return err
	}

	if methane.ProfessorID == nil || *methane.ProfessorID != userID {
		return errors.New("forbidden")
	}

	if methane.Status != "черновик" {
		return errors.New("only draft methane can be formed")
	}

	delete(updates, "id")
	delete(updates, "professor_id")
	delete(updates, "researcher_id")
	delete(updates, "status")
	delete(updates, "date_create")
	delete(updates, "date_form")
	delete(updates, "date_finish")

	now := time.Now()
	updates["status"] = "сформирована"
	updates["date_form"] = &now

	return r.db.Model(&ds.Methane{}).Where("id = ?", id).Updates(updates).Error
}

func (r *Repository) CompleteMethane(id uint, currentUserID uint, status string) error {
	if status != "завершена" && status != "отклонена" {
		return fmt.Errorf("invalid completion status")
	}

	var methane ds.Methane
	err := r.db.Where("id = ?", id).First(&methane).Error
	if err != nil {
		return err
	}

	if methane.Status != "сформирована" {
		return fmt.Errorf("only formed methane can be completed")
	}

	if methane.ResearcherID == currentUserID {
        return fmt.Errorf("forbidden")
    }

	now := time.Now()
	return r.db.Model(&ds.Methane{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"professor_id": currentUserID,
		"date_finish":  &now,
	}).Error
}

// @Security BearerAuth
func (r *Repository) GetMethanes() ([]ds.Methane, error) {
	var methanes []ds.Methane
	err := r.db.Find(&methanes).Error
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	if err != nil {
		return nil, err
	}
	if len(methanes) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return methanes, nil
}

func (r *Repository) AddMethane(methane *ds.Methane) error {
	err := r.db.Model(&ds.Methane{}).Create(methane).Error
	if err != nil {
		return fmt.Errorf("ошибка при добавлении акции: %w", err)
	}

	return nil
}

func (r *Repository) GetMethane(id int) (ds.Methane, error) {
	methane := ds.Methane{}
	err := r.db.Preload("Professor").Where("id = ?", id).First(&methane).Error
	if err != nil {
		return ds.Methane{}, err
	}
	return methane, nil
}

func (r *Repository) DeleteMethane(methaneID uint) error {
	// Обратите внимание: жёсткое удаление, не логическое
	err := r.db.Delete(&ds.Methane{}, methaneID).Error
	if err != nil {
		return fmt.Errorf("ошибка при удалении команды с id %d: %w", methaneID, err)
	}

	return nil
}

func (r *Repository) ModifyMethane(id uint, methane *ds.Methane) error {

	var old_methane ds.Methane

	err := r.db.Model(&ds.Methane{}).Where("id = ?", id).First(&old_methane).Error

	if err != nil {
		return fmt.Errorf("не удалось найти команду с id %d: %w", id, err)
	}

	// Попытки изменения ID акции нужно предотвращать
	// Создаем структуру для обновления без UUID
	updateData := map[string]interface{}{
		"name":        methane.Name,
		"status":      methane.Status,
	}

	err = r.db.Model(&ds.Methane{}).Where("id = ?", id).Updates(&updateData).Error

	if err != nil {
		return fmt.Errorf("ошибка при обновлении акции с id %d: %w", id, err)
	}

	return nil
}

func (r *Repository) AddOrReplaceMethaneImage(methaneID uint, header *multipart.FileHeader) error {
	// Название будущего файла - <номер акции>.png

	filename := strconv.FormatUint(uint64(methaneID), 10)

	// Расширение .png будет даже у картинок, исходно его не имевших.
	// Это не совсем хорошо, но отображаться всё будет.

	// Открываем файл
	file, err := header.Open()
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	// Как бы мы ни вышли из этой функции, файл обязательно будет закрыт
	defer file.Close()

	// header уже содержит заголовок с типом файла, но для пущей уверенности
	// мы получим Content-Type на базе реального содержимого файла

	// Тип полученного файла хранится в его первых 512 байтах
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	// Определяем тип файла по его содержимому
	contentType := http.DetectContentType(buffer)

	// Позиция чтения файла возвращается к исходной
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("ошибка перемещения по файловому потоку: %w", err)
	}

	ctx := context.Background()

	_, err = r.minio.PutObject(
		ctx,
		r.minio_bucket_name,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		})

	if err != nil {
		return fmt.Errorf("не удалось добавить объект в хранилище minio: %w", err)
	}

	// Эндпоинт здесь должен совпадать с эндпоинтом MinIO в конфигурации
	err = r.db.Model(&ds.Methane{}).Where("id = ?", methaneID).UpdateColumn("pic", "http://127.0.0.1:9000/"+r.minio_bucket_name+"/"+filename).Error
	if err != nil {
		// Если не удалось сохранить в БД, удаляем из MinIO
		r.minio.RemoveObject(ctx, r.minio_bucket_name, filename, minio.RemoveObjectOptions{})
		return fmt.Errorf("ошибка сохранения пути к изображению: %w", err)
	}

	return nil
}

// GetMethanesWithFilter получает заявки с фильтрацией по дате и статусу
func (r *Repository) GetMethanesWithFilter(status string, dateFrom, dateTo string) ([]ds.MethaneListSerializer, error) {
	var results []ds.MethaneListSerializer

	query := r.db.Model(&ds.Methane{}).
		Select("methanes.*, COUNT(methane_reagents.id) as reagent_count").
		Joins("LEFT JOIN methane_reagents ON methanes.id = methane_reagents.methane_id").
		Where("methanes.status != ?", "удалена").
		Group("methanes.id")

	if status != "" {
		query = query.Where("methanes.status = ?", status)
	}
	if dateFrom != "" {
		query = query.Where("methanes.date_form >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("methanes.date_form <= ?", dateTo)
	}

	err := query.Find(&results).Error
	return results, err
}

// GetDraftMethane получает черновик текущего пользователя
func (r *Repository) GetDraftMethane(userID uint) (*ds.Methane, error) {
	var methane ds.Methane
	err := r.db.Where("researcher_id = ? AND status = ?", userID, "черновик").First(&methane).Error
	if err != nil {
		return nil, err
	}
	return &methane, nil
}

// CreateDraftMethane создаёт пустую заявку-черновик
func (r *Repository) CreateDraftMethane(userID uint) (*ds.Methane, error) {
	now := time.Now()
	methane := ds.Methane{
		Name:       "Новый эксперимент",
		Status:     "черновик",
		DateCreate: now,
		ResearcherID: userID,
	}
	err := r.db.Create(&methane).Error
	return &methane, err
}

// FormMethane формирует заявку (меняет статус, вычисляет поля)
// func (r *Repository) FormMethane(id uint, updates map[string]interface{}) error {
// 	updates["status"] = "сформирована"
// 	updates["date_form"] = time.Now()
// 	return r.db.Model(&ds.Methane{}).Where("id = ?", id).Updates(updates).Error
// }

// SoftDeleteMethane мягкое удаление (установка статуса)
func (r *Repository) SoftDeleteMethane(id uint) error {
	return r.db.Model(&ds.Methane{}).Where("id = ?", id).Update("status", "удалена").Error
}

// GetMethaneWithReagents получает заявку со всеми реагентами
func (r *Repository) GetMethaneWithReagents(id uint) (ds.Methane, error) {
	var methane ds.Methane
	err := r.db.Preload("Professor").Preload("Researcher").Preload("Reagents").Preload("Reagents.Reagent").
		Where("id = ?", id).First(&methane).Error
	return methane, err
}
