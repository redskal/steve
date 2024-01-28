/*
Released under YOLO licence. Idgaf what you do.
*/
package azrecon

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode"
)

const (
	Name = iota
	Mask = iota
)

type ScoreKeeper struct {
	Count  int
	Values []string
}

// GetNameCombinations gets all combinations of elements in s,
// split by a hyphen.
func GetNameCombinations(s []string) []string {
	doubleSlices, optionals := getColumns(s, Name)
	colCount := len(doubleSlices) - len(optionals)

	var combos []string
	for _, v := range optionals {
		combos = append(combos, getAllCombinations(doubleSlices[:v])...)
	}
	combos = append(combos, getAllCombinations(doubleSlices)...)

	for i := 0; i < len(combos); i++ {
		if strings.Count(combos[i], "-") < colCount-1 {
			if i <= len(combos) {
				combos = slices.Delete(combos, i, i+1)
				i--
			}
		}
	}
	return combos
}

// GetMaskCombinations gets all combinations of hashcat-style masks
// for the given []string.
func GetMaskCombinations(s []string) []string {
	doubleSlices, optionals := getColumns(s, Mask)
	colCount := len(doubleSlices) - len(optionals)

	var combos []string
	for _, v := range optionals {
		combos = append(combos, getAllCombinations(doubleSlices[:v])...)
	}
	combos = append(combos, getAllCombinations(doubleSlices)...)

	for i := 0; i < len(combos); i++ {
		if strings.Count(combos[i], "-") < colCount-1 {
			if i <= len(combos) {
				combos = slices.Delete(combos, i, i+1)
				// if we don't decrement here we'll get an index error later
				i--
			}
		}
	}
	return combos
}

// getAllCombinations returns all combinations of each element within each slice
// with the elements of each other slice.
func getAllCombinations(ss [][]string) []string {
	permutations := []string{}
	if len(ss) == 1 {
		// TODO: this seem ok?
		return ss[0]
	}

	for i, s := range ss {
		for _, t := range s {
			if i+1 < len(ss) {
				subpermutations := getAllCombinations(ss[i+1:])
				for _, v := range subpermutations {
					permutation := fmt.Sprintf("%v-%v", t, v)
					permutations = append(permutations, permutation)
				}
			}
		}
	}

	return permutations
}

// getColumns returns a [][]string sorted by number of occurences
// of each key. Provides a list of low-occurence columns in optionals
func getColumns(r []string, method int) ([][]string, []int) {
	partsCount := getPartsCount(r)

	// this works until they add shit to the beginning of the resource name...
	partsSlices := make([][]string, partsCount)

	// split by hyphen, and make a slice (column) for each parts values
	for _, v := range r {
		parts := strings.Split(v, "-")
		for ii, part := range parts {
			partsSlices[ii] = append(partsSlices[ii], part)
		}
	}
	ret := make([][]string, len(partsSlices))
	optionals := []int{}
	for i, v := range partsSlices {
		ks := keepScore(v)
		keys := sortKeys(ks)
		for ii := 0; ii < len(ks); ii++ {
			switch method {
			case Name:
				// ret[i] = keys by order of occurence
				ret[i] = append(ret[i], keys[ii])
			case Mask:
				// add a hashcat mask, but de-duplicate for this one
				mask := createMask(keys[ii])
				if !slices.Contains(ret[i], mask) {
					ret[i] = append(ret[i], mask)
				}
			}
		}
		// if there's only one occurence, we'll add the index to optionals
		if len(ks) == 1 {
			optionals = append(optionals, i)
		}
	}
	return ret, optionals
}

// createMask generates a hashcat-style mask for the provided string
func createMask(s string) string {
	var out string
	for _, c := range s {
		switch {
		case unicode.IsDigit(c):
			out += "?d"
		case unicode.IsLetter(c):
			out += "?l"
		default:
			out += string(c)
		}
	}
	return out
}

// sortKeys will sort keys ordered by number of occurences
func sortKeys(m map[string]ScoreKeeper) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]].Count > m[keys[j]].Count
	})
	return keys
}

// getPartsCount gets the longest number of parts from each string in s
// when split by '-'
func getPartsCount(s []string) int {
	var highestPartCount int
	for _, v := range s {
		parts := strings.Split(v, "-")
		// highestPartCount to create map later
		if len(parts) > highestPartCount {
			highestPartCount = len(parts)
		}
	}
	return highestPartCount
}

// keepScore counts the number of occurences a string appears, and
// keeps track of strings with strong likeness.
func keepScore(parts []string) map[string]ScoreKeeper {
	test := map[string]ScoreKeeper{}
	for _, part := range parts {
		// is it like an existing key? yes...
		for key := range test {
			tmp := levenshtein([]rune(key), []rune(part))
			// tmp < len(key) ...or... tmp <= len(key) ?!?!?
			if tmp <= len(key) && tmp > 0 {
				// similar
				if wtf, ok := test[key]; ok && !slices.Contains(wtf.Values, part) {
					// wtf...a weird quirk of Golang
					wtf.Values = append(wtf.Values, part)
					test[key] = wtf
				}
				//} else if tmp > len(key) && tmp > 0 {
				// unlikely to exist
				//	if _, ok := test[part]; !ok {
				//		test[part] = ScoreKeeper{
				//			Count: 1,
				//		}
				//	}
			} else if wtf, ok := test[part]; ok {
				// it exists for sure...
				wtf.Count += 1
				test[part] = wtf
			}
		}
		// otherwise it remains empty as we don't create any keys
		if _, ok := test[part]; !ok {
			test[part] = ScoreKeeper{
				Count: 1,
			}
		}
	}
	return test
}

// levenshtein is robbed from here:
// https://github.com/hakluke/hakoriginfinder/blob/main/hakoriginfinder.go#L21
func levenshtein(str1, str2 []rune) int {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return column[s1len]
}

// minimum is robbed from here:
// https://github.com/hakluke/hakoriginfinder/blob/main/hakoriginfinder.go#L47
func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}
