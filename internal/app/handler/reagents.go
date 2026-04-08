package handler

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"lab3/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetReagentsAPI(ctx *gin.Context) {
	search := ctx.Query("search")

	reagents, err := h.Repository.GetReagents(search)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		// "status": "success",
		"data":   reagents,
	})
}

func (h *Handler) GetReagentByIDAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reagent, err := h.Repository.GetReagentByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		// "status": "success",
		"data":   reagent,
	})
}

func (h *Handler) AddReagentAPI(ctx *gin.Context) {
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reagent := ds.Reagent{
		Name:        ctx.Request.FormValue("name"),
		Formula:     ctx.Request.FormValue("formula"),
		Description: ctx.Request.FormValue("description"),
	}

	// Парсим молярную массу
	if m := ctx.Request.FormValue("molar_mass"); m != "" {
		mm, err := strconv.ParseFloat(m, 64)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		reagent.MolarMass = mm
	}

	// Создаём реагент сначала без файлов
	if err := h.Repository.AddReagent(&reagent); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Обработка файлов
	imgHeader, _ := ctx.FormFile("img")
	videoHeader, _ := ctx.FormFile("video")

	if imgHeader != nil || videoHeader != nil {
		imgURL, videoURL, err := h.Repository.UploadReagentFiles(reagent.ID, imgHeader, videoHeader)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		updates := map[string]interface{}{}
		if imgURL != "" {
			updates["img"] = imgURL
		}
		if videoURL != "" {
			updates["video"] = videoURL
		}

		h.Repository.UpdateReagent(reagent.ID, updates)
		reagent.Img = imgURL
		reagent.Video = videoURL
	}

	ctx.JSON(http.StatusCreated, gin.H{
		// "status":  "success",
		"data":    reagent,
		"message": "реагент добавлен",
	})
}

func (h *Handler) UpdateReagentAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректный ID"))
		return
	}

	// текущий реагент
	currentReagent, err := h.Repository.GetReagentByID(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, fmt.Errorf("реагент не найден"))
		return
	}

	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	updates := map[string]interface{}{}

	if name := ctx.Request.FormValue("name"); name != "" {
		updates["name"] = name
	}

	if formula := ctx.Request.FormValue("formula"); formula != "" {
		updates["formula"] = formula
	}

	if m := ctx.Request.FormValue("molar_mass"); m != "" {
		mm, err := strconv.ParseFloat(m, 64)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректная молярная масса"))
			return
		}
		updates["molar_mass"] = mm
	}

	deleteImg := ctx.Request.FormValue("delete_img") == "true"
	deleteVideo := ctx.Request.FormValue("delete_video") == "true"
	imgHeader, imgErr := ctx.FormFile("img")
	videoHeader, videoErr := ctx.FormFile("video")

	if imgErr == nil {
		code, err := validateFileUpload(imgHeader)
		if err != nil {
			h.errorHandler(ctx, code, fmt.Errorf("изображение: %w", err))
			return
		}
	}

	if videoErr == nil {
		if videoHeader.Size > 50<<20 { // 50MB max для видео
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("видео слишком большое (макс 50MB)"))
			return
		}

		file, _ := videoHeader.Open()
		buffer := make([]byte, 512)
		file.Read(buffer)
		file.Close()

		contentType := http.DetectContentType(buffer)
		validVideoTypes := []string{"video/mp4", "video/webm", "video/ogg", "video/quicktime"}
		isValid := false
		for _, t := range validVideoTypes {
			if contentType == t {
				isValid = true
				break
			}
		}
		if !isValid {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("недопустимый формат видео"))
			return
		}
	}

	var newImgURL, newVideoURL string
	var uploadedImg, uploadedVideo bool

	if imgErr == nil && !deleteImg {
		if currentReagent.Img != "" {
			oldFilename := fmt.Sprintf("reagent_%d_img", id)
			h.Repository.DeleteFileFromMinIO(oldFilename)
		}

		imgName := fmt.Sprintf("reagent_%d_img", id)
		newImgURL, err = h.Repository.UploadFileToMinIO(context.Background(), imgName, imgHeader)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("ошибка загрузки изображения: %w", err))
			return
		}
		updates["img"] = newImgURL
		uploadedImg = true
	} else if deleteImg {
		// Удаляем изображение
		if currentReagent.Img != "" {
			oldFilename := fmt.Sprintf("reagent_%d_img", id)
			h.Repository.DeleteFileFromMinIO(oldFilename)
		}
		updates["img"] = ""
	}

	if videoErr == nil && !deleteVideo {
		// Удаляем старое видео если было
		if currentReagent.Video != "" {
			oldFilename := fmt.Sprintf("reagent_%d_video", id)
			h.Repository.DeleteFileFromMinIO(oldFilename)
		}

		videoName := fmt.Sprintf("reagent_%d_video", id)
		newVideoURL, err = h.Repository.UploadFileToMinIO(context.Background(), videoName, videoHeader)
		if err != nil {
			// Откат изображения если видео не загрузилось
			if uploadedImg {
				imgFilename := fmt.Sprintf("reagent_%d_img", id)
				h.Repository.DeleteFileFromMinIO(imgFilename)
			}
			h.errorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("ошибка загрузки видео: %w", err))
			return
		}
		updates["video"] = newVideoURL
		uploadedVideo = true
	} else if deleteVideo {
		// Удаляем видео
		if currentReagent.Video != "" {
			oldFilename := fmt.Sprintf("reagent_%d_video", id)
			h.Repository.DeleteFileFromMinIO(oldFilename)
		}
		updates["video"] = ""
	}

	// Если ничего не обновляется
	if len(updates) == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("нет данных для обновления"))
		return
	}

	// Применяем обновления
	if err := h.Repository.UpdateReagent(uint(id), updates); err != nil {
		// Откат файлов при ошибке БД
		if uploadedImg {
			imgFilename := fmt.Sprintf("reagent_%d_img", id)
			h.Repository.DeleteFileFromMinIO(imgFilename)
		}
		if uploadedVideo {
			videoFilename := fmt.Sprintf("reagent_%d_video", id)
			h.Repository.DeleteFileFromMinIO(videoFilename)
		}
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Формируем ответ с актуальными данными
	updatedReagent, _ := h.Repository.GetReagentByID(uint(id))

	ctx.JSON(http.StatusOK, gin.H{
		// "status":  "success",
		"data":    updatedReagent,
		"message": "реагент обновлён",
	})
}

// Вспомогательная функция валидации видео
func validateVideo(header *multipart.FileHeader) bool {
	if header.Size > 50<<20 {
		return false
	}
	file, _ := header.Open()
	defer file.Close()

	buffer := make([]byte, 512)
	file.Read(buffer)

	contentType := http.DetectContentType(buffer)
	validTypes := []string{"video/mp4", "video/webm", "video/ogg", "video/quicktime"}
	for _, t := range validTypes {
		if contentType == t {
			return true
		}
	}
	return false
}

func (h *Handler) DeleteReagentAPI(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := h.Repository.DeleteReagent(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "реагент удалён"})
}
