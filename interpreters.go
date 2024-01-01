package xl8r

import "fmt"

var _ Interpreter[int, int] = (*convertr[int, int])(nil) //contract

type convertr[P, H any] struct {
	codecs codecMap[P, H]
}

// creates a new Interpreter instance based on the specified Codecs
func New[P, H any](codecs ...Codec[P, H]) (r Interpreter[P, H], e error) {
	cdMap := make(codecMap[P, H])
	cdMap.addCodecs(codecs...)

	if numCodecs := len(cdMap); numCodecs < 2 {
		e = fmt.Errorf("need codecs > 1, received [ %d ]", numCodecs)
		return
	}
	r = &convertr[P, H]{
		codecs: cdMap,
	}
	return
}

func (x *convertr[P, H]) getCodecIf(name string) (c Codec[P, H], b bool) {
	return x.codecs.getIf(name)
}

func (x *convertr[P, H]) To(dest, source string, content P, opts0 ...Opts) (r P, e error) {
	if origin, hasOrigin := x.getCodecIf(source); hasOrigin {
		if destination, hasDestination := x.getCodecIf(dest); hasDestination {
			if hubData, err := origin.Encode(content, opts0...); err == nil {
				r, e = destination.Decode(hubData, opts0...)
				return
			} else {
				e = err
			}
		} else {
			e = fmt.Errorf("no decoder [ '%s'<- ]", dest)
		}
	} else {
		e = fmt.Errorf("no encoder [ <-'%s' ]", source)
	}
	return
}

func (x *convertr[P, H]) Decode(dest string, hubData H, opts0 ...Opts) (r P, e error) {
	if destination, hasDestination := x.getCodecIf(dest); hasDestination {
		r, e = destination.Decode(hubData, opts0...)
	} else {
		e = fmt.Errorf("no decoder [ '%s'<- ]", dest)
	}
	return
}

func (x *convertr[P, H]) Encode(source string, content P, opts0 ...Opts) (r H, e error) {
	if origin, hasOrigin := x.getCodecIf(source); hasOrigin {
		r, e = origin.Encode(content, opts0...)
	} else {
		e = fmt.Errorf("no encoder [ <-'%s' ]", source)
	}
	return
}

func (x *convertr[P, H]) Origins(content0 ...P) (r []string) {
	if len(content0) == 0 {
		r = x.codecs.keys()
		return
	}
	codecs := make(codecMap[P, H])
	for _, content := range content0 {
		for name, origin := range x.codecs {
			if _, exists := codecs.getIf(name); !exists && origin.Evaluate(content) {
				codecs[name] = origin
			}
		}
	}
	r = codecs.keys()
	return
}

func (x *convertr[P, H]) Knows(name string) (r bool) {
	_, r = x.getCodecIf(name)
	return
}
