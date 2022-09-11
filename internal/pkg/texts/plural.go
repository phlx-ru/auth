package texts

// Plural возвращает корректную форму множественного числа: Plural(11, "рыба", "рыбы", "рыб") -> "рыб"
func Plural[T int | int64](count T, formSingular, formPluralWeak, formPluralStrong string) string {
	if count%10 == 1 && count%100 != 11 {
		return formSingular
	}
	if count%10 >= 2 && count%10 <= 4 && (count%100 < 10 || count%100 >= 20) {
		return formPluralWeak
	}
	return formPluralStrong
}
