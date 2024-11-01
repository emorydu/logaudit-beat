package conv

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os"
)

func C2UTF8(enc, src, dst string, lastLine int) (int, error) {
	input, err := os.Open(src)
	if err != nil {
		return lastLine, fmt.Errorf("error opening source file: %w", err)
	}
	defer input.Close()

	target, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return lastLine, fmt.Errorf("error opening target file: %w", err)
	}
	defer target.Close()

	var decoder *encoding.Decoder
	if enc == "" || enc == "gbk" {
		decoder = simplifiedchinese.GBK.NewDecoder()
	} else {
		return lastLine, fmt.Errorf("unsupported encoding formats")
	}

	reader := transform.NewReader(input, decoder)
	writer := bufio.NewWriter(target)

	scanner := bufio.NewScanner(reader)
	currentLine := 0

	for scanner.Scan() {
		if currentLine < lastLine {
			currentLine++
			continue
		}

		line := scanner.Text()
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return lastLine, fmt.Errorf("error writing to target file: %w", err)
		}
		currentLine++
	}

	if err = writer.Flush(); err != nil {
		return lastLine, fmt.Errorf("error flushing writer to target file: %w", err)
	}

	if err = scanner.Err(); err != nil {
		return lastLine, fmt.Errorf("error reading source file: %w", err)
	}

	return currentLine, nil
}
