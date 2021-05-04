package dns

import (
	"bytes"
	"encoding/binary"
	"github.com/rs/zerolog/log"
	"reflect"
)

func (flags DNSFlags) WriteFlags() uint16 {
	result := uint16(0)

	if flags.QueryResponse {
		result |= uint16(0b1000_0000_0000_0000)
	}

	result |= (uint16(flags.OpCode) << 11) & uint16(0b01111_1000_0000_0000)

	if flags.AuthoritativeAnswer {
		result |= uint16(0b0000_0100_0000_0000)
	}

	if flags.Truncated {
		result |= uint16(0b0000_0010_0000_0000)
	}

	if flags.RecursionDesired {
		result |= uint16(0b0000_0001_0000_0000)
	}

	if flags.RecursionAvailable {
		result |= uint16(0b0000_0000_1000_0000)
	}
	if flags.Z {
		result |= uint16(0b0000_0000_0100_0000)
	}
	if flags.AuthenticData {
		result |= uint16(0b0000_0000_0010_0000)
	}
	if flags.CheckingDisabled {
		result |= uint16(0b0000_0000_0001_0000)
	}

	result |= uint16(uint8(flags.ResponseCode) & uint8(0b0000_1111))

	return result
}

func writeLabels(responseBuffer *bytes.Buffer, labels []string) error {
	//TODO: Pointer ??
	if reflect.DeepEqual(labels, []string{"P", "O", "I", "N", "T", "E", "R"}) {
		_, err := responseBuffer.Write([]byte{0xc0, 0x0c})
		return err
	}

	for _, label := range labels {
		labelLength := len(label)
		labelBytes := []byte(label)

		responseBuffer.WriteByte(byte(labelLength))
		responseBuffer.Write(labelBytes)
	}

	err := responseBuffer.WriteByte(byte(0))

	return err
}

func writeResourceRecords(buffer *bytes.Buffer, rrs []DNSResourceRecord) error {
	for _, rr := range rrs {
		err := writeLabels(buffer, rr.Labels)
		if err != nil {
			return err
		}

		err = binary.Write(buffer, binary.BigEndian, rr.Type)
		if err != nil {
			return err
		}

		err = binary.Write(buffer, binary.BigEndian, rr.Class)
		if err != nil {
			return err
		}

		err = binary.Write(buffer, binary.BigEndian, rr.TimeToLive)
		if err != nil {
			return err
		}

		err = binary.Write(buffer, binary.BigEndian, rr.ResourceDataLength)
		if err != nil {
			return err
		}

		err = binary.Write(buffer, binary.BigEndian, rr.ResourceData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pdu DNSPDU) Bytes() ([]byte, error) {

	var responseBuffer = new(bytes.Buffer)

	pdu.Header.Flags = pdu.Flags.WriteFlags()

	err := binary.Write(responseBuffer, binary.BigEndian, &pdu.Header)

	if err != nil {
		return nil, err
	}

	for _, question := range pdu.Questions {
		err := writeLabels(responseBuffer, question.Labels)
		if err != nil {
			return nil, err
		}

		err = binary.Write(responseBuffer, binary.BigEndian, question.Type)
		if err != nil {
			return nil, err
		}

		err = binary.Write(responseBuffer, binary.BigEndian, question.Class)
		if err != nil {
			return nil, err
		}
	}
	log.Debug().Msgf("Generated Answers: %v", pdu.AnswerResourceRecords)
	err = writeResourceRecords(responseBuffer, pdu.AnswerResourceRecords)
	if err != nil {
		return nil, err
	}

	err = writeResourceRecords(responseBuffer, pdu.AuthorityResourceRecords)
	if err != nil {
		return nil, err
	}

	err = writeResourceRecords(responseBuffer, pdu.AdditionalResourceRecords)
	if err != nil {
		return nil, err
	}

	return responseBuffer.Bytes(), nil
}