package xl8r

// a function that converts the specified content for a given point
// into hub data.
//   - returns the hub data and a nil error, if successful
//   - returns a zero value and a non-nil error, if the conversion was not possible
type Encoder[P, H any] func(v P, opts0 ...Opts) (r H, e error)

// a function that converts the specified hub data into content, for a given point
//   - returns a zero value and a non-nil error, if the conversion was not possible
type Decoder[H, P any] func(v H, opts0 ...Opts) (r P, e error)

// a function that returns bool true if the specified content
// is processable by the given Encoder function
type Evaluator[P any] func(v P) (r bool)


// user-defined options for encoder and decoder functions
type Opts struct {
	Enc map[string]any
	Dec map[string]any
}
