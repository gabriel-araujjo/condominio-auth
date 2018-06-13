package security

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"reflect"
	"testing"
)

func TestClientCode(t *testing.T) {
	var cipherKey [32]byte
	rand.Reader.Read(cipherKey[:])

	cipher, err := aes.NewCipher(cipherKey[:])
	if err != nil {
		t.Errorf("can't create cipher due to %q", err.Error())
		return
	}
	notary := &Notary{
		codeCipher: cipher,
	}
	scope := []int64{1, 2, 3}
	var clientID int64 = 1
	var userID int64 = 233
	code, err := notary.NewClientCode(clientID, scope, userID)
	if err != nil {
		t.Errorf("error while creating code %q", err.Error())
		return
	}
	decipheredClientID, decipheredScope, decipheredUserID, err := notary.DecipherCode(code)
	if err != nil {
		t.Errorf("error while deciphering token %q", err.Error())
		return
	}

	if decipheredClientID != clientID {
		t.Error("unmatch clientID")
		return
	}

	if !reflect.DeepEqual(decipheredScope, scope) {
		t.Error("unmatch scope")
	}

	if userID != decipheredUserID {
		t.Error("unmatch userID")
		return
	}
}

func TestCodeUniqueness(t *testing.T) {
	var cipherKey [32]byte
	rand.Reader.Read(cipherKey[:])

	cipher, err := aes.NewCipher(cipherKey[:])
	if err != nil {
		t.Errorf("can't create cipher due to %q", err.Error())
		return
	}
	notary := &Notary{
		codeCipher: cipher,
	}
	var scope [29]int64
	var clientID int64 = 1
	var userID int64 = 233
	var codes map[string]struct{}
	for i := 0; i < 128; i++ {
		code, _ := notary.NewClientCode(clientID, scope[:], userID)
		if _, contains := codes[code]; contains {
			t.Fatalf("code %q is duplicated", code)
		}
		fmt.Printf("%s\n", code)
	}
}
