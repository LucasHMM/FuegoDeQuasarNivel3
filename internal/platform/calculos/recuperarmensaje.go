package calculos

import (
	"errors"
	"strings"
)

func GetMessage(KenoviMessage, SkywalkerMessage, SatoMessage []string) (ShipCleanMessage string, err1 error) {
	// Primero encontrar el máximo largo
	maxLen := max(len(KenoviMessage), len(SkywalkerMessage), len(SatoMessage))

	// Normalizar los mensajes al mismo largo agregando "" al inicio si es necesario
	KenoviMessage = normalize(KenoviMessage, maxLen)
	SkywalkerMessage = normalize(SkywalkerMessage, maxLen)
	SatoMessage = normalize(SatoMessage, maxLen)

	// Reconstruir palabra por palabra
	result := make([]string, maxLen)
	for i := 0; i < maxLen; i++ {
		word := ""
		if KenoviMessage[i] != "" {
			word = KenoviMessage[i]
		} else if SkywalkerMessage[i] != "" {
			word = SkywalkerMessage[i]
		} else if SatoMessage[i] != "" {
			word = SatoMessage[i]
		}
		result[i] = word
	}

	// Comprobar si pudimos reconstruir al menos una palabra
	found := false
	for _, w := range result {
		if w != "" {
			found = true
			break
		}
	}
	if !found {
		return "", errors.New("no se pudo reconstruir ningún mensaje")
	}

	// Unir las palabras con espacio
	ShipCleanMessage = strings.TrimSpace(strings.Join(result, " "))
	return ShipCleanMessage, nil
}

// max devuelve el mayor de tres enteros
func max(a, b, c int) int {
	if a >= b && a >= c {
		return a
	} else if b >= a && b >= c {
		return b
	}
	return c
}

// normalize rellena con "" al inicio para igualar longitud
func normalize(msg []string, maxLen int) []string {
	if len(msg) == maxLen {
		return msg
	}
	// agregar "" al inicio
	diff := maxLen - len(msg)
	prefix := make([]string, diff)
	return append(prefix, msg...)
}
