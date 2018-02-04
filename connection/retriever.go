package connection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"strings"
)

type NoTLSConnErr string

func (f NoTLSConnErr) Error() string {
	return fmt.Sprintf("No TLS Conn Received")
}

func separateDomainPort(domain string) (string, string) {
	var domainPort = strings.SplitN(domain, ":", 2)
	if len(domainPort) == 2 {
		return domainPort[0], domainPort[1]
	}
	return domain, "443"
}

func Connect(domain, cipherscanbinPath string) ([]byte, error) {

	var domainName, port = separateDomainPort(domain)

	ip := getRandomIP(domainName)

	if ip == "" {
		e := fmt.Errorf("Could not resolve ip for: ", domain)
		log.Println(e)
		return nil, e
	}

	cmd := cipherscanbinPath + " --no-tolerance -j --curves -servername " + domain + " " + ip + ":" + port
	log.Println(cmd)
	comm := exec.Command("bash", "-c", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	comm.Stdout = &out
	comm.Stderr = &stderr
	err := comm.Run()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	info := CipherscanOutput{}
	err = json.Unmarshal([]byte(out.String()), &info)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	info.Target = domain
	info.IP = ip

	c, err := info.Stored()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return json.Marshal(c)
}

func getRandomIP(domain string) string {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return ""
	}

	max := len(ips)

	for {
		if max == 0 {
			return ""
		}
		index := rand.Intn(len(ips))

		if ips[index].To4() != nil {
			return ips[index].String()
		} else {
			ips = append(ips[:index], ips[index+1:]...)
		}
		max--
	}
}
