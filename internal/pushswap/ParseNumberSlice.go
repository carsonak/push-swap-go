package pushswap

import (
	"fmt"
	"strconv"
)

type empty struct{}

func ParseNumberSlice(numStrings []string, allowDups bool) ([]float64, error) {
	numList := make([]float64, 0, len(numStrings))
	numSeen := map[float64]empty{}

	for _, numStr := range numStrings {
		n, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing %q: %v", numStr, err)
		}

		if !allowDups {
			_, exists := numSeen[n]
			if exists {
				return nil, fmt.Errorf("duplicate number %.f", n)
			}

			numSeen[n] = empty{}
		}

		numList = append(numList, n)
	}

	return numList, nil
}
