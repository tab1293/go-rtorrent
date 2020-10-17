package rtorrent

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tab1293/go-rtorrent/xmlrpc"
)

// RTorrent is used to communicate with a remote rTorrent instance
type RTorrent struct {
	addr         string
	xmlrpcClient *xmlrpc.Client
}

// Torrent represents a torrent in rTorrent
type Torrent struct {
	Hash              string
	Name              string
	Path              string
	Size              int
	Completed         bool
	CompletedBytes    int
	Ratio             float64
	State             int
	DownRate          int
	UpRate            int
	PeersConnected    int
	PeersNotConnected int
	PeersComplete     int
	PeersAccounted    int
	Hashing           int
	ChunkSize         int
	IsMultiFile       bool
}

// File represents a file in rTorrent
type File struct {
	Path            string
	Size            int
	Priority        int
	Index           int
	TotalChunks     int
	ChunksCompleted int
}

// View represents a "view" within RTorrent
type View string

const (
	// ViewMain represents the "main" view, containing all torrents
	ViewMain View = "main"
	// ViewStarted represents the "started" view, containing only torrents that have been started
	ViewStarted View = "started"
	// ViewStopped represents the "stopped" view, containing only torrents that have been stopped
	ViewStopped View = "stopped"
	// ViewHashing represents the "hashing" view, containing only torrents that are currently hashing
	ViewHashing View = "hashing"
	// ViewSeeding represents the "seeding" view, containing only torrents that are currently seeding
	ViewSeeding View = "seeding"
)

// Pretty returns a formatted string representing this Torrent
func (t *Torrent) Pretty() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}

// Pretty returns a formatted string representing this File
func (f *File) Pretty() string {
	b, _ := json.MarshalIndent(f, "", "  ")
	return string(b)
}

// New returns a new instance of `RTorrent`
// Pass in a true value for `insecure` to turn off certificate verification
func New(addr string, insecure bool) *RTorrent {
	return &RTorrent{
		addr:         addr,
		xmlrpcClient: xmlrpc.NewClient(addr, insecure),
	}
}

// WithHTTPClient allows you to a provide a custom http.Client.
func (r *RTorrent) WithHTTPClient(client *http.Client) *RTorrent {
	r.xmlrpcClient = xmlrpc.NewClientWithHTTPClient(r.addr, client)
	return r
}

func (r *RTorrent) Shutdown() error {
	_, err := r.xmlrpcClient.Call("system.shutdown.normal")
	if err != nil {
		return errors.Wrap(err, "system.shutdown.normal XMLRPC call failed")
	}

	return nil
}

func (r *RTorrent) SetSessionDirectory(d string) error {
	_, err := r.xmlrpcClient.Call("session.path.set", d)
	if err != nil {
		return errors.Wrap(err, "session.path.set XMLRPC call failed")
	}
	return nil
}

func (r *RTorrent) SetDefaultDirectory(t Torrent, d string) error {
	fmt.Printf("setting default directory: %s\n", d)
	_, err := r.xmlrpcClient.Call("d.directory.set", t.Hash, d)
	if err != nil {
		return errors.Wrap(err, "d.directory.set XMLRPC call failed")
	}
	return nil
}

// GetTorrent returns the torrent identified by the given hash
func (r *RTorrent) GetTorrent(t Torrent) (Torrent, error) {
	var ret Torrent
	args := []interface{}{"", string(ViewMain), "d.hash=", "d.complete=", "d.completed_bytes=", "d.down.rate=", "d.up.rate=", "d.ratio=", "d.size_bytes=", "d.state=", "d.peers_connected=", "d.name=", "d.base_path=", "d.hashing=", "d.chunk_size=", "d.peers_not_connected=", "d.peers_accounted=", "d.peers_complete=", "d.is_multi_file="}
	results, err := r.xmlrpcClient.Call("d.multicall2", args...)
	if err != nil {
		return ret, err
	}
	for _, outerResult := range results.([]interface{}) {
		for _, innerResult := range outerResult.([]interface{}) {
			d := innerResult.([]interface{})
			hash := d[0].(string)
			if hash == t.Hash {
				ret.Hash = hash
				ret.Completed = d[1].(int) > 0
				ret.CompletedBytes = d[2].(int)
				ret.DownRate = d[3].(int)
				ret.UpRate = d[4].(int)
				ret.Ratio = float64(d[5].(int)) / float64(1000)
				ret.Size = d[6].(int)
				ret.State = d[7].(int)
				ret.PeersConnected = d[8].(int)
				ret.Name = d[9].(string)
				ret.Path = d[10].(string)
				ret.Hashing = d[11].(int)
				ret.ChunkSize = d[12].(int)
				ret.PeersNotConnected = d[13].(int)
				ret.PeersAccounted = d[14].(int)
				ret.PeersComplete = d[15].(int)
				ret.IsMultiFile = d[16].(int) > 0
				return ret, nil
			}
		}
	}
	return ret, errors.New("torrent not found")
}

