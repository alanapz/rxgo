package rx

import (
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

func sendValuesThenEndOfStreamAsync[T any](env *RxEnvironment, downstream chan T, downstreamUnsubscribed <-chan u.Never, values ...T) {

	env.Execute(func() {
		for value := range slices.Values(values) {
			if u.Selection(u.SelectDone(downstreamUnsubscribed), u.SelectSend(downstream, value)) == u.DoneResult {
				return
			}
		}
	})

}
