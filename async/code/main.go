package main

import (
	"fmt"
	"strings"
)

func checkStrDivider(str1, str2, strDivider string) bool {

	if strDivider == "" {
		return false
	}

	var isIntOffset1, isIntOffset2 bool

	for l, r := 0, len(strDivider); r <= len(str1) || r <= len(str2); l, r = r, r+len(strDivider) {
		if r == len(str1) {
			isIntOffset1 = true
		}
		if r == len(str2) {
			isIntOffset2 = true
		}
		if r <= len(str1) && str1[l:r] != strDivider {
			return false
		}
		if r <= len(str2) && str2[l:r] != strDivider {
			return false
		}
	}

	if !isIntOffset1 || !isIntOffset2 {
		return false
	} else {
		return true
	}

}

func gcdOfStrings(str1 string, str2 string) string {
	var shortestStr string

	if len(str1) < len(str2) {
		shortestStr = str1
	} else {
		shortestStr = str2
	}

	var strNOD string

	for curId := range shortestStr {
		tmp := shortestStr[:curId+1]
		if checkStrDivider(str1, str2, tmp) {
			strNOD = tmp
		}
	}

	return strNOD
}

func canPlaceFlowers(flowerbed []int, n int) bool {
	placesNum := 0
	for i := 0; i != len(flowerbed); i++ {
		switch i {
		case 0:
			if flowerbed[i+1] == 0 {
				flowerbed[i] = 1
				placesNum++
			}
		case len(flowerbed) - 1:
			if flowerbed[i-1] == 0 {
				flowerbed[i] = 1
				placesNum++
			}
		default:
			if flowerbed[i-1] == 0 && flowerbed[i+1] == 0 {
				flowerbed[i] = 1
				placesNum++

			}
		}
	}
	if placesNum >= n {
		return true
	} else {
		return false
	}
}

func main() {
	s := []string{"foo", "bar"}
	fmt.Println(len(strings.Join(s, " ")))
	fmt.Printf("%+v", strings.Split("fo  o", " "))
}
