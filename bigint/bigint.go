package bigint

import "math/big"

func dummy() *big.Int {
	return big.NewInt(0)
}

func dummyF() *big.Float {
	return big.NewFloat(0.0)
}

// LessThan determines whether or not big.Int a is less than big.Int b.
// Returns true if a is smaller than b.
func LessThan(a, b *big.Int) bool {
	return a.Cmp(b) < 0
}

// LTE determines whether or not big.Int a is less than or equal to
// big.Int b.
// Returns true if a is less than or smaller than b.
func LTE(a, b *big.Int) bool {
	return LessThan(a, b) || SameAs(a, b)
}

// GT determines whether or not big.Int a is greater than big.Int b.
// Returns true if a is greater than b, false otherwise.
func GT(a, b *big.Int) bool {
	return a.Cmp(b) > 0
}

func GTE(a, b *big.Int) bool {
	return GT(a, b) || SameAs(a, b)
}

func SameAs(a, b *big.Int) bool {
	return a.Cmp(b) == 0
}

func Copy(a *big.Int) *big.Int {
	x := dummy().Set(a)
	//info("CPY BI [%s] -> %s\n", a, x)
	return x
}

func Add(a, b *big.Int) *big.Int {
	x := dummy().Add(a, b)
	//info("ADD BI [%s + %s] -> %s\n", a, b, x)
	return x
}

func Sub(a, b *big.Int) *big.Int {
	x := dummy().Sub(a, b)
	//info("SUB BI [%s - %s] -> %s\n", a, b, x)
	return x
}

func Mul(a, b *big.Int) *big.Int {
	x := dummy().Mul(a, b)
	//info("MUL BI [%s * %s] -> %s\n", a, b, x)
	return x
}

func MulF(a, b *big.Float) *big.Float {
	return dummyF().Mul(a, b)
}

func Div(a, b *big.Int) *big.Int {
	x := dummy().Div(a, b)
	//info("DIV BI [%s - %s] -> %s\n", a, b, x)
	return x
}

func DivF(a, b *big.Int) *big.Float {
	if SameAs(b, big.NewInt(0)) {
		return new(big.Float).SetInt(big.NewInt(0))
	}
	aF := new(big.Float).SetInt(a)
	bF := new(big.Float).SetInt(b)
	return dummyF().Quo(aF, bF)
}

func Mod(a, b *big.Int) *big.Int {
	x := dummy().Mod(a, b)
	//info("MOD BI [%s %% %s] -> %s\n", a, b, x)
	return x
}

func Pow(a, b *big.Int) *big.Int {
	x := dummy().Exp(a, b, nil)
	//info("POW BI [%s ^ %s] -> %s\n", a, b, x)
	return x
}

func Percent(total, current *big.Int) *big.Float {

	if SameAs(total, big.NewInt(0)) {
		return big.NewFloat(0.0)
	}

	tmp := DivF(current, total)

	return MulF(tmp, big.NewFloat(100.0))
}

func ToBigInt(data string) *big.Int {
	i := new(big.Int)
	i.SetString(data, 10)
	return i
}
