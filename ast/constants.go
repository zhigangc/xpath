package ast

const (
	ast_unknown = iota
	ast_op_or						// left or right
	ast_op_and						// left and right
	ast_op_equal					// left = right
	ast_op_not_equal				// left != right
	ast_op_less					// left < right
	ast_op_greater					// left > right
	ast_op_less_or_equal			// left <= right
	ast_op_greater_or_equal		// left >= right
	ast_op_add						// left + right
	ast_op_subtract				// left - right
	ast_op_multiply				// left * right
	ast_op_divide					// left / right
	ast_op_mod						// left % right
	ast_op_negate					// left - right
	ast_op_union					// left | right
	ast_predicate					// apply predicate to set; next points to next predicate
	ast_filter						// select * from left where right
	ast_filter_posinv				// select * from left where right; proximity position invariant
	ast_string_constant			// string constant
	ast_number_constant			// number constant
	ast_variable					// variable
	ast_func_last					// last()
	ast_func_position				// position()
	ast_func_count					// count(left)
	ast_func_id					// id(left)
	ast_func_local_name_0			// local-name()
	ast_func_local_name_1			// local-name(left)
	ast_func_namespace_uri_0		// namespace-uri()
	ast_func_namespace_uri_1		// namespace-uri(left)
	ast_func_name_0				// name()
	ast_func_name_1				// name(left)
	ast_func_string_0				// string()
	ast_func_string_1				// string(left)
	ast_func_concat				// concat(left right siblings)
	ast_func_starts_with			// starts_with(left right)
	ast_func_contains				// contains(left right)
	ast_func_substring_before		// substring-before(left right)
	ast_func_substring_after		// substring-after(left right)
	ast_func_substring_2			// substring(left right)
	ast_func_substring_3			// substring(left right third)
	ast_func_string_length_0		// string-length()
	ast_func_string_length_1		// string-length(left)
	ast_func_normalize_space_0		// normalize-space()
	ast_func_normalize_space_1		// normalize-space(left)
	ast_func_translate				// translate(left right third)
	ast_func_boolean				// boolean(left)
	ast_func_not					// not(left)
	ast_func_true					// true()
	ast_func_false					// false()
	ast_func_lang					// lang(left)
	ast_func_number_0				// number()
	ast_func_number_1				// number(left)
	ast_func_sum					// sum(left)
	ast_func_floor					// floor(left)
	ast_func_ceiling				// ceiling(left)
	ast_func_round					// round(left)
	ast_step						// process set left with step
	ast_step_root					// select root node
)

const (
	axis_ancestor = iota
	axis_ancestor_or_self
	axis_attribute
	axis_child
	axis_descendant
	axis_descendant_or_self
	axis_following
	axis_following_sibling
	axis_namespace
	axis_parent
	axis_preceding
	axis_preceding_sibling
	axis_self
)
	
const (
	nodetest_none = iota
	nodetest_name
	nodetest_type_node
	nodetest_type_comment
	nodetest_type_pi
	nodetest_type_text
	nodetest_pi
	nodetest_all
	nodetest_all_in_namespace
)