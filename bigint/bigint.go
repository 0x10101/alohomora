package bigint

import "math/big"

func dummy() *big.Int {
	return big.NewInt(0)
}

func dummyF() *big.Float {
	return big.NewFloat(0.0)
}

// Lt determines whether or not big.Int a is less than big.Int b.
// Returns true if a is smaller than b.
func Lt(a, b *big.Int) bool {
	return a.Cmp(b) < 0
}

// LtE determines whether or not big.Int a is less than or equal to
// big.Int b.
// Returns true if a is less than or smaller than b.
func LtE(a, b *big.Int) bool {
	return Lt(a, b) || Eq(a, b)
}

// Gt determines whether or not big.Int a is greater than big.Int b.
// Returns true if a is greater than b, false otherwise.
func Gt(a, b *big.Int) bool {
	return a.Cmp(b) > 0
}

// GtE determines whether big.Int a is greater than or equal to big.Int b.
// Returns true if a is greater than or equal to b, false otherwise.
func GtE(a, b *big.Int) bool {
	return Gt(a, b) || Eq(a, b)
}

// Eq determines whether big.Int a is the same as big.Int b.
// Returns true if a equals b (same value), false otherwise.
func Eq(a, b *big.Int) bool {
	return a.Cmp(b) == 0
}

// Cp copies a big.Int and returns the copy.
func Cp(a *big.Int) *big.Int {
	x := dummy().Set(a)
	return x
}

// Add adds big.Int a to big.Int b and returns the result.
func Add(a, b *big.Int) *big.Int {
	x := dummy().Add(a, b)
	return x
}

// Sub subtracts big.Int b from big.Int a and returns the result.
func Sub(a, b *big.Int) *big.Int {
	x := dummy().Sub(a, b)
	return x
}

// Mul multiplies big.Int a with big.Int b and returns the result.
func Mul(a, b *big.Int) *big.Int {
	x := dummy().Mul(a, b)
	return x
}

// MulF multiplies big.Float a with big.Float b and returns the result.
func MulF(a, b *big.Float) *big.Float {
	return dummyF().Mul(a, b)
}

// Div divides big.Int a by big.Int b and returns the result.
func Div(a, b *big.Int) *big.Int {
	x := dummy().Div(a, b)
	return x
}

// DivF divides big.Float a by big.Float b and returns the result.
func DivF(a, b *big.Int) *big.Float {
	if Eq(b, big.NewInt(0)) {
		return new(big.Float).SetInt(big.NewInt(0))
	}
	aF := new(big.Float).SetInt(a)
	bF := new(big.Float).SetInt(b)
	return dummyF().Quo(aF, bF)
}

// Mod divides big.Int a by big.Int b and returns the remainder.
func Mod(a, b *big.Int) *big.Int {
	x := dummy().Mod(a, b)
	return x
}

// Pow calculates a ** b and returns the result.
func Pow(a, b *big.Int) *big.Int {
	x := dummy().Exp(a, b, nil)
	return x
}

// Percent calculates a percentage based on a total and a current value and returns it.
func Percent(total, current *big.Int) *big.Float {

	if Eq(total, big.NewInt(0)) {
		return big.NewFloat(0.0)
	}

	tmp := DivF(current, total)

	return MulF(tmp, big.NewFloat(100.0))
}

// ToBigInt converts a string to a big.Int and returns it.
func ToBigInt(data string) *big.Int {
	i := new(big.Int)
	i.SetString(data, 10)
	return i
}
