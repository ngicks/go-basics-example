package oneof

import (
	"reflect"
	"slices"
)

func SendEach[T ~[]C, C ~(chan<- E), E any](chans T, fn func() E, cancel <-chan struct{}) (sent []int, completed bool) {
	chans = slices.Clone(chans)
	sent = make([]int, 0, len(chans))
	completed = true

	for len(chans) != len(sent) {
		chosen, ok := Send(chans, fn(), cancel)
		if !ok {
			completed = false
			break
		}
		sent = append(sent, chosen)
		chans[chosen] = nil
	}
	return
}

func Send[T ~[]C, C ~(chan<- E), E any](chans T, v E, cancel <-chan struct{}) (chosen int, sent bool) {
	switch x := len(chans); {
	case x == 0:
		panic("zero chans")
	case x <= 4:
		var c [4]C
		_ = copy(c[:], chans)
		return Send4(c, v, cancel)
	case x <= 8:
		var c [8]C
		_ = copy(c[:], chans)
		return Send8(c, v, cancel)
	case x <= 16:
		var c [16]C
		_ = copy(c[:], chans)
		return Send16(c, v, cancel)
	default:
		return SendN(chans, v, cancel)
	}
}

func Send4[T ~[4]C, C ~(chan<- E), E any](chans T, v E, cancel <-chan struct{}) (chosen int, sent bool) {
	sent = true
	select {
	case <-cancel:
		sent = false
	case chans[0] <- v:
		chosen = 0
	case chans[1] <- v:
		chosen = 1
	case chans[2] <- v:
		chosen = 2
	case chans[3] <- v:
		chosen = 3
	}
	return
}

func Send8[T ~[8]C, C ~(chan<- E), E any](chans T, v E, cancel <-chan struct{}) (chosen int, sent bool) {
	sent = true
	select {
	case <-cancel:
		sent = false
	case chans[0] <- v:
		chosen = 0
	case chans[1] <- v:
		chosen = 1
	case chans[2] <- v:
		chosen = 2
	case chans[3] <- v:
		chosen = 3
	case chans[4] <- v:
		chosen = 4
	case chans[5] <- v:
		chosen = 5
	case chans[6] <- v:
		chosen = 6
	case chans[7] <- v:
		chosen = 7
	}
	return
}

func Send16[T ~[16]C, C ~(chan<- E), E any](chans T, v E, cancel <-chan struct{}) (chosen int, sent bool) {
	sent = true
	select {
	case <-cancel:
		sent = false
	case chans[0] <- v:
		chosen = 0
	case chans[1] <- v:
		chosen = 1
	case chans[2] <- v:
		chosen = 2
	case chans[3] <- v:
		chosen = 3
	case chans[4] <- v:
		chosen = 4
	case chans[5] <- v:
		chosen = 5
	case chans[6] <- v:
		chosen = 6
	case chans[7] <- v:
		chosen = 7
	case chans[8] <- v:
		chosen = 8
	case chans[9] <- v:
		chosen = 9
	case chans[10] <- v:
		chosen = 10
	case chans[11] <- v:
		chosen = 11
	case chans[12] <- v:
		chosen = 12
	case chans[13] <- v:
		chosen = 13
	case chans[14] <- v:
		chosen = 14
	case chans[15] <- v:
		chosen = 15
	}
	return
}

func SendN[T ~[]C, C ~(chan<- E), E any](chans T, v E, cancel <-chan struct{}) (chosen int, sent bool) {
	cases := []reflect.SelectCase{{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(cancel),
	}}
	for _, ch := range chans {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: reflect.ValueOf(v),
		})
	}
	chosen, _, _ = reflect.Select(cases)
	if chosen == 0 {
		return
	}
	return chosen - 1, true
}

func Recv[T ~[]C, C ~(<-chan E), E any](chans T, cancel <-chan struct{}) (v E, chosen int, received bool) {
	switch x := len(chans); {
	case x == 0:
		panic("zero chans")
	case x <= 4:
		var c [4]C
		_ = copy(c[:], chans)
		return Recv4(c, cancel)
	case x <= 8:
		var c [8]C
		_ = copy(c[:], chans)
		return Recv8(c, cancel)
	case x <= 16:
		var c [16]C
		_ = copy(c[:], chans)
		return Recv16(c, cancel)
	default:
		return RecvN(chans, cancel)
	}
}

func Recv4[T ~[4]C, C ~(<-chan E), E any](chans T, cancel <-chan struct{}) (v E, chosen int, received bool) {
	received = true
	select {
	case <-cancel:
		received = false
	case v = <-chans[0]:
		chosen = 0
	case v = <-chans[1]:
		chosen = 1
	case v = <-chans[2]:
		chosen = 2
	case v = <-chans[3]:
		chosen = 3
	}
	return
}

func Recv8[T ~[8]C, C ~(<-chan E), E any](chans T, cancel <-chan struct{}) (v E, chosen int, received bool) {
	received = true
	select {
	case <-cancel:
		received = false
	case v = <-chans[0]:
		chosen = 0
	case v = <-chans[1]:
		chosen = 1
	case v = <-chans[2]:
		chosen = 2
	case v = <-chans[3]:
		chosen = 3
	case v = <-chans[4]:
		chosen = 4
	case v = <-chans[5]:
		chosen = 5
	case v = <-chans[6]:
		chosen = 6
	case v = <-chans[7]:
		chosen = 7
	}
	return
}

func Recv16[T ~[16]C, C ~(<-chan E), E any](chans T, cancel <-chan struct{}) (v E, chosen int, received bool) {
	received = true
	select {
	case <-cancel:
		// chosen should stay zero, to prevent misuse.
		received = false
	case v = <-chans[0]:
		chosen = 0
	case v = <-chans[1]:
		chosen = 1
	case v = <-chans[2]:
		chosen = 2
	case v = <-chans[3]:
		chosen = 3
	case v = <-chans[4]:
		chosen = 4
	case v = <-chans[5]:
		chosen = 5
	case v = <-chans[6]:
		chosen = 6
	case v = <-chans[7]:
		chosen = 7
	case v = <-chans[8]:
		chosen = 8
	case v = <-chans[9]:
		chosen = 9
	case v = <-chans[10]:
		chosen = 10
	case v = <-chans[11]:
		chosen = 11
	case v = <-chans[12]:
		chosen = 12
	case v = <-chans[13]:
		chosen = 13
	case v = <-chans[14]:
		chosen = 14
	case v = <-chans[15]:
		chosen = 15
	}
	return
}

func RecvN[T ~[]C, C ~(<-chan E), E any](chans T, cancel <-chan struct{}) (value E, chosen int, received bool) {
	cases := []reflect.SelectCase{{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(cancel),
	}}
	for _, ch := range chans {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}
	chosen, recv, _ := reflect.Select(cases)
	if chosen == 0 {
		return
	}
	return recv.Interface().(E), chosen - 1, true
}
