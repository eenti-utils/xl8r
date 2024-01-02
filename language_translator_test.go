package xl8r

import (
	"fmt"
	"strings"
	"testing"
)

/***********************************************************************
 * PLEASE BE ADVISED:                                                  *
 * The translations represented here are only for feature illustration *
 *  and testing and purposes.                                          *
 * Some of the translations may not be 100% accurate.                  *
 ***********************************************************************/

/****************************************************************************************
 * The "Spoke-N-Hub Translator" - Translate from any language to any other language ... *
 *                                                                                      *
 * But don't get too exicited, the codecs defined here are just ... "simpletons"        *
 *  as they only translate the names (and in some cases symbols) of numbers             *
 *  from language to language.                                                          *
 *                                                                                      *
 * Adding a new language codec to the collective would provide translatability          *
 *  between the new language, and all previouly defined langages...                     *
 *                                                                                      *
 *  (apologies to Sci-Fi fans, as that might sound a bit "Borg-ish")                    *
 *                                                                                      *
 ****************************************************************************************/

var definedCodecs = []Codec[myLanguageContentType, myLanguageHubDataType]{
	// order is not important here
	// just using alphabetical for readability
	fetchEngCodec(),           // English
	fetchHaitianCreoleCodec(), // Haitian Creole
	fetchKlingonCodec(),       // Klingon
	fetchJapaneseCodec(),      // Japanese
	fetchSpanishCodec(),       // Spanish
	// etc ...
}

func TestTellMeIn_X_Language(t *testing.T) {

	// a demo function that utilizes our Spoke-N-Hub Translator ...
	tellMeIn := func(requestedlanguage string, stuff string) (r string, e error) {
		spokeNHubTranslateLang, err := New(definedCodecs...)
		if err != nil {
			e = err
			return
		}
		language := strings.ToLower(strings.TrimSpace(requestedlanguage))
		if !spokeNHubTranslateLang.Knows(language) {
			e = fmt.Errorf("unknown language [ %s ]", requestedlanguage)
			return
		}
		sourceContent := myLanguageContentType(stuff)
		var sourceLang string
		if origins := spokeNHubTranslateLang.Origins(sourceContent); len(origins) > 0 {
			sourceLang = origins[0]
		} else {
			e = fmt.Errorf(`no available interpreter for "%s"`, stuff)
			return
		}
		if toR, toErr := spokeNHubTranslateLang.To(language, sourceLang, myLanguageContentType(stuff)); toErr == nil {
			r = string(toR)
		} else {
			e = toErr
		}
		return
	}

	oneToFive := "one two three four five"
	oneToFiveKanji := "一 二 三 四 五"

	tt := []struct {
		thingToSay, languageToUse, expected string
		expectedErr                         error
	}{
		{thingToSay: oneToFive, languageToUse: "English", expected: "one two three four five"},
		{thingToSay: oneToFive, languageToUse: "Haitian Creole", expected: "en de twa kat senk"},
		{thingToSay: oneToFive, languageToUse: "Klingon", expected: "wa’ cha’ wej loS vagh"},
		{thingToSay: oneToFive, languageToUse: "Japanese", expected: "ichi ni san shi go"},
		{thingToSay: oneToFive, languageToUse: "Spanish", expected: "uno dos tres cuatro cinco"},
		{thingToSay: oneToFive, languageToUse: "FooBarBaz", expected: "", expectedErr: fmt.Errorf(`unknown language [ FooBarBaz ]`)},
		{thingToSay: "FooBar", languageToUse: "English", expected: "", expectedErr: fmt.Errorf(`no available interpreter for "FooBar"`)},
		{thingToSay: oneToFiveKanji, languageToUse: "English", expected: "one two three four five"},
		{thingToSay: oneToFiveKanji, languageToUse: "Haitian Creole", expected: "en de twa kat senk"},
		{thingToSay: oneToFiveKanji, languageToUse: "Klingon", expected: "wa’ cha’ wej loS vagh"},
		{thingToSay: oneToFiveKanji, languageToUse: "Japanese", expected: "ichi ni san shi go"},
		{thingToSay: oneToFiveKanji, languageToUse: "Spanish", expected: "uno dos tres cuatro cinco"},
	}

	for i, tx := range tt {
		result, err := tellMeIn(tx.languageToUse, tx.thingToSay)
		assrtEqual(t, tx.expected, result)
		assrtEqual(t, err, tx.expectedErr)
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		t.Logf(`# %d: tellMeIn("%s","%s") ==>> "%s" %s`,
			i, tx.languageToUse, tx.thingToSay, result, errMsg)
	}
}

