package xl8r

import (
	"fmt"
	"strings"
	"testing"
)

/***********************************************************************
 * PLEASE BE ADVISED:                                                  *
 * The conversions represented here are only for feature illustration  *
 *  and testing and purposes.                                          *
 * Some of the conversions may not be 100% accurate.                   *
 ***********************************************************************/

var definedDurationFmtTestCodecs = []Codec[durationValue, *durationHubData]{
	// order is not important here
	// just using alphabetical for readability
	fetchColonDlmHMSCodec(), // _:_:_ | _ : _ : _
	fetchHHMMSSCodec(),      // _ hh _ mm _ ss  | _hh _mm _ss | _hh_mm_ss
	fetchMinutesCodec(),     // _ min[s|utes]
	// etc ...
}

func TestDisplayIn_X_Format(t *testing.T) {

	// a demo function that utilizes our Spoke-N-Hub Translator ...
	displayDurationIn := func(requestedFormat string, stuff string, opts0 ...Opts) (r string, e error) {
		translateDurationFormat, err := New(definedDurationFmtTestCodecs...)
		if err != nil {
			e = err
			return
		}
		format := strings.ToLower(strings.TrimSpace(requestedFormat))
		if !translateDurationFormat.Knows(format) {
			e = fmt.Errorf("unknown format [ %s ]", requestedFormat)
			return
		}
		sourceContent := durationValue(stuff)
		var sourceFormat string
		if origins := translateDurationFormat.Origins(sourceContent); len(origins) > 0 {
			sourceFormat = origins[0]
		} else {
			e = fmt.Errorf(`no available interpreter for "%s"`, stuff)
			return
		}
		if toR, toErr := translateDurationFormat.To(format, sourceFormat, durationValue(stuff), opts0...); toErr == nil {
			r = string(toR)
		} else {
			e = toErr
		}
		return
	}

	sixtyMins := "60 mins"
	sixtyMinsAndAHalf := "60.5 minute"
	oneHourAndAHalfHMS := "01:30:00"
	mins2600 := "2600 min"

	formattingOpt := &Opts{
		Dec: map[string]any{
			"precision": 2, // we made up parameter that specifies to use 2 decimal places
		},
	}

	tt := []struct {
		durationValue, formatToUse, expected string
		expectedErr                          error
		useOpt                               *Opts
	}{
		{durationValue: sixtyMins, formatToUse: "minutes", expected: "60.000000 minutes"},
		{durationValue: sixtyMins, formatToUse: "hhmmss", expected: "1hh 0mm 0ss"},
		{durationValue: sixtyMins, formatToUse: "hh:mm:ss", expected: "1:0:0"},
		{durationValue: sixtyMinsAndAHalf, formatToUse: "minutes", expected: "60.500000 minutes"},
		{durationValue: sixtyMinsAndAHalf, formatToUse: "hhmmss", expected: "1hh 0mm 30ss"},
		{durationValue: sixtyMinsAndAHalf, formatToUse: "hh:mm:ss", expected: "1:0:30"},
		{durationValue: oneHourAndAHalfHMS, formatToUse: "minutes", expected: "90.000000 minutes"},
		{durationValue: oneHourAndAHalfHMS, formatToUse: "hhmmss", expected: "1hh 30mm 0ss"},
		{durationValue: oneHourAndAHalfHMS, formatToUse: "hh:mm:ss", expected: "1:30:0"},
		{durationValue: mins2600, formatToUse: "minutes", expected: "2600.000000 minutes"},
		{durationValue: mins2600, formatToUse: "hhmmss", expected: "43hh 20mm 0ss"},
		{durationValue: mins2600, formatToUse: "hh:mm:ss", expected: "43:20:0"},

		{durationValue: sixtyMins, formatToUse: "minutes", expected: "60.00 minutes", useOpt: formattingOpt},
		{durationValue: sixtyMinsAndAHalf, formatToUse: "minutes", expected: "60.50 minutes", useOpt: formattingOpt},
		{durationValue: oneHourAndAHalfHMS, formatToUse: "minutes", expected: "90.00 minutes", useOpt: formattingOpt},
		{durationValue: mins2600, formatToUse: "minutes", expected: "2600.00 minutes", useOpt: formattingOpt},
		{durationValue: mins2600, formatToUse: "hhmmss", expected: "43hh 20mm 0ss", useOpt: formattingOpt}, // our "hhmmss" codec ignores options
		{durationValue: mins2600, formatToUse: "hh:mm:ss", expected: "43:20:0", useOpt: formattingOpt},     // our "hh:mm:ss" codec ignores options
	}

	for i, tx := range tt {
		var result string
		var err error
		if opt := tx.useOpt; opt != nil {
			result, err = displayDurationIn(tx.formatToUse, tx.durationValue, *tx.useOpt)
		} else {
			result, err = displayDurationIn(tx.formatToUse, tx.durationValue)
		}

		assrtEqual(t, tx.expected, result)
		assrtEqual(t, err, tx.expectedErr)
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		t.Logf(`# %d: displayDurationIn("%s","%s") ==>> "%s" %s`,
			i, tx.formatToUse, tx.durationValue, result, errMsg)
	}
}
