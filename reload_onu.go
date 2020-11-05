package main

import (
	"time"
	"fmt"
	"bufio"
	"regexp"
	"flag"
	"log"
	"os"
	"strings"
	"./telnet"
	h "./wrapers"
)

const timeout = 10 * time.Second

var (
	l           h.MyLogger
	required   = []string{"host", "onu-path"}

	logVerbose = flag.Int("log-level", 3, "Log verbose [ 1 - Fatal, 2 - ERROR, 4 - INFO, 8 - MYSQL, 16 - FUNC, 32 - DEBUG ]")
	logPath    = flag.String("log-file", fmt.Sprintf("%s.log", os.Args[0]), "Log path")
	host	   = flag.String("host", "", "BDCom IP address")
	user	   = flag.String("user", "admin", "Username")
	pass	   = flag.String("password", "", "Password")
	enablePass = flag.String("enable", "", "Enable password")
	onuPath	   = flag.String("onu-path", "", "File with ONUs MAC-addresses to reload")
	onuMac	   = flag.String("onu-mac", "", "Single ONUs MAC-addresses to reload")
)

func main() {

	flag.Parse()

	//
	// Init logs
	//

	f, err := os.OpenFile(*logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	l = h.InitLog(f, *logVerbose)
	l.Printf(h.DEBUG, "START")


	//
	// Validate incoming args
	//

	validIP := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)
	validMac := regexp.MustCompile(`^[a-f0-9]{4}\.[a-f0-9]{4}\.[a-f0-9]{4}$`)

	if (!validIP.MatchString(*host)) {
		l.Printf(h.ERROR, "Wrong IP Address")
		os.Exit(1)
	}

	var onuMAC []string
	if (validMac.MatchString(*onuMac)) {
        onuMAC = append(onuMAC, *onuMac)
	} else {
		file, err := os.Open(*onuPath)
		checkErr(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if m := strings.TrimSpace(scanner.Text()); validMac.MatchString(m) {
				onuMAC = append(onuMAC, m)
			}
		}
	}

	if (len(onuMAC) == 0) {
		l.Printf(h.ERROR, "Empty ONUs MAC-addresses list")
		os.Exit(1)
	} else {
		l.Printf(h.INFO, fmt.Sprintf("Found %d ONUs MAC-addresses to reload", len(onuMAC)))
	}

	//
	// Connect to BDCom
	//

	l.Printf(h.INFO, fmt.Sprintf("Try connect to %s:23", *host))
	t, err := telnet.Dial("tcp", fmt.Sprintf("%s:23", *host))
	checkErr(err)
	t.SetUnixWriteMode(true)
	defer t.Close()

	//
	// Authorize
	//

	expect(t, "sername: ")
	sendln(t, *user)
	expect(t, "assword: ")
	sendln(t, *pass)
	expect(t, ">")
	sendln(t, "enable")
	expect(t, "assword:", "#")
	sendln(t, *enablePass)
	expect(t, "#")

	//
	// Payload
	//

	for _, mac := range onuMAC {
		l.Printf(h.INFO, fmt.Sprintf("Try reload ONU with MAC-Address: %s", mac))
		sendln(t, fmt.Sprintf("epon reboot onu mac-address %s", mac))
		expect(t, "?", "#")
		sendln(t, "y")
		expect(t, "#")
	}

	l.Printf(h.DEBUG, "END")

	// var data []byte
	// sendln(t, "show version")
	// data, err = t.ReadBytes('#')
	// checkErr(err)
	// os.Stdout.Write(data)
	// os.Stdout.WriteString("\n")
}

func checkErr(err error) {
	if err != nil {
		l.Printf(h.ERROR, fmt.Sprintf("%s", err))
		os.Exit(1)
	}
}

func expect(t *telnet.Conn, d ...string) {
	checkErr(t.SetReadDeadline(time.Now().Add(timeout)))
	checkErr(t.SkipUntil(d...))
}

func sendln(t *telnet.Conn, s string) {
	checkErr(t.SetWriteDeadline(time.Now().Add(timeout)))
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err := t.Write(buf)
	checkErr(err)
}
