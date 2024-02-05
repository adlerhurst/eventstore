package eventstorev1alpha

import "strings"

type indexes []*index

func (idxs indexes) byName(names ...string) (res []*index) {
	for _, idx := range idxs {
		for _, name := range names {
			if idx.flag != name {
				continue
			}
			res = append(res, idx)
			break
		}
	}
	return res
}

func (idxs indexes) lastByName(names ...string) *index {
	for i := len(idxs) - 1; i >= 0; i-- {
		for _, name := range names {
			if idxs[i].flag != name {
				continue
			}
			return idxs[i]
		}
	}
	return nil
}

func (idxs indexes) last() *index {
	return idxs[len(idxs)-1]
}

func (idxs indexes) primitives() *index {
	return idxs.lastByName("")
}

type index struct {
	flag string
	args []string
}

func fieldIndexes(args []string, names ...string) (res indexes) {
	// index for primitive fields
	res = append(res, new(index))
	for _, arg := range args {
		flag, found := strings.CutPrefix(arg, "-")
		if !found {
			res[len(res)-1].args = append(res[len(res)-1].args, arg)
			continue
		}
		// also trim second -
		flag = strings.TrimPrefix(flag, "-")

		if isFieldFlag(flag, names) {
			res = append(res, &index{flag: flag})
		} else {
			res[len(res)-1].args = append(res[len(res)-1].args, arg)
		}
	}

	return res
}

func isFieldFlag(arg string, names []string) bool {
	for _, name := range names {
		if arg != name {
			continue
		}
		return true
	}
	return false
}
