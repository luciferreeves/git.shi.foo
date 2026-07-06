package sessions

import "github.com/gofiber/fiber/v2/middleware/session"

func CreateSession(sess *session.Session, providerID string) error {
	return Set(sess, SessionAuthKey, providerID)
}

func DestroySession(sess *session.Session) error {
	return Delete(sess, SessionAuthKey)
}

func GetSessionProviderID(sess *session.Session) string {
	value := Get(sess, SessionAuthKey)
	if value == nil {
		return ""
	}

	providerID, ok := value.(string)
	if !ok {
		return ""
	}

	return providerID
}
