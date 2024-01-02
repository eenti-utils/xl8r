package xl8r

import (
	"fmt"
	"strings"
)

/***********************************************************************
 * PLEASE BE ADVISED:                                                  *
 * The translations represented here are only for feature illustration *
 *  and testing and purposes.                                          *
 * Some of the translations may not be 100% accurate.                  *
 ***********************************************************************/

/****************
 * Codecs Setup *
 ****************/

// all "Origins" for our demo language translator use this data type
//  - when encoding, values of this type get converted to the "hub data" type
//  - when decoding, "hub data" values get converted to this type
type myLanguageContentType string

func (s myLanguageContentType) toWords() []string {
	return strings.Split(string(s), " ")
}

// all "Origins" for our demo language translator use this data type
//  - when encoding, values of the "content type" get converted to this type
//  - when decoding, values of this type get converted to the "content type"
type myLanguageHubDataType []int

func (h myLanguageHubDataType) asArray() (r []int) {
	for _, v := range h {
		r = append(r, v)
	}
	return
}

// low-level helper function
func testContentWithEncMap(v myLanguageContentType, encMap map[string]int, f ...func(word string, wordOK bool)) bool {
	act := func(_ string, _ bool) {}

	if len(f) > 0 {
		act = f[0]
	}
	for _, word := range v.toWords() {
		word = strings.ToLower(strings.TrimSpace(word))
		if _, exists := encMap[word]; !exists {
			act(word, false)
			return false
		}
		act(word, true)
	}
	return true
}

// low-level helper function
func testDecodeUsingMap(v myLanguageHubDataType, decMap map[int]string) (r myLanguageContentType, e error) {
	result := &strings.Builder{}
	for i, n := range v {
		if i > 0 {
			result.WriteString(" ")
		}
		if word, exists := decMap[n]; exists {
			result.WriteString(word)
		} else {
			e = fmt.Errorf("{unknown: %d}", n)
			result.WriteString(e.Error())
		}
	}
	r = myLanguageContentType(result.String())
	return
}

func testFixQuoteChars(s, replaceWith string,replacing ... string) (r string) {
	r = s
	for _, replaceable := range replacing {
		r = strings.ReplaceAll(r,replaceable,replaceWith)
	}
	return
}

// make a *Spoke (ie. a codec) using the specified name and helper maps for encoding and decoding values
func fetchXSpoke(id string, encMap map[string]int, decMap map[int]string) (r *Spoke[myLanguageContentType, myLanguageHubDataType]) {

	eval := func(v myLanguageContentType) bool {
		return testContentWithEncMap(v, encMap)
	}

	r = &Spoke[myLanguageContentType, myLanguageHubDataType]{
		Id: id,
		Enc: func(v myLanguageContentType, opts0 ...Opts) (r myLanguageHubDataType, e error) {
			if len(v) == 0 {
				return
			}
			// create the hub data, from the content, using the encoder map
			testContentWithEncMap(v, encMap, func(word string, wordOK bool) {
				switch wordOK {
				case true:
					r = append(r, encMap[word])
				default:
					e = fmt.Errorf("unknown word: '%s'", word)
				}
			})
			return
		},
		Dec: func(v myLanguageHubDataType, opts0 ...Opts) (r myLanguageContentType, e error) {
			if l := len(v); l > 0 {
				// create the content, from the hub data, using the decoder map
				r, e = testDecodeUsingMap(v, decMap)
			}
			return
		},
		Check: eval,
	}
	return
}

// generate the codec for English numbers
func fetchEngCodec() (r Codec[myLanguageContentType, myLanguageHubDataType]) {
	return fetchXSpoke(
		"english",
		map[string]int{
			"zero":  0,
			"one":   1,
			"two":   2,
			"three": 3,
			"four":  4,
			"five":  5,
			"six":   6,
			"seven": 7,
			"eight": 8,
			"nine":  9,
			"ten":   10,
		},
		map[int]string{
			0:  "zero",
			1:  "one",
			2:  "two",
			3:  "three",
			4:  "four",
			5:  "five",
			6:  "six",
			7:  "seven",
			8:  "eight",
			9:  "nine",
			10: "ten",
		},
	)
}

