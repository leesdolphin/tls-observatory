package mozillaEvaluationWorker

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mozilla/tls-observatory/connection"
	"github.com/mozilla/tls-observatory/logger"
	"github.com/mozilla/tls-observatory/worker"
)

var workerName = "mozillaEvaluationWorker"
var workerDesc = `The evaluation worker provided insight on the compliance level of the tls configuration of the audited target.
For more info check https://wiki.mozilla.org/Security/Server_Side_TLS.`

var sstls ServerSideTLSJson
var modern, intermediate, old Configuration

var log = logger.GetLogger()

func init() {
	err := json.Unmarshal([]byte(ServerSideTLSConfiguration), &sstls)
	if err != nil {
		log.Error("Could not load Server Side TLS configuration. Evaluation Worker not available")
		return
	}
	modern = sstls.Configurations["modern"]
	intermediate = sstls.Configurations["intermediate"]
	old = sstls.Configurations["old"]
	worker.RegisterWorker(workerName, worker.Info{Runner: new(eval), Description: workerDesc})
}

type ServerSideTLSJson struct {
	Configurations map[string]Configuration `json:"configurations"`
	Version        float64                  `json:"version"`
}

// Configuration represents configurations levels declared by the Mozilla server-side-tls
// see https://wiki.mozilla.org/Security/Server_Side_TLS
type Configuration struct {
	Ciphersuite          string   `json:"ciphersuite"`
	Ciphers              []string `json:"ciphers"`
	TLSVersions          []string `json:"tls_versions"`
	TLSCurves            []string `json:"tls_curves"`
	CertificateType      string   `json:"certificate_type"`
	CertificateCurve     string   `json:"certificate_curve"`
	CertificateSignature string   `json:"certificate_signature"`
	RsaKeySize           float64  `json:"rsa_key_size"`
	DHParamSize          float64  `json:"dh_param_size"`
	ECDHParamSize        float64  `json:"ecdh_param_size"`
	Hsts                 string   `json:"hsts"`
	OldestClients        []string `json:"oldest_clients"`
}

// EvaluationResults contains the results of the mozillaEvaluationWorker
type EvaluationResults struct {
	Level    string              `json:"level"`
	Failures map[string][]string `json:"failures"`
}

type eval struct {
}

// Run implements the worker interface.It is called to get the worker results.
func (e eval) Run(in worker.Input, resChan chan worker.Result) {

	res := worker.Result{WorkerName: workerName}

	b, err := Evaluate(in.Connection, in.Certificate.SignatureAlgorithm)
	if err != nil {
		res.Success = false
		res.Errors = append(res.Errors, err.Error())
	} else {
		res.Result = b
		res.Success = true
	}

	resChan <- res
}

