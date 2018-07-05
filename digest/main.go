package main

import (
	"crypto/md5"
	"io"
	"encoding/hex"
	"strings"
)

const nonce = "UAZs1dp3wX5BtXEpoCXKO2lHhap564rX"
const cnonce = "oaSHizKi0RcJXmFE2TMtW8IefL799dWU"
const nc = "00000001"
const qop = "auth"
const password = "time"
const method = "POST"

func main() {
	s := "Digest " +
		"username=\"go\", " +
		"realm=\"example@restvoice.org\", " +
		"nonce=\"UAZs1dp3wX5BtXEpoCXKO2lHhap564rX\", " +
		"uri=\"/customers/1/invoices\", " +
		"algorithm=\"MD5\", qop=auth, " +
		"nc=00000001, " +
		"cnonce=\"oaSHizKi0RcJXmFE2TMtW8IefL799dWU\", " +
		"response=\"74e05470ffcf73054539fc6a54e0a2a0\", " +
		"opaque=\"xU2Z4FyqwKUBdwTMRYdGtAG1ppaT0bNm\""

	authFields := digestParts(s)
	step1 := hash(authFields["username"] + ":" + authFields["realm"] + ":" + password)
	step2 := hash(method + ":" + authFields["uri"])
	response := hash(step1 + ":" + authFields["nonce"] + ":" + authFields["nc"] + ":" + authFields["cnonce"]+ ":" + authFields["qop"] + ":" + step2)
	if response == authFields["response"] {
		print("success")
	}
}

func hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

func digestParts(authorization string) map[string]string {
	result := map[string]string{}
	wantedHeaders := []string{"username", "nonce", "realm", "qop", "uri", "nc", "response", "opaque", "cnonce"}
	requestHeaders := strings.Split(authorization, ",")
	for _, r := range requestHeaders {
		for _, w := range wantedHeaders {
			if strings.Contains(r, " " + w) {
				v  := strings.Split(r, "=")[1]
				result[w] = strings.Trim(v, `"`)
			}
		}
	}
	return result
}
