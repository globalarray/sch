package quiz

import (
	"benzo/internal/repository"
	"benzo/internal/repository/repository_model"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
)

const ShortestBlockSize int = sha256.BlockSize / 4

var (
	ErrQuizNotFound   = errors.New("quiz not found")
	ErrInvalidHashSum = errors.New("invalid hash sum")
)

var (
	cachedEncodedQuizzes = map[int64]string{}
	cachedDecodedQuizzes = map[string]int64{}
)

func Encode(q repository_model.Quiz) (string, error) {
	if k, ok := cachedEncodedQuizzes[q.ID]; ok {
		return k, nil
	}

	buf := bytes.NewBuffer([]byte{})

	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(q.Name)))[:ShortestBlockSize]

	if _, err := buf.WriteString(sum); err != nil {
		return "", err
	}

	if err := binary.Write(buf, binary.LittleEndian, q.ID); err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	cachedEncodedQuizzes[q.ID] = encoded
	cachedDecodedQuizzes[encoded] = q.ID

	return encoded, nil
}

func Decode(key string) (int64, error) {
	if id, ok := cachedDecodedQuizzes[key]; ok {
		return id, nil
	}

	buf, err := base64.StdEncoding.DecodeString(key)

	if err != nil {
		return 0, err
	}

	hashSum := buf[:ShortestBlockSize]

	r := bytes.NewBuffer(buf[ShortestBlockSize:])

	var id int64

	if err := binary.Read(r, binary.LittleEndian, &id); err != nil {
		return id, err
	}

	q, err := repository.Repo().GetQuizByID(id)

	if err != nil {
		return 0, err
	}

	if q.ID != id {
		return 0, ErrQuizNotFound
	}

	if fmt.Sprintf("%x", sha256.Sum256([]byte(q.Name)))[:ShortestBlockSize] != string(hashSum) {
		return 0, ErrInvalidHashSum
	}

	cachedEncodedQuizzes[q.ID] = key
	cachedDecodedQuizzes[key] = q.ID

	return q.ID, nil
}
