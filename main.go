package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type Document struct {
	ID      int
	Content string
	Score   float64
}

type InvertedIndex map[string][]int

type SearchEngine struct {
	index        InvertedIndex
	documents    []Document
	avgDocLength float64
	k1, k2       float64
}

func BuildInvertedIndex(documents []Document) InvertedIndex {
	index := make(InvertedIndex)

	for _, doc := range documents {
		tokens := strings.Fields(strings.ToLower(doc.Content))

		for _, token := range tokens {
			if _, ok := index[token]; !ok {
				index[token] = make([]int, 0)
			}
			index[token] = append(index[token], doc.ID)
		}
	}

	return index
}

func (se *SearchEngine) CalculateTFIDFScore(tokens []string) map[int]float64 {
	scores := make(map[int]float64)

	for _, token := range tokens {
		if docSet, ok := se.index[token]; ok {
			idf := math.Log(float64(len(se.documents)) / float64(len(docSet)))
			for _, docID := range docSet {
				tf := float64(strings.Count(strings.ToLower(se.documents[docID].Content), token))
				scores[docID] += tf * idf
			}
		}
	}

	return scores
}

func (se *SearchEngine) CalculateBM25Score(tokens []string) map[int]float64 {
	scores := make(map[int]float64)

	for _, token := range tokens {
		if docSet, ok := se.index[token]; ok {
			idf := math.Log(float64(len(se.documents)-len(docSet))+0.5) / (float64(len(docSet)) + 0.5)
			for _, docID := range docSet {
				tf := float64(strings.Count(strings.ToLower(se.documents[docID].Content), token))
				dl := float64(len(strings.Fields(strings.ToLower(se.documents[docID].Content))))
				numerator := (se.k1 + 1) * tf * (se.k1 + 1) / (tf + se.k1*(1.0-se.k2+se.k2*dl/se.avgDocLength))
				denominator := tf + se.k1*(1.0-se.k2+se.k2*dl/se.avgDocLength)
				scores[docID] += idf * numerator / denominator
			}
		}
	}

	return scores
}

func (se *SearchEngine) Search(query string) []Document {
	tokens := strings.Fields(strings.ToLower(query))
	scores := se.CalculateTFIDFScore(tokens)
	// or, to use bm25 scoring algorithm:
	// scores := se.CalculateBM25Score(tokens)

	var results []Document
	for docID, score := range scores {
		results = append(results, Document{ID: docID, Content: se.documents[docID].Content, Score: score})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) > 10 {
		results = results[:10]
	}
	return results
}

