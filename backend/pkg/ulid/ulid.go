package ulid

import (
	"regexp"
	"time"

	"github.com/o-ga09/zenn-hackthon-2026/pkg/errors"
	"github.com/oklog/ulid/v2"
	"golang.org/x/exp/rand"
)

type ULIDValue struct {
	Value string
}

func NewULID(id string) (ULIDValue, error) {
	if id == "" {
		return ULIDValue{}, errors.ErrEmptyULID
	}

	// 正規表現でチェックする
	regexp := regexp.MustCompile(`^[0-9A-Z]{26}$`)

	if !regexp.MatchString(id) {
		return ULIDValue{}, errors.ErrInvalidULID
	}
	return ULIDValue{Value: id}, nil
}

func GenerateULID() (string, error) {
	entropy := rand.New(rand.NewSource(uint64(time.Now().Nanosecond())))
	ms := ulid.Timestamp(time.Now())
	ulid, err := ulid.New(ms, entropy)

	if err != nil {
		return "", err
	}

	ulidValue, err := NewULID(ulid.String())

	if err != nil {
		return "", err
	}

	return ulidValue.String(), nil
}

func (u *ULIDValue) Equals(target *ULIDValue) bool {
	return u.Value == target.Value
}

func (u *ULIDValue) String() string {
	return u.Value
}
