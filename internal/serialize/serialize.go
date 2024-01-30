package serialize

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

func Byteify(input interface{}) (bytes.Buffer, error) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(input)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return buf, nil
}

func Unbyteify(reader *gzip.Reader, output interface{}) error {
	return json.NewDecoder(reader).Decode(output)
}

func Zipify(buf bytes.Buffer) bytes.Buffer {
	gzippedData := &bytes.Buffer{}

	gw := gzip.NewWriter(gzippedData)
	gw.Write(buf.Bytes())
	gw.Close()

	return *gzippedData
}

func Unzipify(buf bytes.Buffer) (*gzip.Reader, error) {
	return gzip.NewReader(bytes.NewReader(buf.Bytes()))
}