// generate the codec for Spanish numbers
func fetchSpanishCodec() (r Codec[myLanguageContentType, myLanguageHubDataType]) {
	return fetchXSpoke(
		"spanish",
		map[string]int{
			"cero":   0,
			"uno":    1,
			"dos":    2,
			"tres":   3,
			"cuatro": 4,
			"cinco":  5,
			"seis":   6,
			"siete":  7,
			"ocho":   8,
			"nueve":  9,
			"diez":   10,
		},
		map[int]string{
			0:  "cero",
			1:  "uno",
			2:  "dos",
			3:  "tres",
			4:  "cuatro",
			5:  "cinco",
			6:  "seis",
			7:  "siete",
			8:  "ocho",
			9:  "nueve",
			10: "diez",
		},
	)
}

// generate the codec for Haitian Creole numbers
func fetchHaitianCreoleCodec() (r Codec[myLanguageContentType, myLanguageHubDataType]) {
	return fetchXSpoke(
		"haitian creole",
		map[string]int{
			"zèro": 0,
			"en":   1,
			"de":   2,
			"twa":  3,
			"kat":  4,
			"senk": 5,
			"sis":  6,
			"sèt":  7,
			"wit":  8,
			"nèf":  9,
			"dis":  10,
		},
		map[int]string{
			0:  "zèro",
			1:  "en",
			2:  "de",
			3:  "twa",
			4:  "kat",
			5:  "senk",
			6:  "sis",
			7:  "sèt",
			8:  "wit",
			9:  "nèf",
			10: "dis",
		},
	)
}

