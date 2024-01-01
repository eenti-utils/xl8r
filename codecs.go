package xl8r

type codecMap[P, H any] map[string]Codec[P, H]

func (m *codecMap[P, H]) addCodecs(c0 ...Codec[P, H]) {
	for _, c := range c0 {
		if codecIsValid(c) {
			(*m)[c.Name()] = c
		}
	}
}

func (m *codecMap[P, H]) getIf(name string) (r Codec[P, H], b bool) {
	r, b = (*m)[name]
	return
}

func (m *codecMap[P, H]) keys() (r []string) {
	for k := range *m {
		r = append(r, k)
	}
	return
}

func codecIsValid[P, H any](c Codec[P, H]) (r bool) {
	if c == nil {
		return
	}

	if spoke, isSpoke := c.(*Spoke[P, H]); isSpoke {
		r = len(spoke.Id) > 0 &&
			spoke.Enc != nil &&
			spoke.Dec != nil &&
			spoke.Check != nil
		return
	}

	r = len(c.Name()) > 0
	return
}
