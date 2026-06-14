package handler

import (
	"fmt"
	"lab3/internal/app/pkg/auth"
	"mime/multipart"
	"net/http"
	"strconv"

	"lab3/internal/app/ds"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMethanesAPI(ctx *gin.Context) {
	var methanes []ds.Methane
	var err error

	methanes, err = h.Repository.GetMethanes()

	if err != nil {
		if err.Error() != "массив пустой" {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		} else {
			methanes = []ds.Methane{}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":   methanes,
	})
}

func (h *Handler) GetMethaneByIdAPI(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	methane, err := h.Repository.GetMethaneWithReagents(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	fullMethane := ds.FullMethaneSerializer{
		ID:           methane.ID,
		Name:         methane.Name,
		Status:       methane.Status,
		DateCreate:   methane.DateCreate,
		DateForm:     methane.DateForm,
		DateFinish:   methane.DateFinish,
		Temperature:  methane.Temperature,
		MethaneYield: methane.MethaneYield,
		Img:          methane.Img,
		Admin:        methane.Admin,
		Reagents:     methane.Reagents,
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"data":   fullMethane,
	})
}

func (h *Handler) CreateDraftMethaneAPI(ctx *gin.Context) {
	userID := auth.CurrentUserID()

	methane, err := h.Repository.CreateDraftMethane(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data":    methane,
		"message": "черновик создан",
	})
}

func (h *Handler) AddMethaneAPI(ctx *gin.Context) {
	// Прочитаем в ОЗУ 2 Мб данных формы
	// Этого должно хватить для текстовых полей и несложных изображений (логотипов)

	err := ctx.Request.ParseMultipartForm(2 << 20)

	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Пытаемся получить файл из формы
	// header имеет тип *multipart.FileHeader и является указателем
	// на метаданные файла
	header, err := ctx.FormFile("pic")

	// Флаг, устанавливаемый, когда поле pic сопоставлено с неким файлом
	fileFound := false

	if err != nil {
		// Файла в запросе нет - допустимая ситуация
		if err != http.ErrMissingFile {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
	} else {
		fileFound = true
	}

	if fileFound {
		code, err := validateFileUpload(header)

		if err != nil {
			h.errorHandler(ctx, code, err)
			return
		}
	}

	methane := ds.Methane{
		Name:   ctx.Request.FormValue("name"),
		AdminID: 1, // временный хардкод
	}

	if tempStr := ctx.Request.FormValue("temperature"); tempStr != "" {
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		methane.Temperature = temp
	}

	err = h.Repository.AddMethane(&methane)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if fileFound {
		if err = h.Repository.AddOrReplaceMethaneImage(methane.ID, header); err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		updatedMethane, err := h.Repository.GetMethane(int(methane.ID))
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"data":    updatedMethane,
			"message": "метан успешно добавлен",
		})
	} else {
		ctx.JSON(http.StatusCreated, gin.H{
			"data":    methane,
			"message": "метан успешно добавлен",
		})
	}
}

func (h *Handler) FormMethaneAPI(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	methane, err := h.Repository.GetMethane(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Проверка прав
	if methane.AdminID != auth.CurrentUserID() {
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("нет прав"))
		return
	}

	if methane.Status != "черновик" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("можно формировать только черновик"))
		return
	}

	// Проверка обязательных полей
	if methane.Name == "" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("название обязательно"))
		return
	}

	// Подсчёт реагентов для вычисляемого поля
	count, _ := h.Repository.CountReagentsInMethane(uint(id))
	if count == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("добавьте хотя бы один реагент"))
		return
	}

	// Вычисление выхода метана (пример формулы из лабы 2)
	var totalMass float64
	reagents, _ := h.Repository.GetMethaneReagents(uint(id))
	for _, r := range reagents {
		totalMass += r.Volume
	}

	// Формула выхода метана (упрощённая)
	methaneYield := (totalMass * 0.15) / (1 + methane.Temperature/1000)

	updates := map[string]interface{}{
		"methane_yield": methaneYield,
	}

	if err := h.Repository.FormMethane(uint(id), updates); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "заявка сформирована",
	})
}

