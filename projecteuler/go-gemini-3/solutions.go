package main

import (
	"math/big"
	"strconv"
)

func solution_0001() string {
	sum := 0
	for i := 0; i < 1000; i++ {
		if i%3 == 0 || i%5 == 0 {
			sum += i
		}
	}
	return strconv.Itoa(sum)
}

func solution_0002() string {
	sum := 0
	a, b := 1, 2
	for b <= 4000000 {
		if b%2 == 0 {
			sum += b
		}
		a, b = b, a+b
	}
	return strconv.Itoa(sum)
}

func solution_0003() string {
	n := 600851475143
	largest := 0

	// Remove factor 2
	for n%2 == 0 {
		largest = 2
		n /= 2
	}

	// Remove odd factors
	factor := 3
	for factor*factor <= n {
		if n%factor == 0 {
			largest = factor
			n /= factor
		} else {
			factor += 2
		}
	}

	if n > 1 {
		largest = n
	}

	return strconv.Itoa(largest)
}

func isPalindrome(n int) bool {
	original := n
	reversed := 0
	for n > 0 {
		reversed = reversed*10 + n%10
		n /= 10
	}
	return original == reversed
}

func solution_0004() string {
	largest := 0
	for i := 999; i >= 100; i-- {
		// optimization: if i*i <= largest, we can't find a bigger palindrome
		// because j <= i, so i*j <= i*i
		if i*i <= largest {
			break
		}
		for j := i; j >= 100; j-- {
			p := i * j
			if p <= largest {
				break
			}
			if isPalindrome(p) {
				largest = p
			}
		}
	}
	return strconv.Itoa(largest)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return (a * b) / gcd(a, b)
}

func solution_0005() string {
	ans := 1
	for i := 2; i <= 20; i++ {
		ans = lcm(ans, i)
	}
	return strconv.Itoa(ans)
}

func solution_0006() string {
	n := 100
	sum := n * (n + 1) / 2
	sumSq := n * (n + 1) * (2*n + 1) / 6
	diff := sum*sum - sumSq
	return strconv.Itoa(diff)
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

func solution_0007() string {
	count := 1 // We know 2 is prime
	candidate := 1
	for count < 10001 {
		candidate += 2
		if isPrime(candidate) {
			count++
		}
	}
	return strconv.Itoa(candidate)
}

func solution_0008() string {
	s := "73167176531330624919225119674426574742355349194934" +
		"96983520312774506326239578318016984801869478851843" +
		"85861560789112949495459501737958331952853208805511" +
		"12540698747158523863050715693290963295227443043557" +
		"66896648950445244523161731856403098711121722383113" +
		"62229893423380308135336276614282806444486645238749" +
		"30358907296290491560440772390713810515859307960866" +
		"70172427121883998797908792274921901699720888093776" +
		"65727333001053367881220235421809751254540594752243" +
		"52584907711670556013604839586446706324415722155397" +
		"53697817977846174064955149290862569321978468622482" +
		"83972241375657056057490261407972968652414535100474" +
		"82166370484403199890008895243450658541227588666881" +
		"16427171479924442928230863465674813919123162824586" +
		"17866458359124566529476545682848912883142607690042" +
		"24219022671055626321111109370544217506941658960408" +
		"07198403850962455444362981230987879927244284909188" +
		"84580156166097919133875499200524063689912560717606" +
		"05886116467109405077541002256983155200055935729725" +
		"71636269561882670428252483600823257530420752963450"

	largest := int64(0)
	for i := 0; i <= len(s)-13; i++ {
		prod := int64(1)
		for j := 0; j < 13; j++ {
			digit := int64(s[i+j] - '0')
			prod *= digit
		}
		if prod > largest {
			largest = prod
		}
	}
	return strconv.FormatInt(largest, 10)
}

func solution_0009() string {
	for a := 1; a < 1000/3; a++ {
		// b = (500000 - 1000a) / (1000 - a)
		numerator := 500000 - 1000*a
		denominator := 1000 - a

		if numerator%denominator == 0 {
			b := numerator / denominator
			if b > a {
				c := 1000 - a - b
				if c > b {
					return strconv.Itoa(a * b * c)
				}
			}
		}
	}
	return ""
}

func solution_0010() string {
	limit := 2000000
	// sieve[i] == false means i is prime (initially all false)
	sieve := make([]bool, limit)

	sum := int64(0)

	for i := 2; i*i < limit; i++ {
		if !sieve[i] {
			for j := i * i; j < limit; j += i {
				sieve[j] = true
			}
		}
	}

	for i := 2; i < limit; i++ {
		if !sieve[i] {
			sum += int64(i)
		}
	}

	return strconv.FormatInt(sum, 10)
}

func modPow(base, exp, mod int64) int64 {
	res := int64(1)
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp /= 2
	}
	return res
}

