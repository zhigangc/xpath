package lexer

const (
	ct_parse_pcdata = 1	// \0, &, \r, <
	ct_parse_attr = 2		// \0, &, \r, ', "
	ct_parse_attr_ws = 4	// \0, &, \r, ', ", \n, tab
	ct_space = 8			// \r, \n, space, tab
	ct_parse_cdata = 16	// \0, ], >, \r
	ct_parse_comment = 32	// \0, -, >, \r
	ct_symbol = 64			// Any symbol > 127, a-z, A-Z, 0-9, _, :, -, .
	ct_start_symbol = 128	// Any symbol > 127, a-z, A-Z, _, :
)

var chartype_table = []int{
	55,  0,   0,   0,   0,   0,   0,   0,      0,   12,  12,  0,   0,   63,  0,   0,   // 0-15
	0,   0,   0,   0,   0,   0,   0,   0,      0,   0,   0,   0,   0,   0,   0,   0,   // 16-31
	8,   0,   6,   0,   0,   0,   7,   6,      0,   0,   0,   0,   0,   96,  64,  0,   // 32-47
	64,  64,  64,  64,  64,  64,  64,  64,     64,  64,  192, 0,   1,   0,   48,  0,   // 48-63
	0,   192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 192, 192, 192, 192, 192, // 64-79
	192, 192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 0,   0,   16,  0,   192, // 80-95
	0,   192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 192, 192, 192, 192, 192, // 96-111
	192, 192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 0, 0, 0, 0, 0,           // 112-127

	192, 192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 192, 192, 192, 192, 192, // 128+
	192, 192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 192, 192, 192, 192, 192,
	192, 192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 192, 192, 192, 192, 192,
	192, 192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 192, 192, 192, 192, 192,
	192, 192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 192, 192, 192, 192, 192,
	192, 192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 192, 192, 192, 192, 192,
	192, 192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 192, 192, 192, 192, 192,
	192, 192, 192, 192, 192, 192, 192, 192,    192, 192, 192, 192, 192, 192, 192, 192,
}

const (
	ctx_special_pcdata = 1   // Any symbol >= 0 and < 32 (except \t, \r, \n), &, <, >
	ctx_special_attr = 2     // Any symbol >= 0 and < 32 (except \t), &, <, >, "
	ctx_start_symbol = 4	  // Any symbol > 127, a-z, A-Z, _
	ctx_digit = 8			  // 0-9
	ctx_symbol = 16			  // Any symbol > 127, a-z, A-Z, 0-9, _, -, .
)
	
var chartypex_table = []int{
	3,  3,  3,  3,  3,  3,  3,  3,     3,  0,  2,  3,  3,  2,  3,  3,     // 0-15
	3,  3,  3,  3,  3,  3,  3,  3,     3,  3,  3,  3,  3,  3,  3,  3,     // 16-31
	0,  0,  2,  0,  0,  0,  3,  0,     0,  0,  0,  0,  0, 16, 16,  0,     // 32-47
	24, 24, 24, 24, 24, 24, 24, 24,    24, 24, 0,  0,  3,  0,  3,  0,     // 48-63

	0,  20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 20, 20, 20, 20, 20,    // 64-79
	20, 20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 0,  0,  0,  0,  20,    // 80-95
	0,  20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 20, 20, 20, 20, 20,    // 96-111
	20, 20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 0,  0,  0,  0,  0,     // 112-127

	20, 20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 20, 20, 20, 20, 20,    // 128+
	20, 20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 20, 20, 20, 20, 20,
	20, 20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 20, 20, 20, 20, 20,
	20, 20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 20, 20, 20, 20, 20,
	20, 20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 20, 20, 20, 20, 20,
	20, 20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 20, 20, 20, 20, 20,
	20, 20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 20, 20, 20, 20, 20,
	20, 20, 20, 20, 20, 20, 20, 20,    20, 20, 20, 20, 20, 20, 20, 20,
}

func isCharType(char byte, ctype int) bool {
	return isCharTypeByTable(char, ctype, chartype_table)
}

func isCharTypeX(char byte, ctype int) bool {
	return isCharTypeByTable(char, ctype, chartypex_table)
}
func isCharTypeByTable(char byte, ctype int, table []int) bool {
	index := int(char)
	if index >= 128 {
		index = 128
	}
	return (table[index] & ctype) > 0
}