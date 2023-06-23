package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(" Syntax: main.go <domain_name>")
		return
	}

	domainName := os.Args[1]
	whoisServer := "whois.iana.org"

	// Find the Whois server for the top-level domain
	topLevelDomain, _ := getTopLevelDomain(domainName)
	resp, _ := sendWhoisQuery(whoisServer, topLevelDomain)
	whoisServer = extractWhoisServer(resp)

	// Send the Whois query to the appropriate server
	resp, err := sendWhoisQuery(whoisServer, domainName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the response
	fmt.Println(string(resp))
}

func getTopLevelDomain(domainName string) (string, error) {
	parts := strings.Split(domainName, ".")
	if len(parts) < 2 {
		return "", errors.New("Invalid domain name")
	}
	return parts[len(parts)-1], nil
}

func extractWhoisServer(whoisResp []byte) string {
	respStr := string(whoisResp)
	re := regexp.MustCompile(`whois:\s+([\w.-]+)`)
	match := re.FindStringSubmatch(respStr)
	if match == nil {
		return ""
	}
	return match[1]
}

func sendWhoisQuery(whoisServer string, domainName string) ([]byte, error) {
	conn, err := net.Dial("tcp", whoisServer+":43")
	if err != nil {
		return []byte{}, err
	}
	defer conn.Close()

	fmt.Fprintf(conn, domainName+"\r\n")
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		return []byte{}, err
	}

	return resp, nil
}
