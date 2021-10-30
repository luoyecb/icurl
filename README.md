# icurl
Interactive curl tool

# install
```sh
go get github.com/luoyecb/icurl
make
./icurl
```

# run
```
  _____     ____   __    __   ______     _____
 (_   _)   / ___)  ) )  ( (  (   __ \   (_   _)
   | |    / /     ( (    ) )  ) (__) )    | |
   | |   ( (       ) )  ( (  (    __/     | |
   | |   ( (      ( (    ) )  ) \ \  _    | |   __
  _| |__  \ \___   ) \__/ (  ( ( \ \_)) __| |___) )
 /_____(   \____)  \______/   )_) \__/  \________/

You can get help information through the help() function.
icurl>
```

# help
```
icurl> help()
=== context
context = {
	scheme = "http", # http|https
	host   = "localhost", # must string
	port   = 80,     # must number
	path   = "",     # must string
	method = "GET",  # GET|PUT|POST|DELETE
	url    = "",     # must string, if set, other fields are ignored
	data   = "",     # must string, if method=GET, data are ignored, otherwise use data instead of query
	query  = {},     # must table
	header = {},     # must table
}

=== functions
exit|quit                 : exit
reset()                   : reset context
loadf(string)             : load lua file, absolute path
loadh(string)             : load lua file, default in dir ~/.icurl/
listh(string)             : list lua file, default in dir ~/.icurl/
saveh(string, [bool])     : save lua file, default in dir ~/.icurl/, bool arg means whether overwrite existing file or not
debugc()                  : print context information
send([bool])              : send http requeset, http method is context.method, bool arg means json pretty formatting
send_get([bool])          : send http get requeset, bool arg means json pretty formatting
send_post([bool])         : send http post requeset, bool arg means json pretty formatting
send_form([bool])         : send http post requeset, with header "Content-Type:application/x-www-form-urlencoded", bool arg means json pretty formatting
set_query(string, string) : set context.query
set_header(string, string): set context.header
json_encode(table, [bool]): json encode, bool arg means json pretty formatting
shell(string)             : exec shell command
!string                   : exec shell command
help()                    : show this help information

Everything follows Lua grammar.
Good luck.

```
