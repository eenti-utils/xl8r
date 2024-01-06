package xl8r

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

/***********************************************************************
 * PLEASE BE ADVISED:                                                  *
 * The conversions represented here are only for feature illustration  *
 *  and testing and purposes.                                          *
 * Some of the conversions may not be 100% accurate.                   *
 ***********************************************************************/

const (
	weekScnds = 604800
	dayScnds  = 86400
	hourScnds = 3600
	minScnds  = 60
)

/************************************************
 * Some codecs for a few duration value formats *
 ************************************************/

func fetchHHMMSSCodec() (r Codec[durationValue, *durationHubData]) {

	// this codec prefers to have duration values in the _ hh _ mm _ ss format

	re := regexp.MustCompile(`(?i)([\d]+)\s*hh\s*([\d]{1,2})\s*mm\s*([\d]{1,2})\s*ss`)
	r = &Spoke[durationValue, *durationHubData]{
		Id: "hhmmss",
		Enc: func(v durationValue, opts0 ...Opts) (r *durationHubData, e error) {
			var hh, mm, ss float64
			if vals, matches := v.matches(re); matches && len(vals) > 3 {
				hh, _ = strconv.ParseFloat(vals[1], 64)
				mm, _ = strconv.ParseFloat(vals[2], 64)
				ss, _ = strconv.ParseFloat(vals[3], 64)
				r = newDurationHubData(&durationParams{H: hh, M: mm, S: ss})
			} else {
				e = fmt.Errorf("could not parse -- %v", v)
			}
			return
		},
		Dec: func(v *durationHubData, opts0 ...Opts) (r durationValue, e error) {
			var hours, minutes, seconds int
			if math.Floor(v.TotalDays()) > 0 {
				hours = int(math.Floor(v.TotalHours()))
			} else {
				hours = v.Hours()
			}
			minutes = v.Minutes()
			seconds = v.Seconds()
			r = durationValue(fmt.Sprintf("%dhh %dmm %dss", hours, minutes, seconds))
			return
		},
		Check: func(v durationValue) (r bool) {
			_, r = v.matches(re)
			return
		},
	}
	return
}

func fetchMinutesCodec() (r Codec[durationValue, *durationHubData]) {

	// this codec knows and cares only about duration values formatted in minutes

	re := regexp.MustCompile(`(?i)([\d]+|[\d]+\.?[\d]+)\s*min`)
	r = &Spoke[durationValue, *durationHubData]{
		Id: "minutes",
		Enc: func(v durationValue, opts0 ...Opts) (r *durationHubData, e error) {
			var minutes float64
			if vals, matches := v.matches(re); matches && len(vals) > 1 {
				minutes, _ = strconv.ParseFloat(vals[1], 64)
				r = newDurationHubData(&durationParams{M: minutes})
			} else {
				e = fmt.Errorf("could not parse -- %v", v)
			}
			return
		},
		Dec: func(v *durationHubData, opts0 ...Opts) (r durationValue, e error) {
			if len(opts0) > 0 {
				opts := opts0[0]
				if decoderOpts := opts.Dec; len(decoderOpts) > 0 {
					if xValue, exists := decoderOpts["precision"]; exists && xValue != nil {
						switch precision := xValue.(type) {
						case int:
							if precision >= 0 {
								fmtStr := fmt.Sprintf("%%.%df minutes", precision)
								r = durationValue(fmt.Sprintf(fmtStr, v.TotalMinutes()))
								return
							}
							// otherwise- we don't recognize (or can't use) the specified options
						}
					}
				}
			}
			r = durationValue(fmt.Sprintf("%f minutes", v.TotalMinutes()))
			return
		},
		Check: func(v durationValue) (r bool) {
			_, r = v.matches(re)
			return
		},
	}
	return
}

