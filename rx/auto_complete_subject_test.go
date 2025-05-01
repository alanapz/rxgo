package rx

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"alanpinder.com/rxgo/v2/ux"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestAutoCompleteSubjectAlreadyCompleted(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	subject := NewAutoCompleteSubject[string]()

	fmt.Printf("Sending value: 'hello'\n")
	subject.Value("hello")

	wg := &sync.WaitGroup{}

	done := addTestSubscriber(t, wg, "s1", subject, ux.Of("hello"), ux.Of[error]())

	wg.Wait()
	done()
}

func TestAutoCompleteSubjectValue(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	subject := NewAutoCompleteSubject[string]()
	subject.Value("hello")

	wg := &sync.WaitGroup{}

	done1 := addTestSubscriber(t, wg, "s1", subject, ux.Of("hello"), ux.Of[error]())
	done2 := addTestSubscriber(t, wg, "s2", subject, ux.Of("hello"), ux.Of[error]())
	done3 := addTestSubscriber(t, wg, "s3", subject, ux.Of("hello"), ux.Of[error]())

	wg.Wait()
	done1()
	done2()
	done3()
}

func TestAutoCompleteSubjectError(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	err := errors.New("new error")

	subject := NewAutoCompleteSubject[string]()
	subject.Error(err)

	wg := &sync.WaitGroup{}

	done1 := addTestSubscriber(t, wg, "s1", subject, ux.Of[string](), ux.Of(err))
	done2 := addTestSubscriber(t, wg, "s2", subject, ux.Of[string](), ux.Of(err))
	done3 := addTestSubscriber(t, wg, "s3", subject, ux.Of[string](), ux.Of(err))

	wg.Wait()
	done1()
	done2()
	done3()
}
