package agent

import "log"

func Agent() {
	power := 3

	for i := 0; i < power; i++ {
		log.Printf("worker %d starting", i+1)
		go worker()
	}

	select {} // бесконечное ожидание (можно сделать через wg.Wait, но смысла нет)
}
