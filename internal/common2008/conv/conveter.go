package conv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
)

func C2UTF8(_, src, dst string, last int64) (int64, error) {
	f, err := os.Open(src)
	if err != nil {
		return last, fmt.Errorf("error opening source file: %w", err)
	}
	defer f.Close()

	_, err = f.Seek(last, io.SeekStart)
	if err != nil {
		return last, fmt.Errorf("error seeking to last: %w", err)
	}

	inputInfo, err := f.Stat()
	if err != nil {
		return last, fmt.Errorf("error getting file info: %w", err)
	}
	size := inputInfo.Size()

	if size <= last {
		return last, nil
	}

	target, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return last, fmt.Errorf("error opening target file: %w", err)
	}
	defer target.Close()

	writer := bufio.NewWriter(target)
	decoder := simplifiedchinese.GBK.NewDecoder()

	originBuf := make([]byte, 4096)
	for {
		n, err := f.Read(originBuf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}

		last += int64(n)

		reader := transform.NewReader(bytes.NewBuffer(originBuf[:n]), decoder)

		buf := make([]byte, 1024)
		for {
			m, err := reader.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
			}
			_, err = writer.Write(buf[:m])
			if err != nil {
				return last, fmt.Errorf("error writing to file: %w", err)
			}

			if err = writer.Flush(); err != nil {
				return last, fmt.Errorf("error flushing writer: %w", err)
			}
		}

	}

	return last, nil
}

//
//func C2UTF8(enc, src, dst string, last int64) (int64, error) {
//	input, err := os.Open(src)
//	if err != nil {
//		return last, fmt.Errorf("error opening source file: %w", err)
//	}
//	defer input.Close()
//
//	if err = unix.Flock(int(input.Fd()), unix.LOCK_SH); err != nil {
//		return 0, fmt.Errorf("error locking source file: %w", err)
//	}
//	defer unix.Flock(int(input.Fd()), unix.LOCK_UN)
//
//	inputInfo, err := input.Stat()
//	if err != nil {
//		return last, fmt.Errorf("error getting file info: %w", err)
//	}
//	size := inputInfo.Size()
//
//	if size <= last {
//		return last, nil
//	}
//
//	_, err = input.Seek(last, io.SeekStart)
//	if err != nil {
//		return last, fmt.Errorf("error seeking in file: %w", err)
//	}
//
//	target, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
//	if err != nil {
//		return last, fmt.Errorf("error opening target file: %w", err)
//	}
//	defer target.Close()
//
//	var decoder *encoding.Decoder
//	if enc == "" || enc == "gbk" {
//		decoder = simplifiedchinese.GBK.NewDecoder()
//	} else {
//		return last, fmt.Errorf("unsupported encoding formats")
//	}
//
//	reader := transform.NewReader(input, decoder)
//	writer := bufio.NewWriter(target)
//
//	buf := make([]byte, 8*1024*1024)
//	for {
//		n, err := reader.Read(buf)
//		if err != nil {
//			if err == io.EOF {
//				break
//			}
//			return last, fmt.Errorf("error reading from input file: %w", err)
//		}
//		_, err = writer.Write(buf[:n])
//		if err != nil {
//			return last, fmt.Errorf("error writing to target file: %w", err)
//		}
//	}
//
//	if err = writer.Flush(); err != nil {
//		return last, fmt.Errorf("error flushing writer to target file: %w", err)
//	}
//
//	return size, nil
//}

//
//func C2UTF8(enc, src, dst string, last int64) (int64, error) {
//	input, err := os.Open(src)
//	if err != nil {
//		return last, fmt.Errorf("error opening source file: %w", err)
//	}
//	defer input.Close()
//
//	//if err = unix.Flock(int(input.Fd()), unix.LOCK_SH); err != nil {
//	//	return 0, fmt.Errorf("error locking source file: %w", err)
//	//}
//	//defer unix.Flock(int(input.Fd()), unix.LOCK_UN)
//
//	inputInfo, err := input.Stat()
//	if err != nil {
//		return last, fmt.Errorf("error getting file info: %w", err)
//	}
//	size := inputInfo.Size()
//
//	if size <= last {
//		return last, nil
//	}
//
//	_, err = input.Seek(last, io.SeekStart)
//	if err != nil {
//		return last, fmt.Errorf("error seeking in file: %w", err)
//	}
//
//	target, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
//	if err != nil {
//		return last, fmt.Errorf("error opening target file: %w", err)
//	}
//	defer target.Close()
//
//	var decoder *encoding.Decoder
//	if enc == "" || enc == "gbk" {
//		decoder = simplifiedchinese.GBK.NewDecoder()
//	} else {
//		return last, fmt.Errorf("unsupported encoding formats")
//	}
//	//reader := transform.NewReader(input, decoder)
//	//writer := bufio.NewWriter(target)
//
//	total := int64(0)
//	for {
//		originBuf := make([]byte, 4096)
//		n, err := input.Read(originBuf)
//		if err != nil {
//			if errors.Is(err, io.EOF) {
//				break
//			}
//			return last, fmt.Errorf("error reading original file: %w", err)
//		}
//		total += int64(n)
//		reader := transform.NewReader(bytes.NewBuffer(originBuf), decoder)
//		writer := bufio.NewWriter(target)
//		for {
//			buf := make([]byte, 4096)
//			m, err := reader.Read(buf)
//			if err != nil {
//				if err == io.EOF {
//					break
//				}
//				return last, fmt.Errorf("error reading from input file: %w", err)
//			}
//			_, err = writer.Write(buf[:m])
//			if err != nil {
//				return last, fmt.Errorf("error writing to target file: %w", err)
//			}
//		}
//		if err = writer.Flush(); err != nil {
//			return last, fmt.Errorf("error flushing writer to target file: %w", err)
//		}
//
//	}
//
//	//for {
//	//	buf := make([]byte, 8*1024*1024)
//	//	n, err := reader.Read(buf)
//	//	if err != nil {
//	//		if err == io.EOF {
//	//			break
//	//		}
//	//		return last, fmt.Errorf("error reading from input file: %w", err)
//	//	}
//	//	_, err = writer.Write(buf[:n])
//	//	if err != nil {
//	//		return last, fmt.Errorf("error writing to target file: %w", err)
//	//	}
//	//}
//
//	return total + last, nil
//}
