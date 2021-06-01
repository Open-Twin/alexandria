package dns

import "time"

/*
 * definitions of the following structs and its properties can be found at the RFC 1035:
 * https://tools.ietf.org/html/rfc1035
 */

// DNSHeader describes the DNS header as documented in RFC 1035
type DNSHeader struct {
	// A 16 bit identifier assigned by the program that generates any kind of query.
	//This identifier is copied the corresponding reply and can be used by the requester to match up replies to outstanding queries.
	Identifier uint16 `json:"identifier"`
	//A one bit field that specifies whether this message is a query (0), or a response (1).
	Flags                          uint16 `json:"flags"`
	TotalQuestions                 uint16 `json:"num_questions"`
	TotalAnswerResourceRecords     uint16 `json:"num_answers"`
	TotalAuthorityResourceRecords  uint16 `json:"num_authority"`
	TotalAdditionalResourceRecords uint16 `json:"num_additional"`
}

// DNSFlags describe the flags of the DNS header
type DNSFlags struct {
	// Specifies whether this message is a query (false) or a response (true)
	QueryResponse bool `json:"query_response"`
	// Specifies the kind of query
	OpCode uint8 `json:"op_code"`
	// Specifies whether name server is an authority for the domain naim in question section
	AuthoritativeAnswer bool `json:"authoritative_answer"`
	// Specifies whether this message was truncated due to length greater than permitted
	Truncated bool `json:"truncated"`
	// Specifies whether the response should be pursued recursively
	RecursionDesired bool `json:"recursion_desired"`
	// Specifies if recursive support is available on the server
	RecursionAvailable bool `json:"recursion_available"`
	//Z
	Z bool `json:"z"`
	//Authentic Data
	AuthenticData bool `json:"authentic_data"`
	//Checking disabled
	CheckingDisabled bool `json:"checking_disabled"`
	// Response code
	ResponseCode uint8 `json:"response_code"`
}

// DNSResourceRecord describes individual records in the request and response of the DNS payload body
type DNSResourceRecord struct {
	Labels []string `json:"labels"`
	// RR class type c
	Type uint16 `json:"type"`
	// RR class code
	Class uint16 `json:"class"`
	// Time that the resource record can be cached before it is queried again
	TimeToLive uint32 `json:"ttl"`
	// length of the resource data
	ResourceDataLength uint16 `json:"rd_length"`
	ResourceData       []byte `json:"rd"`
}

// DNSQuestion describes individual records in the request and response of the DNS payload body
type DNSQuestion struct {
	Labels []string `json:"labels"`
	// a two octet code which specifies the type of the query.
	//The values for this field include all codes valid for a TYPE field, together with some more general codes which can match more than one type of RR.
	Type uint16 `json:"type"`
	//a two octet code that specifies the class of the query.
	//For example, the QCLASS field is IN for the Internet.
	Class uint16 `json:"class"`
}

// DNSPDU describes the DNS protocol data unit as documented in RFC 1035
type DNSPDU struct {
	Header                    DNSHeader           `json:"header"`
	Flags                     DNSFlags            `json:"flags"`
	Questions                 []DNSQuestion       `json:"questions"`
	AnswerResourceRecords     []DNSResourceRecord `json:"answers"`
	AuthorityResourceRecords  []DNSResourceRecord `json:"authority"`
	AdditionalResourceRecords []DNSResourceRecord `json:"additional"`
}

type NodeHealth struct {
	Healthy     bool
	Connections int
	LastOnline  time.Time
}
