package meilisearch

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex_AddDocuments(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		documentsPtr interface{}
	}
	tests := []struct {
		name          string
		args          args
		wantResp      *AsyncUpdateID
		wantErr       bool
		expectedError Error
	}{
		{
			name: "TestIndexBasicAddDocuments",
			args: args{
				UID:    "1",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexAddDocumentsWithCustomClient",
			args: args{
				UID:    "2",
				client: customClient,
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexMultipleAddDocuments",
			args: args{
				UID:    "2",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
					{"ID": "456", "Name": "Le Petit Prince"},
					{"ID": "1", "Name": "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexBasicAddDocumentsWithIntID",
			args: args{
				UID:    "3",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"BookID": float64(123), "Title": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexAddDocumentsWithIntIDWithCustomClient",
			args: args{
				UID:    "4",
				client: customClient,
				documentsPtr: []map[string]interface{}{
					{"BookID": float64(123), "Title": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexMultipleAddDocumentsWithIntID",
			args: args{
				UID:    "5",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"BookID": float64(123), "Title": "Pride and Prejudice"},
					{"BookID": float64(456), "Title": "Le Petit Prince", "Tag": "Conte"},
					{"BookID": float64(1), "Title": "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexAddDocumentsMissingPrimaryKey",
			args: args{
				UID:    "6",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"Key": float64(123), "Title": "Pride and Prejudice"},
					{"Key": float64(456), "Title": "Le Petit Prince", "Tag": "Conte"},
					{"Key": float64(1), "Title": "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
			wantErr: true,
			expectedError: Error(Error{
				Endpoint:         "/indexes/6/documents",
				Method:           "POST",
				Function:         "AddDocuments",
				RequestToString:  "[{\"Key\":123,\"Title\":\"Pride and Prejudice\"},{\"Key\":456,\"Tag\":\"Conte\",\"Title\":\"Le Petit Prince\"},{\"Key\":1,\"Title\":\"Alice In Wonderland\"}]",
				ResponseToString: "{\"message\":\"schema cannot be built without a primary key\",\"errorCode\":\"missing_primary_key\",\"errorType\":\"invalid_request_error\",\"errorLink\":\"https://docs.meilisearch.com/errors#missing_primary_key\"}",
				MeilisearchApiMessage: meilisearchApiMessage{
					Message:   "schema cannot be built without a primary key",
					ErrorCode: "missing_primary_key",
					ErrorType: "invalid_request_error",
					ErrorLink: "https://docs.meilisearch.com/errors#missing_primary_key",
				},
				StatusCode:         400,
				StatusCodeExpected: []int{202},
				rawMessage:         "unaccepted status code found: ${statusCode} expected: ${statusCodeExpected}, MeilisearchApiError Message: ${message}, ErrorCode: ${errorCode}, ErrorType: ${errorType}, ErrorLink: ${errorLink} (path \"${method} ${endpoint}\" with method \"${function}\")",
				OriginError:        error(nil),
				ErrCode:            4,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.AddDocuments(tt.args.documentsPtr)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, &tt.expectedError, err)
			} else {
				require.GreaterOrEqual(t, gotResp.UpdateID, tt.wantResp.UpdateID)
				require.NoError(t, err)
				_, err = i.DefaultWaitForPendingUpdate(gotResp)
				require.NoError(t, err)

				var documents []map[string]interface{}
				err = i.GetDocuments(&DocumentsRequest{
					Limit: 3,
				}, &documents)
				require.NoError(t, err)
				require.Equal(t, tt.args.documentsPtr, documents)
			}

			_, _ = deleteAllIndexes(c)
		})
	}
}

func TestIndex_AddDocumentsWithPrimaryKey(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		documentsPtr interface{}
		primaryKey   string
	}
	tests := []struct {
		name     string
		args     args
		wantResp *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicAddDocumentsWithPrimaryKey",
			args: args{
				UID:    "1",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"key": "123", "Name": "Pride and Prejudice"},
				},
				primaryKey: "key",
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexAddDocumentsWithPrimaryKeyWithCustomClient",
			args: args{
				UID:    "2",
				client: customClient,
				documentsPtr: []map[string]interface{}{
					{"key": "123", "Name": "Pride and Prejudice"},
				},
				primaryKey: "key",
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexMultipleAddDocumentsWithPrimaryKey",
			args: args{
				UID:    "3",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"key": "123", "Name": "Pride and Prejudice"},
					{"key": "456", "Name": "Le Petit Prince"},
					{"key": "1", "Name": "Alice In Wonderland"},
				},
				primaryKey: "key",
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexAddDocumentsWithPrimaryKeyWithIntID",
			args: args{
				UID:    "4",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"key": float64(123), "Name": "Pride and Prejudice"},
				},
				primaryKey: "key",
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexMultipleAddDocumentsWithPrimaryKeyWithIntID",
			args: args{
				UID:    "5",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"key": float64(123), "Name": "Pride and Prejudice"},
					{"key": float64(456), "Name": "Le Petit Prince"},
					{"key": float64(1), "Name": "Alice In Wonderland"},
				},
				primaryKey: "key",
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotResp, err := i.AddDocumentsWithPrimaryKey(tt.args.documentsPtr, tt.args.primaryKey)
			require.GreaterOrEqual(t, gotResp.UpdateID, tt.wantResp.UpdateID)
			require.NoError(t, err)
			_, err = i.DefaultWaitForPendingUpdate(gotResp)
			require.NoError(t, err)

			var documents []map[string]interface{}
			err = i.GetDocuments(&DocumentsRequest{
				Limit: 3,
			}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.args.documentsPtr, documents)

			_, _ = deleteAllIndexes(c)
		})
	}
}

func TestIndex_DeleteAllDocuments(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicDeleteAllDocuments",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexDeleteAllDocumentsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)

			SetUpBasicIndex()
			gotResp, err := i.DeleteAllDocuments()
			require.NoError(t, err)
			require.Equal(t, gotResp, tt.wantResp)
			_, err = i.DefaultWaitForPendingUpdate(gotResp)
			require.NoError(t, err)

			var documents interface{}
			err = i.GetDocuments(&DocumentsRequest{
				Limit: 5,
			}, &documents)
			require.NoError(t, err)
			require.Empty(t, documents)

			_, _ = deleteAllIndexes(c)
		})
	}
}

