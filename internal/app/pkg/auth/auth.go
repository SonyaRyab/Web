package auth

import "lab3/internal/app/ds"

// CurrentUserID возвращает ID текущего пользователя (заглушка для 3ей лабы)
// func CurrentUserID() uint {
// 	return 1 // Хардкод для тестового пользователя "test"
// }

// // CurrentUserIsModerator проверяет, является ли текущий пользователь модератором
// func CurrentUserIsModerator() bool {
// 	// В реальности нужно доставать из БД, пока хардкод
// 	return false
// }

// // GetCurrentUser возвращает фиксированного пользователя для демонстрации
// func GetCurrentUser() ds.User {
// 	return ds.User{
// 		ID:          1,
// 		Username:    "test",
// 		Login:       "test",
// 		IsModerator: false,
// 	}
// }


package auth

import (
    "sync"
    "lab3/internal/app/ds"
)

// Приватная структура
type userService struct {
    user ds.User
}

var (
    instance *userService
    once     sync.Once  //выполнение только один раз
)

// точка входа
func GetUserService() *userService {
    once.Do(func() {  // Внутри блокировка - только один поток создаст объект
        instance = &userService{
            user: ds.User{
                ID:       1,
                Username: "test",
            },
        }
    })
    return instance
}

func (s *userService) GetUser() ds.User {
    return s.user
}

func (s *userService) IsModerator() bool {
    return s.user.IsModerator
}