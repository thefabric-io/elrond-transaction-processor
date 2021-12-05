package processor

type Nonce int

func (n Nonce) Equals(n2 Nonce) bool {
	return n == n2
}

func (n Nonce) IsGreaterThan(n2 Nonce) bool {
	return n > n2
}

func (n Nonce) Subtract(d Nonce) Nonce {
	return n - d
}

func (n Nonce) Increment() Nonce {
	return n + 1
}

func (n Nonce) Decrement() Nonce {
	r := n - 1
	if r < 0 {
		return 0
	}

	return r
}
