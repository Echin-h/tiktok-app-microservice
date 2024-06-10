package utils

import "time"

const (
	CacheExpire = time.Hour
)

func GenUserInfoCacheKey(userId int64) string {
	return "user_info_" + Int642Str(userId)
}

func GenPopUserListCacheKey() string {
	return "pop_user_list"
}

func GenFollowUserCacheKey(userId, followUserId int64) string {
	return "follow_user_" + Int642Str(userId) + "_" + Int642Str(followUserId)
}
