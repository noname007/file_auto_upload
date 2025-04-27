package file

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Processor interface {
	Process(ctx context.Context, fileName, filePath string) error
}

type File struct {
	path      string
	processor Processor
}

func NewFile(processor Processor) *File {
	return &File{
		processor: processor,
	}
}

func (f *File) Scan(ctx context.Context, srcDir string, archiveDir string) error {
	now := time.Now()

	files, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("scan %s failed, err: %w", srcDir, err)
	}

	for _, file := range files {
		dstDir := fmt.Sprintf("%s%c%s", archiveDir, os.PathSeparator, now.Format("20060102"))
		if _, err := os.Stat(dstDir); os.IsNotExist(err) {
			err := os.Mkdir(dstDir, os.ModeDir)
			if err != nil {
				return fmt.Errorf("mkdir failed, dir: %s, err: %w", dstDir, err)
			}
		}

		name := file.Name()

		filepath := fmt.Sprintf("%s%c%s", srcDir, os.PathSeparator, name)
		if err := f.processor.Process(ctx, name, filepath); err != nil {
			return fmt.Errorf("process file %s failed, err: %w", name, err)
		}

		dateTimeStr := now.Format("20060102_150405")
		dstFile := fmt.Sprintf("%s%c%s_%s", archiveDir, os.PathSeparator, dateTimeStr, name)
		err = os.Rename(filepath, dstFile)
		if err != nil {
			return fmt.Errorf("mv file from %s to %s failed, err: %w", filepath, dstFile, err)
		}
	}

	return nil
}
