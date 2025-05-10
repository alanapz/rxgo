package rx

import "slices"

type OperatorFunction[T any, U any] = func(Observable[T]) Observable[U]

func Pipe[T any](source Observable[T], operators ...OperatorFunction[T, T]) Observable[T] {

	current := source

	for operator := range slices.Values(operators) {
		current = operator(current)
	}

	return current
}

func Pipe2[T1 any, T2 any](source Observable[T1], o1 OperatorFunction[T1, T2]) Observable[T2] {
	return o1(source)
}

func Pipe3[T1 any, T2 any, T3 any](source Observable[T1], o1 OperatorFunction[T1, T2], o2 OperatorFunction[T2, T3]) Observable[T3] {
	return o2(Pipe2(source, o1))
}

func Pipe4[T1 any, T2 any, T3 any, T4 any](source Observable[T1], o1 OperatorFunction[T1, T2], o2 OperatorFunction[T2, T3], o3 OperatorFunction[T3, T4]) Observable[T4] {
	return o3(Pipe3(source, o1, o2))
}

func Pipe5[T1 any, T2 any, T3 any, T4 any, T5 any](source Observable[T1], o1 OperatorFunction[T1, T2], o2 OperatorFunction[T2, T3], o3 OperatorFunction[T3, T4], o4 OperatorFunction[T4, T5]) Observable[T5] {
	return o4(Pipe4(source, o1, o2, o3))
}

func Pipe6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](source Observable[T1], o1 OperatorFunction[T1, T2], o2 OperatorFunction[T2, T3], o3 OperatorFunction[T3, T4], o4 OperatorFunction[T4, T5], o5 OperatorFunction[T5, T6]) Observable[T6] {
	return o5(Pipe5(source, o1, o2, o3, o4))
}
