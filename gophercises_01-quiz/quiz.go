package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io"
  "os"
  "bufio"
  "strings"
  "time"
)

type Question struct {
  question string
  answer string
}


func main() {

  filename := flag.String("filename", "problem.csv", "CSV File that contains quiz questions")
  flag.Parse()

  f, err := os.Open(*filename)
  if err != nil {
    return
  }
  questions, err := readCSV(f)

  if err != nil {
    fmt.Println(err.Error())
    return
  }
  if questions == nil {
    return
  }

  score, err := ask(questions)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  fmt.Printf("Score: %d/%d\n", score, len(questions))
}

func readCSV(f io.Reader) ([]Question, error) {
  allQuestions, err := csv.NewReader(f).ReadAll()
  if err != nil {
    return nil, err
  }
  numQues := len(allQuestions)
  if numQues == 0 {
    return nil, fmt.Errorf("No questions in file")
  }

  var q []Question
  for _, line := range allQuestions {
    ques := Question{}
    ques.question = line[0]
    ques.answer = line[1]
    q = append(q, ques)
  }

  return q, nil
}


func getInput(input chan string) {
  for {
    in := bufio.NewReader(os.Stdin)
    ans, _ := in.ReadString('\n')
    input <- ans
  }
}

func ask(questions []Question) (int, error) {

  timer := time.NewTimer(time.Duration(10) * time.Second)

  done := make(chan string)
  go getInput(done)

  score := 0
  for _, q := range questions {
    fmt.Printf("%s = ?\n", q.question)

    select {
    case <- timer.C:
      fmt.Printf("Oops you ran out of time!\n")
      return score, nil
    case ans := <-done:
      if strings.Compare(strings.Trim(strings.ToLower(ans), "\n"), q.answer) == 0 {
        score++
      }
    }
  }
  return score, nil
}