func (h *Handler) UpdateMethaneAPI(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)

	methane, err := h.Repository.GetMethane(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// только admin может редактировать
	if methane.AdminID != auth.CurrentUserID() {
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("нет прав"))
		return
	}

	if methane.Status != "черновик" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("можно редактировать только черновик"))
		return
	}

	if err := ctx.Request.ParseMultipartForm(2 << 20); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Собираем обновления
	updates := map[string]interface{}{}

	if name := ctx.Request.FormValue("name"); name != "" {
		updates["name"] = name
	}

	// Парсим temperature как float64
	if tempStr := ctx.Request.FormValue("temperature"); tempStr != "" {
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("некорректная температура"))
			return
		}
		updates["temperature"] = temp
	}

	// Обработка файла
	header, err := ctx.FormFile("pic")

	var fileFound bool

	if err != nil {
		// Файла в запросе нет - допустимая ситуация
		if err != http.ErrMissingFile {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "File error: " + err.Error()})
			return
		}
	} else {
		fileFound = true
	}

	if fileFound {
		code, err := validateFileUpload(header)

		if err != nil {
			h.errorHandler(ctx, code, err)
			return
		}
	}

	if err = h.Repository.ModifyMethane(uint(id), &methane); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if fileFound {
		if err = h.Repository.AddOrReplaceMethaneImage(methane.ID, header); err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	updatedMethane, err := h.Repository.GetMethane(int(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    updatedMethane,
		"message": "запись успешно обновлена",
	})
}

func (h *Handler) CompleteMethaneAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Получаем заявку
	methane, err := h.Repository.GetMethane(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}

	// Проверяем, что заявка в статусе "сформирована"
	if methane.Status != "сформирована" {
		h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("можно завершить только сформированную заявку"))
		return
	}

	// Проверяем права модератора (в 3 лабе - заглушка, но можно проверить)
	if !auth.CurrentUserIsModerator() {
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("только модератор может завершить заявку"))
		return
	}

	// Завершаем заявку
	moderatorID := auth.CurrentUserID()
	if err := h.Repository.CompleteMethane(uint(id), moderatorID, "завершена"); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "заявка завершена",
	})
}

// func (h *Handler) DeleteMethaneAPI(ctx *gin.Context) {
// 	strId := ctx.Param("id")
// 	id, err := strconv.Atoi(strId)
// 	if err != nil {
// 		h.errorHandler(ctx, http.StatusInternalServerError, err)
// 		return
// 	}

// 	err = h.Repository.DeleteMethane(uint(id))
// 	if err != nil {
// 		h.errorHandler(ctx, http.StatusInternalServerError, err)
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"date_finish":  date_finish,
// 		"message": "заявка успешно удалена",
// 	})
// }

func (h *Handler) DeleteMethaneAPI(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	// Проверка прав (опционально)
	methane, err := h.Repository.GetMethane(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	
	if methane.AdminID != auth.CurrentUserID() {
		h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("нет прав"))
		return
	}

	// Мягкое удаление
	deletedMethane, err := h.Repository.DeleteMethane(uint(id))
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"date_finish": deletedMethane.DateFinish,
		"message":     "заявка успешно удалена",
	})
}

// Вспомогательная функция, определяющая по заголовку
// Content-File файла, является ли тот изображением
func isImage(contentType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, t := range imageTypes {
		if contentType == t {
			return true
		}
	}
	return false
}

// Вспомогательная функция, выполняющая валидацию загруженного файла
func validateFileUpload(header *multipart.FileHeader) (int, error) {
	// Окрываем чтение файлового потока
	file, err := header.Open()

	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("не удалось получить файл")
	}

	defer file.Close()

	// Знакомая логика определения типа содержимого...

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("не удалось прочитать файл")
	}

	contentType := http.DetectContentType(buffer)

	if !isImage(contentType) {
		return http.StatusBadRequest, fmt.Errorf("файл должен быть изображением")
	}

	// Позиция чтения файла возвращается к исходной
	_, err = file.Seek(0, 0)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("ошибка при обработке файла")
	}

	return 0, nil
}

func (h *Handler) GetCartAPI(ctx *gin.Context) {
	userID := auth.CurrentUserID()

	// Ищем черновик пользователя
	methane, err := h.Repository.GetDraftMethane(userID)
	if err != nil {
		// Создаём новый черновик если нет
		methane, err = h.Repository.CreateDraftMethane(userID)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	count, _ := h.Repository.CountReagentsInMethane(methane.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"methane_id": methane.ID,
			"count":      count,
		},
	})
}
