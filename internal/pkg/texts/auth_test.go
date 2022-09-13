package texts

import (
	"fmt"
	"testing"
	"time"

	"auth/internal/pkg/template"
	"github.com/stretchr/testify/require"
)

const (
	limitMessageInUCS2Encoding = 70
	codeExpiredInterval        = 2 * time.Minute
)

func TestLimits(t *testing.T) {
	m := int(codeExpiredInterval.Minutes())
	interpolation := map[string]any{
		"code":    "1234",
		"minutes": fmt.Sprintf("%d %s", m, Plural(m, `минуту`, `минуты`, `минут`)),
	}
	message := []rune(template.MustInterpolate(AuthCodeSMSBody, interpolation))
	length := len(message)
	require.Truef(
		t,
		length <= limitMessageInUCS2Encoding,
		`message AuthCodeSMSBody length is %d, must be less or equal %d: [%s]`,
		length,
		limitMessageInUCS2Encoding,
		string(message),
	)
}
