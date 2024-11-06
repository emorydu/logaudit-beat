package conv

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/emorydu/dbaudit/internal/common/utils"
	"os"
	"path/filepath"
	"strings"
)

func DiffPosition(rootPath string, convpath []string) error {
	positionFile := filepath.Join(rootPath, "position")
	m := make(map[string]struct{})

	var err error

	for _, path := range convpath {
		if !strings.HasSuffix(path, "*") {
			m[fmt.Sprintf(fmt.Sprintf("%s######%s.utf8######", path, path))] = struct{}{}
		} else {
			path = path[:len(path)-len("*")]
			err = utils.EnsureDir(path + "utf8")
			if err != nil && !errors.Is(err, utils.ErrAlreadyExists) {
				return err
			}

			err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				// todo sub-dir files...
				if !info.IsDir() {
					fmt.Println(p)
					pp := filepath.Join(path+"utf8", info.Name()+".utf8")
					m[fmt.Sprintf(fmt.Sprintf("%s######%s######", p, pp))] = struct{}{}
				}

				return nil
			})

		}

	}
	if err != nil {
		return err
	}

	existing := make([]string, 0)
	pf, err := os.Open(positionFile)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(pf)
	for scanner.Scan() {
		existing = append(existing, strings.TrimSpace(scanner.Text()))

	}
	if err = scanner.Err(); err != nil {
		return err
	}
	pf.Close()

	output, err := os.OpenFile(positionFile, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer output.Close()

	writer := bufio.NewWriter(output)

	for item := range m {
		if !contains(existing, item) {
			_, err = writer.WriteString(fmt.Sprintf("%s%d", item, 0) + "\n")
			if err != nil {
				return err
			}
		}
	}

	for _, item := range existing {
		spans := strings.Split(item, "######")
		if len(spans) != 3 {
			return fmt.Errorf("position file %s contains invalid syntax", positionFile)
		}
		if _, exists := m[fmt.Sprintf("%s######%s######", spans[0], spans[1])]; exists {
			_, err = writer.WriteString(item + "\n")
			if err != nil {
				return err
			}
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}

func contains(l []string, item string) bool {
	for _, v := range l {
		if strings.HasPrefix(v, item) {
			return true
		}
	}

	return false
}
