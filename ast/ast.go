package ast

import "math"
import "strings"
import "xpath/node"

type ast_node struct {
	Type byte
	RetType byte

	axis byte
	test byte

	left *ast_node
	right *ast_node
	next *ast_node

	data interface{}
}

func equal_to(a, b interface{}) bool {
	return true
}

func not_equal_to(a, b interface{}) bool {
	return true
}

func less(a, b interface{}) bool {
	return true
}

func less_equal(a, b interface{}) bool {
	return true
}

func (n *ast_node) eval_boolean(xctx *XPathContext) bool {
	switch n.Type {
	case ast_op_or:
		return n.left.eval_boolean(xctx) || n.right.eval_boolean(xctx)
		
	case ast_op_and:
		return n.left.eval_boolean(xctx) && n.right.eval_boolean(xctx)
		
	case ast_op_equal:
		return compare_eq(n.left, n.right, xctx, equal_to)

	case ast_op_not_equal:
		return compare_eq(n.left, n.right, xctx, not_equal_to)

	case ast_op_less:
		return compare_rel(n.left, n.right, xctx, less)
	
	case ast_op_greater:
		return compare_rel(n.right, n.left, xctx, less)

	case ast_op_less_or_equal:
		return compare_rel(n.left, n.right, xctx, less_equal)
	
	case ast_op_greater_or_equal:
		return compare_rel(n.right, n.left, xctx, less_equal)

	case ast_func_starts_with:
		lr := n.left.eval_string(xctx);
		rr := n.right.eval_string(xctx);
		return strings.HasPrefix(lr, rr)

	case ast_func_contains:
		lr := n.left.eval_string(xctx);
		rr := n.right.eval_string(xctx);
		return strings.Contains(lr, rr)

	case ast_func_boolean:
		return n.left.eval_boolean(xctx);
		
	case ast_func_not:
		return !n.left.eval_boolean(xctx);
		
	case ast_func_true:
		return true;
		
	case ast_func_false:
		return false;

	case ast_func_lang:
		if xctx.Node.Attribute() != nil {
			return false
		}
		/*
		xpath_string lang = _left.eval_string(xctx);
		
		for (xml_node n = c.n.node(); n; n = n.parent())
		{
			xml_attribute a = n.attribute(PUGIXML_TEXT("xml:lang"));
			
			if (a)
			{
				const char_t* value = a.value();
				
				// strnicmp / strncasecmp is not portable
				for (const char_t* lit = lang.c_str(); *lit; ++lit)
				{
					if (tolower_ascii(*lit) != tolower_ascii(*value)) return false;
					++value;
				}
				
				return *value == 0 || *value == '-';
			}
		}
		
		return false;
		*/
	case ast_variable:
		/*
		assert(_rettype == _data.variable.type());

		if (_rettype == xpath_type_boolean)
			return _data.variable.get_boolean();
		*/
		fallthrough

	default:
		switch (n.RetType) {
		case xpath_type_number:
			return n.eval_number(xctx) != 0
		case xpath_type_string:
			return len(n.eval_string(xctx)) > 0
		case xpath_type_node_set:				
			return len(n.eval_node_set(xctx)) > 0
		default:
			//assert(!"Wrong expression for return type boolean");
			return false;
		}
	}
	return false
}


func (n *ast_node) eval_string(xctx *XPathContext) string {
	return ""
}

func (n *ast_node) eval_number(xctx *XPathContext) float64 {
	switch n.Type {
	case ast_op_add:
		return n.left.eval_number(xctx) + n.right.eval_number(xctx)
		
	case ast_op_subtract:
		return n.left.eval_number(xctx) - n.right.eval_number(xctx)

	case ast_op_multiply:
		return n.left.eval_number(xctx) * n.right.eval_number(xctx)

	case ast_op_divide:
		return n.left.eval_number(xctx) / n.right.eval_number(xctx)

	case ast_op_mod:
		return math.Mod(n.left.eval_number(xctx), n.right.eval_number(xctx))

	case ast_op_negate:
		return -n.left.eval_number(xctx)

	case ast_number_constant:
		//return _data.number;

	case ast_func_last:
		return float64(xctx.Size)
	
	case ast_func_position:
		return float64(xctx.Position)

	case ast_func_count:
		return float64(len(n.left.eval_node_set(xctx)))
	
	case ast_func_string_length_0:
		return float64(len(xctx.Node.String()))
	
	case ast_func_string_length_1:
		return float64(len(n.left.eval_string(xctx)))
	
	case ast_func_number_0:
		return float64(convert_string_to_number(xctx.Node.String()))
	
	case ast_func_number_1:
		return n.left.eval_number(xctx)

	case ast_func_sum:
		r := 0.
		
		ns := n.left.eval_node_set(xctx)
		
		for _, it := range ns {
			r += float64(convert_string_to_number(it.String()))
		}
	
		return r

	case ast_func_floor:
		r := n.left.eval_number(xctx)
		
		return math.Floor(r)

	case ast_func_ceiling:
		r := n.left.eval_number(xctx)
		
		return math.Ceil(r)

	case ast_func_round:
		r := n.left.eval_number(xctx)
		if r >= -0.5 && r <= 0 {
			return math.Ceil(r)
		} else {
			return math.Floor(r + 0.5)
		}
	
	case ast_variable:
		//assert(_rettype == _data.variable->type());
		/*
		if n.RetType == xpath_type_number {
			return n.data.variable.get_number()
		}
		*/
		// fallthrough to type conversion
		fallthrough

	default:
		switch n.RetType {
		case xpath_type_boolean:
			if n.eval_boolean(xctx) {
				return 1
			} else {
				return 0
			}

		case xpath_type_string:
			return float64(convert_string_to_number(n.eval_string(xctx)))
			
		case xpath_type_node_set:
			return float64(convert_string_to_number(n.eval_string(xctx)))
			
		default:
			//assert(!"Wrong expression for return type number");
			return 0;
		}
		
	}
	return 0
}

func (n *ast_node) eval_node_set(xctx *XPathContext) []*xnode.XPathNode {
	return nil
}


