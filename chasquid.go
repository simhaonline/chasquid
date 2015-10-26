package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/mail"
	"net/textproto"
	"strings"

	"github.com/golang/glog"
	"golang.org/x/net/trace"
)

func main() {
	flag.Parse()

	monAddr := ":1099"
	glog.Infof("Monitoring HTTP server listening on %s", monAddr)
	go http.ListenAndServe(monAddr, nil)

	ListenAndServe()
}

const (
	// TODO: get this via config/dynamically. It's only used for show.
	hostname = "charqui.com.ar"

	// Maximum data size, in bytes.
	maxDataSize = 52428800
)

func getTLSConfig() (*tls.Config, error) {
	var err error
	conf := &tls.Config{}

	// TODO: Get these from the configuration (we have to support many, not
	// just 1 like here).
	conf.Certificates = make([]tls.Certificate, 1)
	conf.Certificates[0], err = tls.LoadX509KeyPair(".cert.pem", ".key.pem")
	if err != nil {
		return nil, fmt.Errorf("Error loading client certificate: %v", err)
	}

	conf.BuildNameToCertificate()

	return conf, nil
}

func ListenAndServe() {
	// Configure TLS.
	tlsConfig, err := getTLSConfig()
	if err != nil {
		glog.Fatalf("Error loading TLS config: %v", err)
	}

	// Listen.
	addr := ":1025"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("Error listening: %v", err)
	}
	defer l.Close()

	glog.Infof("Server listening on %s", addr)

	// Serve.
	for {
		conn, err := l.Accept()
		if err != nil {
			glog.Fatalf("Error accepting: %v", err)
		}

		sc := &Conn{
			netconn:   conn,
			tc:        textproto.NewConn(conn),
			tlsConfig: tlsConfig,
		}
		go sc.Handle()
	}
}

type Conn struct {
	// Connection information.
	netconn net.Conn
	tc      *textproto.Conn

	// TLS configuration.
	tlsConfig *tls.Config

	// Envelope.
	mail_from string
	rcpt_to   []string
	data      []byte
}

func (c *Conn) Handle() {
	defer c.netconn.Close()

	tr := trace.New("SMTP", "connection")
	defer tr.Finish()

	c.tc.PrintfLine("220 %s ESMTP charquid", hostname)

	var cmd, params string
	var err error

loop:
	for {
		cmd, params, err = c.readCommand()
		if err != nil {
			c.tc.PrintfLine("554 error reading command: %v", err)
			break
		}

		tr.LazyPrintf("-> %s %s", cmd, params)

		var code int
		var msg string

		switch cmd {
		case "HELO":
			code, msg = c.HELO(params)
		case "EHLO":
			code, msg = c.EHLO(params)
		case "HELP":
			code, msg = c.HELP(params)
		case "NOOP":
			code, msg = c.NOOP(params)
		case "RSET":
			code, msg = c.RSET(params)
		case "MAIL":
			code, msg = c.MAIL(params)
		case "RCPT":
			code, msg = c.RCPT(params)
		case "DATA":
			// DATA handles the whole sequence.
			code, msg = c.DATA(params, tr)
		case "STARTTLS":
			code, msg = c.STARTTLS(params, tr)
		case "QUIT":
			c.writeResponse(221, "Be seeing you...")
			break loop
		default:
			code = 500
			msg = "unknown command"
		}

		if code > 0 {
			tr.LazyPrintf("<- %d  %s", code, msg)

			err = c.writeResponse(code, msg)
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		tr.LazyPrintf("exiting with error: %v", err)
	}
}

func (c *Conn) HELO(params string) (code int, msg string) {
	types := []string{
		"general store", "used armor dealership", "second-hand bookstore",
		"liquor emporium", "antique weapons outlet", "delicatessen",
		"jewelers", "quality apparel and accessories", "hardware",
		"rare books", "lighting store"}
	t := types[rand.Int()%len(types)]
	msg = fmt.Sprintf("Hello my friend, welcome to chasqui's %s!", t)

	return 250, msg
}

func (c *Conn) EHLO(params string) (code int, msg string) {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, hostname+" - Your hour of destiny has come.\n")
	fmt.Fprintf(buf, "8BITMIME\n")
	fmt.Fprintf(buf, "PIPELINING\n")
	fmt.Fprintf(buf, "SIZE %d\n", maxDataSize)
	fmt.Fprintf(buf, "STARTTLS\n")
	fmt.Fprintf(buf, "HELP\n")
	return 250, buf.String()
}

func (c *Conn) HELP(params string) (code int, msg string) {
	return 214, "hoy por ti, mañana por mi"
}

func (c *Conn) RSET(params string) (code int, msg string) {
	c.resetEnvelope()

	msgs := []string{
		"Who was that Maud person anyway?",
		"Thinking of Maud you forget everything else.",
		"Your mind releases itself from mundane concerns.",
		"As your mind turns inward on itself, you forget everything else.",
	}
	return 250, msgs[rand.Int()%len(msgs)]
}

func (c *Conn) NOOP(params string) (code int, msg string) {
	return 250, "noooooooooooooooooooop"
}