func fetchColonDlmHMSCodec() (r Codec[durationValue, *durationHubData]) {

	// this codec prefers to have duration values in the _:_:_ format

	re := regexp.MustCompile(`([\d]+)\s*:\s*([\d]{1,2})\s*:\s*([\d]{1,2})`)
	r = &Spoke[durationValue, *durationHubData]{
		Id: "hh:mm:ss",
		Enc: func(v durationValue, opts0 ...Opts) (r *durationHubData, e error) {
			var hh, mm, ss float64
			if vals, matches := v.matches(re); matches && len(vals) > 3 {
				hh, _ = strconv.ParseFloat(vals[1], 64)
				mm, _ = strconv.ParseFloat(vals[2], 64)
				ss, _ = strconv.ParseFloat(vals[3], 64)
				r = newDurationHubData(&durationParams{H: hh, M: mm, S: ss})
			} else {
				e = fmt.Errorf("could not parse -- %v", v)
			}
			return
		},
		Dec: func(v *durationHubData, opts0 ...Opts) (r durationValue, e error) {
			var hours, minutes, seconds int
			if v.TotalDays() > 0 {
				hours = int(math.Floor(v.TotalHours()))
			} else {
				hours = v.Hours()
			}
			minutes = v.Minutes()
			seconds = v.Seconds()
			r = durationValue(fmt.Sprintf("%d:%d:%d", hours, minutes, seconds))
			return
		},
		Check: func(v durationValue) (r bool) {
			_, r = v.matches(re)
			return
		},
	}
	return
}

/************************
 * Type Definitions for *
 * - Content Data Type  *
 * - Hub Data Type      *
 ************************/

// this is the Content Data Type
type durationValue string

func (v durationValue) String() string {
	return string(v)
}

func (v durationValue) matches(re *regexp.Regexp) (r []string, b bool) {
	if re != nil {
		if m := re.FindStringSubmatch(v.String()); m != nil {
			r = m
			b = true
		}
	}
	return
}

// this is the Hub Data Type
type durationHubData struct {
	weekPrtn,
	dayPrtn,
	hourPrtn,
	minutePrtn,
	secondPrtn int
	secondTotl float64
	//sub- second properties
	millisPrtn,
	microsPrtn,
	nanosPrtn int
}

func newDurationHubData(i durationHubInput) (r *durationHubData) {

	r = new(durationHubData)
	if i == nil {
		return
	}

	remaining := i.Seconds()
	if remaining <= 0 {
		return
	}

	r.secondTotl = remaining

	if total := r.TotalWeeks(remaining); total >= 1 {
		v := math.Floor(total)
		r.weekPrtn = int(v)
		remaining -= (v * weekScnds)
	}

	if remaining == 0 {
		return
	}

	if total := r.TotalDays(remaining); total >= 1 {
		ceiling := 7
		v := math.Floor(total)
		if value := int(v); value < ceiling {
			r.dayPrtn = value
		}
		remaining -= (v * dayScnds)
	}

	if remaining == 0 {
		return
	}

	if total := r.TotalHours(remaining); total >= 1 {
		ceiling := 24
		v := math.Floor(total)
		if value := int(v); value < ceiling {
			r.hourPrtn = value
		}

		remaining -= (v * hourScnds)
	}

	if remaining == 0 {
		return
	}

	if total := r.TotalMinutes(remaining); total >= 1 {
		ceiling := 60
		v := math.Floor(total)
		if value := int(v); value < ceiling {
			r.minutePrtn = value
		}
		remaining -= (v * minScnds)
	}

	if remaining == 0 {
		return
	}

	if total := remaining; total >= 1 {
		ceiling := 60
		v := math.Floor(total)
		if value := int(v); value < ceiling {
			r.secondPrtn = int(math.Floor(total))
		}
	}

	if remaining == 0 {
		return
	}

	// sub-second portions ...
	subS := remaining - math.Floor(remaining)
	if subS == 0 {
		return
	}

	const valMax = 1000
	if total := subS * valMax; total >= 1 {
		ceiling := valMax
		v := math.Floor(total)
		subS = total - v
		if value := int(v); value < ceiling {
			r.millisPrtn = value
		}
	}

	if total := subS * valMax; total >= 1 {
		ceiling := valMax
		v := math.Floor(total)
		subS = total - v
		if value := int(v); value < ceiling {
			r.microsPrtn = value
		}
	}

	if total := subS * valMax; total >= 1 {
		ceiling := valMax
		v := math.Floor(total)
		if value := int(v); value < ceiling {
			r.nanosPrtn = value
		}
	}
	return
}

