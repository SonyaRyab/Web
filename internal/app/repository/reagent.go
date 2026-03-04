package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"lab2/internal/app/ds"
	"time"

	"github.com/sirupsen/logrus"
)

func (r *Repository) GetAllReagents() ([]ds.Reagent, error) {
	var reagents []ds.Reagent
	err := r.db.Where("is_deleted = false").Find(&reagents).Error // добавили условие
	if err != nil {
		return nil, err
	}
	return reagents, nil
}

func (r *Repository) GetReagentByID(id int) (*ds.Reagent, error) {
	var reagent ds.Reagent
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&reagent).Error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Возвращаем nil, если записи нет
		}
		return nil, err
	}
	return &reagent, nil
}

func (r *Repository) SearchReagentsByName(name string) ([]ds.Reagent, error) {
	var reagents []ds.Reagent
	err := r.db.Where("name ILIKE ? and is_deleted = ?", "%"+name+"%", false).Find(&reagents).Error // добавили условие
	if err != nil {
		return nil, err
	}
	return reagents, nil
}

// GetCartCount для получения количества услуг в заявке (чатов в сообщении в моем случае)
func (r *Repository) GetCartCount() int64 {
	var experimentID uint
	var count int64
	adminID := 1
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	err := r.db.Model(&ds.Methane{}).
		Where("admin_id = ? AND status = ?", adminID, "черновик").
		Select("id").First(&experimentID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.MethaneReagent{}).
		Where("experiment_id = ?", experimentID).
		Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records:", err)
		return 0
	}
	return count
}

func (r *Repository) DeleteReagent(reagentID uint) error {
	query := "UPDATE reagents SET is_deleted = true WHERE id = $1"
	result := r.db.Exec(query, reagentID)

	if result.Error != nil {
		return fmt.Errorf("ошибка при удалении реагента с id %d: %w", reagentID, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("реагент с id %d не найден", reagentID)
	}

	return nil
}

// GetDraftExperimentID - получить ID черновика или создать новый
func (r *Repository) GetDraftMethaneID(adminID uint) (uint, error) {
	var exp ds.Methane
	err := r.db.Where("admin_id = ? AND status = ?", adminID, "черновик").First(&exp).Error

	if err != nil {
		// Создаём новый черновик
		exp = ds.Methane{
			Status:     "черновик",
			DateCreate: time.Now(),
			AdminID:    adminID,
		}
		err = r.db.Create(&exp).Error
		if err != nil {
			return 0, err
		}
	}

	return exp.ID, nil
}

// AddReagentToExperiment - добавить реагент в заявку (ORM)
func (r *Repository) AddReagentToMethane(methaneID, reagentID uint, quantity float64) error {
	// Проверяем, есть ли уже такой реагент
	var existing ds.MethaneReagent
	err := r.db.Where("methane_id = ? AND reagent_id = ?", methaneID, reagentID).First(&existing).Error

	if err == nil {
		// Уже есть — обновляем количество
		existing.Quantity += quantity
		return r.db.Save(&existing).Error
	}

	// Новая запись
	item := ds.MethaneReagent{
		Methane_id: methaneID,
		Reagent_id: reagentID,
		Quantity:   quantity,
		//OrderNum:  1, // можно добавить логику подсчёта
	}

	return r.db.Create(&item).Error
}

// GetExperimentWithReagents - получить заявку со всеми реагентами
func (r *Repository) GetMethaneWithReagents(expID uint) (*ds.Methane, []ds.MethaneReagent, error) {
	var exp ds.Methane
	err := r.db.First(&exp, expID).Error
	if err != nil {
		return nil, nil, err
	}

	var items []ds.MethaneReagent
	err = r.db.Where("methane_id = ?", expID).Preload("Reagent").Find(&items).Error
	if err != nil {
		return nil, nil, err
	}

	for i := range items {
		items[i].Reagent.Img = r.GetImageURL(items[i].Reagent.Img)
	}

	return &exp, items, nil
}

// CalculateMethaneYield - рассчитать выход метана по реакции Сабатье
// CO2 + 4H2 → CH4 + 2H2O
func (r *Repository) CalculateMethaneYield(methaneID uint) (float64, error) {
	var items []ds.MethaneReagent
	err := r.db.Where("methane_id = ?", methaneID).Preload("Reagent").Find(&items).Error
	if err != nil {
		return 0, err
	}

	var h2Moles, co2Moles float64

	for _, item := range items {
		if item.Reagent.Formula == "H2" {
			h2Moles = item.Quantity / item.Reagent.MolarMass // моль H2
		}
		if item.Reagent.Formula == "CO2" {
			co2Moles = item.Quantity / item.Reagent.MolarMass // моль CO2
		}
	}

	if co2Moles == 0 {
		return 0, fmt.Errorf("нет CO2 в реакции")
	}

	// Стехиометрия: 4 моль H2 на 1 моль CO2
	requiredH2 := co2Moles * 4
	actualH2 := h2Moles

	// Определяем ограничивающий реагент
	var methaneMoles float64
	if actualH2 >= requiredH2 {
		// CO2 ограничивает
		methaneMoles = co2Moles
	} else {
		// H2 ограничивает
		methaneMoles = actualH2 / 4
	}

	// Теоретический выход метана (в % от стехиометрического)
	// Упрощённая формула: (фактический CH4 / теоретический CH4) * 100
	yield := (methaneMoles / co2Moles) * 100

	if yield > 100 {
		yield = 100
	}

	return yield, nil
}

func (r *Repository) DeleteMethane(expID uint) error {
	query := "UPDATE methanes SET status = 'удалён', date_update = NOW() WHERE id = $1 AND status = 'черновик'"
	result := r.db.Exec(query, expID)

	if result.Error != nil {
		return fmt.Errorf("ошибка при удалении заявки: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("заявка не найдена или не в статусе черновик")
	}

	return nil
}
