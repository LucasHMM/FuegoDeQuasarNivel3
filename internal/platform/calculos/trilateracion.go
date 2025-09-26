package calculos

import (
	"fmt"
	"math"
)

type Point struct {
	X float64
	Y float64
}

type Point32 struct {
	X float32
	Y float32
}

// checkPairwiseIntersection devuelve false si dos circunferencias no pueden intersectar
func checkPairwiseIntersection(p1, p2 Point, r1, r2 float64) bool {
	d := math.Hypot(p2.X-p1.X, p2.Y-p1.Y)
	// No intersectan si están demasiado separadas o una está completamente dentro de la otra sin tocar
	if d > r1+r2 {
		return false
	}
	if d < math.Abs(r1-r2) {
		return false
	}
	// si d == 0 y r1 == r2 -> coincidentes (inf soluciones) -> tratamos como degenerado
	if d == 0 && r1 == r2 {
		return false
	}
	return true
}

// GetLocation enmascara la función Trilateracion para que cumpla con la firma requerida.
func GetLocation(p1, p2, p3 Point32, r1, r2, r3 float32) Point32 {
	// Convertimos Point32 → Point (float64)
	p1f := Point{X: float64(p1.X), Y: float64(p1.Y)}
	p2f := Point{X: float64(p2.X), Y: float64(p2.Y)}
	p3f := Point{X: float64(p3.X), Y: float64(p3.Y)}

	// Radios a float64
	r1f := float64(r1)
	r2f := float64(r2)
	r3f := float64(r3)

	// tolerancia para residuos
	tol := 10.0
	// Llamo a Trilateracion
	pos, err := Trilateracion(p1f, p2f, p3f, r1f, r2f, r3f, tol)
	if err != nil {
		// Devolver un valor por defecto cuando falla:
		return Point32{X: 0, Y: 0}
	}

	// Convertimos el resultado a Point32
	return Point32{
		X: float32(pos.X),
		Y: float32(pos.Y),
	}
}

// Trilateracion calcula la posición (x, y) de la fuente
// a partir de tres puntos y sus distancias. Resuelve usando la
// linealizacion (resta de circunferencias) y verifica residuos
// tol es la tolerancia máxima aceptable en unidades de distancia
// (ej. 0.1, 0.5, para datos con mucho ruido o 1e-6 para datos muy precisos).
func Trilateracion(p1, p2, p3 Point, r1, r2, r3 float64, tol float64) (Point, error) {

	if !checkPairwiseIntersection(p1, p2, r1, r2) {
		return Point{}, fmt.Errorf("las circunferencias 1 y 2 no pueden intersectar (pareja incoherente)")
	}
	if !checkPairwiseIntersection(p1, p3, r1, r3) {
		return Point{}, fmt.Errorf("las circunferencias 1 y 3 no pueden intersectar (pareja incoherente)")
	}
	if !checkPairwiseIntersection(p2, p3, r2, r3) {
		return Point{}, fmt.Errorf("las circunferencias 2 y 3 no pueden intersectar (pareja incoherente)")
	}

	// Construimos las ecuaciones lineales (resta de circunferencias)
	A := 2 * (p2.X - p1.X)
	B := 2 * (p2.Y - p1.Y)
	C := r1*r1 - r2*r2 - p1.X*p1.X + p2.X*p2.X - p1.Y*p1.Y + p2.Y*p2.Y

	D := 2 * (p3.X - p1.X)
	E := 2 * (p3.Y - p1.Y)
	F := r1*r1 - r3*r3 - p1.X*p1.X + p3.X*p3.X - p1.Y*p1.Y + p3.Y*p3.Y

	// Determinante de la matriz
	den := A*E - B*D
	if math.Abs(den) < 1e-12 {
		return Point{}, fmt.Errorf("determinante cero o casi cero: puntos colineales o configuración incoherente")
	}

	x := (C*E - B*F) / den
	y := (A*F - C*D) / den

	// Calculamos residuos frente a las ecuaciones originales
	d1 := math.Hypot(x-p1.X, y-p1.Y)
	d2 := math.Hypot(x-p2.X, y-p2.Y)
	d3 := math.Hypot(x-p3.X, y-p3.Y)

	rerr1 := math.Abs(d1 - r1)
	rerr2 := math.Abs(d2 - r2)
	rerr3 := math.Abs(d3 - r3)

	maxErr := math.Max(rerr1, math.Max(rerr2, rerr3))

	if maxErr > tol {
		return Point{}, fmt.Errorf("no hay intersección coherente: max residual = %.6f > tol(%.6f). residuos = [%.6f, %.6f, %.6f]",
			maxErr, tol, rerr1, rerr2, rerr3)
	}

	return Point{X: x, Y: y}, nil
}

/*
// TrilateracionLS resuelve la posición usando mínimos cuadrados
func TrilateracionLS(p1, p2, p3 Point, r1, r2, r3 float64) (Point, error) {
	// Sistema lineal: A * [x y]^T = b
	// Se construye restando las circunferencias
	A := mat.NewDense(2, 2, []float64{
		2 * (p2.X - p1.X), 2 * (p2.Y - p1.Y),
		2 * (p3.X - p1.X), 2 * (p3.Y - p1.Y),
	})

	b := mat.NewVecDense(2, []float64{
		r1*r1 - r2*r2 - p1.X*p1.X + p2.X*p2.X - p1.Y*p1.Y + p2.Y*p2.Y,
		r1*r1 - r3*r3 - p1.X*p1.X + p3.X*p3.X - p1.Y*p1.Y + p3.Y*p3.Y,
	})

	var x mat.VecDense
	if err := x.SolveVec(A, b); err != nil {
		return Point{}, fmt.Errorf("error resolviendo mínimos cuadrados: %v", err)
	}

	return Point{X: x.AtVec(0), Y: x.AtVec(1)}, nil
}
*/
