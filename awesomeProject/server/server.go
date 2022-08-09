package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
)

type quiz struct {
	Qid               int
	username          string
	quizName          string
	time              int
	numberOfQuestions int
	result            int
}

func main() {
	db, err := sql.Open("mysql", "root:Aa@12345@tcp(127.0.0.1:3306)/Azmoon")
	if err != nil {
		Errorr(err)
	}
	err = db.Ping()
	if err != nil {
		Errorr(err)
	}
	defer db.Close()
	server, _ := net.Listen("tcp", "localhost:8000")
	defer server.Close()
	if err != nil {
		Errorr(err)
	}

	fmt.Println("server is ready!")
	for {
		client, err := server.Accept()
		if err != nil {
			Errorr(err)
		}
		go handle(client, db)
	}
}
func Errorr(text error) {
	fmt.Println(text)
	os.Exit(1)
}
func handle(client net.Conn, db *sql.DB) {
	message, err := bufio.NewReader(client).ReadString('\n')
	if err != nil {
		Errorr(err)
	}
	username := strings.TrimSuffix(message, "\n")
	for username == "server" {
		message, err = bufio.NewReader(client).ReadString('\n')
		if err != nil {
			Errorr(err)
		}
		message = strings.TrimSuffix(message, "\n")
		if message == "create" {
			var q quiz
			message, err = bufio.NewReader(client).ReadString('\n')
			if err != nil {
				Errorr(err)
			}
			message = strings.TrimSuffix(message, "\n")
			q.Qid, _ = strconv.Atoi(message)
			message, err = bufio.NewReader(client).ReadString('\n')
			if err != nil {
				Errorr(err)
			}
			message = strings.TrimSuffix(message, "\n")
			q.username = message
			message, err = bufio.NewReader(client).ReadString('\n')
			if err != nil {
				Errorr(err)
			}
			message = strings.TrimSuffix(message, "\n")
			q.quizName = message
			message, err = bufio.NewReader(client).ReadString('\n')
			if err != nil {
				Errorr(err)
			}
			message = strings.TrimSuffix(message, "\n")
			q.time, _ = strconv.Atoi(message)
			message, err = bufio.NewReader(client).ReadString('\n')
			if err != nil {
				Errorr(err)
			}
			message = strings.TrimSuffix(message, "\n")
			q.numberOfQuestions, _ = strconv.Atoi(message)
			db.Query("INSERT INTO quiz(Qid,username,quizName,time,numberOfQuestions,result) VALUES(?,?,?,?,?,?)", q.Qid, q.username, q.quizName, q.time, q.numberOfQuestions, -100)
			for i := 1; i <= q.numberOfQuestions; i++ {
				message, err = bufio.NewReader(client).ReadString('\n')
				if err != nil {
					Errorr(err)
				}
				message = strings.TrimSuffix(message, "\n")
				a, _ := strconv.Atoi(message)
				db.Query("INSERT INTO answers(Qid,id,answer) VALUES(?,?,?)", q.Qid, i, a)
			}
		}
	}
	for {
		message, err = bufio.NewReader(client).ReadString('\n')
		if err != nil {
			Errorr(err)
		}
		arr := strings.Split(message, "#")
		if arr[0] == "showQuiz" {
			result, err := db.Query("SELECT * from quiz   WHERE username = ?", username)
			if err != nil {
				Errorr(err)
			}
			var quizes []quiz
			for result.Next() {
				var q quiz
				err = result.Scan(&q.Qid, &q.username, &q.quizName, &q.time, &q.numberOfQuestions, &q.result)
				if err != nil {
					Errorr(err)
				}
				quizes = append(quizes, q)
				fmt.Println(q)
			}
			temp := ""
			for i := 0; i < len(quizes); i++ {
				temp += strconv.Itoa(quizes[i].Qid) + "." + quizes[i].quizName + "#"
			}
			fmt.Fprintf(client, temp+"\n")
			message, err = bufio.NewReader(client).ReadString('\n')
			chosenQuiz := strings.TrimSuffix(message, "\n")
			if err != nil {
				Errorr(err)
			}
			var theQuiz quiz
			for i := 0; i < len(quizes); i++ {
				if strconv.Itoa(quizes[i].Qid) == chosenQuiz {
					theQuiz = quizes[i]
					break
				}
			}
			fmt.Fprintf(client, strconv.Itoa(theQuiz.time)+"\n")
			var answers []int
			for i := 0; i <= theQuiz.numberOfQuestions; i++ {
				answers = append(answers, 0)
			}
			for message != "$" {
				message, err = bufio.NewReader(client).ReadString('\n')
				arrr := strings.Split(message, "\n")
				if err != nil {
					Errorr(err)
				}
				if arrr[0] == "$" {
					break
				}
				arr := strings.Split(arrr[0], ".")
				q, _ := strconv.Atoi(arr[0])
				a, _ := strconv.Atoi(arr[1])
				answers[q] = a
			}
			corrects := 0
			wrongs := 0
			for i := 1; i <= theQuiz.numberOfQuestions; i++ {
				result, err := db.Query("SELECT answer from answers   WHERE Qid = ? and id = ?", theQuiz.Qid, i)
				if err != nil {
					Errorr(err)
				}
				var theAns int
				for result.Next() {
					err = result.Scan(&theAns)
					if err != nil {
						Errorr(err)
					}
					if answers[i] == 0 {
						continue
					}
					if answers[i] == theAns {
						corrects++
					} else {
						wrongs++
					}
				}

			}
			fmt.Println(corrects)
			fmt.Println(wrongs)

			precentage := float64(3*corrects-wrongs) / float64(3*theQuiz.numberOfQuestions)
			fmt.Println(precentage)
			precentage *= 100
			fmt.Fprintf(client, theQuiz.quizName+"#"+strconv.Itoa(int(math.Abs(precentage)))+"\n")
			_, err = db.Query("UPDATE quiz SET result=? WHERE username = ? and Qid = ?", precentage, theQuiz.username, theQuiz.Qid)
			if err != nil {
				Errorr(err)
			}
		}
	}
}
