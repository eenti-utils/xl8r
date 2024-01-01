package xl8r

type Interpreter[P, H any] interface {
	// translate the specified content
	//   1. first encoding "Origin -->> Hub"
	//   2. then decoding "Hub -->> Destination"
	To(destination, origin string, content P, opts0 ...Opts) (translatedResult P, e error)
	// translate the specified hub data into content, using the decoder for the specified destination
	Decode(destination string, hubData H, opts0 ...Opts) (translatedResult P, e error)
	// translate the specified content into hub data, using the encoder for the specified origin
	Encode(origin string, content P, opts0 ...Opts) (hubDataResult H, e error)
	// returns the names of all codecs with an encoder function that can process the specified content
	Origins(content0 ...P) (r []string)
	// returns bool true, if the specified codec has been registered
	Knows(name string) (r bool)
}

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
