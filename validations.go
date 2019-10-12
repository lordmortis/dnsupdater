package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

const MAX_REQUEST_AGE_IN_SECTIONS = 60

func validTimestamp(timeStamp time.Time) bool {
	now := time.Now()
	earliest := time.Now().Add(-MAX_REQUEST_AGE_IN_SECTIONS * time.Second)
	return timeStamp.Before(now) && timeStamp.After(earliest)
}

func validSignature(request UpdateRequest, secret string) bool {
	hash := md5.New()
	data := fmt.Sprintf("%s\n%d\n", request.Hostname, request.Timestamp.Unix())
	_, err := io.WriteString(hash, fmt.Sprintf("%s\n%d\n", request.Hostname, request.Timestamp.Unix()))
	if err != nil {
		return false
	}

	if len(request.IPV4) > 0 {
		_, err = io.WriteString(hash, fmt.Sprintf("%s\n", request.IPV4))
		if err != nil {
			return false
		}
	}

	data = data + fmt.Sprintf("%s\n", secret)
	_, err = io.WriteString(hash, fmt.Sprintf("%s\n", secret))
	if err != nil {
		return false
	}

	calcSig := fmt.Sprintf("%x", hash.Sum(nil))

	fmt.Printf("Data string: \"%s\"", data)
	fmt.Printf("Comparing: \"%s\" and \"%s\"", calcSig, request.Signature)

	return calcSig == request.Signature
}