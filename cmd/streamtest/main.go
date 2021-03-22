package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	streamURL = "http://localhost:8080/v1/Pro520/MARB-108-CAM1.byu.edu/stream"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	errg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < 200; i++ {
		id := strconv.Itoa(i)
		errg.Go(func() error {
			return makeBadRequest(id, ctx)
		})
	}

	if err := errg.Wait(); err != nil {
		fmt.Printf("group error: %s\n", err)
		return
	}

	fmt.Printf("all done!\n")
}

func makeBadRequest(id string, ctx context.Context) error {
	defer fmt.Printf("[%v] done\n", id)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, streamURL, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to do request: %w", err)
	}
	defer resp.Body.Close()

	contentType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	switch {
	case err != nil:
		return fmt.Errorf("unable to parse media type: %w", err)
	case contentType != "multipart/x-mixed-replace":
		return fmt.Errorf("invalid content type: %s", contentType)
	case params["boundary"] == "":
		return fmt.Errorf("no multipart boundary found")
	}

	mr := multipart.NewReader(resp.Body, params["boundary"])

	// read a few frames
	seen := 0
	for {
		part, err := mr.NextPart()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return fmt.Errorf("unable to get next part: %w", err)
		}
		defer part.Close()

		_, err = ioutil.ReadAll(part)
		if err != nil {
			return fmt.Errorf("unable to read part: %w", err)
		}

		seen++
		fmt.Printf("[%v] read %v frames\n", id, seen)

		if seen > 10 {
			fmt.Printf("[%v] stopping read\n", id)
			time.Sleep(10 * time.Minute)
		}
	}
}