func TestTranslateLang(t *testing.T) {

	spokeNHubTranslateLang, err := New(definedCodecs...)
	//spokeNHubTranslateLang, err := New[myPointDataType,myHubDataType](definedCodecs...)	//alternately

	assrtNotNil(t, spokeNHubTranslateLang)
	assrtNil(t, err)

	var expectedResult myLanguageContentType = "de senk en"

	result, tErr := spokeNHubTranslateLang.To("haitian creole", "spanish", "dos cinco uno")

	assrtNil(t, tErr)
	assrtEqual(t, expectedResult, result)
	t.Log(result)
}

func TestLangOrigins(t *testing.T) {
	spokeNHubTranslateLang, err := New(definedCodecs...)
	assrtNotNil(t, spokeNHubTranslateLang)
	assrtNil(t, err)

	tt := []struct {
		text     myLanguageContentType
		expected []string
	}{
		{text: "dos cinco uno", expected: []string{"spanish"}},
		{text: "三 五 七", expected: []string{"japanese"}},
		{text: "san go nana", expected: []string{"japanese"}},
		{text: "wej vagh Soch", expected: []string{"klingon"}},
		{text: "wej vagh soch"}, // in our klingon codec, case sensitivity is important!
	}

	for i, tx := range tt {
		origins := spokeNHubTranslateLang.Origins(tx.text)
		assrtEqual(t, tx.expected, origins)
		t.Logf(`# %d: Origin points for content "%s" -- %v`, i, tx.text, origins)
	}
}

func TestLangEncode(t *testing.T) {
	spokeNHubTranslateLang, err := New(definedCodecs...)
	assrtNotNil(t, spokeNHubTranslateLang)
	assrtNil(t, err)
	tt := []struct {
		from     string
		text     myLanguageContentType
		expected myLanguageHubDataType
	}{
		{from: "english", text: "one two three four", expected: []int{1, 2, 3, 4}},
		{from: "japanese", text: "一 二 三 四", expected: []int{1, 2, 3, 4}},
		{from: "japanese", text: "ichi ni san shi", expected: []int{1, 2, 3, 4}},
		{from: "spanish", text: "uno dos tres cuatro", expected: []int{1, 2, 3, 4}},
		{from: "haitian creole", text: "en de twa kat", expected: []int{1, 2, 3, 4}},
		{from: "klingon", text: "wa' cha' wej loS", expected: []int{1, 2, 3, 4}},
	}

	for i, tx := range tt {
		hubDataValue, encodeErr := spokeNHubTranslateLang.Encode(tx.from, tx.text)
		assrtNotNil(t, hubDataValue)
		assrtNil(t, encodeErr)

		assrtEqual(t, tx.expected, hubDataValue)
		t.Logf(`# %d: encode from Origin point "%s" to hub data -- "%s" == %v`, i, tx.from, tx.text, hubDataValue)
	}
}

