package dayThree

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
)

func isSymbol(c rune) bool {
	return !unicode.IsDigit(c) && !isSpace(c)
}

func isSpace(c rune) bool {
	return c == '.'
}

func isValidPart(i int, precedingRow, thisRow, followingRow string) bool {
	for _, row := range []string{precedingRow, thisRow, followingRow} {
		if row == "" {
			continue
		}
		if i-1 > 0 && isSymbol(rune(row[i-1])) {
			return true
		}
		if isSymbol(rune(row[i])) {
			return true
		}
		if i+2 < len(row) && isSymbol(rune(row[i+1])) {
			return true
		}
	}

	return false
}

func ScanPuzzle(path string) []int {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	rows := []string{}

	for scanner.Scan() {
		// Not my favorite, but I'm getting lazy today
		rows = append(rows, scanner.Text())
	}

	parts := []int{}
	candidatePartIsValid := false
	previousRow := rows[0]
	nextRow := rows[0]

	for i, row := range rows {
		candidatePart := 0
		if i > 0 {
			previousRow = rows[i-1]
		}
		if i+1 < len(rows) {
			nextRow = rows[i+1]
		}
		for j, c := range row {
			atEnd := j == len(row)-1

			if unicode.IsDigit(c) {
				candidatePart = candidatePart*10 + int(c-'0')
				candidatePartIsValid = candidatePartIsValid || isValidPart(j, previousRow, row, nextRow)
				if !atEnd {
					continue
				}
			}

			if candidatePart != 0 || atEnd {
				if candidatePartIsValid {
					parts = append(parts, candidatePart)
				} else {
					//fmt.Println(i, j, candidatePart)
				}

				candidatePart = 0
				candidatePartIsValid = false
			}
		}
	}

	return parts
}

func PartOne(puzzleFile string) {
	// parse all lines
	// for each line, for each position: check if number is adjacent to a symbol
	// return part numbers slice
	// accumulate
	// return sum

	parts := ScanPuzzle(puzzleFile)

	//fmt.Println(parts)

	partsSum := 0
	for _, p := range parts {
		partsSum += p
	}
	fmt.Println("Parts Sum:", partsSum)
}

func PartTwo(puzzleFile string) {
}
