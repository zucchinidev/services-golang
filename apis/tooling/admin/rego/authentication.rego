package sales.rego

# rego.v1 allows me to use functions like decode_verify
import rego.v1

# set the auth variable with default value to false
default auth := false

auth if {
	[valid, _, _] := verify_jwt

	# if valid is true, then return back the value and then auth will be assigned with true value
	valid = true
}

verify_jwt := io.jwt.decode_verify(input.Token, {
	"cert": input.Key,
	"iss": input.ISS,
})
