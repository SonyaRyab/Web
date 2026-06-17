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
func (r Repository) GetReagents(search string) ([]ds.Reagent, error) {
	var reagents []ds.Reagent

	query := r.db.Where("isdeleted = ?", false)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR formula ILIKE ?", like, like)
	}

	err := query.Order("id asc").Find(&reagents).Error
	return reagents, err
}

func (r Repository) GetReagentsPaged(search string, page int, limit int) ([]ds.Reagent, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 24
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	var items []ds.Reagent
	var total int64

	query := r.db.Model(&ds.Reagent{}).Where("isdeleted = ?", false)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR formula ILIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id asc").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r Repository) GetAllReagentIDs() ([]uint, error) {
	var ids []uint
	err := r.db.Model(&ds.Reagent{}).Where("isdeleted = ?", false).Pluck("id", &ids).Error
	return ids, err
}

// GetReagentByID получает один реагент по ID
func (r *Repository) GetReagentByID(id uint) (ds.Reagent, error) {
	var reagent ds.Reagent
	err := r.db.Where("id = ? AND isdeleted = ?", id, false).First(&reagent).Error
	return reagent, err
}

// UploadReagentFiles загружает изображение и видео в MinIO
func (r Repository) UploadReagentFiles(reagentID uint, imgHeader, videoHeader *multipart.FileHeader) (string, string, error) {
	ctx := context.Background()
	baseName := fmt.Sprintf("reagent-%d", reagentID)

	var imgURL, videoURL string
	var err error

	if imgHeader != nil {
		imgURL, err = r.UploadFileToMinIO(ctx, baseName+"-img", imgHeader)
		if err != nil {
			return "", "", err
		}
	}

	if videoHeader != nil {
		videoURL, err = r.UploadFileToMinIO(ctx, baseName+"-video", videoHeader)
		if err != nil {
			return "", "", err
		}
	}

	return imgURL, videoURL, nil
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

func (r Repository) UploadFileToMinIO(ctx context.Context, filename string, header *multipart.FileHeader) (string, error) {
	file, err := header.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	_, err = r.minio.PutObject(
		ctx,
		r.minio_bucket_name,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s/%s", r.minio.EndpointURL().Host, r.minio_bucket_name, filename), nil
}

// DeleteFileFromMinIO удаляет файл из MinIO 
func (r *Repository) DeleteFileFromMinIO(filename string) error {
	return r.minio.RemoveObject(context.Background(), r.minio_bucket_name, filename, minio.RemoveObjectOptions{})
}

// GetMinioBucketName возвращает имя бакета
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
