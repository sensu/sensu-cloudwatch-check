package main

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Create mockService Object to use in testing.
// FIXME: replace s3 specific items with correct AWS service items
type mockService struct {
	listBucketsOutput      *s3.ListBucketsOutput
	getBucketTaggingOutput *s3.GetBucketTaggingOutput
}

// Create mockService Functions that match functions defined in ServiceAPI interface in main.go
func (m mockService) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.listBucketsOutput, nil
}

func (m mockService) GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
	return m.getBucketTaggingOutput, nil
}

func TestCheckFunction(t *testing.T) {
	cases := []struct {
		client                 mockService
		tags                   []string
		listBucketsOutput      *s3.ListBucketsOutput
		getBucketTaggingOutput *s3.GetBucketTaggingOutput
		expectedState          int
	}{ //start of array
		{ //start of struct
			client:        mockService{},
			expectedState: 0,
			tags: []string{
				"test_key",
			},
			listBucketsOutput: &s3.ListBucketsOutput{
				Buckets: []types.Bucket{
					types.Bucket{
						Name: func() *string { s := "test_bucket"; return &s }(),
					},
					types.Bucket{
						Name: func() *string { s := "second_bucket"; return &s }(),
					},
				},
			},
			getBucketTaggingOutput: &s3.GetBucketTaggingOutput{
				TagSet: []types.Tag{
					types.Tag{
						Key:   func() *string { s := "test_key"; return &s }(),
						Value: func() *string { s := "test_value"; return &s }(),
					},
				},
			},
		},
	}
	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := tt.client
			client.listBucketsOutput = tt.listBucketsOutput
			client.getBucketTaggingOutput = tt.getBucketTaggingOutput
			state, err := checkFunction(client)
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if state != tt.expectedState {
				t.Errorf("expect state: %v, got %v", tt.expectedState, state)
			}
		})

	}
}
