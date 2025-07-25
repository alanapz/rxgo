package utils

import (
	"fmt"
	"reflect"
)

func Drop[T any](val T, err error) error {
	return err
}

func Wrap[T any](value1 T, err error) func(*T) error {
	return func(out1 *T) error {

		if err != nil {
			return err
		}

		*out1 = value1
		return nil
	}
}

func Wrap2[T any, U any](value1 T, value2 U, err error) func(*T, *U) error {
	return func(out1 *T, out2 *U) error {

		if err != nil {
			return err
		}

		*out1 = value1
		*out2 = value2
		return nil
	}

}

func Cast[U any, T any](value1 T, err error) (U, error) {

	var zero U

	if err != nil {
		return zero, err
	}

	result1, ok := reflect.ValueOf(value1).Interface().(U)

	if !ok {
		return zero, fmt.Errorf("cannot cast %v of type %s to type %s", value1, reflect.TypeFor[T](), reflect.TypeFor[U]())
	}

	return result1, nil
}

// func Cast[U any, T any]() func(func(U, error) error) func(T, error) error {

// 	var empty U

// 	return func(next func(U, error) error) func(T, error) error {
// 		return func(t T, err error) error {

// 			if err != nil {
// 				return next(empty, err)
// 			}

// 			return next(reflect.ValueOf(t).Interface().(U), nil)
// 		}
// 	}
// }

func Wrap3[T any, U any, V any](val1 T, val2 U, val3 V, err error) func(func(T, error) error, func(U, error) error, func(V, error) error) error {
	return func(f1 func(T, error) error, f2 func(U, error) error, f3 func(V, error) error) error {

		err1 := f1(val1, err)
		err2 := f2(val2, err)
		err3 := f3(val3, err)

		if err1 != nil {
			return err1
		}

		if err2 != nil {
			return err2
		}

		if err3 != nil {
			return err3
		}

		return nil
	}
}

func Map[T any, U any](mapper func(T) U) func(func(U, error) error) func(T, error) error {

	var empty U

	return func(next func(U, error) error) func(T, error) error {
		return func(t T, err error) error {

			if err != nil {
				return next(empty, err)
			}

			return next(mapper(t), nil)
		}
	}
}
