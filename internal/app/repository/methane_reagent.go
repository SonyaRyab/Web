package repository

import (
	"lab3/internal/app/ds"
)

// AddReagentToMethane добавляет реагент в заявку (создаёт заявку если нужно)
func (r *Repository) AddReagentToMethane(methaneID, reagentID uint, volume float64, orderNum int) error {
	// Проверяем существование связи
	var existing ds.MethaneReagent
	err := r.db.Where("methane_id = ? AND reagent_id = ?", methaneID, reagentID).First(&existing).Error

	if err == nil {
		// Обновляем количество если уже есть
		existing.Volume += volume
		return r.db.Save(&existing).Error
	}

	// Создаём новую связь - простая структура без связей
	mr := ds.MethaneReagent{
		MethaneID: methaneID,
		ReagentID: reagentID,
		Volume:  volume,
	}

	return r.db.Create(&mr).Error
}

// UpdateMethaneReagent изменяет количество/порядок реагента в заявке
func (r *Repository) UpdateMethaneReagent(methaneID, reagentID uint, volume *float64, orderNum *int) error {
	updates := map[string]interface{}{}
	if volume != nil {
		updates["volume"] = *volume
	}
	if orderNum != nil {
		updates["order_num"] = *orderNum
	}

	return r.db.Model(&ds.MethaneReagent{}).
		Where("methane_id = ? AND reagent_id = ?", methaneID, reagentID).
		Updates(updates).Error
}

// RemoveReagentFromMethane удаляет реагент из заявки
func (r *Repository) RemoveReagentFromMethane(methaneID, reagentID uint) error {
	return r.db.Where("methane_id = ? AND reagent_id = ?", methaneID, reagentID).
		Delete(&ds.MethaneReagent{}).Error
}

// GetMethaneReagents получает все реагенты заявки с деталями
func (r *Repository) GetMethaneReagents(methaneID uint) ([]ds.MethaneReagent, error) {
	var items []ds.MethaneReagent
	err := r.db.Preload("Reagent").Where("methane_id = ?", methaneID).Find(&items).Error
	return items, err
}

// CountReagentsInMethane считает количество реагентов в заявке
func (r *Repository) CountReagentsInMethane(methaneID uint) (int64, error) {
	var count int64
	err := r.db.Model(&ds.MethaneReagent{}).Where("methane_id = ?", methaneID).Count(&count).Error
	return count, err
}

func (r *Repository) GetMethaneReagentsWithDetails(methaneID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := r.db.Raw(`
		SELECT mr.id, mr.volume, r.id as reagent_id, r.name, r.img, r.molar_mass
		FROM methane_reagents mr
		JOIN reagents r ON mr.reagent_id = r.id
		WHERE mr.methane_id = ?
	`, methaneID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uint
		var volume float64
		var reagentID uint
		var name, img string
		var molarMass float64

		rows.Scan(&id, &volume, &reagentID, &name, &img, &molarMass)

		results = append(results, map[string]interface{}{
			"id":         id,
			"volume":   volume,
			"reagent_id": reagentID,
			"name":       name,
			"img":        img,
			"molar_mass": molarMass,
		})
	}

	return results, nil
}
