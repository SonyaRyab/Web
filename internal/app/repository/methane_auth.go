package repository

import (
	"lab4/internal/app/ds"
	"lab4/internal/app/role"
)

func (r *Repository) GetMethanesForUser(userID uint, userRole role.Role) ([]ds.Methane, error) {
	var items []ds.Methane

	query := r.db.Preload("Moderator").Preload("Admin")

	switch userRole {
    case role.Admin:
        err := query.Where("admin_id = ?", userID).Find(&items).Error
        return items, err
    case role.Moderator:
        err := query.Find(&items).Error
        return items, err
    default:
        return items, nil
    }
}
