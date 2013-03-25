package parser

import "xpath/lexer"

type XPathParser struct {
	lexer *lexer.XPathLexer
}
struct xpath_parser
	{
		xpath_allocator* _alloc;
		xpath_lexer _lexer;

		const char_t* _query;
		xpath_variable_set* _variables;

		xpath_parse_result* _result;

		char_t _scratch[32];

	#ifdef PUGIXML_NO_EXCEPTIONS
		jmp_buf _error_handler;
	#endif

		void throw_error(const char* message)
		{
			_result->error = message;
			_result->offset = _lexer.current_pos() - _query;

		#ifdef PUGIXML_NO_EXCEPTIONS
			longjmp(_error_handler, 1);
		#else
			throw xpath_exception(*_result);
		#endif
		}

		void throw_error_oom()
		{
		#ifdef PUGIXML_NO_EXCEPTIONS
			throw_error("Out of memory");
		#else
			throw std::bad_alloc();
		#endif
		}

		void* alloc_node()
		{
			void* result = _alloc->allocate_nothrow(sizeof(xpath_ast_node));

			if (!result) throw_error_oom();

			return result;
		}

		const char_t* alloc_string(const xpath_lexer_string& value)
		{
			if (value.begin)
			{
				size_t length = static_cast<size_t>(value.end - value.begin);

				char_t* c = static_cast<char_t*>(_alloc->allocate_nothrow((length + 1) * sizeof(char_t)));
				if (!c) throw_error_oom();

				memcpy(c, value.begin, length * sizeof(char_t));
				c[length] = 0;

				return c;
			}
			else return 0;
		}

		xpath_ast_node* parse_function_helper(ast_type_t type0, ast_type_t type1, size_t argc, xpath_ast_node* args[2])
		{
			assert(argc <= 1);

			if (argc == 1 && args[0]->rettype() != xpath_type_node_set) throw_error("Function has to be applied to node set");

			return new (alloc_node()) xpath_ast_node(argc == 0 ? type0 : type1, xpath_type_string, args[0]);
		}

		xpath_ast_node* parse_function(const xpath_lexer_string& name, size_t argc, xpath_ast_node* args[2])
		{
			switch (name.begin[0])
			{
			case 'b':
				if (name == PUGIXML_TEXT("boolean") && argc == 1)
					return new (alloc_node()) xpath_ast_node(ast_func_boolean, xpath_type_boolean, args[0]);
					
				break;
			
			case 'c':
				if (name == PUGIXML_TEXT("count") && argc == 1)
				{
					if (args[0]->rettype() != xpath_type_node_set) throw_error("Function has to be applied to node set");
					return new (alloc_node()) xpath_ast_node(ast_func_count, xpath_type_number, args[0]);
				}
				else if (name == PUGIXML_TEXT("contains") && argc == 2)
					return new (alloc_node()) xpath_ast_node(ast_func_contains, xpath_type_string, args[0], args[1]);
				else if (name == PUGIXML_TEXT("concat") && argc >= 2)
					return new (alloc_node()) xpath_ast_node(ast_func_concat, xpath_type_string, args[0], args[1]);
				else if (name == PUGIXML_TEXT("ceiling") && argc == 1)
					return new (alloc_node()) xpath_ast_node(ast_func_ceiling, xpath_type_number, args[0]);
					
				break;
			
			case 'f':
				if (name == PUGIXML_TEXT("false") && argc == 0)
					return new (alloc_node()) xpath_ast_node(ast_func_false, xpath_type_boolean);
				else if (name == PUGIXML_TEXT("floor") && argc == 1)
					return new (alloc_node()) xpath_ast_node(ast_func_floor, xpath_type_number, args[0]);
					
				break;
			
			case 'i':
				if (name == PUGIXML_TEXT("id") && argc == 1)
					return new (alloc_node()) xpath_ast_node(ast_func_id, xpath_type_node_set, args[0]);
					
				break;
			
			case 'l':
				if (name == PUGIXML_TEXT("last") && argc == 0)
					return new (alloc_node()) xpath_ast_node(ast_func_last, xpath_type_number);
				else if (name == PUGIXML_TEXT("lang") && argc == 1)
					return new (alloc_node()) xpath_ast_node(ast_func_lang, xpath_type_boolean, args[0]);
				else if (name == PUGIXML_TEXT("local-name") && argc <= 1)
					return parse_function_helper(ast_func_local_name_0, ast_func_local_name_1, argc, args);
			
				break;
			
			case 'n':
				if (name == PUGIXML_TEXT("name") && argc <= 1)
					return parse_function_helper(ast_func_name_0, ast_func_name_1, argc, args);
				else if (name == PUGIXML_TEXT("namespace-uri") && argc <= 1)
					return parse_function_helper(ast_func_namespace_uri_0, ast_func_namespace_uri_1, argc, args);
				else if (name == PUGIXML_TEXT("normalize-space") && argc <= 1)
					return new (alloc_node()) xpath_ast_node(argc == 0 ? ast_func_normalize_space_0 : ast_func_normalize_space_1, xpath_type_string, args[0], args[1]);
				else if (name == PUGIXML_TEXT("not") && argc == 1)
					return new (alloc_node()) xpath_ast_node(ast_func_not, xpath_type_boolean, args[0]);
				else if (name == PUGIXML_TEXT("number") && argc <= 1)
					return new (alloc_node()) xpath_ast_node(argc == 0 ? ast_func_number_0 : ast_func_number_1, xpath_type_number, args[0]);
			
				break;
			
			case 'p':
				if (name == PUGIXML_TEXT("position") && argc == 0)
					return new (alloc_node()) xpath_ast_node(ast_func_position, xpath_type_number);
				
				break;
			
			case 'r':
				if (name == PUGIXML_TEXT("round") && argc == 1)
					return new (alloc_node()) xpath_ast_node(ast_func_round, xpath_type_number, args[0]);

				break;
			
			case 's':
				if (name == PUGIXML_TEXT("string") && argc <= 1)
					return new (alloc_node()) xpath_ast_node(argc == 0 ? ast_func_string_0 : ast_func_string_1, xpath_type_string, args[0]);
				else if (name == PUGIXML_TEXT("string-length") && argc <= 1)
					return new (alloc_node()) xpath_ast_node(argc == 0 ? ast_func_string_length_0 : ast_func_string_length_1, xpath_type_string, args[0]);
				else if (name == PUGIXML_TEXT("starts-with") && argc == 2)
					return new (alloc_node()) xpath_ast_node(ast_func_starts_with, xpath_type_boolean, args[0], args[1]);
				else if (name == PUGIXML_TEXT("substring-before") && argc == 2)
					return new (alloc_node()) xpath_ast_node(ast_func_substring_before, xpath_type_string, args[0], args[1]);
				else if (name == PUGIXML_TEXT("substring-after") && argc == 2)
					return new (alloc_node()) xpath_ast_node(ast_func_substring_after, xpath_type_string, args[0], args[1]);
				else if (name == PUGIXML_TEXT("substring") && (argc == 2 || argc == 3))
					return new (alloc_node()) xpath_ast_node(argc == 2 ? ast_func_substring_2 : ast_func_substring_3, xpath_type_string, args[0], args[1]);
				else if (name == PUGIXML_TEXT("sum") && argc == 1)
				{
					if (args[0]->rettype() != xpath_type_node_set) throw_error("Function has to be applied to node set");
					return new (alloc_node()) xpath_ast_node(ast_func_sum, xpath_type_number, args[0]);
				}

				break;
			
			case 't':
				if (name == PUGIXML_TEXT("translate") && argc == 3)
					return new (alloc_node()) xpath_ast_node(ast_func_translate, xpath_type_string, args[0], args[1]);
				else if (name == PUGIXML_TEXT("true") && argc == 0)
					return new (alloc_node()) xpath_ast_node(ast_func_true, xpath_type_boolean);
					
				break;

			default:
				break;
			}

			throw_error("Unrecognized function or wrong parameter count");

			return 0;
		}

		axis_t parse_axis_name(const xpath_lexer_string& name, bool& specified)
		{
			specified = true;

			switch (name.begin[0])
			{
			case 'a':
				if (name == PUGIXML_TEXT("ancestor"))
					return axis_ancestor;
				else if (name == PUGIXML_TEXT("ancestor-or-self"))
					return axis_ancestor_or_self;
				else if (name == PUGIXML_TEXT("attribute"))
					return axis_attribute;
				
				break;
			
			case 'c':
				if (name == PUGIXML_TEXT("child"))
					return axis_child;
				
				break;
			
			case 'd':
				if (name == PUGIXML_TEXT("descendant"))
					return axis_descendant;
				else if (name == PUGIXML_TEXT("descendant-or-self"))
					return axis_descendant_or_self;
				
				break;
			
			case 'f':
				if (name == PUGIXML_TEXT("following"))
					return axis_following;
				else if (name == PUGIXML_TEXT("following-sibling"))
					return axis_following_sibling;
				
				break;
			
			case 'n':
				if (name == PUGIXML_TEXT("namespace"))
					return axis_namespace;
				
				break;
			
			case 'p':
				if (name == PUGIXML_TEXT("parent"))
					return axis_parent;
				else if (name == PUGIXML_TEXT("preceding"))
					return axis_preceding;
				else if (name == PUGIXML_TEXT("preceding-sibling"))
					return axis_preceding_sibling;
				
				break;
			
			case 's':
				if (name == PUGIXML_TEXT("self"))
					return axis_self;
				
				break;

			default:
				break;
			}

			specified = false;
			return axis_child;
		}

		nodetest_t parse_node_test_type(const xpath_lexer_string& name)
		{
			switch (name.begin[0])
			{
			case 'c':
				if (name == PUGIXML_TEXT("comment"))
					return nodetest_type_comment;

				break;

			case 'n':
				if (name == PUGIXML_TEXT("node"))
					return nodetest_type_node;

				break;

			case 'p':
				if (name == PUGIXML_TEXT("processing-instruction"))
					return nodetest_type_pi;

				break;

			case 't':
				if (name == PUGIXML_TEXT("text"))
					return nodetest_type_text;

				break;
			
			default:
				break;
			}

			return nodetest_none;
		}

		// PrimaryExpr ::= VariableReference | '(' Expr ')' | Literal | Number | FunctionCall
		xpath_ast_node* parse_primary_expression()
		{
			switch (_lexer.current())
			{
			case lex_var_ref:
			{
				xpath_lexer_string name = _lexer.contents();

				if (!_variables)
					throw_error("Unknown variable: variable set is not provided");

				xpath_variable* var = get_variable_scratch(_scratch, _variables, name.begin, name.end);

				if (!var)
					throw_error("Unknown variable: variable set does not contain the given name");

				_lexer.next();

				return new (alloc_node()) xpath_ast_node(ast_variable, var->type(), var);
			}

			case lex_open_brace:
			{
				_lexer.next();

				xpath_ast_node* n = parse_expression();

				if (_lexer.current() != lex_close_brace)
					throw_error("Unmatched braces");

				_lexer.next();

				return n;
			}

			case lex_quoted_string:
			{
				const char_t* value = alloc_string(_lexer.contents());

				xpath_ast_node* n = new (alloc_node()) xpath_ast_node(ast_string_constant, xpath_type_string, value);
				_lexer.next();

				return n;
			}

			case lex_number:
			{
				double value = 0;

				if (!convert_string_to_number_scratch(_scratch, _lexer.contents().begin, _lexer.contents().end, &value))
					throw_error_oom();

				xpath_ast_node* n = new (alloc_node()) xpath_ast_node(ast_number_constant, xpath_type_number, value);
				_lexer.next();

				return n;
			}

			case lex_string:
			{
				xpath_ast_node* args[2] = {0};
				size_t argc = 0;
				
				xpath_lexer_string function = _lexer.contents();
				_lexer.next();
				
				xpath_ast_node* last_arg = 0;
				
				if (_lexer.current() != lex_open_brace)
					throw_error("Unrecognized function call");
				_lexer.next();

				if (_lexer.current() != lex_close_brace)
					args[argc++] = parse_expression();

				while (_lexer.current() != lex_close_brace)
				{
					if (_lexer.current() != lex_comma)
						throw_error("No comma between function arguments");
					_lexer.next();
					
					xpath_ast_node* n = parse_expression();
					
					if (argc < 2) args[argc] = n;
					else last_arg->set_next(n);

					argc++;
					last_arg = n;
				}
				
				_lexer.next();

				return parse_function(function, argc, args);
			}

			default:
				throw_error("Unrecognizable primary expression");

				return 0;
			}
		}
		
		// FilterExpr ::= PrimaryExpr | FilterExpr Predicate
		// Predicate ::= '[' PredicateExpr ']'
		// PredicateExpr ::= Expr
		xpath_ast_node* parse_filter_expression()
		{
			xpath_ast_node* n = parse_primary_expression();

			while (_lexer.current() == lex_open_square_brace)
			{
				_lexer.next();

				xpath_ast_node* expr = parse_expression();

				if (n->rettype() != xpath_type_node_set) throw_error("Predicate has to be applied to node set");

				bool posinv = expr->rettype() != xpath_type_number && expr->is_posinv();

				n = new (alloc_node()) xpath_ast_node(posinv ? ast_filter_posinv : ast_filter, xpath_type_node_set, n, expr);

				if (_lexer.current() != lex_close_square_brace)
					throw_error("Unmatched square brace");
			
				_lexer.next();
			}
			
			return n;
		}
		
		// Step ::= AxisSpecifier NodeTest Predicate* | AbbreviatedStep
		// AxisSpecifier ::= AxisName '::' | '@'?
		// NodeTest ::= NameTest | NodeType '(' ')' | 'processing-instruction' '(' Literal ')'
		// NameTest ::= '*' | NCName ':' '*' | QName
		// AbbreviatedStep ::= '.' | '..'
		xpath_ast_node* parse_step(xpath_ast_node* set)
		{
			if (set && set->rettype() != xpath_type_node_set)
				throw_error("Step has to be applied to node set");

			bool axis_specified = false;
			axis_t axis = axis_child; // implied child axis

			if (_lexer.current() == lex_axis_attribute)
			{
				axis = axis_attribute;
				axis_specified = true;
				
				_lexer.next();
			}
			else if (_lexer.current() == lex_dot)
			{
				_lexer.next();
				
				return new (alloc_node()) xpath_ast_node(ast_step, set, axis_self, nodetest_type_node, 0);
			}
			else if (_lexer.current() == lex_double_dot)
			{
				_lexer.next();
				
				return new (alloc_node()) xpath_ast_node(ast_step, set, axis_parent, nodetest_type_node, 0);
			}
		
			nodetest_t nt_type = nodetest_none;
			xpath_lexer_string nt_name;
			
			if (_lexer.current() == lex_string)
			{
				// node name test
				nt_name = _lexer.contents();
				_lexer.next();

				// was it an axis name?
				if (_lexer.current() == lex_double_colon)
				{
					// parse axis name
					if (axis_specified) throw_error("Two axis specifiers in one step");

					axis = parse_axis_name(nt_name, axis_specified);

					if (!axis_specified) throw_error("Unknown axis");

					// read actual node test
					_lexer.next();

					if (_lexer.current() == lex_multiply)
					{
						nt_type = nodetest_all;
						nt_name = xpath_lexer_string();
						_lexer.next();
					}
					else if (_lexer.current() == lex_string)
					{
						nt_name = _lexer.contents();
						_lexer.next();
					}
					else throw_error("Unrecognized node test");
				}
				
				if (nt_type == nodetest_none)
				{
					// node type test or processing-instruction
					if (_lexer.current() == lex_open_brace)
					{
						_lexer.next();
						
						if (_lexer.current() == lex_close_brace)
						{
							_lexer.next();

							nt_type = parse_node_test_type(nt_name);

							if (nt_type == nodetest_none) throw_error("Unrecognized node type");
							
							nt_name = xpath_lexer_string();
						}
						else if (nt_name == PUGIXML_TEXT("processing-instruction"))
						{
							if (_lexer.current() != lex_quoted_string)
								throw_error("Only literals are allowed as arguments to processing-instruction()");
						
							nt_type = nodetest_pi;
							nt_name = _lexer.contents();
							_lexer.next();
							
							if (_lexer.current() != lex_close_brace)
								throw_error("Unmatched brace near processing-instruction()");
							_lexer.next();
						}
						else
							throw_error("Unmatched brace near node type test");

					}
					// QName or NCName:*
					else
					{
						if (nt_name.end - nt_name.begin > 2 && nt_name.end[-2] == ':' && nt_name.end[-1] == '*') // NCName:*
						{
							nt_name.end--; // erase *
							
							nt_type = nodetest_all_in_namespace;
						}
						else nt_type = nodetest_name;
					}
				}
			}
			else if (_lexer.current() == lex_multiply)
			{
				nt_type = nodetest_all;
				_lexer.next();
			}
			else throw_error("Unrecognized node test");
			
			xpath_ast_node* n = new (alloc_node()) xpath_ast_node(ast_step, set, axis, nt_type, alloc_string(nt_name));
			
			xpath_ast_node* last = 0;
			
			while (_lexer.current() == lex_open_square_brace)
			{
				_lexer.next();
				
				xpath_ast_node* expr = parse_expression();

				xpath_ast_node* pred = new (alloc_node()) xpath_ast_node(ast_predicate, xpath_type_node_set, expr);
				
				if (_lexer.current() != lex_close_square_brace)
					throw_error("Unmatched square brace");
				_lexer.next();
				
				if (last) last->set_next(pred);
				else n->set_right(pred);
				
				last = pred;
			}
			
			return n;
		}
		
		// RelativeLocationPath ::= Step | RelativeLocationPath '/' Step | RelativeLocationPath '//' Step
		xpath_ast_node* parse_relative_location_path(xpath_ast_node* set)
		{
			xpath_ast_node* n = parse_step(set);
			
			while (_lexer.current() == lex_slash || _lexer.current() == lex_double_slash)
			{
				lexeme_t l = _lexer.current();
				_lexer.next();

				if (l == lex_double_slash)
					n = new (alloc_node()) xpath_ast_node(ast_step, n, axis_descendant_or_self, nodetest_type_node, 0);
				
				n = parse_step(n);
			}
			
			return n;
		}
		
		// LocationPath ::= RelativeLocationPath | AbsoluteLocationPath
		// AbsoluteLocationPath ::= '/' RelativeLocationPath? | '//' RelativeLocationPath
		xpath_ast_node* parse_location_path()
		{
			if (_lexer.current() == lex_slash)
			{
				_lexer.next();
				
				xpath_ast_node* n = new (alloc_node()) xpath_ast_node(ast_step_root, xpath_type_node_set);

				// relative location path can start from axis_attribute, dot, double_dot, multiply and string lexemes; any other lexeme means standalone root path
				lexeme_t l = _lexer.current();

				if (l == lex_string || l == lex_axis_attribute || l == lex_dot || l == lex_double_dot || l == lex_multiply)
					return parse_relative_location_path(n);
				else
					return n;
			}
			else if (_lexer.current() == lex_double_slash)
			{
				_lexer.next();
				
				xpath_ast_node* n = new (alloc_node()) xpath_ast_node(ast_step_root, xpath_type_node_set);
				n = new (alloc_node()) xpath_ast_node(ast_step, n, axis_descendant_or_self, nodetest_type_node, 0);
				
				return parse_relative_location_path(n);
			}

			// else clause moved outside of if because of bogus warning 'control may reach end of non-void function being inlined' in gcc 4.0.1
			return parse_relative_location_path(0);
		}
		
		// PathExpr ::= LocationPath
		//				| FilterExpr
		//				| FilterExpr '/' RelativeLocationPath
		//				| FilterExpr '//' RelativeLocationPath
		// UnionExpr ::= PathExpr | UnionExpr '|' PathExpr
		// UnaryExpr ::= UnionExpr | '-' UnaryExpr
		xpath_ast_node* parse_path_or_unary_expression()
		{
			// Clarification.
			// PathExpr begins with either LocationPath or FilterExpr.
			// FilterExpr begins with PrimaryExpr
			// PrimaryExpr begins with '$' in case of it being a variable reference,
			// '(' in case of it being an expression, string literal, number constant or
			// function call.

			if (_lexer.current() == lex_var_ref || _lexer.current() == lex_open_brace || 
				_lexer.current() == lex_quoted_string || _lexer.current() == lex_number ||
				_lexer.current() == lex_string)
			{
				if (_lexer.current() == lex_string)
				{
					// This is either a function call, or not - if not, we shall proceed with location path
					const char_t* state = _lexer.state();
					
					while (PUGI__IS_CHARTYPE(*state, ct_space)) ++state;
					
					if (*state != '(') return parse_location_path();

					// This looks like a function call; however this still can be a node-test. Check it.
					if (parse_node_test_type(_lexer.contents()) != nodetest_none) return parse_location_path();
				}
				
				xpath_ast_node* n = parse_filter_expression();

				if (_lexer.current() == lex_slash || _lexer.current() == lex_double_slash)
				{
					lexeme_t l = _lexer.current();
					_lexer.next();
					
					if (l == lex_double_slash)
					{
						if (n->rettype() != xpath_type_node_set) throw_error("Step has to be applied to node set");

						n = new (alloc_node()) xpath_ast_node(ast_step, n, axis_descendant_or_self, nodetest_type_node, 0);
					}
	
					// select from location path
					return parse_relative_location_path(n);
				}

				return n;
			}
			else if (_lexer.current() == lex_minus)
			{
				_lexer.next();

				// precedence 7+ - only parses union expressions
				xpath_ast_node* expr = parse_expression_rec(parse_path_or_unary_expression(), 7);

				return new (alloc_node()) xpath_ast_node(ast_op_negate, xpath_type_number, expr);
			}
			else
				return parse_location_path();
		}

		struct binary_op_t
		{
			ast_type_t asttype;
			xpath_value_type rettype;
			int precedence;

			binary_op_t(): asttype(ast_unknown), rettype(xpath_type_none), precedence(0)
			{
			}

			binary_op_t(ast_type_t asttype_, xpath_value_type rettype_, int precedence_): asttype(asttype_), rettype(rettype_), precedence(precedence_)
			{
			}

			static binary_op_t parse(xpath_lexer& lexer)
			{
				switch (lexer.current())
				{
				case lex_string:
					if (lexer.contents() == PUGIXML_TEXT("or"))
						return binary_op_t(ast_op_or, xpath_type_boolean, 1);
					else if (lexer.contents() == PUGIXML_TEXT("and"))
						return binary_op_t(ast_op_and, xpath_type_boolean, 2);
					else if (lexer.contents() == PUGIXML_TEXT("div"))
						return binary_op_t(ast_op_divide, xpath_type_number, 6);
					else if (lexer.contents() == PUGIXML_TEXT("mod"))
						return binary_op_t(ast_op_mod, xpath_type_number, 6);
					else
						return binary_op_t();

				case lex_equal:
					return binary_op_t(ast_op_equal, xpath_type_boolean, 3);

				case lex_not_equal:
					return binary_op_t(ast_op_not_equal, xpath_type_boolean, 3);

				case lex_less:
					return binary_op_t(ast_op_less, xpath_type_boolean, 4);

				case lex_greater:
					return binary_op_t(ast_op_greater, xpath_type_boolean, 4);

				case lex_less_or_equal:
					return binary_op_t(ast_op_less_or_equal, xpath_type_boolean, 4);

				case lex_greater_or_equal:
					return binary_op_t(ast_op_greater_or_equal, xpath_type_boolean, 4);

				case lex_plus:
					return binary_op_t(ast_op_add, xpath_type_number, 5);

				case lex_minus:
					return binary_op_t(ast_op_subtract, xpath_type_number, 5);

				case lex_multiply:
					return binary_op_t(ast_op_multiply, xpath_type_number, 6);

				case lex_union:
					return binary_op_t(ast_op_union, xpath_type_node_set, 7);

				default:
					return binary_op_t();
				}
			}
		};

		xpath_ast_node* parse_expression_rec(xpath_ast_node* lhs, int limit)
		{
			binary_op_t op = binary_op_t::parse(_lexer);

			while (op.asttype != ast_unknown && op.precedence >= limit)
			{
				_lexer.next();

				xpath_ast_node* rhs = parse_path_or_unary_expression();

				binary_op_t nextop = binary_op_t::parse(_lexer);

				while (nextop.asttype != ast_unknown && nextop.precedence > op.precedence)
				{
					rhs = parse_expression_rec(rhs, nextop.precedence);

					nextop = binary_op_t::parse(_lexer);
				}

				if (op.asttype == ast_op_union && (lhs->rettype() != xpath_type_node_set || rhs->rettype() != xpath_type_node_set))
					throw_error("Union operator has to be applied to node sets");

				lhs = new (alloc_node()) xpath_ast_node(op.asttype, op.rettype, lhs, rhs);

				op = binary_op_t::parse(_lexer);
			}

			return lhs;
		}

		// Expr ::= OrExpr
		// OrExpr ::= AndExpr | OrExpr 'or' AndExpr
		// AndExpr ::= EqualityExpr | AndExpr 'and' EqualityExpr
		// EqualityExpr ::= RelationalExpr
		//					| EqualityExpr '=' RelationalExpr
		//					| EqualityExpr '!=' RelationalExpr
		// RelationalExpr ::= AdditiveExpr
		//					  | RelationalExpr '<' AdditiveExpr
		//					  | RelationalExpr '>' AdditiveExpr
		//					  | RelationalExpr '<=' AdditiveExpr
		//					  | RelationalExpr '>=' AdditiveExpr
		// AdditiveExpr ::= MultiplicativeExpr
		//					| AdditiveExpr '+' MultiplicativeExpr
		//					| AdditiveExpr '-' MultiplicativeExpr
		// MultiplicativeExpr ::= UnaryExpr
		//						  | MultiplicativeExpr '*' UnaryExpr
		//						  | MultiplicativeExpr 'div' UnaryExpr
		//						  | MultiplicativeExpr 'mod' UnaryExpr
		xpath_ast_node* parse_expression()
		{
			return parse_expression_rec(parse_path_or_unary_expression(), 0);
		}

		xpath_parser(const char_t* query, xpath_variable_set* variables, xpath_allocator* alloc, xpath_parse_result* result): _alloc(alloc), _lexer(query), _query(query), _variables(variables), _result(result)
		{
		}

		xpath_ast_node* parse()
		{
			xpath_ast_node* result = parse_expression();
			
			if (_lexer.current() != lex_eof)
			{
				// there are still unparsed tokens left, error
				throw_error("Incorrect query");
			}
			
			return result;
		}

		static xpath_ast_node* parse(const char_t* query, xpath_variable_set* variables, xpath_allocator* alloc, xpath_parse_result* result)
		{
			xpath_parser parser(query, variables, alloc, result);

		#ifdef PUGIXML_NO_EXCEPTIONS
			int error = setjmp(parser._error_handler);

			return (error == 0) ? parser.parse() : 0;
		#else
			return parser.parse();
		#endif
		}
	};