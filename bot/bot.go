package bot

import "os"
import "log"
import "bufio"
import "strings"
import "net"
import "fmt"

type IrcWriter func(fmt string, a ...interface{})

type Message struct {
	Contents string
	User     string
	Channel  string
}

type Module interface {
	Name() string
	Handle(Message, IrcWriter)
}

type Bot struct {
	network  string
	channel  string
	nick     string
	password string

	keychar string

	modules []Module
	enabled map[string]bool

	conn net.Conn
}

func (b *Bot) Run() {
	//connect to IRC
	var err error
	b.conn, err = net.Dial("tcp", b.network+":6667")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(b.conn, "USER %s 8 * :%s\n", b.nick, b.nick)
	fmt.Fprintf(b.conn, "NICK %s\n", b.nick)
	fmt.Fprintf(b.conn, "PRIVMSG ChanServ IDENTIFY %s\n", b.password)
	fmt.Fprintf(b.conn, "JOIN "+b.channel+"\n")

	ircwriter := func(format string, a ...interface{}) {
		fmt.Fprintf(b.conn, format, a...)
	}

	ircReader := bufio.NewReader(b.conn)
	for {
		//read message
		line, err := ircReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		//respond to a PING
		if strings.HasPrefix(line, "PING") {
			fmt.Fprintf(b.conn, "PONG %s\n", line[4:])
			continue
		}
		//remove \r\n
		line = line[:len(line)-2]

		//parse relevant bits out
		components := strings.SplitN(line, " ", 4)
		if components[1] != "PRIVMSG" {
			continue //ignore
		}
		m := Message{
			components[0][1:strings.Index(components[0], "!")],
			components[2],
			components[3][1:],
		}

		//pass it to modules
		for _, module := range b.modules {
			module.Handle(m, ircwriter)
		}
	}
}

func (b *Bot) Register(m Module) {
	if b.enabled[m.Name()] {
		b.modules = append(b.modules, m)
	}
}

func NewBot(s string) (b *Bot) {
	b = new(Bot)

	//open settings file
	f, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}

	//get settings out of f
	rdr := bufio.NewReader(f)

	for str, err := rdr.ReadString('\n'); err == nil; str, err = rdr.ReadString('\n') {
		if str[0] == '#' {
			continue
		}
		str = str[:len(str)-1]
		line := strings.SplitN(str, ":", 2)
		arg := line[0]
		val := line[1]
		switch arg {
		case "enable":
			b.enabled[val] = true
		case "network":
			b.network = val
		case "channel":
			b.channel = val
		case "nick":
			b.nick = val
		case "password":
			b.password = val
		case "keychar":
			b.keychar = val
		}
	}

	//close settings file
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	return
}
