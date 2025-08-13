package wireguard

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync/atomic"
)

var ipCounter uint32 = 10 // start from .10

func GenerateKeyPair() (privateKey, publicKey string, err error) {
	// ⚠️ ceci n'est PAS une vraie clé WireGuard, mais OK pour un PoC côté FE.
	priv := make([]byte, 32)
	if _, err = rand.Read(priv); err != nil { return "", "", err }
	privateKey = base64.StdEncoding.EncodeToString(priv)
	publicKey = base64.StdEncoding.EncodeToString(reverseBytes(priv))
	return privateKey, publicKey, nil
}

func reverseBytes(b []byte) []byte {
	n := len(b)
	out := make([]byte, n)
	for i := 0; i < n; i++ { out[i] = b[n-1-i] }
	return out
}

func ReserveVPNIP() string {
	n := atomic.AddUint32(&ipCounter, 1)
	return fmt.Sprintf("10.10.10.%d/32", n%250+5)
}

func GenerateMikrotikScript(privateKey, ip string) string {
	return fmt.Sprintf(`# MikroTik WireGuard onboarding script
/interface wireguard add name=wg0 private-key="%s"
/ip address add address=%s interface=wg0
/interface wireguard peers add public-key="SERVER_PUBLIC_KEY" endpoint-address="vpn.example.com" endpoint-port=51820 allowed-address=0.0.0.0/0
`, privateKey, ip)
}
