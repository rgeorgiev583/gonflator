package postcommit

import (
	"io"
	"log"
	"net"
	"url"

	"bitbucket.org/kardianos/rsync"
	"bitbucket.org/kardianos/rsync/proto"
	"github.com/go-ini/ini"
	
	"github.com/rgeorgiev583/gonflator/delta"
	"github.com/rgeorgiev583/gonflator/git/hooks/postcommit"
)

func Push(rdiff <-chan delta.Delta, target io.Writer) {
	go func() {
		delta <- rdiff
		target.Write(delta, 
	}	
}

func main() int {
	var url string

	if (len(os.Args) < 2) {
		cfg, err := ini.Load(".gonflation")
		if err != nil {
			log.Fatalln(err.Error())
			return 1
		}
		
		section, err := cfg.GetSection("")
		if err != nil {
			log.Fatalln(err.Error())
			return 2
		}
		
		key, err := section.GetKey("url")
		if err != nil {
			log.Fatalln(err.Error())
			return 3
		}
		
		url = key.String()
	} else {
		url = os.Args[1]
	}
	
	conn, err := net.Dial("tcp", url)
	if err != nil {
		log.Fatalln(err.Error())
		return 4
	}
	
	opsEncode := &proto.Writer{Writer: conn}
	opsEncode.Header(proto.TypeDelta, proto.CompGZip, rsync.DefaultBlockSize)
	conn.Write(opsEncode.)
	
	return 0
}
	
}