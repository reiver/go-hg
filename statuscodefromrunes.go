package hg

func statuscodeFromRunes(mostsignificant rune, leastsignificant rune) (int, bool) {
	ms := (mostsignificant  - rune('0'))
	ls := (leastsignificant - rune('0'))

	if ms < 0 || 9 < ms {
		return 0, false
	}
	if ls < 0 || 9 < ls {
		return 0, false
	}

	return int(10*ms + ls), true
}

