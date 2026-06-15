package repository

import (
	"context"
	"fmt"
	"lab4/internal/app/ds"
	"mime/multipart"
	"net/http"

	"github.com/minio/minio-go/v7"
)

// GetReagents получает список реагентов с фильтрацией
func (r *Repository) GetReagents(search string) ([]ds.Reagent, error) {
	var reagents []ds.Reagent
	query := r.db.Where("is_deleted = ?", false)

	if search != "" {
		query = query.Where("name ILIKE ? OR formula ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Find(&reagents).Error
	if err != nil {
		return nil, err
	}

	return reagents, nil
}

// GetReagentByID получает один реагент по ID
func (r *Repository) GetReagentByID(id uint) (ds.Reagent, error) {
	var reagent ds.Reagent
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&reagent).Error
	if err != nil {
		return ds.Reagent{}, err
	}
	return reagent, nil
}

// AddReagent добавляет новый реагент
func (r *Repository) AddReagent(reagent *ds.Reagent) error {
	return r.db.Create(reagent).Error
}

// UpdateReagent обновляет реагент
func (r *Repository) UpdateReagent(id uint, updates map[string]interface{}) error {
	return r.db.Model(&ds.Reagent{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteReagent мягкое удаление
func (r *Repository) DeleteReagent(id uint) error {
	return r.db.Model(&ds.Reagent{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// UploadReagentFiles загружает изображение и видео в MinIO
func (r *Repository) UploadReagentFiles(reagent_id uint, imgHeader, videoHeader *multipart.FileHeader) (imgURL, videoURL string, err error) {
	ctx := context.Background()

	// Генерируем имена файлов на латинице
	baseName := fmt.Sprintf("reagent_%d", reagent_id)

	// Загрузка изображения
	if imgHeader != nil {
		imgName := baseName + "_img"
		imgURL, err = r.UploadFileToMinIO(ctx, imgName, imgHeader)
		if err != nil {
			return "", "", fmt.Errorf("ошибка загрузки изображения: %w", err)
		}
	}

	// Загрузка видео
	if videoHeader != nil {
		videoName := baseName + "_video"
		videoURL, err = r.UploadFileToMinIO(ctx, videoName, videoHeader)
		if err != nil {
			// Откат изображения если видео не загрузилось
			if imgURL != "" {
				r.minio.RemoveObject(ctx, r.minio_bucket_name, baseName+"_img", minio.RemoveObjectOptions{})
			}
			return "", "", fmt.Errorf("ошибка загрузки видео: %w", err)
		}
	}

	return imgURL, videoURL, nil
}

func (r *Repository) UploadFileToMinIO(ctx context.Context, filename string, header *multipart.FileHeader) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Определяем Content-Type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	_, err = r.minio.PutObject(
		ctx,
		r.minio_bucket_name,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s/%s/%s", r.minio.EndpointURL().Host, r.minio_bucket_name, filename), nil
}

// DeleteFileFromMinIO удаляет файл из MinIO (ПУБЛИЧНЫЙ МЕТОД)
func (r *Repository) DeleteFileFromMinIO(filename string) error {
	return r.minio.RemoveObject(context.Background(), r.minio_bucket_name, filename, minio.RemoveObjectOptions{})
}

// GetMinioBucketName возвращает имя бакета (ПУБЛИЧНЫЙ МЕТОД)
func (r *Repository) GetMinioBucketName() string {
	return r.minio_bucket_name
}

// загрузка одного файла в MinIO (вспомогательный для handler)
func (r *Repository) UploadSingleFile(filename string, header *multipart.FileHeader) (string, error) {
	return r.UploadFileToMinIO(context.Background(), filename, header)
}

// удаление файла из MinIO
func (r *Repository) DeleteFile(filename string) error {
	return r.minio.RemoveObject(context.Background(), r.minio_bucket_name, filename, minio.RemoveObjectOptions{})
}
