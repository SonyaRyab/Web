package repository

import (
    "lab4/internal/app/ds"
)

func (r *Repository) GetMethanesPaged(limit, offset int) ([]ds.Methane, int64, error) {
    var entities []ds.Methane
    var total int64

    query := r.db.Model(&ds.Methane{})

    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := query.Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
        return nil, 0, err
    }

    return entities, total, nil
}