// Evaluate runs compliance checks of the provided json Stored connection and returns the results
func Evaluate(connInfo connection.Stored, certsigalg string) ([]byte, error) {

	var isO, isI, isM, isB bool

	results := EvaluationResults{}
	results.Failures = make(map[string][]string)

	// assume the worst
	results.Level = "bad"

	isO, results.Failures["old"] = isOld(connInfo, certsigalg)
	if isO {
		results.Level = "old"

		ord, ordres := isOrdered(connInfo, old.Ciphers, "old")
		if !ord {
			results.Level += " with bad ordering"
			results.Failures["old"] = append(results.Failures["old"], ordres...)
		}
	}

	isI, results.Failures["intermediate"] = isIntermediate(connInfo, certsigalg)
	if isI {
		results.Level = "intermediate"

		ord, ordres := isOrdered(connInfo, intermediate.Ciphers, "intermediate")
		if !ord {
			results.Level += " with bad ordering"
			results.Failures["intermediate"] = append(results.Failures["intermediate"], ordres...)
		}
	}

	isM, results.Failures["modern"] = isModern(connInfo, certsigalg)
	if isM {
		results.Level = "modern"

		ord, ordres := isOrdered(connInfo, modern.Ciphers, "modern")
		if !ord {
			results.Level += " with bad ordering"
			results.Failures["modern"] = append(results.Failures["modern"], ordres...)
		}
	}

	isB, results.Failures["bad"] = isBad(connInfo)
	if isB {
		results.Level = "bad"
	}

	js, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func isBad(c connection.Stored) (bool, []string) {

	var failures, allProtos, allCiphers []string
	status := false
	hasSSLv2 := false
	hasBadPFS := false
	hasBadPK := false
	hasMD5 := false

	for _, cs := range c.CipherSuite {

		allCiphers = append(allCiphers, cs.Cipher)

		if contains(cs.Protocols, "SSLv2") {
			hasSSLv2 = true
		}

		for _, proto := range cs.Protocols {
			if !contains(allProtos, proto) {
				allProtos = append(allProtos, proto)
			}
		}

		if cs.PFS != "None" {
			if !hasGoodPFS(cs.PFS, old.DHParamSize, old.ECDHParamSize, false, false) {
				hasBadPFS = true
			}
		}

		if cs.PubKey < old.RsaKeySize {
			hasBadPK = true
		}

		if cs.SigAlg == "md5WithRSAEncryption" {
			hasMD5 = true
		}
	}

	badCiphers := extra(old.Ciphers, allCiphers)
	if len(badCiphers) > 0 {
		for _, c := range badCiphers {
			failures = append(failures, fmt.Sprintf("remove cipher %s", c))
			status = true
		}
	}

	if hasSSLv2 {
		failures = append(failures, "disable SSLv2")
		status = true
	}

	if hasBadPFS {
		status = true
		failures = append(failures,
			fmt.Sprintf("don't use DHE smaller than %.0fbits or ECC smaller than %.0fbits",
				old.DHParamSize, old.ECDHParamSize))
	}

	if hasBadPK {
		failures = append(failures, "don't use a public key shorter than 2048bit")
		status = true
	}

	if hasMD5 {
		failures = append(failures, "don't use an MD5 signature")
		status = true
	}

	return status, failures
}

func isOld(c connection.Stored, certsigalg string) (bool, []string) {

	status := true
	var allProtos []string
	hasSHA1 := true
	has3DES := false
	hasSSLv3 := false
	hasOCSP := true
	hasPFS := true
	var failures []string

	for _, cs := range c.CipherSuite {

		if !contains(old.Ciphers, cs.Cipher) {
			failures = append(failures, fmt.Sprintf("remove %s cipher", cs.Cipher))
			status = false
		}

		if cs.Cipher == "DES-CBC3-SHA" {
			has3DES = true
		}

		if contains(cs.Protocols, "SSLv3") {
			hasSSLv3 = true
		}

		for _, proto := range cs.Protocols {
			if !contains(allProtos, proto) {
				allProtos = append(allProtos, proto)
			}
		}

		if cs.PFS != "None" {
			if !hasGoodPFS(cs.PFS, old.DHParamSize, old.ECDHParamSize, true, false) {
				hasPFS = false
			}
		}

		if certsigalg != old.CertificateSignature {

			fail := fmt.Sprintf("%s is not an old signature", cs.SigAlg)
			if !contains(failures, fail) {
				failures = append(failures, fail)
			}

			hasSHA1 = false
			status = false
		}

		if !cs.OCSPStapling {
			hasOCSP = false
		}
	}

	extraProto := extra(old.TLSVersions, allProtos)
	for _, p := range extraProto {
		failures = append(failures, fmt.Sprintf("disable %s protocol", p))
		status = false
	}

	missingProto := extra(allProtos, old.TLSVersions)
	for _, p := range missingProto {
		failures = append(failures, fmt.Sprintf("consider enabling %s", p))
	}

	if !c.ServerSide {
		failures = append(failures, "enforce ServerSide ordering")
	}

	if !hasSSLv3 {
		failures = append(failures, "add SSLv3 support")
		status = false
	}

	if !hasOCSP {
		failures = append(failures, "consider enabling OCSP stapling")
	}

	if !hasSHA1 {
		failures = append(failures, "it is recommended to use sha1")
		status = false
	}

	if !has3DES {
		failures = append(failures, "add cipher DES-CBC3-SHA")
		status = false
	}

	if !hasPFS {
		status = false
		failures = append(failures,
			fmt.Sprintf("use DHE of %.0fbits and ECC of %.0fbits",
				old.DHParamSize, old.ECDHParamSize))
	}

	return status, failures
}

func isIntermediate(c connection.Stored, certsigalg string) (bool, []string) {

	status := true
	var allProtos []string
	hasTLSv1 := false
	hasAES := false
	hasSHA256 := true
	hasOCSP := true
	hasPFS := true
	var failures []string

	for _, cs := range c.CipherSuite {

		if !contains(intermediate.Ciphers, cs.Cipher) {
			failures = append(failures, fmt.Sprintf("remove %s cipher", cs.Cipher))
			status = false
		}

		for _, proto := range cs.Protocols {
			if !contains(allProtos, proto) {
				allProtos = append(allProtos, proto)
			}
		}

		if contains(cs.Protocols, "TLSv1") {
			hasTLSv1 = true
		}

		if cs.Cipher == "AES128-SHA" {
			hasAES = true
		}

		if cs.PFS != "None" {
			if !hasGoodPFS(cs.PFS, intermediate.DHParamSize, intermediate.ECDHParamSize, false, false) {
				hasPFS = false
			}
		}

		if certsigalg != intermediate.CertificateSignature {

			fail := fmt.Sprintf("%s is not an intermediate signature", cs.SigAlg)
			if !contains(failures, fail) {
				failures = append(failures, fail)
			}
			hasSHA256 = false
		}

		if !cs.OCSPStapling {
			hasOCSP = false
		}
	}

	extraProto := extra(intermediate.TLSVersions, allProtos)
	for _, p := range extraProto {
		failures = append(failures, fmt.Sprintf("disable %s protocol", p))
		status = false
	}

	if !hasAES {
		failures = append(failures, "add cipher AES128-SHA")
		status = false
	}

	if !hasTLSv1 {
		failures = append(failures, "consider adding TLSv1")
		status = false
	}

	if !c.ServerSide {
		failures = append(failures, "enforce ServerSide ordering")
	}

	if !hasOCSP {
		failures = append(failures, "consider enabling OCSP stapling")
	}

	if !hasSHA256 {
		failures = append(failures, "it is recommended to use sha256")
	}

	if !hasPFS {
		failures = append(failures, "use DHE of at least 2048bits and ECC of at least 256bits")
	}

	return status, failures
}

func isModern(c connection.Stored, certsigalg string) (bool, []string) {

	status := true
	var allProtos []string
	hasSHA256 := true
	hasOCSP := true
	hasPFS := true
	var failures []string

	for _, cs := range c.CipherSuite {

		if !contains(modern.Ciphers, cs.Cipher) {
			failures = append(failures, fmt.Sprintf("remove %s cipher", cs.Cipher))
			status = false
		}

		for _, proto := range cs.Protocols {
			if !contains(allProtos, proto) {
				allProtos = append(allProtos, proto)
			}
		}

		if cs.PFS != "None" {
			if !hasGoodPFS(cs.PFS, modern.DHParamSize, modern.ECDHParamSize, false, false) {
				hasPFS = false
			}
		}

		if certsigalg != modern.CertificateSignature {
			fail := fmt.Sprintf("%s is not a modern signature", cs.SigAlg)
			if !contains(failures, fail) {
				failures = append(failures, fail)
			}
			hasSHA256 = false
		}

		if !cs.OCSPStapling {
			hasOCSP = false
		}
	}

	extraProto := extra(modern.TLSVersions, allProtos)
	for _, p := range extraProto {
		failures = append(failures, fmt.Sprintf("disable %s protocol", p))
		status = false
	}

	if !c.ServerSide {
		failures = append(failures, "enforce ServerSide ordering")
	}

	if !hasOCSP {
		failures = append(failures, "consider enabling OCSP stapling")
	}

	if !hasSHA256 {
		failures = append(failures, "it is recommended to use sha256")
		status = false
	}

	if !hasPFS {
		status = false
		failures = append(failures,
			fmt.Sprintf("use DHE of at least %.0fbits and ECC of at least %.0fbits",
				modern.DHParamSize, modern.ECDHParamSize))
	}
	return status, failures
}

func isOrdered(c connection.Stored, conf []string, level string) (bool, []string) {

	var failures []string
	status := true
	prevpos := 0

	for _, ciphersuite := range c.CipherSuite {
		for pos, cipher := range conf {
			if ciphersuite.Cipher == cipher {
				if pos < prevpos {
					failures = append(failures, fmt.Sprintf("increase priority of %s over %s", ciphersuite.Cipher, conf[prevpos]))
					status = false
				}
				prevpos = pos
			}
		}
	}

	if !status {
		failures = append(failures, fmt.Sprintf("fix ciphersuite ordering, use recommended %s ciphersuite", level))
	}
	return status, failures
}

func hasGoodPFS(curPFS string, targetDH, targetECC float64, mustMatchDH, mustMatchECDH bool) bool {
	pfs := strings.Split(curPFS, ",")
	if len(pfs) < 2 {
		return false
	}

	if "ECDH" == pfs[0] {
		bitsStr := strings.TrimRight(pfs[2], "bits")

		bits, err := strconv.ParseFloat(bitsStr, 64)
		if err != nil {
			return false
		}

		if mustMatchECDH {
			if bits != targetECC {
				return false
			}
		} else {
			if bits < targetECC {
				return false
			}
		}

	} else if "DH" == pfs[0] {
		bitsStr := strings.TrimRight(pfs[1], "bits")

		bits, err := strconv.ParseFloat(bitsStr, 64)
		if err != nil {
			return false
		}

		if mustMatchDH {
			if bits != targetDH {
				return false
			}
		} else {
			if bits < targetDH {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func extra(s []string, e []string) []string {

	var extra []string
func (e eval) PrintAnalysis(r []byte) (results []string, err error) {
	var (
		eval           EvaluationResults
		previousissues []string
		prefix         string
	)
	err = json.Unmarshal(r, &eval)
	if err != nil {
		err = fmt.Errorf("Mozilla evaluation worker: failed to parse results: %v", err)
		return
	}
	results = append(results, fmt.Sprintf("* Mozilla evaluation: %s", eval.Level))
	for _, lvl := range []string{"bad", "old", "intermediate", "modern"} {
		if _, ok := eval.Failures[lvl]; ok && len(eval.Failures[lvl]) > 0 {
			for _, issue := range eval.Failures[lvl] {
				for _, previousissue := range previousissues {
					if issue == previousissue {
						goto next
					}
				}
				prefix = "for " + lvl + " level:"
				if lvl == "bad" {
					prefix = "bad configuration:"
				}
				results = append(results, fmt.Sprintf("  - %s %s", prefix, issue))
				previousissues = append(previousissues, issue)
			next:
			}

	for _, str := range e {
		if !contains(s, str) {
			extra = append(extra, str)
		}
	}
	return extra
	if eval.Level != "bad" {
		results = append(results,
			fmt.Sprintf("  - oldest clients: %s", strings.Join(sstls.Configurations[eval.Level].OldestClients, ", ")))
	}
	return
}