func (h *durationHubData) calcTotal(standard float64, useSeconds ...float64) (r float64) {
	scnds := h.secondTotl
	if len(useSeconds) > 0 {
		scnds = math.Abs(useSeconds[0])
	}
	if standard > 0 {
		r = scnds / standard
	}
	return
}

func (h *durationHubData) Weeks() (r int) {
	r = h.weekPrtn
	return
}

// returns the entire duration in weeks
func (h *durationHubData) TotalWeeks(useSeconds ...float64) (r float64) {
	r = h.calcTotal(weekScnds, useSeconds...)
	return
}

func (h *durationHubData) Days() (r int) {
	r = h.dayPrtn
	return
}

// returns the entire duration in days
func (h *durationHubData) TotalDays(useSeconds ...float64) (r float64) {
	r = h.calcTotal(dayScnds, useSeconds...)
	return
}

func (h *durationHubData) Hours() (r int) {
	r = h.hourPrtn
	return
}

// returns the entire duration in hours
func (h *durationHubData) TotalHours(useSeconds ...float64) (r float64) {
	r = h.calcTotal(hourScnds, useSeconds...)
	return
}

func (h *durationHubData) Minutes() (r int) {
	r = h.minutePrtn
	return
}

// returns the entire duration in minutes
func (h *durationHubData) TotalMinutes(useSeconds ...float64) (r float64) {
	r = h.calcTotal(minScnds, useSeconds...)
	return
}

func (h *durationHubData) Seconds() (r int) {
	r = h.secondPrtn
	return
}

// returns the entire duration in seconds
func (h *durationHubData) TotalSeconds() (r float64) {
	r = h.calcTotal(1)
	return
}

func (h *durationHubData) MilliSeconds() (r int) {
	r = h.millisPrtn
	return
}

// returns the entire duration in milliseconds
func (h *durationHubData) TotalMilliSeconds(useSeconds ...float64) (r float64) {
	r = h.calcTotal(1/(1*1000), useSeconds...)
	return
}

func (h *durationHubData) MicroSeconds() (r int) {
	r = h.microsPrtn
	return
}

// returns the entire duration in microseconds
func (h *durationHubData) TotalMicroSeconds(useSeconds ...float64) (r float64) {
	r = h.calcTotal(1/(1*1000000), useSeconds...)
	return
}

func (h *durationHubData) NanoSeconds() (r int) {
	r = h.nanosPrtn
	return
}

// returns the entire duration in nanoseconds
func (h *durationHubData) TotalNanoSeconds(useSeconds ...float64) (r float64) {
	r = h.calcTotal(1/(1*1000000000), useSeconds...)
	return
}

// convenience interface
type durationHubInput interface {
	Seconds() (r float64)
}

var _ durationHubInput = (*durationParams)(nil) //contract

// convenience envelope
type durationParams struct {
	//Week
	W,
	//Day
	D,
	//Hour
	H,
	//Minute
	M,
	//Second
	S,
	//Millisecond
	Ms,
	//Microsecond
	Mu,
	//Nanosecond
	N float64
}

func (p *durationParams) Seconds() (r float64) {
	r += math.Abs(p.W * weekScnds)
	r += math.Abs(p.D * dayScnds)
	r += math.Abs(p.H * hourScnds)
	r += math.Abs(p.M * minScnds)
	r += math.Abs(p.S)
	r += math.Abs(p.Ms * (1 / (1000 * 1)))
	r += math.Abs(p.Mu * (1 / (1000 * 1000)))
	r += math.Abs(p.N * (1 / (1000 * 1000 * 1000)))

	return
}
