package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type Manifest struct {
	Name string
	Seq  int
	Type string
}

var manifestRe, _ = regexp.Compile(`^file\s+(\S+)\s+seq\s+(\d+)\s+type\s+([bhi])\s*$`)

func parseManifest(s string) (Manifest, error) {
	matches := manifestRe.FindStringSubmatch(s)
	if len(matches) != 4 {
		return Manifest{}, fmt.Errorf("incorrect manifest string: `%s`", s)
	}
	seq, _ := strconv.Atoi(matches[2])
	return Manifest{
		Name: matches[1],
		Seq:  seq,
		Type: matches[3],
	}, nil
}

func readManifestFile(dir, name string) (Manifest, error) {
	f, err := os.OpenFile(filepath.Join(dir, name), os.O_RDONLY|os.O_EXCL, 0o644)
	if err != nil {
		return Manifest{}, err
	}
	data, err := bufio.NewReader(f).ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			if len(data) == 0 {
				return Manifest{}, fmt.Errorf("empty manifest file")
			}
		} else {
			return Manifest{}, err
		}
	}
	return parseManifest(data)
}

func createManifest(dir, name string, aofName string) error {
	f, err := os.OpenFile(filepath.Join(dir, name)+".manifest", os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(f, "file %s seq 1 type i", aofName); err != nil {
		return err
	}
	return nil
}

func createAOF(dir, name string) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	newName, err := newAOFFileName(dir, name)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(filepath.Join(dir, newName), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("cannot create file: %w", err)
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return newName, nil
}

func newAOFFileName(aofPath, name string) (string, error) {
	fileNames, err := filepath.Glob(fmt.Sprintf("%s.*.incr.aof", name))
	if err != nil {
		return "", err
	}
	next, err := getNextAOFIncr(fileNames, name)
	if err != nil {
		return "", nil
	}
	newName := fmt.Sprintf("%s.%d.incr.aof", name, next)
	return newName, nil
}

func getNextAOFIncr(prev []string, fname string) (int, error) {
	re, err := regexp.Compile(fmt.Sprintf("^%s\\.(\\d+)\\.incr.aof$", fname))
	if err != nil {
		return -1, err
	}
	res := 0
	for _, s := range prev {
		sm := re.FindStringSubmatch(s)
		if len(sm) == 0 {
			continue
		}
		num, err := strconv.Atoi(sm[1])
		if err != nil {
			return -1, err
		}
		res = max(res, num)
	}
	return 1, nil
	// return res + 1, nil
}

func startAOFWrite(filename string, c chan string) {
	for {
		s := <-c
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
		if err != nil {
			// env.Log.Error("problem opening AOF", "err", err)
			f.Close()
			panic("cannot write manifest")
		}
		if _, err := fmt.Fprint(f, s); err != nil {
			// env.Log.Error("problem writing to AOF", "err", err, "s", s)
			f.Close()
			panic("cannot write manifest")
		}
		f.Close()
	}
}
