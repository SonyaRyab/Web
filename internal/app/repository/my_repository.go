package repository

import (
    "lab4/internal/app/ds"
)

func (r Repository) GetMethanesPaged(status string, limit, offset int) ([]ds.Methane, int64, error) {
	var entities []ds.Methane
	var total int64

	query := r.db.Model(&ds.Methane{}).
		Preload("Moderator").
		Preload("Admin")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("date_create DESC").
		Offset(offset).
		Limit(limit).
		Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}