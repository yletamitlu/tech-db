package helpers

import (
	"github.com/yletamitlu/tech-db/internal/consts"
	"strings"
	"sync"
)

type UserCache struct {
	idNickname   map[int]string
	idNicknameMu sync.RWMutex

	nicknameId   map[string]int
	nicknameIdMu sync.RWMutex
}

func NewUserCache() *UserCache {
	return &UserCache{
		idNickname: make(map[int]string),
		nicknameId: make(map[string]int),
	}
}

func (uc *UserCache) getNicknameById(id int) (string, error) {
	uc.idNicknameMu.RLock()
	nickname, ok := uc.idNickname[id]
	uc.idNicknameMu.RUnlock()

	if !ok {
		return "", consts.ErrNotFound
	}

	return nickname, nil
}

func (uc *UserCache) GetNickname(nickname string) (string, error) {
	uc.nicknameIdMu.RLock()
	id, ok := uc.nicknameId[strings.ToLower(nickname)]
	uc.nicknameIdMu.RUnlock()

	if !ok {
		return "", consts.ErrNotFound
	}

	return uc.getNicknameById(id)
}

func (uc *UserCache) Set(id int, nickname string) {
	uc.nicknameIdMu.Lock()
	uc.nicknameId[strings.ToLower(nickname)] = id
	uc.nicknameIdMu.Unlock()

	uc.idNicknameMu.Lock()
	uc.idNickname[id] = nickname
	uc.idNicknameMu.Unlock()
}
