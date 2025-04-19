package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mislu/market-api/internal/utils/lib"
	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	salt, err := lib.GenerateSalt()
	if err != nil {
		t.Fatal(err)
	}

	text := "How original"
	now := time.Now()
	encrypt, err := lib.EncryptPassword(text, salt)
	fmt.Println(time.Since(now))
	if err != nil {
		t.Fatal(err)
	}

	reencrypt, err := lib.EncryptPassword(text, salt)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, encrypt, reencrypt)
}
