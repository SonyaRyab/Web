package repository

import (
	"errors"
	"fmt"
	"time"

	"lab4/internal/app/ds"
)

func (r Repository) GetMethaneByID(id uint) (ds.Methane, error) {
	var methane ds.Methane

	err := r.db.
		// Preload("Moderator").
		// Preload("Admin").
		Preload("Reagents").
		Preload("Reagents.Reagent").
		Where("id = ?", id).
		First(&methane).Error
	if err != nil {
		return ds.Methane{}, err
	}

	return methane, nil
}

func (r Repository) GetDraftMethane(userID uint) (ds.Methane, error) {
	var methane ds.Methane

	err := r.db.
		Preload("Moderator").
		Preload("Admin").
		Preload("Reagents").
		Preload("Reagents.Reagent").
		Where("admin_id = ? AND status = ?", userID, "draft").
		First(&methane).Error
	if err != nil {
		return ds.Methane{}, err
	}

	return methane, nil
}

func (r Repository) CreateDraftMethane(userID uint) (ds.Methane, error) {
	now := time.Now()

	methane := ds.Methane{
		Name:       "",
		Status:     "draft",
		DateCreate: now,
		AdminID:    userID,
	}

	err := r.db.Create(&methane).Error
	if err != nil {
		return ds.Methane{}, err
	}

	return methane, nil
}

func (r Repository) FormMethaneByOwner(id uint, userID uint, updates map[string]interface{}) error {
	var methane ds.Methane

	err := r.db.Where("id = ?", id).First(&methane).Error
	if err != nil {
		return err
	}

	if methane.AdminID != userID {
		return errors.New("forbidden")
	}

	if methane.Status != "draft" {
		return errors.New("only draft methane can be formed")
	}

	delete(updates, "id")
	delete(updates, "moderator_id")
	delete(updates, "admin_id")
	delete(updates, "date_create")
	delete(updates, "date_form")
	delete(updates, "date_finish")

	now := time.Now()
	updates["status"] = "formed"
	updates["date_form"] = now

	return r.db.Model(&ds.Methane{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r Repository) FormMethane(id uint, updates map[string]interface{}) error {
	updates["status"] = "formed"
	updates["date_form"] = time.Now()

	return r.db.Model(&ds.Methane{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r Repository) SoftDeleteMethane(id uint) error {
	return r.db.Model(&ds.Methane{}).
		Where("id = ?", id).
		Update("status", "deleted").Error
}

func (r Repository) GetMethaneWithReagents(id uint) (ds.Methane, error) {
	var methane ds.Methane

	err := r.db.
		Preload("Moderator").
		Preload("Admin").
		Preload("Reagents").
		Preload("Reagents.Reagent").
		Where("id = ?", id).
		First(&methane).Error
	if err != nil {
		return ds.Methane{}, err
	}

	return methane, nil
}

func (r Repository) CompleteMethane(id uint, currentUserID uint, status string) error {
	if status != "completed" && status != "rejected" && status != "formed" {
		return fmt.Errorf("invalid completion status")
	}

	var methane ds.Methane
	err := r.db.Where("id = ?", id).First(&methane).Error
	if err != nil {
		return err
	}

	if status == "formed" {
		if methane.Status != "draft" {
			return fmt.Errorf("only draft methane can be formed")
		}

		return r.db.Model(&ds.Methane{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"status":    "formed",
				"date_form": time.Now(),
			}).Error
	}

	if methane.Status != "formed" {
		return fmt.Errorf("only formed methane can be completed")
	}

	now := time.Now()
	return r.db.Model(&ds.Methane{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"moderator_id": currentUserID,
			"date_finish":  now,
		}).Error
}

func (r Repository) AddReagentToDraft(methaneID uint, userID uint, reagent_id uint, volume int) error {
    var methane ds.Methane
    err := r.db.Where("id = ?", methaneID).First(&methane).Error
    if err != nil {
        return err
    }

    if methane.AdminID != userID {
        return errors.New("forbidden")
    }

    if methane.Status != "draft" {
        return errors.New("only draft methane can be modified")
    }

    var reagent ds.Reagent
    err = r.db.Where("id = ?", reagent_id).First(&reagent).Error
    if err != nil {
        return errors.New("reagent not found")
    }

    var existing ds.MethaneReagent
    err = r.db.Where("methane_id = ? AND reagent_id = ?", methaneID, reagent_id).First(&existing).Error
    
    if err == nil {
        return r.db.Model(&ds.MethaneReagent{}).
            Where("methane_id = ? AND reagent_id = ?", methaneID, reagent_id).
            Update("volume", existing.Volume + float64(volume)).Error
    }

    methaneReagent := ds.MethaneReagent{
        MethaneID: methaneID,
        Reagent_id: reagent_id,
        Volume:  float64(volume),
    }

    // return r.db.Create(&methaneReagent).Error
	return r.db.Omit("Reagent").Create(&methaneReagent).Error
}

func (r Repository) UpdateReagentVolume(methaneID uint, userID uint, reagent_id uint, volume int) error {
    var methane ds.Methane
    err := r.db.Where("id = ?", methaneID).First(&methane).Error
    if err != nil {
        return err
    }

    if methane.AdminID != userID {
        return errors.New("forbidden")
    }

    if methane.Status != "draft" {
        return errors.New("only draft methane can be modified")
    }

    return r.db.Model(&ds.MethaneReagent{}).
        Where("methane_id = ? AND reagent_id = ?", methaneID, reagent_id).
        Update("volume", volume).Error
}

func (r Repository) RemoveReagentFromDraft(methaneID uint, userID uint, reagent_id uint) error {
    var methane ds.Methane
    err := r.db.Where("id = ?", methaneID).First(&methane).Error
    if err != nil {
        return err
    }

    if methane.AdminID != userID {
        return errors.New("forbidden")
    }

    if methane.Status != "draft" {
        return errors.New("only draft methane can be modified")
    }

    return r.db.Where("methane_id = ? AND reagent_id = ?", methaneID, reagent_id).
        Delete(&ds.MethaneReagent{}).Error
}