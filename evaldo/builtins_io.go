// +build !b_no_io

package evaldo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"rye/env"
	"strings"

	//	"crypto/tls"

	"net/http"
	//"net/smtp"
	//	gomail "gopkg.in/mail.v2"
	// "net/url"
	//"strconv"
	"net/http/cgi"
)

func __input(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch str := arg0.(type) {
	case env.String:
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(str)
		inp, _ := reader.ReadString('\n')
		fmt.Println(inp)
		return env.String{inp}
	default:
		//env1.ReturnFlag = true
		env1.FailureFlag = true
		return env.NewError("arg 1 should be string")
	}
}

func __open(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch s := arg0.(type) {
	case env.Uri:
		path := strings.Split(s.Path, "://")
		file, err := os.Open(path[1])
		if err != nil {
			//env1.ReturnFlag = true
			env1.FailureFlag = true
			return *env.NewError(err.Error())
		}
		return *env.NewNative(env1.Idx, file, "rye-file")
	default:
		//env1.ReturnFlag = true
		env1.FailureFlag = true
		return *env.NewError("just accepting Uri-s")
	}
}

func __create(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch s := arg0.(type) {
	case env.Uri:
		path := strings.Split(s.Path, "://")
		file, err := os.Create(path[1])
		if err != nil {
			env1.ReturnFlag = true
			env1.FailureFlag = true
			return *env.NewError(err.Error())
		}
		return *env.NewNative(env1.Idx, file, "rye-file")
	default:
		env1.ReturnFlag = true
		env1.FailureFlag = true
		return *env.NewError("just accepting Uri-s")
	}
}

func __open_reader(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch s := arg0.(type) {
	case env.Uri:
		path := strings.Split(s.Path, "://")
		file, err := os.Open(path[1])
		//trace3(path)
		if err != nil {
			env1.FailureFlag = true
			return *env.NewError("Error opening file")
		}
		return *env.NewNative(env1.Idx, bufio.NewReader(file), "rye-reader")
	case env.String:
		return *env.NewNative(env1.Idx, strings.NewReader(s.Value), "rye-reader")
	default:
		env1.FailureFlag = true
		return *env.NewError("just accepting Uri-s")
	}
}

func __read_all(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch s := arg0.(type) {
	case env.Native:
		data, err := ioutil.ReadAll(s.Value.(io.Reader))
		if err != nil {
			env1.FailureFlag = true
			return *env.NewError("Error reading file")
		}
		return env.String{string(data)}
	}
	env1.FailureFlag = true
	return *env.NewError("Failed")
}

func __close(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch s := arg0.(type) {
	case env.Native:
		err := s.Value.(*os.File).Close()
		if err != nil {
			env1.FailureFlag = true
			return *env.NewError(err.Error())
		}
		return env.String{""}
	}
	env1.FailureFlag = true
	return *env.NewError("Failed")
}

func __write(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch f := arg0.(type) {
	case env.Native:
		switch s := arg1.(type) {
		case env.String:

			bytesWritten, err := f.Value.(io.Writer).Write([]byte(s.Value))
			if err != nil {
				env1.FailureFlag = true
				return *env.NewError(err.Error())
			}
			return env.Integer{int64(bytesWritten)}
			//log.Printf("Wrote %d bytes.\n", bytesWritten)
		}
	}
	env1.FailureFlag = true
	return *env.NewError("Failed")
}

func __fs_read(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch f := arg0.(type) {
	case env.Uri:

		data, err := ioutil.ReadFile(f.GetPath())
		if err != nil {
			env1.FailureFlag = true
			return *env.NewError(err.Error())
		}
		return env.String{string(data)}
	}
	env1.FailureFlag = true
	return *env.NewError("Failed")

	// Read file to byte slice
}

func __fs_write(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch f := arg0.(type) {
	case env.Uri:
		switch s := arg1.(type) {
		case env.String:

			err := ioutil.WriteFile(f.GetPath(), []byte(s.Value), 0644)
			if err != nil {
				env1.FailureFlag = true
				return *env.NewError(err.Error())
			}
			return arg1
		}
	}
	env1.FailureFlag = true
	return *env.NewError("Failed")
}

