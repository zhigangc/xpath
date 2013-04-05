package lexer

const (
	lex_none = iota
	lex_equal
	lex_not_equal
	lex_less
	lex_greater
	lex_less_or_equal
	lex_greater_or_equal
	lex_plus
	lex_minus
	lex_multiply
	lex_union
	lex_var_ref
	lex_open_brace
	lex_close_brace
	lex_quoted_string
	lex_number
	lex_slash
	lex_double_slash
	lex_open_square_brace
	lex_close_square_brace
	lex_string
	lex_comma
	lex_axis_attribute
	lex_dot
	lex_double_dot
	lex_double_colon
	lex_eof
)

const strterm = 0

type lexeme_t int

type Lexer struct {
	query string
	cur int
	cur_lexeme_pos int
	cur_lexeme_content_begin int
	cur_lexeme_content_end int
	cur_lexeme lexeme_t
}

func NewLexer(query string) *Lexer {
	lx := &Lexer{query: query}
	lx.next()
	return lx
}

func (lx *Lexer) at(i int) byte {
	if i < len(lx.query) {
		return lx.query[i]
	}
	return strterm
}

func (lx *Lexer) next() {
	cur := lx.cur

	for ; isCharType(lx.at(cur), ct_space); cur++ {}

	// save lexeme position for error reporting
	lx.cur_lexeme_pos = cur

	switch lx.at(cur) {
	case strterm:
		lx.cur_lexeme = lex_eof
	case '>':
		if lx.at(cur+1) == '=' {
			cur += 2
			lx.cur_lexeme = lex_greater_or_equal
		} else {
			cur += 1
			lx.cur_lexeme = lex_greater
		}
	case '<':
		if lx.at(cur+1) == '=' {
			cur += 2
			lx.cur_lexeme = lex_less_or_equal
		} else {
			cur += 1
			lx.cur_lexeme = lex_less
		}
	case '!':
		if lx.at(cur+1) == '=' {
			cur += 2
			lx.cur_lexeme = lex_not_equal
		} else {
			lx.cur_lexeme = lex_none
		}
	case '=':
		cur += 1
		lx.cur_lexeme = lex_equal
	case '+':
		cur += 1
		lx.cur_lexeme = lex_plus
	case '-':
		cur += 1
		lx.cur_lexeme = lex_minus
	case '*':
		cur += 1
		lx.cur_lexeme = lex_multiply
	case '|':
		cur += 1
		lx.cur_lexeme = lex_union
	case '$':
		cur += 1
		if isCharTypeX(lx.at(cur), ctx_start_symbol) {
			lx.cur_lexeme_content_begin = cur

			for ; isCharTypeX(lx.at(cur), ctx_symbol); cur++ {}

			if lx.at(cur) == ':' && isCharTypeX(lx.at(cur+1), ctx_symbol) { // qname
				cur ++ // :
				for ; isCharTypeX(lx.at(cur), ctx_symbol); cur++ {}
			}

			lx.cur_lexeme_content_end = cur
			lx.cur_lexeme = lex_var_ref
		} else {
			lx.cur_lexeme = lex_none
		}
	case '(':
		cur += 1
		lx.cur_lexeme = lex_open_brace
	case ')':
		cur += 1
		lx.cur_lexeme = lex_close_brace
	case '[':
		cur += 1
		lx.cur_lexeme = lex_open_square_brace
	case ']':
		cur += 1
		lx.cur_lexeme = lex_close_square_brace
	case ',':
		cur += 1
		lx.cur_lexeme = lex_comma
	case '/':
		if lx.at(cur+1) == '/' {
			cur += 2
			lx.cur_lexeme = lex_double_slash
		} else {
			cur += 1
			lx.cur_lexeme = lex_slash
		}
	case '.':
		if lx.at(cur+1) == '.' {
			cur += 2
			lx.cur_lexeme = lex_double_dot
		} else if isCharTypeX(lx.at(cur+1), ctx_digit) {
			lx.cur_lexeme_content_begin = cur // .
			cur ++
			for ; isCharTypeX(lx.at(cur), ctx_digit); cur++ {}
			lx.cur_lexeme_content_end = cur
			lx.cur_lexeme = lex_number
		} else {
			cur += 1;
			lx.cur_lexeme = lex_dot
		}
	case '@':
		cur += 1;
		lx.cur_lexeme = lex_axis_attribute;
	case '"', '\'':
		terminator := lx.at(cur)
		cur ++

		lx.cur_lexeme_content_begin = cur

		for ; lx.query[cur] != terminator; cur++ {}
		
		lx.cur_lexeme_content_end = cur
		
		if lx.at(cur) == strterm {
			lx.cur_lexeme = lex_none
		} else {
			cur += 1;
			lx.cur_lexeme = lex_quoted_string;
		}
	case ':':
		if lx.at(cur+1) == ':' {
			cur += 2
			lx.cur_lexeme = lex_double_colon
		} else {
			lx.cur_lexeme = lex_none
		}
	default:
		if isCharTypeX(lx.at(cur), ctx_digit) {
			lx.cur_lexeme_content_begin = cur
			
			for ; isCharTypeX(lx.at(cur), ctx_digit); cur++ {}
			
			if lx.at(cur) == '.' {
				cur ++
				for ; isCharTypeX(lx.at(cur), ctx_digit); cur++ {}
			}

			lx.cur_lexeme_content_end = cur

			lx.cur_lexeme = lex_number
		} else if isCharTypeX(lx.at(cur), ctx_start_symbol) {
			lx.cur_lexeme_content_begin = cur;

			for ; isCharTypeX(lx.at(cur), ctx_symbol); cur++ {}
			
			if lx.at(cur) == ':'	{
				if lx.at(cur+1) == '*' { // namespace test ncname:*
					cur += 2 // :*
				} else if isCharTypeX(lx.at(cur+1), ctx_symbol) { // namespace test qname
					cur++ // :
					for ; isCharTypeX(lx.at(cur), ctx_symbol); cur++ {}
				}
			}

			lx.cur_lexeme_content_end = cur
		
			lx.cur_lexeme = lex_string
		} else {
			lx.cur_lexeme = lex_none;
		}
	}
	lx.cur = cur
}

func (lx *Lexer) current() lexeme_t {
	return lx.cur_lexeme
}
/*
		const char_t* current_pos() const
		{
			return _cur_lexeme_pos;
		}

		const xpath_lexer_string& contents() const
		{
			assert(_cur_lexeme == lex_var_ref || _cur_lexeme == lex_number || _cur_lexeme == lex_string || _cur_lexeme == lex_quoted_string);

			return _cur_lexeme_contents;
		}
*/