package xl8r

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var definedBaseCodecs []Codec[baseContentData, baseHubData] = genBaseCodecs()

func TestDisplayIn_X_Base(t *testing.T) {

	// a demo function that utilizes our Spoke-N-Hub Translator ...
	displayNumberIn := func(requestedBase, originalBase string, numeral string, opts0 ...Opts) (r string, e error) {
		const errMsgUnhandledFmt = "unknown base [ %s ]"
		baseName := func(s string) (r string) {
			r = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s), " ", ""))
			if num, err := strconv.Atoi(r); err == nil {
				r = fmt.Sprintf("base%d",num)
			}
			return
		}

		convertBase, err := New(definedBaseCodecs...)
		if err != nil {
			e = err
			return
		}

		rqstdBase := baseName(requestedBase)
		originBase := baseName(originalBase)

		if !convertBase.Knows(rqstdBase) {
			e = fmt.Errorf(errMsgUnhandledFmt, requestedBase)
			return
		}

		if !convertBase.Knows(originBase) {
			e = fmt.Errorf(errMsgUnhandledFmt, originalBase)
			return
		}
		sourceContent := baseContentData(numeral)
		if toR, toErr := convertBase.To(rqstdBase, originBase, sourceContent, opts0...); toErr == nil {
			r = string(toR)
		} else {
			e = toErr
		}
		return
	}


	tt := []struct {
		number, requestBase, originalBase, expected string
		expectedErr                                 error
		useOpt                                      *Opts
	}{
		{number: "32", requestBase: "10", originalBase: "-10", expected: "", expectedErr: fmt.Errorf(`unknown base [ -10 ]`)},
		{number: "32", requestBase: "Decimal", originalBase: "Base10", expected: "32"},
		{number: "33", requestBase: "Binary", originalBase: "Base10", expected: "100001"},
		{number: "15", requestBase: "hex", originalBase: "Base 10", expected: "f"},
		{number: "f", requestBase: "base2", originalBase: "hexadecimal", expected: "1111"},
		{number: "c0c0c0", requestBase: "hex", originalBase: "base 16", expected: "c0c0c0"},
		{number: "-c0c0c0", requestBase: "hex", originalBase: "base 16", expected: "-c0c0c0"},
		{number: "111", requestBase: "unary", originalBase: "binary", expected: "1111111" },
		{number: "1111111", requestBase: "Decimal", originalBase: "base 1", expected: "7" },
		{number: "-3", requestBase: "Unary", originalBase: "base 10", expected: "", expectedErr: fmt.Errorf(`base 1 can only represent non-negative integers`) },
	}

	for i, tx := range tt {
		var result string
		var err error
		if opt := tx.useOpt; opt != nil {
			result, err = displayNumberIn(tx.requestBase, tx.originalBase, tx.number, *tx.useOpt)
		} else {
			result, err = displayNumberIn(tx.requestBase, tx.originalBase, tx.number)
		}

		assrtEqual(t, tx.expected, result)
		assrtEqual(t, err, tx.expectedErr)
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		t.Logf(`# %d: displayNumberIn("%s","%s","%s") ==>> "%s" %s`,
			i, tx.requestBase, tx.originalBase, tx.number, result, errMsg)
	}
}
