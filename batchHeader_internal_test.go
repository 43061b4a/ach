// Copyright 2017 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"strings"
	"testing"
)

// mockBatchheader creates a batch header
func mockBatchHeader() *BatchHeader {
	bh := NewBatchHeader()
	bh.ServiceClassCode = 220
	bh.StandardEntryClassCode = "PPD"
	bh.CompanyName = "ACME Corporation"
	bh.CompanyIdentification = "121042882"
	bh.CompanyEntryDescription = "PAYROLL"
	bh.ODFIIdentification = "12104288"
	return bh
}

// testMockBatchHeader creates a batch header
func testMockBatchHeader(t testing.TB) {
	bh := mockBatchHeader()
	if err := bh.Validate(); err != nil {
		t.Error("mockBatchHeader does not validate and will break other tests")
	}
	if bh.ServiceClassCode != 220 {
		t.Error("ServiceClassCode dependent default value has changed")
	}
	if bh.StandardEntryClassCode != "PPD" {
		t.Error("StandardEntryClassCode dependent default value has changed")
	}
	if bh.CompanyName != "ACME Corporation" {
		t.Error("CompanyName dependent default value has changed")
	}
	if bh.CompanyIdentification != "121042882" {
		t.Error("CompanyIdentification dependent default value has changed")
	}
	if bh.CompanyEntryDescription != "PAYROLL" {
		t.Error("CompanyEntryDescription dependent default value has changed")
	}
	if bh.ODFIIdentification != "12104288" {
		t.Error("ODFIIdentification dependent default value has changed")
	}
}

// testParseBatchHeader parses a known batch header record string
func testParseBatchHeader(t testing.TB) {
	var line = "5225companyname                         origid    PPDCHECKPAYMT000002080730   1076401250000001"
	r := NewReader(strings.NewReader(line))
	r.line = line
	if err := r.parseBatchHeader(); err != nil {
		t.Errorf("%T: %s", err, err)
	}
	record := r.currentBatch.GetHeader()

	if record.recordType != "5" {
		t.Errorf("RecordType Expected '5' got: %v", record.recordType)
	}
	if record.ServiceClassCode != 225 {
		t.Errorf("ServiceClassCode Expected '225' got: %v", record.ServiceClassCode)
	}
	if record.CompanyNameField() != "companyname     " {
		t.Errorf("CompanyName Expected 'companyname    ' got: '%v'", record.CompanyNameField())
	}
	if record.CompanyDiscretionaryDataField() != "                    " {
		t.Errorf("CompanyDiscretionaryData Expected '                    ' got: %v", record.CompanyDiscretionaryDataField())
	}
	if record.CompanyIdentificationField() != "origid    " {
		t.Errorf("CompanyIdentification Expected 'origid    ' got: %v", record.CompanyIdentificationField())
	}
	if record.StandardEntryClassCode != "PPD" {
		t.Errorf("StandardEntryClassCode Expected 'PPD' got: %v", record.StandardEntryClassCode)
	}
	if record.CompanyEntryDescriptionField() != "CHECKPAYMT" {
		t.Errorf("CompanyEntryDescription Expected 'CHECKPAYMT' got: %v", record.CompanyEntryDescriptionField())
	}
	if record.CompanyDescriptiveDate != "000002" {
		t.Errorf("CompanyDescriptiveDate Expected '000002' got: %v", record.CompanyDescriptiveDate)
	}
	if record.EffectiveEntryDateField() != "080730" {
		t.Errorf("EffectiveEntryDate Expected '080730' got: %v", record.EffectiveEntryDateField())
	}
	if record.settlementDate != "   " {
		t.Errorf("SettlementDate Expected '   ' got: %v", record.settlementDate)
	}
	if record.OriginatorStatusCode != 1 {
		t.Errorf("OriginatorStatusCode Expected 1 got: %v", record.OriginatorStatusCode)
	}
	if record.ODFIIdentificationField() != "07640125" {
		t.Errorf("OdfiIdentification Expected '07640125' got: %v", record.ODFIIdentificationField())
	}
	if record.BatchNumberField() != "0000001" {
		t.Errorf("BatchNumber Expected '0000001' got: %v", record.BatchNumberField())
	}
}