// generate the codec for Japanese numbers
func fetchJapaneseCodec() (r Codec[myLanguageContentType, myLanguageHubDataType]) {

	/**********************************************************************
	 * This codec has special handling for optional translations          *
	 * A custom decoder function is used, in place of                     *
	 *  the default decoder provided by the fetchXSpoke(..) function      *
	 * The custom decoder honors the user-defined Dec option called "use" *
	 *  and will employ the appropriate decoder map, as directed.         *
	 *  If a value is not recognized, then the common decoder map is used.*
	 **********************************************************************/

	/***************************************************************
	 * PLEASE BE ADVISED:                                          *
	 * The maps represented here are only for feature illustration *
	 *  and testing and purposes.                                  *
	 * Some of the maps are not 100% accurate.                     *
	 ***************************************************************/

	toCommonMap := map[int]string{
		0:  "rei",
		1:  "ichi",
		2:  "ni",
		3:  "san",
		4:  "shi",
		5:  "go",
		6:  "roku",
		7:  "nana",
		8:  "hachi",
		9:  "kyu",
		10: "ju",
	}
	toKunyomiMap := map[int]string{
		0:  "zero",
		1:  "hito",
		2:  "futa",
		3:  "mi",
		4:  "yon",
		5:  "itsu",
		6:  "mu",
		7:  "nana",
		8:  "ya",
		9:  "kokono",
		10: "to",
	}
	toKunyomi2Map := map[int]string{
		0:  "zero",
		1:  "hito",
		2:  "futa",
		3:  "mi",
		4:  "yo",
		5:  "itsu",
		6:  "mu",
		7:  "nano",
		8:  "yo",
		9:  "kokono",
		10: "so",
	}
	toOnyomiMap := map[int]string{
		0:  "rei",
		1:  "ichi",
		2:  "ni",
		3:  "san",
		4:  "shi",
		5:  "go",
		6:  "roku",
		7:  "shichi",
		8:  "hachi",
		9:  "kyu",
		10: "ju",
	}
	toOnyomi2Map := map[int]string{
		0:  "rei",
		1:  "itsu",
		2:  "ji",
		3:  "zo",
		4:  "shi",
		5:  "go",
		6:  "riku",
		7:  "shichi",
		8:  "hachi",
		9:  "ku",
		10: "ju",
	}
	toKanjiMap := map[int]string{
		0:  "零",
		1:  "一",
		2:  "二",
		3:  "三",
		4:  "四",
		5:  "五",
		6:  "六",
		7:  "七",
		8:  "八",
		9:  "九",
		10: "十",
	}

	spoke := fetchXSpoke(
		"japanese",
		map[string]int{
			"rei":    0,
			"ichi":   1,
			"ni":     2,
			"ji":     2,
			"san":    3,
			"zo":     3,
			"shi":    4,
			"go":     5,
			"roku":   6,
			"riku":   6,
			"shichi": 7,
			"hachi":  8,
			"kyu":    9,
			"ku":     9,
			"ju":     10,
			"zero":   0,
			"hito":   1,
			"futa":   2,
			"mi":     3,
			"yon":    4,
			"itsu":   5,
			"mu":     6,
			"nana":   7,
			"nano":   7,
			"ya":     8,
			"yo":     8,
			"kokono": 9,
			"to":     10,
			"so":     10,
			"零":      0,
			"一":      1,
			"二":      2,
			"三":      3,
			"四":      4,
			"五":      5,
			"六":      6,
			"七":      7,
			"八":      8,
			"九":      9,
			"十":      10,
		},
		toCommonMap,
	)

	// customize the decoder function to honor user-defined options, when received...
	spoke.Dec = func(v myLanguageHubDataType, opts0 ...Opts) (r myLanguageContentType, e error) {
		if l := len(v); l > 0 {

			if len(opts0) == 0 {
				// no user-defined options were specified ...
				return testDecodeUsingMap(v, toCommonMap)
			}

			// check for a recognized user-defined option
			//  (in this case, the decoder option, "use")
			if decoderOpts := opts0[0].Dec; len(decoderOpts) > 0 {
				if useOpt, exists := decoderOpts["use"]; exists {
					switch useVal := useOpt.(type) {
					case string:
						switch useVal {
						case "kunyomi":
							return testDecodeUsingMap(v, toKunyomiMap) // use: "kunyomi"
						case "kunyomi2":
							return testDecodeUsingMap(v, toKunyomi2Map) // use: "kunyomi2"
						case "onyomi":
							return testDecodeUsingMap(v, toOnyomiMap) // use: "onnyomi"
						case "onyomi2":
							return testDecodeUsingMap(v, toOnyomi2Map) // use: "onnyomi2"
						case "kanji":
							return testDecodeUsingMap(v, toKanjiMap) // use: "kanji"
						}
					}
				}
			}

			// we didn't recognize the options that were specified
			//  so use the default map, for decoding...
			return testDecodeUsingMap(v, toCommonMap)
		}
		return
	}
	return spoke
}

// generate the codec for Klingon numbers
func fetchKlingonCodec() (r Codec[myLanguageContentType, myLanguageHubDataType]) {

	fromKlingonMap := map[string]int{
		"pagh":   0,
		"wa’":    1,
		"cha’":   2,
		"wej":    3,
		"loS":    4,
		"vagh":   5,
		"jav":    6,
		"Soch":   7,
		"chorgh": 8,
		"Hut":    9,
		"wa’maH": 10,
	}

	spoke := fetchXSpoke(
		"klingon",
		fromKlingonMap,
		map[int]string{
			0:  "pagh",
			1:  "wa’",
			2:  "cha’",
			3:  "wej",
			4:  "loS",
			5:  "vagh",
			6:  "jav",
			7:  "Soch",
			8:  "chorgh",
			9:  "Hut",
			10: "wa’maH",
		},
	)

	testContentWithKlingonEncMap := func(v myLanguageContentType, f ...func(word string, wordOK bool)) bool {
		act := func(_ string, _ bool) {}

		if len(f) > 0 {
			act = f[0]
		}
		for _, word := range v.toWords() {
			word = testFixQuoteChars(strings.TrimSpace(word),"’","'","`")
			// word = strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(word), "'", "’"), "`", "’")
			if _, exists := fromKlingonMap[word]; !exists {
				act(word, false)
				return false
			}
			act(word, true)
		}
		return true
	}

	//use a custom encoder and a evaluater that handle case sensitivity differently than
	// the default encoder, provided by the fetchXSpoke(...) function
	spoke.Enc = func(v myLanguageContentType, opts0 ...Opts) (r myLanguageHubDataType, e error) {
		if len(v) == 0 {
			return
		}
		// create the hub data, from the content, using the encoder map
		testContentWithKlingonEncMap(v, func(word string, wordOK bool) {
			switch wordOK {
			case true:
				r = append(r, fromKlingonMap[word])
			default:
				e = fmt.Errorf("unknown word: '%s'", word)
			}
		})
		return
	}

	spoke.Check = func(v myLanguageContentType) bool {
		return testContentWithKlingonEncMap(v)
	}

	return spoke
}

