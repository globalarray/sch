package quiz

import (
	"benzo/internal/repository/repository_model"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	q := repository_model.NewQuiz("test", 123)

	e, err := Encode(q)

	if err != nil {
		t.Fatal(err)
	}

	buf := []byte(e[2:])

	hashSum := fmt.Sprintf("%x", sha256.Sum256([]byte(q.Name)))

	if string(buf[:sha256.BlockSize]) != hashSum {
		t.Errorf("got %s, want %s", string(buf[:64]), hashSum)
	}
}
