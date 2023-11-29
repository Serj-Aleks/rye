//go:build b_nng
// +build b_nng

package evaldo

import "C"

import (
	//	"fmt"
	"rye/env"

	"go.nanomsg.org/mangos"

	// register transports
	"go.nanomsg.org/mangos/protocol/rep"
	"go.nanomsg.org/mangos/protocol/req"
	_ "go.nanomsg.org/mangos/transport/all"
)

/*

Basic example req/rep

sock: open nng://rep |^check "can't get new socket"
listen sock tcp://127.0.0.1:40404 |^check "can't listen on rep socket"
forever {
		msg: read sock |^ "cannot receive on rep socket"

		if to-string msg == "DATE" {
			print "NODE0: RECEIVED DATE REQUEST"
			send sock to-byte now/date |^check "can't send reply"
		}
}


open nng://rep |^check "can't get new socket" :sock
  |listen tcp://127.0.0.1:40404 |^check "can't listen on rep socket"

forever {
	receive sock |^ "cannot receive on rep socket"
	  |to-string == "DATE"
	  |if {
		  print "NODE0: RECEIVED DATE REQUEST"
		  now/date
		    |to-bytes
		    |send sock |^check "can't send reply"
	}
}


if sock, err = rep.NewSocket(); err != nil {
		die("can't get new rep socket: %s", err)
	}
	if err = sock.Listen(url); err != nil {
		die("can't listen on rep socket: %s", err.Error())
	}
	for {
		// Could also use sock.RecvMsg to get header
		msg, err = sock.Recv()
		if err != nil {
			die("cannot receive on rep socket: %s", err.Error())
		}
		if string(msg) == "DATE" { // no need to terminate
			fmt.Println("NODE0: RECEIVED DATE REQUEST")
			d := date()
			fmt.Printf("NODE0: SENDING DATE %s\n", d)
			err = sock.Send([]byte(d))
			if err != nil {
				die("can't send reply: %s", err.Error())
			}
		}
	}
	}


*/

var Builtins_nng = map[string]*env.Builtin{

	"nng-schema//open": {
		Argsn: 1,
		Doc:   "TODODOC",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch uri := arg0.(type) {
			case env.Uri:
				// TODO -- switch over socket type nng://rep req ...
				var sock mangos.Socket
				var err error
				switch uri.GetPath() {
				case "rep":
					if sock, err = rep.NewSocket(); err != nil {
						ps.FailureFlag = true
						return MakeBuiltinError(ps, err.Error(), "nng-schema//open")
					}
				case "req":
					if sock, err = req.NewSocket(); err != nil {
						ps.FailureFlag = true
						return MakeBuiltinError(ps, err.Error(), "nng-schema//open")
					}
				default:
					return MakeBuiltinError(ps, "No matching path found.", "nng-schema//open")
				}
				return *env.NewNative(ps.Idx, sock, "Nng-socket")
			default:
				ps.FailureFlag = true
				return MakeArgError(ps, 1, []env.Type{env.UriType}, "nng-schema//open")
			}
		},
	},

	"Nng-socket//listen": {
		Argsn: 2,
		Doc:   "TODODOC",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch sock := arg0.(type) {
			case env.Native:
				switch url := arg1.(type) {
				case env.Uri:
					var err error
					if err = sock.Value.(mangos.Socket).Listen(url.Path); err != nil {
						ps.FailureFlag = true
						return *env.NewError(err.Error())
					}
					return arg0
				default:
					ps.FailureFlag = true
					return MakeArgError(ps, 2, []env.Type{env.UriType}, "Nng-socket//listen")
				}
			default:
				ps.FailureFlag = true
				return MakeArgError(ps, 1, []env.Type{env.NativeType}, "Nng-socket//listen")
			}
		},
	},
	"Nng-socket//dial": {
		Argsn: 2,
		Doc:   "TODODOC",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch sock := arg0.(type) {
			case env.Native:
				switch url := arg1.(type) {
				case env.Uri:
					var err error
					if err = sock.Value.(mangos.Socket).Dial(url.Path); err != nil {
						ps.FailureFlag = true
						return MakeBuiltinError(ps, err.Error(), "Nng-socket//dial")
					}
					return arg0
				default:
					ps.FailureFlag = true
					return MakeArgError(ps, 2, []env.Type{env.UriType}, "Nng-socket//dial")
				}
			default:
				ps.FailureFlag = true
				return MakeArgError(ps, 1, []env.Type{env.NativeType}, "Nng-socket//dial")
			}
		},
	},
	"Nng-socket//receive": {
		Argsn: 1,
		Doc:   "TODODOC",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch sock := arg0.(type) {
			case env.Native:
				msg, err := sock.Value.(mangos.Socket).Recv()
				if err != nil {
					ps.FailureFlag = true
					return MakeBuiltinError(ps, err.Error(), "Nng-socket//receive")
				}
				return env.String{string(msg)}
			default:
				ps.FailureFlag = true
				return MakeArgError(ps, 1, []env.Type{env.NativeType}, "Nng-socket//receive")
			}
		},
	},
	"Nng-socket//send": {
		Argsn: 2,
		Doc:   "TODODOC",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch sock := arg0.(type) {
			case env.Native:
				switch d := arg1.(type) {
				case env.String:
					err := sock.Value.(mangos.Socket).Send([]byte(d.Value))
					if err != nil {
						ps.FailureFlag = true
						return MakeBuiltinError(ps, err.Error(), "Nng-socket//send")
					}
					return env.String{string(d.Value)}
				default:
					ps.FailureFlag = true
					return MakeArgError(ps, 2, []env.Type{env.StringType}, "Nng-socket//send")
				}
			default:
				ps.FailureFlag = true
				return MakeArgError(ps, 1, []env.Type{env.NativeType}, "Nng-socket//send")
			}
		},
	},
}
