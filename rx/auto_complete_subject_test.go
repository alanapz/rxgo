package rx

// import (
// 	"errors"
// 	"fmt"
// 	"sync"
// 	"testing"

// 	u "alanpinder.com/rxgo/v2/utils"
// )

// // TestHelloName calls greetings.Hello with a name, checking
// // for a valid return value.
// func TestAutoCompleteSubjectAlreadyCompleted(t *testing.T) {

// 	cleanupTest := tu.PrepareTest(t)
// 	defer cleanupTest()

// 	subject := NewAutoCompleteSubject[string]()

// 	fmt.Printf("Sending value: 'hello'\n")
// 	subject.Value("hello")

// 	var wg sync.WaitGroup

// 	done := addTestSubscriber(t, &wg, "s1", subject, u.Of("hello"), u.Of[error]())

// 	wg.Wait()
// 	done()
// }

// func TestAutoCompleteSubjectValue(t *testing.T) {

// 	cleanupTest := PrepareTest(t)
// 	defer cleanupTest()

// 	subject := NewAutoCompleteSubject[string]()
// 	subject.Value("hello")

// 	var wg sync.WaitGroup

// 	done1 := tu.AddSubscriber(t, &wg, "s1", subject, u.Of("hello"), u.Of[error]())
// 	done2 := tu.AddSubscriber(t, &wg, "s2", subject, u.Of("hello"), u.Of[error]())
// 	done3 := tu.AddSubscriber(t, &wg, "s3", subject, u.Of("hello"), u.Of[error]())

// 	wg.Wait()
// 	done1()
// 	done2()
// 	done3()
// }

// func TestAutoCompleteSubjectError(t *testing.T) {

// 	cleanupTest := tu.PrepareTest(t)
// 	defer cleanupTest()

// 	err := errors.New("new error")

// 	subject := NewAutoCompleteSubject[string]()
// 	subject.Error(err)

// 	var wg sync.WaitGroup

// 	done1 := addTestSubscriber(t, &wg, "s1", subject, u.Of[string](), u.Of(err))
// 	done2 := addTestSubscriber(t, &wg, "s2", subject, u.Of[string](), u.Of(err))
// 	done3 := addTestSubscriber(t, &wg, "s3", subject, u.Of[string](), u.Of(err))

// 	wg.Wait()
// 	done1()
// 	done2()
// 	done3()
// }
