package gtp

import (
	"fmt"
	"testing"
)

type v2PDUComparable struct {
	testName     string
	pduOctets    []byte
	matchingPdu  *V2PDU
	piggybackPdu *V2PDU
}

type v2PDUNamesComparable struct {
	expectedName string
	pduType      V2MessageType
}

func TestPDUNames(t *testing.T) {
	// This test set is mostly to make sure the list doesn't accidentally
	// get shifted if values are changed
	testCases := []v2PDUNamesComparable{
		v2PDUNamesComparable{"Reserved", 0},
		v2PDUNamesComparable{"Echo Request", 1},
		v2PDUNamesComparable{"Create Session Response", 33},
		v2PDUNamesComparable{"Resume Notification", 164},
		v2PDUNamesComparable{"Reserved", 172},
		v2PDUNamesComparable{"Modify Access Bearers Request", 211},
		v2PDUNamesComparable{"MBMS Session Stop Response", 236},
	}

	for _, testCase := range testCases {
		if NameOfV2MessageForType(testCase.pduType) != testCase.expectedName {
			t.Errorf("For PDU Message Type (%d), expected name = (%s), got = (%s)", testCase.pduType, testCase.expectedName, NameOfV2MessageForType(testCase.pduType))
		}
	}
}

func TestV2PDUDecodeValidCasesNoPiggyback(t *testing.T) {
	testCases := []v2PDUComparable{
		v2PDUComparable{
			testName: "Valid Modify Bearer Request",
			pduOctets: []byte{
				// PDU Header
				0x48, 0x22, 0x00, 0x3e, 0x05, 0x40, 0x3b, 0x2e, 0x00, 0x1a, 0xcc, 0x00,
				// ULI
				0x56, 0x00, 0x0d, 0x00, 0x18, 0x00, 0x11, 0x00, 0xff, 0x00, 0x00, 0x11,
				0x00, 0x0f, 0x42, 0x4d, 0x00,
				// RATType
				0x52, 0x00, 0x01, 0x00, 0x06,
				// Delay Value
				0x5c, 0x00, 0x01, 0x00, 0x00,
				// Bearer Context
				0x5d, 0x00, 0x12, 0x00, 0x49, 0x00, 0x01, 0x00, 0x05, 0x57, 0x00, 0x09,
				0x00, 0x80, 0xe4, 0x03, 0xfb, 0x94, 0xac, 0x13, 0x01, 0xb2,
				// Recovery
				0x03, 0x00, 0x01, 0x00, 0x95,
			},
			matchingPdu: &V2PDU{
				Type:                     ModifyBearerRequest,
				IsCarryingPiggybackedPDU: false,
				PriorityFieldIsPresent:   false,
				TEIDFieldIsPresent:       true,
				SequenceNumber:           0x00001acc,
				Priority:                 0,
				TEID:                     0x05403b2e,
				TotalLength:              0x0042,
				InformationElements: []*V2IE{
					&V2IE{
						Type:           UserLocationInformation,
						InstanceNumber: 0,
						TotalLength:    17,
						DataLength:     13,
						Data: []byte{
							0x18, 0x00, 0x11, 0x00, 0xff, 0x00, 0x00, 0x11,
							0x00, 0x0f, 0x42, 0x4d, 0x00,
						},
					},
					&V2IE{
						Type:           RATType,
						InstanceNumber: 0,
						TotalLength:    5,
						DataLength:     1,
						Data:           []byte{0x06},
					},
					&V2IE{
						Type:           DelayValue,
						InstanceNumber: 0,
						TotalLength:    5,
						DataLength:     1,
						Data:           []byte{0x00},
					},
					&V2IE{
						Type:           BearerContext,
						InstanceNumber: 0,
						TotalLength:    22,
						DataLength:     18,
						Data: []byte{
							0x49, 0x00, 0x01, 0x00, 0x05, 0x57, 0x00, 0x09,
							0x00, 0x80, 0xe4, 0x03, 0xfb, 0x94, 0xac, 0x13,
							0x01, 0xb2,
						},
					},
					&V2IE{
						Type:           RecoveryRestartCounter,
						InstanceNumber: 0,
						TotalLength:    5,
						DataLength:     1,
						Data:           []byte{0x95},
					},
				},
			},
			piggybackPdu: nil,
		},
		v2PDUComparable{
			testName: "Truncated Modify Bearer Requests Piggybacked with anoth MBR",
			pduOctets: []byte{
				// PDU Header
				0x58, 0x22, 0x00, 0x1e, 0x05, 0x40, 0x3b, 0x2e, 0x00, 0x1a, 0xcc, 0x00,
				// ULI
				0x56, 0x00, 0x0d, 0x00, 0x18, 0x00, 0x11, 0x00, 0xff, 0x00, 0x00, 0x11,
				0x00, 0x0f, 0x42, 0x4d, 0x00,
				// RATType
				0x52, 0x00, 0x01, 0x00, 0x06,
				// Piggybacked PDU Header
				0x48, 0x22, 0x00, 0x28, 0x05, 0x40, 0x3b, 0x2e, 0x00, 0x1a, 0xcc, 0x00,
				// Delay Value
				0x5c, 0x00, 0x01, 0x00, 0x00,
				// Bearer Context
				0x5d, 0x00, 0x12, 0x00, 0x49, 0x00, 0x01, 0x00, 0x05, 0x57, 0x00, 0x09,
				0x00, 0x80, 0xe4, 0x03, 0xfb, 0x94, 0xac, 0x13, 0x01, 0xb2,
				// Recovery
				0x03, 0x00, 0x01, 0x00, 0x95,
			},
			matchingPdu: &V2PDU{
				Type:                     ModifyBearerRequest,
				IsCarryingPiggybackedPDU: true,
				PriorityFieldIsPresent:   false,
				TEIDFieldIsPresent:       true,
				SequenceNumber:           0x00001acc,
				Priority:                 0,
				TEID:                     0x05403b2e,
				TotalLength:              0x0022,
				InformationElements: []*V2IE{
					&V2IE{
						Type:           UserLocationInformation,
						InstanceNumber: 0,
						TotalLength:    17,
						DataLength:     13,
						Data: []byte{
							0x18, 0x00, 0x11, 0x00, 0xff, 0x00, 0x00, 0x11,
							0x00, 0x0f, 0x42, 0x4d, 0x00,
						},
					},
					&V2IE{
						Type:           RATType,
						InstanceNumber: 0,
						TotalLength:    5,
						DataLength:     1,
						Data:           []byte{0x06},
					},
				},
			},
			piggybackPdu: &V2PDU{
				Type:                     ModifyBearerRequest,
				IsCarryingPiggybackedPDU: false,
				PriorityFieldIsPresent:   false,
				TEIDFieldIsPresent:       true,
				SequenceNumber:           0x00001acc,
				Priority:                 0,
				TEID:                     0x05403b2e,
				TotalLength:              0x002c,
				InformationElements: []*V2IE{
					&V2IE{
						Type:           DelayValue,
						InstanceNumber: 0,
						TotalLength:    5,
						DataLength:     1,
						Data:           []byte{0x00},
					},
					&V2IE{
						Type:           BearerContext,
						InstanceNumber: 0,
						TotalLength:    22,
						DataLength:     18,
						Data: []byte{
							0x49, 0x00, 0x01, 0x00, 0x05, 0x57, 0x00, 0x09,
							0x00, 0x80, 0xe4, 0x03, 0xfb, 0x94, 0xac, 0x13,
							0x01, 0xb2,
						},
					},
					&V2IE{
						Type:           RecoveryRestartCounter,
						InstanceNumber: 0,
						TotalLength:    5,
						DataLength:     1,
						Data:           []byte{0x95},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		pdu, piggybackPdu, err := DecodeV2PDU(testCase.pduOctets)

		if err != nil {
			t.Errorf("(%s) Failed to decode, err = (%s)", testCase.testName, err)
			continue
		}

		if err = compareTwoV2PDUObjects(testCase.matchingPdu, pdu); err != nil {
			t.Errorf("(%s) %s", testCase.testName, err)
		}

		if piggybackPdu != nil {
			if testCase.piggybackPdu == nil {
				t.Errorf("(%s) On decode, received unexpected piggybacked PDU", testCase.testName)
			} else {
				if err = compareTwoV2PDUObjects(testCase.piggybackPdu, piggybackPdu); err != nil {
					t.Errorf("(%s) piggyback PDU: %s", testCase.testName, err)
				}
			}
		} else {
			if testCase.piggybackPdu != nil {
				t.Errorf("(%s) On decode, should have received piggyback PDU, but did not", testCase.testName)
			}
		}

	}
}

func compareTwoV2PDUObjects(expected *V2PDU, got *V2PDU) error {
	if expected.Type != got.Type {
		return fmt.Errorf("Expected Type = (%d) [%s], got = (%d) [%s]", expected.Type, NameOfV2MessageForType(expected.Type), got.Type, NameOfV2MessageForType(got.Type))
	}

	if expected.IsCarryingPiggybackedPDU != got.IsCarryingPiggybackedPDU {
		return fmt.Errorf("Expected IsCarryingPiggybackedPDU = (%t), got = (%t)", expected.IsCarryingPiggybackedPDU, got.IsCarryingPiggybackedPDU)
	}

	if expected.PriorityFieldIsPresent != got.PriorityFieldIsPresent {
		return fmt.Errorf("Expected PriorityFieldIsPresent = (%t), got = (%t)", expected.PriorityFieldIsPresent, got.PriorityFieldIsPresent)
	}

	if expected.TEIDFieldIsPresent != got.TEIDFieldIsPresent {
		return fmt.Errorf("Expected TEIDFieldIsPresent = (%t), got = (%t)", expected.TEIDFieldIsPresent, got.TEIDFieldIsPresent)
	}

	if expected.SequenceNumber != got.SequenceNumber {
		return fmt.Errorf("Expected SequenceNumber = (%d), got = (%d)", expected.SequenceNumber, got.SequenceNumber)
	}

	if expected.TEID != got.TEID {
		return fmt.Errorf("Expected TEID = (%d), got = (%d)", expected.TEID, got.TEID)
	}

	if expected.Priority != got.Priority {
		return fmt.Errorf("Expected Priority = (%d), got = (%d)", expected.Priority, got.Priority)
	}

	if expected.TotalLength != got.TotalLength {
		return fmt.Errorf("Expected TotalLength = (%d), got = (%d)", expected.TotalLength, got.TotalLength)
	}

	if len(expected.InformationElements) != len(got.InformationElements) {
		return fmt.Errorf("Expected (%d) IEs, got = (%d)", len(expected.InformationElements), len(got.InformationElements))
	}

	for ieIndex, expectedIE := range expected.InformationElements {
		if err := compareTwoV2IEObjects(expectedIE, got.InformationElements[ieIndex]); err != nil {
			return fmt.Errorf("For IE (%d): %s", ieIndex, err)
		}
	}

	return nil
}