func solution_0350() string {
	// f(G, L, N) mod 101^4
	// G = 10^6, L = 10^12, N = 10^18
	// Answer is sum over k=1..L/G of (product over p^alpha || k of ((alpha+1)^N - alpha^N))

	const MOD = 104060401

	// G := int64(1000000)
	// L := int64(1000000000000)
	// LP = L/G = 10^6
	LP := 1000000
	N := int64(1000000000000000000)

	// SPF sieve
	spf := make([]int, LP+1)
	for i := range spf {
		spf[i] = i
	}
	for i := 2; i*i <= LP; i++ {
		if spf[i] == i {
			for j := i * i; j <= LP; j += i {
				if spf[j] == j {
					spf[j] = i
				}
			}
		}
	}

	totalSum := int64(0)

	for k := 1; k <= LP; k++ {
		term := int64(1)
		temp := k
		for temp > 1 {
			p := spf[temp]
			alpha := 0
			for temp%p == 0 {
				temp /= p
				alpha++
			}
			// (alpha+1)^N - alpha^N
			term1 := modPow(int64(alpha+1), N, MOD)
			term2 := modPow(int64(alpha), N, MOD)
			val := (term1 - term2 + MOD) % MOD
			term = (term * val) % MOD
		}
		totalSum = (totalSum + term) % MOD
	}

	return strconv.FormatInt(totalSum, 10)
}

func isPrimeMR(n int64) bool {
	if n < 2 {
		return false
	}
	// ProbablyPrime(0) is 100% accurate for input < 2^64?
	// Documentation says "ProbablyPrime(0) is the same as ProbablyPrime(10) in some versions" or essentially reliable.
	// For competitive programming, Usually ProbablyPrime(0) is enough.
	return big.NewInt(n).ProbablyPrime(0)
}

func solution_0387() string {
	sum := int64(0)
	limit := int64(100000000000000) // 10^14

	var dfs func(current int64, digitSum int64)
	dfs = func(current int64, digitSum int64) {
		// Check if current is Strong Harshad
		if current%digitSum == 0 {
			if isPrimeMR(current / digitSum) {
				// Try to form primes
				for d := int64(0); d <= 9; d++ {
					p := current*10 + d
					if p >= limit {
						continue
					}
					if isPrimeMR(p) {
						sum += p
					}
				}
			}
		} else {
			// If current is not Harshad, it's not RTH, so we shouldn't have reached here via recursion
			// But initial calls 1..9 are Harshad.
			// Recursive calls only happen if nextVal is Harshad.
			// So this branch shouldn't be taken unless called incorrectly.
			return
		}

		// Continue generating RTH
		if current*10 >= limit {
			return
		}

		for d := int64(0); d <= 9; d++ {
			nextVal := current*10 + d
			nextSum := digitSum + d
			if nextVal < limit && nextVal%nextSum == 0 {
				dfs(nextVal, nextSum)
			}
		}
	}

	for i := int64(1); i <= 9; i++ {
		dfs(i, i)
	}

	return strconv.FormatInt(sum, 10)
}

func modInverse(a, m int) int {
	// Extended Euclidean
	m0, x0, x1 := m, 0, 1
	if m == 1 {
		return 0
	}
	for a > 1 {
		q := a / m
		m, a = a%m, m
		x0, x1 = x1-q*x0, x0
	}
	if x1 < 0 {
		x1 += m0
	}
	return x1
}

func solution_0407() string {
	limit := 10000000
	spf := make([]int32, limit+1)
	for i := range spf {
		spf[i] = int32(i)
	}
	for i := 2; i*i <= limit; i++ {
		if spf[i] == int32(i) {
			for j := i * i; j <= limit; j += i {
				if spf[j] == int32(j) {
					spf[j] = int32(i)
				}
			}
		}
	}

	sum := int64(0)

	for n := 1; n <= limit; n++ {
		if n == 1 {
			// M(1): largest a < 1 such that a^2=a mod 1.
			// a=0. 0^2=0=0 mod 1. M(1)=0.
			continue
		}

		// Factor n
		temp := n
		var qs []int
		for temp > 1 {
			p := int(spf[temp])
			q := 1
			for temp%p == 0 {
				q *= p
				temp /= p
			}
			qs = append(qs, q)
		}

		// Find max solution
		// We have system: x = r_i mod q_i, r_i in {0, 1}
		// We want largest x < n.
		// Solution is unique modulo product(q_i) = n.

		maxVal := 0

		// Iterate 2^k solutions
		// Use simple recursion or bitmask
		numFactors := len(qs)

		coeffs := make([]int, numFactors)
		for i, q := range qs {
			Ni := n / q
			// inv(Ni, q)
			yi := modInverse(Ni, q)
			coeffs[i] = (Ni * yi) % n
		}

		// Try all subsets
		// 2^k iterations
		limitMask := 1 << numFactors
		for mask := 0; mask < limitMask; mask++ {
			val := 0
			for i := 0; i < numFactors; i++ {
				if (mask & (1 << i)) != 0 {
					val = (val + coeffs[i])
					if val >= n {
						val -= n
					}
				}
			}
			// val is a solution.
			// We want a < n.
			// val is in [0, n-1] already because modulo n.
			// Max value
			if val > maxVal {
				maxVal = val
			}
		}

		sum += int64(maxVal)
	}

	return strconv.FormatInt(sum, 10)
}
