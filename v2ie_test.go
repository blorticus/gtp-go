package gtp

import (
	"fmt"
	"testing"
)

type v2IEComparable struct {
	ieOctets   []byte
	matchingIE *V2IE
}

type v2IEFailCase struct {
	name        string
	inputStream []byte
}

func TestV2IEDecodeInvalidCases(t *testing.T) {
	cases := []v2IEFailCase{
		v2IEFailCase{
			name:        "Empty stream",
			inputStream: []byte{},
		},
		v2IEFailCase{
			name:        "Stream too short for header",
			inputStream: []byte{0x01, 0x00, 0x06},
		},
		v2IEFailCase{
			name:        "Header only",
			inputStream: []byte{0x01, 0x00, 0x06, 0x00},
		},
		v2IEFailCase{
			name:        "Insufficient byte stream length",
			inputStream: []byte{0x01, 0x00, 0x06, 0x00, 0x12, 0x34, 0x56, 0x78},
		},
	}

	for _, testCase := range cases {
		_, err := DecodeV2IE(testCase.inputStream)

		if err == nil {
			t.Errorf("(%s) Expected error on DecodeV2IE(), but received none", testCase.name)
		}
	}
}

func TestV2IEDecodeValidCases(t *testing.T) {
	cases := []v2IEComparable{
		v2IEComparable{
			ieOctets: []byte{0x56, 0x00, 0x0d, 0x00, 0x18, 0x01, 0x00, 0x01, 0xff, 0x00, 0x01, 0x00, 0x01, 0x0f, 0x42, 0x4d, 0x00},
			matchingIE: &V2IE{
				Type:           UserLocationInformation,
				DataLength:     13,
				TotalLength:    17,
				InstanceNumber: 0,
				Data:           []byte{0x18, 0x01, 0x00, 0x01, 0xff, 0x00, 0x01, 0x00, 0x01, 0x0f, 0x42, 0x4d, 0x00},
			},
		},
		v2IEComparable{
			ieOctets: []byte{0x52, 0x00, 0x01, 0x03, 0x06},
			matchingIE: &V2IE{
				Type:           RATType,
				DataLength:     1,
				TotalLength:    5,
				InstanceNumber: 3,
				Data:           []byte{0x06},
			},
		},
	}

	testCaseNumber := 0
	for _, testCase := range cases {
		testCaseNumber++
		ie, err := DecodeV2IE(testCase.ieOctets)

		if err != nil {
			t.Errorf("For test case number [%d], received error on decode: %s", testCaseNumber, err)
			continue
		}

		if err = compareTwoV2IEObjects(testCase.matchingIE, ie); err != nil {
			t.Errorf("For test case number [%d]: %s", testCaseNumber, err)
		}
	}
}

func compareTwoV2IEObjects(expected *V2IE, got *V2IE) error {
	if expected.Type != got.Type {
		return fmt.Errorf("Expected IE Type [%d] (%s), got [%d] (%s)", expected.Type, NameOfV2IEForType(expected.Type), got.Type, NameOfV2IEForType(got.Type))
	}

	if expected.DataLength != got.DataLength {
		return fmt.Errorf("Expected IE DataLength = %d, got = %d", expected.DataLength, got.DataLength)
	}

	if expected.TotalLength != got.TotalLength {
		return fmt.Errorf("Expected IE TotalLength = %d, got = %d", expected.TotalLength, got.TotalLength)
	}

	if expected.InstanceNumber != got.InstanceNumber {
		return fmt.Errorf("Expected IE InstanceNumber = %d, got = %d", expected.InstanceNumber, got.InstanceNumber)
	}

	if len(expected.Data) != len(got.Data) {
		return fmt.Errorf("Expected IE Data has (%d) bytes, got = (%d) bytes", len(expected.Data), len(got.Data))
	}

	for i := 0; i < len(expected.Data); i++ {
		if expected.Data[i] != got.Data[i] {
			return fmt.Errorf("Expected IE Data byte index (%d) is (0x%02x), got (0x%02x)", i, expected.Data[i], got.Data[i])
		}
	}

	return nil
}

///
// 0020    5c 00 01 00 00 5d 00 12 00 49 00 01 00 05
// 0030   57 00 09 00 80 e4 03 fb 94 ac 13 01 b2 03 00 01
// 0040   00 95