func (c *Conn) MAIL(params string) (code int, msg string) {
	// params should be: "FROM:<name@host>"
	// First, get rid of the "FROM:" part (but check it, it's mandatory).
	sp := strings.SplitN(params, ":", 2)
	if len(sp) != 2 || sp[0] != "FROM" {
		return 500, "unknown command"
	}

	e, err := mail.ParseAddress(sp[1])
	if err != nil || e.Address == "" {
		return 501, "malformed address"
	}

	if !strings.Contains(e.Address, "@") {
		return 501, "sender address must contain a domain"
	}

	c.resetEnvelope()
	c.mail_from = e.Address
	return 250, "You feel like you are being watched"
}

func (c *Conn) RCPT(params string) (code int, msg string) {
	// params should be: "TO:<name@host>"
	// First, get rid of the "TO:" part (but check it, it's mandatory).
	sp := strings.SplitN(params, ":", 2)
	if len(sp) != 2 || sp[0] != "TO" {
		return 500, "unknown command"
	}

	e, err := mail.ParseAddress(sp[1])
	if err != nil || e.Address == "" {
		return 501, "malformed address"
	}

	if c.mail_from == "" {
		return 503, "sender not yet given"
	}

	// RFC says 100 is the minimum limit for this, but it seems excessive.
	if len(c.rcpt_to) > 100 {
		return
	}

	// TODO: do we allow receivers without a domain?
	// TODO: check the case:
	//  - local recipient, always ok
	//  - external recipient, only ok if mail_from is local (needs auth)

	c.rcpt_to = append(c.rcpt_to, e.Address)
	return 250, "You have an eerie feeling..."
}

func (c *Conn) DATA(params string, tr trace.Trace) (code int, msg string) {
	if c.mail_from == "" {
		return 503, "sender not yet given"
	}

	if len(c.rcpt_to) == 0 {
		return 503, "need an address to send to"
	}

	// We're going ahead.
	err := c.writeResponse(354, "You suddenly realize it is unnaturally quiet")
	if err != nil {
		return 554, fmt.Sprintf("error writing DATA response: %v", err)
	}

	tr.LazyPrintf("<- 354  You experience a strange sense of peace")

	dotr := io.LimitReader(c.tc.DotReader(), maxDataSize)
	c.data, err = ioutil.ReadAll(dotr)
	if err != nil {
		return 554, fmt.Sprintf("error reading DATA: %v", err)
	}

	tr.LazyPrintf("-> ... %d bytes of data", len(c.data))

	// TODO: here is where we queue/send/process the message!

	// It is very important that we reset the envelope before returning,
	// so clients can send other emails right away without needing to RSET.
	c.resetEnvelope()

	msgs := []string{
		"You offer the Amulet of Yendor to Anhur...",
		"An invisible choir sings, and you are bathed in radiance...",
		"The voice of Anhur booms out: Congratulations, mortal!",
		"In return to thy service, I grant thee the gift of Immortality!",
		"You ascend to the status of Demigod(dess)...",
	}
	return 250, msgs[rand.Int()%len(msgs)]
}

func (c *Conn) STARTTLS(params string, tr trace.Trace) (code int, msg string) {
	err := c.writeResponse(220, "You experience a strange sense of peace")
	if err != nil {
		return 554, fmt.Sprintf("error writing STARTTLS response: %v", err)
	}

	tr.LazyPrintf("<- 220  You experience a strange sense of peace")

	client := tls.Server(c.netconn, c.tlsConfig)
	err = client.Handshake()
	if err != nil {
		return 554, fmt.Sprintf("error in client handshake: %v", err)
	}

	tr.LazyPrintf("<> ...  jump to TLS was successful")

	// Override the connections. We don't need the older ones anymore.
	c.netconn = client
	c.tc = textproto.NewConn(client)

	// Reset the envelope; clients must start over after switching to TLS.
	c.resetEnvelope()

	// 0 indicates not to send back a reply.
	return 0, ""
}

func (c *Conn) resetEnvelope() {
	c.mail_from = ""
	c.rcpt_to = nil
	c.data = nil
}

func (c *Conn) readCommand() (cmd, params string, err error) {
	var msg string

	msg, err = c.tc.ReadLine()
	if err != nil {
		return "", "", err
	}

	sp := strings.SplitN(msg, " ", 2)
	cmd = strings.ToUpper(sp[0])
	if len(sp) > 1 {
		params = sp[1]
	}

	return cmd, params, err
}

func (c *Conn) writeResponse(code int, msg string) error {
	defer c.tc.W.Flush()

	return writeResponse(c.tc.W, code, msg)
}

// writeResponse writes a multi-line response to the given writer.
// This is the writing version of textproto.Reader.ReadResponse().
func writeResponse(w io.Writer, code int, msg string) error {
	var i int
	lines := strings.Split(msg, "\n")

	// The first N-1 lines use "<code>-<text>".
	for i = 0; i < len(lines)-2; i++ {
		_, err := w.Write([]byte(fmt.Sprintf("%d-%s\r\n", code, lines[i])))
		if err != nil {
			return err
		}
	}

	// The last line uses "<code> <text>".
	_, err := w.Write([]byte(fmt.Sprintf("%d %s\r\n", code, lines[i])))
	if err != nil {
		return err
	}

	return nil
}