// TestParseBatchHeader tests parsing a known batch header record string
func TestParseBatchHeader(t *testing.T) {
	testParseBatchHeader(t)
}

// BenchmarkParseBatchHeader benchmarks parsing a known batch header record string
func BenchmarkParseBatchHeader(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testParseBatchHeader(b)
	}
}

// testBHString validates that a known parsed file can be return to a string of the same value
func testBHString(t testing.TB) {
	var line = "5225companyname                         origid    PPDCHECKPAYMT000002080730   1076401250000001"
	r := NewReader(strings.NewReader(line))
	r.line = line
	if err := r.parseBatchHeader(); err != nil {
		t.Errorf("%T: %s", err, err)
	}
	record := r.currentBatch.GetHeader()

	if record.String() != line {
		t.Errorf("Strings do not match")
	}
}

// TestBHString tests validating that a known parsed file can be return to a string of the same value
func TestBHString(t *testing.T) {
	testBHString(t)
}

// BenchmarkBHString benchmarks validating that a known parsed file can be return to a string of the same value
func BenchmarkBHString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBHString(b)
	}
}

// testValidateBHRecordType verifies error if recordType is not 5
func testValidateBHRecordType(t testing.TB) {
	bh := mockBatchHeader()
	bh.recordType = "2"
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "recordType" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestValidateBHRecordType tests verifying error if recordType is not 5
func TestValidateBHRecordType(t *testing.T) {
	testValidateBHRecordType(t)
}

// BenchmarkValidateBHRecordType benchmarks verifying error if recordType is not 5
func BenchmarkValidateBHRecordType(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testValidateBHRecordType(b)
	}
}

// testInvalidServiceCode verifies error if service code is not valid
func testInvalidServiceCode(t testing.TB) {
	bh := mockBatchHeader()
	bh.ServiceClassCode = 123
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "ServiceClassCode" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestInvalidServiceCode tests verifying error if service code is not valid
func TestInvalidServiceCode(t *testing.T) {
	testInvalidServiceCode(t)
}

// BenchmarkInvalidServiceCode benchmarks verifying error if service code is not valid
func BenchmarkInvalidServiceCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testInvalidServiceCode(b)
	}
}


// testValidateInvalidSECCode verifies error if service class is not valid
func testInvalidSECCode(t testing.TB) {
	bh := mockBatchHeader()
	bh.StandardEntryClassCode = "123"
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "StandardEntryClassCode" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestInvalidSECCode tests verifying error if service class is not valid
func TestInvalidSECCode(t *testing.T) {
	testInvalidSECCode(t)
}

// BenchmarkInvalidSECCode benchmarks verifying error if service class is not valid
func BenchmarkInvalidSECCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testInvalidSECCode(b)
	}
}

// testInvalidOrigStatusCode verifies error if originator status code is not valid
func testInvalidOrigStatusCode(t testing.TB) {
	bh := mockBatchHeader()
	bh.OriginatorStatusCode = 3
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "OriginatorStatusCode" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestInvalidOrigStatusCode tests verifying error if originator status code is not valid
func TestInvalidOrigStatusCode(t *testing.T) {
	testInvalidOrigStatusCode(t)
}

// BenchmarkInvalidOrigStatusCode benchmarks  verifying error if originator status code is not valid
func BenchmarkInvalidOrigStatusCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testInvalidOrigStatusCode(b)
	}
}

// testBatchHeaderFieldInclusion verifies batch header field inclusion
func testBatchHeaderFieldInclusion(t testing.TB) {
	bh := mockBatchHeader()
	bh.BatchNumber = 0
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "BatchNumber" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBatchHeaderFieldInclusion tests verifying batch header field inclusion
func TestBatchHeaderFieldInclusion(t *testing.T) {
	testBatchHeaderFieldInclusion(t)
}

// BenchmarkBatchHeaderFieldInclusion benchmarks verifying batch header field inclusion
func BenchmarkBatchHeaderFieldInclusion(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchHeaderFieldInclusion(b)
	}
}

