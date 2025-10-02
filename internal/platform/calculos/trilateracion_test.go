package calculos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrilateracion_PosicionCorrecta(t *testing.T) {
	p1 := Point{X: -500, Y: -200}
	p2 := Point{X: 100, Y: -100}
	p3 := Point{X: 500, Y: 100}
	r1 := 927.75
	r2 := 360.0
	r3 := 360.0
	tol := 1.0

	pos, err := Trilateracion(p1, p2, p3, r1, r2, r3, tol)
	assert.NoError(t, err)
	assert.InDelta(t, 426.4, pos.X, 0.1)
	assert.InDelta(t, -252.8, pos.Y, 0.1)
}

func TestTrilateracion_ErrorPorCircunferenciasNoIntersectan(t *testing.T) {
	p1 := Point{X: 0, Y: 0}
	p2 := Point{X: 0, Y: 0}
	p3 := Point{X: 0, Y: 0}
	r1 := 1
	r2 := 1
	r3 := 1
	tol := 1.0

	_, err := Trilateracion(p1, p2, p3, r1, r2, r3, tol)
	assert.Error(t, err)
}
