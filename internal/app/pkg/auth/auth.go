package auth

import "lab3/internal/app/ds"

// CurrentUserID возвращает ID текущего пользователя (заглушка для 3ей лабы)
func CurrentUserID() uint {
	return 1 // Хардкод для тестового пользователя "test"
}

// CurrentUserIsModerator проверяет, является ли текущий пользователь модератором
func CurrentUserIsModerator() bool {
	// В реальности нужно доставать из БД, пока хардкод
	return false
}

// GetCurrentUser возвращает фиксированного пользователя для демонстрации
func GetCurrentUser() ds.User {
	return ds.User{
		ID:          1,
		Username:    "test",
		Login:       "test",
		IsModerator: false,
	}
}