// generate the codec for Ga numbers
func fetchGaCodec() (r Codec[myLanguageContentType, myLanguageHubDataType]) {
	return fetchXSpoke(
		"ga",
		map[string]int{
			"ekobɛ":  0,
			"ekome":   1,
			"enyɔ":   2,
			"etɛ": 3,
			"ejwɛ":  4,
			"enumɔ":  5,
			"ekpaa":   6,
			"kpawo": 7,
			"kpaanyɔ": 8,
			"nɛɛhu":  9,
			"nyɔŋma":   10,
		},
		map[int]string{
			0:  "ekobɛ",
			1:  "ekome",
			2:  "enyɔ",
			3:  "etɛ",
			4:  "ejwɛ",
			5:  "enumɔ",
			6:  "ekpaa",
			7:  "kpawo",
			8:  "kpaanyɔ",
			9:  "nɛɛhu",
			10: "nyɔŋma",
		},
	)
}

// generate the codec for Ga numbers
func fetchHawaiianCodec() (r Codec[myLanguageContentType, myLanguageHubDataType]) {

	fromHawaiianMap := map[string]int{
		"῾ole":   0,
		"῾ekahi":    1,
		"akahi": 1,
		"῾elua":   2,
		"῾ekolu":    3,
		"῾ehā":    4,
		"῾elima":   5,
		"῾alima": 5,
		"῾eono":    6,
		"῾ehiku":   7,
		"῾ewalu": 8,
		"῾awalu": 8,
		"῾eiwa":    9,
		"iwa": 9,
		"῾aiwa": 9,
		"῾umi": 10,
	}

	spoke := fetchXSpoke(
		"hawaiian",
		fromHawaiianMap,
		map[int]string{
			0:  "῾ole",
			1:  "῾ekahi",
			2:  "῾elua",
			3:  "῾ekolu",
			4:  "῾ehā",
			5:  "῾elima",
			6:  "῾eono",
			7:  "῾ehiku",
			8:  "῾ewalu",
			9:  "῾eiwa",
			10: "῾umi",
		},
	)
	testContentWithHawaiianEncMap := func(v myLanguageContentType, f ...func(word string, wordOK bool)) bool {
		act := func(_ string, _ bool) {}

		if len(f) > 0 {
			act = f[0]
		}
		for _, word := range v.toWords() {
			word = strings.ToLower(testFixQuoteChars(strings.TrimSpace(word),"῾","'","’","`"))
			if _, exists := fromHawaiianMap[word]; !exists {
				act(word, false)
				return false
			}
			act(word, true)
		}
		return true
	}

	//use a custom encoder and a evaluater that handle case sensitivity differently than
	// the default encoder, provided by the fetchXSpoke(...) function
	spoke.Enc = func(v myLanguageContentType, opts0 ...Opts) (r myLanguageHubDataType, e error) {
		if len(v) == 0 {
			return
		}
		// create the hub data, from the content, using the encoder map
		testContentWithHawaiianEncMap(v, func(word string, wordOK bool) {
			switch wordOK {
			case true:
				r = append(r, fromHawaiianMap[word])
			default:
				e = fmt.Errorf("unknown word: '%s'", word)
			}
		})
		return
	}

	spoke.Check = func(v myLanguageContentType) bool {
		return testContentWithHawaiianEncMap(v)
	}

	return spoke
}
