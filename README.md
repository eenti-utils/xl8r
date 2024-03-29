# xl8r

The xl8r package is a library written in Go that facilitates development of "spoke and hub" translators.

### Ideal Use Case:

xl8r is ideal for building scalable translator utilities to convert values of any type.

Some examples being:
- language text translators
- duration format converters
- base converters

### Usage Example:

```go
  /***********************************
   * assuming user-defined types for *
   *   - LanguageData                *
   *   - CommonHubData               *
   ***********************************/

  var englishCodec, spanishCodec, japaneseCodec xl8r.Codec[LanguageData, CommonHubData]
  // define codecs ...
  // note that each codec has a unique name

  var err error
  var translateLang xl8r.Interpreter[LanguageData, CommonHubData]
  translateLang, err = xl8r.New(englishCodec, spanishCodec, japaneseCodec)

  if err != nil {
    panic(err)  // something failed
  }

  // assuming here that one of the codecs is named "english" and the other is named "japanese"
  result, toErr := translateLang.To("english","japanese","一 二 三 四 五")
  if toErr != nil {
    panic(toErr)  // something failed
  }

  fmt.Sprintf("%v",result)  // one two three four five

  // user-defined options may also be along passed to codecs
  // codecs can be designed to honor options when passed-in
  var useKanji, useOnyomi xl8r.Opts
  // define options here ...

  result2, toErr2 := translateLang.To("japanese","english","one two three four five", useOnyomi)
  if toErr2 != nil {
    panic(toErr2)  // something failed
  }

  fmt.Sprintf("%v",result2)  // ichi ni san shi go

  result3, toErr3 := translateLang.To("japanese","english","one two three four five", useKanji)
  if toErr3 != nil {
    panic(toErr3)  // something failed
  }

  fmt.Sprintf("%v",result3)  // 一 二 三 四 五

```

In the example above, the translator was instantiated with three codecs.

Each codec is responsible for:
- converting its own content data into hub data
- converting hub data into its own content data

Content data can be converted between any of the three codecs.

**Scalability:**
Adding a fourth and fifth codec, makes translations possible between all five codecs, and so forth.

Note the xl8r language translator unit tests provide more detailed examples for usage.

xl8r solves the more mundane problem of code organization, when writing translator or converter utilities.

The more interesting challenge of performing the conversions for a particular use case, are left to developers.

### Spoke and Hub:
```mermaid
flowchart TB;
  X((Point<br>Px))<-->H(((Hub)));
  A((Point<br>P1))<-->H;
  B((Point<br>P2))<-->H;
  C((Point<br>P3))<-->H;
  N((Point<br>Pn))<-->H;

  style H fill:#FFD9B8,stroke:black,color:black
  style X fill:#F0F0F0,stroke:#0000FF,color:#0000FF
  style A fill:#F0F0F0,stroke:#FF9900,color:#FF9900
  style B fill:#F0F0F0,stroke:#cc00ff,color:#cc00ff
  style C fill:#F0F0F0,stroke:#663300,color:#663300
  style N fill:#F0F0F0,stroke:#00F000,color:#00F000

```
<details>

<summary>More on xl8r "spoke and hub" paradigm</summary>

In this paradigm, each point represents a different Origin and/or Destination for data translations.

For example, some content may be translated from `Point P1` to `Point P3`, where:
- `Point P1` is called "english"
- `Point P3` is called "spanish"
- the content to be translated is the value `string` "four"

In this example (and the `xl8r` package), `Point P1` is considered as the _Origin_ and `Point P3` as the _Destination_.

The _Hub_ represents a commonality between _all points_ in the system.
- the hub data, in this example, is the value `int` 4

The _Spoke_ represents the path to and from `Point` and `Hub`.
- from `Point` (_Origin_) to `Hub`, "content data" is converted to "hub data"  (ie. _Encoded_)
- from `Hub` to `Point` (_Destination_), "hub data" is converted to "content data"  (ie. _Decoded_)

Summarizing the "english" to "spanish" translation, in _spoke and hub_ terms:
- from "english": (`Point P1`) convert value `string` "four" to value `int` 4 (`Hub`)
- to "spanish": (`Hub`) convert value `int` 4 to value `string` "cuatro" (`Point P3`)

</details>

```mermaid
graph LR
  Dx(Content Data<br>Px) --> X((Point<br>Px))
  X -- Encode --> H(((Hub Data)))
  H -- Decode --> N((Point<br>Pn))
  N --> Dn(Content Data<br>Pn)

  style H fill:#FFD9B8,stroke:black,color:black
  style Dx fill:#F0F0F0,stroke:#0000FF,color:#0000FF
  style X fill:#F0F0F0,stroke:#0000FF,color:#0000FF
  style N fill:#F0F0F0,stroke:#00F000,color:#00F000
  style Dn fill:#F0F0F0,stroke:#00F000,color:#00F000
```

In theory, translation between any two points within the system is possible.

### Spokes:

The spokes in the xl8r package are called _codecs_, and handle:
- _encoding_ - converting Content Data values to Hub Data values
- _decoding_ - converting Hub Data to Content Data values

```mermaid
graph LR
  P((Point)) -- Encode --> H(((Hub Data)))
  H -- Decode --> P

  style H fill:#FFD9B8,stroke:black,color:black
  style P fill:#F0F0F0,stroke:black,color:black
```

Any object implementing the `xl8r.Codec[P,H any]` interface, is considered by xl8r to be a spoke.

```go
// a Codec handles conversion of
//   - origin content to hub data (encoding)
//   - hub data to destination content (decoding)
type Codec[P, H any] interface {
	// name of the codec
	Name() string
	// function that converts origin content into hub data (ie. the encoder)
	Encode(v P, opts0 ...Opts) (r H, e error)
	// function that converts hub data into destination content (ie. the decoder)
	Decode(v H, opts0 ...Opts) (r P, e error)
	// function that returns bool true, if the specified content
	// is processable by the given encoder function
	Evaluate(v P) (r bool)
}
```

In addition to the xl8r package-provided struct `xl8r.Spoke[P,H any]`, users are free to implement the `xl8r.Codec[P,H any]` interface as desired.

### Content and Hub Data:

**Content Data** is any data that:
- comes inbound from an Origin Point within the system (ie. a value to be translated)
- goes outbound from a Destination Point within the system (ie. the resulting value of something translated)

The _data type_ for content data:
- can be anything
- is the same for _all_ Origin and Destination Points within the system
- may differ from (but need not be different from) the hub data type
- supports the diverse values of all points within the system
- is represented generically in xl8r code as the letter **P**

**Hub Data** is the pivot between all points within the system.

The _data type_ for hub data:
- can be anything
- supports values that are understood by all points within the system
- is represented generically in xl8r code as the letter **H**
  
```go
	type myContentDataType string
	type myHubDataType int

	// create a Codec to process myContentDataType and myHubDataType -- xl8r.Spoke[P,H any]
	myAwesomeCodec := &xl8r.Spoke[myContentDataType, myHubDataType] {
		Id: "awesomeness",	// a unique name, within the system ...
		Enc: func(v myContentDataType, opts0 ...Opts) (r myHubDataType, e error) {
			// define the  greatest encoder, here...
		},
		Dec: (v myHubDataType, opts0 ...Opts) (r myContentDataType, e error) {
			// define the  greatest decoder, here...
		},
		Check: (v myContentDataType) (r bool) {
			// define the content data evaluator, here...
		},
	}
```


