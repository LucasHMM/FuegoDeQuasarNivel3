package calculos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessage_ReconstruyeCorrectamente(t *testing.T) {
	msg1 := []string{"este", "", "", "mensaje", ""}
	msg2 := []string{"", "es", "", "", "secreto"}
	msg3 := []string{"este", "", "un", "", ""}

	result, err := GetMessage(msg1, msg2, msg3)
	assert.NoError(t, err)
	assert.Equal(t, "este es un mensaje secreto", result)
}

func TestGetMessage_ErrorSiNoHayPalabras(t *testing.T) {
	msg1 := []string{"", "", ""}
	msg2 := []string{"", "", ""}
	msg3 := []string{"", "", ""}

	result, err := GetMessage(msg1, msg2, msg3)
	assert.Error(t, err)
	assert.Empty(t, result)
}