func main() {
	documents := []Document{
		{ID: 0, Content: "Lorem ipsum blah blah fox"},
		{ID: 1, Content: "The quick brown fox jumped over the lazy dog. The dog slept peacefully."},
		{ID: 2, Content: "I have a dream that one day this nation will rise up and live out the true meaning of its creed: 'We hold these truths to be self-evident, that all men are created equal.'"},
		{ID: 3, Content: "To be, or not to be, that is the question: Whether 'tis nobler in the mind to suffer The slings and arrows of outrageous fortune, Or to take arms against a sea of troubles And by opposing end them."},
		{ID: 4, Content: "In a hole in the ground there lived a hobbit. Not a nasty, dirty, wet hole, filled with the ends of worms and an oozy smell, nor yet a dry, bare, sandy hole with nothing in it to sit down on or to eat: it was a hobbit-hole, and that means comfort."},
		{ID: 5, Content: "The only way to do great work is to love what you do. If you haven't found it yet, keep looking. Don't settle. As with all matters of the heart, you'll know when you find it."},
		{ID: 6, Content: "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife."},
		{ID: 7, Content: "It was the best of times, it was the worst of times, it was the age of wisdom, it was the age of foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of Darkness, it was the spring of hope, it was the winter of despair."},
		{ID: 8, Content: "Two households, both alike in dignity, In fair Verona, where we lay our scene, From ancient grudge break to new mutiny, Where civil blood makes civil hands unclean."},
		{ID: 9, Content: "Once upon a time in a far-off land, there was a princess who was very beautiful and very kind, but also very sad."},
		{ID: 10, Content: "It is not in the stars to hold our destiny but in ourselves."},
		{ID: 11, Content: "In the beginning God created the heaven and the earth. And the earth was without form, and void; and darkness was upon the face of the deep. And the Spirit of God moved upon the face of the waters."},
		{ID: 12, Content: "There are known knowns; there are things we know we know. We also know there are known unknowns; that is to say we know there are some things we do not know. But there are also unknown unknowns – the ones we don't know we don't know."},
		{ID: 13, Content: "When I consider how my light is spent Ere half my days in this dark world and wide, And that one talent which is death to hide Lodg'd with me useless, though my soul more bent To serve therewith my Maker, and present My true account, lest he returning chide;"},
		{ID: 14, Content: "I wandered lonely as a cloud That floats on high o'er vales and hills, When all at once I saw a crowd, A host, of golden daffodils; Beside the lake, beneath the trees, Fluttering and dancing in the breeze."},
		{ID: 15, Content: "Do not go gentle into that good night, Old age should burn and rave at close of day; Rage, rage against the dying of the light."},
		{ID: 16, Content: "The sun was shining on the sea, Shining with all his might: He did his very best to make The billows smooth and bright."},
		{ID: 17, Content: "In Xanadu did Kubla Khan A stately pleasure-dome decree: Where Alph, the sacred river, ran Through caverns measureless to man Down to a sunless sea."},
		{ID: 18, Content: "I celebrate myself, and sing myself, And what I assume you shall assume, For every atom belonging to me as good belongs to you."},
		{ID: 19, Content: "The love that moves the sun and all the stars."},
		{ID: 20, Content: "It was a bright cold day in April, and the clocks were striking thirteen. Winston Smith, his chin nuzzled into his breast in an effort to escape the vile wind, slipped quickly through the glass doors of Victory Mansions, though not quickly enough to prevent a swirl of gritty dust from entering along with him."},
		{ID: 21, Content: "It was a pleasure to burn. It was a special pleasure to see things eaten, to see things blackened and changed."},
		{ID: 22, Content: "The human race, to which so many of my readers belong, has been playing at children's games from the beginning, and will probably do it till the end, which is a nuisance for the few people who grow up. And one of the games to which it is most attached is called 'Keep to-morrow dark,' and which is also sometimes called 'Cheat the Prophet.'"},
		{ID: 23, Content: "Happy families are all alike; every unhappy family is unhappy in its own way."},
		{ID: 24, Content: "I am an invisible man. No, I am not a spook like those who haunted Edgar Allan Poe; nor am I one of your Hollywood-movie ectoplasms. I am a man of substance, of flesh and bone, fiber and liquids—and I might even be said to possess a mind. I am invisible, understand, simply because people refuse to see me."},
		{ID: 25, Content: "It was a dark and stormy night; the rain fell in torrents, except at occasional intervals, when it was checked by a violent gust of wind which swept up the streets (for it is in London that our scene lies), rattling along the housetops, and fiercely agitating the scanty flame of the lamps that struggled against the darkness."},
		{ID: 26, Content: "The sky above the port was the color of television, tuned to a dead channel."},
		{ID: 27, Content: "All children, except one, grow up. They soon know that they will grow up, and the way Wendy knew was this. One day when she was two years old she was playing in a garden, and she plucked another flower and ran with it to her mother. I suppose she must have looked rather delightful, for Mrs. Darling put her hand to her heart and cried, 'Oh, why can't you remain like this for ever!' This was all that passed between them on the subject, but henceforth Wendy knew that she must grow up. You always know after you are two. Two is the beginning of the end."},
		{ID: 28, Content: "As Gregor Samsa awoke one morning from uneasy dreams he found himself transformed in his bed into a gigantic insect."},
		{ID: 29, Content: "Call me Ishmael. Some years ago—never mind how long precisely—having little or no money in my purse, and nothing particular to interest me on shore, I thought I would sail about a little and see the watery part of the world."},
		{ID: 30, Content: "It was the day my grandmother exploded."},
	}

	index := BuildInvertedIndex(documents)
	docLength := 0.
	for _, doc := range documents {
		docLength += float64(len(doc.Content))
	}
	searchEngine := SearchEngine{index: index, documents: documents, avgDocLength: docLength / float64(len(documents)), k1: 1.2, k2: 0.75}

	for {
		fmt.Print("Enter a search query: ")
		query, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		query = strings.TrimSpace(query)
		if query == "" {
			break
		}
		results := searchEngine.Search(query)
		fmt.Printf("%d results for query '%s':\n", len(results), query)
		for _, result := range results {
			fmt.Printf("- %s (score=%.2f)\n", result.Content, result.Score)
		}
	}
}
