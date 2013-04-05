package ast

type compareFunc func(interface{}, interface{}) bool

func compare_eq(lhs, rhs *AstNode, xctx *Context, comp compareFunc) bool {
	lt := lhs.RetType
	rt := rhs.RetType

	if lt != xpath_type_node_set && rt != xpath_type_node_set {
		if lt == xpath_type_boolean || rt == xpath_type_boolean {
			return comp(lhs.eval_boolean(xctx), rhs.eval_boolean(xctx))
		} else if lt == xpath_type_number || rt == xpath_type_number {
			return comp(lhs.eval_number(xctx), rhs.eval_number(xctx))
		} else if (lt == xpath_type_string || rt == xpath_type_string) {
			ls := lhs.eval_string(xctx);
			rs := rhs.eval_string(xctx);
			return comp(ls, rs);
		}
	} else if (lt == xpath_type_node_set && rt == xpath_type_node_set) {
		ls := lhs.eval_node_set(xctx);
		rs := rhs.eval_node_set(xctx);

		for _, li := range ls {
			for _, ri := range rs {
				if comp(li.String(), ri.String()) {
					return true
				}

			}
		}
		return false
	} else {
		if (lt == xpath_type_node_set) {
			tmp := lhs
			lhs = rhs
			rhs = tmp
			tmp2 := lt
			lt = rt
			rt = tmp2
		}

		if (lt == xpath_type_boolean) {
			return comp(lhs.eval_boolean(xctx), rhs.eval_boolean(xctx))
		} else if lt == xpath_type_number {
			l := lhs.eval_number(xctx);
			rs := rhs.eval_node_set(xctx);

			for _, ri := range rs {
				if (comp(l, convert_string_to_number(ri.String()))) {
					return true
				}
			}

			return false
		} else if (lt == xpath_type_string) {
			l := lhs.eval_string(xctx)
			rs := rhs.eval_node_set(xctx)

			for _, ri := range rs {
				if comp(l, ri.String()) {
					return true
				}
			}
			return false
		}
	}
	return false;
}

func compare_rel(lhs, rhs *AstNode, xctx *Context, comp compareFunc) bool {
	lt := lhs.RetType
	rt := rhs.RetType

	if (lt != xpath_type_node_set && rt != xpath_type_node_set) {
		return comp(lhs.eval_number(xctx), rhs.eval_number(xctx))
	} else if (lt == xpath_type_node_set && rt == xpath_type_node_set) {
		ls := lhs.eval_node_set(xctx);
		rs := rhs.eval_node_set(xctx);

		for _, li := range ls {
			l := convert_string_to_number(li.String())

			for _, ri := range rs {
				if comp(l, convert_string_to_number(ri.String())) {
					return true
				}
			}
		}

		return false
	} else if (lt != xpath_type_node_set && rt == xpath_type_node_set) {
		l := lhs.eval_number(xctx);
		rs := rhs.eval_node_set(xctx);
		for _, ri := range rs {
			if (comp(l, convert_string_to_number(ri.String()))) {
				return true
			}
		}
		return false;
	} else if (lt == xpath_type_node_set && rt != xpath_type_node_set) {
		ls := lhs.eval_node_set(xctx)
		r := rhs.eval_number(xctx);

		for _, li := range ls {
			if comp(convert_string_to_number(li.String()), r) {
				return true
			}
		}
		return false;
	} 
	return false;
}

func convert_string_to_number(s string) int {
	return 0
}

