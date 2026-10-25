// vapidkeys generates a VAPID key pair for Web Push. Run it once per environment and store
// the output as the VAPID_PUBLIC_KEY and VAPID_PRIVATE_KEY environment variables / secrets.
//
//	go run ./cmd/vapidkeys
package main

import (
	"fmt"
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

func main() {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to generate VAPID keys:", err)
		os.Exit(1)
	}

	fmt.Println("VAPID_PUBLIC_KEY=" + publicKey)
	fmt.Println("VAPID_PRIVATE_KEY=" + privateKey)
}
