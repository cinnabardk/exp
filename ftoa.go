package main

import (
	"fmt"
	"math"
	"strconv"
)

func ftoa(f float64, prec int) string {
	fr, exp := math.Frexp(f)
	v := int64(fr * (1 << 53))
	e := exp - 53

	buf := make([]byte, 1000)
	n := copy(buf, strconv.FormatInt(v, 10))
	//n := copy(buf, strconv.Itoa64(v))

	for ; e > 0; e-- {
		δ := 0
		if buf[0] >= '5' {
			δ = 1
		}
		x := byte(0)
		for i := n - 1; i >= 0; i-- {
			x += 2 * (buf[i] - '0')
			x, buf[i+δ] = x/10, x%10+'0'
		}
		if δ == 1 {
			buf[0] = '1'
			n++
		}
	}
	dp := n

	for ; e < 0; e++ {
		if buf[n-1]%2 != 0 {
			buf[n] = '0'
			n++
		}
		δ, x := 0, byte(0)
		if buf[0] < '2' {
			δ, x = 1, buf[0]-'0'
			n--
			dp--
		}
		for i := 0; i < n; i++ {
			x = x*10 + buf[i+δ] - '0'
			buf[i], x = x/2+'0', x%2
		}
	}

	if prec > 0 {
		if n > prec {
			if buf[prec] > '5' || buf[prec] == '5' && (nonzero(buf[prec+1:n]) || buf[prec-1]%2 == 1) {
				i := prec - 1
				for i >= 0 && buf[i] == '9' {
					buf[i] = '0'
					i--
				}
				if i >= 0 {
					buf[i]++
				} else {
					buf[0] = '1'
					dp++
				}
			}
			n = prec
		}
		for n < prec {
			buf[n] = '0'
			n++
		}
	}
	return fmt.Sprintf("%c.%se%+d", buf[0], buf[1:n], dp-1)
}

func nonzero(buf []byte) bool {
	for _, c := range buf {
		if c != '0' {
			return true
		}
	}
	return false
}

// Difficult boundary cases, derived from tables given in
//	Vern Paxson, A Program for Testing IEEE Decimal-Binary Conversion
//	ftp://ftp.ee.lbl.gov/testbase-report.ps.Z
//
var ftoaTests = []struct {
	N int
	F float64
	A string
}{
	// Table 3: Stress Inputs for Converting 53-bit Binary to Decimal, < 1/2 ULP
	{0, math.Ldexp(8511030020275656, -342), "9.e-88"},
	{1, math.Ldexp(5201988407066741, -824), "4.6e-233"},
	{2, math.Ldexp(6406892948269899, +237), "1.41e+87"},
	{3, math.Ldexp(8431154198732492, +72), "3.981e+37"},
	{4, math.Ldexp(6475049196144587, +99), "4.1040e+45"},
	{5, math.Ldexp(8274307542972842, +726), "2.92084e+234"},
	{6, math.Ldexp(5381065484265332, -456), "2.891946e-122"},
	{7, math.Ldexp(6761728585499734, -1057), "4.3787718e-303"},
	{8, math.Ldexp(7976538478610756, +376), "1.22770163e+129"},
	{9, math.Ldexp(5982403858958067, +377), "1.841552452e+129"},
	{10, math.Ldexp(5536995190630837, +93), "5.4835744350e+43"},
	{11, math.Ldexp(7225450889282194, +710), "3.89190181146e+229"},
	{12, math.Ldexp(7225450889282194, +709), "1.945950905732e+229"},
	{13, math.Ldexp(8703372741147379, +117), "1.4460958381605e+51"},
	{14, math.Ldexp(8944262675275217, -1001), "4.17367747458531e-286"},
	{15, math.Ldexp(7459803696087692, -707), "1.107950772878888e-197"},
	{16, math.Ldexp(6080469016670379, -381), "1.2345501366327440e-99"},
	{17, math.Ldexp(8385515147034757, +721), "9.25031711960365024e+232"},
	{18, math.Ldexp(7514216811389786, -828), "4.198047150284889840e-234"},
	{19, math.Ldexp(8397297803260511, -345), "1.1716315319786511046e-88"},
	{20, math.Ldexp(6733459239310543, +202), "4.32810072844612493629e+76"},
	{21, math.Ldexp(8091450587292794, -473), "3.317710118160031081518e-127"},

	// Table 4: Stress Inputs for Converting 53-bit Binary to Decimal, > 1/2 ULP
	{0, math.Ldexp(6567258882077402, +952), "3.e+302"},
	{1, math.Ldexp(6712731423444934, +535), "7.6e+176"},
	{2, math.Ldexp(6712731423444934, +534), "3.78e+176"},
	{3, math.Ldexp(5298405411573037, -957), "4.350e-273"},
	{4, math.Ldexp(5137311167659507, -144), "2.3037e-28"},
	{5, math.Ldexp(6722280709661868, +363), "1.26301e+125"},
	{6, math.Ldexp(5344436398034927, -169), "7.142211e-36"},
	{7, math.Ldexp(8369123604277281, -853), "1.3934574e-241"},
	{8, math.Ldexp(8995822108487663, -780), "1.41463449e-219"},
	{9, math.Ldexp(8942832835564782, -383), "4.539277920e-100"},
	{10, math.Ldexp(8942832835564782, -384), "2.2696389598e-100"},
	{11, math.Ldexp(8942832835564782, -385), "1.13481947988e-100"},
	{12, math.Ldexp(6965949469487146, -249), "7.700366561890e-60"},
	{13, math.Ldexp(6965949469487146, -250), "3.8501832809448e-60"},
	{14, math.Ldexp(6965949469487146, -251), "1.92509164047238e-60"},
	{15, math.Ldexp(7487252720986826, +548), "6.898586531774201e+180"},
	{16, math.Ldexp(5592117679628511, +164), "1.3076622631878654e+65"},
	{17, math.Ldexp(8887055249355788, +665), "1.36052020756121240e+216"},
	{18, math.Ldexp(6994187472632449, +690), "3.592810217475959676e+223"},
	{19, math.Ldexp(8797576579012143, +588), "8.9125197712484551899e+192"},
	{20, math.Ldexp(7363326733505337, +272), "5.58769757362301140950e+97"},
	{21, math.Ldexp(8549497411294502, -448), "1.176257830728540379990e-119"},

	{3, 12345000, "1.234e+7"},
}

func testFToa() {
	fmt.Printf("%s\n", ftoa(math.Pi, 50)) // 50 digits of math.Pi (not 50 digits of π)
	//fmt.Printf("%s\n", ftoa(math.Nextafter(0, 1), 0))
	for _, tt := range ftoaTests {
		if a := ftoa(tt.F, tt.N+1); a != tt.A {
			fmt.Printf("ftoa(%g, %d) = %q, want %q\n", tt.F, tt.N+1, a, tt.A)
		}
	}
}