// AddTorrentURL adds a new torrent by URL
func (r *RTorrent) AddTorrentURL(url string) error {
	_, err := r.xmlrpcClient.Call("load.normal", "", url)
	if err != nil {
		return errors.Wrap(err, "load.normal XMLRPC call failed")
	}
	return nil
}

// AddTorrent adds a new torrent by the torrent files data
func (r *RTorrent) AddTorrent(data []byte) error {
	_, err := r.xmlrpcClient.Call("load.raw", "", data)
	if err != nil {
		return errors.Wrap(err, "load.raw XMLRPC call failed")
	}
	return nil
}

func (r *RTorrent) StartTorrent(t Torrent) error {
	_, err := r.xmlrpcClient.Call("d.start", t.Hash)
	if err != nil {
		return errors.Wrap(err, "d.start XMLRPC call failed")
	}
	return nil
}

func (r *RTorrent) StopTorrent(t Torrent) error {
	_, err := r.xmlrpcClient.Call("d.stop", t.Hash)
	if err != nil {
		return errors.Wrap(err, "d.stop XMLRPC call failed")
	}
	return nil
}

func (r *RTorrent) CloseTorrent(t Torrent) error {
	_, err := r.xmlrpcClient.Call("d.close", t.Hash)
	if err != nil {
		return errors.Wrap(err, "d.close XMLRPC call failed")
	}
	return nil
}

// Delete removes the torrent
func (r *RTorrent) Delete(t Torrent) error {
	_, err := r.xmlrpcClient.Call("d.erase", t.Hash)
	if err != nil {
		return errors.Wrap(err, "d.erase XMLRPC call failed")
	}
	return nil
}

// GetTorrents returns all of the torrents reported by this RTorrent instance
func (r *RTorrent) GetTorrents(view View) ([]Torrent, error) {
	args := []interface{}{"", string(view), "d.name=", "d.size_bytes=", "d.hash=", "d.custom1=", "d.base_path=", "d.is_active=", "d.complete=", "d.ratio="}
	results, err := r.xmlrpcClient.Call("d.multicall2", args...)
	var torrents []Torrent
	if err != nil {
		return torrents, errors.Wrap(err, "d.multicall2 XMLRPC call failed")
	}
	for _, outerResult := range results.([]interface{}) {
		for _, innerResult := range outerResult.([]interface{}) {
			torrentData := innerResult.([]interface{})
			torrents = append(torrents, Torrent{
				Hash:      torrentData[2].(string),
				Name:      torrentData[0].(string),
				Path:      torrentData[4].(string),
				Size:      torrentData[1].(int),
				Completed: torrentData[6].(int) > 0,
				Ratio:     float64(torrentData[7].(int)) / float64(1000),
			})
		}
	}
	return torrents, nil
}

// GetFiles returns all of the files for a given `Torrent`
func (r *RTorrent) GetFiles(t Torrent) ([]File, error) {
	args := []interface{}{t.Hash, 0, "f.path=", "f.size_bytes=", "f.priority=", "f.completed_chunks=", "f.size_chunks="}
	results, err := r.xmlrpcClient.Call("f.multicall", args...)
	var files []File
	if err != nil {
		return files, errors.Wrap(err, "f.multicall XMLRPC call failed")
	}
	for _, outerResult := range results.([]interface{}) {
		for i, innerResult := range outerResult.([]interface{}) {
			fileData := innerResult.([]interface{})
			files = append(files, File{
				Path:            fileData[0].(string),
				Size:            fileData[1].(int),
				Priority:        fileData[2].(int),
				Index:           i,
				ChunksCompleted: fileData[3].(int),
				TotalChunks:     fileData[4].(int),
			})
		}
	}
	return files, nil
}

