package controllers

import "strconv"

func toAtomic(amount string) (uint64, error) {
	toFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}
	toAtomic := uint64(toFloat * 1000000000000)
	return toAtomic, nil
}
