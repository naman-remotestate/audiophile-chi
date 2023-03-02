package utils

import "strconv"

func GetLimitAndPage(limit, page string) (int64, int64, error) {
	newLimit := 10
	newPage := 0
	var err error
	if page != "" {
		newPage, err = strconv.Atoi(page)
		if err != nil {
			return -1, -1, err
		}
	} else {
		newLimit, err = strconv.Atoi(limit)
		if err != nil {
			return -1, -1, err
		}
		newPage, err = strconv.Atoi(page)
		if err != nil {
			return -1, -1, err
		}
	}

	return int64(newLimit), int64(newPage), nil
}
