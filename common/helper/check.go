package helper

func IsDigit(ch rune) bool {
	if ch == '0' || ch == '1' || ch == '2' || ch == '3' || ch == '4' || ch == '5' || ch == '6' || ch == '7' || ch == '8' || ch == '9' {
		return true
	}
	return false
}

func IsSymbol(ch rune) bool {
	if ch == ';' || ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' || ch == '(' || ch == ')' || ch == ',' {
		return true
	}
	return false
}

func IsWhiteSpace(ch rune) bool {
	if ch == ' ' || ch == '\t' || ch == '\n' {
		return true
	}
	return false
}