func (r *RTorrent) SetFilePrority(t Torrent, i int, p int) error {
	_, err := r.xmlrpcClient.Call("f.priority.set", fmt.Sprintf("%s:f%d", t.Hash, i), p)
	if err != nil {
		return errors.Wrap(err, "f.priority.set XMLRPC call failed")
	}
	return nil
}

// DownTotal returns the total downloaded metric reported by this RTorrent instance (bytes)
func (r *RTorrent) DownTotal() (int, error) {
	result, err := r.xmlrpcClient.Call("throttle.global_down.total")
	if err != nil {
		return 0, errors.Wrap(err, "throttle.global_down.total XMLRPC call failed")
	}
	if totals, ok := result.([]interface{}); ok {
		result = totals[0]
	}
	if total, ok := result.(int); ok {
		return total, nil
	}
	return 0, errors.Errorf("result isn't int: %v", result)
}

// DownRate returns the current download rate reported by this RTorrent instance (bytes/s)
func (r *RTorrent) DownRate() (int, error) {
	result, err := r.xmlrpcClient.Call("throttle.global_down.rate")
	if err != nil {
		return 0, errors.Wrap(err, "throttle.global_down.rate XMLRPC call failed")
	}
	if totals, ok := result.([]interface{}); ok {
		result = totals[0]
	}
	if total, ok := result.(int); ok {
		return total, nil
	}
	return 0, errors.Errorf("result isn't int: %v", result)
}

// UpTotal returns the total uploaded metric reported by this RTorrent instance (bytes)
func (r *RTorrent) UpTotal() (int, error) {
	result, err := r.xmlrpcClient.Call("throttle.global_up.total")
	if err != nil {
		return 0, errors.Wrap(err, "throttle.global_up.total XMLRPC call failed")
	}
	if totals, ok := result.([]interface{}); ok {
		result = totals[0]
	}
	if total, ok := result.(int); ok {
		return total, nil
	}
	return 0, errors.Errorf("result isn't int: %v", result)
}

// UpRate returns the current upload rate reported by this RTorrent instance (bytes/s)
func (r *RTorrent) UpRate() (int, error) {
	result, err := r.xmlrpcClient.Call("throttle.global_up.rate")
	if err != nil {
		return 0, errors.Wrap(err, "throttle.global_up.rate XMLRPC call failed")
	}
	if totals, ok := result.([]interface{}); ok {
		result = totals[0]
	}
	if total, ok := result.(int); ok {
		return total, nil
	}
	return 0, errors.Errorf("result isn't int: %v", result)
}

// IP returns the IP reported by this RTorrent instance
func (r *RTorrent) IP() (string, error) {
	result, err := r.xmlrpcClient.Call("network.bind_address")
	if err != nil {
		return "", errors.Wrap(err, "network.bind_address XMLRPC call failed")
	}
	if ips, ok := result.([]interface{}); ok {
		result = ips[0]
	}
	if ip, ok := result.(string); ok {
		return ip, nil
	}
	return "", errors.Errorf("result isn't string: %v", result)
}

// Name returns the name reported by this RTorrent instance
func (r *RTorrent) Name() (string, error) {
	result, err := r.xmlrpcClient.Call("system.hostname")
	if err != nil {
		return "", errors.Wrap(err, "system.hostname XMLRPC call failed")
	}
	if names, ok := result.([]interface{}); ok {
		result = names[0]
	}
	if name, ok := result.(string); ok {
		return name, nil
	}
	return "", errors.Errorf("result isn't string: %v", result)
}

func (r *RTorrent) ListMethods() ([]string, error) {
	results, err := r.xmlrpcClient.Call("system.listMethods")
	if err != nil {
		return []string{}, errors.Wrap(err, "system.listMethods XMLRPC call failed")
	}

	methods := make([]string, 0)
	for _, outerResult := range results.([]interface{}) {
		for _, innerResult := range outerResult.([]interface{}) {
			method := innerResult.(string)
			methods = append(methods, method)
		}
	}

	return methods, nil

}

func (r *RTorrent) MethodSignature(methodName string) (string, error) {
	result, err := r.xmlrpcClient.Call("system.methodHelp", methodName)
	if err != nil {
		return "", errors.Wrap(err, "system.methodHelp XMLRPC call failed")
	}

	fmt.Printf("%+v\n", result)

	if signatures, ok := result.([]interface{}); ok {
		result = signatures[0]
	}
	if signature, ok := result.(string); ok {
		return signature, nil
	}
	return "", errors.Errorf("result isn't string: %v", result)
}
