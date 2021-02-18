package dns

// DNSHeader describes the DNS header as documented in RFC 1035
type DNSHeader struct {
	Identifier                     uint16 `json:"identifier"`
	Flags                          uint16 `json:"flags"`
	TotalQuestions                 uint16 `json:"num_questions"`
	TotalAnswerResourceRecords     uint16 `json:"num_answers"`
	TotalAuthorityResourceRecords  uint16 `json:"num_authority"`
	TotalAdditionalResourceRecords uint16 `json:"num_additional"`
}

// DNSFlags describe the flags of the DNS header
type DNSFlags struct {
	QueryResponse       bool  `json:"query_response"`
	OpCode              uint8       `json:"op_code"`
	AuthoritativeAnswer bool                `json:"authoritative_answer"`
	Truncated           bool                `json:"truncated"`
	RecursionDesired    bool                `json:"recursion_desired"`
	RecursionAvailable  bool                `json:"recursion_available"`
	ResponseCode        uint8 `json:"response_code"`
}

// DNSPDU describes the DNS protocol data unit as documented in RFC 1035
type DNSPDU struct {
	Header                    DNSHeader           `json:"header"`
	Flags                     DNSFlags            `json:"flags"`
	/*Questions                 []DNSQuestion       `json:"questions"`
	AnswerResourceRecords     []DNSResourceRecord `json:"answers"`
	AuthorityResourceRecords  []DNSResourceRecord `json:"authority"`
	AdditionalResourceRecords []DNSResourceRecord `json:"additional"`*/
}