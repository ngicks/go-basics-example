package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
)

func startTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		bin, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for k, v := range r.Header {
			fmt.Printf("header: %s = %s\n", k, strings.Join(v, ", "))
		}
		fmt.Printf("bytes read: %d\n", len(bin))
		fmt.Printf("body=\n%s\n", bin)
		w.WriteHeader(http.StatusOK)
		return
	}))
}

func sendMultipart(ctx context.Context, url string, client *http.Client) error {
	randomMsg1 := strings.NewReader("foobarbaz")
	randomMgs2 := strings.NewReader("quxquuxcorge")

	/*
		1.
		multipart.NewWriter(w)でio.Writerをラップした*multipart.Writerを得る。
		各種メソッドの呼び出しの結果はwに逐次書き込まれる。
	*/
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	/*
		2.
		CreateFormField, CreateFormFile, CreateFormFieldなど各種メソッドを呼び出すと
		--<<boundary>>
		content-disposition: form-data; name="foo"

		のようなヘッダーを書いたうえで、そのセクションの内容を書き込むためのio.Writerを返す。
	*/
	w, err := mw.CreateFormField("foo")
	if err != nil {
		return err
	}
	_, err = io.Copy(w, randomMsg1)
	if err != nil {
		return err
	}

	w, err = mw.CreateFormFile("bar", "bar.tar.gz")
	if err != nil {
		return err
	}
	_, err = io.Copy(w, randomMgs2)
	if err != nil {
		return err
	}

	/*
		3．
		*multipart.Writerは必ず閉じる。
		閉じなければ最後の--<<boundary>>--が書かれない。
	*/
	err = mw.Close()
	if err != nil {
		return err
	}

	/*
		4.
		http.NewRequestWithContextにデータを受けたbytes.Bufferを渡す。
		(*Writer).FormDataContentType()でboundary込みの"Content-Type"を得られるのせセットする。
	*/
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", mw.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("not 200: code =%d, status = %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func sendMultipartStream(ctx context.Context, url string, client *http.Client) error {
	randomMsg1 := strings.NewReader("foobarbaz")
	randomMgs2 := strings.NewReader("quxquuxcorge")

	writeMultipart := func(mw *multipart.Writer, buf []byte, content1, content2 io.Reader) error {
		var (
			w   io.Writer
			err error
		)

		w, err = mw.CreateFormField("foo")
		if err != nil {
			return err
		}

		_, err = io.CopyBuffer(w, content1, buf)
		if err != nil {
			return err
		}

		w, err = mw.CreateFormFile("bar", "bar.tar.gz")
		if err != nil {
			return err
		}

		_, err = io.CopyBuffer(w, content2, buf)
		if err != nil {
			return err
		}

		err = mw.Close()
		if err != nil {
			return err
		}

		return err
	}

	/*
		1.
		一旦でセクション内容以外を書き出してデータサイズをえる。
		各セクションの内容がファイルなどであれば、
		サイズは既知であるのでこれでContent-Lengthを適切に設定できる。
	*/
	metaData := bytes.Buffer{}
	err := writeMultipart(multipart.NewWriter(&metaData), nil, bytes.NewBuffer(nil), bytes.NewBuffer(nil))
	if err != nil {
		return err
	}
	fmt.Printf("meta data size = %d\n\n", metaData.Len())

	/*
		2.
		io.Pipeでin-memory pipeを取得し、
		別goroutineの中でpw(*io.PipeWriter)にストリーム書き込みする。
		pwに書かれた内容はpr(*io.PipeReader)から読み出せるので、
		これをhttp.NewRequestWithContextに渡す。
	*/
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	defer pr.Close()
	go func() {
		/*
			3.
			goroutineの中で書き込む処理そのものは
			non-stream版とさほど変わらない
		*/
		var err error
		defer func() {
			_ = pw.CloseWithError(err)
		}()
		err = writeMultipart(mw, nil, randomMsg1, randomMgs2)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, pr)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", mw.FormDataContentType())
	/*
		4.
		ContentLengthを設定する。
		req.Header.Set("Content-Length", num)は無視されるので、必ずこちらを設定する。
	*/
	req.ContentLength = int64(metaData.Len()) + randomMsg1.Size() + randomMgs2.Size()

	// Doして終わり。
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("not 200: code =%d, status = %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func main() {
	client := http.DefaultClient
	server := startTestServer()
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("buffer:\n\n")
	if err := sendMultipart(ctx, server.URL, client); err != nil {
		panic(err)
	}

	fmt.Printf("\n")
	fmt.Printf("streaming:\n\n")
	if err := sendMultipartStream(ctx, server.URL, client); err != nil {
		panic(err)
	}
}