func TestLangDecode(t *testing.T) {
	spokeNHubTranslateLang, err := New(definedCodecs...)
	assrtNotNil(t, spokeNHubTranslateLang)
	assrtNil(t, err)

	userOptsKanji := &Opts{Dec: map[string]any{"use": "kanji"}}

	tt := []struct {
		to           string
		expected     myLanguageContentType
		hubDataValue myLanguageHubDataType
		opts         *Opts
	}{
		{to: "english", expected: "one two three four", hubDataValue: []int{1, 2, 3, 4}},
		{to: "japanese", expected: "一 二 三 四", hubDataValue: []int{1, 2, 3, 4}, opts: userOptsKanji},
		{to: "japanese", expected: "ichi ni san shi", hubDataValue: []int{1, 2, 3, 4}},
		{to: "spanish", expected: "uno dos tres cuatro", hubDataValue: []int{1, 2, 3, 4}},
		{to: "haitian creole", expected: "en de twa kat", hubDataValue: []int{1, 2, 3, 4}},
		{to: "klingon", expected: "wa’ cha’ wej loS", hubDataValue: []int{1, 2, 3, 4}},
	}

	for i, tx := range tt {
		var decodedResult myLanguageContentType
		var decodeErr error
		if tx.opts == nil {
			decodedResult, decodeErr = spokeNHubTranslateLang.Decode(tx.to, tx.hubDataValue)
		} else {
			decodedResult, decodeErr = spokeNHubTranslateLang.Decode(tx.to, tx.hubDataValue, *tx.opts)
		}

		assrtNotNil(t, decodedResult)
		assrtNil(t, decodeErr)

		assrtEqual(t, tx.expected, decodedResult)
		t.Logf(`# %d: decode from hub data to Destination point "%s" -- %v ==>> "%s"`, i, tx.to, tx.hubDataValue, decodedResult)
	}
}

func TestTranslations(t *testing.T) {

	spokeNHubTranslateLang, err := New(definedCodecs...)
	assrtNotNil(t, spokeNHubTranslateLang)
	assrtNil(t, err)

	userOptsKanji := &Opts{Dec: map[string]any{"use": "kanji"}}
	userOptsKunyomi := &Opts{Dec: map[string]any{"use": "kunyomi"}}
	userOptsOnyomi := &Opts{Dec: map[string]any{"use": "onyomi"}}

	tt := []struct {
		to, from       string
		text, expected myLanguageContentType
		opts           *Opts
	}{
		{to: "haitian creole", from: "english", text: "ten nine eight", expected: "dis nèf wit"},
		{to: "haitian creole", from: "english", text: "Ten Nine Eight", expected: "dis nèf wit"},
		{to: "english", from: "english", text: "TEN nine eight", expected: "ten nine eight"},
		{to: "english", from: "haitian creole", text: "nèf senk twa", expected: "nine five three"},
		{to: "japanese", from: "spanish", text: "cuatro cinco uno cero", expected: "shi go ichi rei"},
		{to: "japanese", from: "spanish", text: "cuatro cinco uno cero", expected: "shi go ichi rei", opts: userOptsOnyomi},
		{to: "japanese", from: "spanish", text: "cuatro cinco uno cero", expected: "yon itsu hito zero", opts: userOptsKunyomi},
		{to: "japanese", from: "spanish", text: "cuatro cinco uno cero", expected: "四 五 一 零", opts: userOptsKanji},
		{to: "english", from: "japanese", text: "四 五 一 零", expected: "four five one zero"},
		{to: "klingon", from: "japanese", text: "四 五 一 零", expected: "loS vagh wa’ pagh"},
		{to: "english", from: "klingon", text: "loS vagh wa’ pagh", expected: "four five one zero"},
	}

	for i, tx := range tt {
		var result myLanguageContentType
		var tErr error
		if tx.opts == nil {
			result, tErr = spokeNHubTranslateLang.To(tx.to, tx.from, tx.text)
		} else {
			result, tErr = spokeNHubTranslateLang.To(tx.to, tx.from, tx.text, *tx.opts)
		}
		assrtNil(t, tErr)
		assrtEqual(t, tx.expected, result)
		t.Logf(`# %d: from %s to %s -- "%s" == "%s"`, i, tx.from, tx.to, tx.text, result)
	}
}
