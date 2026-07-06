package sessions

import "github.com/gofiber/fiber/v2/middleware/session"

func Set(sess *session.Session, key string, value any) error {
	sess.Set(key, value)
	return sess.Save()
}

func Get(sess *session.Session, key string) any {
	return sess.Get(key)
}

func Delete(sess *session.Session, key string) error {
	sess.Delete(key)
	return sess.Save()
}
