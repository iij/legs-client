package socket

import (
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/iij/legs-client/daemon/context"
	"github.com/iij/legs-client/daemon/golum"
	"github.com/iij/legs-client/daemon/handler/socket/dto"
	"github.com/iij/legs-client/daemon/log"
	"github.com/iij/legs-client/util"
	message "github.com/iij/legs-message"
)

// GetHandler is struct for holding context.LegscContext
type GetHandler struct {
	ctx *context.LegscContext
}

// NewGetHandler initialize GetHandler.
func NewGetHandler(ctx *context.LegscContext) *GetHandler {
	return &GetHandler{
		ctx: ctx,
	}
}

// Get handles the request from UNIX domain socket.
func (g *GetHandler) Get(w io.Writer, r *golum.Request) {
	data, err := dto.ParseGetRequest(r.Body)
	if err != nil {
		log.Error("failed to parse request")
		golum.RenderError(w, err)
		return
	}

	action := message.TransferActionGet

	msgData := message.TransferRequestData{
		Action: action,
		Target: data.Target,
	}
	msg := message.NewTransferRequest(data.SessionID, msgData)
	g.ctx.SendMessage(msg)

	g.handleResponse(data, w)
}

func (g *GetHandler) handleResponse(req *dto.GetRequest, w io.Writer) {
	timer := time.NewTimer(20 * time.Second)
	resCh := make(chan []message.TransferResponseData, 1)
	setCh(req.SessionID, resCh)
	for {
		select {
		case res := <-resCh:
			success := true
			var sl saveList = []saveTarget{}
			for _, r := range res {
				if r.StatusCode < 200 || r.StatusCode >= 300 {
					success = false
					continue
				}

				path, err := makeSavePath(r.URL, req.SaveTo)
				if err != nil {
					log.Error(err.Error())
					golum.RenderError(w, err)
					return
				}

				sl.set(path, r.Body)
			}

			if err := sl.save(req.Options.IsDryRun); err != nil {
				log.Error(err.Error())
				golum.RenderError(w, err)
				return
			}

			data := dto.GetResponse{Responses: res}
			b, err := data.ToBinary()
			if err != nil {
				log.Error("failed to parse request")
				golum.RenderError(w, errors.New("failed to parse request"))
				return
			}

			if success {
				golum.RenderResponse(w, b)
			} else {
				golum.RenderError(w, errors.New(string(b)))
			}
			return

		case <-timer.C:
			log.Error("request timeout")
			golum.RenderError(w, NewTimeoutError())
			return
		}
	}
}

type saveTarget struct {
	path string
	data [][]byte
}

func (st *saveTarget) add(data []byte) {
	st.data = append(st.data, data)
}

type saveList []saveTarget

func (l *saveList) searchIndex(path string) (int, bool) {
	for i, st := range *l {
		if st.path == path {
			return i, true
		}
	}
	return 0, false
}

func (l *saveList) set(path string, data []byte) {
	if idx, isExist := l.searchIndex(path); isExist {
		(*l)[idx].add(data)
		return
	}
	st := saveTarget{
		path: path,
		data: [][]byte{data},
	}
	*l = append(*l, st)
}

func (l *saveList) save(isDryRun bool) error {
	if isDryRun {
		return nil
	}
	for _, st := range *l {
		if err := save(st.path, st.data[0]); err != nil {
			return err
		}
		for i := 1; i < len(st.data); i++ {
			path := st.path + "_" + strconv.Itoa(i)
			if err := save(path, st.data[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func save(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

// makeSavePath returns the absolute path of the save destination.
// If saveTo is a directory path, the save destination is /saveTo/URLPathname.
// If saveTo is a file path, the save destination is the absolute path of saveTo.
func makeSavePath(urlStr, saveTo string) (string, error) {
	filePath, err := util.ExpandPath(saveTo)
	if err != nil {
		return "", err
	}
	fInfo, _ := os.Stat(filePath)
	if fInfo == nil {
		return filePath, nil
	}
	if fInfo.IsDir() {
		u, err := url.Parse(urlStr)
		if err != nil {
			return "", err
		}
		_, fileName := filepath.Split(u.Path)

		return filepath.Join(filePath, fileName), nil
	}

	return filePath, nil
}
