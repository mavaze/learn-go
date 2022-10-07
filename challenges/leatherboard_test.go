package challenges

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

const MAX_RANK = 10

type Quiz struct {
	Id        int
	Questions int
	Interval  int
	Users     []*User
	index     int
}

type Answer struct {
	UserId    int
	IfCorrect bool
	TimeTaken int
}

type User struct {
	UserId       int
	TotalCorrect int
	TimeTaken    int
}

func NewQuiz(id, questions, interval, users int) *Quiz {
	quiz := &Quiz{Id: id, Questions: questions, Interval: interval}
	quiz.Users = make([]*User, users)
	return quiz
}

func (q *Quiz) NextQuestion(answers map[int]Answer) bool {
	for i, u := range q.Users {
		if u == nil {
			u = &User{UserId: i}
			q.Users[i] = u
		}
		if answers[u.UserId].IfCorrect {
			u.TotalCorrect++
		}
		u.TimeTaken += answers[u.UserId].TimeTaken
	}
	sort.Slice(q.Users, func(i, j int) bool {
		if q.Users[i].TotalCorrect == q.Users[j].TotalCorrect {
			return q.Users[i].TimeTaken < q.Users[j].TimeTaken
		}
		return q.Users[i].TotalCorrect > q.Users[j].TotalCorrect
	})

	q.index++
	return q.index < q.Questions
}

func (q *Quiz) GetRanking() []int {
	ranking := make([]int, len(q.Users))
	for i := 0; i < len(q.Users) && i < MAX_RANK; i++ {
		ranking[i] = q.Users[i].UserId
	}
	return ranking
}

func TestLeatherboard(t *testing.T) {
	quiz := NewQuiz(1, 15, 30, 10)
	mockAnswers := make(map[int]Answer, len(quiz.Users))
	for q := 1; ; q++ {
		for u := 0; u < len(quiz.Users); u++ {
			mockAnswers[u] = Answer{u, rand.Intn(10)%2 == 0, rand.Intn(quiz.Interval)}
		}
		fmt.Printf("Question: %d, Answers: %v\n", q, mockAnswers)

		next := quiz.NextQuestion(mockAnswers)
		fmt.Printf("Ranking after %d question: %v\n\n", q, quiz.GetRanking())
		if !next {
			break
		}
	}
	fmt.Println("Final Ranking:", quiz.GetRanking())
}
