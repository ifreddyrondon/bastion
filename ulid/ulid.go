package ulid

import (
	"database/sql/driver"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
	uuid "github.com/satori/go.uuid"
)

var randPool = &sync.Pool{
	New: func() interface{} {
		seed := time.Now().UnixNano() + rand.Int63()
		return rand.NewSource(seed)
	},
}

// ULID is an ID type provided by bastion that is a lexically sortable UUID.
// The internal representation is an ULID (https://github.com/oklog/ulid).
type ULID uuid.UUID

// New returns a new ULID, which is a lexically sortable UUID.
func New() ULID {
	entropy := randPool.Get().(rand.Source)
	id := ULID(ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(entropy)))
	randPool.Put(entropy)

	return id
}

// NewFromText creates a new ULID from its string representation. Will
// return an error if the text is not a valid ULID.
func NewFromText(text string) (ULID, error) {
	var id ULID
	err := id.UnmarshalText([]byte(text))
	return id, err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Following formats are supported:
// "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
// "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
// "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8"
// Implements the exact same code as the UUID UnmarshalText.
func (id *ULID) UnmarshalText(text []byte) error {
	u := uuid.UUID(*id)
	err := u.UnmarshalText(text)
	*id = ULID(u)
	return err
}

// IsEmpty returns whether the ID is empty or not. An empty ID means it has not
// been set yet.
func (id ULID) IsEmpty() bool {
	return uuid.Equal(uuid.UUID(id), uuid.Nil)
}

// String returns the string representation of the ID.
func (id ULID) String() string {
	return uuid.UUID(id).String()
}

// Value implements the Valuer interface.
func (id ULID) Value() (driver.Value, error) {
	return uuid.UUID(id).Value()
}

// Scan implements the Scanner interface.
func (id *ULID) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		if len(src) != uuid.Size {
			return id.UnmarshalText(src)
		}

		var ulid ulid.ULID
		ulid.UnmarshalBinary(src)
		*id = ULID(ulid)
		return nil
	case string:
		return id.Scan([]byte(src))
	default:
		return fmt.Errorf("cannot scan %T into ULID", src)
	}
}
