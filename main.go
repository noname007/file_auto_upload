package main

import (
	"context"
	"fmt"
	"github.com/noname007/file_auto_upload/file"
	"github.com/noname007/file_auto_upload/file/repo"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	srcDir := viper.GetString("srcDir")
	uploadedPath := viper.GetString("uploadedPath")

	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		panic(fmt.Sprintf("srcDir {%s} not exist.", srcDir))
	}

	if _, err = os.Stat(uploadedPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("storage uploaded file's directory not exist: %s", uploadedPath))
	}

	cos, err := repo.NewCos(repo.CosOption{
		SecretIdValue:  viper.GetString("secretIdValue"),
		SecretKeyValue: viper.GetString("secretKeyValue"),
		BucketURL:      viper.GetString("bucketUrl"),
		ServiceURL:     viper.GetString("serviceUrl"),
	})
	fs := file.NewFile(cos)
	wg := sync.WaitGroup{}

	ctx := context.Background()

	go func() {
		wg.Add(1)
		defer wg.Done()
		t := time.NewTicker(3 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				err := fs.Scan(ctx, srcDir, uploadedPath)
				if err != nil {
					fmt.Printf("err: %s\n", err)
				}
			}
		}
	}()

	wg.Wait()
}
