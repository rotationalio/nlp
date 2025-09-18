package main

/*
##################################################
## USAGE: `go run docs/examples/readability.go` ##
##################################################
*/

import (
	"fmt"
	"os"
	"strings"

	"go.rtnl.ai/nlp/stats"
	"go.rtnl.ai/nlp/text"
)

// An example for how to perform readability scoring comparisons where there is
// human written text and an LLM rewrite.
func RewrittenExample() {
	// Human written
	humanText, _ := text.New("Added Flesch-Kincaid reading ease and grade level metrics and some related functionality: a sentence segmentizer, a word syllable tokenizer, an optimized whitespace tokenizer (and additional regex for the same result with the RegexTokenizer), and some additional helpful functions on the Text API that were related. I also started to add some useful functions in the ngrams package which I ended up not using in this project, however I didn't feel like I should delete them as they may come in handy at a later date.")

	// Rewritten by an LLM
	llmText, _ := text.New("I added tools to measure how easy it is to read text (like Flesch-Kincaid scores) and what grade level it's at. I also added features like splitting text into sentences, breaking words into syllables, and handling spaces more efficiently. There's also a regular expression version of the space tokenizer for the same result. I added some extra functions on the Text API that were useful. I also started adding tools to the ngrams package, but didn't use them in this project --- they might be helpful later.")

	// Get scores
	humanEase, _ := humanText.FleschKincaidReadingEase()
	llmEase, _ := llmText.FleschKincaidReadingEase()

	// Compare
	if humanEase > llmEase {
		fmt.Printf("The human written text was simpler than the LLM re-written text: %.2f < %.2f\n", humanEase, llmEase)
	} else {
		fmt.Printf("The LLM re-written text was simpler than the human written text: %.2f > %.2f\n", humanEase, llmEase)
	}

	// Print some other stats
	fmt.Printf(
		"Words/Sentence: Human: %.2f / LLM: %.2f\n",
		float64(humanText.WordCount())/float64(humanText.SentenceCount()),
		float64(llmText.WordCount())/float64(llmText.SentenceCount()),
	)
	fmt.Printf(
		"Syllables/Word: Human: %.2f / LLM: %.2f\n",
		float64(humanText.SyllableCount())/float64(humanText.WordCount()),
		float64(llmText.SyllableCount())/float64(llmText.WordCount()),
	)

	// Newline
	fmt.Println()
}

// An example for how to perform readability scoring comparisons where there is
// two paragraphs written at different levels by an LLM.
func ELIFive_and_ELIHighSchool() {
	// "Explain electricity to a 5 year old in a single paragraph"
	eliFiveText, _ := text.New("Using a water slide as an analogy, electricity is the flow of tiny particles called electrons. A battery or power plant is like a big pump that pushes the electrons down a long, winding \"water slide\" made of wires. When you flip a switch, it's like opening a gate that lets the electrons rush through the slide, carrying energy that can light up a light bulb, turn on a TV, or power your toys.")

	// "Explain electricity to a HS senior in a single paragraph"
	eliHSText, _ := text.New("Electricity is the movement of electrons, which are negatively charged particles, through a conductor like a metal wire. This flow is known as electric current, and it requires a closed, continuous loop called a circuit to travel. The \"push\" that forces electrons to move is called voltage, or potential difference, and can be compared to the pressure in a water pipe. Resistance, measured in ohms, is the opposition to this flow, acting like a narrow section or kink in the pipe that slows down the current. Power plants generate electricity by converting other forms of energy (such as the kinetic energy from steam, wind, or flowing water) into electrical energy, often by spinning a turbine that drives a generator.")

	// Get grades
	eliFiveGrade, _ := eliFiveText.FleschKincaidGradeLevel()
	eliHSGrade, _ := eliHSText.FleschKincaidGradeLevel()

	// Compare
	if eliFiveGrade < eliHSGrade {
		fmt.Printf("ELI 5 grade is lower than the ELI High School grade: %.2f < %.2f\n", eliFiveGrade, eliHSGrade)
	} else {
		fmt.Printf("ELI High School grade is lower than the ELI 5 grade: %.2f > %.2f\n", eliFiveGrade, eliHSGrade)
	}

	// Print some other stats
	fmt.Printf(
		"Words/Sentence: ELI5: %.2f / ELIHS: %.2f\n",
		float64(eliFiveText.WordCount())/float64(eliFiveText.SentenceCount()),
		float64(eliHSText.WordCount())/float64(eliHSText.SentenceCount()),
	)
	fmt.Printf(
		"Syllables/Word: ELI5: %.2f / ELIHS: %.2f\n",
		float64(eliFiveText.SyllableCount())/float64(eliFiveText.WordCount()),
		float64(eliHSText.SyllableCount())/float64(eliHSText.WordCount()),
	)

	// Newline
	fmt.Println()
}

// An example for how to perform readability scoring to get statistics on an
// array of text "chunks".
func StatsExample() {
	// Load some movie reviews from a file
	data, _ := os.ReadFile("docs/examples/data/movie_reviews.txt")

	// Track the stats
	easeStats := new(stats.Statistics)
	gradeStats := new(stats.Statistics)
	wordsPerSentenceStats := new(stats.Statistics)
	syllablesPerWordStats := new(stats.Statistics)

	// Collect the stats for the reviews
	for review := range strings.Lines(string(data)) {
		if review != "" {
			myText, _ := text.New(string(review))

			ease, _ := myText.FleschKincaidReadingEase()
			easeStats.Update(ease)

			grade, _ := myText.FleschKincaidGradeLevel()
			gradeStats.Update(grade)

			wordsPerSentenceStats.Update(float64(myText.WordCount()) / float64(myText.SentenceCount()))

			syllablesPerWordStats.Update(float64(myText.SyllableCount()) / float64(myText.WordCount()))
		}
	}

	// Print stats info
	fmt.Println("easeStats:")
	fmt.Println(easeStats.Serialize())
	fmt.Println() // Newline

	fmt.Println("gradeStats:")
	fmt.Println(gradeStats.Serialize())
	fmt.Println() // Newline

	fmt.Println("wordsPerSentenceStats:")
	fmt.Println(wordsPerSentenceStats.Serialize())
	fmt.Println() // Newline

	fmt.Println("syllablesPerWordStats:")
	fmt.Println(syllablesPerWordStats.Serialize())
	fmt.Println() // Newline
}

func main() {
	RewrittenExample()
	ELIFive_and_ELIHighSchool()
	StatsExample()
}
