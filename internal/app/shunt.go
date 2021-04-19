package app

import (
	"bufio"
	"bytes"
	_ "embed" //embed for bypass file
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"unsafe"

	"github.com/Asutorufa/yuhaiin/internal/app/component"
	"github.com/Asutorufa/yuhaiin/pkg/net/mapper"
)

//go:embed yuhaiin.conf
var bypassData []byte

type Shunt struct {
	component.Mapper
	file   string
	mapper *mapper.Mapper

	fileLock sync.RWMutex
}

//NewShunt file: bypass file; lookup: domain resolver, can be nil
func NewShunt(file string, lookup func(string) ([]net.IP, error)) (*Shunt, error) {
	s := &Shunt{
		file:   file,
		mapper: mapper.NewMapper(lookup),
	}
	err := s.RefreshMapping()
	if err != nil {
		return nil, fmt.Errorf("refresh mapping failed: %v", err)
	}
	return s, nil
}

func (s *Shunt) RefreshMapping() error {
	s.fileLock.RLock()
	defer s.fileLock.RUnlock()
	fmt.Println(s.file)
	_, err := os.Stat(s.file)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path.Dir(s.file), os.ModePerm)
		if err != nil {
			return fmt.Errorf("make dir all failed: %v", err)
		}
		err = ioutil.WriteFile(s.file, bypassData, os.ModePerm)
		if err != nil {
			return fmt.Errorf("write bypass file failed: %v", err)
		}
	}

	f, err := os.Open(s.file)
	if err != nil {
		return fmt.Errorf("open bypass file failed: %v", err)
	}
	defer f.Close()

	s.mapper.Clear()

	re, _ := regexp.Compile("^([^ ]+) +([^ ]+) *$") // already test that is right regular expression, so don't need to check error
	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if bytes.HasPrefix(a, []byte("#")) {
			continue
		}
		result := re.FindSubmatch(a)
		if len(result) != 3 {
			continue
		}
		mode := component.Mode[strings.ToLower(*(*string)(unsafe.Pointer(&result[2])))]
		if mode == component.OTHERS {
			continue
		}
		s.mapper.Insert(string(result[1]), mode)
	}
	return nil
}

func (s *Shunt) SetFile(f string) error {
	if s.file == f {
		return nil
	}
	s.fileLock.Lock()
	s.file = f
	s.fileLock.Unlock()

	return s.RefreshMapping()
}

func getType(b bool) component.RespType {
	if b {
		return component.IP
	}
	return component.DOMAIN
}

func (s *Shunt) Get(domain string) component.MapperResp {
	mark, markType := s.mapper.Search(domain)
	x, ok := mark.(component.MODE)
	if !ok {
		return component.MapperResp{
			Mark: component.OTHERS,
			IP:   getType(markType),
		}
	}

	if component.ModeMapping[x] == "" {
		x = component.OTHERS
	}

	return component.MapperResp{
		Mark: x,
		IP:   getType(markType),
	}
}

func (s *Shunt) SetLookup(f func(string) ([]net.IP, error)) {
	s.mapper.SetLookup(f)
}