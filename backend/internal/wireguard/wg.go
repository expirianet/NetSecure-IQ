package wireguard

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "github.com/google/uuid"
)

type Peer struct {
    ID         uuid.UUID
    PrivateKey string
    PublicKey  string
    IPAddress  string
    Enabled    bool
}

func GenerateKeyPair() (privateKey, publicKey string, err error) {
    priv := make([]byte, 32)
    _, err = rand.Read(priv)
    if err != nil {
        return "", "", err
    }
    privateKey = base64.StdEncoding.EncodeToString(priv)
    publicKey = base64.StdEncoding.EncodeToString([]byte(reverse(string(priv))))
    return privateKey, publicKey, nil
}

func reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func CreatePeerAndScript(mac string) (*Peer, string, error) {
    priv, pub, err := GenerateKeyPair()
    if err != nil {
        return nil, "", err
    }
    peer := &Peer{
        ID:         uuid.New(),
        PrivateKey: priv,
        PublicKey:  pub,
        IPAddress:  fmt.Sprintf("10.10.10.%d/32", randByte()),
        Enabled:    false,
    }
    script := GenerateMikrotikScript(peer, mac)
    return peer, script, nil
}

func randByte() int {
    b := make([]byte, 1)
    rand.Read(b)
    return int(b[0])
}

func GenerateMikrotikScript(peer *Peer, mac string) string {
    return fmt.Sprintf(`
# MikroTik WireGuard onboarding script
/interface wireguard add name=wg0 private-key="%s"
/ip address add address=%s interface=wg0
/interface wireguard peers add public-key="SERVER_PUBLIC_KEY" endpoint-address="vpn.example.com" endpoint-port=51820 allowed-address=0.0.0.0/0
# MAC: %s
`, peer.PrivateKey, peer.IPAddress, mac)
}
