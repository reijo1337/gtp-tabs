package main

import (
	"gtp-tabs/GtpStorage/protocol"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Server структура для grpc сервера
type Server struct {
	db *Database
}

// MakeServer возвращает новый объект Server, который представляет определения для grpc
func MakeServer(db *Database) (*Server, error) {
	log.Println("Server: Set up book service...")
	return &Server{db: db}, nil
}

// GetAuthorsByLetter возвращает список музыкантов и количество их исполнителей через поиск по первой букве
func (s *Server) GetAuthorsByLetter(in *protocol.SearchString, p protocol.Tabs_GetAuthorsByLetterServer) error {
	log.Println("New request for searching musicians by letter", in.GetSearch())
	result, err := s.db.getMusiciansByLetter(in.GetSearch())
	if err != nil {
		log.Println("Can't get musicians by letter", in.GetSearch(), "from database.", err)
		return err
	}

	for _, res := range result {
		mwc := &protocol.MusicianWithCount{
			Name:  res.Name,
			Count: res.Count,
		}
		if err := p.Send(mwc); err == nil {
			log.Println("Can't send info about musician", res.Name)
			return err
		}
	}
	log.Println("Request for searcing musicians by letter", in.GetSearch(), "processes succsesfuly")
	return nil
}

// GetAuthorsByName возвращает список музыкантов и количество их исполнителей через поиск по подстроке
func (s *Server) GetAuthorsByName(in *protocol.SearchString, p protocol.Tabs_GetAuthorsByNameServer) error {
	log.Println("New request for searching musicians by substing", in.GetSearch())
	result, err := s.db.getMusicians(in.GetSearch())
	if err != nil {
		log.Println("Can't get musicians by substring", in.GetSearch(), "from database.", err)
		return err
	}

	for _, res := range result {
		mwc := &protocol.MusicianWithCount{
			Name:  res.Name,
			Count: res.Count,
		}
		if err := p.Send(mwc); err == nil {
			log.Println("Can't send info about musician", res.Name, err)
			return err
		}
	}
	log.Println("Request for searcing musicians by letter", in.GetSearch(), "processes succsesfuly")
	return nil
}

// FindTabsByName возвращает список табулатур и количество их исполнителей через поиск по подстроке
func (s *Server) FindTabsByName(in *protocol.SearchString, p protocol.Tabs_FindTabsByNameServer) error {
	log.Println("New request for searching tabs by substring", in.GetSearch())
	results, err := s.db.getTabsByName(in.GetSearch())
	if err != nil {
		return err
	}
	for _, res := range results {
		mwc := &protocol.TabWithSize{
			Musician: res.Musician,
			Name:     res.Name,
			Size:     res.Size,
		}
		if err := p.Send(mwc); err == nil {
			log.Println("Can't send info about tab", res.Name, err)
			return err
		}
	}
	return nil
}

// GetAuthorsByCategory поиск по категориям
func (s *Server) GetAuthorsByCategory(in *protocol.Category, p protocol.Tabs_GetAuthorsByCategoryServer) error {
	log.Println("New request for searching by category", in.GetName())
	results, err := s.db.getMusiciansByCategory(in.GetName())
	if err != nil {
		return err
	}
	for _, res := range results {
		mwc := &protocol.MusicianWithCount{
			Name:  res.Name,
			Count: res.Count,
		}
		if err := p.Send(mwc); err == nil {
			log.Println("Can't send info about musician", res.Name, err)
			return err
		}
	}
	return nil
}

// Upload загрузка файла на сервер через последовательность чанков, которые складываются в файл
func (s *Server) Upload(stream protocol.Tabs_UploadServer) (err error) {
	// while there are messages coming
	buffer := make(map[string][]byte, 0)
	for {
		chunck, err := stream.Recv()
		filename := chunck.GetFilename()
		if _, ok := buffer[filename]; !ok {
			buffer[filename] = make([]byte, 0)
		}
		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Wrapf(err,
				"неудачное чтение чанков")
			return err
		}
		buffer[filename] = append(buffer[filename], chunck.GetContent()...)
	}

	files := make(map[string]*os.File, 0)
	for filename := range buffer {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		files[filename] = file
	}

	for filename, value := range buffer {
		_, err = files[filename].Write(value)
		if err != nil {
			return
		}
	}
	// once the transmission finished, send the
	// confirmation if nothign went wrong
	err = stream.SendAndClose(&protocol.UploadStatus{
		Message: "Файл успешно получен",
		Code:    protocol.UploadStatusCode_Ok,
	})

	return
}

// Download скачивание файла
func (s *Server) Download(in *protocol.DownloadRequest, p protocol.Tabs_DownloadServer) error {
	f, err := os.Open(filepath.Join(in.GetFilename()))
	if err != nil {
		return err
	}
	defer f.Close()

	var b [4096 * 1000]byte
	for {
		n, err := f.Read(b[:])
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		err = p.Send(&protocol.DownloadResponse{
			Data: b[:n],
		})
		if err != nil {
			return err
		}
	}
	return nil
}
