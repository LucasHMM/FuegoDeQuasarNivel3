package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_NewYGetSatellite(t *testing.T) {
	repo := New()
	sat, err := repo.GetSatellite("kenobi")
	assert.NoError(t, err)
	assert.Equal(t, "kenobi", sat.Name)
	assert.Equal(t, float32(-500), sat.Position.X)
}

func TestRepository_SaveYGetAllSatellites(t *testing.T) {
	repo := New()
	sat := Satellite{
		Name:     "kenobi",
		Position: Point{X: -500, Y: -200},
		Distance: 100,
		Message:  []string{"hola"},
	}
	err := repo.SaveSatellite(sat)
	assert.NoError(t, err)

	all, err := repo.GetAllSatellites()
	assert.NoError(t, err)
	assert.True(t, len(all) >= 1)
}
