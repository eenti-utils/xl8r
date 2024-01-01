package xl8r

import "fmt"

var _ Codec[int,int] = (*Spoke[int,int])(nil)	//contract

// a named codec that handles conversion of
//   - point content to hub data (encoding)
//   - hub data to point content (decoding)
type Spoke[P, H any] struct {
	// name of the codec
	Id    string
	// function that converts content into hub data (ie. the encoder)
	Enc   Encoder[P, H]
	// function that converts hub data into content (ie. the decoder)
	Dec   Decoder[H, P]
	// function that returns bool true, if the specified content 
	// is processable by the given encoder function
	Check Evaluator[P]
}

// name of the codec
func (s *Spoke[P, H]) Name() string {
	return s.Id
}

// function that converts content into hub data (ie. the encoder)
func (s *Spoke[P, H]) Encode(v P, opts0 ...Opts) (r H, e error) {
	if encode := s.Enc; encode != nil {
		r, e = encode(v, opts0...)
	} else {
		e = fmt.Errorf("nil encoder function")
	}
	return
}

// function that converts hub data into content (ie. the decoder)
func (s *Spoke[P, H]) Decode(v H, opts0 ...Opts) (r P, e error) {
	if decode := s.Dec; decode != nil {
		r, e = decode(v, opts0...)
	} else {
		e = fmt.Errorf("nil decoder function")
	}
	return
}

// function that returns bool true, if the specified content 
// is processable by the given encoder function
func (s *Spoke[P, H]) Evaluate(v P) (r bool) {
	if eval := s.Check; eval != nil {
		r = eval(v)
	}
	return
}