func TestIndex_DeleteOneDocument(t *testing.T) {
	type args struct {
		UID          string
		PrimaryKey   string
		client       *Client
		identifier   string
		documentsPtr interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantResp *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicDeleteOneDocument",
			args: args{
				UID:        "1",
				client:     defaultClient,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexDeleteOneDocumentWithCustomClient",
			args: args{
				UID:        "2",
				client:     customClient,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexDeleteOneDocumentinMultiple",
			args: args{
				UID:        "3",
				client:     defaultClient,
				identifier: "456",
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
					{"ID": "456", "Name": "Le Petit Prince"},
					{"ID": "1", "Name": "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexBasicDeleteOneDocumentWithIntID",
			args: args{
				UID:        "4",
				client:     defaultClient,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"BookID": 123, "Title": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexDeleteOneDocumentWithIntIDWithCustomClient",
			args: args{
				UID:        "5",
				client:     customClient,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"BookID": 123, "Title": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexDeleteOneDocumentWithIntIDinMultiple",
			args: args{
				UID:        "6",
				client:     defaultClient,
				identifier: "456",
				documentsPtr: []map[string]interface{}{
					{"BookID": 123, "Title": "Pride and Prejudice"},
					{"BookID": 456, "Title": "Le Petit Prince"},
					{"BookID": 1, "Title": "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotAddResp, err := i.AddDocuments(tt.args.documentsPtr)
			require.GreaterOrEqual(t, gotAddResp.UpdateID, tt.wantResp.UpdateID)
			require.NoError(t, err)
			_, err = i.DefaultWaitForPendingUpdate(gotAddResp)
			require.NoError(t, err)

			gotResp, err := i.DeleteDocument(tt.args.identifier)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.UpdateID, tt.wantResp.UpdateID)
			_, err = i.DefaultWaitForPendingUpdate(gotResp)
			require.NoError(t, err)

			var document []map[string]interface{}
			err = i.GetDocument(tt.args.identifier, &document)
			require.NoError(t, err)
			require.Empty(t, document)

			_, _ = deleteAllIndexes(c)
		})
	}
}

func TestIndex_DeleteDocuments(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		identifier   []string
		documentsPtr []docTest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicDeleteDocument",
			args: args{
				UID:        "1",
				client:     defaultClient,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexDeleteDocumentWithCustomClient",
			args: args{
				UID:        "2",
				client:     customClient,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexBasicDeleteDocument",
			args: args{
				UID:        "3",
				client:     defaultClient,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexBasicDeleteDocument",
			args: args{
				UID:        "4",
				client:     defaultClient,
				identifier: []string{"123", "456", "1"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)

			gotAddResp, err := i.AddDocuments(tt.args.documentsPtr)
			require.NoError(t, err)
			_, err = i.DefaultWaitForPendingUpdate(gotAddResp)
			require.NoError(t, err)

			gotResp, err := i.DeleteDocuments(tt.args.identifier)
			require.NoError(t, err)
			require.Equal(t, gotResp, tt.wantResp)
			_, err = i.DefaultWaitForPendingUpdate(gotResp)
			require.NoError(t, err)

			var document docTest
			for _, identifier := range tt.args.identifier {
				err = i.GetDocument(identifier, &document)
				require.NoError(t, err)
				require.Empty(t, document)
			}

			_, _ = deleteAllIndexes(c)
		})
	}
}

func TestIndex_GetDocument(t *testing.T) {
	type args struct {
		UID         string
		client      *Client
		identifier  string
		documentPtr *docTestBooks
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestIndexBasicGetDocument",
			args: args{
				UID:         "indexUID",
				client:      defaultClient,
				identifier:  "123",
				documentPtr: &docTestBooks{},
			},
			wantErr: false,
		},
		{
			name: "TestIndexGetDocumentWithCustomClient",
			args: args{
				UID:         "indexUID",
				client:      customClient,
				identifier:  "123",
				documentPtr: &docTestBooks{},
			},
			wantErr: false,
		},
		{
			name: "TestIndexGetDocumentWithNoExistingDocument",
			args: args{
				UID:         "indexUID",
				client:      defaultClient,
				identifier:  "125",
				documentPtr: &docTestBooks{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			SetUpBasicIndex()

			require.Empty(t, tt.args.documentPtr)
			err := i.GetDocument(tt.args.identifier, tt.args.documentPtr)
			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, tt.args.documentPtr)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, tt.args.documentPtr)
				require.Equal(t, strconv.Itoa(tt.args.documentPtr.BookID), tt.args.identifier)
			}

			_, _ = deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateDocuments(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		documentsPtr []docTestBooks
	}
	tests := []struct {
		name string
		args args
		want *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicUpdateDocument",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentOnMultipleDocuments",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
					{BookID: 1344, Title: "Harry Potter and the Half-Blood Prince"},
					{BookID: 4, Title: "The Hobbit"},
					{BookID: 42, Title: "The Great Gatsby"},
				},
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentWithNoExistingDocument",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 237, Title: "One Hundred Years of Solitude"},
				},
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentWithNoExistingMultipleDocuments",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 246, Title: "One Hundred Years of Solitude"},
					{BookID: 834, Title: "To Kill a Mockingbird"},
					{BookID: 44, Title: "Don Quixote"},
					{BookID: 594, Title: "The Great Gatsby"},
				},
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			SetUpBasicIndex()

			got, err := i.UpdateDocuments(tt.args.documentsPtr)
			require.NoError(t, err)
			require.Equal(t, got, tt.want)
			_, err = i.DefaultWaitForPendingUpdate(got)
			require.NoError(t, err)

			var document docTestBooks
			for _, identifier := range tt.args.documentsPtr {
				err = i.GetDocument(strconv.Itoa(identifier.BookID), &document)
				require.NoError(t, err)
				require.Equal(t, identifier.BookID, document.BookID)
				require.Equal(t, identifier.Title, document.Title)
			}

			_, _ = deleteAllIndexes(c)
		})
	}
}

func TestIndex_UpdateDocumentsWithPrimaryKey(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		documentsPtr []docTestBooks
		primaryKey   string
	}
	tests := []struct {
		name string
		args args
		want *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicUpdateDocumentsWithPrimaryKey",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
				primaryKey: "book_id",
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
				primaryKey: "book_id",
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyOnMultipleDocuments",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
					{BookID: 1344, Title: "Harry Potter and the Half-Blood Prince"},
					{BookID: 4, Title: "The Hobbit"},
					{BookID: 42, Title: "The Great Gatsby"},
				},
				primaryKey: "book_id",
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyWithNoExistingDocument",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 237, Title: "One Hundred Years of Solitude"},
				},
				primaryKey: "book_id",
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyWithNoExistingMultipleDocuments",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 246, Title: "One Hundred Years of Solitude"},
					{BookID: 834, Title: "To Kill a Mockingbird"},
					{BookID: 44, Title: "Don Quixote"},
					{BookID: 594, Title: "The Great Gatsby"},
				},
				primaryKey: "book_id",
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			SetUpBasicIndex()

			got, err := i.UpdateDocumentsWithPrimaryKey(tt.args.documentsPtr, tt.args.primaryKey)
			require.NoError(t, err)
			require.Equal(t, got, tt.want)
			_, err = i.DefaultWaitForPendingUpdate(got)
			require.NoError(t, err)

			var document docTestBooks
			for _, identifier := range tt.args.documentsPtr {
				err = i.GetDocument(strconv.Itoa(identifier.BookID), &document)
				require.NoError(t, err)
				require.Equal(t, identifier.BookID, document.BookID)
				require.Equal(t, identifier.Title, document.Title)
			}

			_, _ = deleteAllIndexes(c)
		})
	}
}
