package redis

import (
	"encoding/json"
	"fmt"
	"time"
)

// CacheMessages stores messages between two users in Redis
func CacheMessages(senderID, receiverID string, messages interface{}) error {
	key := fmt.Sprintf("messages:%s:%s", senderID, receiverID)
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	return Client.Set(key, data, 10*time.Minute).Err()
}

// GetCachedMessages retrieves cached messages between two users
func GetCachedMessages(senderID, receiverID string, result interface{}) error {
	key := fmt.Sprintf("messages:%s:%s", senderID, receiverID)
	data, err := Client.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), result)
}

// InvalidateMessagesCache removes cached messages for both directions of conversation
func InvalidateMessagesCache(senderID, receiverID string) error {
	key1 := fmt.Sprintf("messages:%s:%s", senderID, receiverID)
	key2 := fmt.Sprintf("messages:%s:%s", receiverID, senderID)
	return Client.Del(key1, key2).Err()
}

// CacheUser stores user data in Redis
func CacheUser(userID string, userData interface{}) error {
	key := fmt.Sprintf("user_data:%s", userID)
	data, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	return Client.Set(key, data, 30*time.Minute).Err()
}

// GetCachedUser retrieves cached user data
func GetCachedUser(userID string, result interface{}) error {
	key := fmt.Sprintf("user_data:%s", userID)
	data, err := Client.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), result)
}

// DeleteCachedUser removes cached user data
func DeleteCachedUser(userID string) error {
	key := fmt.Sprintf("user_data:%s", userID)
	return Client.Del(key).Err()
}

// CacheMessageCount stores message count for pagination
func CacheMessageCount(senderID, receiverID string, count int64) error {
	key := fmt.Sprintf("message_count:%s:%s", senderID, receiverID)
	return Client.Set(key, count, 5*time.Minute).Err()
}

// GetCachedMessageCount retrieves cached message count
func GetCachedMessageCount(senderID, receiverID string) (int64, error) {
	key := fmt.Sprintf("message_count:%s:%s", senderID, receiverID)
	return Client.Get(key).Int64()
}

// DeleteMessageCount removes cached message count
func DeleteMessageCount(senderID, receiverID string) error {
	key1 := fmt.Sprintf("message_count:%s:%s", senderID, receiverID)
	key2 := fmt.Sprintf("message_count:%s:%s", receiverID, senderID)
	return Client.Del(key1, key2).Err()
}

// CacheConversationList stores user's conversation list
func CacheConversationList(userID string, conversations interface{}) error {
	key := fmt.Sprintf("conversations:%s", userID)
	data, err := json.Marshal(conversations)
	if err != nil {
		return err
	}
	return Client.Set(key, data, 15*time.Minute).Err()
}

// GetCachedConversationList retrieves cached conversation list
func GetCachedConversationList(userID string, result interface{}) error {
	key := fmt.Sprintf("conversations:%s", userID)
	data, err := Client.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), result)
}

// InvalidateConversationList removes cached conversation list
func InvalidateConversationList(userID string) error {
	key := fmt.Sprintf("conversations:%s", userID)
	return Client.Del(key).Err()
}

// CacheLastMessage stores the last message in a conversation
func CacheLastMessage(senderID, receiverID string, message interface{}) error {
	key := fmt.Sprintf("last_message:%s:%s", senderID, receiverID)
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return Client.Set(key, data, 20*time.Minute).Err()
}

// GetCachedLastMessage retrieves the last message in a conversation
func GetCachedLastMessage(senderID, receiverID string, result interface{}) error {
	key := fmt.Sprintf("last_message:%s:%s", senderID, receiverID)
	data, err := Client.Get(key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), result)
}

// InvalidateLastMessage removes cached last message
func InvalidateLastMessage(senderID, receiverID string) error {
	key1 := fmt.Sprintf("last_message:%s:%s", senderID, receiverID)
	key2 := fmt.Sprintf("last_message:%s:%s", receiverID, senderID)
	return Client.Del(key1, key2).Err()
}
