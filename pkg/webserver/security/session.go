package security

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	ac "github.com/praserx/atomic-cache"
	"github.com/praserx/gobarista/pkg/logger"
)

var Sessions *ac.AtomicCache

const SessionKey = "barista"

type Session struct {
	UserID      uint
	UserRole    string
	Code        string
	CodeUsed    bool
	CodeValidTo time.Time
	Logged      bool
}

func init() {
	Sessions = ac.New(ac.OptionMaxRecords(512))
}

func SessionSet(key string, session Session) {
	var binarySessionObject bytes.Buffer
	enc := gob.NewEncoder(&binarySessionObject) // Will write to network.

	if err := enc.Encode(session); err != nil {
		logger.Error(fmt.Sprintf("cannot encode session: %v", err))
	}

	Sessions.Set(key, binarySessionObject.Bytes(), 48*time.Hour)
}

func SessionGet(key string) (Session, bool) {
	var binarySessionObject bytes.Buffer
	dec := gob.NewDecoder(&binarySessionObject) // Will write to network.

	sbs, err := Sessions.Get(key)
	if errors.Is(err, ac.ErrNotFound) {
		return Session{}, false
	} else if err != nil {
		logger.Error(fmt.Sprintf("cannot get session: %v", err))
		return Session{}, false
	}

	binarySessionObject.Write(sbs)

	var session Session
	err = dec.Decode(&session)
	if err != nil {
		logger.Error(fmt.Sprintf("cannot get session: %v", err))
	}

	return session, true
}
