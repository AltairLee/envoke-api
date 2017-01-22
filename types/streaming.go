package types

import (
	"github.com/julienschmidt/httprouter"
	"github.com/zballs/go_resonate/crypto"
	. "github.com/zballs/go_resonate/util"
	"net/http"
	"path/filepath"
	"strings"
	"sync/atomic"
)

// HTTP

func FlushWriter(w Writer) flushWriter {
	fw := flushWriter{w: w}
	if f, ok := w.(http.Flusher); ok {
		fw.f = f
	}
	return fw
}

type flushWriter struct {
	f http.Flusher
	w Writer
}

func (fw flushWriter) Write(p []byte) error {
	n, err := fw.w.Write(p)
	if err != nil {
		return err
	} else if size := len(p); n != size {
		return Errorf("Only wrote %d of %d bytes\n", n, size)
	}
	if fw.f != nil {
		fw.f.Flush()
	}
	return nil
}

type HttpService struct {
	dir    string
	files  http.FileSystem
	logger Logger
	router *httprouter.Router
}

func NewHttpService(dir string) *HttpService {
	files := http.Dir(dir)
	logger := NewLogger("streaming_server")
	router := httprouter.New()
	hs := &HttpService{
		dir:    dir,
		files:  files,
		logger: logger,
		router: router,
	}
	hs.SetPlayHandler()
	return hs
}

func (serv *HttpService) Path(args ...string) string {
	args = append([]string{serv.dir}, args...)
	return filepath.Join(args...)
}

func (serv *HttpService) SetPlayHandler() {
	serv.router.GET("/play/:project_title/:song_title/*",
		func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
			// Get values
			projectTitle := params.ByName("project_title")
			songTitle := params.ByName("song_title")
			pub58 := params.ByName("public_key")
			if pub58 == "" {
				errMsg := "Could not find public key"
				serv.logger.Error(errMsg)
				http.Error(w, errMsg, http.StatusBadRequest)
				return
			}
			sig58 := params.ByName("signature")
			if sig58 == "" {
				errMsg := "Could not find signature"
				serv.logger.Error(errMsg)
				http.Error(w, errMsg, http.StatusBadRequest)
				return
			}
			// Public key
			pub := new(crypto.PublicKey)
			if err := pub.FromB58(pub58); err != nil {
				serv.logger.Error(err.Error())
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			// Signature
			sig := new(crypto.Signature)
			if err := sig.FromB58(sig58); err != nil {
				serv.logger.Error(err.Error())
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			// Verify signature
			if !pub.Verify([]byte(projectTitle+songTitle), sig) {
				errMsg := "Signature verification failed"
				serv.logger.Error(errMsg)
				http.Error(w, errMsg, http.StatusUnauthorized)
				return
			}
			// TODO: verify payment
			// Send file bytes
			path := filepath.Join(projectTitle, songTitle, ".mp3")
			file, err := serv.files.Open(path)
			if err != nil {
				serv.logger.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fw := FlushWriter(w)
			bytes := ReadAll(file)
			fw.Write(bytes)
		})
}

// Socket

const LOCALHOST = "localhost"
const PORT = ":8889"

type PlayRequest struct {
	// Payment
	ProjectTitle string            `json:"project_title"`
	SongTitle    string            `json:"song_title"`
	PublicKey    *crypto.PublicKey `json:"public_key"`
	Signature    *crypto.Signature `json:"signature"`
}

func NewPlayRequest(projectTitle, songTitle string, pub *crypto.PublicKey, sig *crypto.Signature) *PlayRequest {
	return &PlayRequest{
		ProjectTitle: projectTitle,
		SongTitle:    songTitle,
		PublicKey:    pub,
		Signature:    sig,
	}
}

type PlayResponse struct {
	Data  []byte `json:"data, omitempty"`
	Error error  `json:"error, omitempty"`
}

type SocketService struct {
	dir      string
	files    http.FileSystem
	lis      Listener
	logger   Logger
	shutdown int32
}

func NewSocketService(addr, dir string) (*SocketService, error) {
	files := http.Dir(dir)
	if addr == "" {
		addr = LOCALHOST
	}
	lis, err := ListenTCP(addr + ":" + PORT)
	if err != nil {
		return nil, err
	}
	logger := NewLogger("socket_service")
	return &SocketService{
		dir:    dir,
		files:  files,
		lis:    lis,
		logger: logger,
	}, nil
}

func (serv *SocketService) Path(args ...string) string {
	args = append([]string{serv.dir}, args...)
	return filepath.Join(args...)
}

func (serv *SocketService) PlayAddr() string {
	ipAddr, err := GetMyIP()
	Check(err)
	playAddr := ipAddr + ":" + PORT
	return playAddr
}

func (serv *SocketService) Shutdown() error {
	if !atomic.CompareAndSwapInt32(&serv.shutdown, 0, 1) {
		return Error("Socket service is already stopped")
	}
	return nil
}

func (serv *SocketService) AcceptConnections() {
	for {
		if atomic.LoadInt32(&serv.shutdown) == 1 {
			return
		}
		conn, err := serv.lis.Accept()
		if err != nil {
			serv.logger.Error(err.Error())
			continue
		}
		go serv.HandleConn(conn)
	}
}

func (serv *SocketService) HandleConn(conn Conn) {
	playRequest := new(PlayRequest)
	playResponse := new(PlayResponse)
	if err := ReadJSON(conn, playRequest); err != nil {
		serv.logger.Error(err.Error())
		playResponse.Error = err
		WriteJSON(conn, playResponse)
		return
	}
	projectTitle := playRequest.ProjectTitle
	pub := playRequest.PublicKey
	sig := playRequest.Signature
	songTitle := playRequest.SongTitle
	// Verify signature
	if !pub.Verify([]byte(projectTitle+songTitle), sig) {
		err := Error("Signature verification failed")
		serv.logger.Error(err.Error())
		playResponse.Error = err
		WriteJSON(conn, playResponse)
		return
	}
	// TODO: verify payment
	// Find file in project
	dir, err := OpenFile(projectTitle)
	if err != nil {
		serv.logger.Error(err.Error())
		playResponse.Error = err
		WriteJSON(conn, playResponse)
		return
	}
	filenames, err := dir.Readdirnames(0)
	if err != nil {
		serv.logger.Error(err.Error())
		playResponse.Error = err
		WriteJSON(conn, playResponse)
		return
	}
	var filename string
	for _, name := range filenames {
		if strings.Contains(name, songTitle) {
			filename = name
			break
		}
	}
	if filename == "" {
		err = Error("Could not find song with title: " + songTitle)
		serv.logger.Error(err.Error())
		playResponse.Error = err
		WriteJSON(conn, playResponse)
		return
	}
	path := filepath.Join(projectTitle, filename)
	file, err := serv.files.Open(path)
	if err != nil {
		serv.logger.Error(err.Error())
		playResponse.Error = err
		WriteJSON(conn, playResponse)
		return
	}
	playResponse.Data = ReadAll(file)
	WriteJSON(conn, playResponse)
}