func __https_s_get(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch f := arg0.(type) {
	case env.Uri:

		resp, err := http.Get(f.GetProtocol() + "://" + f.GetPath())
		if err != nil {
			env1.FailureFlag = true
			return *env.NewError(err.Error())
		}

		// Print the HTTP Status Code and Status Name
		//mt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			return env.String{string(body)}
		} else {
			env1.FailureFlag = true
			return *env.NewError2(resp.StatusCode, string(body))
		}

		// log.Printf("Data read: %s\n", data)

	}
	env1.FailureFlag = true
	return *env.NewError("Failed")

	// Read file to byte slice
}

func __http_s_post(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch f := arg0.(type) {
	case env.Uri:

		switch t := arg2.(type) {
		case env.Tagword:
			switch d := arg1.(type) {
			case env.String:

				var tt string
				tidx, terr := env1.Idx.GetIndex("json")
				tidx2, terr2 := env1.Idx.GetIndex("text")
				if terr && t.Index == tidx {
					//if t.Value == "json" {
					tt = "application/json"
				} else if terr2 && t.Index == tidx2 {
					tt = "text/plain"
				} else {
					env1.FailureFlag = true
					return *env.NewError("wrong content type")
				}
				// TODO -- add other cases
				// fmt.Println("BEFORE")

				resp, err := http.Post(f.GetProtocol()+"://"+f.GetPath(), tt, bytes.NewBufferString(d.Value))
				if err != nil {
					// fmt.Println("ERR")
					env1.FailureFlag = true
					return *env.NewError(err.Error())
				}

				// Print the HTTP Status Code and Status Name
				// fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					// fmt.Println("ERR")
					env1.FailureFlag = true
					return *env.NewError(err.Error())
				}

				if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
					return env.String{string(body)}
				} else {
					// fmt.Println("ERR33")
					env1.FailureFlag = true
					return env.NewError2(resp.StatusCode, string(body))
				}
			}
		}
	}
	// fmt.Println("ERR44")

	env1.FailureFlag = true
	return *env.NewError("Failed 123")

	// Read file to byte slice
}

func __email_send(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch to_ := arg0.(type) {
	case env.Email:

		switch msg := arg1.(type) {
		case env.String:

			idx, _ := env1.Idx.GetIndex("user-profile")
			uctx_, _ := env1.Ctx.Get(idx)
			uctx := uctx_.(env.RyeCtx)
			fmt.Println(to_)
			fmt.Println(msg)
			fmt.Println(uctx)
			// TODO continue: uncomment and make it work
			/*
				from, _ := uctx.Get(env1.Idx.GetIndex("smtp-from"))
				password, _ := uctx.Get(env1.Idx.GetIndex("smtp-password"))
				server, _ := uctx.Get(env1.Idx.GetIndex("smtp-server"))
				port, _ := uctx.Get(env1.Idx.GetIndex("smtp-port"))
				// Receiver email address.
				// to := []string{
				//	to_.Value,
				//}
				// Message.
				// message := []byte(msg.Value)
				m := gomail.NewMessage()

				// Set E-Mail sender
				m.SetHeader("From", from)

				// Set E-Mail receivers
				m.SetHeader("To", to_.Address)

				// Set E-Mail subject
				m.SetHeader("Subject", msg.Value)

				// Set E-Mail body. You can set plain text or html with text/html
				m.SetBody("text/plain", msg.Value)

				// Settings for SMTP server
				d := gomail.NewDialer(server, port, from, password)

				// This is only needed when SSL/TLS certificate is not valid on server.
				// In production this should be set to false.
				//			d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

				// Now send E-Mail
				if err := d.DialAndSend(m); err != nil {
					env1.FailureFlag = true
					return env.NewError(err.Error())
				}
			*/
			return env.Integer{1}
		}
	}
	env1.FailureFlag = true
	return *env.NewError("Failed")

	// Read file to byte slice
}

func __https_s__new_request(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch uri := arg0.(type) {
	case env.Uri:
		switch method := arg1.(type) {
		case env.Tagword:
			method1 := env1.Idx.GetWord(method.Index)
			if !(method1 == "GET" || method1 == "POST") {
				env1.FailureFlag = true
				return *env.NewError("Wrong method")
			}
			switch data := arg2.(type) {
			case env.String:
				data1 := strings.NewReader(data.Value)
				req, err := http.NewRequest(method1, uri.GetProtocol()+"://"+uri.GetPath(), data1)
				if err != nil {
					env1.FailureFlag = true
					return *env.NewError("Error: " + err.Error())
				}
				return *env.NewNative(env1.Idx, req, "https-request")
			default:
				env1.FailureFlag = true
				return *env.NewError("Arg 3 not String")
			}
		default:
			env1.FailureFlag = true
			return *env.NewError("Arg 2 not Word")
		}
	default:
		env1.FailureFlag = true
		return *env.NewError("Arg 1 not Uri")
	}
}

