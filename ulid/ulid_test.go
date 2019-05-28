package ulid_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ulidok "github.com/oklog/ulid"

	"github.com/ifreddyrondon/bastion/ulid"
)

func TestULID_IsEmpty(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	var id ulid.ULID
	a.True(id.IsEmpty())

	id = ulid.New()
	a.False(id.IsEmpty())
}

func TestNewULID_FromText(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	u, err := ulid.NewFromText("0167c8a5-d308-8692-809d-b1ad4a2d9562")
	a.Nil(err)
	a.False(u.IsEmpty())
}

func TestNewULID_FromTextFails(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	u, err := ulid.NewFromText("a")
	a.True(u.IsEmpty())
	a.EqualError(err, "uuid: incorrect UUID length: a")
}

func TestULID_String(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	u, err := ulid.NewFromText("0167c8a5-d308-8692-809d-b1ad4a2d9562")
	a.Nil(err)
	a.Equal("0167c8a5-d308-8692-809d-b1ad4a2d9562", u.String())
}

func TestULID_Value(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	id := ulid.New()
	v, _ := id.Value()
	a.Equal(id.String(), v)
}

func TestULID_ScanValue(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	expected := ulid.New()
	v, err := expected.Value()
	a.NoError(err)

	var id ulid.ULID
	a.NoError(id.Scan(v))

	a.Equal(expected, id)
	a.Equal(expected.String(), id.String())
}

func TestULID_ScanValue_binary(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	expected := ulid.New()
	b, err := ulidok.ULID(expected).MarshalBinary()
	a.NoError(err)

	var id ulid.ULID
	a.NoError(id.Scan(b))
	a.Equal(expected, id)
	a.Equal(expected.String(), id.String())
}

func TestULID_ScanValue_err(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	var id ulid.ULID
	a.EqualError(id.Scan(12), "cannot scan int into ULID")
}