// testBatchHeaderCompanyNameAlphaNumeric verifies batch header company name is alphanumeric
func testBatchHeaderCompanyNameAlphaNumeric(t testing.TB) {
	bh := mockBatchHeader()
	bh.CompanyName = "AT&T®"
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "CompanyName" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBatchHeaderCompanyNameAlphaNumeric tests verifying batch header company name is alphanumeric
func TestBatchHeaderCompanyNameAlphaNumeric(t *testing.T) {
	testBatchHeaderCompanyNameAlphaNumeric(t)
}

// BenchmarkBatchHeaderCompanyNameAlphaNumeric benchmarks verifying batch header company name is alphanumeric
func BenchmarkBatchHeaderCompanyNameAlphaNumeric(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchHeaderCompanyNameAlphaNumeric(b)
	}
}

// testBatchCompanyDiscretionaryDataAlphaNumeric verifies company discretionary data is alphanumeric
func testBatchCompanyDiscretionaryDataAlphaNumeric(t testing.TB) {
	bh := mockBatchHeader()
	bh.CompanyDiscretionaryData = "®"
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "CompanyDiscretionaryData" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBatchCompanyDiscretionaryDataAlphaNumeric tests verifying company discretionary data is alphanumeric
func TestBatchCompanyDiscretionaryDataAlphaNumeric(t *testing.T) {
	testBatchCompanyDiscretionaryDataAlphaNumeric(t)
}

// BenchmarkBatchCompanyDiscretionaryDataAlphaNumeric benchmarks verifying company discretionary data is alphanumeric
func BenchmarkBatchCompanyDiscretionaryDataAlphaNumeric(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchCompanyDiscretionaryDataAlphaNumeric(b)
	}
}

// testBatchCompanyIdentificationAlphaNumeric verifies company identification is alphanumeric
func testBatchCompanyIdentificationAlphaNumeric(t testing.TB) {
	bh := mockBatchHeader()
	bh.CompanyIdentification = "®"
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "CompanyIdentification" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBatchCompanyIdentificationAlphaNumeric tests verifying company identification is alphanumeric
func TestBatchCompanyIdentificationAlphaNumeric(t *testing.T) {
	testBatchCompanyIdentificationAlphaNumeric(t)
}

// BenchmarkBatchCompanyIdentificationAlphaNumeric benchmarks verifying company identification is alphanumeric
func BenchmarkBatchCompanyIdentificationAlphaNumeric(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchCompanyIdentificationAlphaNumeric(b)
	}
}

// testBatchCompanyEntryDescriptionAlphaNumeric verifies company entry description is alphanumeric
func testBatchCompanyEntryDescriptionAlphaNumeric(t testing.TB) {
	bh := mockBatchHeader()
	bh.CompanyEntryDescription = "P®YROLL"
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "CompanyEntryDescription" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBatchCompanyEntryDescriptionAlphaNumeric tests verifying company entry description is alphanumeric
func TestBatchCompanyEntryDescriptionAlphaNumeric(t *testing.T) {
	testBatchCompanyEntryDescriptionAlphaNumeric(t)
}

// BenchmarkBatchCompanyEntryDescriptionAlphaNumeric benchmarks verifying company entry description is alphanumeric
func BenchmarkBatchCompanyEntryDescriptionAlphaNumeric(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchCompanyEntryDescriptionAlphaNumeric(b)
	}
}

// testBHFieldInclusionRecordType verifies record type field inclusion
func testBHFieldInclusionRecordType(t testing.TB) {
	bh := mockBatchHeader()
	bh.recordType = ""
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "recordType" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBHFieldInclusionRecordType tests verifying record type field inclusion
func TestBHFieldInclusionRecordType(t *testing.T) {
	testBHFieldInclusionRecordType(t)
}

// BenchmarkBHFieldInclusionRecordType benchmarks verifying record type field inclusion
func BenchmarkBHFieldInclusionRecordType(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBHFieldInclusionRecordType(b)
	}
}

