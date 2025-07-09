package indexerutil

import "sort"

func CountWords(text []string) map[string]int {
	wordsFrequency := make(map[string]int)

	for _, word := range text {
		_, exists := wordsFrequency[word]
		if !exists {
			wordsFrequency[word] = 1
		} else {
			wordsFrequency[word] += 1
		}
	}

	return wordsFrequency
}

func MostCommonWords(wordsFrequency map[string]int, n int) map[string]int {
	type kv struct {
		Key   string
		Value int
	}

	var freqList []kv
	for key, value := range wordsFrequency {
		freqList = append(freqList, kv{key, value})
	}

	// sort by descending
	sort.Slice(freqList, func(i, j int) bool {
		// if same frequency on two words, sort alphabetically
		if freqList[i].Value == freqList[j].Value {
			return freqList[i].Key < freqList[j].Key
		}

		return freqList[i].Value > freqList[j].Value
	})

	topWords := make(map[string]int)
	for i := 0; i < len(freqList) && i < n; i++ {
		topWords[freqList[i].Key] = freqList[i].Value
	}

	return topWords
}
