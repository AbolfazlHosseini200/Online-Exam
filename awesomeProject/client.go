package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	time2 "time"
)

func main() {
	var username string
	fmt.Print("username:")
	fmt.Scan(&username)
	n := 1
	if n == 1 {
		connection, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			Errorrr(err)
		} else {
			fmt.Println("connected to the server!!!!!!!")
		}
		fmt.Fprintf(connection, username+"\n")
		if username == "server" {
			fmt.Println("Choose Your Action:\n1.create a quiz")
			fmt.Scan(&n)
			if n == 1 {
				fmt.Fprintf(connection, "create\n")
				var Qid int
				var username string
				var quizName string
				var time int
				var numberOfQuestions int
				fmt.Println("enter qid:")
				fmt.Scan(&Qid)
				fmt.Fprintf(connection, strconv.Itoa(Qid)+"\n")
				fmt.Println("enter username:")
				fmt.Scan(&username)
				fmt.Fprintf(connection, username+"\n")
				fmt.Println("enter quiz name:")
				fmt.Scan(&quizName)
				fmt.Fprintf(connection, quizName+"\n")
				fmt.Println("enter time:")
				fmt.Scan(&time)
				fmt.Fprintf(connection, strconv.Itoa(time)+"\n")
				fmt.Println("enter number of questions:")
				fmt.Scan(&numberOfQuestions)
				fmt.Fprintf(connection, strconv.Itoa(numberOfQuestions)+"\n")
				for i := 1; i <= numberOfQuestions; i++ {
					var a string
					fmt.Print(strconv.Itoa(i) + ":")
					fmt.Scan(&a)
					fmt.Fprintf(connection, a+"\n")
				}
			}
		} else {
			for {
				fmt.Println("Choose Your Action:\n1.Attend A Quiz")
				fmt.Scan(&n)
				if n == 1 {
					fmt.Fprintf(connection, "showQuiz#"+"\n")
					message, err := bufio.NewReader(connection).ReadString('\n')
					if err != nil {
						Errorrr(err)
					}
					message = strings.TrimSuffix(message, "\n")
					quizes := strings.Split(message, "#")
					for i := 0; i < len(quizes); i++ {
						fmt.Println(quizes[i])
					}
					fmt.Scan(&n)
					fmt.Fprintf(connection, strconv.Itoa(n)+"\n")
					message, err = bufio.NewReader(connection).ReadString('\n')
					if err != nil {
						Errorrr(err)
					}
					message = strings.TrimSuffix(message, "\n")
					time, err := strconv.Atoi(message)
					if err != nil {
						Errorrr(err)
					}
					fmt.Println("You have " + message + " minutes to answer questions\nwhenever you're ready press enter!")
					fmt.Scan()
					fmt.Scan()
					fmt.Scan()
					start := time2.Now()
					fmt.Println("enter your answers with this format:\nquestionNumber.answer\nif you want to delete your answer:\nquestionNumber.0\nat the end enter $")
					for start.Before(start.Add(time2.Duration(time) * time2.Minute)) {
						var ans string
						fmt.Scan(&ans)
						if len(strings.Split(ans, ".")) == 2 {
							fmt.Fprintf(connection, ans+"\n")
						} else {
							fmt.Println("wrong command!!!!!!!")
						}
						if ans == "$" {
							break
						}
					}
					if start.Before(start.Add(time2.Duration(time) * time2.Minute)) {
						fmt.Fprintf(connection, "$"+"\n")
					}
					message, err = bufio.NewReader(connection).ReadString('\n')
					if err != nil {
						Errorrr(err)
					}
					message = strings.TrimSuffix(message, "\n")
					arr := strings.Split(message, "#")
					f, err := os.Create(arr[0] + ".txt")
					defer f.Close()
					_, err2 := f.WriteString("percentage:" + arr[1])
					if err2 != nil {
						Errorrr(err2)
					}
				}
			}
		}
	}
}
func Errorrr(text error) {
	fmt.Println(text)
	os.Exit(1)
}
