package xl8r

import (
	"fmt"
	"strconv"
	"strings"
)

type baseContentData string

func (d baseContentData) String() string {
	return strings.ToLower(strings.TrimSpace(string(d)))
}

type baseHubData int64

func (h baseHubData) AsBase(b int) (r string, e error) {
	r = strconv.FormatInt(int64(h), b)
	return
}

// generates many codecs for base to base conversions ...
func genBaseCodecs() (r []Codec[baseContentData, baseHubData]) {
	aliases := make(map[int][]string)
	aliases[2] = []string{"binary"}
	aliases[10] = []string{"dec", "decimal"}
	aliases[16] = []string{"hex", "hexadecimal"}

	for i := 2; i < 37; i++ {
		var ids []string
		ids = append(ids, fmt.Sprintf("base%d", i))
		if moreIds, exist := aliases[i]; exist {
			ids = append(ids, moreIds...)
		}

		enc := testFetchBaseEnc(i)
		dec := testFetchBaseDec(i)
		chk := testFetchBaseChk(i)

		for _, id := range ids {
			r = append(r, &Spoke[baseContentData, baseHubData]{
				Id:    id,
				Enc:   enc,
				Dec:   dec,
				Check: chk,
			})
		}
	}

	// add a special codec for base 1 ...
	chkB1 := func(v baseContentData) (r bool) {
		nStr := v.String()
		if nStr = strings.ReplaceAll(nStr, "1", ""); len(nStr) == 0 {
			r = true
		}
		return
	}
	encB1 := func(v baseContentData, _ ...Opts) (r baseHubData, e error) {
		if chkB1(v) {
			r = baseHubData(len(v.String()))
		} else {
			e = fmt.Errorf("invalid base 1 number [ %s ]",v.String())
		}
		return
	}
	decB1 := func(v baseHubData, _ ...Opts) (r baseContentData, e error) {
		if num := int(v); num > 0 {
			r = baseContentData(strings.Repeat("1", num))
			return
		} else if num == 0 {
			r = baseContentData("")
			return
		}
		e = fmt.Errorf("base 1 can only represent non-negative integers")
		return
	}

	for _, b1Id := range []string{"base1", "unary"} {
		r = append(r, &Spoke[baseContentData, baseHubData]{
			Id:    b1Id,
			Enc:   encB1,
			Dec:   decB1,
			Check: chkB1,
		})
	}
	return
}

func testFetchBaseEnc(b int) func(v baseContentData, _ ...Opts) (r baseHubData, e error) {
	return func(v baseContentData, _ ...Opts) (r baseHubData, e error) {
		hd, err := strconv.ParseInt(v.String(), b, 64)
		r = baseHubData(hd)
		e = err
		return
	}
}

func testFetchBaseDec(b int) func(v baseHubData, _ ...Opts) (r baseContentData, e error) {
	return func(v baseHubData, _ ...Opts) (r baseContentData, e error) {
		if result, err := v.AsBase(b); err == nil {
			r = baseContentData(result)
		} else {
			e = err
		}
		return
	}
}

func testFetchBaseChk(b int) func(v baseContentData) (r bool) {
	return func(v baseContentData) (r bool) {
		if _, err := strconv.ParseInt(v.String(), b, 64); err == nil {
			r = true
		}
		return
	}
}