// testBHFieldInclusionCompanyName verifies company name field inclusion
func testBHFieldInclusionCompanyName(t testing.TB) {
	bh := mockBatchHeader()
	bh.CompanyName = ""
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "CompanyName" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBHFieldInclusionCompanyName tests verifying company name field inclusion
func TestBHFieldInclusionCompanyName(t *testing.T) {
	testBHFieldInclusionCompanyName(t)
}

// BenchmarkBHFieldInclusionCompanyName benchmarks verifying company name field inclusion
func BenchmarkBHFieldInclusionCompanyName(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBHFieldInclusionCompanyName(b)
	}
}

// testBHFieldInclusionCompanyIdentification verifies company identification field inclusion
func testBHFieldInclusionCompanyIdentification(t testing.TB) {
	bh := mockBatchHeader()
	bh.CompanyIdentification = ""
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "CompanyIdentification" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBHFieldInclusionCompanyIdentification tests verifying company identification field inclusion
func TestBHFieldInclusionCompanyIdentification(t *testing.T) {
	testBHFieldInclusionCompanyIdentification(t)
}

// BenchmarkBHFieldInclusionCompanyIdentification benchmarks verifying company identification field inclusion
func BenchmarkBHFieldInclusionCompanyIdentification(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBHFieldInclusionCompanyIdentification(b)
	}
}

// testBHFieldInclusionStandardEntryClassCode verifies SEC Code field inclusion
func testBHFieldInclusionStandardEntryClassCode(t testing.TB) {
	bh := mockBatchHeader()
	bh.StandardEntryClassCode = ""
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "StandardEntryClassCode" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBHFieldInclusionStandardEntryClassCode tests verifying SEC Code field inclusion
func TestBHFieldInclusionStandardEntryClassCode(t *testing.T) {
	testBHFieldInclusionStandardEntryClassCode(t)
}

// BenchmarkBHFieldInclusionStandardEntryClassCode benchmarks verifying SEC Code field inclusion
func BenchmarkBHFieldInclusionStandardEntryClassCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBHFieldInclusionStandardEntryClassCode(b)
	}
}

// testBHFieldInclusionCompanyEntryDescription verifies Company Entry Description field inclusion
func testBHFieldInclusionCompanyEntryDescription(t testing.TB) {
	bh := mockBatchHeader()
	bh.CompanyEntryDescription = ""
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "CompanyEntryDescription" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBHFieldInclusionCompanyEntryDescription tests verifying Company Entry Description field inclusion
func Test(t *testing.T) {
	testBHFieldInclusionCompanyEntryDescription(t)
}

// BenchmarkBHFieldInclusionCompanyEntryDescription benchmarks verifying Company Entry Description field inclusion
func BenchmarkBHFieldInclusionCompanyEntryDescription(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBHFieldInclusionCompanyEntryDescription(b)
	}
}


// testBHFieldInclusionOriginatorStatusCode verifies Originator Status Code field inclusion
func testBHFieldInclusionOriginatorStatusCode(t testing.TB) {
	bh := mockBatchHeader()
	bh.OriginatorStatusCode = 0
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "OriginatorStatusCode" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBHFieldInclusionOriginatorStatusCode tests verifying Originator Status Code field inclusion
func TestBHFieldInclusionOriginatorStatusCode(t *testing.T) {
	testBHFieldInclusionOriginatorStatusCode(t)
}

// BenchmarkBHFieldInclusionOriginatorStatusCode benchmarks verifying Originator Status Code field inclusion
func BenchmarkBHFieldInclusionOriginatorStatusCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBHFieldInclusionOriginatorStatusCode(b)
	}
}


// testBHFieldInclusionODFIIdentification verifies ODFIIdentification field inclusion
func testBHFieldInclusionODFIIdentification(t testing.TB) {
	bh := mockBatchHeader()
	bh.ODFIIdentification = ""
	if err := bh.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "ODFIIdentification" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestBHFieldInclusionODFIIdentification tests verifying ODFIIdentification field inclusion
func TestBHFieldInclusionODFIIdentification(t *testing.T) {
	testBHFieldInclusionODFIIdentification(t)
}

// BenchmarkBHFieldInclusionODFIIdentification benchmarks verifying ODFIIdentification field inclusion
func BenchmarkBHFieldInclusionODFIIdentification(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBHFieldInclusionODFIIdentification(b)
	}
}