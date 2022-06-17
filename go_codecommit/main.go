package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codecommit"
)

func ExtractTarGz(reader io.Reader, target string) error {
	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	tarReader := tar.NewReader(archive)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	mySession := session.Must(session.NewSession())

	// Create a CodeCommit client from just a session.
	svc := codecommit.New(mySession)

	// Create a CodeCommit client with additional configuration
	//svc := codecommit.New(mySession, aws.NewConfig().WithRegion("ap-northeast-1"))

	// CodeCommitのリポジトリから必要なtargzファイルをダウンロード
	var input codecommit.GetFileInput
	input.SetRepositoryName("my-repo-1")
	input.SetCommitSpecifier("main")
	input.SetFilePath("dummy.tgz")

	output, err := svc.GetFile(&input)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	// targzを指定パスの下に解凍
	var buff bytes.Buffer
	buff.Write(output.FileContent)
	err = ExtractTarGz(&buff, "f:/dev/golang")
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
		f, err := os.Create("f:/dev/golang/dummy2.tgz")
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		defer f.Close()
		defer func() {
			err := f.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()
		f.Write(output.FileContent)
	*/
	fmt.Println("complete")
}
