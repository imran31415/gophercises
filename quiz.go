package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type QuizAnswer struct {
	SubmittedAnswer string
	IsCorrect bool
	QuizQuestion
}

type QuizQuestion struct {
	Question string
	Answer string
}

//flag vars
var filename string
var shuffleQuestions bool
var timerLength int

func init(){
	flag.StringVar(&filename , "filename", "./problems.csv", "The csv filename of questions")
	flag.BoolVar(&shuffleQuestions, "shuffle", false, "Whether to shuffle the order the questions are asked")
	flag.IntVar(&timerLength, "timer", 5, "Set a timer to the quiz, default is 5")
	flag.Parse()
}

// Ask questions from a csv within a given time limit
func main() {
	timer := time.NewTimer(time.Duration(timerLength) * time.Second)
	// 'answers' gets updated in the go routine on a per question basis.
	var answers []QuizAnswer
	// 'numQuestions' gets updated when the csv is parsed
	var numQuestions int
	// exit early if we hit the time limit
	go func() {
		<- timer.C
		fmt.Printf("\n Time limit of %d Expired \n",timerLength )
		finishQuiz(&answers, numQuestions)
		return
	}()
	runQuiz(&answers, filename, shuffleQuestions, &numQuestions)
	// if all questions are answered in the time limit
	finishQuiz(&answers, numQuestions)
}

func runQuiz(submittedAnswers *[]QuizAnswer, filename string, shuffleQuestions bool, numQuestions *int)  {
	lines, err := parseCsv(filename, shuffleQuestions)
	if err != nil {
		panic(err)
	}
	*numQuestions = len(lines)
	// Loop through lines questions and ask each one and update the submittedAnswers slice
	for i, line := range lines {
		question := QuizQuestion{
			Question: line[0],
			Answer: line[1],
		}
		quizAnswer, err := askQuestion(question, i)
		if err != nil{
			panic(err)
		}
		*submittedAnswers = append(*submittedAnswers, quizAnswer)
	}

}

func finishQuiz(answers *[]QuizAnswer, numQuestions int){
	correct, answered := evaluateAnswers(answers)
	fmt.Printf("\n Total Questions in Quiz  %d \n", numQuestions)
	fmt.Printf("\n Total Answers Attempted %d \n", answered)
	fmt.Printf("\n Total Correct Answers %d \n", correct)
	os.Exit(3)
}

func parseCsv(filename string, shuffleQuestions bool) (questionArray [][]string, err error){
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	// convert to a nested string slice to enable shuffling
	for _, line := range lines {
		questionArray = append(questionArray, line)

	}
	if shuffleQuestions == true {
		Shuffle(questionArray)
		fmt.Println("Shuffled!")

	}
	return questionArray, nil

}



func askQuestion(question QuizQuestion, index int)  (quizAnswer QuizAnswer, err error){
	isCorrect := false
	fmt.Println("Question " + string(index))
	fmt.Println(question.Question)
	answer := getUserInput()
	fmt.Printf("\n You Answered: %s \n", answer)
	fmt.Printf("Correct Answer %s \n", question.Answer)

	if answer == question.Answer {
		isCorrect = true
		fmt.Println("Correct!")
	} else {
		isCorrect = false
	}
	quizAnswer = QuizAnswer {
		IsCorrect: isCorrect,
		SubmittedAnswer: answer,
		QuizQuestion: question,
	}
	return quizAnswer, nil
}

// Keep trying to get user input until it is validly parsed
func getUserInput() (inputAnswer string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Answer: ")
	inputAnswer, err := reader.ReadString('\n')
	// Re ask the question if we get invalid input we cant parse.
	if err != nil {
		getUserInput()
	}
	inputAnswer = strings.TrimSpace(inputAnswer)
	return inputAnswer
}


func evaluateAnswers(answers *[]QuizAnswer) (int, int)  {
	totalCorrect := 0
	totalAnswered := 0
	for _, answer := range *answers{
		if answer.IsCorrect {
			totalCorrect +=1
		}
		totalAnswered +=1
	}
	return totalCorrect, totalAnswered
}

// utility for shuffling a slice
func Shuffle(slice [][]string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}