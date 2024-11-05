package conv

import (
	"bufio"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func C2UTF8(enc, src, dst string, last int64) (int64, error) {
	input, err := os.Open(src)
	if err != nil {
		return last, fmt.Errorf("error opening source file: %w", err)
	}
	defer input.Close()

	if err = unix.Flock(int(input.Fd()), unix.LOCK_SH); err != nil {
		return 0, fmt.Errorf("error locking source file: %w", err)
	}
	defer unix.Flock(int(input.Fd()), unix.LOCK_UN)

	inputInfo, err := input.Stat()
	if err != nil {
		return last, fmt.Errorf("error getting file info: %w", err)
	}
	size := inputInfo.Size()

	if size <= last {
		return last, nil
	}

	_, err = input.Seek(last, io.SeekStart)
	if err != nil {
		return last, fmt.Errorf("error seeking in file: %w", err)
	}

	target, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return last, fmt.Errorf("error opening target file: %w", err)
	}
	defer target.Close()

	var decoder *encoding.Decoder
	if enc == "" || enc == "gbk" {
		decoder = simplifiedchinese.GBK.NewDecoder()
	} else {
		return last, fmt.Errorf("unsupported encoding formats")
	}

	reader := transform.NewReader(input, decoder)
	writer := bufio.NewWriter(target)

	buf := make([]byte, 8*1024*1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return last, fmt.Errorf("error reading from input file: %w", err)
		}
		_, err = writer.Write(buf[:n])
		if err != nil {
			return last, fmt.Errorf("error writing to target file: %w", err)
		}
	}

	if err = writer.Flush(); err != nil {
		return last, fmt.Errorf("error flushing writer to target file: %w", err)
	}

	return size, nil
}