func __https_request__set_header(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch req := arg0.(type) {
	case env.Native:
		switch method := arg1.(type) {
		case env.Tagword:
			name := env1.Idx.GetWord(method.Index)
			switch data := arg2.(type) {
			case env.String:
				req.Value.(*http.Request).Header.Set(name, data.Value)
				return arg0
			default:
				return makeError(env1, "Arg 3 not String")
			}
		default:
			return makeError(env1, "Arg 2 not Word")
		}
	default:
		return makeError(env1, "Arg 1 not Native")
	}
	return makeError(env1, "Failed")
}

func makeError(env1 *env.ProgramState, msg string) *env.Error {
	env1.FailureFlag = true
	return env.NewError(msg)
}

func __https_request__do(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch req := arg0.(type) {
	case env.Native:
		client := &http.Client{}
		resp, err := client.Do(req.Value.(*http.Request))
		if err != nil {
			return makeError(env1, "Error: "+err.Error())
		}
		return *env.NewNative(env1.Idx, resp, "https-response")
	default:
		return makeError(env1, "Arg 1 not Native")
	}
	return makeError(env1, "Failed")
}

func __https_response__read_body(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
	switch resp := arg0.(type) {
	case env.Native:
		data, err := ioutil.ReadAll(resp.Value.(*http.Response).Body)
		if err != nil {
			makeError(env1, "Error: "+err.Error())
		}
		return env.String{string(data)}
	default:
		env1.FailureFlag = true
		return *env.NewError("Arg 1 not Native")
	}
	env1.FailureFlag = true
	return *env.NewError("Failed")
}

var Builtins_io = map[string]*env.Builtin{

	"input": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __input(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"file-schema//open": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __open(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"file-schema//create": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __create(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	// should this be generic method or not?
	"reader": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __open_reader(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"rye-file//read-all": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __read_all(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"rye-file//write": {
		Argsn: 2,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __write(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"rye-file//close": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __close(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"file-schema//read": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __fs_read(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"file-schema//write": {
		Argsn: 2,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __fs_write(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"https-schema//get": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __https_s_get(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"https-schema//post": {
		Argsn: 3,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __http_s_post(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"http-schema//get": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __https_s_get(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"http-schema//post": {
		Argsn: 3,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __http_s_post(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"https-schema//new-request": {
		Argsn: 3,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __https_s__new_request(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"https-request//set-header": {
		Argsn: 3,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __https_request__set_header(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"https-request//call": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __https_request__do(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"https-response//read-body": {
		Argsn: 1,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __https_response__read_body(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"email//send": {
		Argsn: 2,
		Fn: func(env1 *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			return __email_send(env1, arg0, arg1, arg2, arg3, arg4)
		},
	},

	"serve-cgi": {
		Argsn: 3,
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch rword := arg0.(type) {
			case env.Tagword:
				switch wword := arg1.(type) {
				case env.Tagword:
					switch bloc := arg2.(type) {
					case env.Block:
						var rctx *env.RyeCtx
						if err := cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							ser := ps.Ser
							ctx := ps.Ctx
							ps.Ser = bloc.Series
							ps.Ctx = env.NewEnv(ps.Ctx) // make new context with no parent
							ps.Ctx.Set(rword.Index, *env.NewNative(ps.Idx, w, "Go-server-response-writer"))
							ps.Ctx.Set(wword.Index, *env.NewNative(ps.Idx, r, "Go-server-request"))
							EvalBlock(ps)
							rctx = ps.Ctx
							ps.Ctx = ctx
							ps.Ser = ser
						})); err != nil {
							return makeError(ps, "Error: "+err.Error())
						}
						return *rctx
					default:
						return makeError(ps, "Arg 3 not Block")
					}
				default:
					return makeError(ps, "Arg 2 not Tagword")
				}
			default:
				return makeError(ps, "Arg 1 not Tagword")
			}
		},
	},
}
