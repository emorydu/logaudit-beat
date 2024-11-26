package main

//
//import (
//	"bufio"
//	"fmt"
//	"os"
//	"strings"
//)
//
//func diffPosition(rootPath string, convpath []string) error {
//	positionFile := rootPath + "position"
//	m := make(map[string]struct{})
//	for _, path := range convpath {
//		m[fmt.Sprintf(fmt.Sprintf("%s######%s.utf8######", path, path))] = struct{}{}
//	}
//
//	existing := make([]string, 0)
//	pf, err := os.Open(positionFile)
//	if err != nil {
//		return err
//	}
//	scanner := bufio.NewScanner(pf)
//	for scanner.Scan() {
//		existing = append(existing, strings.TrimSpace(scanner.Text()))
//
//	}
//	if err = scanner.Err(); err != nil {
//		return err
//	}
//	pf.Close()
//
//	output, err := os.OpenFile(positionFile, os.O_WRONLY|os.O_TRUNC, 0644)
//	if err != nil {
//		return err
//	}
//	defer output.Close()
//
//	writer := bufio.NewWriter(output)
//
//	for item := range m {
//		if !contains(existing, item) {
//			_, err = writer.WriteString(fmt.Sprintf("%s%d", item, 0) + "\n")
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	for _, item := range existing {
//		spans := strings.Split(item, "######")
//		if len(spans) != 3 {
//			return fmt.Errorf("position file %s contains invalid syntax", positionFile)
//		}
//		if _, exists := m[fmt.Sprintf("%s######%s######", spans[0], spans[1])]; exists {
//			_, err = writer.WriteString(item + "\n")
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	if err = writer.Flush(); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func contains(l []string, item string) bool {
//	for _, v := range l {
//		if strings.HasPrefix(v, item) {
//			return true
//		}
//	}
//
//	return false
//}
//
//func main() {
//	diffPosition("", []string{"/var/logs/xxxx.log", "/var/logs/command.log"})
//}
